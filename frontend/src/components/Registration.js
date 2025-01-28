import React, { useState } from 'react';
import '../styles/Registration.css';

const Registration = () => {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSignUp = async () => {
        if (!username || !email || !password) {
            setMessage("All fields are required!");
            return;
        }
        if (!/\S+@\S+\.\S+/.test(email)) {
            setMessage("Please enter a valid email address!");
            return;
        }
        if (password.length < 6) {
            setMessage("Password must be at least 6 characters long!");
            return;
        }

        console.log("Starting registration...");
        setLoading(true);
        setMessage("");

        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, email, password }),
                credentials: "include",
            });

            console.log("Response:", response);
            const data = await response.json();
            console.log("Response Data:", data);

            if (response.ok) {
                setMessage(data.message || "Registration successful!");
                setTimeout(() => {
                    window.location.href = "/login";
                }, 2000);
            } else {
                setMessage(data.message || "Registration failed!");
            }
        } catch (error) {
            console.error("Error during registration:", error);
            setMessage("An error occurred. Please try again later.");
        } finally {
            setLoading(false);
        }
    };


    return (
        <div className="auth-container">
            <h1>Create your account</h1>
            <input
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
            />
            <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button
                className="auth-button"
                onClick={handleSignUp}
                disabled={loading}
            >
                {loading ? 'Signing up...' : 'Sign up'}
            </button>
            {message && <p className="message">{message}</p>}
        </div>
    );
};

export default Registration;
