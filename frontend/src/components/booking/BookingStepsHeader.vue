<script setup>
const props = defineProps({
  flowSteps: { type: Array, required: true },
});
</script>

<template>
  <section class="hero-strip">
    <div>
      <p class="eyebrow">Cinema Booking</p>
      <h1>เลือกหนัง รอบฉาย และที่นั่งในหน้าเดียว</h1>
      <p class="hero-copy">
        ระบบจองที่ออกแบบให้เห็นสถานะครบในจอเดียว ตั้งแต่เลือกเรื่องไปจนถึงล็อกที่นั่งแบบเรียลไทม์
      </p>
    </div>

    <div class="flow">
      <div
        v-for="step in flowSteps"
        :key="step.number"
        class="flow-step"
        :class="{ active: step.active }"
      >
        <span>{{ step.number }}</span>
        <p>{{ step.label }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero-strip {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 24px;
  align-items: end;
  padding: 34px 0 28px;
}

.eyebrow {
  color: #f0bb62;
  font-size: 0.78rem;
  font-weight: 900;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero-strip h1 {
  max-width: 700px;
  margin-top: 8px;
  font-size: clamp(2rem, 4.2vw, 4rem);
  line-height: 1.04;
  text-wrap: balance;
}

.hero-copy {
  max-width: 620px;
  margin-top: 14px;
  color: rgba(247, 242, 233, 0.68);
  font-size: 1rem;
}

.flow {
  display: grid;
  grid-template-columns: repeat(4, minmax(74px, 1fr));
  gap: 8px;
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 999px;
  background: rgba(8, 12, 18, 0.54);
  box-shadow: 0 18px 50px rgba(0, 0, 0, 0.24);
  backdrop-filter: blur(12px);
}

.flow-step {
  display: grid;
  justify-items: center;
  gap: 7px;
  color: rgba(247, 242, 233, 0.42);
  font-size: 0.78rem;
}

.flow-step span {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.18);
  background: rgba(255, 255, 255, 0.06);
  font-weight: 900;
}

.flow-step.active {
  color: #fff;
}

.flow-step.active span {
  background: linear-gradient(135deg, #e66b43, #b72f24);
  border-color: rgba(255, 255, 255, 0.28);
  box-shadow: 0 10px 24px rgba(217, 83, 51, 0.28);
}

@media (max-width: 980px) {
  .hero-strip {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .hero-strip {
    padding: 16px 0;
  }
  .flow {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>

