package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func ListEntriesHandle(server *Server, w http.ResponseWriter, r *http.Request) {
	userId, err := readUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	entries, err := getEntriesFromDb(userId, r.Context(), server.measurements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		log.Print("Failed to serialize entries\n", err)
		http.Error(w, "Please try again later", http.StatusInternalServerError)
		return
	}
}

func readUserId(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	var jsonBody map[string]string
	err := decoder.Decode(&jsonBody)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("cannot parse request, check provided parameters")
	}
	return jsonBody["userId"], nil
}

func getEntriesFromDb(userId string, context context.Context, measurements *azcosmos.ContainerClient) ([]Entry, error) {
	pk := azcosmos.NewPartitionKeyString(userId)
	queryPager := measurements.NewQueryItemsPager("SELECT * FROM c", pk, nil)

	var entries = make([]Entry, 0, 10)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context)
		if err != nil {
			log.Print(err)
			return entries, fmt.Errorf("Server encounterd an error, please try again later")
		}

		for _, item := range queryResponse.Items {
			var entry Entry
			err = json.Unmarshal(item, &entry)
			if err != nil {
				log.Print("Failed to unpack item from database\n", err)
				continue
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}
