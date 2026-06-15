<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "../composables/useAuth";

const router = useRouter();
const { loginGoogle, loginDev, isFirebaseConfigured } = useAuth();
const error = ref("");
const loading = ref(false);

const handleGoogleLogin = async () => {
  loading.value = true;
  error.value = "";
  try {
    await loginGoogle();
    router.push("/");
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};

const handleAdminLogin = async () => {
  loading.value = true;
  error.value = "";
  try {
    await loginDev("admin001", "admin@test.com");
    router.push("/admin");
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <main class="login-page">
    <div class="login-light login-light-left"></div>
    <div class="login-light login-light-right"></div>

    <button
      class="admin-login-button"
      type="button"
      :disabled="loading"
      @click="handleAdminLogin"
    >
      Admin Login
    </button>

    <section class="login-shell">
      <div class="brand-panel">
        <p class="kicker">CineBook</p>
        <h1>จองตั๋วหนังให้จบในไม่กี่ขั้นตอน</h1>
        <p class="lead">
          เลือกหนัง ดูรอบฉาย แล้วล็อกที่นั่งแบบเรียลไทม์ก่อนยืนยันการจอง
        </p>
        <div class="mini-flow" aria-label="booking steps">
          <span>Login</span>
          <span>Movie</span>
          <span>Showtime</span>
          <span>Seat</span>
        </div>
      </div>

      <div class="login-card">
        <div class="card-heading">
          <p>ยินดีต้อนรับ</p>
          <h2>เข้าสู่ระบบ</h2>
        </div>

        <button
          class="login-button google"
          :disabled="loading || !isFirebaseConfigured"
          @click="handleGoogleLogin"
        >
          <span class="button-icon">G</span>
          <span>{{ loading ? "กำลังเข้าสู่ระบบ..." : "เข้าสู่ระบบด้วย Google" }}</span>
        </button>

        <p v-if="!isFirebaseConfigured" class="dev-note">
          กรุณาตั้งค่า Firebase OAuth ก่อนใช้งานระบบ login ของผู้ใช้
        </p>

        <p v-if="error" class="error">{{ error }}</p>
      </div>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  position: relative;
  overflow: hidden;
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 28px;
  background:
    linear-gradient(rgba(8, 12, 18, 0.58), rgba(8, 12, 18, 0.88)),
    url("../assets/hero.png") center / cover;
}

.login-page::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
  background-size: 84px 100%;
  opacity: 0.24;
}

.login-light {
  position: absolute;
  pointer-events: none;
  border-radius: 999px;
  filter: blur(26px);
}

.login-light-left {
  width: 420px;
  height: 420px;
  left: -160px;
  bottom: 40px;
  background: rgba(217, 83, 51, 0.38);
}

.login-light-right {
  width: 360px;
  height: 360px;
  right: -120px;
  top: 80px;
  background: rgba(240, 187, 98, 0.2);
}

.admin-login-button {
  position: absolute;
  top: 24px;
  right: 28px;
  z-index: 2;
  min-height: 40px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 8px;
  padding: 0 16px;
  background: rgba(10, 14, 20, 0.72);
  color: #fff;
  font-size: 0.9rem;
  font-weight: 900;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.24);
  backdrop-filter: blur(14px);
  transition: transform 0.18s ease, background 0.18s ease, border-color 0.18s ease;
}

.admin-login-button:hover:not(:disabled) {
  transform: translateY(-1px);
  border-color: rgba(240, 187, 98, 0.48);
  background: rgba(20, 28, 38, 0.84);
}

.admin-login-button:disabled {
  cursor: not-allowed;
  opacity: 0.72;
}

.login-shell {
  position: relative;
  z-index: 1;
  width: min(1040px, 100%);
  min-height: 620px;
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) 420px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.11), rgba(255, 255, 255, 0.035)),
    rgba(12, 16, 22, 0.72);
  box-shadow: 0 30px 100px rgba(0, 0, 0, 0.48);
  backdrop-filter: blur(16px);
}

.brand-panel {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 54px;
  color: #fff;
  background:
    linear-gradient(90deg, rgba(8, 12, 18, 0.08), rgba(8, 12, 18, 0.72)),
    radial-gradient(circle at 22% 24%, rgba(240, 187, 98, 0.18), transparent 18rem);
}

.kicker {
  width: max-content;
  padding: 7px 10px;
  border: 1px solid rgba(255, 255, 255, 0.24);
  border-radius: 999px;
  color: #f7c873;
  font-size: 0.78rem;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

h1 {
  max-width: 620px;
  margin-top: 18px;
  font-size: clamp(2.4rem, 5vw, 4.5rem);
  line-height: 1.02;
  letter-spacing: 0;
}

.lead {
  max-width: 520px;
  margin-top: 18px;
  color: rgba(255, 255, 255, 0.76);
  font-size: 1.05rem;
}

.mini-flow {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 34px;
}

.mini-flow span {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.86);
  font-size: 0.84rem;
}

.login-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 42px;
  background: linear-gradient(180deg, #fff8eb, #f3eadc);
  color: #171b22;
}

.card-heading {
  margin-bottom: 26px;
}

.card-heading p {
  color: #9f493c;
  font-size: 0.82rem;
  font-weight: 800;
  text-transform: uppercase;
}

.card-heading h2 {
  margin-top: 4px;
  font-size: 2rem;
}

.login-button {
  min-height: 54px;
  display: flex;
  align-items: center;
  gap: 12px;
  border: 0;
  border-radius: 8px;
  padding: 0 16px;
  font-size: 1rem;
  font-weight: 800;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.14);
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.14);
  filter: brightness(1.03);
}

.login-button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.button-icon {
  width: 30px;
  height: 30px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.18);
  font-weight: 900;
}

.google {
  background: linear-gradient(135deg, #df5b35, #a72d24);
  color: #fff;
}

.dev-note {
  margin-top: 14px;
  color: #6d6a64;
  font-size: 0.88rem;
}

.error {
  margin-top: 18px;
  color: #b3261e;
  font-weight: 700;
}

@media (max-width: 820px) {
  .login-page {
    padding: 76px 16px 16px;
  }

  .admin-login-button {
    top: 16px;
    right: 16px;
  }

  .login-shell {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    min-height: 360px;
    padding: 34px;
  }

  .login-card {
    padding: 30px;
  }
}
</style>
