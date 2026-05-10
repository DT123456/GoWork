import { defineStore } from 'pinia'
import Cookies from 'js-cookie'
import { login as loginApi, getUserInfo } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: Cookies.get('Admin-Token') || '',
    userId: '',
    username: '',
    nickname: '',
    roles: [],
    permissions: []
  }),

  actions: {
    // 登录
    async login(userInfo) {
      const { data } = await loginApi(userInfo)
      this.token = data.token
      this.userId = data.user_id
      this.username = data.username
      this.nickname = data.nickname
      this.roles = data.roles || []
      this.permissions = data.permissions || []
      Cookies.set('Admin-Token', data.token, { expires: 7 })
      return data
    },

    // 获取用户信息
    async getInfo() {
      if (!this.userId) return
      const { data } = await getUserInfo(this.userId)
      this.nickname = data.nickname
      this.roles = data.roles || []
      this.permissions = data.permissions || []
      return data
    },

    // 登出
    logout() {
      this.token = ''
      this.userId = ''
      this.username = ''
      this.nickname = ''
      this.roles = []
      this.permissions = []
      Cookies.remove('Admin-Token')
    },

    // 重置 token
    resetToken() {
      this.token = ''
      Cookies.remove('Admin-Token')
    }
  },

  persist: true
})
