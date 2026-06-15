<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import Header from "../components/Header.vue";
import { useAuth } from "../composables/useAuth";
import { useWebSocket } from "../composables/useWebSocket";
import { apiRequest } from "../lib/api";

import MovieSelect from "../components/booking/MovieSelect.vue";
import ShowtimeSelect from "../components/booking/ShowtimeSelect.vue";
import SeatsSelection from "../components/booking/SeatsSelection.vue";
import BookingStepsHeader from "../components/booking/BookingStepsHeader.vue";

const router = useRouter();
const { user, isLoggedIn, isAdmin, logout, loading: authLoading } = useAuth();

const movies = ref([]);

const seats = ref([]);
const showtime = ref(null);
const selectedSeats = ref([]);
const selectedMovieId = ref("");
const selectedShowtimeId = ref("");
const currentStep = ref(2);
const loading = ref(false);
const message = ref("");
const messageType = ref("success");

const selectedMovie = computed(() => movies.value.find((movie) => movie.id === selectedMovieId.value) || movies.value[0]);

// แต่ละหนัง = 1 showtime (1 รอบ) ตามโครงสร้างที่ admin สร้าง
const showtimeOptions = computed(() => {
  const movie = selectedMovie.value;
  if (!movie?._showtime) return [];
  const item = movie._showtime;
  return [
    {
      id: item.id,
      time: item.time || "",
      label: item.date || "",
      hall: "",
      seats: item.total_seats ? `${item.total_seats} ที่นั่ง` : "",
    },
  ];
});

const selectedShowtime = computed(
  () => showtimeOptions.value.find((item) => item.id === selectedShowtimeId.value) || showtimeOptions.value[0]
);

const flowSteps = computed(() => [
  { number: 1, label: "Login", active: true },
  { number: 2, label: "เลือกหนัง", active: currentStep.value >= 2 },
  { number: 3, label: "รอบฉาย", active: currentStep.value >= 3 },
  { number: 4, label: "ที่นั่ง", active: currentStep.value >= 4 },
]);

const displayDate = computed(() => showtime.value?.date || "วันนี้");

const myLockedSeats = computed(() =>
  selectedSeats.value.filter((seatId) => {
    const seat = seats.value.find((s) => s.id === seatId);
    return seat?.status === "LOCKED" && seat?.locked_by === user.value?.id;
  })
);

// ดึงรายการหนังทั้งหมดที่ admin สร้างไว้
const fetchShowtimes = async () => {
  const data = await apiRequest("/api/showtimes"); // array ของ { id, movie, date, time, total_seats }

  movies.value = (data || []).map((item) => ({
    id: item.id,
    title: item.movie,
    thaiTitle: item.movie,
    genre: "",
    length: "",
    rating: "",
    accent: "#df5b35",
    poster: "linear-gradient(145deg, #1b2735 0%, #090a0f 54%, #df5b35 100%)",
    _showtime: item,
  }));

  if (movies.value.length > 0) {
    selectedMovieId.value = movies.value[0].id;
    showtime.value = movies.value[0]._showtime;
    selectedShowtimeId.value = movies.value[0]._showtime?.id || "";
  } else {
    selectedMovieId.value = "";
    showtime.value = null;
    selectedShowtimeId.value = "";
  }
};

const fetchSeats = async () => {
  const showtimeId = showtime.value?.id || selectedShowtimeId.value || "show_001";
  seats.value = await apiRequest(`/api/seats?showtime_id=${showtimeId}`);
  selectedSeats.value = selectedSeats.value.filter((id) => {
    const seat = seats.value.find((s) => s.id === id);
    return seat && seat.status === "LOCKED" && seat.locked_by === user.value?.id;
  });
};

const handleSeatUpdate = (data) => {
  if (data.type !== "seat_update") return;
  const index = seats.value.findIndex((s) => s.id === data.seat_id);
  if (index === -1) return;

  seats.value[index] = {
    ...seats.value[index],
    status: data.status,
    locked_by: data.locked_by || "",
  };

  if (data.status !== "LOCKED" || data.locked_by !== user.value?.id) {
    selectedSeats.value = selectedSeats.value.filter((id) => id !== data.seat_id);
  }
};

useWebSocket(handleSeatUpdate);

const pickMovie = async (movieId) => {
  selectedMovieId.value = movieId;
  const movie = movies.value.find((m) => m.id === movieId);
  showtime.value = movie?._showtime || null;
  selectedShowtimeId.value = movie?._showtime?.id || "";
  currentStep.value = Math.max(currentStep.value, 3);
  selectedSeats.value = [];

  try {
    await fetchSeats();
  } catch (error) {
    message.value = error.message;
    messageType.value = "error";
  }
};

const pickShowtime = async (showtimeId) => {
  selectedShowtimeId.value = showtimeId;
  currentStep.value = 4;
  selectedSeats.value = [];

  try {
    await fetchSeats();
  } catch (error) {
    message.value = error.message;
    messageType.value = "error";
  }
};

