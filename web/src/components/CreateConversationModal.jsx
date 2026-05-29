import React, { useState } from "react";
import { api } from "../api/client";

export function CreateConversationModal({ isOpen, onClose, onCreated, type }) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [input, setInput] = useState("");

  const handleCreate = async () => {
    try {
      setError(null);
      setLoading(true);
      let result;
      if (type === "dm") {
        const user = await api.getUserByUsername(input);
        result = await api.createDM(user.id);
      } else {
        result = await api.createGroup(input);
      }
      setInput("");
      onCreated(result);
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg p-6 w-96 border border-gray-700">
        <h2 className="text-xl font-bold mb-4">
          Create {type === "dm" ? "DM" : "Group"}
        </h2>

        <input
          type={type === "text"}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={type === "dm" ? "Username" : "Group name"}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white mb-4 focus:outline-none focus:border-blue-500"
        />

        {error && <p className="text-red-400 text-sm mb-4">{error}</p>}

        <div className="flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition"
          >
            Cancel
          </button>
          <button
            onClick={handleCreate}
            disabled={loading || !input}
            className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white rounded-lg transition"
          >
            {loading ? "Creating..." : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
