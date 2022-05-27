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
	id        string
	sys       int
	dia       int
	pulse     int
	timestamp int
}

func parseIntFromUrl(parameter_name string, r *http.Request) (int, error) {
	providedValue := r.URL.Query().Get(parameter_name)
	parsedValue, err := strconv.Atoi(providedValue)
	if err != nil {
		return 0, errors.New(fmt.Sprint("Cannot parse:", providedValue))
	}
	return parsedValue, nil
}

func parseEntry(r *http.Request) (Entry, error) {
	sys_parameter_name, dia_parameter_name, pulse_parameter_name := "sys", "dia", "pulse"

	parsedSys, err := parseIntFromUrl(sys_parameter_name, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", sys_parameter_name, err)
	}

	parsedDia, err := parseIntFromUrl(dia_parameter_name, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", dia_parameter_name, err)
	}

	parsedPulse, err := parseIntFromUrl(pulse_parameter_name, r)
	if err != nil {
		return Entry{}, fmt.Errorf("please provide valid value for %v\n%v", pulse_parameter_name, err)
	}

	return Entry{uuid.NewString(), parsedSys, parsedDia, parsedPulse, int(time.Now().Unix())}, nil
}

func HandleAddEntry(env *Server, w http.ResponseWriter, r *http.Request) {
	entry, err := parseEntry(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	log.Print(entry)
	pushToDb(env.containerClient)

}

func pushToDb(containerClient *azcosmos.ContainerClient) {
	id := uuid.NewString()
	testItem := map[string]string{"id": id, "otherValue": "10"}
	marshalled, _ := json.Marshal(testItem)
	pk := azcosmos.NewPartitionKeyString(id)
	itemResponse, _ := containerClient.CreateItem(context.TODO(), pk, marshalled, nil)
	log.Print(itemResponse)
}
