import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Home from './components/Home';
import Login from './components/Login';
import Logout from './components/Logout';
import Registration from './components/Registration';

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false); // Login state
    const [csrfToken, setCsrfToken] = useState(null); // CSRF token state

    // Fetch the CSRF token on component mount
    useEffect(() => {
        const fetchCsrfToken = async () => {
            try {
                const response = await fetch("http://localhost:8080/csrf-token", {
                    method: "GET",
                    credentials: "include", // Include cookies
                });

                if (response.ok) {
                    const data = await response.json();
                    console.log("CSRF Token fetched:", data.csrf_token);
                    setCsrfToken(data.csrf_token);
                } else {
                    console.error("Failed to fetch CSRF token. Status:", response.status);
                }
            } catch (error) {
                console.error("Error fetching CSRF token:", error);
            }
        };

        fetchCsrfToken();
    }, []);



    return (
        <Router>
            <Routes>
                {/* Home Route */}
                <Route
                    path="/"
                    element={
                        isLoggedIn ? (
                            <Home isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />

                {/* Login Route */}
                <Route
                    path="/login"
                    element={<Login setIsLoggedIn={setIsLoggedIn} setCsrfToken={setCsrfToken} />}
                />

                {/* Logout Route */}
                <Route
                    path="/logout"
                    element={
                        isLoggedIn ? (
                            <Logout setIsLoggedIn={setIsLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />

                {/* Registration Route */}
                <Route path="/register" element={<Registration />} />
            </Routes>
        </Router>
    );
};

export default App;
