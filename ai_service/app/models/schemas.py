from pydantic import BaseModel
from typing import List, Optional

class Word(BaseModel):
    text: str
    start: float
    end: float
    confidence: float
    pii_type: Optional[str] = None

class Transcript(BaseModel):
    full_text: str
    clean_text: Optional[str] = None
    words: List[Word]
    language: str
    duration: float

class PIIFound(BaseModel):
    type: str
    text: str
    audio_start: float
    audio_end: float
    confidence: float

class AIServiceRequest(BaseModel):
    file_path: str
    job_id: str

class AIServiceResponse(BaseModel):
    transcript: Optional[Transcript]
    pii_found: List[PIIFound]
    output_file: str
    error: Optional[str] = None
