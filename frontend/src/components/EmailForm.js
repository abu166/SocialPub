import React, { useState } from "react";
import "../styles/EmailForm.css"; // Add styles here

const EmailForm = () => {
    const [email, setEmail] = useState("");
    const [message, setMessage] = useState("");
    const [attachment, setAttachment] = useState(null);
    const [responseMessage, setResponseMessage] = useState("");
    const [error, setError] = useState("");

    const handleFileChange = (e) => {
        setAttachment(e.target.files[0]);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setResponseMessage("");

        const formData = new FormData();
        formData.append("email", email);
        formData.append("message", message);
        if (attachment) {
            formData.append("attachment", attachment);
        }

        try {
            const response = await fetch("http://localhost:8000/send-email", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Something went wrong.");
            }

            const responseData = await response.json();
            setResponseMessage(responseData.message);
            setEmail("");
            setMessage("");
            setAttachment(null);
        } catch (err) {
            setError(err.message);
        }
    };

    return (
        <div className="email-form-container">
            <h1>Send Support Request</h1>
            {responseMessage && <p className="success">{responseMessage}</p>}
            {error && <p className="error">{error}</p>}
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Email:</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Your email"
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Message:</label>
                    <textarea
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Your message"
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Attachment (optional):</label>
                    <input type="file" onChange={handleFileChange} />
                </div>
                <button type="submit" className="submit-button">
                    Send Email
                </button>
            </form>
        </div>
    );
};

export default EmailForm;
