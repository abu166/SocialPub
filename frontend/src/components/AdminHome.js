import React from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles/Admin.css";

const AdminHome = ({ setIsLoggedIn }) => {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            const csrfToken = localStorage.getItem("csrf_token");
            const response = await fetch("http://localhost:8080/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": csrfToken,
                },
                credentials: "include",
            });

            if (response.ok) {
                setIsLoggedIn(false);
                navigate("/login");
            } else {
                alert("Logout failed!");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    };

    return (
        <div className="admin-container">
            <nav className="sidebar">
                <ul>
                    <li><Link to="/admin">Admin Dashboard</Link></li>
                    <li><Link to="/manage-users">Manage Users</Link></li>
                    <li><Link to="/reports">View Reports</Link></li>
                </ul>
            </nav>
            <main className="admin-main">
                <header className="admin-header">
                    <h2>Admin Dashboard</h2>
                    <button className="logout-button" onClick={handleLogout}>Log out</button>
                </header>
                <div className="admin-content">
                    <p>Welcome, Admin! Manage users, view reports, and control the system.</p>
                </div>
            </main>
        </div>
    );
};

export default AdminHome;
