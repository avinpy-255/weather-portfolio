package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Current struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
		Weather  []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
}

func GetWeather(lat, lon string) (*WeatherResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%s&lon=%s&exclude=minutely,hourly,daily,alerts&appid=%s&units=metric", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API response status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	var data WeatherResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
