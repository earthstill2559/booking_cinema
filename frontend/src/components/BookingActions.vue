<template>
  <div class="actions">
    <button
      v-if="selectedSeats.length > 0"
      class="btn secondary"
      :disabled="isLoading"
      @click="cancelSelection"
    >
      ยกเลิกที่เลือก
    </button>
    <button
      class="btn primary"
      :disabled="isLoading || selectedSeats.length === 0"
      @click="confirmBooking"
    >
      <span v-if="!isLoading">ยืนยันการจอง</span>
      <span v-else class="loading-text">
        <span class="spinner"></span>
        กำลังจอง...
      </span>
    </button>
  </div>
</template>

<script setup>
defineProps({
  selectedSeats: {
    type: Array,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["confirm-booking", "cancel-selection"]);

const confirmBooking = () => {
  emit("confirm-booking");
};

const cancelSelection = () => {
  emit("cancel-selection");
};
</script>

<style scoped>
.actions {
  display: grid;
  gap: 10px;
  width: 100%;
}

.btn {
  min-height: 48px;
  border: 0;
  border-radius: 8px;
  padding: 0 18px;
  font-size: 0.96rem;
  font-weight: 900;
  transition: transform 0.18s ease, opacity 0.18s ease;
}

.btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.primary {
  background: #d95333;
  color: white;
}

.secondary {
  background: rgba(255, 255, 255, 0.08);
  color: #f7f2e9;
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.loading-text {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
