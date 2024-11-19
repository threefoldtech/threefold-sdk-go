package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	prometheus_integration "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/tests/load_testing/prometheus"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	data, err := os.ReadFile("test.yml")
	if err != nil {
		panic(fmt.Sprintf("Failed to read YAML file: %v", err))
	}
	var loadTest LoadTest
	err = yaml.Unmarshal(data, &loadTest)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal YAML: %v", err))
	}

	var targets []vegeta.Target
	for _, endpoint := range loadTest.Endpoints {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    loadTest.BaseUrl + endpoint,
		})
	}

	attacker := vegeta.NewAttacker()

	rate := vegeta.Rate{Freq: loadTest.Freq, Per: time.Second}
	duration := loadTest.Duration * time.Second

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
	pm.AvgLatency.Set(results.Latencies.Mean.Seconds())
	pm.MaxLatency.Set(results.Latencies.Max.Seconds())

	fmt.Printf("Success Rate: %.2f%%\n", successRate)
	fmt.Printf("Error Rate: %.2f%%\n", errorRate)
	fmt.Printf("Average Latency: %v\n", results.Latencies.Mean.Seconds())
	fmt.Printf("Maximum Latency: %v\n", results.Latencies.Max.Seconds())

	time.Sleep(10 * time.Second)
}

type LoadTest struct {
	BaseUrl   string        `yaml:"base_url"`
	Endpoints []string      `yaml:"endpoints"`
	Freq      int           `yaml:"rate_second"`
	Duration  time.Duration `yaml:"duration"`
}
