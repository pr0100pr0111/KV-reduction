<template>
  <div class="container">
    <header class="header">
      <h1>KV Redaction</h1>
      <p>Анонимизация персональных данных в аудиосообщениях</p>
    </header>

    <main>
      <UploadZone
        v-if="currentView === 'upload'"
        @file-selected="handleFileSelected"
      />

      <ProcessingView
        v-if="currentView === 'processing'"
        :file="selectedFile"
        :progress="progress"
        :status="status"
        :logs="logs"
      />

      <ResultsView
        v-if="currentView === 'results'"
        :results="results"
        :task-id="taskId"
        @new-file="resetApp"
      />
    </main>
  </div>
</template>

<script>
import { ref } from 'vue'
import UploadZone from './components/UploadZone.vue'
import ProcessingView from './components/ProcessingView.vue'
import ResultsView from './components/ResultsView.vue'
import { audioService } from './services/api'

export default {
  name: 'App',
  components: {
    UploadZone,
    ProcessingView,
    ResultsView
  },
  setup() {
    const currentView = ref('upload')
    const selectedFile = ref(null)
    const taskId = ref(null)
    const progress = ref(0)
    const status = ref('')
    const logs = ref([])
    const results = ref(null)
    let pollInterval = null

    const addLog = (message, level = 'info') => {
      logs.value.push({
        time: new Date().toLocaleTimeString('ru-RU'),
        message,
        level
      })
    }

    const handleFileSelected = async (file) => {
      selectedFile.value = file
      currentView.value = 'processing'

      try {
        addLog('Загрузка файла на сервер...', 'info')
        progress.value = 10
        status.value = 'Загрузка...'

        const uploadResponse = await audioService.uploadFile(file)
        taskId.value = uploadResponse.task_id

        addLog(`Файл загружен. ID задачи: ${taskId.value}`, 'success')
        progress.value = 30
        status.value = 'Файл загружен'

        addLog('Запуск обработки...', 'info')
        await audioService.startProcessing(taskId.value)

        addLog('Обработка запущена', 'success')
        progress.value = 40
        status.value = 'Обработка...'

        startPolling()
      } catch (error) {
        addLog(`Ошибка: ${error.message}`, 'error')
        status.value = 'Ошибка'
      }
    }

    const startPolling = () => {
      pollInterval = setInterval(async () => {
        try {
          const statusData = await audioService.getStatus(taskId.value)

          if (statusData.logs && statusData.logs.length > 0) {
            statusData.logs.forEach(log => {
              const exists = logs.value.some(l => l.message === log.message)
              if (!exists) {
                addLog(log.message, log.level || 'info')
              }
            })
          }

          if (statusData.progress !== undefined) {
            progress.value = 40 + (statusData.progress * 0.5)
          }

          if (statusData.status) {
            status.value = statusData.status
          }

          if (statusData.status === 'completed') {
            clearInterval(pollInterval)
            addLog('Обработка завершена успешно', 'success')
            progress.value = 100
            status.value = 'Готово'
            await loadResults()
          } else if (statusData.status === 'failed') {
            clearInterval(pollInterval)
            addLog('Обработка завершилась с ошибкой', 'error')
            status.value = 'Ошибка'
          }
        } catch (error) {
          clearInterval(pollInterval)
          addLog(`Ошибка получения статуса: ${error.message}`, 'error')
        }
      }, 1000)
    }

    const loadResults = async () => {
      try {
        const resultsData = await audioService.getResults(taskId.value)
        results.value = resultsData
        currentView.value = 'results'
        addLog('Результаты загружены', 'success')
      } catch (error) {
        addLog(`Ошибка загрузки результатов: ${error.message}`, 'error')
      }
    }

    const resetApp = () => {
      currentView.value = 'upload'
      selectedFile.value = null
      taskId.value = null
      progress.value = 0
      status.value = ''
      logs.value = []
      results.value = null
      if (pollInterval) {
        clearInterval(pollInterval)
      }
    }

    return {
      currentView,
      selectedFile,
      taskId,
      progress,
      status,
      logs,
      results,
      handleFileSelected,
      resetApp
    }
  }
}
</script>
