import React, {
  createContext,
  useContext,
  useRef,
  useState,
  useCallback,
  useEffect,
} from "react";
import { ConnectionState } from "../types";
import type { WebSocketContextType, Message } from "../types";
import { apiService } from "../services/api";

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};

interface WebSocketProviderProps {
  children: React.ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  children,
}) => {
  const ws = useRef<WebSocket | null>(null);
  const [connectionState, setConnectionState] = useState<ConnectionState>(
    ConnectionState.DISCONNECTED
  );
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(
    null
  );
  const reconnectAttemptsRef = useRef(0);
  const shouldReconnectRef = useRef(true);
  const currentRoomRef = useRef<string | null>(null);
  const messageQueueRef = useRef<Partial<Message>[]>([]);
  const presenceHeartbeatRef = useRef<ReturnType<typeof setInterval> | null>(
    null
  );
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const lastPongRef = useRef<number>(Date.now());

  const maxReconnectAttempts = 10;
  const baseReconnectDelay = 1000;
  const maxReconnectDelay = 30000;
  const heartbeatInterval = 120000; // 2 minutes
  const pingInterval = 30000; // 30 seconds - ping to keep connection alive

  const getWebSocketUrl = useCallback(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const host = import.meta.env.VITE_WS_HOST || window.location.hostname;
    const port = import.meta.env.VITE_WS_PORT || "8080";
    const token = apiService.getAuthToken();

    console.log("[WebSocket] Building URL with:", {
      protocol,
      host,
      port,
      hasToken: !!token,
    });

    if (!token) {
      console.log(
        "[WebSocket] No token available, skipping WebSocket connection"
      );
      return null;
    }

    const wsUrl = `${protocol}//${host}:${port}/ws/chat?token=${token}`;
    console.log("[WebSocket] Connecting to:", wsUrl);
    return wsUrl;
  }, []);

  const calculateReconnectDelay = useCallback(() => {
    const exponentialDelay = Math.min(
      baseReconnectDelay * Math.pow(2, reconnectAttemptsRef.current),
      maxReconnectDelay
    );
    const jitter = Math.random() * 1000;
    return exponentialDelay + jitter;
  }, []);

  const startPresenceHeartbeat = useCallback(() => {
    if (presenceHeartbeatRef.current) {
      clearInterval(presenceHeartbeatRef.current);
    }

    if (currentRoomRef.current) {
      presenceHeartbeatRef.current = setInterval(() => {
        if (
          ws.current?.readyState === WebSocket.OPEN &&
          currentRoomRef.current
        ) {
          ws.current.send(
            JSON.stringify({
              type: "presence_heartbeat",
              room: currentRoomRef.current,
            })
          );
        }
      }, heartbeatInterval);
    }
  }, []);

  const stopPresenceHeartbeat = useCallback(() => {
    if (presenceHeartbeatRef.current) {
      clearInterval(presenceHeartbeatRef.current);
      presenceHeartbeatRef.current = null;
    }
  }, []);

  const startPingPong = useCallback(() => {
    // Clear any existing ping interval
    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current);
    }

    lastPongRef.current = Date.now();

    // Send ping every 30 seconds
    pingIntervalRef.current = setInterval(() => {
      if (ws.current?.readyState === WebSocket.OPEN) {
        // Send ping message
        try {
          ws.current.send(JSON.stringify({ type: "ping" }));
          console.log("[WebSocket] 🏓 Ping sent");
          
          // Check if we received a pong from the previous ping
          // If more than 60 seconds since last pong, connection is likely dead
          const timeSinceLastPong = Date.now() - lastPongRef.current;
          if (timeSinceLastPong > 60000) {
            console.log(
              `[WebSocket] ⚠️ Connection appears dead (no pong for ${Math.round(timeSinceLastPong / 1000)}s), forcing reconnect`
            );
            ws.current.close();
          }
        } catch (error) {
          console.error("[WebSocket] Failed to send ping:", error);
          ws.current.close();
        }
      }
    }, pingInterval);
  }, [pingInterval]);

  const stopPingPong = useCallback(() => {
    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current);
      pingIntervalRef.current = null;
    }
  }, []);

  const connect = useCallback(() => {
    const wsUrl = getWebSocketUrl();

    if (!wsUrl) {
      setConnectionState(ConnectionState.DISCONNECTED);
      return;
    }

    if (
      ws.current?.readyState === WebSocket.OPEN ||
      ws.current?.readyState === WebSocket.CONNECTING
    ) {
      console.log("WebSocket already connected or connecting");
      return;
    }

    console.log(
      "Connecting to WebSocket...",
      reconnectAttemptsRef.current > 0
        ? `(Attempt ${reconnectAttemptsRef.current + 1})`
        : ""
    );

    setConnectionState(
      reconnectAttemptsRef.current > 0
        ? ConnectionState.RECONNECTING
        : ConnectionState.CONNECTING
    );

    try {
      ws.current = new WebSocket(wsUrl);

      // Handle pong messages to update lastPongRef
      const handlePong = (event: MessageEvent) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === "pong") {
            lastPongRef.current = Date.now();
            console.log("[WebSocket] 🏓 Pong received");
          }
        } catch (error) {
          // Ignore parse errors, let ChatRoom handle them
        }
      };

      ws.current.addEventListener("message", handlePong);

      ws.current.onopen = () => {
        console.log("[WebSocket] Connected successfully");
        setConnectionState(ConnectionState.CONNECTED);
        reconnectAttemptsRef.current = 0;
        lastPongRef.current = Date.now();

        if (currentRoomRef.current) {
          console.log("[WebSocket] Rejoining room:", currentRoomRef.current);
          ws.current?.send(
            JSON.stringify({
              type: "join_room",
              room: currentRoomRef.current,
            })
          );
        }

        if (messageQueueRef.current.length > 0) {
          console.log(
            "Sending queued messages:",
            messageQueueRef.current.length
          );
          messageQueueRef.current.forEach((msg) => {
            if (ws.current?.readyState === WebSocket.OPEN) {
              ws.current.send(JSON.stringify(msg));
            }
          });
          messageQueueRef.current = [];
        }

        startPresenceHeartbeat();
        startPingPong();
      };

      ws.current.onerror = (error) => {
        console.error("[WebSocket] Error:", error);
        setConnectionState(ConnectionState.ERROR);
      };

      ws.current.onclose = (event) => {
        console.log("[WebSocket] Closed:", event.code, event.reason);
        setConnectionState(ConnectionState.DISCONNECTED);
        stopPresenceHeartbeat();
        stopPingPong();

        if (
          shouldReconnectRef.current &&
          reconnectAttemptsRef.current < maxReconnectAttempts
        ) {
          const delay = calculateReconnectDelay();
          console.log(`Reconnecting in ${Math.round(delay / 1000)}s...`);

          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttemptsRef.current += 1;
            connect();
          }, delay);
        } else if (reconnectAttemptsRef.current >= maxReconnectAttempts) {
          console.error("Max reconnection attempts reached");
          setConnectionState(ConnectionState.ERROR);
        }
      };
    } catch (error) {
      console.error("Failed to create WebSocket:", error);
      setConnectionState(ConnectionState.ERROR);
    }
  }, [
    getWebSocketUrl,
    calculateReconnectDelay,
    startPresenceHeartbeat,
    stopPresenceHeartbeat,
    startPingPong,
    stopPingPong,
  ]);

  const disconnect = useCallback(() => {
    shouldReconnectRef.current = false;
    stopPresenceHeartbeat();
    stopPingPong();

    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }

    setConnectionState(ConnectionState.DISCONNECTED);
  }, [stopPresenceHeartbeat, stopPingPong]);

  const sendMessage = useCallback((message: Partial<Message>): boolean => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));
      return true;
    } else {
      console.log("WebSocket not connected, queuing message");
      messageQueueRef.current.push(message);
      return false;
    }
  }, []);

  const setCurrentRoom = useCallback(
    (roomName: string) => {
      currentRoomRef.current = roomName;

      // Send join_room message if WebSocket is connected
      if (ws.current?.readyState === WebSocket.OPEN) {
        console.log("[WebSocket] Sending join_room for:", roomName);
        ws.current.send(
          JSON.stringify({
            type: "join_room",
            room: roomName,
          })
        );
      }

      startPresenceHeartbeat();
    },
    [startPresenceHeartbeat]
  );

  const clearCurrentRoom = useCallback(() => {
    currentRoomRef.current = null;
    stopPresenceHeartbeat();
  }, [stopPresenceHeartbeat]);

  useEffect(() => {
    const token = apiService.getAuthToken();
    if (token) {
      shouldReconnectRef.current = true;
      connect();
    }

    return () => {
      shouldReconnectRef.current = false;
      stopPresenceHeartbeat();
      stopPingPong();

      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }

      if (ws.current) {
        ws.current.close();
      }
    };
  }, [connect, stopPresenceHeartbeat, stopPingPong]);

  const value: WebSocketContextType = {
    ws,
    connectionState,
    connect,
    disconnect,
    sendMessage,
    setCurrentRoom,
    clearCurrentRoom,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketContext;
