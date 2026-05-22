import React from "react";
import { Clock, Check, CheckCheck } from "lucide-react";
import type { Message } from "../types";
import { formatMessageTime } from "../utils/dateUtils";
import MediaMessage from "./MediaMessage";

interface MessageBubbleProps {
  message: Message;
  isOwn: boolean;
  showSender?: boolean;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({
  message,
  isOwn,
  showSender = true,
}) => {
  // Check if message has media
  const hasMedia = message.fileUrl && message.fileType;

  return (
    <div
      className={`flex ${
        isOwn ? "justify-end" : "justify-start"
      } mb-4 animate-fade-in`}
    >
      <div
        className={`max-w-xs lg:max-w-md xl:max-w-lg ${
          isOwn ? "ml-auto" : "mr-auto"
        }`}
      >
        {showSender && !isOwn && (
          <div className="text-xs text-gray-400 mb-1 ml-3 font-medium">
            {message.senderName}
          </div>
        )}

        <div
          className={`relative overflow-hidden rounded-2xl p-4 ${
            isOwn 
              ? "bg-gradient-to-br from-gray-700/40 to-gray-800/40 backdrop-blur-xl border border-gray-600/30" 
              : "bg-gradient-to-br from-gray-800/30 to-gray-900/30 backdrop-blur-xl border border-gray-700/40"
          } ${message.pending ? "opacity-70" : ""} transition-all duration-300 hover:scale-[1.01]`}
        >
          {/* Media content */}
          {hasMedia && (
            <div className="mb-2 rounded-xl overflow-hidden">
              <MediaMessage message={message} isOwn={isOwn} />
            </div>
          )}

          {/* Text content */}
          {message.content && (
            <p className={`break-words whitespace-pre-wrap ${
              isOwn ? "text-white/95" : "text-white/90"
            }`}>
              {message.content}
            </p>
          )}

          <div className={`flex items-center gap-1 mt-2 ${
            isOwn ? "justify-between" : "justify-end"
          }`}>
            <span className="text-xs text-white/60">
              {formatMessageTime(message.timestamp)}
            </span>

            {isOwn && (
              <span className="text-white/70">
                {message.pending ? (
                  <Clock className="w-3.5 h-3.5 animate-pulse" />
                ) : message.isHistorical ? (
                  <Check className="w-3.5 h-3.5" />
                ) : (
                  <CheckCheck className="w-3.5 h-3.5 text-gray-300" />
                )}
              </span>
            )}
          </div>
          
          {/* Decorative corner elements for visual enhancement */}
          <div className={`absolute top-0 ${isOwn ? 'right-0' : 'left-0'} opacity-10`}>
            <div className="w-8 h-8 bg-gradient-to-br from-gray-600/30 to-gray-700/30 rounded-full -m-2"></div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;