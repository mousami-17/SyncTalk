import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Hash, ArrowRight, LogOut, MessageSquare } from "lucide-react";
import { apiService } from "../services/api";
import { validateRoomName } from "../utils/validation";
import BackgroundEffects from "../components/BackgroundEffects";

const RoomSelection: React.FC = () => {
  const navigate = useNavigate();
  const [roomName, setRoomName] = useState("");
  const [error, setError] = useState("");

  const popularRooms = [
    { name: "general", description: "General discussions" },
    { name: "random", description: "Random conversations" },
    { name: "tech", description: "Technology talks" },
    { name: "gaming", description: "Gaming community" },
    { name: "music", description: "Music discussions" }
  ];

  const handleJoinRoom = (room: string) => {
    const validation = validateRoomName(room);

    if (!validation.isValid) {
      setError(validation.error || "Invalid room name");
      return;
    }

    navigate(`/chat/${room}`);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    handleJoinRoom(roomName);
  };

  const handleLogout = () => {
    apiService.removeAuthToken();
    navigate("/");
  };

  return (
    <div className="min-h-screen flex flex-col relative overflow-hidden">
      <BackgroundEffects />
      
      {/* Header with logout */}
      <div className="p-6 flex justify-end">
        <button
          onClick={handleLogout}
          className="btn-logout flex items-center gap-2"
        >
          <LogOut className="w-4 h-4" />
          <span>Logout</span>
        </button>
      </div>

      <div className="flex-1 flex flex-col items-center justify-center px-4 relative z-10">
        <div className="w-full max-w-2xl mb-8 text-center">
          <div className="flex justify-center mb-6">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center shadow-soft">
              <MessageSquare className="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Choose a Room</h1>
          <p className="text-white/60 text-lg">
            Join or create a chat room
          </p>
        </div>

        <form onSubmit={handleSubmit} className="w-full max-w-lg mb-12">
          {error && (
            <div className="glass bg-red-500/20 border-red-500/50 rounded-xl p-3 mb-4 animate-fade-in">
              <p className="text-red-200 text-sm">{error}</p>
            </div>
          )}

          <div className="flex gap-3">
            <input
              type="text"
              value={roomName}
              onChange={(e) => {
                setRoomName(e.target.value);
                setError("");
              }}
              placeholder="Enter room name"
              className="input-glass flex-1"
            />
            <button
              type="submit"
              className="btn-glass px-6 flex items-center gap-2"
            >
              <span>Join</span>
              <ArrowRight className="w-5 h-5" />
            </button>
          </div>
        </form>

        <div className="w-full max-w-2xl">
          <h2 className="text-xl font-semibold text-white mb-6 text-center">
            Popular Rooms
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {popularRooms.map((roomObj) => (
              <button
                key={roomObj.name}
                onClick={() => handleJoinRoom(roomObj.name)}
                className="glass neuro rounded-xl p-5 flex items-center gap-3 hover:bg-white/20 transition-all duration-300 group"
              >
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-gray-700/30 to-gray-800/30 flex items-center justify-center group-hover:from-gray-600/50 group-hover:to-gray-700/50 transition-all">
                  <Hash className="w-6 h-6 text-gray-400" />
                </div>
                <div className="flex-1 text-left">
                  <p className="text-white font-semibold text-lg">#{roomObj.name}</p>
                  <p className="text-white/70">
                    {roomObj.description}
                  </p>
                </div>
                <ArrowRight className="w-6 h-6 text-white/60 group-hover:text-gray-400 transition-colors" />
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default RoomSelection;
