package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type propertyInfo struct {
	Location string `json:"Location"`
}

var properties map[string]propertyInfo

func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" { //this key should be stored in a database
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Property Portal!")
}

func allproperty(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	json.NewEncoder(w).Encode(properties)
}

func property(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)
	if r.Method == "GET" {
		if _, ok := properties[params["PropertyName"]]; ok {
			json.NewEncoder(w).Encode(
				properties[params["PropertyName"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No property found"))
		}
	}

	if r.Method == "DELETE" {
		if _, ok := properties[params["PropertyName"]]; ok {
			delete(properties, params["PropertyName"])
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No property found"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {

		// POST is for creating new entry
		if r.Method == "POST" {
			// read the string sent to the service
			var newProperty propertyInfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newProperty)
				if newProperty.Location == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply property " +
							"information " + "in JSON format"))
					return
				}
				// check if entry exists; add only if
				// entry does not exist
				if _, ok := properties[params["PropertyName"]]; !ok {
					properties[params["PropertyName"]] = newProperty
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Property added: " +
						params["PropertyName"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate property name"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply property information " +
					"in JSON format"))
			}
		}

		if r.Method == "PUT" {
			var newProperty propertyInfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &newProperty)
				if newProperty.Location == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply property " +
							" information " +
							"in JSON format"))
					return
				}
				// check if entry exists; add only if
				// entry does not exist
				if _, ok := properties[params["PropertyName"]]; !ok {
					properties[params["PropertyName"]] =
						newProperty
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Property added: " +
						params["PropertyName"]))
				} else {
					// update entry
					properties[params["PropertyName"]] = newProperty
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " +
					"property information " +
					"in JSON format"))
			}
		}
	}
}

func main() {
	// instantiate properties
	properties = make(map[string]propertyInfo)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/properties", allproperty)
	router.HandleFunc("/api/v1/properties/{PropertyName}", property).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	http.ListenAndServe(":5000", router)
}
