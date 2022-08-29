package main

import (
	"log"

	"github.com/nivanov045/curly-waffle/cmd/calendar/api"
	"github.com/nivanov045/curly-waffle/cmd/calendar/config"
	"github.com/nivanov045/curly-waffle/cmd/calendar/service"
	"github.com/nivanov045/curly-waffle/cmd/calendar/storage"
)

func main() {
	cfg, err := config.BuildConfig()
	if err != nil {
		log.Fatalln("calendar::main::error: in env parsing:", err)
	}
	log.Println("calendar::main::info: cfg:", cfg)
	myStorage := storage.New()
	serv := service.New(myStorage)
	myapi := api.New(serv)
	log.Fatalln(myapi.Run(cfg.Address))
}
