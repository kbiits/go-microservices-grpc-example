# Microservices in Golang using gRPC

## Pre-Install
You must have mongodb installed in your machine and already configured to allow connections from docker container.  
To do that, modify file `/etc/mongodb.conf` and add docker `bridge ip` in `net` section of the configuration


```conf
# bind ip to 172.17.0.1 to allow connections from docker container
# network interfaces
net:
  port: 27017
  bindIp: 127.0.0.1,172.17.0.1
```

## Installation
Clone this repo
```bash
git clone https://github.com/kbiits/go-microservices-grpc-example
```

Create user for mongodb
```bash
mongo user/mongo-init.js
```

Create docker network
```bash
docker network create microservices
```

Build golang microservices app  
note: I don't add build steps to Dockerfile because I want to reduce the size of created images

Build advice service :
```bash
cd advice
make build
cd ..
``` 
Build user service :
```bash
cd user
make build
cd ..
``` 

Build gateway :
```bash
cd gateway
make build
cd ..
```

Run docker
```bash
docker compose up -d
```

And taraa, you can see the website running on `localhost:3000` 