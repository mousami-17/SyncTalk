import React from "react";
import { Wifi, WifiOff, Loader2 } from "lucide-react";
import { ConnectionState } from "../types";

interface ConnectionStatusProps {
  connectionState: ConnectionState;
}

const ConnectionStatus: React.FC<ConnectionStatusProps> = ({
  connectionState,
}) => {
  const getStatusConfig = () => {
    switch (connectionState) {
      case ConnectionState.CONNECTED:
        return {
          icon: Wifi,
          text: "Online",
          className: "text-green-400",
        };
      case ConnectionState.CONNECTING:
        return {
          icon: Loader2,
          text: "Connecting...",
          className: "text-yellow-400",
          animate: "animate-spin",
        };
      case ConnectionState.RECONNECTING:
        return {
          icon: Loader2,
          text: "Reconnecting...",
          className: "text-orange-400",
          animate: "animate-spin",
        };
      case ConnectionState.DISCONNECTED:
        return {
          icon: WifiOff,
          text: "Offline",
          className: "text-gray-400",
        };
      case ConnectionState.ERROR:
        return {
          icon: WifiOff,
          text: "Error",
          className: "text-red-400",
        };
      default:
        return {
          icon: WifiOff,
          text: "Unknown",
          className: "text-gray-400",
        };
    }
  };

  const config = getStatusConfig();
  const Icon = config.icon;

  return (
    <div className={`glass shadow-soft px-4 py-2 rounded-xl flex items-center gap-2 ${connectionState === ConnectionState.CONNECTED ? 'opacity-50 hover:opacity-100 transition-opacity' : ''}`}>
      <Icon className={`w-4 h-4 ${config.className} ${config.animate || ""}`} />
      <span className={`text-sm font-medium ${config.className}`}>
        {config.text}
      </span>
    </div>
  );
};

export default ConnectionStatus;
