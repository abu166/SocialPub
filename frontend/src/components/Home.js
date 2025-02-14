import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/Home.css';

const Home = ({ isLoggedIn, setIsLoggedIn }) => {
    const navigate = useNavigate();

    // Logout handler
    const handleLogout = async () => {
        try {
            const csrfToken = localStorage.getItem('csrf_token');
            console.log("Attempting to log out with CSRF Token:", csrfToken);
            const response = await fetch("http://localhost:8080/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": csrfToken,
                },
                credentials: "include",
            });
            if (response.ok) {
                console.log("Logout successful!");
                setIsLoggedIn(false);
                navigate("/");
            } else {
                const data = await response.json();
                console.error("Logout failed:", data.message);
                alert(data.message || "Logout failed!");
            }
        } catch (error) {
            console.error("Error during logout:", error);
            alert("An error occurred. Please try again later.");
        }
    };

    return (
        <div className="home-container">
            {/* Left Sidebar */}
            <nav className="sidebar">
                <ul>
                    <li><Link to="/">Home</Link></li>
                    <li><Link to="/search">Search</Link></li>
                    <li><Link to="/create-post">+</Link></li>
                    <li><Link to="/profile">Profile</Link></li>
                    <li><Link to="/cart">Cart</Link></li> {/* Add Cart route */}
                </ul>
            </nav>

            {/* Main Content Area */}
            <main className="main-content">
                <header className="main-header">
                    <h2>Home</h2>
                    {/* Conditional rendering of Login/Logout button */}
                    {!isLoggedIn ? (
                        <Link to="/login" className="login-button">Log in</Link>
                    ) : (
                        <button className="logout-button" onClick={handleLogout}>Log out</button>
                    )}
                </header>

                {/* Posts Section */}
                <div className="posts">
                    <div className="post">
                        <div className="post-header">
                            <strong>@nishajwrites</strong> <span>11h</span>
                        </div>
                        <p>If he didn't mean it that way how come he hasn't apologized?</p>
                        <div className="post-footer">
                            <span>❤️ 479</span> <span>💬 7</span> <span>🔄 36</span>
                        </div>
                    </div>
                    <div className="post">
                        <div className="post-header">
                            <strong>@morallygreykay</strong> <span>16h</span>
                        </div>
                        <p>If you're really paying attention...</p>
                        <div className="post-footer">
                            <span>❤️ 704</span> <span>💬 9</span> <span>🔄 55</span>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Home;