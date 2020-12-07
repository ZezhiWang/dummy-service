package main

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://35.224.223.252:8888/1",
		Body:   nil,
		Header: nil,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Bang Bang Bang") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Max: %s\n", metrics.Latencies.Max)
	fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
