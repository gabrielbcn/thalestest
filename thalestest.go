// Sample application coded in go...
// HTTP server that serves three different things, all of them quite easy to guess...

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const defaultAddr = ":8080"

// main starts an http server on the $PORT environment variable.
func main() {
	addr := defaultAddr
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("server starting to listen on %s", addr)
	log.Printf("http://localhost%s", addr)
	log.Printf("http://localhost%s/test", addr)
	log.Printf("http://localhost%s/ip", addr)
	http.HandleFunc("/", home)
	http.HandleFunc("/ip/", getip)
	http.HandleFunc("/test/", test)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}

// home logs the received request and returns a simple response.
func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request main: %s %s", r.Method, r.URL.Path)
	rand.Seed(time.Now().UnixNano())
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	fmt.Fprintf(w, "Magic 8-Ball says:", answers[rand.Intn(len(answers))])
}

func getip(w http.ResponseWriter, r *http.Request) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Oops: " + err.Error() + "\n")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Fprintf(w, ipnet.IP.String()+"\n")
			}
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request test: %s %s", r.Method, r.URL.Path)
	fmt.Fprintf(w, "The test page")
}
