<template>
  <div class="legend">
    <div class="legend-item">
      <div class="legend-box available"></div>
      <span>ว่าง</span>
    </div>
    <div class="legend-item">
      <div class="legend-box locked"></div>
      <span>มีคนเลือกอยู่</span>
    </div>
    <div class="legend-item">
      <div class="legend-box selected"></div>
      <span>ที่นั่งของคุณ</span>
    </div>
    <div class="legend-item">
      <div class="legend-box booked"></div>
      <span>จองแล้ว</span>
    </div>
  </div>
</template>

<style scoped>
.legend {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  width: 100%;
  padding: 14px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.07);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(247, 242, 233, 0.78);
  font-size: 0.82rem;
}

.legend-box {
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  border-radius: 5px 5px 2px 2px;
}

.legend-box.available {
  background: #263746;
}

.legend-box.locked {
  background: #c8994d;
}

.legend-box.selected {
  background: #d95333;
}

.legend-box.booked {
  background: rgba(255, 255, 255, 0.16);
}
</style>
