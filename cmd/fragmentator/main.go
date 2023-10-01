package main

import (
	"fragmentator/internal/composer"
	"github.com/charmbracelet/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	app, err := composer.NewComposer()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Run application")
	app.RunApplication()

	log.Info("Press (Ctrl+C) to finish")
	<-sigs

	log.Info("Shutdown application")
	app.StopApplication()
	log.Info("Exit")
}
