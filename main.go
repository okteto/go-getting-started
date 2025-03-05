package main

import (
        "fmt"
        "net/http"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        // Contador para las llamadas al endpoint de health
        healthCheckCounter = promauto.NewCounter(prometheus.CounterOpts{
                Name: "hello_world_health_checks_total",
                Help: "El número total de veces que se ha llamado al endpoint de health",
        })
)

func main() {
        fmt.Println("Starting hello-world server...")
        
        // Configurar rutas
        http.HandleFunc("/", helloServer)
        http.HandleFunc("/health", healthServer)
        
        // Exponer métricas de Prometheus en /metrics
        http.Handle("/metrics", promhttp.Handler())
        
        fmt.Println("Server listening on :8080")
        fmt.Println("Metrics available at /metrics")
        fmt.Println("Health endpoint at /health")
        
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        message := `
        ¡Desarrollar con Okteto es increíble! 🚀
        
        Hot reloading, entornos de desarrollo en la nube,
        y despliegues instantáneos hacen que el desarrollo
        sea más rápido y divertido que nunca.
        
        ¡Okteto mola muchísimo! 😎
        `
        fmt.Fprint(w, message)
}

func healthServer(w http.ResponseWriter, r *http.Request) {
        // Incrementar el contador cada vez que se llama al endpoint de health
        healthCheckCounter.Inc()
        
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
}
