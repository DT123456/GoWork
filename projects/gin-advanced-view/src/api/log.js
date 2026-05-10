import request from '@/utils/request'

// 获取系统日志列表
export function getSystemLogs(params) {
  return request({
    url: '/admin/logs/system',
    method: 'get',
    params
  })
}

// 获取系统日志详情
export function getSystemLogDetail(id) {
  return request({
    url: `/admin/logs/system/${id}`,
    method: 'get'
  })
}

// 获取系统日志统计
export function getSystemLogStats(params) {
  return request({
    url: '/admin/logs/system/stats',
    method: 'get',
    params
  })
}

// 删除系统日志
export function deleteSystemLogs(ids) {
  return request({
    url: '/admin/logs/system',
    method: 'delete',
    data: { ids }
  })
}

// 获取访问日志列表
export function getAccessLogs(params) {
  return request({
    url: '/admin/logs/access',
    method: 'get',
    params
  })
}

// 获取访问统计
export function getAccessLogStats(params) {
  return request({
    url: '/admin/logs/access/stats',
    method: 'get',
    params
  })
}

// 清理访问日志
export function cleanAccessLogs(days) {
  return request({
    url: '/admin/logs/access',
    method: 'delete',
    data: { days }
  })
}
