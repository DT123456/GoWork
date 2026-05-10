<template>
  <div class="role-permission-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色权限管理</span>
          <div>
            <el-button type="primary" @click="handleCreate">新增权限</el-button>
            <el-button type="success" @click="handleSeed" style="margin-left: 10px">初始化数据</el-button>
          </div>
        </div>
      </template>

      <el-table :data="permissionList" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限名称" />
        <el-table-column prop="path" label="路径" />
        <el-table-column prop="method" label="方法" width="100" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEditPerm(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDeletePerm(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分配权限对话框 -->
      <el-dialog v-model="dialogVisible" title="分配权限" width="600px">
        <el-form ref="formRef" :model="form" label-width="100px">
          <el-form-item label="角色">
            <el-select v-model="form.role_id" placeholder="请选择角色" style="width: 100%">
              <el-option
                v-for="role in roleList"
                :key="role.id"
                :label="role.name"
                :value="role.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="选择权限">
            <el-checkbox-group v-model="form.permission_ids">
              <el-checkbox
                v-for="perm in permissionList"
                :key="perm.id"
                :label="perm.id"
                :value="perm.id"
              >
                {{ perm.name }}
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </template>
      </el-dialog>

      <!-- 新增/编辑权限对话框 -->
      <el-dialog v-model="permDialogVisible" :title="permForm.id ? '编辑权限' : '新增权限'" width="500px">
        <el-form ref="permFormRef" :model="permForm" :rules="permRules" label-width="100px">
          <el-form-item label="权限名称" prop="name">
            <el-input v-model="permForm.name" placeholder="如: user:create" />
          </el-form-item>
          <el-form-item label="请求路径" prop="path">
            <el-input v-model="permForm.path" placeholder="如: /user/*" />
          </el-form-item>
          <el-form-item label="请求方法" prop="method">
            <el-select v-model="permForm.method" placeholder="请选择" style="width: 100%">
              <el-option label="全部" value="*" />
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="permForm.description" type="textarea" placeholder="请输入描述" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="permDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPermForm">确定</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPermissions, assignRole, seedRoles, createPermission, updatePermission, deletePermission } from '@/api/role'
import { getRoles } from '@/api/role'

const loading = ref(false)
const permissionList = ref([])
const roleList = ref([])
const dialogVisible = ref(false)
const permDialogVisible = ref(false)
const formRef = ref()
const permFormRef = ref()
const form = reactive({
  role_id: '',
  permission_ids: []
})
const permForm = reactive({
  id: null,
  name: '',
  path: '',
  method: '*',
  description: ''
})
const permRules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  path: [{ required: true, message: '请输入请求路径', trigger: 'blur' }],
  method: [{ required: true, message: '请选择请求方法', trigger: 'change' }]
}

const fetchPermissions = async () => {
  loading.value = true
  try {
    const res = await getPermissions()
    permissionList.value = res.data || []
  } catch (error) {
    console.error('获取权限列表失败', error)
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await getRoles()
    roleList.value = res.data || []
  } catch (error) {
    console.error('获取角色列表失败', error)
  }
}

const handleEdit = (row) => {
  form.role_id = row.id
  form.permission_ids = []
  dialogVisible.value = true
}

const handleCreate = () => {
  permForm.id = null
  permForm.name = ''
  permForm.path = ''
  permForm.method = '*'
  permForm.description = ''
  permDialogVisible.value = true
}

const handleEditPerm = (row) => {
  permForm.id = row.id
  permForm.name = row.name
  permForm.path = row.path
  permForm.method = row.method
  permForm.description = row.description || ''
  permDialogVisible.value = true
}

const submitPermForm = async () => {
  if (!permFormRef.value) return
  await permFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (permForm.id) {
        await updatePermission(permForm)
        ElMessage.success('更新成功')
      } else {
        await createPermission(permForm)
        ElMessage.success('创建成功')
      }
      permDialogVisible.value = false
      fetchPermissions()
    } catch (error) {
      console.error('操作失败', error)
    }
  })
}

const handleDeletePerm = (row) => {
  ElMessageBox.confirm('确定要删除该权限吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePermission(row.id)
      ElMessage.success('删除成功')
      fetchPermissions()
    } catch (error) {
      console.error('删除失败', error)
    }
  }).catch(() => {})
}

const handleSeed = async () => {
  try {
    await seedRoles()
    ElMessage.success('初始化成功')
    fetchPermissions()
  } catch (error) {
    console.error('初始化失败', error)
  }
}

const submitForm = async () => {
  if (!form.role_id) {
    ElMessage.warning('请选择角色')
    return
  }
  try {
    await assignRole({
      role_id: form.role_id,
      permission_ids: form.permission_ids
    })
    ElMessage.success('分配成功')
    dialogVisible.value = false
  } catch (error) {
    console.error('分配失败', error)
  }
}

onMounted(() => {
  fetchPermissions()
  fetchRoles()
})
</script>

<style lang="scss" scoped>
.role-permission-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
