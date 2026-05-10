<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>Gin-Advanced</h2>
        <p>后台管理系统</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        autocomplete="on"
        label-position="left"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            placeholder="密码"
            name="password"
            type="password"
            tabindex="2"
            autocomplete="on"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item prop="captcha">
          <el-input
            v-model="loginForm.captcha"
            placeholder="验证码"
            name="captcha"
            type="text"
            tabindex="3"
            style="width: 60%"
            size="large"
            @keyup.enter="handleLogin"
          />
          <img
            :src="captchaUrl"
            class="captcha-image"
            @click="refreshCaptcha"
            alt="验证码"
          />
        </el-form-item>

        <el-button
          :loading="loading"
          type="primary"
          size="large"
          style="width: 100%"
          @click="handleLogin"
        >
          登 录
        </el-button>

        <div class="tips">
          <router-link to="/register">还没有账号？去注册</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getCaptcha } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)
const captchaUrl = ref('')
const captchaId = ref('')

const loginForm = reactive({
  username: 'admin',
  password: 'admin123',
  captcha: '',
  captcha_id: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const refreshCaptcha = async () => {
  try {
    const res = await getCaptcha()
    if (res.data) {
      captchaUrl.value = res.data.image || ''
      captchaId.value = res.data.captcha_id || ''
      loginForm.captcha_id = captchaId.value
    }
  } catch (error) {
    console.error('获取验证码失败', error)
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await userStore.login({
        username: loginForm.username,
        password: loginForm.password,
        captcha: loginForm.captcha,
        captcha_id: loginForm.captcha_id
      })
      ElMessage.success('登录成功')
      router.push('/')
    } catch (error) {
      console.error('登录失败', error)
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;

  .login-box {
    width: 450px;
    padding: 40px;
    background: #fff;
    border-radius: 10px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);

    .login-header {
      text-align: center;
      margin-bottom: 30px;

      h2 {
        margin: 0 0 10px;
        font-size: 28px;
        color: #333;
      }

      p {
        margin: 0;
        font-size: 14px;
        color: #666;
      }
    }

    .login-form {
      .captcha-image {
        width: 35%;
        height: 40px;
        margin-left: 10px;
        cursor: pointer;
        border-radius: 4px;
      }

      .tips {
        margin-top: 20px;
        text-align: right;

        a {
          color: #409eff;
          text-decoration: none;

          &:hover {
            text-decoration: underline;
          }
        }
      }
    }
  }
}
</style>
