import React from "react";
import type { TypingIndicator as TypingIndicatorType } from "../types";

interface TypingIndicatorProps {
  typingUsers: TypingIndicatorType[];
  currentUserId: string;
}

const TypingIndicator: React.FC<TypingIndicatorProps> = ({
  typingUsers,
  currentUserId,
}) => {
  const otherTypingUsers = typingUsers.filter(
    (user) => user.isTyping && user.userId !== currentUserId
  );

  if (otherTypingUsers.length === 0) {
    return null;
  }

  const getTypingText = () => {
    if (otherTypingUsers.length === 1) {
      return `${otherTypingUsers[0].userName} is typing`;
    } else if (otherTypingUsers.length === 2) {
      return `${otherTypingUsers[0].userName} and ${otherTypingUsers[1].userName} are typing`;
    } else {
      return `${otherTypingUsers.length} people are typing`;
    }
  };

  return (
    <div className="flex items-center gap-2 px-4 py-2 animate-fade-in">
      <div className="flex gap-1">
        <span className="w-2 h-2 bg-gray-400 rounded-full animate-typing-dots" />
        <span
          className="w-2 h-2 bg-gray-400 rounded-full animate-typing-dots"
          style={{ animationDelay: "0.2s" }}
        />
        <span
          className="w-2 h-2 bg-gray-400 rounded-full animate-typing-dots"
          style={{ animationDelay: "0.4s" }}
        />
      </div>
      <span className="text-sm text-white/60">{getTypingText()}</span>
    </div>
  );
};

export default TypingIndicator;
