package main

import (
	"encoding/json"
	"os"
)

type Transaction struct {
	Tiers map[int]map[string]TransactionReq `json:"tier"`
}

type TransactionReq struct {
	PartialReq Request `json:"partial_req"`
	CompReq    Request `json:"comp_req"`
}

type Request struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

func generatePostBody(services []string, val string) []byte {
	requests := make(map[string]TransactionReq)
	for _, svc := range services {
		url := os.Getenv(svc) + "/base/" + val
		requests[svc] = TransactionReq{
			PartialReq: Request{
				Method: "POST",
				URL:    url,
				Body:   "",
			},
			CompReq:    Request{
				Method: "DELETE",
				URL:    url,
				Body:   "",
			},
		}
	}

	body, _ := json.Marshal(Transaction{Tiers: map[int]map[string]TransactionReq{
		0: requests,
	}})

	return body
}