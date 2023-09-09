package main

import (
	"AbdelrahmanDwedar/blogo/tables"
	"net/http"
)

type APIServer struct {
	store *tables.PostgresStore
}

func NewServer(store *tables.PostgresStore) *APIServer {
	return &APIServer{
		store,
	}
}

func (s APIServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`
		{
			"message": "pong"
		}
	`))
}

func (s APIServer) CreateNewUser()  {
	
}
