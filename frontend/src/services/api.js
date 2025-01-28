const BASE_URL = "http://localhost:8080";

export const api = {
    register: async (data) => {
        const res = await fetch(`${BASE_URL}/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
            credentials: "include",
        });

        if (!res.ok) {
            const error = await res.json();
            throw new Error(error.message || "Failed to register");
        }

        return res.json();
    },

    login: async (data) => {
        const res = await fetch(`${BASE_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
            credentials: "include", // Include credentials to manage session cookies
        });

        if (!res.ok) {
            const error = await res.json();
            throw new Error(error.message || "Login failed");
        }

        // Return the response, assuming it includes necessary tokens or session info
        return res.json();
    },

    logout: async () => {
        const csrfToken = document.cookie
            .split("; ")
            .find((row) => row.startsWith("csrf_token="))
            ?.split("=")[1];

        if (!csrfToken) {
            throw new Error("CSRF token is missing!");
        }

        const headers = {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        };

        const res = await fetch(`${BASE_URL}/logout`, {
            method: "POST",
            headers,
            credentials: "include",
        });

        const contentType = res.headers.get("Content-Type");

        if (!res.ok) {
            if (contentType && contentType.includes("application/json")) {
                const error = await res.json();
                throw new Error(error.message || "Failed to log out");
            } else {
                const errorText = await res.text();
                throw new Error(errorText.trim() || "Failed to log out");
            }
        }

        return contentType && contentType.includes("application/json")
            ? res.json()
            : res.text();
    },

    csrf_token: async () => {
        const response = await fetch(`${BASE_URL}/csrf-token`, {
            method: 'GET',
            credentials: "include",
            headers: {
                'X-CSRF-Token': localStorage.getItem('csrfToken') || '',
            },
        });
        const data = await response.json();
        return data.csrfToken;
    },

    protected: async (username) => {
        // Extract the CSRF token from cookies
        const csrfToken = document.cookie
            .split("; ")
            .find((row) => row.startsWith("X-CSRF-Token="))
            ?.split("=")[1];

        if (!csrfToken) {
            throw new Error("CSRF token is missing!");
        }

        const res = await fetch(`${BASE_URL}/protected`, {
            method: "POST",
            headers: {
                "X-CSRF-Token": csrfToken, // Include the CSRF token in the headers
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: new URLSearchParams({ username }),
            credentials: "include", // Ensure credentials (cookies) are sent
        });

        if (!res.ok) {
            const error = await res.json();
            throw new Error(error.message || "Failed to fetch protected resource");
        }

        return res.json();
    },
};

// Example logout handler
export const handleLogout = async () => {
    try {
        await api.logout();
        window.location.href = "/login"; // Redirect to login page
    } catch (error) {
        console.error("Logout failed:", error.message);
    }
};
