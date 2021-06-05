package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
)

// VERSION of the service.
const VERSION = "0.3.5"

type params struct {
	Port       string `env:"PORT" envDefault:"9191"`
	CorsOrigin string `env:"CORS_ORIGIN" envDefault:"*"`
}

func main() {
	var cfg params

	log.Println("App version", VERSION)

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	log.Printf("Server run at http://localhost:%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, &router{cfg: cfg}); err != nil {
		panic(err)
	}
}
