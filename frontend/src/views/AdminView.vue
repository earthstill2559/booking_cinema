<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import Header from "../components/Header.vue";
import { useAuth } from "../composables/useAuth";
import { apiRequest } from "../lib/api";

const router = useRouter();
const { user, token, isLoggedIn, isAdmin, logout, loading: authLoading } = useAuth();

const bookings = ref([]);
const loading = ref(false);
const error = ref("");
const filters = ref({
  movie: "",
  date: "",
  user_id: "",
});

// --- เพิ่มรอบหนังใหม่ ---
const newShowtime = ref({
  movie: "",
  date: "",
  time: "",
  total_seats: 50,
});
const createLoading = ref(false);
const createError = ref("");
const createSuccess = ref("");

const createShowtime = async () => {
  createError.value = "";
  createSuccess.value = "";

  if (!newShowtime.value.movie || !newShowtime.value.date || !newShowtime.value.time) {
    createError.value = "กรุณากรอกชื่อหนัง วันที่ และเวลาให้ครบ";
    return;
  }
  if (!newShowtime.value.total_seats || newShowtime.value.total_seats < 1) {
    createError.value = "จำนวนที่นั่งต้องมากกว่า 0";
    return;
  }

  createLoading.value = true;
  try {
    await apiRequest("/api/admin/showtimes", {
      method: "POST",
      token: token.value,
      body: {
        movie: newShowtime.value.movie,
        date: newShowtime.value.date,
        time: newShowtime.value.time,
        total_seats: newShowtime.value.total_seats,
      },
    });

    createSuccess.value = "สร้างรอบหนังสำเร็จแล้ว";
    newShowtime.value = { movie: "", date: "", time: "", total_seats: 50 };

    await fetchBookings();
  } catch (e) {
    createError.value = e.message;
  } finally {
    createLoading.value = false;
  }
};
// --- จบส่วนเพิ่มรอบหนังใหม่ ---

const fetchBookings = async () => {
  loading.value = true;
  error.value = "";
  try {
    const params = new URLSearchParams();
    if (filters.value.movie) params.set("movie", filters.value.movie);
    if (filters.value.date) params.set("date", filters.value.date);
    if (filters.value.user_id) params.set("user_id", filters.value.user_id);

    const query = params.toString();
    bookings.value = await apiRequest(`/api/admin/bookings${query ? `?${query}` : ""}`, {
      token: token.value,
    });
  } catch (e) {
    error.value = e.message;
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
  if (!isAdmin.value) {
    router.push("/");
    return;
  }
  await fetchBookings();
});
</script>

