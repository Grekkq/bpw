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

func EntryAddHandle(server *Server, w http.ResponseWriter, r *http.Request) {
	entry, err := parseEntry(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = pushToDb(entry, r.Context(), server.containerClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Successfully saved entry in database ", entry)
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

func parseIntFromUrl(parameterName string, r *http.Request) (int, error) {
	providedValue := r.URL.Query().Get(parameterName)
	parsedValue, err := strconv.Atoi(providedValue)
	if err != nil {
		return 0, errors.New(fmt.Sprint("Cannot parse:", providedValue))
	}
	return parsedValue, nil
}

func pushToDb(entry Entry, context context.Context, containerClient *azcosmos.ContainerClient) error {
	marshalled, err := json.Marshal(entry)
	if err != nil {
		log.Print("Cannot serialize provided object.\n", err)
		return fmt.Errorf("cannot save in database\nPlease contact administrator")
	}

	pk := azcosmos.NewPartitionKeyString(entry.Id)
	_, err = containerClient.CreateItem(context, pk, marshalled, nil)
	if err != nil {
		log.Print("Failed saving data to db.\n", err)
		return fmt.Errorf("cannot save in database\nPlease try again later")
	}
	return nil
}
