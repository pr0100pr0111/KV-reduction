# KV Redaction - Frontend

## стек

- **Vue 3**
- **Vite** - сборщик и dev-сервер
- **Axios** - HTTP клиент для API запросов


## Структура фронта

```
frontend/
├── src/
│   ├── components/
│   │   ├── UploadZone.vue       # Компонент загрузки файлов
│   │   ├── ProcessingView.vue   # Компонент отображения процесса обработки
│   │   └── ResultsView.vue      # Компонент отображения результатов
│   ├── services/
│   │   └── api.js               # API сервис для взаимодействия с бэкендом
│   ├── App.vue                  # Главный компонент приложения
│   ├── main.js                  # Точка входа
│   └── style.css                # Глобальные стили
├── index.html                   # HTML шаблон
├── vite.config.js               # Конфигурация Vite
├── package.json                 # Зависимости проекта
└── README.md                    # Документация
```

## Компоненты

### 1. App.vue (Главный компонент)
Управляет состоянием всего приложения и координирует взаимодействие между компонентами.

**Состояния:**
- `currentView` - текущий экран ('upload' | 'processing' | 'results')
- `selectedFile` - выбранный аудиофайл
- `taskId` - ID задачи обработки
- `progress` - прогресс обработки (0-100)
- `status` - текстовый статус обработки
- `logs` - массив логов обработки
- `results` - результаты обработки

**Методы:**
- `handleFileSelected()` - обработка выбранного файла
- `startPolling()` - запуск опроса статуса обработки
- `loadResults()` - загрузка результатов
- `resetApp()` - сброс приложения к начальному состоянию

### 2. UploadZone.vue
Компонент для загрузки аудиофайлов с поддержкой drag & drop.

**События:**
- `@file-selected` - эмитится при выборе файла

**Функции:**
- Drag & drop интерфейс
- Валидация типа файла (только audio/*)
- Визуальная обратная связь при наведении

### 3. ProcessingView.vue
Отображает процесс обработки файла.

**Props:**
- `file` - объект File
- `progress` - прогресс (0-100)
- `status` - текстовый статус
- `logs` - массив логов

**Функции:**
- Отображение информации о файле
- Прогресс-бар с анимацией
- Логи в реальном времени с автоскроллом
- Цветовая индикация уровня логов (info, success, warning, error)

### 4. ResultsView.vue
Отображает результаты обработки.

**Props:**
- `results` - объект с результатами
- `taskId` - ID задачи

**События:**
- `@new-file` - эмитится при нажатии "Загрузить новый файл"

**Функции:**
- Отображение оригинальной и анонимизированной транскрипции
- Аудиоплееры для прослушивания
- Скачивание транскрипций и аудиофайлов

## API Сервис (api.js)

Централизованный сервис для всех HTTP запросов к бэкенду.

**Методы:**

```javascript
audioService.uploadFile(file)
// POST /upload
// Загружает аудиофайл на сервер
// Возвращает: { task_id: string }

audioService.startProcessing(taskId)
// POST /process/{taskId}
// Запускает обработку загруженного файла
// Возвращает: { status: string }

audioService.getStatus(taskId)
// GET /status/{taskId}
// Получает текущий статус обработки
// Возвращает: {
//   status: 'processing' | 'completed' | 'failed',
//   progress: number (0-100),
//   logs: Array<{ message: string, level: string }>
// }

audioService.getResults(taskId)
// GET /results/{taskId}
// Получает результаты обработки
// Возвращает: {
//   original_transcript: string,
//   redacted_transcript: string,
//   original_audio_url: string,
//   redacted_audio_url: string
// }

audioService.downloadAudio(taskId, type)
// GET /download/{taskId}/{type}
// Скачивает аудиофайл (type: 'original' | 'redacted')
// Возвращает: Blob

audioService.getAudioUrl(path)
// Формирует полный URL для аудиофайла
```


## API Endpoints

### POST /upload
Загрузка аудиофайла.

**Request:**
```
Content-Type: multipart/form-data
Body: FormData { file: File }
```

**Response:**
```json
{
  "task_id": "uuid-string"
}
```

### POST /process/{task_id}
Запуск обработки.

**Response:**
```json
{
  "status": "processing"
}
```

### GET /status/{task_id}
Получение статуса (polling endpoint).

**Response:**
```json
{
  "status": "processing|completed|failed",
  "progress": 75,
  "logs": [
    {
      "message": "Транскрибация завершена",
      "level": "success",
      "timestamp": "2026-04-18T10:05:30Z"
    }
  ]
}
```

### GET /results/{task_id}
Получение результатов.

**Response:**
```json
{
  "original_transcript": "Полный текст с персональными данными",
  "redacted_transcript": "Текст с [REDACTED] вместо ПД (заменить на черные квадратики как у эпштейна)",
  "original_audio_url": "/files/original_uuid.mp3",
  "redacted_audio_url": "/files/redacted_uuid.mp3"
}
```

### GET /download/{task_id}/{type}
Скачивание аудиофайла.

**Parameters:**
- `type`: "original" | "redacted"

**Response:**
```
Content-Type: audio/mpeg
Body: binary audio data
```

### GET /files/{filename}
Статическая раздача аудиофайлов для плеера.

**Response:**
```
Content-Type: audio/mpeg
Body: binary audio data
```

## Установка зависимостей

```bash
cd frontend
npm install
```

## Запуск (dev mode)

```bash
npm run dev
```

Приложение будет доступно на `http://localhost:3000`

## Прод сборка

```bash
npm run build
```

Собранные файлы будут в директории `dist/`

## Preview production (смайлик с ноготочками)

```bash
npm run preview
```

## Конфиг

### Переменные окружения

надо создать файл `.env` в корне `frontend/`:

```env
VITE_API_URL=http://localhost:8000
```

### Vite конфигурация

`vite.config.js` настроен с:
- Vue plugin
- Dev server на порту 3000
- Proxy для API запросов (`/api` → `http://localhost:8000`)

### Как работает поллинг

Фронтенд опрашивает `/status/{task_id}` каждую секунду для получения:
- Актуального статуса обработки
- Прогресса (0-100%)
- Новых логов

Polling останавливается когда `status === 'completed'` или `status === 'failed'`

### по поводу обработки ошибок

- Try-catch блоки для всех API запросов
- Логирование ошибок в консоль логов
- Визуальная индикация ошибок (красным)

## TO DO (как вариант)

- WebSocket вместо polling для real-time обновлений
- Поддержка множественной загрузки файлов
- Темная/светлая тема (сейчас только темная)
- Тесты
