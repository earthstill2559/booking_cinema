import { ref, computed } from "vue";
import {
  isFirebaseConfigured,
  loginWithGoogle,
  logoutFirebase,
  getIdToken,
  createDevToken,
  getAuthInstance,
  isDevMode,
} from "../lib/firebase";
import { onAuthStateChanged } from "firebase/auth";
import { apiRequest } from "../lib/api";

const user = ref(null);
const token = ref(null);
const role = ref("USER");
const loading = ref(true);

let initialized = false;

export function useAuth() {
  if (!initialized) {
    initialized = true;

    const auth = getAuthInstance();
    console.log(
      "[useAuth] isFirebaseConfigured =",
      isFirebaseConfigured,
      "auth present =",
      Boolean(auth),
      "isDevMode =",
      isDevMode
    );

    // DEV_MODE: bypass Firebase, use localStorage token if available
    if (isDevMode) {
      const saved = localStorage.getItem("cinema_dev_auth");
      if (saved) {
        const parsed = JSON.parse(saved);
        user.value = parsed.user;
        token.value = parsed.token;
        role.value = parsed.role || "USER";
      }
      loading.value = false;
    } else if (isFirebaseConfigured && auth) {
      onAuthStateChanged(auth, async (firebaseUser) => {
        if (firebaseUser) {
          user.value = {
            id: firebaseUser.uid,
            email: firebaseUser.email,
            name: firebaseUser.displayName,
          };

          token.value = await firebaseUser.getIdToken();
          console.log(
            "[useAuth] firebaseUser token set? =",
            Boolean(token.value)
          );

          await fetchProfile();

          localStorage.setItem(
            "cinema_dev_auth",
            JSON.stringify({
              user: user.value,
              token: token.value,
              role: role.value,
            })
          );
        } else {
          user.value = null;
          token.value = null;
          role.value = "USER";
          localStorage.removeItem("cinema_dev_auth");
          console.log("[useAuth] onAuthStateChanged firebaseUser is null");
        }

        loading.value = false;
      });
    } else {
      const saved = localStorage.getItem("cinema_dev_auth");
      if (saved) {
        const parsed = JSON.parse(saved);
        user.value = parsed.user;
        token.value = parsed.token;
        role.value = parsed.role || "USER";
      }
      loading.value = false;
    }
  }

  const isLoggedIn = computed(() => Boolean(user.value && token.value));
  const isAdmin = computed(() => role.value === "ADMIN");

  async function fetchProfile() {
    if (!token.value) return;
    try {
      const profile = await apiRequest("/api/me", { token: token.value });
      role.value = profile.role || "USER";
    } catch {
      role.value = "USER";
    }
  }

  async function loginGoogle() {
    try {
      if (isDevMode) {
        throw new Error(
          "[useAuth] DEV_MODE=true => use loginDev() / local auth token instead of loginGoogle()"
        );
      }

      const auth = getAuthInstance();
      // guard กัน signIn ด้วย auth ที่ไม่พร้อม/ไม่ได้ init
      if (!isFirebaseConfigured || !auth) {
        throw new Error(
          "[useAuth] Firebase auth not initialized (isFirebaseConfigured/auth=false) - set DEV_MODE=true for local auth"
        );
      }

      // (ถ้า auth ถูก init แล้ว แต่ยังไม่มี currentUser ก็ถือว่าเป็นเคส login ใหม่ได้)
      const firebaseUser = await loginWithGoogle();

      user.value = {
        id: firebaseUser.uid,
        email: firebaseUser.email,
        name: firebaseUser.displayName,
      };

      token.value = await getIdToken();
      console.log("[useAuth] loginGoogle() token set? =", Boolean(token.value));
      await fetchProfile();
    } catch (e) {
      console.log(
        "[useAuth] loginGoogle() error =",
        e?.code || e?.name || "unknown",
        "-",
        e?.message || e
      );
      throw e;
    }
  }

  async function loginDev(userId, email) {
    const devToken = createDevToken(userId, email);
    user.value = { id: userId, email, name: email || userId };
    token.value = devToken;
    await fetchProfile();
    localStorage.setItem(
      "cinema_dev_auth",
      JSON.stringify({ user: user.value, token: devToken, role: role.value })
    );
  }

  async function logout() {
    await logoutFirebase();
    user.value = null;
    token.value = null;
    role.value = "USER";
    localStorage.removeItem("cinema_dev_auth");
  }

  return {
    user,
    token,
    role,
    loading,
    isLoggedIn,
    isAdmin,
    isFirebaseConfigured,
    loginGoogle,
    loginDev,
    logout,
    refreshToken: async () => {
      token.value = (await getIdToken()) || token.value;
    },
  };
}
