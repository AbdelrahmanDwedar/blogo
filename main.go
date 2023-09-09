package main

import (
	"AbdelrahmanDwedar/blogo/tables"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	r := mux.NewRouter()

	store := tables.NewPostgresStore()
	defer func() {
		err := store.Client.Close()
		if err != nil {
			return
		}
	}()

	server := NewServer(store)

	r.HandleFunc("/ping", server.Ping)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
