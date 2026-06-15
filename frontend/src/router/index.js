import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import BookingView from "../views/BookingView.vue";
import AdminView from "../views/AdminView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "booking", component: BookingView },
    { path: "/login", name: "login", component: LoginView },
    { path: "/admin", name: "admin", component: AdminView, meta: { requiresAdmin: true } },
  ],
});

router.beforeEach(async (to, _from, next) => {
  const { useAuth } = await import("../composables/useAuth");
  const { isLoggedIn, isAdmin, loading } = useAuth();

  // ไม่ให้ UI ค้างถ้า auth/firebase ไม่ callback เสร็จ
  if (loading.value) {
    const timeoutMs = 5000;
    await Promise.race([
      new Promise((resolve) => {
        const stop = setInterval(() => {
          if (!loading.value) {
            clearInterval(stop);
            resolve();
          }
        }, 50);
      }),
      new Promise((resolve) => setTimeout(resolve, timeoutMs)),
    ]);
  }

  if (to.meta.requiresAdmin) {
    if (!isLoggedIn.value) return next("/login");
    if (!isAdmin.value) return next("/");
  }

  next();
});

export default router;
