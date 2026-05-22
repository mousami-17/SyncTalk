import React, { useState, useEffect, useRef, useCallback } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Send, LogOut, Hash, Loader2 } from "lucide-react";
import { useWebSocket } from "../contexts/WebSocketContext";
import { apiService } from "../services/api";
import { ConnectionState } from "../types";
import type {
  Message,
  User,
  TypingIndicator as TypingIndicatorType,
} from "../types";
import MessageBubble from "../components/MessageBubble";
import UserList from "../components/UserList";
import TypingIndicator from "../components/TypingIndicator";
import ConnectionStatus from "../components/ConnectionStatus";
import BackgroundEffects from "../components/BackgroundEffects";
import FileUpload from "../components/FileUpload";
import FilePreview from "../components/FilePreview";
import { getDateLabel, isDifferentDay } from "../utils/dateUtils";

const ChatRoom: React.FC = () => {
  const { roomName } = useParams<{ roomName: string }>();
  const navigate = useNavigate();
  const { ws, connectionState, sendMessage, setCurrentRoom, clearCurrentRoom } =
    useWebSocket();

  const [messages, setMessages] = useState<Message[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [typingUsers, setTypingUsers] = useState<TypingIndicatorType[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [showScrollButton, setShowScrollButton] = useState(false);

  // File upload state
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isTypingRef = useRef(false);
  const currentUserId = apiService.getUserId() || "";

  const scrollToBottom = useCallback((smooth = true) => {
    if (messagesContainerRef.current) {
      // For smooth scroll to bottom
      if (smooth) {
        messagesContainerRef.current.scrollTo({
          top: messagesContainerRef.current.scrollHeight,
          behavior: "smooth",
        });
      } else {
        messagesContainerRef.current.scrollTop =
          messagesContainerRef.current.scrollHeight;
      }
    } else {
      // Fallback to element ref if container isn't available
      messagesEndRef.current?.scrollIntoView({
        behavior: smooth ? "smooth" : "auto",
      });
    }
  }, []);

  const handleScroll = useCallback(() => {
    if (!messagesContainerRef.current) return;

    const { scrollTop, scrollHeight, clientHeight } =
      messagesContainerRef.current;
    const isNearBottom = scrollHeight - scrollTop - clientHeight < 100;
    setShowScrollButton(!isNearBottom);
  }, []);

  useEffect(() => {
    if (!roomName) {
      navigate("/chat");
      return;
    }

    setCurrentRoom(roomName);

    const loadHistory = async () => {
      setIsLoadingHistory(true);
      try {
        const history = await apiService.fetchRoomMessages(roomName, 50, 0);
        const historicalMessages = history.messages.map((msg) => ({
          ...msg,
          isHistorical: true,
        }));
        setMessages(historicalMessages);
        setTimeout(() => scrollToBottom(false), 100);
      } catch (error) {
        console.error("Failed to load message history:", error);
      } finally {
        setIsLoadingHistory(false);
      }
    };

    loadHistory();

    return () => {
      clearCurrentRoom();
    };
  }, [roomName, navigate, setCurrentRoom, clearCurrentRoom, scrollToBottom]);

  useEffect(() => {
    if (!ws.current || connectionState !== ConnectionState.CONNECTED) {
      console.log("[ChatRoom] WebSocket not ready, state:", connectionState);
      return;
    }

    console.log(
      "[ChatRoom] Setting up message listener, WebSocket state:",
      ws.current.readyState
    );

    const handleMessage = (event: MessageEvent) => {
      console.log("[ChatRoom] 🎯 handleMessage called with event:", event);
      try {
        const data = JSON.parse(event.data);
        console.log("[ChatRoom] Received WebSocket message:", data);

        switch (data.type) {
          case "chat_message":
            console.log("[ChatRoom] Processing chat_message:", data);
            setMessages((prev) => {
              // Check if this exact message already exists (by ID or exact match)
              const isDuplicate = prev.some(
                (m) =>
                  (m.id && m.id === data.id) || // Match by database ID if available
                  (m.sender === data.sender &&
                    m.content === data.content &&
                    m.timestamp === data.timestamp &&
                    !m.pending &&
                    !m.tempId) // Only consider non-pending, non-temp messages as duplicates
              );

              if (isDuplicate) {
                console.log("[ChatRoom] Duplicate message detected, skipping");
                return prev;
              }

              // Remove pending message with matching tempId (for sender's optimistic update)
              // We need to find and remove the temp message that matches this real message
              // Match by: same sender, same content, and has a tempId (is pending)
              const withoutTemp = prev.filter((m) => {
                // Keep all messages that are NOT the pending version of this message
                const isPendingVersionOfThisMessage =
                  m.tempId && // Has a tempId (is pending)
                  m.pending && // Is marked as pending
                  m.sender === data.sender && // Same sender
                  m.content === data.content; // Same content

                return !isPendingVersionOfThisMessage; // Keep everything except the pending version
              });

              console.log(
                "[ChatRoom] Adding message to state. Previous:",
                prev.length,
                "After removing temp:",
                withoutTemp.length,
                "Sender:",
                data.sender,
                "Content:",
                data.content
              );
              // Preserve existing isHistorical status, default to false for new messages
              const updatedMessages = [
                ...withoutTemp,
                { ...data, isHistorical: data.isHistorical || false },
              ];
              // Check if we were already near the bottom before adding the message
              if (messagesContainerRef.current) {
                // Update state first
                return updatedMessages;
              }
              return updatedMessages;
            });
            // Scroll to bottom after message is added to maintain chat position
            setTimeout(() => {
              if (messagesContainerRef.current) {
                messagesContainerRef.current.scrollTo({
                  top: messagesContainerRef.current.scrollHeight,
                  behavior: "smooth",
                });
              } else {
                messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
              }
            }, 10);
            break;

          case "join_room":
            // Update users list when receiving join_room response
            if (data.users) {
              console.log(
                "[ChatRoom] Updating users from join_room:",
                data.users
              );
              // Map backend format (userId, userName) to frontend format (id, name)
              const mappedUsers = data.users.map((user: any) => ({
                id: user.userId || user.id,
                name: user.userName || user.name,
                status: user.status || "online",
              }));
              setUsers(mappedUsers);
            }
            break;

          case "presence_update":
            if (data.users) {
              console.log(
                "[ChatRoom] Updating users from presence_update:",
                data.users
              );
              // Map backend format (userId, userName) to frontend format (id, name)
              const mappedUsers = data.users.map((user: any) => ({
                id: user.userId || user.id,
                name: user.userName || user.name,
                status: user.status || "online",
              }));
              setUsers(mappedUsers);
            }
            break;

          case "typing_indicator":
            setTypingUsers((prev) => {
              const filtered = prev.filter((u) => u.userId !== data.userId);
              if (data.isTyping) {
                return [
                  ...filtered,
                  {
                    userId: data.userId,
                    userName: data.userName,
                    isTyping: true,
                  },
                ];
              }
              return filtered;
            });
            break;
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    ws.current.addEventListener("message", handleMessage);

    return () => {
      ws.current?.removeEventListener("message", handleMessage);
    };
  }, [ws, connectionState]); // Re-run when connection state changes (reconnects)

  const handleTyping = useCallback(() => {
    if (!isTypingRef.current && roomName) {
      isTypingRef.current = true;
      sendMessage({
        type: "typing_indicator",
        room: roomName,
        content: "",
        sender: currentUserId,
        senderName: "",
        timestamp: new Date().toISOString(),
      });
    }

    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    typingTimeoutRef.current = setTimeout(() => {
      isTypingRef.current = false;
    }, 3000);
  }, [roomName, sendMessage, currentUserId]);

  const handleFileSelect = useCallback((file: File) => {
    setSelectedFile(file);
  }, []);

  const handleCancelUpload = useCallback(() => {
    setSelectedFile(null);
    setUploadProgress(0);
  }, []);

  const handleSendFile = useCallback(async () => {
    if (!selectedFile || !roomName) return;

    setUploading(true);
    setUploadProgress(0);

    try {
      // Simulate upload progress (you can implement real progress tracking)
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return 90;
          }
          return prev + 10;
        });
      }, 200);

      // Upload file
      const response = await apiService.uploadFile(selectedFile, roomName);

      clearInterval(progressInterval);
      setUploadProgress(100);

      if (!response.error) {
        // Send message with file metadata
        const fileMessage: Message = {
          type: "chat_message",
          content: `Sent ${response.fileType}`,
          room: roomName,
          sender: currentUserId,
          senderName: "",
          timestamp: new Date().toISOString(),
          fileUrl: response.fileUrl,
          fileType: response.fileType,
          fileName: response.fileName,
          fileSize: response.fileSize,
          thumbnail: response.thumbnail,
          duration: response.duration,
          width: response.width,
          height: response.height,
        };

        sendMessage(fileMessage);
        setMessages((prev) => [...prev, { ...fileMessage, pending: false }]);
      }

      // Reset state
      setSelectedFile(null);
      setUploadProgress(0);
      // Scroll to bottom smoothly after adding the message
      setTimeout(() => {
        if (messagesContainerRef.current) {
          messagesContainerRef.current.scrollTo({
            top: messagesContainerRef.current.scrollHeight,
            behavior: "smooth",
          });
        } else {
          scrollToBottom(true); // Use smooth scroll
        }
      }, 10);
    } catch (error) {
      console.error("File upload failed:", error);
      alert("Failed to upload file. Please try again.");
    } finally {
      setUploading(false);
    }
  }, [selectedFile, roomName, currentUserId, sendMessage, scrollToBottom]);

  const handleSendMessage = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();

      if (!inputMessage.trim() || !roomName) return;

      const tempId = `temp-${Date.now()}`;
      const newMessage: Message = {
        tempId,
        type: "chat_message",
        content: inputMessage.trim(),
        room: roomName,
        sender: currentUserId,
        senderName: "",
        timestamp: new Date().toISOString(),
        pending: true,
      };

      setMessages((prev) => {
        const newMessages = [...prev, newMessage];
        return newMessages;
      });
      sendMessage(newMessage);
      setInputMessage("");
      isTypingRef.current = false;

      // Scroll to bottom smoothly after adding the message
      setTimeout(() => {
        if (messagesContainerRef.current) {
          messagesContainerRef.current.scrollTo({
            top: messagesContainerRef.current.scrollHeight,
            behavior: "smooth",
          });
        } else {
          scrollToBottom(true); // Use smooth scroll
        }
      }, 10);
    },
    [inputMessage, roomName, currentUserId, sendMessage, scrollToBottom]
  );

  const handleLogout = () => {
    apiService.removeAuthToken();
    navigate("/");
  };

  const renderMessages = () => {
    const messageElements: React.ReactElement[] = [];

    messages.forEach((message, index) => {
      const prevMessage = index > 0 ? messages[index - 1] : null;
      const showDateLabel =
        !prevMessage ||
        isDifferentDay(prevMessage.timestamp, message.timestamp);

      if (showDateLabel) {
        messageElements.push(
          <div
            key={`date-${message.timestamp}`}
            className="flex justify-center my-4"
          >
            <div className="glass shadow-soft px-4 py-1 rounded-full">
              <span className="text-xs text-white/60 font-medium">
                {getDateLabel(message.timestamp)}
              </span>
            </div>
          </div>
        );
      }

      const isOwn = message.sender === currentUserId;
      const showSender = !prevMessage || prevMessage.sender !== message.sender;

      messageElements.push(
        <MessageBubble
          key={message.id || message.tempId || index}
          message={message}
          isOwn={isOwn}
          showSender={showSender}
        />
      );
    });

    return messageElements;
  };

  return (
    <div className="h-screen flex relative overflow-hidden">
      <BackgroundEffects />

      <div className="flex-1 flex flex-col relative z-10">
        <header className="glass-strong shadow-inset border-b border-gray-700/40 p-4 flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center shadow-soft">
                <Hash className="w-5 h-5 text-white" />
              </div>
              <div>
                <h1 className="text-xl font-bold text-white">{roomName}</h1>
                <p className="text-sm text-white/60">
                  {users.length} members online
                </p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <ConnectionStatus connectionState={connectionState} />
              <button
                onClick={handleLogout}
                className="btn-logout flex items-center gap-2"
              >
                <LogOut className="w-4 h-4" />
                <span className="hidden sm:inline">Logout</span>
              </button>
            </div>
          </div>
        </header>

        <div className="flex-1 flex gap-4 p-4 overflow-hidden">
          <div className="flex-1 flex flex-col glass-strong neuro rounded-2xl overflow-hidden">
            <div
              ref={messagesContainerRef}
              onScroll={handleScroll}
              className="flex-1 overflow-y-auto p-4 space-y-2"
              style={{ minHeight: 0 }} // This allows the flex item to shrink below its content size
            >
              {isLoadingHistory ? (
                <div className="flex items-center justify-center h-full">
                  <Loader2 className="w-8 h-8 text-violet-400 animate-spin" />
                </div>
              ) : messages.length === 0 ? (
                <div className="flex items-center justify-center h-full">
                  <div className="text-center">
                    <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-700/20 to-gray-800/20 flex items-center justify-center mx-auto mb-4 shadow-soft">
                      <Hash className="w-8 h-8 text-gray-400" />
                    </div>
                    <p className="text-white/60">
                      No messages yet. Start the conversation!
                    </p>
                  </div>
                </div>
              ) : (
                renderMessages()
              )}
              <div ref={messagesEndRef} />
            </div>

            <TypingIndicator
              typingUsers={typingUsers}
              currentUserId={currentUserId}
            />

            {showScrollButton && (
              <button
                onClick={() => scrollToBottom()}
                className="absolute bottom-24 right-8 btn-glass w-10 h-10 rounded-full flex items-center justify-center animate-fade-in shadow-soft"
              >
                ↓
              </button>
            )}

            <form
              onSubmit={handleSendMessage}
              className="p-4 border-t border-white/10 flex-shrink-0"
            >
              <div className="flex gap-2">
                <FileUpload
                  onFileSelect={handleFileSelect}
                  disabled={connectionState !== ConnectionState.CONNECTED}
                  buttonClass="btn-attach p-3"
                />
                <div className="relative flex-1">
                  <input
                    type="text"
                    value={inputMessage}
                    onChange={(e) => {
                      setInputMessage(e.target.value);
                      handleTyping();
                    }}
                    placeholder="Type a message..."
                    className="input-glass flex-1 w-full"
                    disabled={connectionState !== ConnectionState.CONNECTED}
                  />
                </div>
                <button
                  type="submit"
                  disabled={
                    !inputMessage.trim() ||
                    connectionState !== ConnectionState.CONNECTED
                  }
                  className="btn-send px-6 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Send className="w-5 h-5" />
                </button>
              </div>
            </form>
          </div>

          <div className="hidden lg:block w-80 flex-shrink-0">
            <UserList
              users={users}
              currentUserId={currentUserId}
              typingUsers={typingUsers}
            />
          </div>
        </div>
      </div>

      {/* File Preview Modal */}
      {selectedFile && (
        <FilePreview
          file={selectedFile}
          onSend={handleSendFile}
          onCancel={handleCancelUpload}
          uploading={uploading}
          uploadProgress={uploadProgress}
        />
      )}
    </div>
  );
};

export default ChatRoom;
