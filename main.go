package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("hello world")
	sigStop := make(chan os.Signal, 1)
	signal.Notify(sigStop, syscall.SIGINT, syscall.SIGTERM)
	<-sigStop
}
