package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	prometheus_integration "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/tests/load_testing/prometheus"
	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	baseUrl := "http://localhost:8080"
	targets := []vegeta.Target{
		{Method: "GET", URL: baseUrl + "/contracts"},
		{Method: "GET", URL: baseUrl + "/farms"},
		{Method: "GET", URL: baseUrl + "/gateways"},
		{Method: "GET", URL: baseUrl + "/nodes"},
		{Method: "GET", URL: baseUrl + "/public_ips"},
		{Method: "GET", URL: baseUrl + "/stats"},
		{Method: "GET", URL: baseUrl + "/twins"},
	}

	attacker := vegeta.NewAttacker()

	rate := vegeta.Rate{Freq: 5000, Per: time.Second}
	duration := 10 * time.Second

	reg := prometheus.NewRegistry()
	pm := prometheus_integration.NewMetrics()

	if err := pm.Register(reg); err != nil {
		log.Fatal("error registering metrics", err)
	}

	go func() {
		fmt.Println("Prometheus metrics server running on :9090/metrics")
		log.Fatal(http.ListenAndServe(":9090", prometheus_integration.NewHandler(reg, time.Now().UTC())))
	}()

	var results vegeta.Metrics
	for res := range attacker.Attack(vegeta.NewStaticTargeter(targets...), rate, duration, "load testing") {
		results.Add(res)
		pm.Observe(res)
	}

	results.Close()

	successRate := results.Success * 100
	errorRate := (1.0 - results.Success) * 100
	pm.SuccessRate.Set(successRate)
	pm.ErrorRate.Set(errorRate)

	fmt.Printf("Success Rate: %.2f%%\n", successRate)
	fmt.Printf("Error Rate: %.2f%%\n", errorRate)

	time.Sleep(10 * time.Second)

}
