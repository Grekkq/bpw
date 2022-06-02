package main

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Server struct {
	onePass        *OnePass
	databaseClient *azcosmos.DatabaseClient
	users          *azcosmos.ContainerClient
	measurements   *azcosmos.ContainerClient
}

type DatabaseDetails struct {
	endpoint     string
	key          string
	dbName       string
	users        string
	measurements string
}

type OnePass struct {
	HttpAddress              string
	ApiToken                 string
	VaultName                string
	DbConnectionDetailsEntry string
}

type Handler struct {
	*Server
	H func(e *Server, w http.ResponseWriter, r *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: change * to normal url https://developer.mozilla.org/en-US/docs/web/http/cors
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else {
		h.H(h.Server, w, r)
	}

}

type Entry struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Sys       int    `json:"sys"`
	Dia       int    `json:"dia"`
	Pulse     int    `json:"pulse"`
	Timestamp int    `json:"timestamp"`
}
