package external

import (
	"context"
	"encoding/json"
	"fmt"
	"http_server/internal/service"
	"net/http"
	"net/url"
)

type GeolocationClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGeolocationClient(httpClient *http.Client) *GeolocationClient {
	return &GeolocationClient{
		httpClient: httpClient,
		baseURL:    "https://geocoding-api.open-meteo.com/v1/search",
	}
}

type openMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

func (c *GeolocationClient) GetLatitudeLongitudeByCity(ctx context.Context, city string) (float64, float64, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, 0, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call external api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	var result openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode external api response: %w", err)
	}

	return &service.ProviderWeatherResponse{
		Temperature: result.CurrentWeather.Temperature,
		WindSpeed:   result.CurrentWeather.Windspeed,
		WeatherCode: result.CurrentWeather.Weathercode,
		Time:        result.CurrentWeather.Time,
	}, nil
}