const toggleSeatSelection = async (seatId) => {
  if (!isLoggedIn.value) {
    router.push("/login");
    return;
  }

  const seat = seats.value.find((s) => s.id === seatId);
  if (!seat || seat.status === "BOOKED") return;

  if (seat.status === "LOCKED" && seat.locked_by !== user.value?.id) {
    message.value = "ที่นั่งนี้มีคนเลือกอยู่";
    messageType.value = "warning";
    return;
  }

  const isSelected = selectedSeats.value.includes(seatId);

  loading.value = true;
  message.value = "";

  try {
    if (isSelected) {
      await apiRequest(`/api/seats/${seatId}/unlock`, {
        method: "POST",
        body: { showtime_id: showtime.value?.id },
      });
      selectedSeats.value = selectedSeats.value.filter((id) => id !== seatId);
    } else {
      await apiRequest(`/api/seats/${seatId}/lock`, {
        method: "POST",
        body: { showtime_id: showtime.value?.id },
      });
      selectedSeats.value.push(seatId);
    }
    await fetchSeats();
  } catch (error) {
    message.value = error.message;
    messageType.value = "error";
  } finally {
    loading.value = false;
  }
};

const confirmBooking = async () => {
  if (myLockedSeats.value.length === 0) {
    message.value = "กรุณาเลือกที่นั่งก่อนยืนยัน";
    messageType.value = "warning";
    return;
  }

  loading.value = true;
  message.value = "";

  try {
    await apiRequest("/api/bookings/confirm", {
      method: "POST",
      body: {
        showtime_id: showtime.value?.id,
        seat_ids: myLockedSeats.value,
        movie: selectedMovie.value.title,
        date: showtime.value?.date,
      },
    });

    message.value = `จองสำเร็จ: ${selectedMovie.value.thaiTitle} รอบ ${selectedShowtime.value.time} ที่นั่ง ${myLockedSeats.value.join(", ")}`;
    messageType.value = "success";
    selectedSeats.value = [];
    await fetchSeats();
  } catch (error) {
    message.value = error.message;
    messageType.value = "error";
  } finally {
    loading.value = false;
  }
};

const cancelSelection = async () => {
  loading.value = true;
  try {
    for (const seatId of [...selectedSeats.value]) {
      await apiRequest(`/api/seats/${seatId}/unlock`, {
        method: "POST",
        body: { showtime_id: showtime.value?.id },
      });
    }
    selectedSeats.value = [];
    message.value = "";
    await fetchSeats();
  } catch (error) {
    message.value = error.message;
    messageType.value = "error";
  } finally {
    loading.value = false;
  }
};

const handleLogout = async () => {
  await logout();
  router.push("/login");
};

onMounted(async () => {
  while (authLoading.value) {
    await new Promise((r) => setTimeout(r, 50));
  }
  if (!isLoggedIn.value) {
    router.push("/login");
    return;
  }

  try {
    await fetchShowtimes();
    await fetchSeats();
  } catch {
    message.value = "ไม่สามารถโหลดข้อมูลที่นั่งได้ กรุณาตรวจสอบ backend";
    messageType.value = "error";
  }
});
</script>

<template>
  <div class="booking-page">
    <div class="ambient ambient-one"></div>
    <div class="ambient ambient-two"></div>

    <Header
      :user="user"
      :is-admin="isAdmin"
      @logout="handleLogout"
      @go-admin="router.push('/admin')"
    />

    <main class="booking-shell">
      <BookingStepsHeader :flow-steps="flowSteps" />

      <MovieSelect
        :movies="movies"
        :selected-movie-id="selectedMovieId"
        @pickMovie="pickMovie"
      />

      <ShowtimeSelect
        :showtime-options="showtimeOptions"
        :selected-showtime-id="selectedShowtimeId"
        @pickShowtime="pickShowtime"
      />

      <SeatsSelection
        :seats="seats"
        :selected-seats="selectedSeats"
        :current-user-id="user?.id"
        :my-locked-seats="myLockedSeats"
:selected-movie-thai-title="selectedMovie?.thaiTitle || ''"
        :selected-showtime="selectedShowtime"
        :display-date="displayDate"
        :loading="loading"
        :message="message"
        :message-type="messageType"
        @toggle-seat="toggleSeatSelection"
        @confirm-booking="confirmBooking"
        @cancel-selection="cancelSelection"
      />
    </main>
  </div>
</template>

<style scoped>
.booking-page {
  position: relative;
  overflow: hidden;
  min-height: 100vh;
  color: #f7f2e9;
  background:
    linear-gradient(rgba(11, 16, 23, 0.88), rgba(11, 16, 23, 0.95)),
    url("../assets/hero.png") top center / cover fixed,
    linear-gradient(180deg, #121822 0%, #0d1117 54%, #17120f 100%);
}

.booking-page::before {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.035) 1px, transparent 1px),
    linear-gradient(rgba(255, 255, 255, 0.025) 1px, transparent 1px);
  background-size: 72px 72px;
  mask-image: linear-gradient(to bottom, black, transparent 78%);
}

.ambient {
  position: fixed;
  pointer-events: none;
  border-radius: 999px;
  filter: blur(24px);
  opacity: 0.42;
}

.ambient-one {
  width: 420px;
  height: 420px;
  left: -150px;
  top: 110px;
  background: rgba(217, 83, 51, 0.42);
}

.ambient-two {
  width: 320px;
  height: 320px;
  right: -110px;
  top: 360px;
  background: rgba(58, 143, 123, 0.28);
}

.booking-shell {
  position: relative;
  z-index: 1;
  width: min(1180px, calc(100% - 32px));
  margin: 0 auto;
  padding: 34px 0 58px;
}

@media (max-width: 640px) {
  .booking-shell {
    width: min(100% - 20px, 1180px);
    padding-top: 16px;
  }
}
</style>

