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
        :add-log="addLog"
        :logs="logs"
      />
    </main>
  </div>
</template>

<script>
import { ref } from 'vue'
import UploadZone from './components/UploadZone.vue'
import ProcessingView from './components/ProcessingView.vue'
import ResultsView from './components/ResultsView.vue'
import { jobService } from './services/api'

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

        const uploadResponse = await jobService.uploadFile(file, (event) => {
          progress.value = Math.min(90, Math.round((10 + event.loaded * 80) / event.total))
        })
        taskId.value = uploadResponse.job_id

        addLog(`Файл загружен. ID задачи: ${taskId.value}`, 'success')
        progress.value = 40
        status.value = 'Файл загружен, ожидание обработки'

        startPolling()
      } catch (error) {
        addLog(`Ошибка: ${error.message}`, 'error')
        status.value = 'Ошибка загрузки'
      }
    }

    const startPolling = () => {
      pollInterval = setInterval(async () => {
        try {
          const jobData = await jobService.getStatus(taskId.value)

          if (jobData.logs && Array.isArray(jobData.logs)) {
              jobData.logs.forEach(newLog => {
                  if (!logs.value.some(existingLog => existingLog.message === newLog.message)) {
                      addLog(newLog.message, newLog.level);
                  }
              });
          }

          if (jobData.stage === 'stt') {
              progress.value = 50;
              status.value = 'Распознавание речи...';
          } else if (jobData.stage === 'ner') {
              progress.value = 70;
              status.value = 'Поиск персональных данных...';
          } else if (jobData.stage === 'redaction') {
              progress.value = 90;
              status.value = 'Анонимизация аудио...';
          } else if (jobData.progress !== undefined) {
              progress.value = 40 + (jobData.progress * 0.5);
          }

          if (jobData.status === 'completed') {
            clearInterval(pollInterval)
            addLog('Обработка завершена успешно', 'success')
            progress.value = 100
            status.value = 'Готово'
            await loadResults(jobData)
          } else if (jobData.status === 'failed') {
            clearInterval(pollInterval)
            addLog(`Обработка завершилась с ошибкой: ${jobData.error}`, 'error')
            status.value = 'Ошибка обработки'
          }
        } catch (error) {
          clearInterval(pollInterval)
          addLog(`Ошибка получения статуса: ${error.message}`, 'error')
          status.value = 'Ошибка'
        }
      }, 2000)
    }

    const loadResults = async (jobData) => {
      try {
        results.value = jobData
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
        pollInterval = null
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
