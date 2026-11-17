<template>
  <div class="login-page">
    <div class="bg-layer" aria-hidden="true">
      <div class="bg-stars" />
      <div class="bg-grid" />
      <div class="bg-hex" />
      <div class="bg-scan" />
      <div class="bg-orb bg-orb-a" />
      <div class="bg-orb bg-orb-b" />
      <div class="bg-ring" />
      <div class="bg-noise" />
    </div>

    <div class="login-shell">
      <span class="hud-corner hud-tl" aria-hidden="true" />
      <span class="hud-corner hud-tr" aria-hidden="true" />
      <span class="hud-corner hud-bl" aria-hidden="true" />
      <span class="hud-corner hud-br" aria-hidden="true" />

      <section class="brand-panel">
        <div class="brand-mark">
          <span class="brand-dot" />
          <span>SYSTEM ONLINE</span>
        </div>

        <h1 class="brand-title">
          <span class="brand-wtv">WTV</span>
          <span class="brand-sub">管理系统</span>
        </h1>

        <p class="brand-desc">安全接入 · 身份校验 · 会话建立</p>

        <div class="tech-meter" aria-hidden="true">
          <div class="meter-row">
            <span>LINK</span>
            <div class="meter-bar"><i style="--w: 86%" /></div>
            <em>86%</em>
          </div>
          <div class="meter-row">
            <span>NODE</span>
            <div class="meter-bar"><i style="--w: 72%" /></div>
            <em>72%</em>
          </div>
          <div class="meter-row">
            <span>AUTH</span>
            <div class="meter-bar"><i style="--w: 94%" /></div>
            <em>94%</em>
          </div>
        </div>

        <div class="tech-code" aria-hidden="true">
          <span>> init secure channel...</span>
          <span>> handshake ok</span>
          <span class="blink">> awaiting credentials_</span>
        </div>
      </section>

      <section class="form-panel">
        <div class="form-head">
          <div class="form-chip">ACCESS GATE</div>
          <h2>身份认证</h2>
          <p>请输入管理员账号与密码</p>
        </div>

        <el-form class="login-form" :model="form" @submit.prevent="submit">
          <label class="field-label" for="login-username">USERNAME</label>
          <el-form-item>
            <el-input
              id="login-username"
              v-model="form.username"
              placeholder="用户名"
              size="large"
              autocomplete="username"
              @keyup.enter="submit"
            />
          </el-form-item>

          <label class="field-label" for="login-password">PASSWORD</label>
          <el-form-item>
            <el-input
              id="login-password"
              v-model="form.password"
              type="password"
              placeholder="密码"
              show-password
              size="large"
              autocomplete="current-password"
              @keyup.enter="submit"
            />
          </el-form-item>

          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="submit"
          >
            启动会话
          </el-button>
        </el-form>

        <div class="form-foot">
          <span class="foot-id">CH-07</span>
          <span class="foot-sep" />
          <span>ENCRYPTED CHANNEL</span>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { setSession } from '../auth'

const emit = defineEmits(['success'])

const form = reactive({
  username: '',
  password: ''
})
const loading = ref(false)

