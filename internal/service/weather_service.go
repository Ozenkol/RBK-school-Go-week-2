package service

import (
	"context"
	"fmt"
)

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*ProviderWeatherResponse, error)
	GetLatitudeLongitudeByCity(ctx context.Context, city string) (float64, float64, error)
}

type ProviderWeatherResponse struct {
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	Time        string
}

type WeatherResult struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
}

type WeatherService struct {
	weatherProvider WeatherProvider
	geolocationProvider WeatherProvider
}

func NewWeatherService(weatherProvider WeatherProvider, geolocationProvider WeatherProvider) *WeatherService {
	return &WeatherService{
		weatherProvider: weatherProvider,
		geolocationProvider: geolocationProvider,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float64) (*WeatherResult, error) {
	resp, err := s.weatherProvider.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("get weather from provider: %w", err)
	}

	return &WeatherResult{
		Latitude:    lat,
		Longitude:   lon,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
	}, nil
}

func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*WeatherResult, error) {
	lat, lon, err := s.geolocationProvider.GetLatitudeLongitudeByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("get latitude and longitude for city: %w", err)
	}

	return s.GetWeather(ctx, lat, lon)
}

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}
