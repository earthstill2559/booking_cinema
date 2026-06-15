<template>
  <div class="screen-section">
    <div class="screen-label">SCREEN</div>
    <div class="screen-bar"></div>
  </div>
</template>

<script setup></script>

<style scoped>
.screen-section {
  width: 100%;
  text-align: center;
}

.screen-label {
  color: rgba(247, 242, 233, 0.52);
  font-size: 0.75rem;
  font-weight: 900;
  letter-spacing: 0.18em;
}

.screen-bar {
  width: min(520px, 100%);
  height: 42px;
  margin: 12px auto 26px;
  border-top: 4px solid rgba(240, 187, 98, 0.78);
  border-radius: 50% 50% 0 0;
  background: linear-gradient(rgba(255, 255, 255, 0.12), transparent);
  filter: drop-shadow(0 12px 22px rgba(240, 187, 98, 0.16));
}
</style>
