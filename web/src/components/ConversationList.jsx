import React from "react";

export function ConversationList({ conversations, active, onSelect, user }) {
  if (!conversations?.length) {
    return (
      <div className="p-4 text-center text-gray-500 text-sm">
        No conversations
      </div>
    );
  }

  return (
    <div className="space-y-1">
      {conversations.map((conv) => (
        <button
          key={conv.id}
          onClick={() => onSelect(conv)}
          className={`w-full text-left px-4 py-3 rounded-lg transition ${
            active?.id === conv.id
              ? "bg-blue-600 text-white"
              : "text-gray-300 hover:bg-gray-700"
          }`}
        >
          <div className="font-medium text-sm truncate">
            {conv.type === "dm"
              ? `DM with #${conv.userId1 === user?.id ? conv.userId2 : conv.userId1}`
              : conv.name || "Group"}
          </div>
          <div className="text-xs text-gray-400 truncate mt-1">
            {new Date(conv.updatedAt).toLocaleDateString()}
          </div>
        </button>
      ))}
    </div>
  );
}
