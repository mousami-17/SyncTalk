import React from "react";
import { User as UserIcon, Users } from "lucide-react";
import type { User, TypingIndicator } from "../types";

interface UserListProps {
  users: User[];
  currentUserId: string;
  typingUsers: TypingIndicator[];
}

const UserList: React.FC<UserListProps> = ({
  users,
  currentUserId,
  typingUsers,
}) => {
  const isUserTyping = (userId: string) => {
    return typingUsers.some((tu) => tu.userId === userId && tu.isTyping);
  };

  return (
    <div className="glass-strong shadow-soft rounded-2xl p-4 border border-gray-700/40">
      <div className="flex items-center gap-2 mb-4">
        <Users className="w-5 h-5 text-gray-400" />
        <h3 className="text-lg font-semibold text-white">
          Online Users ({users.length})
        </h3>
      </div>

      <div className="space-y-2 max-h-[calc(100vh-300px)] overflow-y-auto">
        {users.map((user) => (
          <div
            key={user.id}
            className="glass shadow-inset rounded-xl p-3 flex items-center gap-3 hover:bg-gray-700/30 transition-all duration-300"
          >
            <div className="relative">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-gray-600 to-gray-700 flex items-center justify-center">
                <UserIcon className="w-5 h-5 text-white" />
              </div>
              <div className="absolute bottom-0 right-0 w-3 h-3 rounded-full border border-gray-800"
                style={{
                  backgroundColor: 
                    user.status === "online" 
                      ? "#10b981" // green-500 
                      : user.status === "typing"
                      ? "#f59e0b" // amber-500
                      : "#6b7280" // gray-500
                }}
              />
            </div>

            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <p className="text-white font-medium truncate">
                  {user.name}
                  {user.id === currentUserId && (
                    <span className="text-gray-400 text-sm ml-1">(You)</span>
                  )}
                </p>
              </div>
              {isUserTyping(user.id) && (
                <p className="text-xs text-gray-400 animate-pulse">
                  typing...
                </p>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default UserList;
