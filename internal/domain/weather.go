package domain

type Weather struct {
	Country     string `json:"country"`
	Temperature string `json:"temperature"`
	Condition   string `json:"condition"`
}
