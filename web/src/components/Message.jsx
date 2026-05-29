import React from "react";

export function Message({ message, isOwn }) {
  const time = new Date(message.createdAt).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

  let content = "";
  try {
    const decoded = atob(message.content);
    const parsed = JSON.parse(decoded);
    content = parsed.text || "";
  } catch {
    content = message.content || "";
  }

  return (
    <div
      className={`flex ${isOwn ? "justify-end" : "justify-start"} mb-3 animate-in`}
    >
      <div
        className={`max-w-xs ${isOwn ? "bg-blue-600" : "bg-gray-700"} rounded-lg px-4 py-2`}
      >
        <p className="text-white text-sm break-words">{content}</p>
        <p
          className={`text-xs mt-1 ${isOwn ? "text-blue-100" : "text-gray-400"}`}
        >
          {time}
        </p>
      </div>
    </div>
  );
}

export function MessageList({ messages, userId }) {
  if (!messages?.length) {
    return (
      <div className="flex items-center justify-center h-full text-gray-500">
        <p>No messages yet</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-1">
      {messages.map((msg) => (
        <Message key={msg.id} message={msg} isOwn={msg.senderId === userId} />
      ))}
    </div>
  );
}
