# A small Azure and AKS challenge to facilitate the discussion

Dear Candidate,

We want to have an open exchange with you to assess your level of
resourcefulness and flexibility when it comes to DevOps and Cloud, and
in order to facilitate that discussion we have thought of a small
challenge for you...

The challenge will be to deploy a K8s cluster from scratch in Azure,
using the AKS service and to adapt the deployment automation baseline to
address some specific problem statements.

The files used are all in this repo which you
can clone. Since we're giving you this document in advance it is fine if
you get it done first and just comment it to us. If you prefer, we will
give you time to do it and check back with you later.

This exercise is not about replicating the steps but about discussing
the process, showing understanding and discussing possible ways of
improving it.

In order to deploy the K8s cluster, you will be given an Azure account

> Username: (to be provided by email)
>
> Password: (to be provided by email)

And a service principal:

> Client ID: (to be provided by email)
>
> Secret: (to be provided by email)

It is suggested to use cloud shell, and, once in cloud shell, automation to deploy a cluster. But you can also other approaches as long as they are representative of a DevOps approach. 

We're proposing the following ansible playbook:

``` yaml
# Description
# ===========
# This playbook first creates an Azure Kubernetes Service cluster of 2 nodes with provided username and ssh key and scales it to 3 nodes.
# Change variables below to customize your AKS deployment.
# It also requires a valid service principal to create AKS cluster, so fill:
# - client_id
# - client_secret
# This sample requires Ansible 2.6 

- name: Create Azure Kubernetes Service
  hosts: localhost
  connection: local
  vars:
    ssh_key: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDH14vPP1NE+qc3SgpoSPYMVaPLj8ouA1azujnMPbVKtN8mSg0KBWVWebx+GLqGI9nrKwmVwPHVzM8LhbbpYhfnN0UA2OV7sxBIshy7BshDTewiLyR9FXuNjuf8Cwd7WI+vgkKO0Yh6EdgwZ2hz+yVTn+5umCRDOg0ht0i2kyPMn1MHv61aN3noYp5zV8W3Tg8oJ2DPei6FxcNheRWkyvjan73/rISQJYWURlI90b3MCQOTXuI2w6BIxF3kHIHLUtzD7FqQhZdEkM/8cA9wUDgnS6+3ySkGOuQazkEWYn60Ll/597u2fLsq6P/pOCM1+KWLf2kY6mTl/dNxiIdV1d01
    resource_group: "{{ resource_group_name }}"
    location: eastus
    aks_name: myAKSCluster
    username: azureuser
    client_id: "{{ lookup('env', 'AZURE_CLIENT_ID') }}"
    client_secret: "{{ lookup('env', 'AZURE_CLIENT_SECRET') }}"
  tasks:
  - name: Create resource group
    azure_rm_resourcegroup:
      name: "{{ resource_group }}"
      location: "{{ location }}"
  - name: Create a managed Azure Container Services (AKS) cluster
    azure_rm_aks:
      name: "{{ aks_name }}"
      location: "{{ location }}"
      resource_group: "{{ resource_group }}"
      dns_prefix: "{{ aks_name }}"
      linux_profile:
        admin_username: "{{ username }}"
        ssh_key: "{{ ssh_key }}"
      service_principal:
        client_id: "{{ client_id }}"
        client_secret: "{{ client_secret }}"
      agent_pool_profiles:
        - name: default
          count: 2
          vm_size: Standard_D2_v2
      tags:
        Environment: Production
  - name: Scale created Azure Container Services (AKS) cluster
    azure_rm_aks:
      name: "{{ aks_name }}"
      location: "{{ location }}"
      resource_group: "{{ resource_group }}"
      dns_prefix: "{{ aks_name }}"
      linux_profile:
        admin_username: "{{ username }}"
        ssh_key: "{{ ssh_key }}"
      service_principal:
        client_id: "{{ client_id }}"
        client_secret: "{{ client_secret }}"
      agent_pool_profiles:
        - name: default
          count: 3
          vm_size: Standard_D2_v2
```



## The Challenge:

In this challenge we want you use the baseline given above as a starting
point.

We then want you to investigate and propose appropriate means for
adapting the cluster deployment, using automation and abstraction of
variable elements where possible, in order to solve the following set of
"problem statements". Each of these problems can be addressed in various
ways, however using automation and abstraction principles are valued
higher than (for example) "hard-coded" solutions.

