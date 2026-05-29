import { useEffect, useRef, useCallback, useState } from "react";

export function useWebSocket(conversationId, onMessage) {
    const [connected, setConnected] = useState(false);
    const wsRef = useRef(null);
    const onMessageRef = useRef(onMessage);

    // Keep the ref current without it being a dependency
    useEffect(() => {
        onMessageRef.current = onMessage;
    }, [onMessage]);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (!conversationId || !token) return;

        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const socket = new WebSocket(
            `${protocol}//localhost:8000/api/ws?room=${conversationId}&token=${token}`,
        );
        wsRef.current = socket;

        socket.addEventListener("open", () => setConnected(true));
        socket.addEventListener("message", (event) => {
            try {
                onMessageRef.current(JSON.parse(event.data));
            } catch (err) {
                console.error("Failed to parse WS message:", err);
            }
        });
        socket.addEventListener("close", () => setConnected(false));
        socket.addEventListener("error", (err) => {
            console.error("WebSocket error:", err);
            setConnected(false);
        });

        return () => socket.close();
    }, [conversationId]); // onMessage intentionally excluded

    const send = useCallback((data) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify(data));
        }
    }, []);

    return { connected, send };
}
