package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
	"github.com/ziad-eliwa/jit-version-control-system/internal/route"
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
	defer app.DB.Close()

	r := route.SetupRoutes(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      r,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	app.Logger.Info("Server is running on port", "port", port)
	err = server.ListenAndServe()
	log.Fatalf("Listen and Serve: %v", err)
}
