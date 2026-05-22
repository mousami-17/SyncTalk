import React, { useRef, useState } from "react";
import { Paperclip } from "lucide-react";

interface FileUploadProps {
  onFileSelect: (file: File) => void;
  disabled?: boolean;
  buttonClass?: string;
}

const FileUpload: React.FC<FileUploadProps> = ({ onFileSelect, disabled, buttonClass = "btn-glass p-3 disabled:opacity-50 disabled:cursor-not-allowed" }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      onFileSelect(file);
    }
    // Reset input
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const file = e.dataTransfer.files?.[0];
    if (file) {
      onFileSelect(file);
    }
  };

  return (
    <>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*,video/*,audio/*"
        onChange={handleFileChange}
        className="hidden"
        disabled={disabled}
      />
      
      <button
        type="button"
        onClick={() => fileInputRef.current?.click()}
        disabled={disabled}
        className={buttonClass}
        title="Upload image, video, or audio"
      >
        <Paperclip className="w-5 h-5" />
      </button>

      {isDragging && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <div className="card-glass p-8 text-center shadow-soft">
            <Paperclip className="w-16 h-16 mx-auto mb-4 text-gray-400" />
            <p className="text-xl text-white font-semibold">Drop file here</p>
            <p className="text-white/60 mt-2">Images, videos, or audio files</p>
          </div>
        </div>
      )}
    </>
  );
};

export default FileUpload;
