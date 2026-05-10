import request from '@/utils/request'

// 上传文件
export function uploadFile(formData) {
  return request({
    url: '/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 上传图片
export function uploadImage(formData) {
  return request({
    url: '/upload/image',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 上传头像
export function uploadAvatar(formData) {
  return request({
    url: '/upload/avatar',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 批量上传
export function uploadMultiple(formData) {
  return request({
    url: '/upload/multiple',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 获取上传配置
export function getUploadConfig() {
  return request({
    url: '/upload/config',
    method: 'get'
  })
}
