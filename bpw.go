package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/1Password/connect-sdk-go/connect"
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

func addEntry(w http.ResponseWriter, r *http.Request) {
	log.Print("Received http request.")
	sys_parameter_name, dia_parameter_name, pulse_parameter_name := "sys", "dia", "pulse"
	parsedSys, err := parseIntFromUrl(sys_parameter_name, r)
	if err != nil {
		log.Print(err)
		panic(fmt.Sprintf("Please provide valid value for %v.", sys_parameter_name))
	}
	parsedDia, err := parseIntFromUrl(dia_parameter_name, r)
	if err != nil {
		log.Print(err)
		panic(fmt.Sprintf("Please provide valid value for %v.", dia_parameter_name))
	}
	parsedPulse, err := parseIntFromUrl(pulse_parameter_name, r)
	if err != nil {
		log.Print(err)
		panic(fmt.Sprintf("Please provide valid value for %v.", pulse_parameter_name))
	}
	entry := Entry{parsedSys, parsedDia, parsedPulse, int(time.Now().Unix())}

	dockerSecrets, _ := secrets.NewDockerSecrets("")
	address, _ := dockerSecrets.Get("1password_http_address")
	token, _ := dockerSecrets.Get("1password_api_token")
	vaultName, _ := dockerSecrets.Get("1password_vault_name")
	dbConnectionStringName, _ := dockerSecrets.Get("db_connection_string_name")

	client := connect.NewClient(address, token)
	item, err := client.GetItemByTitle(dbConnectionStringName, vaultName)
	if err != nil {
		log.Print(err)
		panic("Failed to connect to database.")
	}

	for i := range item.Fields {
		if item.Fields[i].Label == "password" {
			log.Print(item.Fields[i].Value)
			break
		}
	}

	fmt.Fprintf(w, "Entry added test: %v", entry)
}

func handleAddEntry(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)
	addEntry(w, r)
}

func handlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, fmt.Sprintf("%v\nServer encounter an error, please try again later.", r), http.StatusInternalServerError)
	}
}

func main() {
	t := time.Now()
	fmt.Println(t.Unix())

	http.HandleFunc("/addEntry/", handleAddEntry)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
