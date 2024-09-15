package main

import (
	"os"
	"os/signal"
	"rest-arch-training/internal/app"
	"syscall"
)

func main() {
	application := app.New()

	application.Config()

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}
