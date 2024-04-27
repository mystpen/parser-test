package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/mystpen/parser-test/config"
	"github.com/mystpen/parser-test/internal/parcer"
	"github.com/mystpen/parser-test/internal/storage"
)

var ErrNoResponce = errors.New("no responce from API")

func main() {
	resp, err := http.Get(config.Config.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error:", "no responce from API")
		return
	}

	Influencers, err := parcer.Parse(resp)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.CreateCSV(Influencers)
	if err != nil {
		log.Fatal(err)
	}
}
