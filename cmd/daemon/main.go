package main

import (
	"log"
	"time"

	"github.com/sahilium/florence/internal/config"
	"github.com/sahilium/florence/internal/matrix"
	"github.com/sahilium/florence/sources/hackernews"
)

func main() {
	cfg := config.Load()

	matrixClient := &matrix.Client{
		Homeserver: cfg.MatrixHomeserver,
		Token:      cfg.MatrixToken,
		RoomID:     cfg.MatrixRoomID,
	}

	hn := hackernews.New(cfg.HNMinScore)

	log.Println("🚀 daemon started")

	for {
		events, err := hn.Fetch()
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, ev := range events {
			err := matrixClient.Send(ev)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("sent:", ev.Title)
		}

		time.Sleep(5 * time.Minute)
	}
}
