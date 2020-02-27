package main

import (
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

const VERSION = "0.3.5"

type Params struct {
	Port    string `env:"PORT" envDefault:"9191"`
	CorsOrigin    string `env:"CORS_ORIGIN" envDefault:"*"`
}

func main() {
	log.Println("App version", VERSION)
	cfg := Params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	log.Printf("Server run at http://localhost:%s", cfg.Port)

	err = http.ListenAndServe(":" + cfg.Port, Router{cfg:cfg})
	if err != nil {
		panic(err)
	}
}
