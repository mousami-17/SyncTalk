export interface User {
  id: string;
  name: string;
  status: "online" | "offline" | "typing";
  lastSeen?: string;
}

export interface Message {
  id?: string;
  tempId?: string;
  sender: string;
  senderName: string;
  room: string;
  type:
    | "chat_message"
    | "join_room"
    | "leave_room"
    | "presence_update"
    | "typing_indicator";
  content: string;
  timestamp: string;
  server?: string;
  pending?: boolean;
  isHistorical?: boolean;
  
  // Media fields (optional)
  fileUrl?: string;
  fileType?: "image" | "video" | "audio";
  fileName?: string;
  fileSize?: number;
  thumbnail?: string;
  duration?: number;  // For audio/video in seconds
  width?: number;     // For images/videos
  height?: number;    // For images/videos
}

export interface Room {
  id: string;
  name: string;
  userCount: number;
  users: User[];
}

export interface AuthRequest {
  username: string;
  password: string;
}

export interface AuthResponse {
  error: boolean;
  message: string;
  token?: string;
}

export const ConnectionState = {
  CONNECTING: "connecting",
  CONNECTED: "connected",
  DISCONNECTED: "disconnected",
  RECONNECTING: "reconnecting",
  ERROR: "error",
} as const;

export type ConnectionState =
  (typeof ConnectionState)[keyof typeof ConnectionState];

export interface WebSocketContextType {
  ws: React.MutableRefObject<WebSocket | null>;
  connectionState: ConnectionState;
  connect: () => void;
  disconnect: () => void;
  sendMessage: (message: Partial<Message>) => boolean;
  setCurrentRoom: (roomName: string) => void;
  clearCurrentRoom: () => void;
}

export interface MessageHistoryResponse {
  messages: Message[];
  totalCount: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

export interface PresenceUpdate {
  userId: string;
  userName: string;
  status: "online" | "offline";
  timestamp: number;
}

export interface TypingIndicator {
  userId: string;
  userName: string;
  isTyping: boolean;
}

export interface ApiResponse<T = any> {
  error: boolean;
  message: string;
  data?: T;
}

export interface PasswordStrength {
  score: number;
  label: string;
  color: string;
  requirements: {
    length: boolean;
    uppercase: boolean;
    lowercase: boolean;
    number: boolean;
    special?: boolean;
  };
}
