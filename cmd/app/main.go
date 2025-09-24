package main

import (
	"github.com/I-Van-Radkov/kaspersky_1/internal/app"
	"github.com/I-Van-Radkov/kaspersky_1/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.MustLoad()

	application := app.New(cfg)

	application.Run()
}
