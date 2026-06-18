package main

type Organization struct {
	Name        string            `json:"name"`
	URL         map[string]string `json:"url"`
	Description map[string]string `json:"description"`
}

type StreetResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type ZipResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}
