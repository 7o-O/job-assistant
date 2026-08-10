<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const jobDescription = ref('')
const question = ref('')
const answer = ref('')
const loading = ref(false)
const historyList = ref([])
const keyword = ref('')

const loadHistory = async () => {
  const res = await axios.get('/api/analyze/history', {
    params: {
      page: 1,
      pageSize: 10,
      keyword: keyword.value,
    },
  })

  if (res.data.success) {
    historyList.value = res.data.data.list
  }
}

const analyzeJob = async () => {
  if (!jobDescription.value.trim()) {
    alert('请输入岗位描述')
    return
  }

  if (!question.value.trim()) {
    alert('请输入问题')
    return
  }

  loading.value = true
  answer.value = ''

  try {
    const res = await axios.post('/api/analyze', {
      jobDescription: jobDescription.value,
      question: question.value,
    },{timeout:60000})

    if (res.data.success) {
      answer.value = res.data.data.answer
      await loadHistory()
    } else {
      alert(res.data.message)
    }
  } catch (err) {
  if (err.code === 'ECONNABORTED') {
    alert('请求超时，大模型响应太慢，请稍后再试')
  } else if (err.response) {
    alert(err.response.data.message || '后端返回错误')
  } else {
    alert('请求失败，请检查后端服务是否启动')
  }
} finally {
    loading.value = false
  }
}

const deleteHistory = async (id) => {
  if (!confirm('确定要删除这条记录吗？')) {
    return
  }

  const res = await axios.delete(`/api/analyze/history/${id}`)

  if (res.data.success) {
    await loadHistory()
  } else {
    alert(res.data.message)
  }
}

onMounted(() => {
  loadHistory()
})
</script>

<template>
  <div class="page">
    <h1>AI 岗位分析助手</h1>

    <section class="card">
      <h2>岗位分析</h2>

      <label>岗位描述</label>
      <textarea
        v-model="jobDescription"
        placeholder="例如：负责 Go 后端开发，熟悉 MySQL、Redis、RESTful API"
      ></textarea>

      <label>问题</label>
      <input
        v-model="question"
        placeholder="例如：这个岗位需要掌握哪些技能？"
      />

      <button @click="analyzeJob" :disabled="loading">
        {{ loading ? '分析中...' : '开始分析' }}
      </button>
    </section>

    <section v-if="answer" class="card">
      <h2>分析结果</h2>
      <pre>{{ answer }}</pre>
    </section>

    <section class="card">
      <h2>历史记录</h2>

      <div class="search">
        <input
          v-model="keyword"
          placeholder="输入关键词搜索历史记录"
          @keyup.enter="loadHistory"
        />
        <button @click="loadHistory">搜索</button>
      </div>

      <div v-if="historyList.length === 0" class="empty">
        暂无历史记录
      </div>

      <div
        v-for="item in historyList"
        :key="item.ID || item.Id || item.id"
        class="history-item"
      >
        <div class="history-title">
          <strong>{{ item.question || item.Question }}</strong>
          <button @click="deleteHistory(item.ID || item.Id || item.id)">
            删除
          </button>
        </div>

        <p>{{ item.jobDescription || item.JobDescription }}</p>
        <pre>{{ item.answer || item.Answer }}</pre>
      </div>
    </section>
  </div>
</template>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px;
  color: #1f2937;
}

h1 {
  margin-bottom: 24px;
  text-align: center;
}

.card {
  margin-bottom: 24px;
  padding: 24px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  background: #ffffff;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

h2 {
  margin-top: 0;
  margin-bottom: 16px;
}

label {
  display: block;
  margin: 14px 0 8px;
  font-weight: 600;
}

textarea,
input {
  width: 100%;
  box-sizing: border-box;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
}

textarea {
  min-height: 130px;
  resize: vertical;
}

button {
  margin-top: 16px;
  padding: 10px 18px;
  border: none;
  border-radius: 8px;
  background: #2563eb;
  color: #ffffff;
  cursor: pointer;
}

button:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

pre {
  white-space: pre-wrap;
  line-height: 1.7;
  font-family: inherit;
}

.search {
  display: flex;
  gap: 12px;
}

.search input {
  flex: 1;
}

.search button {
  margin-top: 0;
}

.empty {
  margin-top: 16px;
  color: #6b7280;
}

.history-item {
  margin-top: 16px;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f9fafb;
}

.history-title {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.history-title button {
  margin-top: 0;
  background: #dc2626;
}

.history-item p {
  color: #4b5563;
}
</style>
