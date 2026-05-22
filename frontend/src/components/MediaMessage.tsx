import React, { useState } from "react";
import { Volume2, Maximize2 } from "lucide-react";
import type { Message } from "../types";

interface MediaMessageProps {
  message: Message;
  isOwn: boolean;
}

const MediaMessage: React.FC<MediaMessageProps> = ({ message }) => {
  const [showFullscreen, setShowFullscreen] = useState(false);

  const formatDuration = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  const handleDownload = () => {
    if (message.fileUrl) {
      window.open(message.fileUrl, "_blank");
    }
  };

  // Image
  if (message.fileType === "image") {
    return (
      <div className="max-w-sm">
        <div
          className="relative group cursor-pointer rounded-xl overflow-hidden shadow-soft"
          onClick={() => setShowFullscreen(true)}
        >
          <img
            src={message.fileUrl}
            alt={message.fileName || "Image"}
            className="w-full h-auto rounded-xl"
            loading="lazy"
          />
          <div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-all flex items-center justify-center opacity-0 group-hover:opacity-100">
            <Maximize2 className="w-8 h-8 text-white" />
          </div>
        </div>
        {message.fileName && (
          <p className="text-xs text-white/60 mt-2 truncate">
            {message.fileName}
          </p>
        )}

        {/* Fullscreen Modal */}
        {showFullscreen && (
          <div
            className="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4 backdrop-blur-sm"
            onClick={() => setShowFullscreen(false)}
          >
            <div className="relative w-full max-w-6xl max-h-[90vh]">
              <img
                src={message.fileUrl}
                alt={message.fileName || "Image"}
                className="max-w-full max-h-full object-contain rounded-lg"
              />
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  setShowFullscreen(false);
                }}
                className="absolute top-4 right-4 btn-glass p-2 rounded-full w-10 h-10 flex items-center justify-center"
              >
                ✕
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleDownload();
                }}
                className="absolute bottom-4 right-4 btn-glass p-2 rounded-full w-10 h-10 flex items-center justify-center"
              >
                ⭳
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }

  // Video
  if (message.fileType === "video") {
    return (
      <div className="max-w-md">
        <div className="relative rounded-xl overflow-hidden bg-black shadow-soft">
          <video
            src={message.fileUrl}
            poster={message.thumbnail}
            controls
            className="w-full h-auto rounded-xl"
            preload="metadata"
          >
            Your browser does not support video playback.
          </video>
        </div>
        <div className="flex items-center justify-between mt-2">
          {message.fileName && (
            <p className="text-xs text-white/60 truncate flex-1">
              {message.fileName}
            </p>
          )}
          {message.duration && (
            <span className="text-xs text-white/60 ml-2">
              {formatDuration(message.duration)}
            </span>
          )}
        </div>
      </div>
    );
  }

  // Audio
  if (message.fileType === "audio") {
    return (
      <div className="w-80">
        <div className="glass shadow-inset rounded-xl p-4 border border-gray-700/40">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-12 h-12 rounded-full bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center shadow-soft">
              <Volume2 className="w-6 h-6 text-white" />
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-white font-medium truncate text-sm">
                {message.fileName || "Audio Message"}
              </p>
              {message.duration && (
                <p className="text-white/60 text-xs">
                  {formatDuration(message.duration)}
                </p>
              )}
            </div>
          </div>
          <audio
            src={message.fileUrl}
            controls
            className="w-full rounded-lg"
            preload="metadata"
          >
            Your browser does not support audio playback.
          </audio>
        </div>
      </div>
    );
  }

  return null;
};

export default MediaMessage;
