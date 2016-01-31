# Example TCP server

This experiment was born out of a question: what does it take to deploy a modern TCP server in a real production environment or at least something as close to it as possible?

Docker, Heroku, godeps, etc — all wonderful buzzwords, that smell of artisan San Fransisco salads and overblown tech salaries, but how do you actually configure and use them?

There's only one way to find out.

(Everything is work in progress and is probably outdated by the time anyone reads it).

## The server

The server itself is as primitive as it gets: TCP echo. When you run it locally, just use `nc` to connect to it. Here's how an example session looks on the "server" side:

```bash
$ go run go-tcp-echo.go
2016/01/31 01:26:53 Listening to connections on port 3333
2016/01/31 01:27:00 Accepted new connection.
2016/01/31 01:27:02 Read new data from connection [104 105 10]
2016/01/31 01:27:03 Closed connection.
```

And here's how it looks on the "client":

```bash
nc localhost 3333
hi
hi
^C
```

(Just to clarify, the second "hi" is an echo from the server.)

## Docker

Docker is awesome. Or so they say. It's so easy to deploy, you don't ship just your app, you ship the whole VM image! Sounds interesting. Let's try it.

### Building Docker images

OK, so to deploy this project somewhere with Docker I found out that I need to create a Docker image out of it. After installing Docker on my machine, I can run `docker build` to create an image:

```bash
$ docker build -t go-tcp-echo .
Cannot connect to the Docker daemon. Is the docker daemon running on this host?
```

And, of course, fail.

(Forgot to tell you: this is not a tutorial. More, like, a diary. With a lot of frustration and angry ramblings.)

`docker` isn't actually Docker; it's just the CLI that tries to communicate with background daemon. To actually start the daemon, you have to type this [`docker daemon`](https://docs.docker.com/engine/reference/commandline/daemon/):

```bash
$ docker daemon
docker: 'daemon' is not a docker command.
```

Wait. What the fuck?

So, to be honest, I don't understand why this doesn't work, exactly. But hey! I installed the pretty Docker Toolbox, and it gave me two cool shortucts on my launchpad. And when I launch first one, Docker Quickstart Terminal, it comes with all path variables pre-configured, so when I launch a new terminal using it, my `docker build` command actually works:

```bash
docker build -t go-tcp-echo .
Sending build context to Docker daemon 128.5 kB
Step 1 : FROM golang:onbuild
# Executing 3 build triggers...
Step 1 : COPY . /go/src/app
Step 1 : RUN go-wrapper download
 ---> Running in a23a192142a6
+ exec go get -v -d
Step 1 : RUN go-wrapper install
 ---> Running in 1c263c4b7d0a
+ exec go install -v
app
 ---> 7845878d6f49
Removing intermediate container 3c0092cad181
Removing intermediate container a23a192142a6
Removing intermediate container 1c263c4b7d0a
Step 2 : EXPOSE 3333
 ---> Running in 3919284e2075
 ---> 8da10e7896e5
Removing intermediate container 3919284e2075
Successfully built 8da10e7896e5
```

Now, let's just launch it and repeat the same experiment! Here's the command I use to launch it, give it a name ("test" seems appropriate) and listen on the port number 3333 (same one inside and outside). The last option says it to destroy the image after it's done.

```bash
docker run --publish 3333:3333 --name test --rm go-tcp-echo
+ exec app
2016/01/30 23:54:25 Listening to connections on port 3333
```

So, my app is actually launched and is writing some logs again, yay!

Let's try to connect from the other terminal:

```bash
nc -v localhost 3333
nc: connectx to localhost port 3333 (tcp) failed: Connection refused
nc: connectx to localhost port 3333 (tcp) failed: Connection refused
```

Wait, why?

Oh, right, something's to do with IP stuff. As you may have noticed, I don't have a lot of experience with this networking stuff. But hey — when I launched the "Quickstart Terminal", it did say something interesting:

```bash

                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/


docker is configured to use the default machine with IP 192.168.99.100
```

So, let's may be try using my IP in the local network instead, shall we?

```bash
nc 192.168.99.100 3333
hi
hi
^C
```

Cool. Meanwhile, on the server:

```bash
2016/01/30 23:59:24 Accepted new connection.
2016/01/30 23:59:27 Read new data from connection [104 105 10]
2016/01/30 23:59:28 Closed connection.
```

So, everything's going according to plan: I was able to build a Docker image, to launch it, and to connect to it, on a local machine. So far, so good.

### Automatic Build

Building stuff manually is boring. What if I wanted to automate it instead?

Surely, you wouldn't think that someone would develop such a great piece of infrastructure as Docker completely open-source without making sure that they have an essential SAAS on their hands that is mentioned in all Docker tutorials and resources, would you? So, there's a little cute thing called Docker Hub that is created just for that: automating builds of Docker images. They have repositories, just like on Github; and a special kind of repository, that is named Automated Build, which connects to your actual git repo (like this one) and builds a new Docker image each time your git repo updates.

The only thing that irritates me about this is the fact that they decided to use the word "Build" to call not a process, but an object.

Thankfully, just like on Github, you can create an unlimited amount of Automated Builds as long as they're public. So, you can enjoy a public Automated Build for this Github repo here:

[[https://hub.docker.com/r/golergka/go-tcp-echo/]]
