import logging
import os
from fastapi import FastAPI, HTTPException
from .models.schemas import AIServiceRequest, AIServiceResponse
from .services import stt, ner, redactor

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Voice Data Redaction AI Service",
    description="A service to perform STT, NER, and audio redaction.",
    version="0.1.0"
)

@app.get("/health")
def health_check():
    return {"status": "ok"}

@app.post("/process", response_model=AIServiceResponse)
async def process_audio(request: AIServiceRequest):
    logger.info(f"Received processing request for job_id: {request.job_id}")
    logger.info(f"File path: {request.file_path}")

    try:
        transcript_result = stt.transcribe_audio(request.file_path)

        pii_results = ner.find_pii(transcript_result.words)

        base_name = os.path.basename(request.file_path)
        file_name, _ = os.path.splitext(base_name)
        output_dir = os.path.join(os.path.dirname(os.path.dirname(request.file_path)), "processed")
        output_file_path = os.path.join(output_dir, f"{file_name}.mp3")
        
        redacted_file_path = redactor.redact_audio(
            input_path=request.file_path,
            pii_to_redact=pii_results,
            output_path=output_file_path
        )

        logger.info(f"Processing for job_id: {request.job_id} completed successfully.")

        return AIServiceResponse(
            transcript=transcript_result,
            pii_found=pii_results,
            output_file=redacted_file_path,
        )

    except Exception as e:
        logger.error(f"Error processing job {request.job_id}: {e}", exc_info=True)
        return AIServiceResponse(
            transcript={},
            pii_found=[],
            output_file="",
            error=f"An internal error occurred: {str(e)}"
        )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("ai_service.app.main:app", host="0.0.0.0", port=5000, reload=True)
