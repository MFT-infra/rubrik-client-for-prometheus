package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/MFT-infra/rubrik-client-for-prometheus/src/golang/jobs"
	"github.com/MFT-infra/rubrik-client-for-prometheus/src/golang/livemount"
	"github.com/MFT-infra/rubrik-client-for-prometheus/src/golang/objectprotection"
	"github.com/MFT-infra/rubrik-client-for-prometheus/src/golang/stats"
)

func main() {
	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	rubrik, err := connectRubrik(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Rubrik: %v", err)
	}

	clusterDetails, err := rubrik.Get("v1", "/cluster/me", 60)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster details: %v", err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"].(string)
	log.Printf("Connected to Rubrik cluster: %s", clusterName)

	// Background metric collectors
	go func() {
		for {
			stats.GetStorageSummaryStats(rubrik, clusterName)
			stats.GetRunwayRemaining(rubrik, clusterName)
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		for {
			stats.GetNodeStats(rubrik, clusterName)
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		for {
			stats.Get24HJobStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			stats.GetSlaComplianceStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			jobs.GetMssqlFailedJobs(rubrik, clusterName)
			jobs.GetVmwareVmFailedJobs(rubrik, clusterName)
			time.Sleep(5 * time.Minute)
		}
	}()

	go func() {
		for {
			stats.GetMssqlCapacityStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			stats.GetOracleCapacityStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			stats.GetVSphereVmCapacityStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			objectprotection.GetSnappableEffectiveSlaDomain(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			objectprotection.GetSlaDomainSummary(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			livemount.GetMssqlLiveMountAges(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			stats.GetRelicStorageStats(rubrik, clusterName)
			time.Sleep(1 * time.Hour)
		}
	}()

	// Start Prometheus metrics HTTP server
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Prometheus exporter on port %s", cfg.PrometheusPort)
	log.Fatal(http.ListenAndServe(":"+cfg.PrometheusPort, nil))
}
