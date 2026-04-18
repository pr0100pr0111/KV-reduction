<template>
  <div class="results-section">
    <h2>Результаты обработки</h2>

    <div class="results-grid">
      <div class="result-card">
        <h3>Оригинальная транскрипция</h3>
        <div class="transcript">{{ results.transcript?.full_text || 'Транскрипция недоступна' }}</div>
        <button class="btn" @click="downloadTranscript('original')">
          Скачать транскрипцию
        </button>
      </div>

      <div class="result-card">
        <h3>Анонимизированная транскрипция</h3>
        <div class="transcript">{{ results.transcript?.clean_text || 'Транскрипция недоступна' }}</div>
        <button class="btn" @click="downloadTranscript('redacted')">
          Скачать транскрипцию
        </button>
      </div>
    </div>

    <div class="results-grid">
      <div class="result-card">
        <h3>Оригинальное аудио</h3>
        <audio
          v-if="results.input_file"
          class="audio-player"
          controls
          :src="getDownloadUrl(taskId, 'original_audio', results.input_file)"
        ></audio>
        <button class="btn" @click="downloadAudio('original_audio')">
          Скачать аудио
        </button>
      </div>

      <div class="result-card">
        <h3>Анонимизированное аудио</h3>
        <audio
          v-if="results.output_file"
          class="audio-player"
          controls
          :src="getDownloadUrl(taskId, 'audio', results.output_file)"
        ></audio>
        <button class="btn btn-primary" @click="downloadAudio('audio')">
          Скачать аудио
        </button>
      </div>
    </div>

    <div class="logs-container">
      <h3 class="logs-header">Логи</h3>
      <div class="logs" ref="logsContainer">
        <div
          v-for="(log, index) in logs"
          :key="index"
          class="log-entry"
          :class="log.level"
        >
          <span class="log-time">{{ log.time }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </div>

    <button class="btn new-file-btn" @click="$emit('new-file')">
      Загрузить новый файл
    </button>
  </div>
</template>

<script>
import { ref, watch, nextTick } from 'vue'
import { jobService, getDownloadUrl } from '../services/api'

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
    },
    addLog: {
      type: Function,
      required: true
    },
    logs: {
      type: Array,
      required: true
    }
  },
  emits: ['new-file'],
  setup(props) {
    const logsContainer = ref(null)

    watch(() => props.logs.length, async () => {
      await nextTick()
      if (logsContainer.value) {
        logsContainer.value.scrollTop = logsContainer.value.scrollHeight
      }
    })

    const downloadTranscript = (type) => {
      const text = type === 'original'
        ? props.results.transcript?.full_text
        : props.results.transcript?.clean_text

      if (!text) {
        props.addLog(`No ${type} transcript available for download.`, 'warn')
        return;
      }

      const blob = new Blob([text], { type: 'text/plain' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${type}_transcript_${props.taskId}.txt`
      a.click()
      URL.revokeObjectURL(url)
    }

    const downloadAudio = async (type) => {
      props.addLog(`Попытка скачивания ${type}...`, 'info')
      try {
        const blob = await jobService.downloadBlob(props.taskId, type)
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        const fileName = type === 'audio' ? props.results.output_file : props.results.input_file;
        a.download = `${type}_${props.taskId}_${fileName}`
        a.click()
        URL.revokeObjectURL(url)
        props.addLog(`Аудио "${fileName}" успешно скачано.`, 'success')
      } catch (error) {
        console.error('Ошибка скачивания аудио:', error)
        props.addLog(`Ошибка скачивания аудио: ${error.message}`, 'error')
      }
    }

    return {
      getDownloadUrl, // Make it available in the template
      downloadTranscript,
      downloadAudio,
      logsContainer
    }
  }
}
</script>
