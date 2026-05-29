const API_BASE_URL = "http://localhost:8000/api";

class APIClient {
    setToken(token) {
        // Compatibility - use localStorage only now
        if (token) {
            localStorage.setItem("token", token);
        } else {
            localStorage.removeItem("token");
        }
    }

    getHeaders() {
        // Always live-read token for latest value
        const token = localStorage.getItem("token");
        return {
            "Content-Type": "application/json",
            ...(token && { Authorization: `Bearer ${token}` }),
        };
    }

    async request(endpoint, options = {}) {
        const url = `${API_BASE_URL}${endpoint}`;
        const response = await fetch(url, {
            ...options,
            headers: {
                ...this.getHeaders(),
                ...options.headers,
            },
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `API error: ${response.status}`);
        }

        return response.json();
    }

    // Auth
    login(username, password) {
        return this.request("/login", {
            method: "POST",
            body: JSON.stringify({ username, password }),
        });
    }

    register(username, password, email) {
        return this.request("/register", {
            method: "POST",
            body: JSON.stringify({ username, password, email }),
        });
    }

    // Conversations
    getConversations() {
        return this.request("/conversations");
    }

    getConversation(id) {
        return this.request(`/conversations/${id}`);
    }

    createDM(targetUserId) {
        return this.request("/conversations/dm", {
            method: "POST",
            body: JSON.stringify({ targetUserId }),
        });
    }

    createGroup(name, isPrivate = false) {
        return this.request("/conversations/group", {
            method: "POST",
            body: JSON.stringify({ name, isPrivate }),
        });
    }

    // Messages
    sendMessage(conversationId, text) {
        return this.request(`/conversations/${conversationId}/messages`, {
            method: "POST",
            body: JSON.stringify({
                type: 0,
                content: { text }, // ← was JSON.stringify({ text }), now a plain object
            }),
        });
    }

    getMessages(conversationId, limit = 50, offset = 0) {
        return this.request(
            `/conversations/${conversationId}/messages?limit=${limit}&offset=${offset}`,
        );
    }

    getRecentMessages(conversationId, limit = 50) {
        return this.request(
            `/conversations/${conversationId}/messages/recent?limit=${limit}`,
        );
    }

    getUserByUsername(username) {
        return this.request(`/users?username=${encodeURIComponent(username)}`);
    }
}

export const api = new APIClient();
