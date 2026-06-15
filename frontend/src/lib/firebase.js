import { initializeApp } from "firebase/app";
import {
  getAuth,
  GoogleAuthProvider,
  signInWithPopup,
  signOut,
} from "firebase/auth";

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

// DEV_MODE: ให้ frontend bypass Firebase เพื่อใช้ local/dev auth
const isDevMode =
  String(import.meta.env.VITE_DEV_MODE ?? "").toLowerCase() === "true";

// DEBUG: ดูว่า Vite inject env เข้า runtime ได้จริงไหม
console.log("[firebase.js] VITE_DEV_MODE =", import.meta.env.VITE_DEV_MODE);
console.log(
  "[firebase.js] import.meta.env.VITE_FIREBASE_PROJECT_ID =",
  import.meta.env.VITE_FIREBASE_PROJECT_ID
);
console.log(
  "[firebase.js] import.meta.env.VITE_FIREBASE_API_KEY (prefix) =",
  (import.meta.env.VITE_FIREBASE_API_KEY || "").slice(0, 12)
);
console.log(
  "[firebase.js] import.meta.env.VITE_FIREBASE_AUTH_DOMAIN =",
  import.meta.env.VITE_FIREBASE_AUTH_DOMAIN
);
console.log(
  "[firebase.js] import.meta.env.VITE_FIREBASE_APP_ID (prefix) =",
  (import.meta.env.VITE_FIREBASE_APP_ID || "").slice(0, 12)
);

const isFirebaseConfigured = !isDevMode && Boolean(
  firebaseConfig.apiKey &&
    firebaseConfig.authDomain &&
    firebaseConfig.projectId &&
    firebaseConfig.appId
);

let app = null;
let authInstance = null;

if (isFirebaseConfigured) {
  app = initializeApp(firebaseConfig);
} else {
  if (isDevMode) {
    console.warn(
      "[firebase.js] DEV_MODE=true => bypass Firebase initialization (ใช้ local/dev auth แทน)"
    );
  } else {
    console.warn(
      "[firebase.js] Firebase config ไม่ครบ - ตรวจสอบ env vars (apiKey/authDomain/projectId/appId)"
    );
  }
}

function getAuthInstance() {
  if (!app) return null;
  // singleton: สร้าง auth instance ครั้งแรกครั้งเดียว แล้ว reuse ตลอด
  if (!authInstance) {
    authInstance = getAuth(app);
  }
  return authInstance;
}

export { isFirebaseConfigured, getAuthInstance, isDevMode };

export async function loginWithGoogle() {
  const auth = getAuthInstance();
  if (!auth) {
    throw new Error("Firebase is not configured");
  }
  const provider = new GoogleAuthProvider();
  const result = await signInWithPopup(auth, provider);
  return result.user;
}

export async function logoutFirebase() {
  const auth = getAuthInstance();
  if (auth?.currentUser) {
    await signOut(auth);
  }
}

export async function getIdToken() {
  const auth = getAuthInstance();
  if (!auth?.currentUser) return null;
  return auth.currentUser.getIdToken();
}

export function createDevToken(userId, email = "") {
  return `dev:${userId}:${email}`;
}

