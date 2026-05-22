import React, { useState } from "react";
import { MessageSquare, Eye, EyeOff, Loader2 } from "lucide-react";
import { apiService } from "../services/api";
import { validateUsername, validatePassword } from "../utils/validation";
import PasswordStrength from "../components/PasswordStrength";
import BackgroundEffects from "../components/BackgroundEffects";

const Auth: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const usernameValidation = validateUsername(username);
    if (!usernameValidation.isValid) {
      setError(usernameValidation.error || "Invalid username");
      return;
    }

    const passwordValidation = validatePassword(password);
    if (!passwordValidation.isValid) {
      setError(passwordValidation.error || "Invalid password");
      return;
    }

    setLoading(true);

    try {
      const response = isLogin
        ? await apiService.login({ username, password })
        : await apiService.signup({ username, password });

      if (response.error) {
        setError(response.message);
      } else if (response.token) {
        console.log("[Auth] Token received:", response.token);
        apiService.setAuthToken(response.token);
        const userId = apiService.getUserIdFromToken(response.token);
        console.log("[Auth] Extracted userId:", userId);
        if (userId) {
          apiService.setUserId(userId);
          console.log("[Auth] UserId saved to localStorage");
        } else {
          console.error("[Auth] Failed to extract userId from token");
        }
        // Force page reload to initialize WebSocket with new token
        window.location.href = "/chat";
      }
    } catch (err: any) {
      setError(err.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
      <BackgroundEffects />

      <div className="w-full max-w-md relative z-10">
        <div className="card-glass glow animate-fade-in">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-700 to-gray-800 mb-4 shadow-soft">
              <MessageSquare className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-3xl font-bold text-white mb-2">
              Welcome to GoChat
            </h1>
            <p className="text-white/60">
              {isLogin ? "Sign in to continue" : "Create your account"}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="glass bg-red-500/20 border-red-500/50 rounded-xl p-3 animate-fade-in">
                <p className="text-red-200 text-sm">{error}</p>
              </div>
            )}

            <div>
              <label className="block text-white/80 text-sm font-medium mb-2">
                Username
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="input-glass w-full"
                placeholder="Enter your username"
                required
                disabled={loading}
              />
            </div>

            <div>
              <label className="block text-white/80 text-sm font-medium mb-2">
                Password
              </label>
              <div className="relative">
                <input
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="input-glass w-full pr-12"
                  placeholder="Enter your password"
                  required
                  disabled={loading}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 top-1/2 -translate-y-1/2 text-white/50 hover:text-white/80 transition-colors"
                >
                  {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
                </button>
              </div>
              {!isLogin && <PasswordStrength password={password} />}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn-glass w-full flex items-center justify-center gap-2"
            >
              {loading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span>
                    {isLogin ? "Signing in..." : "Creating account..."}
                  </span>
                </>
              ) : (
                <span>{isLogin ? "Sign In" : "Sign Up"}</span>
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <div className="text-white/60 text-sm">
              {isLogin
                ? "Don't have an account? "
                : "Already have an account? "}
              <button
                onClick={() => {
                  setIsLogin(!isLogin);
                  setError("");
                }}
                className="text-gray-400 hover:text-white font-medium transition-colors"
                disabled={loading}
              >
                {isLogin ? "Sign Up" : "Sign In"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Auth;
