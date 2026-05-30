import React, { useRef, useState } from "react";
import { api } from "../api/client";

export function MessageInput({ onSend, disabled }) {
  const [text, setText] = useState("");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef(null);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!text.trim() || disabled) return;
    onSend({ type: 0, content: { text } });
    setText("");
  };

  const handleFileSelect = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    setUploading(true);
    try {
      const { fileUrl, fileName } = await api.uploadFile(file);
      let type = 3;
      if (file.type.startsWith("image/")) type = 1;
      else if (file.type.startsWith("video/")) type = 2;
      else if (file.type.startsWith("audio/")) type = 2;
      onSend({ type, content: { url: fileUrl, name: fileName } });
    } catch (err) {
      console.error("Upload failed:", err);
    } finally {
      setUploading(false);
      e.target.value = "";
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="px-6 py-4 border-t border-gray-700 bg-gray-800 flex gap-3 items-center"
    >
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileSelect}
        className="hidden"
        accept="image/*,video/*,audio/*,.pdf,.doc,.docx"
      />
      <button
        type="button"
        onClick={() => fileInputRef.current?.click()}
        disabled={uploading || disabled}
        className="px-3 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-gray-300 rounded-lg transition"
      >
        {uploading ? "⏳" : "📎"}
      </button>
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Type a message..."
        className="flex-1 px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500 placeholder-gray-500"
      />
      <button
        type="submit"
        disabled={!text.trim() || disabled}
        className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-medium rounded-lg transition"
      >
        Send
      </button>
    </form>
  );
}
