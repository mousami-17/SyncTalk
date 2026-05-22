import React, { useState, useEffect } from "react";
import { X, Send, Loader2, Music } from "lucide-react";

interface FilePreviewProps {
  file: File;
  onSend: () => void;
  onCancel: () => void;
  uploading: boolean;
  uploadProgress?: number;
}

const FilePreview: React.FC<FilePreviewProps> = ({
  file,
  onSend,
  onCancel,
  uploading,
  uploadProgress = 0,
}) => {
  const [preview, setPreview] = useState<string>("");
  const fileType = file.type.split("/")[0]; // "image", "video", "audio"

  useEffect(() => {
    if (fileType === "image" || fileType === "video") {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }

    return () => {
      if (preview) {
        URL.revokeObjectURL(preview);
      }
    };
  }, [file]);

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <div className="card-glass max-w-2xl w-full animate-fade-in">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-white">Send File</h3>
          <button
            onClick={onCancel}
            disabled={uploading}
            className="btn-logout p-2 disabled:opacity-50 !bg-red-500/30 !border-red-400/30"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="mb-4">
          {fileType === "image" && preview && (
            <img
              src={preview}
              alt={file.name}
              className="w-full h-64 object-contain rounded-xl bg-black/20"
            />
          )}

          {fileType === "video" && preview && (
            <video
              src={preview}
              controls
              className="w-full h-64 rounded-xl bg-black/20"
            />
          )}

          {fileType === "audio" && (
            <div className="flex items-center justify-center h-32 bg-gradient-to-br from-gray-700/20 to-gray-800/20 rounded-xl shadow-soft">
              <Music className="w-16 h-16 text-gray-400" />
            </div>
          )}
        </div>

        <div className="mb-4">
          <p className="text-white font-medium truncate">{file.name}</p>
          <p className="text-white/60 text-sm">
            {formatFileSize(file.size)} • {file.type}
          </p>
        </div>

        {uploading && (
          <div className="mb-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-white/60">Uploading...</span>
              <span className="text-sm text-white/60">{uploadProgress}%</span>
            </div>
            <div className="h-2 bg-gray-700/40 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-gray-500 to-gray-600 transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}

        <div className="flex gap-2">
          <button
            onClick={onCancel}
            disabled={uploading}
            className="btn-glass flex-1 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={onSend}
            disabled={uploading}
            className="btn-send flex-1 flex items-center justify-center gap-2 disabled:opacity-50"
          >
            {uploading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                <span>Uploading...</span>
              </>
            ) : (
              <>
                <Send className="w-5 h-5" />
                <span>Send</span>
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default FilePreview;
