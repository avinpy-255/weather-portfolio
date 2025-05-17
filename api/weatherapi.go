package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "os"
)

var weatherCodeMap = map[int]string{
	0:  "Clear",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Freezing light drizzle",
	57: "Freezing dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Freezing slight rain",
	67: "Freezing heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
	} `json:"current_weather"`
	Description string
}

func GetWeather(lat, lon string) (*WeatherResponse, error) {
	// apiKey := os.Getenv("OPENWEATHER_API_KEY")
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,weathercode&timezone=auto", lat, lon)

	fmt.Println("Calling weather API:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API response status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature_2m"`
			Humidity    float64 `json:"relative_humidity_2m"`
			Weathercode int     `json:"weathercode"`
			Interval    int     `json:"interval"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	desc := weatherCodeMap[parsed.CurrentWeather.Weathercode]

	return &WeatherResponse{
		CurrentWeather: struct {
			Temperature float64 `json:"temperature"`
			Windspeed   float64 `json:"windspeed"`
			Weathercode int     `json:"weathercode"`
		}{
			Temperature: parsed.CurrentWeather.Temperature,
			Windspeed:   0,
			Weathercode: parsed.CurrentWeather.Weathercode,
		},
		Description: desc,
	}, nil
}
