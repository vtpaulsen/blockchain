# Dev Environment
Please use the Visual Studio Code workspace `ws.code-workspace`

Otherwise gopls will get confused over the two golang modules. See [Golang Issue #32394](https://github.com/golang/go/issues/32394)

Requirements:

- Golang 1.17+
- Docker (With Docker Compose)
- Make

# Making it work with your project
In order to make it work with your project, your project has to satisfy the requirements of Exercise 4.6, Forming a Network, number 2 and 3.

You'll have to match 2 constants in setupbooter/docker/consts.go with your project:

*READY_MESSAGE*: "What's the IP address?" is what your implementation prints, right before it listens for user input.
When it recieves an IP from the user, it should then write something similiar to "What's the port number?". The second string does not really matter, as what the testing framework is going to do, is just write "{IP}\n{PORT}\n" to your program.

*LISTENER_IP_PRINT_MESSAGE*: "Starting on IP: " is what your implementation prints, right before it prints the IP and port of the port it listens to. The testing framwork just trims out this message, and gets the IP and port.

To recap; your implementation should print something like:

What's the IP address?

> 127.0.0.1\n

What's the port number?

> 42980\n

Starting on IP: 127.0.1.1:33935


# Testing
Testing requires a Docker daemon running

If you make changes to production files (in src), a docker-compose build is necessary.

Notice that the make targets automatically builds for you.

# Manually booting up the docker-compose service

## With Visual Studio Code

#### Start an automatically configured graph/setup of peers
`ctrl+shift+p` &#8594; `Tasks: Run Task` &#8594; `Boot Setup X`

#### To inspect a peer, attach from shell to a docker container *n* with
`ctrl+shift+p` &#8594; `Tasks: Run Task` &#8594; `Inspect container n`

## With a shell

##### Start *n* peer services with
`docker-compose --project-name dissy up --scale peer=n`

##### To inspect a peer, attach from shell to a docker container *n* with
`docker attach dissy_peer_n`

# Cleaning up
To stop running containers, run:

`make Cleanup`

As each consequent build will produce a dangling image, your disk can fill up rather quickly. To removing dangling images, run:  

`docker image prune`
