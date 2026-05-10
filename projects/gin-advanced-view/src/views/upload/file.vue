<template>
  <div class="upload-container">
    <el-card>
      <template #header>
        <span>文件上传</span>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="单文件上传" name="single">
          <el-upload
            class="upload-demo"
            drag
            :action="uploadUrl"
            :headers="headers"
            :on-success="handleSingleSuccess"
            :on-error="handleError"
            :before-upload="beforeUpload"
            accept=".jpg,.jpeg,.png,.gif,.webp,.pdf,.doc,.docx,.xls,.xlsx"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 jpg、png、gif、webp、pdf、doc、docx、xls、xlsx 格式，文件大小不超过 10MB
              </div>
            </template>
          </el-upload>
        </el-tab-pane>

        <el-tab-pane label="图片上传" name="image">
          <el-upload
            class="avatar-uploader"
            :action="imageUploadUrl"
            :headers="headers"
            :show-file-list="false"
            :on-success="handleImageSuccess"
            :on-error="handleError"
            :before-upload="beforeImageUpload"
            accept=".jpg,.jpeg,.png,.gif,.webp"
          >
            <img v-if="imageUrl" :src="imageUrl" class="avatar" />
            <el-icon v-else class="avatar-uploader-icon"><plus /></el-icon>
          </el-upload>
        </el-tab-pane>

        <el-tab-pane label="批量上传" name="multiple">
          <el-upload
            class="upload-demo"
            drag
            multiple
            :action="uploadUrl"
            :headers="headers"
            :on-success="handleMultiSuccess"
            :on-error="handleError"
            :before-upload="beforeUpload"
            accept=".jpg,.jpeg,.png,.gif,.webp,.pdf,.doc,.docx,.xls,.xlsx"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将多个文件拖到此处，或<em>点击上传</em>
            </div>
          </el-upload>
        </el-tab-pane>
      </el-tabs>

      <el-divider />

      <h4>上传历史</h4>
      <el-table :data="uploadHistory" stripe>
        <el-table-column prop="filename" label="文件名" />
        <el-table-column prop="url" label="链接">
          <template #default="{ row }">
            <el-link :href="row.url" target="_blank" type="primary">{{ row.url }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="上传时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import Cookies from 'js-cookie'

const activeTab = ref('single')
const imageUrl = ref('')
const uploadHistory = ref([])
const uploadUrl = '/api/upload'
const imageUploadUrl = '/api/upload/image'

const headers = {
  'Authorization': `Bearer ${Cookies.get('Admin-Token')}`
}

const beforeUpload = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

const beforeImageUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

const handleSingleSuccess = (response) => {
  ElMessage.success('上传成功')
  addToHistory(response.data)
}

const handleImageSuccess = (response) => {
  ElMessage.success('上传成功')
  imageUrl.value = response.data.url
  addToHistory(response.data)
}

const handleMultiSuccess = (response) => {
  addToHistory(response.data)
}

const handleError = (err) => {
  ElMessage.error('上传失败')
  console.error(err)
}

const addToHistory = (data) => {
  uploadHistory.value.unshift({
    filename: data.filename || data.name,
    url: data.url,
    time: new Date().toLocaleString('zh-CN')
  })
}
</script>

<style lang="scss" scoped>
.upload-container {
  .avatar-uploader {
    width: 178px;
    height: 178px;
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: border-color 0.3s;

    &:hover {
      border-color: #409eff;
    }

    .avatar {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .avatar-uploader-icon {
      font-size: 28px;
      color: #8c939d;
      width: 178px;
      height: 178px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
