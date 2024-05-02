// Package main is the entry point for this application.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Payload is the structure for the payload of a coupon.
type Payload struct {
	Discount       int    `json:"discount"`
	Code           string `json:"code"`
	MinBasketValue int    `json:"minBasketValue"`
}

// GetPayload is the structure for the payload of a get request for coupons.
type GetPayload struct {
	Codes []string `json:"codes"`
}

// ApplyPayload is the structure for the payload of an apply request.
type ApplyPayload struct {
	Basket struct {
		Value int `json:"value"`
	} `json:"basket"`
	Code string `json:"code"`
}

// randomString generates a random string of length n.
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// main is the main function for this application.
func main() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())

	// Parse the number of requests from the command-line argument.
	numRequests, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error parsing command-line argument:", err)
		return
	}

	// Make the specified number of requests.
	for i := 0; i < numRequests; i++ {
		// Generate a random coupon code.
		code := randomString(10)

		// Create the payload for the coupon.
		payload := Payload{
			Discount:       10,
			Code:           code,
			MinBasketValue: 50,
		}

		// Marshal the payload into JSON.
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		// Create a new reader with the JSON payload.
		body := bytes.NewReader(payloadBytes)

		// Create a new POST request.
		req, err := http.NewRequest("POST", "http://localhost:8080/api/create", body)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the content type of the request to JSON.
		req.Header.Set("Content-Type", "application/json")

		// Send the request.
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Print the response body.
		fmt.Println(string(respBody))

		// Create the payload for the get request.
		getPayload := GetPayload{
			Codes: []string{code},
		}

		// Marshal the get payload into JSON.
		getPayloadBytes, err := json.Marshal(getPayload)
		if err != nil {
			fmt.Println("Error marshalling get payload:", err)
			return
		}

		// Create a new reader with the get payload.
		getBody := bytes.NewReader(getPayloadBytes)

		// Create a new GET request.
		getReq, err := http.NewRequest("GET", "http://localhost:8080/api/coupons", getBody)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			return
		}

		// Set the content type of the get request to JSON.
		getReq.Header.Set("Content-Type", "application/json")

		// Send the get request.
		getResp, err := http.DefaultClient.Do(getReq)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}
		defer getResp.Body.Close()

		// Read the get response body.
		getRespBody, err := ioutil.ReadAll(getResp.Body)
		if err != nil {
			fmt.Println("Error reading GET response body:", err)
			return
		}

		// Print the get response body.
		fmt.Println(string(getRespBody))

		// Create the payload for the apply request.
		applyPayload := ApplyPayload{
			Basket: struct {
				Value int `json:"value"`
			}{
				Value: 30 + rand.Intn(40),
			},
			Code: code,
		}

		// Marshal the apply payload into JSON.
		applyPayloadBytes, err := json.Marshal(applyPayload)
		if err != nil {
			fmt.Println("Error marshalling apply payload:", err)
			return
		}

		// Create a new reader with the apply payload.
		applyBody := bytes.NewReader(applyPayloadBytes)

		// Create a new POST request for the apply endpoint.
		applyReq, err := http.NewRequest("POST", "http://localhost:8080/api/apply", applyBody)
		if err != nil {
			fmt.Println("Error creating apply request:", err)
			return
		}

		// Set the content type of the apply request to JSON.
		applyReq.Header.Set("Content-Type", "application/json")

		// Send the apply request.
		applyResp, err := http.DefaultClient.Do(applyReq)
		if err != nil {
			fmt.Println("Error making apply request:", err)
			return
		}
		defer applyResp.Body.Close()

		// Read the apply response body.
		applyRespBody, err := ioutil.ReadAll(applyResp.Body)
		if err != nil {
			fmt.Println("Error reading apply response body:", err)
			return
		}

		// Print the apply response body.
		fmt.Println(string(applyRespBody))
	}
}
