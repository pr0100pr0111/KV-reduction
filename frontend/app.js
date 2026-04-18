const API_BASE_URL = 'http://localhost:8000';

class AudioRedactionApp {
    constructor() {
        this.uploadZone = document.getElementById('upload-zone');
        this.fileInput = document.getElementById('file-input');
        this.processingZone = document.getElementById('processing-zone');
        this.resultsZone = document.getElementById('results-zone');
        this.logsContainer = document.getElementById('logs');

        this.currentFile = null;
        this.taskId = null;

        this.initEventListeners();
    }

    initEventListeners() {
        this.uploadZone.addEventListener('click', () => this.fileInput.click());
        this.fileInput.addEventListener('change', (e) => this.handleFileSelect(e.target.files[0]));

        this.uploadZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            this.uploadZone.classList.add('dragover');
        });

        this.uploadZone.addEventListener('dragleave', () => {
            this.uploadZone.classList.remove('dragover');
        });

        this.uploadZone.addEventListener('drop', (e) => {
            e.preventDefault();
            this.uploadZone.classList.remove('dragover');
            const file = e.dataTransfer.files[0];
            if (file && file.type.startsWith('audio/')) {
                this.handleFileSelect(file);
            } else {
                this.addLog('Пожалуйста, загрузите аудиофайл', 'error');
            }
        });

        document.getElementById('new-file-btn').addEventListener('click', () => this.reset());

        document.getElementById('download-original-transcript').addEventListener('click', () =>
            this.downloadTranscript('original'));
        document.getElementById('download-redacted-transcript').addEventListener('click', () =>
            this.downloadTranscript('redacted'));
        document.getElementById('download-original-audio').addEventListener('click', () =>
            this.downloadAudio('original'));
        document.getElementById('download-redacted-audio').addEventListener('click', () =>
            this.downloadAudio('redacted'));
    }

    handleFileSelect(file) {
        if (!file) return;

        this.currentFile = file;
        this.showProcessing();
        this.displayFileInfo(file);
        this.uploadFile(file);
    }

    showProcessing() {
        this.uploadZone.classList.add('hidden');
        this.processingZone.classList.remove('hidden');
        this.resultsZone.classList.add('hidden');
    }

    showResults() {
        this.processingZone.classList.add('hidden');
        this.resultsZone.classList.remove('hidden');
    }

    reset() {
        this.uploadZone.classList.remove('hidden');
        this.processingZone.classList.add('hidden');
        this.resultsZone.classList.add('hidden');
        this.fileInput.value = '';
        this.logsContainer.innerHTML = '';
        this.currentFile = null;
        this.taskId = null;
        this.updateProgress(0, 'Готов к загрузке');
    }

    displayFileInfo(file) {
        document.getElementById('file-name').textContent = file.name;
        document.getElementById('file-size').textContent = this.formatFileSize(file.size);
    }

    formatFileSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
        return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
    }

    updateProgress(percent, text) {
        document.getElementById('progress-fill').style.width = percent + '%';
        document.getElementById('progress-text').textContent = text;
    }

    addLog(message, type = 'info') {
        const logEntry = document.createElement('div');
        logEntry.className = `log-entry ${type}`;
        const timestamp = new Date().toLocaleTimeString('ru-RU');
        logEntry.textContent = `[${timestamp}] ${message}`;
        this.logsContainer.appendChild(logEntry);
        this.logsContainer.scrollTop = this.logsContainer.scrollHeight;
    }

    async uploadFile(file) {
        try {
            this.addLog('Начало загрузки файла...', 'info');
            this.updateProgress(10, 'Загрузка файла...');

            const formData = new FormData();
            formData.append('file', file);

            const response = await fetch(`${API_BASE_URL}/upload`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error(`Ошибка загрузки: ${response.status}`);
            }

            const data = await response.json();
            this.taskId = data.task_id;

            this.addLog(`Файл загружен. ID задачи: ${this.taskId}`, 'success');
            this.updateProgress(30, 'Файл загружен');

            await this.startProcessing();
        } catch (error) {
            this.addLog(`Ошибка: ${error.message}`, 'error');
            this.updateProgress(0, 'Ошибка загрузки');
        }
    }

    async startProcessing() {
        try {
            this.addLog('Запуск обработки...', 'info');
            this.updateProgress(40, 'Обработка аудио...');

            const response = await fetch(`${API_BASE_URL}/process/${this.taskId}`, {
                method: 'POST'
            });

            if (!response.ok) {
                throw new Error(`Ошибка обработки: ${response.status}`);
            }

            this.addLog('Обработка запущена', 'success');
            this.pollStatus();
        } catch (error) {
            this.addLog(`Ошибка: ${error.message}`, 'error');
            this.updateProgress(0, 'Ошибка обработки');
        }
    }

    async pollStatus() {
        const pollInterval = setInterval(async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/status/${this.taskId}`);

                if (!response.ok) {
                    throw new Error(`Ошибка получения статуса: ${response.status}`);
                }

                const data = await response.json();

                if (data.logs && data.logs.length > 0) {
                    data.logs.forEach(log => {
                        if (!this.logsContainer.textContent.includes(log.message)) {
                            this.addLog(log.message, log.level || 'info');
                        }
                    });
                }

                if (data.progress !== undefined) {
                    this.updateProgress(40 + (data.progress * 0.5), data.status || 'Обработка...');
                }

                if (data.status === 'completed') {
                    clearInterval(pollInterval);
                    this.addLog('Обработка завершена успешно', 'success');
                    this.updateProgress(100, 'Готово');
                    await this.loadResults();
                } else if (data.status === 'failed') {
                    clearInterval(pollInterval);
                    this.addLog('Обработка завершилась с ошибкой', 'error');
                    this.updateProgress(0, 'Ошибка');
                }
            } catch (error) {
                clearInterval(pollInterval);
                this.addLog(`Ошибка: ${error.message}`, 'error');
            }
        }, 1000);
    }

    async loadResults() {
        try {
            const response = await fetch(`${API_BASE_URL}/results/${this.taskId}`);

            if (!response.ok) {
                throw new Error(`Ошибка получения результатов: ${response.status}`);
            }

            const data = await response.json();

            document.getElementById('original-transcript').textContent =
                data.original_transcript || 'Транскрипция недоступна';
            document.getElementById('redacted-transcript').textContent =
                data.redacted_transcript || 'Транскрипция недоступна';

            if (data.original_audio_url) {
                document.getElementById('original-audio').src =
                    `${API_BASE_URL}${data.original_audio_url}`;
            }

            if (data.redacted_audio_url) {
                document.getElementById('redacted-audio').src =
                    `${API_BASE_URL}${data.redacted_audio_url}`;
            }

            this.showResults();
            this.addLog('Результаты загружены', 'success');
        } catch (error) {
            this.addLog(`Ошибка загрузки результатов: ${error.message}`, 'error');
        }
    }

    downloadTranscript(type) {
        const transcriptElement = document.getElementById(`${type}-transcript`);
        const text = transcriptElement.textContent;
        const blob = new Blob([text], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${type}_transcript_${this.taskId}.txt`;
        a.click();
        URL.revokeObjectURL(url);
        this.addLog(`Транскрипция (${type}) скачана`, 'success');
    }

    async downloadAudio(type) {
        try {
            const response = await fetch(`${API_BASE_URL}/download/${this.taskId}/${type}`);

            if (!response.ok) {
                throw new Error(`Ошибка скачивания: ${response.status}`);
            }

            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${type}_audio_${this.taskId}.${this.getFileExtension(this.currentFile.name)}`;
            a.click();
            URL.revokeObjectURL(url);
            this.addLog(`Аудио (${type}) скачано`, 'success');
        } catch (error) {
            this.addLog(`Ошибка скачивания аудио: ${error.message}`, 'error');
        }
    }

    getFileExtension(filename) {
        return filename.split('.').pop() || 'mp3';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new AudioRedactionApp();
});
