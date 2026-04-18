package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pr0100pr0111/KV-redaction/config"
	"github.com/pr0100pr0111/KV-redaction/internal/api"
	"github.com/pr0100pr0111/KV-redaction/internal/clients"
	"github.com/pr0100pr0111/KV-redaction/internal/jobs"
)

func main() {
	if err := os.MkdirAll("storage/uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	if err := os.MkdirAll("storage/processed", 0755); err != nil {
		log.Fatalf("Failed to create processed directory: %v", err)
	}

	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	aiClient := clients.NewAIServiceClient(cfg.AIServiceURL)
	if err := aiClient.Health(); err != nil {
		log.Printf("Warning: AI service not available: %v", err)
	}

	jobStore := jobs.NewJobStore()
	jobWorker := jobs.NewJobWorker(jobStore, aiClient)
	handler := api.NewHandler(jobStore, jobWorker)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(api.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/upload", handler.HandleUpload)
		v1.GET("/job/:id", handler.HandleGetJobStatus)
		v1.GET("/jobs", handler.HandleGetJobs)
		v1.GET("/download/:id/:type", handler.HandleDownload)
		v1.DELETE("/job/:id", handler.HandleDeleteJob)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Backend starting on port %s", cfg.Port)
		log.Printf("🤖 AI Service: %s", cfg.AIServiceURL)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
