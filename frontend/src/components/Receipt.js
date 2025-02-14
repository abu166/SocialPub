import React, { useEffect, useState } from 'react';

const Receipt = ({ isLoggedIn }) => {
    const [receipt, setReceipt] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchReceipt = async () => {
            if (!isLoggedIn) {
                setError("Please log in to view the receipt.");
                return;
            }

            try {
                const response = await fetch('http://localhost:8081/receipt?transactionID=1', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': localStorage.getItem('csrf_token'),
                    },
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }

                const data = await response.json();
                setReceipt(data);
            } catch (err) {
                console.error("Error fetching receipt:", err);
                setError("Failed to fetch receipt. Please try again later.");
            }
        };

        fetchReceipt();
    }, [isLoggedIn]);

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div>
            <h2>Receipt</h2>
            {receipt ? (
                <div>
                    <p>Transaction ID: {receipt.transactionID}</p>
                    <p>Status: {receipt.status}</p>
                </div>
            ) : (
                <p>Loading receipt...</p>
            )}
        </div>
    );
};

export default Receipt;