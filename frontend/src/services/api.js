import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const audioService = {
  async uploadFile(file) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  async startProcessing(taskId) {
    const response = await api.post(`/process/${taskId}`)
    return response.data
  },

  async getStatus(taskId) {
    const response = await api.get(`/status/${taskId}`)
    return response.data
  },

  async getResults(taskId) {
    const response = await api.get(`/results/${taskId}`)
    return response.data
  },

  async downloadAudio(taskId, type) {
    const response = await api.get(`/download/${taskId}/${type}`, {
      responseType: 'blob'
    })
    return response.data
  },

  getAudioUrl(path) {
    return `${API_BASE_URL}${path}`
  }
}

export default api
