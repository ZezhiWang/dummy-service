package main

import (
	"encoding/json"
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
	partialMethod := "POST"
	compMethod := "DELETE"
	return generateBody(services, val, partialMethod, compMethod)
}

func generateDeleteBody(services []string, val string) []byte {
	partialMethod := "DELETE"
	compMethod := "POST"
	return generateBody(services, val, partialMethod, compMethod)
}

func generateBody(services []string, val, partial, comp string) []byte {
	requests := make(map[string]TransactionReq)
	for _, svc := range services {
		url := "http://" + svc + ".default.svc.cluster.local:8888/" + val
		requests[svc] = TransactionReq{
			PartialReq: Request{
				Method: partial,
				URL:    url,
				Body:   "",
			},
			CompReq:    Request{
				Method: comp,
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