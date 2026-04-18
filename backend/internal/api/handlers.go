package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pr0100pr0111/KV-redaction/internal/jobs"
	"github.com/pr0100pr0111/KV-redaction/internal/models"
)

type Handler struct {
	store  *jobs.JobStore
	worker *jobs.JobWorker
}

func NewHandler(store *jobs.JobStore, worker *jobs.JobWorker) *Handler {
	return &Handler{
		store:  store,
		worker: worker,
	}
}

func (h *Handler) HandleUpload(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	ext := filepath.Ext(file.Filename)
	validExts := []string{".wav", ".mp3", ".ogg", ".flac"}
	valid := false
	for _, v := range validExts {
		if ext == v {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	job := h.store.Create(file.Filename)

	uploadPath := filepath.Join("storage", "uploads", job.ID+ext)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		h.store.Update(job.ID, func(j *models.ProcessingJob) {
			j.Status = "failed"
			j.Error = err.Error()
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Save failed"})
		return
	}

	go h.worker.ProcessJob(job, uploadPath)

	c.JSON(http.StatusOK, gin.H{
		"job_id":  job.ID,
		"status":  "processing",
		"message": "File uploaded successfully",
	})
}

func (h *Handler) HandleGetJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	job, exists := h.store.Get(jobID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) HandleGetJobs(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetAll())
}

func (h *Handler) HandleDeleteJob(c *gin.Context) {
	jobID := c.Param("id")

	job, exists := h.store.Get(jobID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if job.InputFile != "" {
		ext := filepath.Ext(job.InputFile)
		inputPath := filepath.Join("storage", "uploads", job.ID+ext)
		os.Remove(inputPath)
	}
	if job.OutputFile != "" {
		os.Remove(job.OutputFile)
	}

	h.store.Delete(jobID)

	c.JSON(http.StatusOK, gin.H{"message": "Job " + jobID + " deleted"})
}

func (h *Handler) HandleDownload(c *gin.Context) {
	jobID := c.Param("id")
	fileType := c.Param("type")

	job, exists := h.store.Get(jobID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if job.Status != "completed" {
		c.JSON(http.StatusConflict, gin.H{"error": "Job is not completed yet"})
		return
	}

	var filePath string
	var fileName string

	switch fileType {
	case "audio":
		filePath = job.OutputFile
		fileName = "redacted_" + job.InputFile
	case "original_audio":
		ext := filepath.Ext(job.InputFile)
		filePath = filepath.Join("storage", "uploads", job.ID+ext)
		fileName = "original_" + job.InputFile
	case "transcript":
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Transcript download not implemented"})
		return
	case "log":
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Log download not implemented"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type: " + fileType})
		return
	}

	c.FileAttachment(filePath, fileName)
}
