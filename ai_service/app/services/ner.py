import re
from typing import List
from ..models.schemas import Word, PIIFound

def find_pii(words: List[Word]) -> List[PIIFound]:
    print("Mock NER: Finding PII in the transcript.")
    pii_found: List[PIIFound] = []

    full_text = "".join([word.text for word in words])

    pii_patterns = {
        "passport": r"\d{4}\s?\d{6}",
        "phone": r"(\+7|8)\s?\(?\d{3}\)?\s?\d{3}\s?\d{2}\s?\d{2}",
        "email": r"[\w.-]+@[\w.-]+\.w+",
        "inn": r"\d{10}|\d{12}",
        "snils": r"\d{3}-\d{3}-\d{3}\s?\d{2}"
    }

    for word in words:
        for pii_type, pattern in pii_patterns.items():
            if re.search(pattern, word.text):
                pii_found.append(PIIFound(
                    type=pii_type,
                    text=word.text,
                    audio_start=word.start,
                    audio_end=word.end,
                    confidence=word.confidence
                ))
    print(f"Mock NER: Found {len(pii_found)} PII items.")
    return pii_found
