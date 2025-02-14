import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Payment = ({ isLoggedIn }) => {
    const navigate = useNavigate();
    const [cartId, setCartId] = useState('');

    const handleInitiatePayment = async () => {
        if (!isLoggedIn) {
            alert("Please log in to initiate payment.");
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/initiate-payment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': localStorage.getItem('csrf_token'),
                },
                body: JSON.stringify({
                    user_id: 4, // Replace with dynamic user ID
                    cart_id: cartId,
                }),
            });

            if (response.ok) {
                alert("Payment initiated successfully!");
                navigate('/receipt');
            } else {
                const data = await response.json();
                alert(`Error: ${data.message}`);
            }
        } catch (error) {
            console.error("Error initiating payment:", error);
            alert("An error occurred while initiating payment.");
        }
    };

    return (
        <div>
            <h2>Payment</h2>
            <input type="number" placeholder="Cart ID" value={cartId} onChange={(e) => setCartId(e.target.value)} />
            <button onClick={handleInitiatePayment}>Pay Now</button>
        </div>
    );
};

export default Payment;