package main

import (
	"log"
)

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}
	 store.Init();
	
	server := APIServer{":8080", store}
	server.Run()
}