You are not limited to using specific solution areas or tools -- you are
free to propose what you assess as most effective and flexible.
Obviously, we would like you to present and explain your choices and
achievements after the challenge.

## Problem Statements:

We don't want this cluster to be created as it is in the script by
default.

-   Our applications are memory intensive. Sometimes pod deployments are
    stuck due to no available nodes. We are also facing intermittent
    eviction of some of our pods. We would want to prevent this from
    happening, especially to pods providing critical services.

-   We are in Singapore, and we want to be sure to cater for sovereignity issues.

-   We are not sure what the RSA key is for. 

-   We would also like to choose proper names for the cluster and the
    group of resources.

-   We are not happy with a cluster that has a fixed number of nodes: we
    would like the cluster to have a minimum of three nodes, but be able
    to grow.

-   We would like to build more resilience to our cluster and we're
    thinking of availability zones.

What are your recommendations? We would like some quick and easy
solutions for a start. Include one or more of such solutions into your
answer.

## Getting Started:

In order to get the file mentioned, you can use:

`git clone ...`

You will be able to execute it as:

`ansible-playbook myplaybook.yml`

But first please remember to change it to our needs and later on to get
credentials

`az aks get-credentials --resource-group yourResourceGroup --name yourAKSname`

And you should be able to see the nodes created. We leave this to you.

Once we have a Kubernetes cluster we need to give it some content. An
example would be a service that would return some useful information. In
this case we're giving you a sample app in Go that does a few things:


``` go
// HTTP server, coded in Go, that does three different things, all of them quite easy to guess...

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

```
You should create a Dockerfile that gets a container made from this
code, here you have a sample that works:

``` Dockerfile
FROM golang:latest
WORKDIR /app
COPY ./ /app
RUN go run thalestest.go
```

Time to build the container:

`docker build . -t thalestest`

Which will start executing instantaneously from cloud shell. You will be
able to interact with it with the proper controls. Image should be
tagged as thalestest, make sure it has been done.

Now create a deployment, an example of yaml follows. But take into account that you need to deploy the container we have created, not the sample given below.

``` yaml
apiVersion: apps/v1 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: frontend
  labels:
    app: guestbook
spec:
  selector:
    matchLabels:
      app: guestbook
      tier: frontend
  replicas: 3
  template:
    metadata:
      labels:
        app: guestbook
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: gcr.io/google-samples/gb-frontend:v4
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        env:
        - name: GET_HOSTS_FROM
          value: dns
          # Using `GET_HOSTS_FROM=dns` requires your cluster to
          # provide a dns service. As of Kubernetes 1.3, DNS is a built-in
          # service launched automatically. However, if the cluster you are using
          # does not have a built-in DNS service, you can instead
          # access an environment variable to find the master
          # service's host. To do so, comment out the 'value: dns' line above, and
          # uncomment the line below:
          # value: env
        ports:
        - containerPort: 80
```
`kubectl apply -f ./my-manifest.yaml`

Now the deployment should be made accessible to the internet. The following sample exposes the deployment as a nodeport, but we'd rather use a load balancer. Comment on the available options and whether you will have to create a load balancer manually. If it's an entity creating it (recommended), comment on the identity of the entity.

``` yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # comment or delete the following line if you want to use a LoadBalancer
  type: NodePort 
  # if your cluster supports it, uncomment the following to automatically create
  # an external load-balanced IP for the frontend service.
  # type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: guestbook
    tier: frontend
```

Now the service should be accessible through the internet. We should be able to see the IP inside the pods using the given go code. It is interesting to comment about the IP address we are getting, whether it is always the same and why.

### A few hints you may need:

If you have a problem: check that the ports are right. If you open an
interactive session with the container in the node it will tell you to
which port is the http server attached. You must properly map this to
the externally exposed port.

If you want to do something different to make this challenge work, go ahead, just explain us why.

You can also fork this repo, make all the fine tuning, and show us your own version of the exercise.

For uploading the container in the kubernetes cluster you can choose to store it in dockerhub or any other way you see fit. And remember the entry point.

Singapore is in South East Asia, also known as "southeastasia"