<template>
  <transition name="fade-slide">
    <div v-if="message" :class="['notification', messageType]">
      <span class="notification-icon">{{ icon }}</span>
      <span>{{ message }}</span>
    </div>
  </transition>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  message: {
    type: String,
    required: true,
  },
  messageType: {
    type: String,
    default: "success",
    validator: (value) => ["success", "error", "warning"].includes(value),
  },
});

const icon = computed(() => {
  if (props.messageType === "success") return "✓";
  if (props.messageType === "error") return "!";
  return "i";
});
</script>

<style scoped>
.notification {
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 13px 14px;
  border-radius: 8px;
  font-size: 0.88rem;
  font-weight: 700;
  line-height: 1.45;
}

.notification-icon {
  width: 22px;
  height: 22px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.16);
  font-size: 0.76rem;
}

.success {
  background: rgba(58, 143, 123, 0.18);
  color: #a6e3d1;
}

.error {
  background: rgba(185, 58, 58, 0.18);
  color: #ffb0a6;
}

.warning {
  background: rgba(240, 187, 98, 0.18);
  color: #f5d49b;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
