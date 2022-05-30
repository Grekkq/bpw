package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
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

func parseIntFromUrl(parameterName string, r *http.Request) (int, error) {
	providedValue := r.URL.Query().Get(parameterName)
	parsedValue, err := strconv.Atoi(providedValue)
	if err != nil {
		return 0, errors.New(fmt.Sprint("Cannot parse:", providedValue))
	}
	return parsedValue, nil
}

func parseEntry(r *http.Request) (Entry, error) {
	sysParameterName, diaParameterName, pulseParameterName := "sys", "dia", "pulse"

	parsedSys, err := parseIntFromUrl(sysParameterName, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", sysParameterName, err)
	}

	parsedDia, err := parseIntFromUrl(diaParameterName, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", diaParameterName, err)
	}

	parsedPulse, err := parseIntFromUrl(pulseParameterName, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", pulseParameterName, err)
	}

	return Entry{uuid.NewString(), parsedSys, parsedDia, parsedPulse, int(time.Now().Unix())}, nil
}

func HandleAddEntry(env *Server, w http.ResponseWriter, r *http.Request) {
	entry, err := parseEntry(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = pushToDb(entry, env.containerClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Successfully saved entry in database ", entry)
}

func pushToDb(entry Entry, containerClient *azcosmos.ContainerClient) error {
	marshalled, err := json.Marshal(entry)
	if err != nil {
		log.Print("Cannot serialize provided object.\n", err)
		return fmt.Errorf("cannot save in database\nPlease contact administrator")
	}
	pk := azcosmos.NewPartitionKeyString(entry.Id)
	_, err = containerClient.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		log.Print("Failed saving data to db.\n", err)
		return fmt.Errorf("cannot save in database\nPlease try again later")
	}
	return nil
}
