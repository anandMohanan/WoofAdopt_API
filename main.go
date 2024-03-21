package main

import (
	"log"

	"github.com/anandMohanan/WoofAdopt_API/api"
	"github.com/anandMohanan/WoofAdopt_API/storage"
)

func main() {
	store, err := storage.NewStore()
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	server := api.APIServer{":8080", store}
	server.Run()
}
