<template>
  <div class="results-section">
    <h2>Результаты обработки</h2>

    <div class="results-grid">
      <div class="result-card">
        <h3>Оригинальная транскрипция</h3>
        <div class="transcript">{{ results.original_transcript || 'Транскрипция недоступна' }}</div>
        <button class="btn" @click="downloadTranscript('original')">
          Скачать транскрипцию
        </button>
      </div>

      <div class="result-card">
        <h3>Анонимизированная транскрипция</h3>
        <div class="transcript">{{ results.redacted_transcript || 'Транскрипция недоступна' }}</div>
        <button class="btn" @click="downloadTranscript('redacted')">
          Скачать транскрипцию
        </button>
      </div>
    </div>

    <div class="results-grid">
      <div class="result-card">
        <h3>Оригинальное аудио</h3>
        <audio
          v-if="results.original_audio_url"
          class="audio-player"
          controls
          :src="getAudioUrl(results.original_audio_url)"
        ></audio>
        <button class="btn" @click="downloadAudio('original')">
          Скачать аудио
        </button>
      </div>

      <div class="result-card">
        <h3>Анонимизированное аудио</h3>
        <audio
          v-if="results.redacted_audio_url"
          class="audio-player"
          controls
          :src="getAudioUrl(results.redacted_audio_url)"
        ></audio>
        <button class="btn btn-primary" @click="downloadAudio('redacted')">
          Скачать аудио
        </button>
      </div>
    </div>

    <button class="btn new-file-btn" @click="$emit('new-file')">
      Загрузить новый файл
    </button>
  </div>
</template>

<script>
import { audioService } from '../services/api'

export default {
  name: 'ResultsView',
  props: {
    results: {
      type: Object,
      required: true
    },
    taskId: {
      type: String,
      required: true
    }
  },
  emits: ['new-file'],
  setup(props) {
    const getAudioUrl = (path) => {
      return audioService.getAudioUrl(path)
    }

    const downloadTranscript = (type) => {
      const text = type === 'original'
        ? props.results.original_transcript
        : props.results.redacted_transcript

      const blob = new Blob([text], { type: 'text/plain' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${type}_transcript_${props.taskId}.txt`
      a.click()
      URL.revokeObjectURL(url)
    }

    const downloadAudio = async (type) => {
      try {
        const blob = await audioService.downloadAudio(props.taskId, type)
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `${type}_audio_${props.taskId}.mp3`
        a.click()
        URL.revokeObjectURL(url)
      } catch (error) {
        console.error('Ошибка скачивания аудио:', error)
      }
    }

    return {
      getAudioUrl,
      downloadTranscript,
      downloadAudio
    }
  }
}
</script>
