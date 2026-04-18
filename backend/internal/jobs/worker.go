package jobs

import (
	"log"
	"path/filepath"

	"github.com/pr0100pr0111/KV-redaction/internal/clients"
	"github.com/pr0100pr0111/KV-redaction/internal/models"
)

type JobWorker struct {
	store    *JobStore
	aiClient *clients.AIServiceClient
}

func NewJobWorker(store *JobStore, aiClient *clients.AIServiceClient) *JobWorker {
	return &JobWorker{
		store:    store,
		aiClient: aiClient,
	}
}

func (w *JobWorker) ProcessJob(job *models.ProcessingJob, inputPath string) {
	w.store.Update(job.ID, func(j *models.ProcessingJob) {
		j.Status = "processing"
		j.Stage = "stt"
		j.Progress = 10
	})

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		log.Printf("Error getting absolute path for %s: %s", inputPath, err.Error())
		w.store.Update(job.ID, func(j *models.ProcessingJob) {
			j.Status = "failed"
			j.Error = "Could not determine absolute file path"
		})
		return
	}

	req := models.AIServiceRequest{
		FilePath: absPath,
		JobID:    job.ID,
	}

	result, err := w.aiClient.ProcessAudio(req)
	if err != nil || result.Error != "" {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = result.Error
		}
		log.Printf("Error processing job %s: %s", job.ID, errMsg)
		w.store.Update(job.ID, func(j *models.ProcessingJob) {
			j.Status = "failed"
			j.Error = errMsg
		})
		return
	}

	w.store.Update(job.ID, func(j *models.ProcessingJob) {
		j.Transcript = result.Transcript
		j.PIIFound = result.PIIFound
		j.OutputFile = result.OutputFile
		j.Progress = 100
		j.Stage = "completed"
		j.Status = "completed"
	})
	log.Printf("Job %s completed successfully", job.ID)
}
