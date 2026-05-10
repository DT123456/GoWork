<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #409eff">
            <el-icon :size="30"><User /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">用户总数</p>
            <p class="stat-value">{{ stats.userCount || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #67c23a">
            <el-icon :size="30"><UserFilled /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">角色数量</p>
            <p class="stat-value">{{ stats.roleCount || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #e6a23c">
            <el-icon :size="30"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">系统日志</p>
            <p class="stat-value">{{ stats.systemLogCount || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #f56c6c">
            <el-icon :size="30"><Connection /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">访问日志</p>
            <p class="stat-value">{{ stats.accessLogCount || 0 }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">
              {{ userStore.username }}
            </el-descriptions-item>
            <el-descriptions-item label="昵称">
              {{ userStore.nickname }}
            </el-descriptions-item>
            <el-descriptions-item label="角色">
              <el-tag v-for="role in userStore.roles" :key="role" size="small" style="margin-right: 5px">
                {{ role }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="权限">
              <el-tag
                v-for="perm in userStore.permissions"
                :key="perm"
                size="small"
                type="info"
                style="margin-right: 5px; margin-bottom: 2px"
              >
                {{ perm }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最近系统日志</span>
          </template>
          <el-table :data="recentLogs" style="width: 100%">
            <el-table-column prop="username" label="用户" width="100" />
            <el-table-column prop="action" label="操作" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getSystemLogs } from '@/api/log'
import { getRoles } from '@/api/role'

const userStore = useUserStore()
const stats = ref({})
const recentLogs = ref([])

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

onMounted(async () => {
  try {
    const [logsRes, rolesRes] = await Promise.all([
      getSystemLogs({ page: 1, page_size: 5 }),
      getRoles()
    ])
    stats.value = {
      roleCount: rolesRes.data?.length || 0,
      systemLogCount: logsRes.data?.total || 0,
      userCount: '-'
    }
    recentLogs.value = logsRes.data?.list || []
  } catch (error) {
    console.error('获取数据失败', error)
  }
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 20px;
}

.stat-card {
  display: flex;
  align-items: center;

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    margin-right: 20px;
  }

  .stat-info {
    .stat-label {
      margin: 0;
      font-size: 14px;
      color: #909399;
    }

    .stat-value {
      margin: 5px 0 0;
      font-size: 28px;
      font-weight: 600;
      color: #303133;
    }
  }
}
</style>
