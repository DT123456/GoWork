import request from '@/utils/request'

// 获取角色列表
export function getRoles() {
  return request({
    url: '/admin/roles',
    method: 'get'
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/admin/roles',
    method: 'post',
    data
  })
}

// 分配角色
export function assignRole(data) {
  return request({
    url: '/admin/roles/assign',
    method: 'post',
    data
  })
}

// 添加权限
export function addPermission(data) {
  return request({
    url: '/admin/roles/permission',
    method: 'post',
    data
  })
}

// 批量设置角色权限
export function setRolePermissions(data) {
  return request({
    url: '/admin/roles/permission',
    method: 'put',
    data
  })
}

// 创建权限
export function createPermission(data) {
  return request({
    url: '/admin/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(data) {
  return request({
    url: `/admin/permissions/${data.id}`,
    method: 'put',
    data
  })
}

// 删除权限
export function deletePermission(id) {
  return request({
    url: `/admin/permissions/${id}`,
    method: 'delete'
  })
}

// 初始化角色
export function seedRoles() {
  return request({
    url: '/admin/seed',
    method: 'get'
  })
}

// 获取权限列表
export function getPermissions() {
  return request({
    url: '/admin/roles/permission',
    method: 'get'
  })
}
