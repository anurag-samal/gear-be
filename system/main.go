package main

import (
	"github/anurag/altar-be/system/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatal("Error in server configuration: ", err)
	}

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error running server: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	srv.Shutdown()

}
