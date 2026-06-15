<template>
  <div class="seats-wrapper">
    <div class="seats-container">
      <div class="seats-grid">
        <button
          v-for="seat in seats"
          :key="seat.id"
          class="seat"
          :class="seatClass(seat)"
          :disabled="isDisabled(seat)"
          :title="seatTitle(seat)"
          @click="toggleSeatSelection(seat.id)"
        >
          {{ seat.id }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  seats: { type: Array, required: true },
  selectedSeats: { type: Array, required: true },
  currentUserId: { type: String, default: "" },
});

const emit = defineEmits(["toggle-seat"]);

const toggleSeatSelection = (seatId) => {
  emit("toggle-seat", seatId);
};

const seatClass = (seat) => {
  if (seat.status === "BOOKED") return "booked";
  if (seat.status === "LOCKED") {
    if (seat.locked_by === props.currentUserId) return "selected";
    return "locked";
  }
  if (props.selectedSeats.includes(seat.id)) return "selected";
  return "available";
};

const isDisabled = (seat) => {
  if (seat.status === "BOOKED") return true;
  if (seat.status === "LOCKED" && seat.locked_by !== props.currentUserId) return true;
  return false;
};

const seatTitle = (seat) => {
  if (seat.status === "BOOKED") return "ที่นั่งนี้ถูกจองแล้ว";
  if (seat.status === "LOCKED" && seat.locked_by !== props.currentUserId) return "มีคนเลือกอยู่";
  return "เลือกที่นั่ง";
};
</script>

<style scoped>
.seats-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
  margin: 18px 0 0;
}

.seats-container {
  width: 100%;
  max-width: 690px;
}

.seats-grid {
  display: grid;
  grid-template-columns: repeat(8, minmax(0, 1fr));
  gap: 10px;
  width: 100%;
}

.seat {
  width: 100%;
  aspect-ratio: 1;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 7px 7px 3px 3px;
  background: rgba(255, 255, 255, 0.08);
  color: #f7f2e9;
  font-size: 0.82rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.18s ease, border-color 0.18s ease, background 0.18s ease;
}

.seat.available {
  background: #263746;
}

.seat.available:hover {
  transform: translateY(-2px);
  border-color: rgba(240, 187, 98, 0.7);
}

.seat.selected {
  background: #d95333;
  color: white;
  border-color: #f0bb62;
  transform: scale(1.05);
}

.seat.locked {
  background: #c8994d;
  color: #171b22;
  border-color: #c8994d;
  cursor: not-allowed;
  opacity: 0.9;
}

.seat.booked {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.08);
  cursor: not-allowed;
}

.seat:disabled {
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .seats-grid {
    gap: 7px;
  }

  .seat {
    font-size: 0.68rem;
  }
}
</style>
