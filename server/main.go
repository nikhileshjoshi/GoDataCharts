package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type statePopulation struct {
	State      string `json:"state"`
	Population int64  `json:"population"`
}

func main() {
	http.HandleFunc("/data", getData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getData(w http.ResponseWriter, req *http.Request) {
	data, err := readData("../data/StatePopulation.csv")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func readData(filePath string) ([]statePopulation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	data := []statePopulation{}

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		//fmt.Println(record[0], record[1])
		if population, err := strconv.Atoi(record[1]); err != nil {
			log.Println("Could not convert", record[1], "to int, So record for state", record[0], "ignored.")
			continue
		} else {
			p := statePopulation{record[0], int64(population)}
			data = append(data, p)
		}

	}
	return data, nil
}
