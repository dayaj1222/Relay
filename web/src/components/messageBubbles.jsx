import React from "react";

// Explicit matching mapping enum
export const MessageType = {
  TEXT: 0,
  VIDEO: 1,
  AUDIO: 2,
  DOCUMENT: 3,
};

const TextMessage = ({ content }) => {
  // If backend fallback maps raw text string, render safely
  const text = typeof content === "string" ? content : content?.text || "";
  return <p className="break-words text-sm leading-relaxed">{text}</p>;
};

const VideoMessage = ({ content }) => (
  <div className="rounded-md overflow-hidden max-w-xs md:max-w-sm mt-1 border border-black/10">
    <video
      src={content?.fileUrl}
      controls
      className="w-full max-h-60 bg-black"
    />
  </div>
);

const AudioMessage = ({ content }) => (
  <div className="mt-1 w-full max-w-xs bg-black/5 rounded-md p-1 border">
    <audio src={content?.fileUrl} controls className="w-full h-8" />
  </div>
);

const DocumentMessage = ({ content }) => (
  <div className="flex items-center gap-3 p-2.5 bg-black/5 dark:bg-white/5 rounded-md border text-sm max-w-sm">
    <span className="text-2xl select-none">📄</span>
    <div className="flex flex-col min-w-0 flex-1">
      <span className="font-medium truncate text-gray-900 dark:text-gray-100">
        {content?.fileName || "Attachment Document"}
      </span>
      {content?.fileSize && (
        <span className="text-[11px] opacity-60">{content.fileSize}</span>
      )}
    </div>
    <a
      href={content?.fileUrl}
      target="_blank"
      rel="noreferrer"
      className="p-1 px-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded text-xs font-semibold tracking-wide transition"
    >
      OPEN
    </a>
  </div>
);

export const MessageBubbles = ({ messages, currentUserId }) => {
  return (
    <div className="flex flex-col gap-4 p-4 overflow-y-auto h-full w-full">
      {messages.map((msg, index) => {
        // Enforce fallback checks to protect against empty records
        const sender = msg.senderId || "Unknown System User";
        const isMe = String(sender) === String(currentUserId);

        return (
          <div
            key={index}
            className={`flex flex-col max-w-[75%] p-3 rounded-2xl shadow-sm tracking-wide transition-all ${
              isMe
                ? "bg-blue-600 text-white self-end rounded-tr-none"
                : "bg-gray-200 text-gray-800 self-start rounded-tl-none"
            }`}
          >
            {/* Header Identity Label */}
            <span
              className={`text-[10px] font-bold uppercase tracking-wider mb-1 opacity-60 ${
                isMe ? "text-right" : "text-left"
              }`}
            >
              {isMe ? "You" : `User: ${sender}`}
            </span>

            {/* Dynamic Content Component Multi-Router */}
            <div className="w-full">
              {(() => {
                // Ensure parsing comparison checks type cleanly
                const typeNum = Number(msg.type);
                switch (typeNum) {
                  case MessageType.TEXT:
                    return <TextMessage content={msg.content} />;
                  case MessageType.VIDEO:
                    return <VideoMessage content={msg.content} />;
                  case MessageType.AUDIO:
                    return <AudioMessage content={msg.content} />;
                  case MessageType.DOCUMENT:
                    return <DocumentMessage content={msg.content} />;
                  default:
                    return (
                      <div className="flex flex-col text-xs italic opacity-60 p-1">
                        <span>Unsupported view format</span>
                        <span className="text-[10px] font-mono opacity-50">
                          Type Code: {msg.type}
                        </span>
                      </div>
                    );
                }
              })()}
            </div>
          </div>
        );
      })}
    </div>
  );
};
