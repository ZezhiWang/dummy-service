package main

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/lib"
	"time"
)

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 1 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080",
		Body:   nil,
		Header: nil,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Bang Bang Bang") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("90th percentile: %s\n", metrics.Latencies.P90)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}