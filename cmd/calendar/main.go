package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/nivanov045/calendar/internal/api"
	"github.com/nivanov045/calendar/internal/config"
	"github.com/nivanov045/calendar/internal/service"
	"github.com/nivanov045/calendar/internal/storage"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.BuildConfig()
	if err != nil {
		log.Panic().Err(err).Stack()
	}

	myStorage := storage.New()

	serv := service.New(myStorage)

	myapi := api.New(serv)

	log.Panic().Err(myapi.Run(cfg.Address)).Stack()
}