const submit = async () => {
  const username = form.username.trim()
  if (!username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const response = await axios.post('/admin/login', {
      username,
      password: form.password
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '登录失败')
      return
    }
    setSession(response.data.data.token, response.data.data.admin)
    emit('success', response.data.data.admin)
  } catch {
    ElMessage.error('登录失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  document.title = 'WTV管理系统'
})
</script>

<style scoped>
/* 管理后台多为 http 公网 IP 访问；外链字体（fonts.bunny.net）会被浏览器
   Private Network Access / 非安全上下文策略拦截，故改用系统字体栈。 */
.login-page {
  --font-sans: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", system-ui, -apple-system, "Segoe UI", sans-serif;
  --font-display: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", system-ui, -apple-system, "Segoe UI", sans-serif;
  --font-mono: ui-monospace, "SF Mono", Menlo, Consolas, "Courier New", monospace;
  --bg: #05080f;
  --panel: rgba(8, 14, 28, 0.82);
  --line: rgba(120, 240, 255, 0.18);
  --text: #EAF0F8;
  --muted: #7E8BA3;
  --accent: #78F0FF;

  position: relative;
  min-height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 28px 16px;
  box-sizing: border-box;
  background:
    radial-gradient(ellipse 80% 50% at 50% -10%, rgba(40, 90, 120, 0.35), transparent 55%),
    var(--bg);
  color: var(--text);
  font-family: var(--font-sans);
}

.bg-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.bg-stars {
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(1px 1px at 12% 18%, rgba(234, 240, 248, 0.55), transparent),
    radial-gradient(1px 1px at 28% 62%, rgba(120, 240, 255, 0.45), transparent),
    radial-gradient(1.5px 1.5px at 68% 24%, rgba(234, 240, 248, 0.4), transparent),
    radial-gradient(1px 1px at 82% 70%, rgba(120, 240, 255, 0.35), transparent),
    radial-gradient(1px 1px at 46% 40%, rgba(234, 240, 248, 0.3), transparent);
  animation: starsTwinkle 5s ease-in-out infinite;
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(120, 240, 255, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(120, 240, 255, 0.045) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse 75% 65% at 50% 45%, #000 15%, transparent 78%);
  animation: gridPulse 8s ease-in-out infinite;
}

.bg-hex {
  position: absolute;
  inset: -10%;
  opacity: 0.12;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='84' height='48'%3E%3Cpath d='M21 0 L42 12 L42 36 L21 48 L0 36 L0 12 Z' fill='none' stroke='%2378F0FF' stroke-width='1'/%3E%3C/svg%3E");
  background-size: 84px 48px;
  animation: hexDrift 40s linear infinite;
}

.bg-scan {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgba(120, 240, 255, 0.03) 48%,
    rgba(120, 240, 255, 0.08) 50%,
    rgba(120, 240, 255, 0.03) 52%,
    transparent 100%
  );
  background-size: 100% 220%;
  animation: scanMove 4.5s linear infinite;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.bg-orb-a {
  width: 460px;
  height: 460px;
  top: -160px;
  left: -100px;
  background: radial-gradient(circle, rgba(74, 174, 214, 0.4), transparent 68%);
  animation: orbFloat 11s ease-in-out infinite;
}

.bg-orb-b {
  width: 380px;
  height: 380px;
  right: -90px;
  bottom: -120px;
  background: radial-gradient(circle, rgba(120, 240, 255, 0.22), transparent 70%);
  animation: orbFloat 13s ease-in-out infinite reverse;
}

.bg-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  width: min(720px, 90vw);
  height: min(720px, 90vw);
  transform: translate(-50%, -50%);
  border: 1px solid rgba(120, 240, 255, 0.08);
  border-radius: 50%;
  box-shadow:
    inset 0 0 60px rgba(120, 240, 255, 0.04),
    0 0 80px rgba(120, 240, 255, 0.03);
  animation: ringSpin 36s linear infinite;
}

.bg-ring::before {
  content: '';
  position: absolute;
  inset: 18%;
  border: 1px dashed rgba(120, 240, 255, 0.12);
  border-radius: 50%;
}

.bg-noise {
  position: absolute;
  inset: 0;
  opacity: 0.04;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
}

.login-shell {
  position: relative;
  z-index: 1;
  width: min(900px, 100%);
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: var(--panel);
  backdrop-filter: blur(20px);
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(120, 240, 255, 0.04),
    0 24px 80px rgba(0, 0, 0, 0.45);
  animation: shellIn 0.65s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.hud-corner {
  position: absolute;
  width: 18px;
  height: 18px;
  z-index: 2;
  border-color: var(--accent);
  border-style: solid;
  opacity: 0.75;
}

.hud-tl { top: 10px; left: 10px; border-width: 2px 0 0 2px; }
.hud-tr { top: 10px; right: 10px; border-width: 2px 2px 0 0; }
.hud-bl { bottom: 10px; left: 10px; border-width: 0 0 2px 2px; }
.hud-br { bottom: 10px; right: 10px; border-width: 0 2px 2px 0; }

.brand-panel {
  padding: 44px 40px;
  border-right: 1px solid var(--line);
  background:
    linear-gradient(145deg, rgba(120, 240, 255, 0.07), transparent 40%),
    linear-gradient(180deg, rgba(10, 18, 34, 0.95), rgba(5, 8, 15, 0.7));
}

.brand-mark {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 26px;
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.24em;
}

.brand-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--accent);
  animation: pulseDot 2s ease-out infinite;
}

.brand-title {
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.brand-wtv {
  font-family: var(--font-display);
  font-size: clamp(54px, 8vw, 80px);
  font-weight: 800;
  line-height: 0.9;
  letter-spacing: -0.045em;
  background: linear-gradient(115deg, #F4F8FF 8%, #78F0FF 48%, #3D9EC8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  text-shadow: 0 0 40px rgba(120, 240, 255, 0.15);
}

.brand-sub {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--text);
}

.brand-desc {
  margin: 18px 0 0;
  color: var(--muted);
  font-size: 13px;
  letter-spacing: 0.04em;
}

.tech-meter {
  margin-top: 34px;
  display: grid;
  gap: 10px;
}

.meter-row {
  display: grid;
  grid-template-columns: 42px 1fr 36px;
  align-items: center;
  gap: 10px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--muted);
}

