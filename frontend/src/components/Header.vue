<template>
  <header class="app-header">
    <div class="header-inner">
      <button
        v-if="showBack"
        class="icon-button back-button"
        type="button"
        aria-label="กลับไปหน้าจอง"
        title="กลับไปหน้าจอง"
        @click="$emit('go-back')"
      >
        <span aria-hidden="true">‹</span>
        <strong>กลับ</strong>
      </button>

      <div class="brand">
        <span class="brand-mark">CB</span>
        <div class="brand-copy">
          <strong>CineBook</strong>
          <p>Real-time cinema booking</p>
        </div>
      </div>

      <div class="header-actions">
        <button
          v-if="isAdmin && !showBack"
          class="header-button"
          type="button"
          @click="$emit('go-admin')"
        >
          Admin
        </button>
        <span v-if="user" class="user-chip">{{ user.name || user.email || user.id }}</span>
        <button v-if="user" class="header-button ghost" type="button" @click="$emit('logout')">
          Logout
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  user: { type: Object, default: null },
  isAdmin: { type: Boolean, default: false },
  showBack: { type: Boolean, default: false },
});

defineEmits(["logout", "go-admin", "go-back"]);
</script>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 20;
  width: 100%;
  min-height: 76px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.11);
  background:
    linear-gradient(135deg, rgba(20, 28, 38, 0.94), rgba(10, 14, 20, 0.92)),
    rgba(10, 14, 20, 0.94);
  box-shadow: 0 14px 40px rgba(0, 0, 0, 0.24);
  backdrop-filter: blur(16px);
}

.header-inner {
  width: min(1180px, calc(100% - 32px));
  min-height: 76px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: auto minmax(220px, 1fr) auto;
  align-items: center;
  gap: 16px;
}

.brand {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #fff;
}

.brand-mark {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 8px;
  background: linear-gradient(135deg, #df5b35, #a72d24);
  color: #fff;
  font-size: 0.9rem;
  font-weight: 900;
  box-shadow: 0 12px 24px rgba(217, 83, 51, 0.26);
}

.brand-copy {
  min-width: 0;
}

.brand-copy strong {
  display: block;
  color: #fff;
  font-size: 1.08rem;
  line-height: 1.15;
}

.brand-copy p {
  margin-top: 3px;
  color: rgba(247, 242, 233, 0.58);
  font-size: 0.78rem;
  white-space: nowrap;
}

.header-actions {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.user-chip {
  max-width: 240px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 8px 11px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.06);
  color: rgba(247, 242, 233, 0.82);
  font-size: 0.85rem;
  font-weight: 700;
}

.header-button,
.icon-button {
  min-height: 38px;
  border: 1px solid rgba(255, 255, 255, 0.13);
  border-radius: 8px;
  padding: 0 14px;
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  font-weight: 900;
  transition: transform 0.18s ease, background 0.18s ease, border-color 0.18s ease;
}

.header-button:hover,
.icon-button:hover {
  transform: translateY(-1px);
  border-color: rgba(240, 187, 98, 0.38);
  background: rgba(255, 255, 255, 0.12);
}

.ghost {
  background: transparent;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.back-button span {
  font-size: 1.5rem;
  line-height: 0;
  transform: translateY(-1px);
}

@media (max-width: 700px) {
  .app-header {
    min-height: 0;
  }

  .header-inner {
    width: min(100% - 20px, 1180px);
    min-height: 0;
    grid-template-columns: 1fr auto;
    padding: 12px 0;
  }

  .back-button {
    grid-column: 1 / -1;
    width: max-content;
  }

  .brand-copy p {
    white-space: normal;
  }

  .header-actions {
    justify-content: flex-end;
  }

  .user-chip {
    display: none;
  }
}
</style>
