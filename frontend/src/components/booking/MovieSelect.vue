<script setup>
import { computed } from "vue";

const props = defineProps({
  movies: { type: Array, required: true },
  selectedMovieId: { type: String, required: true },
});

const emit = defineEmits(["pickMovie"]);

const posters = computed(() => props.movies);
</script>

<template>
  <section class="section-block movie-section">
    <div class="section-heading">
      <span class="step-number">02</span>
      <div>
        <h2>เลือกหนัง</h2>
        <p>เลือกเรื่องที่ต้องการจอง วันนี้มี 3 เรื่องแนะนำ</p>
      </div>
    </div>

    <div class="movie-grid">
      <button
        v-for="movie in posters"
        :key="movie.id"
        class="movie-card"
        :class="{ selected: movie.id === selectedMovieId }"
        :style="{ '--accent': movie.accent }"
        @click="$emit('pickMovie', movie.id)"
      >
        <div class="poster" :style="{ background: movie.poster }">
          <div class="poster-glow"></div>
          <strong>{{ movie.title.split(" ")[0] }}</strong>
          <span>{{ movie.rating }}</span>
        </div>
        <div class="movie-copy">
          <p>{{ movie.genre }} · {{ movie.length }}</p>
          <h3>{{ movie.thaiTitle }}</h3>
          <strong>{{ movie.title }}</strong>
        </div>
      </button>
    </div>
  </section>
</template>

<style scoped>
.movie-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.movie-card {
  position: relative;
  overflow: hidden;
  display: grid;
  grid-template-columns: 124px minmax(0, 1fr);
  gap: 16px;
  align-items: stretch;
  min-height: 188px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  padding: 14px;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.105), rgba(255, 255, 255, 0.04)),
    rgba(255, 255, 255, 0.055);
  color: #fff;
  text-align: left;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08);
  transition: border-color 0.18s ease, transform 0.18s ease, background 0.18s ease, box-shadow 0.18s ease;
}

.movie-card:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.22);
  box-shadow: 0 24px 46px rgba(0, 0, 0, 0.28);
}

.movie-card.selected {
  border-color: var(--accent);
  box-shadow: 0 24px 54px rgba(0, 0, 0, 0.28), inset 0 0 0 1px var(--accent);
}

.movie-card.selected::before {
  content: "";
  position: absolute;
  inset: 0;
  border-left: 4px solid var(--accent);
  pointer-events: none;
}

.poster {
  position: relative;
  overflow: hidden;
  border-radius: 6px;
  box-shadow: 0 18px 26px rgba(0, 0, 0, 0.28);
}

.poster::before {
  content: "";
  position: absolute;
  inset: 0;
  background:
    linear-gradient(155deg, rgba(255, 255, 255, 0.28), transparent 34%),
    repeating-linear-gradient(90deg, rgba(255, 255, 255, 0.08) 0 2px, transparent 2px 16px);
  mix-blend-mode: screen;
  opacity: 0.48;
}

.poster::after {
  content: "";
  position: absolute;
  inset: 16px;
  border: 1px solid rgba(255, 255, 255, 0.24);
}

.poster-glow {
  position: absolute;
  width: 80px;
  height: 80px;
  left: 50%;
  top: 44%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.18);
  filter: blur(2px);
}

.poster > strong {
  position: absolute;
  left: 10px;
  right: 10px;
  top: 14px;
  z-index: 1;
  color: rgba(255, 255, 255, 0.86);
  font-size: 0.78rem;
  line-height: 1.1;
  text-transform: uppercase;
}

.poster span {
  position: absolute;
  right: 8px;
  bottom: 8px;
  z-index: 1;
  padding: 4px 7px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.48);
  font-size: 0.76rem;
  font-weight: 800;
}

.movie-copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.movie-copy p {
  color: #f0bb62;
  font-size: 0.78rem;
  font-weight: 800;
}

.movie-copy h3 {
  margin-top: 8px;
  font-size: 1.28rem;
  line-height: 1.22;
}

.movie-copy strong {
  margin-top: 4px;
  color: rgba(255, 255, 255, 0.55);
  font-size: 0.86rem;
}
</style>

