import React, { useEffect, useState, useRef } from "react";
import { api } from "../api/client";
import { useWebSocket } from "../hooks/useWebSocket";
import { MessageList } from "../components/Message";
import { ConversationList } from "../components/ConversationList";
import { CreateConversationModal } from "../components/CreateConversationModal";

export function ChatPage({ user, onLogout }) {
  const [conversations, setConversations] = useState([]);
  const [activeConversation, setActiveConversation] = useState(null);
  const [messages, setMessages] = useState([]);
  const [messageText, setMessageText] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [tab, setTab] = useState("dm"); // dm or group
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createType, setCreateType] = useState("dm");
  const messagesEndRef = useRef(null);

  const { connected, send } = useWebSocket(activeConversation?.id, (data) => {
    if (data.type === "message") {
      setMessages((prev) => [...prev, data]);
      scrollToBottom();
    }
  });

  useEffect(() => {
    loadConversations();
  }, [tab]);

  useEffect(() => {
    if (activeConversation) {
      loadMessages();
    }
  }, [activeConversation]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const loadConversations = async () => {
    try {
      setError("");
      const data = await api.getConversations();
      const filtered = data.filter((c) =>
        tab === "dm" ? c.type === "dm" : c.type === "group",
      );
      // Deduplicate by id
      const seen = new Set();
      const unique = filtered.filter((c) => {
        if (seen.has(c.id)) return false;
        seen.add(c.id);
        return true;
      });
      setConversations(unique);
    } catch (err) {
      setError("Failed to load conversations");
    }
  };

  const loadMessages = async () => {
    if (!activeConversation) return;
    try {
      setLoading(true);
      const data = await api.getRecentMessages(activeConversation.id, 50);
      setMessages(data || []);
    } catch (err) {
      setError("Failed to load messages");
    } finally {
      setLoading(false);
    }
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!messageText.trim() || !activeConversation) return;

    const text = messageText;
    setMessageText("");

    try {
      const msg = await api.sendMessage(activeConversation.id, text);
      setMessages((prev) => [...prev, msg]); // ← add sent message immediately
    } catch (err) {
      setError("Failed to send message");
      setMessageText(text);
    }
  };

  const handleConversationCreated = (newConversation) => {
    setConversations((prev) => [newConversation, ...prev]);
    setActiveConversation(newConversation);
  };

  const openCreateModal = (type) => {
    setCreateType(type);
    setShowCreateModal(true);
  };

  const getConversationName = (conv) => {
    if (conv.type === "dm") {
      return `DM #${conv.userId1 === user.id ? conv.userId2 : conv.userId1}`;
    }
    return conv.name || "Group";
  };

  return (
    <div className="flex h-screen bg-gray-900">
      {/* Sidebar */}
      <div className="w-80 bg-gray-800 border-r border-gray-700 flex flex-col">
        {/* Header */}
        <div className="p-4 border-b border-gray-700">
          <h1 className="text-2xl font-bold text-white">Relay Chat</h1>
          <p className="text-xs text-gray-400 mt-1">User #{user.id}</p>
        </div>

        {/* Tabs */}
        <div className="px-4 py-3 border-b border-gray-700 flex gap-2">
          <button
            onClick={() => setTab("dm")}
            className={`flex-1 px-3 py-2 rounded-lg font-medium text-sm transition ${
              tab === "dm"
                ? "bg-blue-600 text-white"
                : "bg-gray-700 text-gray-300 hover:bg-gray-600"
            }`}
          >
            DMs
          </button>
          <button
            onClick={() => setTab("group")}
            className={`flex-1 px-3 py-2 rounded-lg font-medium text-sm transition ${
              tab === "group"
                ? "bg-blue-600 text-white"
                : "bg-gray-700 text-gray-300 hover:bg-gray-600"
            }`}
          >
            Groups
          </button>
        </div>

        {/* Create Button */}
        <div className="px-4 py-3 border-b border-gray-700">
          <button
            onClick={() => openCreateModal(tab)}
            className="w-full px-4 py-2 bg-green-600 hover:bg-green-700 text-white font-medium text-sm rounded-lg transition"
          >
            + New {tab === "dm" ? "DM" : "Group"}
          </button>
        </div>

        {/* Conversations List */}
        <div className="flex-1 overflow-y-auto">
          <ConversationList
            conversations={conversations}
            active={activeConversation}
            onSelect={setActiveConversation}
          />
        </div>

        {/* Logout */}
        <div className="p-4 border-t border-gray-700">
          <button
            onClick={onLogout}
            className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium text-sm rounded-lg transition"
          >
            Logout
          </button>
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col bg-gray-900">
        {activeConversation ? (
          <>
            {/* Chat Header */}
            <div className="px-6 py-4 border-b border-gray-700 flex justify-between items-center bg-gray-800">
              <div>
                <h2 className="text-xl font-bold text-white">
                  {getConversationName(activeConversation)}
                </h2>
                <p className="text-xs text-gray-400 mt-1">
                  {connected ? "● Connected" : "◯ Connecting..."}
                </p>
              </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto px-6 py-4">
              {loading ? (
                <div className="flex items-center justify-center h-full text-gray-500">
                  <p>Loading messages...</p>
                </div>
              ) : (
                <>
                  <MessageList messages={messages} userId={user.id} />
                  <div ref={messagesEndRef} />
                </>
              )}
            </div>

            {/* Input */}
            <form
              onSubmit={handleSendMessage}
              className="px-6 py-4 border-t border-gray-700 bg-gray-800 flex gap-3"
            >
              <input
                type="text"
                value={messageText}
                onChange={(e) => setMessageText(e.target.value)}
                placeholder="Type a message..."
                className="flex-1 px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500 placeholder-gray-500"
              />
              <button
                type="submit"
                disabled={!messageText.trim() || !connected}
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-medium rounded-lg transition"
              >
                Send
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            <p className="text-lg">Select a conversation to start chatting</p>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="fixed bottom-4 right-4 px-4 py-3 bg-red-900 border border-red-700 rounded-lg text-red-200 text-sm">
          {error}
          <button
            onClick={() => setError("")}
            className="ml-4 text-red-400 hover:text-red-300"
          >
            ✕
          </button>
        </div>
      )}

      {/* Create Conversation Modal */}
      <CreateConversationModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onCreated={handleConversationCreated}
        type={createType}
      />
    </div>
  );
}
