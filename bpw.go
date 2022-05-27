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

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	secrets "github.com/ijustfool/docker-secrets"
)

type Entry struct {
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

func addEntry(w http.ResponseWriter, r *http.Request, dockerSecrets *secrets.DockerSecrets) {
	log.Print("Received http request.")
	sys_parameter_name, dia_parameter_name, pulse_parameter_name := "sys", "dia", "pulse"
	parsedSys, err := parseIntFromUrl(sys_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid value for %v.\n%v", sys_parameter_name, err), http.StatusUnprocessableEntity)
		return
	}
	parsedDia, err := parseIntFromUrl(dia_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid value for %v.\n%v", dia_parameter_name, err), http.StatusUnprocessableEntity)
		return
	}
	parsedPulse, err := parseIntFromUrl(pulse_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid value for %v.\n%v", pulse_parameter_name, err), http.StatusUnprocessableEntity)
		return
	}
	entry := Entry{parsedSys, parsedDia, parsedPulse, int(time.Now().Unix())}

	fmt.Fprintf(w, "Entry added test: %v", entry)
}

func handleAddEntry(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)
	dockerSecrets, err := secrets.NewDockerSecrets("")
	if err != nil {
		log.Print("Failed to initialize docker secrets\n", err)
		panic("Cannot access database.\n")
	}
	addEntry(w, r, dockerSecrets)

	pushToDb(dockerSecrets)

}

func pushToDb(dockerSecrets *secrets.DockerSecrets) {
	address, _ := dockerSecrets.Get("1password_http_address")
	token, _ := dockerSecrets.Get("1password_api_token")
	vaultName, _ := dockerSecrets.Get("1password_vault_name")
	dbConnectionStringName, _ := dockerSecrets.Get("db_connection_string_name")

	client := connect.NewClient(address, token)
	item, err := client.GetItemByTitle(dbConnectionStringName, vaultName)
	if err != nil {
		log.Print(err)
		panic("Cannot access database.\n")
	}
	var dbKey, dbEndpoint, dbName, containerName string
	for i := range item.Fields {
		if item.Fields[i].Label == "password" {
			dbKey = item.Fields[i].Value
		}
		if item.Fields[i].Label == "endpointAddress" {
			dbEndpoint = item.Fields[i].Value
		}
		if item.Fields[i].Label == "databaseName" {
			dbName = item.Fields[i].Value
		}
		if item.Fields[i].Label == "containerName" {
			containerName = item.Fields[i].Value
		}
	}
	cred, err := azcosmos.NewKeyCredential(dbKey)
	handle(err)
	dbClient, err := azcosmos.NewClientWithKey(dbEndpoint, cred, nil)
	handle(err)
	database, err := dbClient.NewDatabase(dbName)
	handle(err)
	container, err := database.NewContainer(containerName)
	handle(err)
	log.Print(container)
	id := uuid.NewString()
	testItem := map[string]string{"id": id, "otherValue": "10"}
	marshalled, err := json.Marshal(testItem)
	handle(err)
	pk := azcosmos.NewPartitionKeyString(id)
	itemResponse, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
	handle(err)
	log.Print(itemResponse)
}

func handle(e error) {
	if e != nil {
		log.Print(e)
	}
}

func handlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, fmt.Sprint(r, "Server encounter an error, please try again later."), http.StatusInternalServerError)
	}
}

func main() {
	t := time.Now()
	fmt.Println(t.Unix())

	http.HandleFunc("/addEntry/", handleAddEntry)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
