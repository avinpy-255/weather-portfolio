package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LocationData struct {
	IP          string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
}

type IPAPI struct{}

func (api *IPAPI) GetLocation(ctx context.Context, ip string) (*LocationData, error) {

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get location: %v", err)
	}
	defer resp.Body.Close()

	var data LocationData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Status != "success" {
		return nil, fmt.Errorf("IP lookup failed")
	}

	return &data, nil
}
