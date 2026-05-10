<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <aside class="sidebar-container" :class="{ 'is-collapse': isCollapse }">
      <div class="logo-container">
        <h1 class="logo-title">Gin-Admin</h1>
      </div>
      <el-scrollbar>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :background-color="'#304156'"
          :text-color="'#bfcbd9'"
          :unique-opened="true"
          :collapse-transition="false"
          mode="vertical"
          router
        >
          <sidebar-item
            v-for="route in sidebarRoutes"
            :key="route.path"
            :item="route"
          />
        </el-menu>
      </el-scrollbar>
    </aside>

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <div class="navbar">
        <div class="left">
          <hamburger
            :is-active="isCollapse"
            class="hamburger"
            @toggleClick="toggleSideBar"
          />
          <breadcrumb />
        </div>
        <div class="right">
          <el-dropdown @command="handleCommand">
            <span class="el-dropdown-link">
              <el-avatar :size="32" icon="UserFilled" />
              <span class="username">{{ userStore.nickname || userStore.username }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 内容区 -->
      <app-main />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { asyncRoutes } from '@/stores/permission'
import SidebarItem from './components/SidebarItem.vue'
import Hamburger from './components/Hamburger.vue'
import Breadcrumb from './components/Breadcrumb.vue'
import AppMain from './components/AppMain.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)

const activeMenu = computed(() => {
  const { path } = route
  return path
})

const sidebarRoutes = computed(() => {
  const userRoles = userStore.roles
  // 获取根路由的 children 作为侧边栏菜单
  const rootRoute = asyncRoutes.find(r => r.path === '/')
  if (!rootRoute) return []
  
  const routes = rootRoute.children || []
  
  // 根据角色过滤
  return routes.filter(item => {
    if (!item.meta?.roles) return true
    return item.meta.roles.some(role => userRoles.includes(role))
  })
})

const toggleSideBar = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.app-wrapper {
  display: flex;
  height: 100%;
}

.sidebar-container {
  width: 210px;
  height: 100%;
  background: #304156;
  transition: width 0.28s;
  overflow: hidden;

  &.is-collapse {
    width: 64px;
  }

  .logo-container {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #263445;

    .logo-title {
      color: #fff;
      font-size: 20px;
      font-weight: 600;
      margin: 0;
    }
  }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.navbar {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

  .left {
    display: flex;
    align-items: center;
  }

  .hamburger {
    cursor: pointer;
    margin-right: 20px;
  }

  .right {
    .el-dropdown-link {
      display: flex;
      align-items: center;
      cursor: pointer;

      .username {
        margin-left: 8px;
        font-size: 14px;
      }
    }
  }
}
</style>
