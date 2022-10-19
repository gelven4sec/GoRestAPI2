package main

import (
	"encoding/json"
	"log"
	"os"
)

type Attraction struct {
	Id           uint
	Name         string
	InPark       string
	Place        string
	Manufacturer string
}

func initAttractions() []Attraction {
	rawContent, err := os.ReadFile("./attractions.json")
	if err != nil {
		rawContent = []byte("[]")
		err = os.WriteFile("./attractions.json", rawContent, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	var content []Attraction
	err = json.Unmarshal(rawContent, &content)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return content
}

func syncFile() {
	jsonData, _ := json.Marshal(listAttractions)
	os.WriteFile("./attractions.json", jsonData, 0700)
}

func getAttractions() string {
	jsonData, _ := json.Marshal(listAttractions)
	return string(jsonData)
}

func getAttraction(id uint) (string, bool) {
	for _, attraction := range listAttractions {
		if attraction.Id == id {
			jsonData, _ := json.Marshal(attraction)
			return string(jsonData), true
		}
	}
	return "Error: Not found", false
}

func getLastId() uint {
	if len(listAttractions) == 0 {
		return 1
	} else {
		return listAttractions[len(listAttractions)-1].Id
	}
}

func createAttraction(attraction Attraction) {
	listAttractions = append(listAttractions, attraction)
	syncFile()
}

func getAttractionIndex(id uint) int {
	for index, attraction := range listAttractions {
		if attraction.Id == id {
			return index
		}
	}
	return -1
}
