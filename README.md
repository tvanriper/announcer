# Announcer

Provides a write-once-read-many channel of communications in pure Golang.

## Usage

As an overview:

* Create a new announcer.
* Let different components listen to the announcer.
* Send information through the announcer to the listeners.

This looks like:

```golang
sender := announcer.New(5)
listenA := sender.Listen()
listenB := sender.Listen()
sender.Send("hullo, world")
if hulloA, ok := listenA.Listen(); ok {
    fmt.Println(hulloA)
}
if hulloB, ok := listenB.Listen(); ok {
    fmt.Println(hulloB)
}
sender.Close()
```

This is obviously much more interesting when using Golang threads.

## Installation

`go get github.com/tvanriper/announcer`
