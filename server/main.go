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

type PopulationDataForChart struct {
	States     []string `json:"states"`
	Population []int64  `json:"population"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../web")))
	http.HandleFunc("/data", getData)
	log.Fatal(http.ListenAndServe(":8080", nil))
	//readData("../data/StatePopulation.csv")
}

func getData(w http.ResponseWriter, req *http.Request) {
	chartData, err := readData("../data/StatePopulation.csv")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(chartData)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chartData)
}

func readData(filePath string) (*PopulationDataForChart, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	linecount := 0
	var states []string
	var populations []int64
	for {
		record, err := r.Read()
		linecount++
		if linecount == 1 {
			continue //remove header
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		//fmt.Println(record[0], record[1])
		if population, err := strconv.Atoi(record[1]); err != nil {
			log.Println("Could not convert population value: ", record[1], "to int, So record for state", record[0], "ignored.")
			continue
		} else {
			populations = append(populations, int64(population))
			states = append(states, record[0])
		}

	}
	//fmt.Println(states)
	return &PopulationDataForChart{states, populations}, nil
	//return nil, nil, nil
}
