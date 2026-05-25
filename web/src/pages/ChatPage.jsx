import React, { useState } from "react";
import { useWebsocket } from "../hooks/useWebsocket";
import { MessageBubbles, MessageType } from "../components/messageBubbles.jsx";

export const ChatPage = () => {
  // Replace with your current actual WebSocket URL config
  const { messages, sendMessage } = useWebsocket(
    "ws://[::1]:8080/api/ws?userId=2020",
  );
  const [textInput, setTextInput] = useState("");
  const currentUserId = "2020";

  // 1. Text Submission Handler
  const handleSendText = (e) => {
    e.preventDefault();
    if (!textInput.trim()) return;

    const payload = {
      senderId: "", // Backend overrides this
      type: MessageType.TEXT,
      content: textInput.trim(),
    };

    sendMessage(JSON.stringify(payload));
    setTextInput("");
  };

  // 2. File Upload Orchestrator (Bridges HTTP and WebSockets)
  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Determine target enum type by file extension/mime type
    let messageType = MessageType.DOCUMENT;
    if (file.type.startsWith("video/")) messageType = MessageType.VIDEO;
    if (file.type.startsWith("audio/")) messageType = MessageType.AUDIO;

    // Package bytes into standard multi-part form payload
    const formData = new FormData();
    formData.append("file", file);

    try {
      // Send raw file bytes over HTTP rather than WebSocket
      const response = await fetch("http://[::1]:8080/api/upload", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) throw new Error("Upload failed");
      const data = await response.json();
      // Expected backend response: { fileUrl: "...", fileName: "..." }

      // Send the lightweight signaling payload over WebSocket
      const payload = {
        senderId: "",
        type: messageType,
        content: {
          fileUrl: data.fileUrl,
          fileName: file.name,
        },
      };

      sendMessage(JSON.stringify(payload));
    } catch (err) {
      console.error("File upload network error:", err);
    }
  };

  return (
    <div className="flex flex-col h-screen max-w-2xl mx-auto border shadow-md">
      {/* Scrollable Feed Container */}
      <div className="flex-1 overflow-hidden bg-gray-50">
        <MessageBubbles messages={messages} currentUserId={currentUserId} />
      </div>

      {/* Input Action Form Block */}
      <form
        onSubmit={handleSendText}
        className="flex items-center gap-2 p-3 border-t bg-white"
      >
        <label className="p-2 bg-gray-100 hover:bg-gray-200 rounded cursor-pointer text-xl">
          📎
          <input
            type="file"
            onChange={handleFileUpload}
            className="hidden"
            accept="video/*,audio/*,application/pdf,.doc,.docx"
          />
        </label>

        <input
          type="text"
          value={textInput}
          onChange={(e) => setTextInput(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 p-2 border rounded-md outline-none text-sm"
        />

        <button
          type="submit"
          className="p-2 px-4 bg-blue-600 text-white rounded-md text-sm font-medium"
        >
          Send
        </button>
      </form>
    </div>
  );
};
