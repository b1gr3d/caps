package main

import (
	"caps/internal/caps"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	go caps.Run()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Hasta la vista, baby")

}