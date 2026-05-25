import { useState, useMemo } from "react";
import { useWebsocket } from "../hooks/useWebSocket";

export const ChatPage = () => {
  const [input, setInput] = useState("");

  const wsUrl = useMemo(() => {
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
    const wsProtocol = apiUrl.startsWith("https") ? "wss" : "ws";
    const host = new URL(apiUrl).host;
    return `${wsProtocol}://${host}/api/ws?userId=2020`;
  }, []);

  const { messages, isConnected, sendMessage } = useWebsocket(wsUrl);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    // Send the raw string down the WebSocket pipe
    sendMessage(input);
    setInput("");
  };

  return (
    <div className="max-w-2xl mx-auto mt-10 p-6 bg-white border border-gray-200 rounded-xl shadow-sm font-sans">
      <h2 className="text-2xl font-bold text-gray-800 mb-2">
        Global Chat Pool
      </h2>

      {/* Connection Status Indicator */}
      <div className="text-sm text-gray-600 mb-4">
        Status:{" "}
        <span
          className={`font-bold ${isConnected ? "text-green-600" : "text-red-500"}`}
        >
          {isConnected ? "Connected" : "Disconnected"}
        </span>
      </div>

      {/* Message Feed Display */}
      <div className="h-[350px] border border-gray-200 rounded-lg p-4 bg-gray-50 overflow-y-auto flex flex-col gap-2">
        {messages.length === 0 ? (
          <p className="text-gray-400 text-center italic my-auto">
            No messages yet. Type below to broadcast...
          </p>
        ) : (
          messages.map((msg, index) => (
            <div
              key={index}
              className="flex items-center gap-2 p-2 bg-white rounded-md shadow-sm border border-gray-100"
            >
              <span className="text-base">👤</span>
              <span className="text-gray-700 break-all text-sm">{msg}</span>
            </div>
          ))
        )}
      </div>

      {/* Outbound Message Form */}
      <form onSubmit={handleSubmit} className="mt-4 flex gap-2">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message to broadcast..."
          className="flex-grow px-4 py-2 border border-gray-300 rounded-lg text-sm outline-none focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
          disabled={!isConnected}
        />
        <button
          type="submit"
          className="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium text-sm rounded-lg transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
          disabled={!isConnected}
        >
          Send
        </button>
      </form>
    </div>
  );
};
