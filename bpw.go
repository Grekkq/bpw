package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	sys_parameter_name, dia_parameter_name, pulse_parameter_name := "sys", "dia", "pulse"
	parsedSys, err := parseIntFromUrl(sys_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid %v. %v", sys_parameter_name, err.Error()), http.StatusInternalServerError)
		return
	}
	parsedDia, err := parseIntFromUrl(dia_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid %v. %v", dia_parameter_name, err.Error()), http.StatusInternalServerError)
		return
	}
	parsedPulse, err := parseIntFromUrl(pulse_parameter_name, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Please provide valid %v. %v", pulse_parameter_name, err.Error()), http.StatusInternalServerError)
		return
	}
	entry := Entry{parsedSys, parsedDia, parsedPulse, int(time.Now().Unix())}
	fmt.Fprintf(w, "Entry added test: %v", entry)
}

func main() {
	t := time.Now()
	fmt.Println(t.Unix())
	http.HandleFunc("/addEntry/", addEntry)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
