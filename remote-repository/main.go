package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "GO Backend Server Port")
	flag.Parse()

	logger := slog.Default()

	app, err := app.NewApplication(logger)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		WriteTimeout: time.Minute,
	}

	app.Logger.Info("Server is running on port", "port", port)
	err = server.ListenAndServe()
	log.Fatalf("Listen and Serve: %v", err)
}
