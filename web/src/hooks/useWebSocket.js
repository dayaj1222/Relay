import { useCallback, useEffect, useRef, useState } from "react";

export const useWebsocket = (url) => {
    const [messages, setMessages] = useState([]);
    const [isConnected, setIsConnected] = useState(false);
    const ws = useRef(null);

    const connect = useCallback(() => {
        if (ws.current?.readyState === WebSocket.OPEN) return;

        const socket = new WebSocket(url);

        ws.current = socket;

        socket.onopen = () => {
            setIsConnected(true);
            console.log("Connected to Websocket server");
        };

        socket.onmessage = (event) => {
            setMessages((prev) => [...prev, event.data]);
        };

        socket.onclose = () => {
            setIsConnected(false);
            console.log("Disconnected from Websocket Server");

            setTimeout(connect, 3000);
        };

        socket.onerror = (error) => {
            console.error("Websocket error:", error);
            socket.close();
        };
    }, [url]);

    useEffect(() => {
        connect();
        return () => {
            ws.current?.close();
        };
    }, [connect]);

    const sendMessage = useCallback((message) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(message);
        } else {
            console.error("Cannot send message: Websocket is not connected");
        }
    }, []);

    return { messages, isConnected, sendMessage };
};
