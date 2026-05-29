import { useEffect, useState } from "react";
import { api } from "../api/client";

export function useAuth() {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem("token");
        const userId = localStorage.getItem("userId");

        if (token && userId) {
            setUser({ id: parseInt(userId), token });
        }
        setLoading(false);
    }, []);

    const login = async (username, password) => {
        try {
            setError(null);
            const data = await api.login(username, password);
            api.setToken(data.api_key);
            localStorage.setItem("token", data.api_key);
            localStorage.setItem("userId", data.user_id);
            setUser({ id: data.user_id, token: data.api_key });
            return data;
        } catch (err) {
            setError(err.message);
            throw err;
        }
    };

    const register = async (username, password, email) => {
        try {
            setError(null);
            const data = await api.register(username, password, email);
            api.setToken(data.api_key);
            localStorage.setItem("token", data.api_key);
            localStorage.setItem("userId", data.user_id);
            setUser({ id: data.user_id, token: data.api_key });
            return data;
        } catch (err) {
            setError(err.message);
            throw err;
        }
    };

    const logout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("userId");
        setUser(null);
        api.setToken(null);
    };

    return { user, loading, error, login, register, logout };
}
