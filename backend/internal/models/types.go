package models

import "time"

type Word struct {
	Text       string  `json:"text"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Confidence float64 `json:"confidence"`
	PIIType    string  `json:"pii_type,omitempty"`
}

type Transcript struct {
	FullText  string  `json:"full_text"`
	CleanText string  `json:"clean_text"`
	Words     []Word  `json:"words"`
	Language  string  `json:"language"`
	Duration  float64 `json:"duration"`
}

type PIIFound struct {
	Type       string  `json:"type"`
	Text       string  `json:"text"`
	AudioStart float64 `json:"audio_start"`
	AudioEnd   float64 `json:"audio_end"`
	Confidence float64 `json:"confidence"`
}

type AIServiceRequest struct {
	FilePath string `json:"file_path"`
	JobID    string `json:"job_id"`
}

type AIServiceResponse struct {
	Transcript Transcript `json:"transcript"`
	PIIFound   []PIIFound `json:"pii_found"`
	OutputFile string     `json:"output_file"`
	Error      string     `json:"error,omitempty"`
}

type ProcessingJob struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	InputFile  string     `json:"input_file"`
	OutputFile string     `json:"output_file,omitempty"`
	Transcript Transcript `json:"transcript,omitempty"`
	PIIFound   []PIIFound `json:"pii_found,omitempty"`
	Error      string     `json:"error,omitempty"`
	Progress   int        `json:"progress"`
	Stage      string     `json:"stage"`
}
