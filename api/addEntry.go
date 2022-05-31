package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	err = pushToDb(entry, r.Context(), server.measurements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Successfully saved entry in database ", entry)
}

func parseEntry(r *http.Request) (Entry, error) {
	decoder := json.NewDecoder(r.Body)
	var entry Entry
	err := decoder.Decode(&entry)
	if err != nil {
		log.Print(err)
		switch t := err.(type) {
		default:
			return Entry{}, fmt.Errorf("failed to parse request, check provided vaules")
		case *json.UnmarshalTypeError:
			return Entry{}, fmt.Errorf("cannot parse value for %v, please check if you provided correct value", t.Field)
		}
	}

	entry.Id = uuid.NewString()
	entry.Timestamp = int(time.Now().Unix())
	return entry, nil
}

func pushToDb(entry Entry, context context.Context, measurements *azcosmos.ContainerClient) error {
	marshalled, err := json.Marshal(entry)
	if err != nil {
		log.Print("Cannot serialize provided object.\n", err)
		return fmt.Errorf("cannot save in database\nPlease contact administrator")
	}

	pk := azcosmos.NewPartitionKeyString(entry.UserId)
	_, err = measurements.CreateItem(context, pk, marshalled, nil)
	if err != nil {
		log.Print("Failed saving data to db.\n", err)
		return fmt.Errorf("cannot save in database\nPlease try again later")
	}
	return nil
}
