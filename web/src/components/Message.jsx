import React from "react";
import { Check, Clock, FileText, Play } from "lucide-react";

function parseContent(content) {
  if (typeof content === "object" && content !== null) return content;
  if (typeof content === "string") {
    try {
      return JSON.parse(atob(content));
    } catch {}
    try {
      return JSON.parse(content);
    } catch {}
  }
  return {};
}

function TextContent({ content }) {
  return <p className="text-white text-sm break-words">{content.text || ""}</p>;
}

function ImageContent({ content }) {
  const handleClick = () => window.open(content.url, "_blank");
  return (
    <div>
      <img
        src={content.url}
        alt={content.name || "image"}
        className="max-w-xs rounded-lg cursor-pointer"
        onClick={handleClick}
      />
      {content.caption && (
        <p className="text-xs text-gray-300 mt-1">{content.caption}</p>
      )}
    </div>
  );
}

function MediaContent({ content }) {
  const url = content.url || "";
  const isVideo = url.match(/\.(mp4|webm)$/i);
  if (isVideo) {
    return <video controls src={url} className="max-w-xs rounded-lg" />;
  }
  return (
    <div className="flex items-center gap-2">
      <Play size={16} className="text-white" />
      <audio controls src={url} />
    </div>
  );
}

function DocumentContent({ content }) {
  const url = content.url || "";
  const name = content.name || "Document";
  return (
    <a
      href={url}
      target="_blank"
      rel="noreferrer"
      className="flex items-center gap-2 text-blue-300 underline text-sm"
    >
      <FileText size={16} />
      {name}
    </a>
  );
}

function MessageContent({ type, content }) {
  console.log("MessageContent type:", type, "content:", content);
  if (type === 1) return <ImageContent content={content} />;
  if (type === 2) return <MediaContent content={content} />;
  if (type === 3) return <DocumentContent content={content} />;
  return <TextContent content={content} />;
}

function MessageStatus({ isTemp }) {
  if (isTemp) return <Clock size={12} className="text-gray-400" />;
  return <Check size={12} className="text-blue-100" />;
}

export function Message({ message, isOwn }) {
  const time = new Date(message.createdAt).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
  const content = parseContent(message.content);
  const isTemp = String(message.id).startsWith("temp-");

  return (
    <div className={`flex ${isOwn ? "justify-end" : "justify-start"} mb-3`}>
      <div
        className={`max-w-xs lg:max-w-md ${isOwn ? "bg-blue-600" : "bg-gray-700"} rounded-2xl ${isOwn ? "rounded-tr-sm" : "rounded-tl-sm"} px-4 py-2`}
      >
        <MessageContent type={message.type ?? 0} content={content} />
        <div className="flex items-center justify-end gap-1 mt-1">
          <span
            className={`text-xs ${isOwn ? "text-blue-100" : "text-gray-400"}`}
          >
            {time}
          </span>
          {isOwn && <MessageStatus isTemp={isTemp} />}
        </div>
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
        <Message key={msg.id} message={msg} isOwn={msg.senderId == userId} />
      ))}
    </div>
  );
}
