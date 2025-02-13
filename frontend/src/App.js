import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Home from './components/Home';
import AdminHome from './components/AdminHome';
import ManageUsers from './components/ManageUsers';
import Reports from './components/Reports';
import Login from './components/Login';
import Logout from './components/Logout';
import Registration from './components/Registration';
import PermissionDenied from './components/PermissionDenied';

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(localStorage.getItem('isLoggedIn') === 'true');
    const [isAdmin, setIsAdmin] = useState(localStorage.getItem('isAdmin') === 'true');
    const [csrfToken, setCsrfToken] = useState(localStorage.getItem('csrf_token'));
    const [loading, setLoading] = useState(true);  // Loading state

    const BASE_URL = "http://localhost:8080";

    useEffect(() => {
        const fetchCsrfToken = async () => {
            try {
                const res = await fetch(`${BASE_URL}/csrf-token`, {
                    method: 'GET',
                    credentials: 'include',
                });
                if (!res.ok) throw new Error('Failed to fetch CSRF token');
                const data = await res.json();
                setCsrfToken(data.csrf_token);
                localStorage.setItem('csrf_token', data.csrf_token);
            } catch (error) {
                console.error('Failed to fetch CSRF token:', error);
            }
        };

        const checkAuth = async () => {
            try {
                const response = await fetch(`${BASE_URL}/check-auth`, {
                    method: "GET",
                    credentials: "include",
                });

                if (!response.ok) {
                    console.error("Auth check failed: Response not OK", response);
                    setIsLoggedIn(false);
                    setIsAdmin(false);
                    localStorage.setItem('isLoggedIn', 'false');
                    localStorage.setItem('isAdmin', 'false');
                    setLoading(false);
                    return;
                }

                const data = await response.json();
                console.log("Auth check response:", data); // Debug log

                setIsLoggedIn(data.is_logged_in);
                setIsAdmin(data.is_admin);
                localStorage.setItem('isLoggedIn', data.is_logged_in ? 'true' : 'false');
                localStorage.setItem('isAdmin', data.is_admin ? 'true' : 'false');
                setLoading(false);
            } catch (error) {
                console.error("Auth check failed", error);
                setIsLoggedIn(false);
                setIsAdmin(false);
                localStorage.setItem('isLoggedIn', 'false');
                localStorage.setItem('isAdmin', 'false');
                setLoading(false);
            }
        };

        fetchCsrfToken();
        checkAuth();
    }, []);

    useEffect(() => {
        const hasReloaded = localStorage.getItem('hasReloaded') === 'true';

        // Reload logic: only once per session
        if (window.location.pathname === "/login" && !hasReloaded) {
            localStorage.setItem('hasReloaded', 'true');
            window.location.reload();
        }

        // Reload logic for Home or Admin page
        if ((window.location.pathname === "/" && isLoggedIn) || window.location.pathname === "/admin") {
            if (!hasReloaded) {
                localStorage.setItem('hasReloaded', 'true');
                window.location.reload();
            }
        }

        // Reset hasReloaded flag on logout or page refresh
        return () => {
            localStorage.removeItem('hasReloaded');
        };
    }, [isLoggedIn, isAdmin]);

    if (loading) {
        return <div>Loading...</div>;  // Show loading state while waiting for the authentication check
    }

    return (
        <Router>
            <Routes>
                <Route
                    path="/login"
                    element={<Login setIsLoggedIn={setIsLoggedIn} setIsAdmin={setIsAdmin} setCsrfToken={setCsrfToken} />}
                />

                <Route
                    path="/"
                    element={
                        isLoggedIn ? (
                            isAdmin ? <Navigate to="/admin" /> : <Home isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />

                <Route
                    path="/admin"
                    element={
                        isLoggedIn
                            ? (isAdmin ? <AdminHome setIsLoggedIn={setIsLoggedIn} /> : <Navigate to="/permission-denied" />)
                            : <Navigate to="/login" />
                    }
                />

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

                <Route path="/register" element={<Registration />} />
                <Route path="/permission-denied" element={<PermissionDenied />} />
                <Route path="/manage-users" element={isLoggedIn && isAdmin ? <ManageUsers /> : <Navigate to="/permission-denied" />} />
                <Route path="/reports" element={isLoggedIn && isAdmin ? <Reports /> : <Navigate to="/permission-denied" />} />
            </Routes>
        </Router>
    );
};

export default App;
