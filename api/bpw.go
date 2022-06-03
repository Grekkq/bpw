package main

import (
	"log"
	"net/http"
	"time"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	secrets "github.com/ijustfool/docker-secrets"
)

func main() {
	server := &Server{}
	loadDockerSecrets(server)
	loadDatabaseConnection(server)

	http.Handle("/addEntry", Handler{server, EntryAddHandle})
	http.Handle("/listEntries", Handler{server, ListEntriesHandle})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getOnePassSecrets(secrets *OnePass) DatabaseDetails {
	client := connect.NewClient(secrets.HttpAddress, secrets.ApiToken)
	item, err := client.GetItemByTitle(secrets.DbConnectionDetailsEntry, secrets.VaultName)
	if err != nil {
		log.Print("Failed to get data from 1Password.\nRetrying in 5s", err)
		time.Sleep(time.Second * 5)
		item, err = client.GetItemByTitle(secrets.DbConnectionDetailsEntry, secrets.VaultName)
		if err != nil {
			log.Panic("Failed to get data from 1Password.\nCheck 1password server", err)
		}
	}

	dbDetails := DatabaseDetails{}
	for i := range item.Fields {
		if item.Fields[i].Label == "password" {
			dbDetails.key = item.Fields[i].Value
			continue
		}
		if item.Fields[i].Label == "endpointAddress" {
			dbDetails.endpoint = item.Fields[i].Value
			continue
		}
		if item.Fields[i].Label == "databaseName" {
			dbDetails.dbName = item.Fields[i].Value
			continue
		}
		if item.Fields[i].Label == "usersContainerName" {
			dbDetails.users = item.Fields[i].Value
			continue
		}
		if item.Fields[i].Label == "measurementsContainerName" {
			dbDetails.measurements = item.Fields[i].Value
			continue
		}
	}

	return dbDetails
}

func loadDatabaseConnection(server *Server) {
	dbDetails := getOnePassSecrets(server.onePass)

	cred, err := azcosmos.NewKeyCredential(dbDetails.key)
	if err != nil {
		log.Panic("Failed to get cosmos key\n", err)
	}

	dbClient, err := azcosmos.NewClientWithKey(dbDetails.endpoint, cred, nil)
	if err != nil {
		log.Panic("Failed to get cosmos client\n", err)
	}

	database, err := dbClient.NewDatabase(dbDetails.dbName)
	if err != nil {
		log.Panic("Failed to get cosmos database\n", err)
	}

	users, err := database.NewContainer(dbDetails.users)
	if err != nil {
		log.Panic("Failed to get users cosmos container\n", err)
	}

	measurements, err := database.NewContainer(dbDetails.measurements)
	if err != nil {
		log.Panic("Failed to get measurements cosmos container\n", err)
	}

	server.databaseClient = database
	server.users = users
	server.measurements = measurements
}

func loadDockerSecrets(env *Server) {
	dockerSecrets, err := secrets.NewDockerSecrets("")
	if err != nil {
		log.Panic("Failed to initialize docker secrets\n", err)
	}

	onepass := OnePass{}
	onepass.HttpAddress, err = dockerSecrets.Get("1password_http_address")
	if err != nil {
		log.Panic("Failed to load docker secret 1password_http_address\n", err)
	}

	onepass.ApiToken, err = dockerSecrets.Get("1password_api_token")
	if err != nil {
		log.Panic("Failed to load docker secret 1password_api_token\n", err)
	}

	onepass.VaultName, err = dockerSecrets.Get("1password_vault_name")
	if err != nil {
		log.Panic("Failed to load docker secret 1password_vault_name\n", err)
	}

	onepass.DbConnectionDetailsEntry, err = dockerSecrets.Get("db_connection_string_name")
	if err != nil {
		log.Panic("Failed to load docker secret db_connection_string_name\n", err)
	}

	env.onePass = &onepass
}
