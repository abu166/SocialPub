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
import Cart from './components/Cart';
import Payment from './components/Payment';
import Receipt from './components/Receipt';

const App = () => {
    // State for authentication and CSRF token
    const [isLoggedIn, setIsLoggedIn] = useState(localStorage.getItem('isLoggedIn') === 'true');
    const [isAdmin, setIsAdmin] = useState(localStorage.getItem('isAdmin') === 'true');
    const [csrfToken, setCsrfToken] = useState(localStorage.getItem('csrf_token'));
    const [loading, setLoading] = useState(true);

    const BASE_URL = "http://localhost:8080";

    // Fetch CSRF token on app load
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

        // Check user authentication status
        const checkAuth = async () => {
            try {
                const response = await fetch(`${BASE_URL}/check-auth`, {
                    method: "GET",
                    credentials: "include",
                });
                if (!response.ok) {
                    setIsLoggedIn(false);
                    setIsAdmin(false);
                    localStorage.setItem('isLoggedIn', 'false');
                    localStorage.setItem('isAdmin', 'false');
                    setLoading(false);
                    return;
                }
                const data = await response.json();
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

    // Show loading state while checking auth and fetching CSRF token
    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <Router>
            <Routes>
                {/* Public Routes */}
                <Route path="/login" element={<Login setIsLoggedIn={setIsLoggedIn} setIsAdmin={setIsAdmin} setCsrfToken={setCsrfToken} />} />
                <Route path="/register" element={<Registration />} />
                <Route path="/permission-denied" element={<PermissionDenied />} />

                {/* Protected Routes */}
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
                <Route
                    path="/cart"
                    element={
                        isLoggedIn ? (
                            <Cart isLoggedIn={isLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />
                <Route
                    path="/payment"
                    element={
                        isLoggedIn ? (
                            <Payment isLoggedIn={isLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />
                <Route
                    path="/receipt"
                    element={
                        isLoggedIn ? (
                            <Receipt isLoggedIn={isLoggedIn} csrfToken={csrfToken} />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                />

                {/* Admin Routes */}
                <Route
                    path="/admin"
                    element={
                        isLoggedIn && isAdmin ? (
                            <AdminHome setIsLoggedIn={setIsLoggedIn} />
                        ) : (
                            <Navigate to="/permission-denied" />
                        )
                    }
                />
                <Route
                    path="/manage-users"
                    element={
                        isLoggedIn && isAdmin ? (
                            <ManageUsers />
                        ) : (
                            <Navigate to="/permission-denied" />
                        )
                    }
                />
                <Route
                    path="/reports"
                    element={
                        isLoggedIn && isAdmin ? (
                            <Reports />
                        ) : (
                            <Navigate to="/permission-denied" />
                        )
                    }
                />
            </Routes>
        </Router>
    );
};

export default App;