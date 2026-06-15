<script setup>
import Screen from "../../components/Screen.vue";
import SeatsGrid from "../../components/SeatsGrid.vue";
import Legend from "../../components/Legend.vue";
import BookingSummary from "../../components/BookingSummary.vue";
import Notification from "../../components/Notification.vue";
import BookingActions from "../../components/BookingActions.vue";

const props = defineProps({
  seats: { type: Array, required: true },
  selectedSeats: { type: Array, required: true },
  currentUserId: { type: String, default: "" },
  myLockedSeats: { type: Array, required: true },
  selectedMovieThaiTitle: { type: String, required: true },
  selectedShowtime: { type: Object, required: true },
  displayDate: { type: String, required: true },
  loading: { type: Boolean, required: true },
  message: { type: String, required: true },
  messageType: { type: String, required: true },
});

const emit = defineEmits(["toggleSeat", "confirmBooking", "cancelSelection"]);
</script>

<template>
  <section class="section-block seat-section">
    <div class="section-heading">
      <span class="step-number">04</span>
      <div>
        <h2>เลือกที่นั่ง</h2>
        <p>
          {{ selectedMovieThaiTitle }}
          · รอบ {{ selectedShowtime?.time ?? "-" }}
          · {{ selectedShowtime?.hall ?? "-" }}
        </p>
      </div>
    </div>

    <div class="seat-layout">
      <div class="screen-column">
        <Screen />
        <SeatsGrid
          :seats="seats"
          :selected-seats="selectedSeats"
          :current-user-id="currentUserId"
          @toggle-seat="(seatId) => emit('toggleSeat', seatId)"
        />
      </div>

      <aside class="checkout-panel">
        <div class="ticket-stub">
          <div class="ticket-topline">
            <p>ตั๋วของคุณ</p>
            <span>{{ selectedShowtime?.hall ?? "-" }}</span>
          </div>
          <h3>{{ selectedMovieThaiTitle }}</h3>
          <span>{{ displayDate }} · {{ selectedShowtime?.time ?? "-" }}</span>
        </div>

        <Legend />
        <BookingSummary :selected-seats="myLockedSeats" />
        <Notification :message="message" :message-type="messageType" />

        <BookingActions
          :selected-seats="myLockedSeats"
          :is-loading="loading"
          @confirm-booking="() => emit('confirmBooking')"
          @cancel-selection="() => emit('cancelSelection')"
        />
      </aside>
    </div>
  </section>
</template>

<style scoped>
.seat-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 350px;
  gap: 24px;
  align-items: start;
}

.screen-column {
  min-width: 0;
  border-radius: 8px;
  background:
    linear-gradient(180deg, rgba(8, 12, 18, 0.86), rgba(8, 12, 18, 0.5)),
    rgba(8, 12, 18, 0.54);
  padding: 22px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.06);
}

.checkout-panel {
  position: sticky;
  top: 16px;
  display: grid;
  gap: 14px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  background: rgba(8, 12, 18, 0.5);
}

.ticket-stub {
  position: relative;
  overflow: hidden;
  padding: 20px;
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.58), transparent 38%),
    linear-gradient(135deg, #fff8eb, #ecd4b4);
  color: #171b22;
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.18);
}

.ticket-topline {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.ticket-topline p {
  color: #ad3527;
  font-size: 0.8rem;
  font-weight: 900;
  text-transform: uppercase;
}

.ticket-topline > span {
  margin: 0;
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(173, 53, 39, 0.11);
  color: #8f2e23;
  font-size: 0.76rem;
  font-weight: 900;
}

.ticket-stub h3 {
  margin-top: 16px;
  font-size: 1.42rem;
  line-height: 1.2;
}

.ticket-stub span {
  display: block;
  margin-top: 8px;
  color: #65615a;
}

@media (max-width: 980px) {
  .seat-layout {
    grid-template-columns: 1fr;
  }

  .checkout-panel {
    position: static;
  }
}
</style>

