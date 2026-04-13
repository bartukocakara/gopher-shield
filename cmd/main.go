package main

import (
    "log"
    "net/http"
    "os" // Bunu ekle

    "github.com/bartukocakara/gopher-shield/internal/metrics"
    "github.com/bartukocakara/gopher-shield/internal/proxy"
    "github.com/bartukocakara/gopher-shield/internal/resilience"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // 1. Setup Metrics
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Println("Metrics exporting on :9090")
        log.Fatal(http.ListenAndServe(":9090", nil))
    }()

    // 2. Read Upstream URL from Env (Docker-compose'dan gelen değer)
    targetURL := os.Getenv("UPSTREAM_URL")
    if targetURL == "" {
        targetURL = "http://order-service:3000" // Fallback (Yedek)
    }

    // 3. Setup Breaker with Metric Hook
    cb := resilience.NewCircuitBreaker(3)
    cb.OnStateChange = func(state string) {
        metrics.CBTransitions.WithLabelValues(state).Inc()
        log.Printf("STATE CHANGE: Circuit is now %s", state)
    }

    // 4. Start Proxy using the dynamic targetURL
    shield, _ := proxy.NewShieldProxy(targetURL, cb)
    log.Printf("Proxy listening on :8080 (Targeting: %s)", targetURL)
    log.Fatal(http.ListenAndServe(":8080", shield))
}