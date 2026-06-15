<template>
  <div class="booking-summary">
    <div class="summary-item">
      <span>ที่นั่งที่เลือก</span>
      <strong>{{ selectedSeats.length ? selectedSeats.join(", ") : "ยังไม่ได้เลือก" }}</strong>
    </div>
    <div class="summary-item">
      <span>จำนวน</span>
      <strong>{{ selectedSeats.length }} ที่นั่ง</strong>
    </div>
  </div>
</template>

<script setup>
defineProps({
  selectedSeats: {
    type: Array,
    required: true,
  },
});
</script>

<style scoped>
.booking-summary {
  display: grid;
  gap: 12px;
  width: 100%;
  padding: 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.07);
}

.summary-item {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.summary-item:last-child {
  padding-bottom: 0;
  border-bottom: 0;
}

.summary-item span {
  color: rgba(247, 242, 233, 0.62);
  font-size: 0.84rem;
}

.summary-item strong {
  color: #fff;
  font-size: 1rem;
  text-align: right;
}
</style>
