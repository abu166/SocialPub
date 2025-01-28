import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Login.css';

const Login = ({ setIsLoggedIn, setCsrfToken }) => { // Accept setCsrfToken as a prop
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const navigate = useNavigate();

    const handleLogin = async () => {
        if (!username || !password) {
            setMessage("All fields are required!");
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
                credentials: "include",
            });

            console.log("Response status:", response.status); // Debugging
            const data = await response.json();
            console.log("Response data:", data); // Debugging
            console.log("CSRF Token received:", data.csrf_token);
            localStorage.setItem('csrf_token', data.csrf_token);

            if (response.ok) {
                setCsrfToken(data.csrfToken);
                setIsLoggedIn(true);
                setMessage(data.message || "Login successful!");
                navigate("/");
            } else {
                setMessage(data.message || "Login failed!");
            }
        } catch (error) {
            console.error("Login error:", error); // Debugging
            setMessage("An error occurred. Please try again later.");
        }
    };


    return (
        <div className="auth-container">
            <h1>Log in with your Instagram account</h1>
            <input
                type="text"
                placeholder="Username, phone or email"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button className="auth-button" onClick={handleLogin}>
                Log in
            </button>
            {message && <p className="message">{message}</p>}
            <p>Forgotten password?</p>

            {/* Add a Registration Link */}
            <div className="signup-link">
                <p>Don't have an account? <a href="/register">Sign up</a></p>
            </div>
        </div>
    );
};

export default Login;
