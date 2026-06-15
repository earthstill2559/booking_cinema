<script setup>
const props = defineProps({
  showtimeOptions: { type: Array, required: true },
  selectedShowtimeId: { type: String, required: true },
});

const emit = defineEmits(["pickShowtime"]);
</script>

<template>
  <section class="section-block showtime-section">
    <div class="section-heading">
      <span class="step-number">03</span>
      <div>
        <h2>ดูเวลาที่หนังฉาย</h2>
        <p>เลือกได้ตามรอบที่ต้องการ</p>
      </div>
    </div>

    <div class="time-grid">
      <button
        v-for="time in showtimeOptions"
        :key="time.id"
        class="time-card"
        :class="{ selected: time.id === selectedShowtimeId }"
        @click="$emit('pickShowtime', time.id)"
      >
        <span>{{ time.label }}</span>
        <strong>{{ time.time }}</strong>
        <p>{{ time.hall }} · {{ time.seats }}</p>
      </button>
    </div>
  </section>
</template>

<style scoped>
.time-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.time-card {
  position: relative;
  min-height: 116px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.065);
  color: #fff;
  text-align: left;
  transition: transform 0.18s ease, background 0.18s ease, border-color 0.18s ease;
}

.time-card:hover {
  transform: translateY(-2px);
  border-color: rgba(240, 187, 98, 0.36);
}

.time-card.selected {
  background: linear-gradient(135deg, #fff8eb, #f0dfc5);
  color: #141820;
  box-shadow: 0 18px 38px rgba(240, 187, 98, 0.14);
}

.time-card.selected::after {
  content: "";
  position: absolute;
  right: 14px;
  top: 14px;
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #d95333;
}

.time-card span {
  color: #f0bb62;
  font-weight: 800;
}

.time-card.selected span {
  color: #ad3527;
}

.time-card strong {
  margin-top: 4px;
  font-size: 2rem;
  line-height: 1;
}

.time-card p {
  margin-top: 10px;
  color: currentColor;
  opacity: 0.72;
  font-size: 0.86rem;
}
</style>

