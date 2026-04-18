<template>
  <div class="processing-section">
    <div class="file-info">
      <svg class="file-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M9 18V5l12-2v13M9 13c0 1.66-1.34 3-3 3s-3-1.34-3-3 1.34-3 3-3 3 1.34 3 3z"/>
      </svg>
      <div class="file-details">
        <p class="file-name">{{ file.name }}</p>
        <p class="file-size">{{ formatFileSize(file.size) }}</p>
      </div>
    </div>

    <div class="progress-container">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="progress-text">{{ status }}</p>
    </div>

    <div class="logs-container">
      <h3 class="logs-header">Логи обработки</h3>
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
  </div>
</template>

<script>
import { ref, watch, nextTick } from 'vue'

export default {
  name: 'ProcessingView',
  props: {
    file: {
      type: File,
      required: true
    },
    progress: {
      type: Number,
      required: true
    },
    status: {
      type: String,
      required: true
    },
    logs: {
      type: Array,
      required: true
    }
  },
  setup(props) {
    const logsContainer = ref(null)

    const formatFileSize = (bytes) => {
      if (bytes < 1024) return bytes + ' B'
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
      return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
    }

    watch(() => props.logs.length, async () => {
      await nextTick()
      if (logsContainer.value) {
        logsContainer.value.scrollTop = logsContainer.value.scrollHeight
      }
    })

    return {
      logsContainer,
      formatFileSize
    }
  }
}
</script>
