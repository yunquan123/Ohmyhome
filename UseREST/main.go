package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const baseURL = "http://localhost:5000/api/v1/properties"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

type property struct {
	PropertyName string
	Location     string
}

func addProperty(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(baseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func getProperty(code string) {
	url := baseURL
	if code != "" {
		url = baseURL + "/" + code + "?key=" + key
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func updateProperty(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest(http.MethodPut,
		baseURL+"/"+code+"?key="+key,
		bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}
func deleteProperty(code string) {
	request, err := http.NewRequest(http.MethodDelete,
		baseURL+"/"+code+"?key="+key, nil)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func main() {

	for {
		var choice int
		fmt.Println("choose your choice")
		fmt.Println("1: add Property")
		fmt.Println("2: delete Property")
		fmt.Println("3: retrieve Property")
		fmt.Println("4: update Property")
		fmt.Println("5: get all Properties")
		fmt.Scanln(&choice)

		if choice == 1 {
			var PropertyName string
			var Location string
			fmt.Println("key in your Property Name")
			fmt.Scanln(&PropertyName)
			fmt.Println("key in your Location")
			fmt.Scanln(&Location)
			jsonData := map[string]string{"Location": Location}
			addProperty(PropertyName, jsonData)
			fmt.Println("Successfully added!")
			fmt.Println(jsonData)
		}

		if choice == 2 {
			var PropertyName string
			fmt.Println("key in your Property Name to delete")
			fmt.Scanln(&PropertyName)
			deleteProperty(PropertyName)
			fmt.Println("Successfully deleted")
		}

		if choice == 3 {
			var PropertyName string
			fmt.Println("key in your Property Name")
			fmt.Scanln(&PropertyName)
			getProperty(PropertyName)
		}

		if choice == 4 {
			var PropertyName string
			var Location string
			fmt.Println("key in your Property Name")
			fmt.Scanln(&PropertyName)
			fmt.Println("key in your Location")
			fmt.Scanln(&Location)
			jsonData := map[string]string{"Location": Location}
			updateProperty(PropertyName, jsonData)
			fmt.Println("Successfully updated!")
		}

		if choice == 5 {
			getProperty("")
		}
	}
}
