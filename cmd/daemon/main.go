package main

import (
	"log"
	"time"

	"github.com/sahilium/florence/internal/config"
	"github.com/sahilium/florence/internal/matrix"
	"github.com/sahilium/florence/internal/storage"
	"github.com/sahilium/florence/sources/hackernews"
)

func main() {
	cfg := config.Load()

	db, err := storage.Open(cfg.TursoDatabaseURL, cfg.TursoAuthToken)
	if err != nil {
		log.Fatalf("db: %s", err)
	}
	defer db.Close()

	log.Print("db connected")

	if err := storage.RunMigrations(db.DB, "migrations"); err != nil {
		log.Fatalf("migrations: %s", err)
	}

	log.Print("migrations applied")

	matrixClient := matrix.New(
		cfg.MatrixHomeserver,
		cfg.MatrixToken,
		cfg.Rooms,
	)

	hn := hackernews.New(cfg.HNMinScore)

	cursor, err := db.GetSourceState("hackernews")
	firstRun := cursor == ""
	if err != nil {
		log.Printf("hackernews: get source state: %s", err)
	}

	if firstRun {
		log.Print("hackernews: first run, setting initial state (no notifications)")
		if err := db.SetSourceState("hackernews", "bootstrapped"); err != nil {
			log.Printf("hackernews: set source state: %s", err)
		}
	}

	log.Print("daemon started")

	for {
		events, err := hn.Fetch()
		if err != nil {
			log.Printf("hackernews: fetch: %s", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		if firstRun {
			firstRun = false
			time.Sleep(5 * time.Minute)
			continue
		}

		for _, ev := range events {
			exists, err := db.EventExists(ev.ID)
			if err != nil {
				log.Printf("%s: dedup check: %s", ev.ID, err)
				continue
			}
			if exists {
				log.Printf("%s: already processed, skipping", ev.ID)
				continue
			}

			err = matrixClient.Send(ev)
			if err != nil {
				log.Printf("%s: matrix delivery failed: %s", ev.ID, err)
				if rErr := db.RecordDelivery(ev.ID, "matrix", false, err.Error()); rErr != nil {
					log.Printf("%s: record delivery: %s", ev.ID, rErr)
				}
				continue
			}

			log.Printf("%s: matrix delivery succeeded", ev.ID)

			if err := db.StoreEvent(ev); err != nil {
				log.Printf("%s: store: %s", ev.ID, err)
				continue
			}

			if err := db.RecordDelivery(ev.ID, "matrix", true, ""); err != nil {
				log.Printf("%s: record delivery: %s", ev.ID, err)
				continue
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
