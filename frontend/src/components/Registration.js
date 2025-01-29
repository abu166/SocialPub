import React, { useState } from "react";
import "../styles/Registration.css";

const Registration = () => {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [message, setMessage] = useState("");
    const [loading, setLoading] = useState(false);
    const [verificationCode, setVerificationCode] = useState("");
    const [isVerificationStep, setIsVerificationStep] = useState(false);

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

            const text = await response.text();
            let data;
            try {
                data = JSON.parse(text);
            } catch {
                throw new Error(text);
            }

            if (response.ok) {
                setMessage("A verification code has been sent to your email.");
                setIsVerificationStep(true);
            } else {
                setMessage(data.message || "Registration failed!");
            }
        } catch (error) {
            console.error("Error during registration:", error);
            setMessage(error.message || "An error occurred. Please try again later.");
        } finally {
            setLoading(false);
        }
    };

    const handleVerifyCode = async () => {
        if (!verificationCode) {
            setMessage("Please enter the verification code!");
            return;
        }

        setLoading(true);
        setMessage("");

        try {
            const response = await fetch("http://localhost:8080/verify-email", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, code: verificationCode }),
                credentials: "include",
            });

            const text = await response.text();
            let data;
            try {
                data = JSON.parse(text);
            } catch {
                throw new Error(text);
            }

            if (response.ok) {
                setMessage("Email verified successfully! Redirecting to login...");
                setTimeout(() => {
                    window.location.href = "/login";
                }, 2000);
            } else {
                setMessage(data.message || "Invalid verification code!");
            }
        } catch (error) {
            console.error("Error verifying code:", error);
            setMessage(error.message || "An error occurred. Please try again later.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-container">
            {!isVerificationStep ? (
                <>
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
                        {loading ? "Signing up..." : "Sign up"}
                    </button>
                </>
            ) : (
                <>
                    <h1>Email Verification</h1>
                    <p>A verification code has been sent to your email.</p>
                    <input
                        type="text"
                        placeholder="Enter verification code"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                    />
                    <button
                        className="auth-button"
                        onClick={handleVerifyCode}
                        disabled={loading}
                    >
                        {loading ? "Verifying..." : "Verify"}
                    </button>
                </>
            )}
            {message && <p className="message">{message}</p>}
        </div>
    );
};

export default Registration;
