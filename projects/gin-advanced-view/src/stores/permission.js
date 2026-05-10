import { createRouter, createWebHistory } from 'vue-router'
import Cookies from 'js-cookie'

// 静态路由
export const constantRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404' }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

// 动态路由
export const asyncRoutes = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页', icon: 'Odometer' }
      },
      {
        path: '/user',
        meta: { title: '用户管理', icon: 'User' },
        children: [
          {
            path: '/user/list',
            name: 'UserList',
            component: () => import('@/views/user/list.vue'),
            meta: { title: '用户列表', icon: 'List' }
          }
        ]
      },
      {
        path: '/role',
        meta: { title: '角色管理', icon: 'UserFilled', roles: ['admin'] },
        children: [
          {
            path: '/role/list',
            name: 'RoleList',
            component: () => import('@/views/role/list.vue'),
            meta: { title: '角色列表', icon: 'List', roles: ['admin'] }
          },
          {
            path: '/role/permission',
            name: 'RolePermission',
            component: () => import('@/views/role/permission.vue'),
            meta: { title: '角色权限', icon: 'Key', roles: ['admin'] }
          }
        ]
      },
      {
        path: '/log',
        meta: { title: '日志管理', icon: 'Document', roles: ['admin'] },
        children: [
          {
            path: '/log/system',
            name: 'SystemLog',
            component: () => import('@/views/log/system.vue'),
            meta: { title: '系统日志', icon: 'DocumentCopy', roles: ['admin'] }
          },
          {
            path: '/log/access',
            name: 'AccessLog',
            component: () => import('@/views/log/access.vue'),
            meta: { title: '访问日志', icon: 'Connection', roles: ['admin'] }
          }
        ]
      },
      {
        path: '/upload/file',
        name: 'FileUpload',
        component: () => import('@/views/upload/file.vue'),
        meta: { title: '文件上传', icon: 'UploadFilled' }
      }
    ]
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes: [...constantRoutes, ...asyncRoutes]
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  const hasToken = Cookies.get('Admin-Token')

  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      next()
    }
  } else {
    if (to.path === '/login' || to.path === '/register') {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router
