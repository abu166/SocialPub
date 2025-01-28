import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Logout.css';

const Logout = ({ setIsLoggedIn, csrfToken }) => { // Accept csrfToken as a prop
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            const response = await fetch("http://localhost:8080/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": csrfToken, // Use the CSRF token passed as a prop
                    // "X-Username": username,
                },
                credentials: "include", // Ensure cookies are included
            });


            if (response.ok) {
                // Successfully logged out, reset logged-in state
                setIsLoggedIn(false);

                // Redirect to home page after logout
                navigate("/");
            } else {
                const data = await response.json();
                alert(data.message || "Logout failed!");
            }
        } catch (error) {
            console.error("Error logging out:", error);
            alert("An error occurred. Please try again later.");
        }
    };

    return (
        <div className="auth-container">
            <h1>Are you sure you want to log out?</h1>
            <button className="auth-button" onClick={handleLogout}>
                Log out
            </button>
        </div>
    );
};

export default Logout;
