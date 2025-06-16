package main

import (
    "assignment/pkg/config"
    "assignment/pkg/health"
    "assignment/pkg/logger"
    "assignment/pkg/metrics"
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Initialize logger
    logConfig := logger.LogConfig{
        Level:      getEnvOrDefault("LOG_LEVEL", "info"),
        Format:     getEnvOrDefault("LOG_FORMAT", "json"),
        OutputPath: getEnvOrDefault("LOG_PATH", ""),
    }
    if err := logger.InitLogger(logConfig); err != nil {
        panic(err)
    }

    // Load and validate configuration
    cfg, err := config.LoadConfig(getEnvOrDefault("CONFIG_FILE", "config/app_configuration.json"))
    if err != nil {
        logger.Fatal("Failed to load configuration", logger.Fields{"error": err})
    }
    
    if err := cfg.ValidateConfig(); err != nil {
        logger.Fatal("Invalid configuration", logger.Fields{"error": err})
    }

    // Initialize health checker
    healthChecker := health.NewHealthChecker(30 * time.Second)
    
    // Add health checks
    healthChecker.AddCheck("config", func() error {
        if err := cfg.ValidateConfig(); err != nil {
            return err
        }
        return nil
    })
    
    // Start health check server
    go func() {
        http.Handle("/health", healthChecker)
        if err := http.ListenAndServe(":8080", nil); err != nil {
            logger.Error("Health check server failed", logger.Fields{"error": err})
        }
    }()

    // Start metrics server
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        if err := http.ListenAndServe(":8081", nil); err != nil {
            logger.Error("Metrics server failed", logger.Fields{"error": err})
        }
    }()

    // Setup graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Start your application processing here
    // ... (your existing processing logic) ...

    // Wait for shutdown signal
    <-ctx.Done()
    logger.Info("Shutting down gracefully", logger.Fields{})
    
    // Add cleanup logic here
    // ... (your cleanup logic) ...
}

func getEnvOrDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
} 