.meter-row em {
  font-style: normal;
  color: var(--accent);
  text-align: right;
}

.meter-bar {
  height: 4px;
  border-radius: 2px;
  background: rgba(120, 240, 255, 0.1);
  overflow: hidden;
}

.meter-bar i {
  display: block;
  width: var(--w);
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, #4AAED6, #78F0FF);
  box-shadow: 0 0 10px rgba(120, 240, 255, 0.45);
  animation: meterFill 1.4s ease-out both;
}

.tech-code {
  margin-top: 28px;
  padding: 14px 14px 12px;
  border: 1px solid rgba(120, 240, 255, 0.14);
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.28);
  display: grid;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: rgba(120, 240, 255, 0.72);
  letter-spacing: 0.02em;
}

.tech-code .blink {
  animation: caretBlink 1.1s step-end infinite;
}

.form-panel {
  padding: 44px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-chip {
  display: inline-flex;
  margin-bottom: 14px;
  padding: 4px 10px;
  border: 1px solid rgba(120, 240, 255, 0.28);
  border-radius: 6px;
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.18em;
  background: rgba(120, 240, 255, 0.06);
}

.form-head h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--text);
  letter-spacing: 0.04em;
}

.form-head p {
  margin: 8px 0 26px;
  color: var(--muted);
  font-size: 13px;
}

.field-label {
  display: block;
  margin-bottom: 8px;
  color: rgba(120, 240, 255, 0.75);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.14em;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.login-form :deep(.el-input__wrapper) {
  background: rgba(6, 12, 24, 0.92);
  box-shadow: 0 0 0 1px rgba(120, 240, 255, 0.16) inset;
  border-radius: 10px;
  transition: box-shadow 0.2s ease, background 0.2s ease;
}

.login-form :deep(.el-input__wrapper:hover),
.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow:
    0 0 0 1px rgba(120, 240, 255, 0.55) inset,
    0 0 18px rgba(120, 240, 255, 0.12);
  background: rgba(8, 18, 34, 0.96);
}

.login-button {
  width: 100%;
  margin-top: 8px;
  height: 46px;
  border: 1px solid rgba(120, 240, 255, 0.35);
  border-radius: 10px;
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.12em;
  background: linear-gradient(120deg, #2F8FB8 0%, #78F0FF 100%);
  color: #031018;
  transition: transform 0.2s ease, filter 0.2s ease, box-shadow 0.2s ease;
  box-shadow: 0 0 24px rgba(120, 240, 255, 0.18);
}

.login-button:hover {
  filter: brightness(1.07);
  transform: translateY(-1px);
  box-shadow: 0 0 32px rgba(120, 240, 255, 0.28);
}

.form-foot {
  margin-top: 22px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #536074;
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.16em;
}

.foot-id {
  color: var(--accent);
}

.foot-sep {
  width: 18px;
  height: 1px;
  background: rgba(120, 240, 255, 0.25);
}

@keyframes shellIn {
  from { opacity: 0; transform: translateY(16px) scale(0.985); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes gridPulse {
  0%, 100% { opacity: 0.55; }
  50% { opacity: 0.9; }
}

@keyframes hexDrift {
  from { transform: translateY(0); }
  to { transform: translateY(-48px); }
}

@keyframes scanMove {
  from { background-position: 0 -40%; }
  to { background-position: 0 140%; }
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, 16px); }
}

@keyframes ringSpin {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
}

@keyframes pulseDot {
  0% { box-shadow: 0 0 0 0 rgba(120, 240, 255, 0.55); }
  70% { box-shadow: 0 0 0 10px rgba(120, 240, 255, 0); }
  100% { box-shadow: 0 0 0 0 rgba(120, 240, 255, 0); }
}

@keyframes meterFill {
  from { width: 0; }
  to { width: var(--w); }
}

@keyframes caretBlink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.35; }
}

@keyframes starsTwinkle {
  0%, 100% { opacity: 0.75; }
  50% { opacity: 1; }
}

@media (max-width: 820px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    padding: 28px 24px 20px;
    border-right: none;
    border-bottom: 1px solid var(--line);
  }

  .brand-wtv {
    font-size: 52px;
  }

  .brand-sub {
    font-size: 20px;
  }

  .tech-meter,
  .tech-code {
    display: none;
  }

  .form-panel {
    padding: 24px;
  }
}
</style>
