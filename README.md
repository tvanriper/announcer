# Announcer

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/tvanriper/announcer)
[![Coverage Status](https://coveralls.io/repos/github/tvanriper/announcer/badge.svg?branch=main)](https://coveralls.io/github/tvanriper/announcer?branch=main)

Provides a write-once-read-many channel of communications in pure Golang.

## Usage

As an overview:

* Create a new announcer.
* Let different components listen to the announcer.
* Send information through the announcer to the listeners.

This looks like (without proper error checking):

```golang
sender := announcer.New(5)
listenA := sender.Listen()
listenB := sender.Listen()
sender.Send("hullo, world")
if hulloA, ok := listenA.Listen(); ok {
    fmt.Println(hulloA.(string))
}
if hulloB, ok := listenB.Listen(); ok {
    fmt.Println(hulloB.(string))
}
sender.Close()
```

This is obviously much more interesting when using Golang threads.

## Installation

`go get github.com/tvanriper/announcer`
