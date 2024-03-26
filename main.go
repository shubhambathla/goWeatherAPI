package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define the base URL for the GoWeather API
const weatherAPIBaseURL = "https://goweather.herokuapp.com/weather/"

// WeatherInfo represents the structure of the weather information we want to capture
// Note: Adjust these fields based on the actual API response structure
type WeatherInfo struct {
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	Description string `json:"description"`
}

func fetchWeather(cityName string) (*WeatherInfo, error) {
	// Construct the full API request URL with the city name
	requestURL := fmt.Sprintf("%s%s", weatherAPIBaseURL, cityName)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error retrieving weather data: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response Content-Type is JSON before parsing
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		return nil, fmt.Errorf("expected JSON response, got: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response data: %v", err)
	}

	var weather WeatherInfo
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("error parsing weather data: %v", err)
	}

	return &weather, nil
}

func cityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	var cityName string

	// Determine the request type
	if r.Method == "GET" {
		// Extract the city name from the query parameter
		cityName = r.URL.Query().Get("name")
	} else if r.Method == "POST" {
		// Extract the city name from the JSON body
		var requestData struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}
		cityName = requestData.Name
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Fetch weather information
	weather, err := fetchWeather(cityName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the weather information as JSON
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/city", cityWeatherHandler)
	fmt.Println("Server is started on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
