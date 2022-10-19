package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ErrorResponse Error response structure
type ErrorResponse struct {
	Error string
}

// IdResponse Id response structure
type IdResponse struct {
	Id uint
}

// Helper function to return error response
func returnError(err string) (string, int) {
	errorResponse := ErrorResponse{Error: err}
	jsonData, _ := json.Marshal(errorResponse)
	return string(jsonData), 400
}

func checkIfIdExist(id uint) bool {
	exist := true
	for _, attr := range listAttractions {
		if attr.Id == id {
			exist = false
		}
	}
	return exist
}

// Handle GET request
func getHandler(args url.Values) (string, int) {
	// If param id in form then return the requested attraction
	if id, ok := args["id"]; ok {

		// Convert id to number
		idNumber, err := strconv.Atoi(id[0])
		if err != nil {
			return returnError(err.Error())
		}
		// Check if id is negative
		if idNumber < 0 {
			return returnError("Id can not be negative")
		}

		response, found := getAttraction(uint(idNumber))
		if found {
			return response, 200
		} else {
			return returnError(response)
		}
	}

	// Else return every attraction
	return getAttractions(), 200
}

// Handle POST request
func postHandler(args url.Values) (string, int) {
	newAttraction := Attraction{}

	// Generate id
	for {
		newAttraction.Id = uint(rand.Intn(99999999999-10000000000)+9999999999) % 100000000000
		if checkIfIdExist(newAttraction.Id) {
			break
		}
	}

	// Get Name
	if arg, ok := args["name"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'name' too long")
		}
		newAttraction.Name = arg[0]
	} else {
		return returnError("Missing parameter 'name'")
	}

	// Get InPark
	if arg, ok := args["in_park"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'in_park' too long")
		}
		newAttraction.InPark = arg[0]
	} else {
		return returnError("Missing parameter 'in_park'")
	}

	// Get Place
	if arg, ok := args["place"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'place' too long")
		}
		newAttraction.Place = arg[0]
	} else {
		return returnError("Missing parameter 'place'")
	}

	// Get Manufacturer
	if arg, ok := args["manufacturer"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'manufacturer' too long")
		}
		newAttraction.Manufacturer = arg[0]
	} else {
		return returnError("Missing parameter 'manufacturer'")
	}

	// Append to list and sync file
	createAttraction(newAttraction)

	// Return newly created attraction's id
	idResponse := IdResponse{Id: newAttraction.Id}
	jsonData, _ := json.Marshal(idResponse)
	return string(jsonData), 200
}

// Handle PUT request
func putHandler(args url.Values) (string, int) {
	var id uint
	if _, ok := args["id"]; ok {

		// Convert id to number
		idNumber, err := strconv.Atoi(args["id"][0])
		if err != nil {
			return returnError(err.Error())
		}
		// Check if id is negative
		if idNumber < 0 {
			return returnError("Id can not be negative")
		}

		id = uint(idNumber)
	} else {
		return returnError("Missing id parameter")
	}

	index := getAttractionIndex(id)
	if index == -1 {
		return returnError("Not found")
	}

	if arg, ok := args["name"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'name' too long")
		}
		listAttractions[index].Name = arg[0]
	}

	if arg, ok := args["in_park"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'in_park' too long")
		}
		listAttractions[index].InPark = arg[0]
	}

	if arg, ok := args["place"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'place' too long")
		}
		listAttractions[index].Place = arg[0]
	}

	if arg, ok := args["manufacturer"]; ok {
		if len(arg[0]) > 255 {
			return returnError("Parameter 'manufacturer' too long")
		}
		listAttractions[index].Manufacturer = arg[0]
	}

	syncFile()

	// Return attraction's id
	idResponse := IdResponse{Id: listAttractions[index].Id}
	jsonData, _ := json.Marshal(idResponse)
	return string(jsonData), 200
}

// Handle DELETE request
func deleteHandler(args url.Values) (string, int) {
	var id uint
	if _, ok := args["id"]; ok {

		// Convert id to number
		idNumber, err := strconv.Atoi(args["id"][0])
		if err != nil {
			return returnError(err.Error())
		}
		// Check if id is negative
		if idNumber < 0 {
			return returnError("Id can not be negative")
		}

		id = uint(idNumber)
	} else {
		return returnError("Missing id parameter")
	}

	index := getAttractionIndex(id)
	if index == -1 {
		return returnError("Not found")
	}

	// Keep id to return later
	id = listAttractions[index].Id

	// Remove attraction from list (ptdr la syntax)
	listAttractions = append(listAttractions[:index], listAttractions[index+1:]...)

	syncFile()

	// Return attraction's id
	idResponse := IdResponse{Id: id}
	jsonData, _ := json.Marshal(idResponse)
	return string(jsonData), 200
}

// Handle Attraction request
func handleAttractions(w http.ResponseWriter, r *http.Request) {
	// Log
	fmt.Println(time.Now().Format("02-01-2006 15:04:05") + " " + r.Method + " " + r.RequestURI)

	response := ""    // Default response
	statusCode := 400 // Default status code

	err := r.ParseForm()
	if err != nil {
		response = "Error: " + err.Error()
		fmt.Fprintf(w, response)
		w.WriteHeader(statusCode) // Bad request
		return
	}

	header := w.Header()
	header.Add("Content-Type", "text/json")

	switch r.Method {
	case "GET":
		response, statusCode = getHandler(r.Form)
		w.WriteHeader(statusCode)
	case "POST":
		response, statusCode = postHandler(r.PostForm)
		w.WriteHeader(statusCode)
	case "PUT":
		response, statusCode = putHandler(r.PostForm)
		w.WriteHeader(statusCode)
	case "DELETE":
		response, statusCode = deleteHandler(r.Form)
		w.WriteHeader(statusCode)

	default:
		w.WriteHeader(405) // Bad method
		response = "Error: Unhandled method"
	}

	// Write body response
	fmt.Fprintf(w, response)
}
