package models

import "Clone-TokoOnline/pkg/responses"

type SearchResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TotalResults int  `json:"total_results"`
		Hits         []responses.Product `json:"hits"`
	} `json:"data"`
}