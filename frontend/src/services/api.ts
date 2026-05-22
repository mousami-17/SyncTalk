import type {
  AuthRequest,
  AuthResponse,
  MessageHistoryResponse,
  ApiResponse,
} from "../types";

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

class ApiService {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem("go-chat-token");
    return {
      "Content-Type": "application/json",
      ...(token && { Authorization: token }),
    };
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    const config: RequestInit = {
      headers: this.getAuthHeaders(),
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: true,
          message: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(errorData.message || "Request failed");
      }

      return await response.json();
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  async signup(credentials: AuthRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/auth/signup", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
  }

  async login(credentials: AuthRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
  }

  async fetchRoomMessages(
    roomId: string,
    limit: number = 50,
    offset: number = 0
  ): Promise<MessageHistoryResponse> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });

    return this.request<MessageHistoryResponse>(
      `/api/rooms/${roomId}/messages?${params}`
    );
  }

  async healthCheck(): Promise<ApiResponse> {
    return this.request<ApiResponse>("/api/health");
  }

  async databaseHealth(): Promise<ApiResponse> {
    return this.request<ApiResponse>("/api/health/database");
  }

  async databaseMetrics(): Promise<ApiResponse> {
    return this.request<ApiResponse>("/api/metrics/database");
  }

  async uploadFile(file: File, roomId: string): Promise<any> {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("roomId", roomId);

    const token = localStorage.getItem("go-chat-token");
    const url = `${API_BASE_URL}/api/upload`;

    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          ...(token && { Authorization: token }),
        },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: true,
          message: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(errorData.message || "Upload failed");
      }

      return await response.json();
    } catch (error) {
      console.error("File upload failed:", error);
      throw error;
    }
  }

  setAuthToken(token: string): void {
    localStorage.setItem("go-chat-token", token);
  }

  getAuthToken(): string | null {
    return localStorage.getItem("go-chat-token");
  }

  removeAuthToken(): void {
    localStorage.removeItem("go-chat-token");
    localStorage.removeItem("go-chat-userId");
  }

  getUserId(): string | null {
    return localStorage.getItem("go-chat-userId");
  }

  setUserId(userId: string): void {
    localStorage.setItem("go-chat-userId", userId);
  }

  isAuthenticated(): boolean {
    return !!this.getAuthToken();
  }

  parseJWT(token: string): any {
    try {
      console.log("[API] Parsing JWT token:", token);
      const base64Url = token.split(".")[1];
      console.log("[API] Base64 URL part:", base64Url);
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
          .join("")
      );
      console.log("[API] Decoded JSON:", jsonPayload);
      const parsed = JSON.parse(jsonPayload);
      console.log("[API] Parsed payload:", parsed);
      return parsed;
    } catch (error) {
      console.error("[API] Failed to parse JWT:", error);
      return null;
    }
  }

  getUserIdFromToken(token: string): string | null {
    const payload = this.parseJWT(token);
    console.log("[API] JWT Payload:", payload);
    console.log("[API] Checking userID:", payload?.userID);
    console.log("[API] Checking user_id:", payload?.user_id);
    return payload?.userID || payload?.user_id || null;
  }
}

export const apiService = new ApiService();
export default apiService;
