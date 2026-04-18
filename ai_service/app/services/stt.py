from typing import List
from ..models.schemas import Transcript, Word

def transcribe_audio(audio_path: str) -> Transcript:
    print(f"Mock STT: Transcribing audio from {audio_path}")
    mock_words: List[Word] = [
        Word(text="Привет", start=0.0, end=0.5, confidence=0.9),
        Word(text="меня", start=0.6, end=0.9, confidence=0.85),
        Word(text="зовут", start=1.0, end=1.3, confidence=0.88),
        Word(text="Дмитрий", start=1.4, end=1.9, confidence=0.92),
        Word(text="мой", start=2.0, end=2.2, confidence=0.8),
        Word(text="телефон", start=2.3, end=2.8, confidence=0.9),
        Word(text="89211234567", start=2.9, end=3.8, confidence=0.95),
        Word(text="моя", start=4.0, end=4.2, confidence=0.82),
        Word(text="почта", start=4.3, end=4.7, confidence=0.87),
        Word(text="dmitry.ivanov@example.com", start=4.8, end=6.5, confidence=0.93),
        Word(text="номер", start=6.7, end=7.0, confidence=0.8),
        Word(text="паспорта", start=7.1, end=7.7, confidence=0.9),
        Word(text="1234567890", start=7.8, end=8.8, confidence=0.96),
    ]

    full_text = " ".join([word.text for word in mock_words])
    clean_text = full_text

    return Transcript(
        full_text=full_text,
        clean_text=clean_text,
        words=mock_words,
        language="ru",
        duration=mock_words[-1].end
    )
