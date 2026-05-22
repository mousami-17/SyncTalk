import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { WebSocketProvider } from "./contexts/WebSocketContext";
import { apiService } from "./services/api";
import Auth from "./pages/Auth";
import RoomSelection from "./pages/RoomSelection";
import ChatRoom from "./pages/ChatRoom";

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const isAuthenticated = apiService.isAuthenticated();
  return isAuthenticated ? <>{children}</> : <Navigate to="/" replace />;
};

const App: React.FC = () => {
  return (
    <Router>
      <WebSocketProvider>
        <Routes>
          <Route path="/" element={<Auth />} />
          <Route
            path="/chat"
            element={
              <ProtectedRoute>
                <RoomSelection />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chat/:roomName"
            element={
              <ProtectedRoute>
                <ChatRoom />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </WebSocketProvider>
    </Router>
  );
};

export default App;
