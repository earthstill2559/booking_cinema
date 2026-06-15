import { ref, onUnmounted } from "vue";
import { WS_URL } from "../lib/api";

export function useWebSocket(onMessage) {
  const connected = ref(false);
  let socket = null;
  let reconnectTimer = null;

  function connect() {
    if (socket) {
      socket.close();
    }

    socket = new WebSocket(WS_URL);

    socket.onopen = () => {
      connected.value = true;
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onMessage(data);
      } catch {
        // ignore malformed messages
      }
    };

    socket.onclose = () => {
      connected.value = false;
      reconnectTimer = setTimeout(connect, 3000);
    };
  }

  connect();

  onUnmounted(() => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
    }
    if (socket) {
      socket.close();
    }
  });

  return { connected };
}
