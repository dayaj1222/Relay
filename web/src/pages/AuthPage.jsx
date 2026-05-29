import React, { useState } from "react";
import { useAuth } from "../hooks/useAuth";

export function AuthPage() {
  const { login, register, error: authError, loading: authLoading } = useAuth();
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setError("");
      setLoading(true);
      if (!username || !password) {
        setError("Please fill all fields");
        return;
      }
      if (!isLogin) {
        if (!email) {
          setError("Please enter your email");
          return;
        }
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
          setError("Please enter a valid email");
          return;
        }
        if (password.length < 8) {
          setError("Password must be at least 8 characters");
          return;
        }
        if (!/[A-Z]/.test(password)) {
          setError("Password must contain uppercase letter");
          return;
        }
        if (!/[a-z]/.test(password)) {
          setError("Password must contain lowercase letter");
          return;
        }
        if (!/[0-9]/.test(password)) {
          setError("Password must contain number");
          return;
        }
      }
      if (isLogin) {
        await login(username, password);
      } else {
        await register(username, password, email);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="bg-gray-800 rounded-lg p-8 w-full max-w-sm border border-gray-700">
        <h1 className="text-3xl font-bold text-white mb-2">Relay Chat</h1>
        <p className="text-gray-400 mb-6">Connect & chat in real-time</p>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Username
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
              disabled={loading}
            />
          </div>
          {!isLogin && (
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Email
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                disabled={loading}
              />
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
              disabled={loading}
            />
            {!isLogin && (
              <p className="text-xs text-gray-400 mt-2">
                Min 8 chars, uppercase, lowercase, number
              </p>
            )}
          </div>
          {(error || authError) && (
            <div className="p-3 bg-red-900 border border-red-700 rounded-lg text-red-200 text-sm">
              {error || authError}
            </div>
          )}
          <button
            type="submit"
            disabled={loading || authLoading}
            className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-medium rounded-lg transition"
          >
            {loading ? "Loading..." : isLogin ? "Login" : "Register"}
          </button>
        </form>
        <p className="mt-6 text-center text-gray-400 text-sm">
          {isLogin ? "Don't have an account? " : "Already have an account? "}
          <button
            onClick={() => {
              setIsLogin(!isLogin);
              setError("");
              setUsername("");
              setPassword("");
              setEmail("");
            }}
            className="text-blue-400 hover:text-blue-300 font-medium"
          >
            {isLogin ? "Register" : "Login"}
          </button>
        </p>
      </div>
    </div>
  );
}
