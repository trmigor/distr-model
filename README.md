# distr-model

[![Build Status](https://travis-ci.com/trmigor/distr-model.svg?branch=master)](https://travis-ci.com/trmigor/distr-model)
[![codecov](https://codecov.io/gh/trmigor/distr-model/branch/master/graph/badge.svg?token=IUPSSTH99O)](https://codecov.io/gh/trmigor/distr-model)

Simple model of a distributed system.

## Description

**Goal**: to provide users with a small infrastructure that allows them to create models of distributed processes.

**Requirements**:
* Maximum ease of entry (you need to write a minimum of code and correct/add a minimum number of files).
* Modeling of distributed processes that exchange messages.
* Simulation of synchronous operation mode (processes that received messages during the clock cycle synchronizations send messages to other processes, which receive them at the beginning of the next clock cycle).
* Simulation of an asynchronous operation (each message is delivered to the process within the time specified by the link weight between processes).
* Simulation of message loss. The `errorRate` parameter (`0 <= errorRate <= 1`) determines the probability of message loss. 
* Robust algorithms should be relatively resistant to message loss.

**Types of algorithms for modeling**:
* Topological algorithms;
  * Building spanning trees;
  * Finding the reachability of a node;
  * Finding the shortest path;
* Election algorithms;
* Synchronization algorithms;
* Completion detection algorithms;
* Algorithms for ordered distribution;
* ...

## Project layout

Project layout is corresponding to the [standard layout](https://github.com/golang-standards/project-layout) for Go projects.

* [cmd](cmd) directory contains `main` package, which is compiled into the resulting executable and can be modified by users;
* [configs](configs) directory contains configuration files that can be also modified by user;
* [internal](internal) directory contains packages with internal application logic:
  * [errors](internal/errors) package contains error codes for clarification of arisen errors;
  * [messages](internal/messages) package contains implementation of types related to message passing:
    * [MessageArg.go](internal/messages/MessageArg.go) contains implementation of message argument type;
    * [Message.go](internal/messages/Message.go) contains implementation of message type;
    * [MessageQueue.go](internal/messages/MessageQueue.go) contains implementation of message queue type;
  * [network](internal/network) package contains implementation of the network communication model;
  * [process](internal/process) package contains implementation of the distibuted process model;
  * [world](internal/world) package contains implementation of distributed environment model;
* [pkg](pkg) directory contains export-free packages implementing special data structures, used in the project:
  * [priorityq](pkg/priorityq) package contains implementation of the priority queue data structure;
  * [set](pkg/set) package contains implementation of the set data structure;
* [test](test) directory contains additional testing supplies;
* [user](user) directory contains auxilliary packages that can be modified by users:
  * [context](user/context) package contains contextes for working functions of the distributed processes;
* [vendor](vendor) directory contains dependencies;
* [.travis.yml](.travis.yml) is a Travis CI configuration file;
* [Gopkg.lock](Gopkg.lock) and [Gopkg.toml](Gopkg.toml) are dependency configuration files;
* [Makefile](Makefile) is a script file for GNU Make.

## How it is implemented:

The entire model is implemented in Go.

There are several simple types for implementing common functions. They are located in the different packages in [internal](internal) directory. The purpose of the lab work is to write a suitable [`main`](cmd/main.go) and a [message handler function](cmd/main.go) (more on this later). An example is given in the project.

The main type of the project is [`World`](internal/world/World.go). It creates models of distributed processes (hereinafter - just processes), registers handler functions (common to all system processes) and assigns a handler for a specific process.

The [`Network`](internal/network/NetworkLayer.go) type is responsible for inter-process communication and inter-process message delivery. There can be multiple networks in one world and each process can belong to multiple networks. It simulates asynchronous and synchronous modes of sending messages between processes and can also introduce errors in transmission (for example, a message may be lost with some probability).

The [`Process`](/internal/process/Process.go) type models the distributed process itself. Each process must register on its network in order to notify the network of its appearance. The network now knows where to send messages intended for this process. There are one incoming message queue.  One execution thread is started - `WorkerThread`.  The workflow analyzes the message which handler function this message corresponds to and calls the corresponding handler.

The names of the other types speak for themselves: [`ErrorCode`](internal/errors/Errors.go), [`MessageArg`](internal/messages/MessageArg.go), [`Message`](internal/messages/Message.go), [`MessageQueue`](internal/messages/MessageQueue.go).

A little more about the working function.

It is called with two arguments. The first is the context of the `Process` type, which makes it possible to determine the network topology (immediate neighbors) and its number. The user can add their own context to the `Process` class, which the working function will use. This requires:
1. Describe your context type and put this description in a [`context`](user/context/Context.go) package.
2. Add an instance of this type to the map variable [`Contexts`](user/context/Context.go).
3. Use it as showed in [`workFunctionSETX`](cmd/main.go) function in [`main`](cmd/main.go) package.

This context will be included in the general context of the `Process` class and can be used both by the working function itself and by any other functions (this allows, for example, in one working function to define a list of all available processes, not just neighbors, and in another working function, send messages to these processes).

A worker function should check the message, return `true` if it is ready and can process the message, and `false` if it cannot process it (for example, if the message is intended for another worker function).

The network topology is described in the [config.data](configs/config.data) file. Its commands:

```
; create processes from 1 to 11 
processes 1 11

bidirected 1

errorRate 0.5

link from 1 to 2 [latency 10]

link from 1 to all [latency 5]

link from all to 3 [latency 2]

link from all to all [latency 1]

setprocesses 2 5 TEST

send from 4 to 10 TEST_BEGIN 1

send from -1 to 1 TEST_BEGIN

launch timer 3

wait 10
```

For example, there is a ready-made working function [`workFunctionSETX`](cmd/main.go).

## Requirements

* [**Go**](https://golang.org/) of version 1.16.x;
* [**GNU make**](https://www.gnu.org/software/make/) of version 3.81 and above;
* [**dep**](https://github.com/golang/dep) of version 0.5.4 and above;
* [**golangci-lint**](https://github.com/golangci/golangci-lint) of version 1.39.0 and above.

Also, if you want to contribute, you should use formatting tools:

* [**gofmt**](https://golang.org/pkg/cmd/gofmt/);
* [**goimports**](https://pkg.go.dev/golang.org/x/tools/cmd/goimports).

## Tests

To perform tests, use

```
make test
```
The coverage report will be placed in the [coverage](coverage) directory. To see coverage results use

```
go tool cover -func=coverage/count.out
```

## Build

To build the model use

```
make
```

The resulting executable will be placed in the [bin](bin) directory.

## Usage

To execute, use

```
bin/model
```

## Dependencies

All dependencies are managed by [**dep**](https://github.com/golang/dep). Here are all of them:

* [mt19937](https://github.com/seehuhn/mt19937) - an implementation of Takuji Nishimura's and Makoto Matsumoto's [Mersenne Twister](https://en.wikipedia.org/wiki/Mersenne_twister) pseudo random number generator in Go.
