package handler

import (
	"encoding/json"
	"net/http"
)

type WeatherHandler struct{
	*wu WeatherUseCase
}


