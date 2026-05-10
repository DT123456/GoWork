<template>
  <div class="role-list-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary" @click="handleCreate">新增角色</el-button>
        </div>
      </template>

      <el-table :data="roleList" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="permissions" label="权限" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="perm in row.permissions"
              :key="perm.id"
              size="small"
              type="info"
              style="margin-right: 5px; margin-bottom: 2px"
            >
              {{ perm.name }}
            </el-tag>
            <span v-if="!row.permissions?.length" style="color: #999">无</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleAssignPerm(row)">
              分配权限
            </el-button>
            <el-button type="danger" size="small" :disabled="row.name === 'admin'" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增对话框 -->
    <el-dialog v-model="createDialogVisible" title="新增角色" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="rules" label-width="80px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="createForm.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配权限对话框 -->
    <el-dialog v-model="permDialogVisible" title="分配权限" width="600px">
      <el-form ref="permFormRef" :model="permForm" label-width="80px">
        <el-form-item label="角色">
          <el-select v-model="permForm.role_id" placeholder="请选择角色" style="width: 100%" disabled>
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择权限">
          <el-checkbox-group v-model="permForm.permission_ids">
            <el-checkbox
              v-for="perm in permissionList"
              :key="perm.id"
              :label="perm.id"
              :value="perm.id"
            >
              {{ perm.name }} ({{ perm.method }} {{ perm.path }})
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPermAssign">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoles, createRole, getPermissions, setRolePermissions } from '@/api/role'

const loading = ref(false)
const roleList = ref([])
const permissionList = ref([])
const createDialogVisible = ref(false)
const permDialogVisible = ref(false)
const createFormRef = ref()
const permFormRef = ref()
const createForm = reactive({
  name: '',
  description: ''
})
const permForm = reactive({
  role_id: '',
  permission_ids: []
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }]
}

const fetchRoleList = async () => {
  loading.value = true
  try {
    const res = await getRoles()
    roleList.value = res.data || []
  } catch (error) {
    console.error('获取角色列表失败', error)
  } finally {
    loading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    const res = await getPermissions()
    permissionList.value = res.data || []
  } catch (error) {
    console.error('获取权限列表失败', error)
  }
}

const handleCreate = () => {
  createForm.name = ''
  createForm.description = ''
  createDialogVisible.value = true
}

const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      await createRole(createForm)
      ElMessage.success('创建成功')
      createDialogVisible.value = false
      fetchRoleList()
    } catch (error) {
      console.error('创建失败', error)
    }
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该角色吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).catch(() => {})
}

const handleAssignPerm = (row) => {
  permForm.role_id = row.id
  permForm.permission_ids = row.permissions?.map(p => p.id) || []
  permDialogVisible.value = true
}

const submitPermAssign = async () => {
  if (!permForm.role_id) {
    ElMessage.warning('请选择角色')
    return
  }
  try {
    await setRolePermissions({
      role_id: permForm.role_id,
      permission_ids: permForm.permission_ids
    })
    ElMessage.success('分配成功')
    permDialogVisible.value = false
    fetchRoleList()
  } catch (error) {
    console.error('分配失败', error)
  }
}

onMounted(() => {
  fetchRoleList()
  fetchPermissions()
})
</script>

<style lang="scss" scoped>
.role-list-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
