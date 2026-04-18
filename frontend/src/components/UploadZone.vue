<template>
  <div
    class="upload-zone"
    :class="{ dragover: isDragging }"
    @click="triggerFileInput"
    @dragover.prevent="isDragging = true"
    @dragleave.prevent="isDragging = false"
    @drop.prevent="handleDrop"
  >
    <svg class="upload-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/>
    </svg>
    <p class="upload-text">Перетащите аудиофайл сюда или нажмите для выбора</p>
    <p class="upload-hint">Поддерживаются форматы: MP3, WAV, OGG, M4A</p>
    <input
      ref="fileInput"
      type="file"
      accept="audio/*"
      style="display: none"
      @change="handleFileInput"
    >
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'UploadZone',
  emits: ['file-selected'],
  setup(props, { emit }) {
    const isDragging = ref(false)
    const fileInput = ref(null)

    const triggerFileInput = () => {
      fileInput.value.click()
    }

    const handleFileInput = (event) => {
      const file = event.target.files[0]
      if (file && file.type.startsWith('audio/')) {
        emit('file-selected', file)
      }
    }

    const handleDrop = (event) => {
      isDragging.value = false
      const file = event.dataTransfer.files[0]
      if (file && file.type.startsWith('audio/')) {
        emit('file-selected', file)
      }
    }

    return {
      isDragging,
      fileInput,
      triggerFileInput,
      handleFileInput,
      handleDrop
    }
  }
}
</script>
