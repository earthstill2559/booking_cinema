const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
const WS_URL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws";

export { API_URL, WS_URL };


function getStoredToken() {

  try {
    const saved = localStorage.getItem("cinema_dev_auth");
    if (!saved) return null;
    const parsed = JSON.parse(saved);
    return parsed?.token || null;
  } catch {
    return null;
  }
}

export async function apiRequest(path, { method = "GET", token, body } = {}) {
  const headers = { "Content-Type": "application/json" };

  // ให้ Authorization อัตโนมัติ ถ้า caller ไม่ส่ง token มา
  const effectiveToken = token || getStoredToken();
  if (effectiveToken) {
    headers.Authorization = `Bearer ${effectiveToken}`;
  }


  const res = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || `Request failed (${res.status})`);
  }
  return data;
}
