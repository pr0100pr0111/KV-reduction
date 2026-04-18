import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: 60000,
})

export const jobService = {
  async uploadFile(file, onUploadProgress) {
    const formData = new FormData()
    formData.append('audio', file)

    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-0.5data',
      },
      onUploadProgress,
    })
    return response.data
  },

  async getStatus(jobId) {
    const response = await api.get(`/job/${jobId}`)
    return response.data
  },

  async getResults(jobId) {
    const response = await api.get(`/job/${jobId}`)
    return response.data
  },

  async downloadBlob(jobId, type) {
    const response = await api.get(`/download/${jobId}/${type}`, {
      responseType: 'blob',
    })
    return response.data
  },
}

export const getDownloadUrl = (jobId, type, fileName) => {
  return `${API_BASE_URL}/api/v1/download/${jobId}/${type}?fileName=${encodeURIComponent(fileName)}`
}

export default api
