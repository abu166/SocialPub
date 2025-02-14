import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const Cart = ({ isLoggedIn }) => {
    const navigate = useNavigate();

    // State for cart functionality
    const [cartItems, setCartItems] = useState([]); // State to store cart items
    const [productId, setProductId] = useState('');
    const [quantity, setQuantity] = useState(1);
    const [cartId, setCartId] = useState('');

    // Fetch cart items from the backend
    useEffect(() => {
        const fetchCartItems = async () => {
            if (!isLoggedIn) return;

            try {
                const response = await fetch('http://localhost:8080/cart', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': localStorage.getItem('csrf_token'),
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    setCartItems(data.cartItems || []);
                } else {
                    console.error("Failed to fetch cart items:", response.statusText);
                    alert("Failed to load cart items. Please try again later.");
                }
            } catch (error) {
                console.error("Error fetching cart items:", error);
                alert("An unexpected error occurred while fetching cart items.");
            }
        };

        fetchCartItems();
    }, [isLoggedIn]);

    // Add item to cart
    const handleAddToCart = async () => {
        if (!isLoggedIn) {
            alert("Please log in to add items to the cart.");
            return;
        }

        if (!productId || !quantity || isNaN(productId) || isNaN(quantity) || quantity <= 0) {
            alert("Please enter a valid product ID and quantity.");
            return;
        }

        const payload = {
            user_id: 4, // Replace with dynamic user ID
            product_id: parseInt(productId, 10),
            quantity: parseInt(quantity, 10),
        };

        try {
            const response = await fetch('http://localhost:8080/cart/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': localStorage.getItem('csrf_token'),
                },
                body: JSON.stringify(payload),
            });

            let data;
            try {
                data = await response.json();
            } catch (parseError) {
                console.error("Failed to parse JSON response:", parseError);
                alert("An unexpected error occurred. Please try again later.");
                return;
            }

            if (response.ok) {
                alert(data.message || "Item added to cart successfully!");
                setProductId('');
                setQuantity(1);

                // Refresh cart items
                const updatedResponse = await fetch('http://localhost:8080/cart', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': localStorage.getItem('csrf_token'),
                    },
                });
                const updatedData = await updatedResponse.json();
                setCartItems(updatedData.cartItems || []);
            } else {
                alert(`Error: ${data.error || "Unknown error"}`);
            }
        } catch (error) {
            console.error("Error adding item to cart:", error);
            alert("An error occurred while adding the item to the cart.");
        }
    };

    // Initiate payment
    const handleInitiatePayment = async () => {
        if (!isLoggedIn) {
            alert("Please log in to initiate payment.");
            return;
        }

        if (!cartId || isNaN(cartId)) {
            alert("Please enter a valid cart ID.");
            return;
        }

        const payload = {
            user_id: 4, // Replace with dynamic user ID
            cart_id: parseInt(cartId, 10),
        };

        try {
            const response = await fetch('http://localhost:8080/initiate-payment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': localStorage.getItem('csrf_token'),
                },
                body: JSON.stringify(payload),
            });

            let data;
            try {
                data = await response.json();
            } catch (parseError) {
                console.error("Failed to parse JSON response:", parseError);
                alert("An unexpected error occurred. Please try again later.");
                return;
            }

            if (response.ok) {
                alert(data.message || "Payment initiated successfully!");
                navigate('/receipt');
            } else {
                alert(`Error: ${data.error || "Unknown error"}`);
            }
        } catch (error) {
            console.error("Error initiating payment:", error);
            alert("An error occurred while initiating payment.");
        }
    };

    return (
        <div className="cart-container">
            <h2>Cart</h2>

            {/* Display Cart Items */}
            <div className="cart-items">
                {cartItems.length === 0 ? (
                    <p>Your cart is empty.</p>
                ) : (
                    cartItems.map((item, index) => (
                        <div key={`${item.product_id}-${index}`} className="cart-item">
                            <strong>{item.name}</strong> - ${item.price} x {item.quantity}
                        </div>
                    ))
                )}
            </div>

            {/* Add Item to Cart */}
            <div className="add-to-cart">
                <h3>Add Item to Cart</h3>
                <input
                    type="number"
                    placeholder="Product ID"
                    value={productId}
                    onChange={(e) => setProductId(e.target.value)}
                />
                <input
                    type="number"
                    placeholder="Quantity"
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                />
                <button onClick={handleAddToCart}>Add to Cart</button>
            </div>

            {/* Initiate Payment */}
            <div className="initiate-payment">
                <h3>Initiate Payment</h3>
                <input
                    type="number"
                    placeholder="Cart ID"
                    value={cartId}
                    onChange={(e) => setCartId(e.target.value)}
                />
                <button onClick={handleInitiatePayment}>Pay Now</button>
            </div>
        </div>
    );
};

export default Cart;