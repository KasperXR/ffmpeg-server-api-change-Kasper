package main

type JSONObj struct {
	Payload []VideoObj `json:"payload"`
}

type VideoObj struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	FileName      string         `json:"fileName"`
	ParentOptions []ParentOption `json:"parentOptions"`
}

type ParentOption struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DescTime     int    `json:"time"`
	AudioName    string `json:"audioName"`
	Introduction string `json:"introduction"`
	Active       bool   `json:"active"`
	//NegativeName string   `json:"negativeName"`
	//IsNegative   bool     `json:"isNegative"`
	Options []Option `json:"options"`
}

type Option struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	AudioName string  `json:"audioName"`
	Time      int     `json:"time"`
	Delay     float64 `json:"delay"`
	Active    bool    `json:"active"`
}

// For internal use
type SanitizedOption struct {
	Id        int     `json:"id"`
	Text      string  `json:"text"`
	AudioName string  `json:"audioName"`
	Duration  float64 `json:"duration"`
	Delay     float64 `json:"delay"`
}

type OptionTxtFile struct {
	Title string
	Text  string
}
