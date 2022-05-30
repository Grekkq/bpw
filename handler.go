package main

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Server struct {
	onePass         *OnePass
	databaseClient  *azcosmos.DatabaseClient
	containerClient *azcosmos.ContainerClient
}

type DatabaseDetails struct {
	endpoint      string
	key           string
	dbName        string
	containerName string
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
	h.H(h.Server, w, r)
}

type Entry struct {
	Id        string `json:"id"`
	Sys       int    `json:"sys"`
	Dia       int    `json:"dia"`
	Pulse     int    `json:"pulse"`
	Timestamp int    `json:"timestamp"`
}
