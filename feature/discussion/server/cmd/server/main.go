package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi"
)

func main() {
	stop := httpapi.Start(httpapi.Settings{
		Port:  8080,
		UserH: nil,
		AuthH: nil,
		JwtH:  nil,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := stop(ctx)
	if err != nil {
		log.Fatalf("stopping server: %s", err)
	}
}