<template>
  <div class="admin-page">
    <Header
      :user="user"
      :is-admin="isAdmin"
      show-back
      @logout="handleLogout"
      @go-back="router.push('/')"
    />

    <main class="content">
      <section class="admin-hero">
        <p class="eyebrow">Admin Console</p>
        <h1>ตรวจสอบรายการจอง</h1>
        <p class="subtitle">ดูรายการจองทั้งหมด พร้อมตัวกรองสำหรับค้นหาตามหนัง วันที่ หรือผู้ใช้</p>
      </section>

      <!-- เพิ่มรอบหนังใหม่ -->
      <section class="create-card" aria-label="create showtime">
        <h2>เพิ่มรอบหนังใหม่</h2>
        <div class="create-form">
          <input v-model="newShowtime.movie" placeholder="ชื่อหนัง" />
          <input v-model="newShowtime.date" type="date" />
          <input v-model="newShowtime.time" type="time" />
          <input
            v-model.number="newShowtime.total_seats"
            type="number"
            min="1"
            placeholder="จำนวนที่นั่ง"
          />
          <button type="button" :disabled="createLoading" @click="createShowtime">
            {{ createLoading ? "กำลังบันทึก..." : "สร้างรอบหนัง" }}
          </button>
        </div>
        <p v-if="createError" class="error">{{ createError }}</p>
        <p v-if="createSuccess" class="success">{{ createSuccess }}</p>
      </section>

      <section class="filters" aria-label="booking filters">
        <input v-model="filters.movie" placeholder="ชื่อหนัง" />
        <input v-model="filters.date" placeholder="วันที่ เช่น 2026-06-13" />
        <input v-model="filters.user_id" placeholder="User ID" />
        <button type="button" :disabled="loading" @click="fetchBookings">
          {{ loading ? "กำลังค้นหา..." : "ค้นหา" }}
        </button>
      </section>

      <p v-if="error" class="error">{{ error }}</p>

      <section class="table-card">
        <table>
          <thead>
            <tr>
              <th>Booking ID</th>
              <th>User ID</th>
              <th>Email</th>
              <th>Seat</th>
              <th>Movie</th>
              <th>Date</th>
              <th>Status</th>
              <th>Created</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8">Loading...</td>
            </tr>
            <tr v-else-if="bookings.length === 0">
              <td colspan="8">No bookings found</td>
            </tr>
            <tr v-for="b in bookings" :key="b.id">
              <td>{{ b.id.slice(0, 8) }}...</td>
              <td>{{ b.user_id }}</td>
              <td>{{ b.user_email || "-" }}</td>
              <td>{{ b.seat_id }}</td>
              <td>{{ b.movie }}</td>
              <td>{{ b.date }}</td>
              <td><span class="badge">{{ b.status }}</span></td>
              <td>{{ new Date(b.created_at).toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>
      </section>
    </main>
  </div>
</template>

<style scoped>
.admin-page {
  min-height: 100vh;
  color: #f7f2e9;
  background:
    radial-gradient(circle at top left, rgba(217, 83, 51, 0.2), transparent 32rem),
    linear-gradient(180deg, #111820 0%, #0d1117 100%);
}

.content {
  width: min(1200px, calc(100% - 32px));
  margin: 0 auto;
  padding: 34px 0 54px;
}

.admin-hero {
  margin-bottom: 24px;
}

.eyebrow {
  color: #f0bb62;
  font-size: 0.78rem;
  font-weight: 900;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

h1 {
  margin: 8px 0 0;
  font-size: clamp(2rem, 4vw, 3.2rem);
  line-height: 1.08;
}

.subtitle {
  max-width: 660px;
  color: rgba(247, 242, 233, 0.66);
  margin-top: 8px;
}

/* เพิ่มรอบหนังใหม่ */
.create-card {
  margin-bottom: 22px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(12px);
}

.create-card h2 {
  margin: 0 0 12px;
  font-size: 1.1rem;
  font-weight: 900;
}

.create-form {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.create-form input {
  min-height: 42px;
  min-width: 160px;
  flex: 1 1 160px;
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  background: rgba(8, 12, 18, 0.52);
  color: #fff;
}

.create-form input::placeholder {
  color: rgba(247, 242, 233, 0.42);
}

.create-form button {
  min-height: 42px;
  padding: 10px 20px;
  border: 0;
  border-radius: 8px;
  background: #d95333;
  color: #fff;
  font-weight: 900;
  white-space: nowrap;
}

.create-form button:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

.success {
  margin-top: 10px;
  margin-bottom: 0;
  color: #8be3b0;
  font-weight: 800;
}

.filters {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 22px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(12px);
}

.filters input {
  min-height: 42px;
  min-width: 180px;
  flex: 1 1 180px;
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  background: rgba(8, 12, 18, 0.52);
  color: #fff;
}

.filters input::placeholder {
  color: rgba(247, 242, 233, 0.42);
}

.filters button {
  min-height: 42px;
  padding: 10px 20px;
  border: 0;
  border-radius: 8px;
  background: #d95333;
  color: #fff;
  font-weight: 900;
}

.filters button:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

.table-card {
  overflow-x: auto;
  border-radius: 8px;
  background: #f8f4ed;
  color: #171b22;
  box-shadow: 0 22px 70px rgba(0, 0, 0, 0.22);
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 13px 14px;
  text-align: left;
  border-bottom: 1px solid #e4d8c8;
  font-size: 0.9rem;
  white-space: nowrap;
}

th {
  background: #eee1ce;
  font-weight: 900;
}

.badge {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 2px 9px;
  border-radius: 999px;
  background: #dff2e7;
  color: #1f7547;
  font-size: 0.78rem;
  font-weight: 900;
}

.error {
  margin-bottom: 14px;
  color: #ffb0a6;
  font-weight: 800;
}

@media (max-width: 640px) {
  .content {
    width: min(100% - 20px, 1200px);
    padding-top: 24px;
  }

  .filters,
  .create-card {
    padding: 12px;
  }

  .filters button,
  .create-form button {
    width: 100%;
  }
}
</style>