import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Home from './components/Home';
import Login from './components/Login';
import Logout from './components/Logout';
import Registration from './components/Registration';
import EmailForm from "./components/EmailForm";

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false); // Login state
    const [csrfToken, setCsrfToken] = useState(null); // CSRF token state

    // Fetch the CSRF token on component mount
    useEffect(() => {

        const BASE_URL = "http://localhost:8080";

        const fetchCsrfToken = async () => {
            try {
                const res = await fetch(`${BASE_URL}/csrf-token`, {
                    method: 'GET',
                    credentials: 'include',
                });
                if (!res.ok) throw new Error('Failed to fetch CSRF token');
                const data = await res.json();
                setCsrfToken(data.csrf_token);
            } catch (error) {
                console.error('Failed to fetch CSRF token:', error);
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

                <Route path="/send-email" element={<EmailForm />} />
            </Routes>
        </Router>
    );
};

export default App;
