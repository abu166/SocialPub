import React from 'react';
import ReactDOM from 'react-dom/client'; // Use the new root API
import App from './App';
import './index.css';

const root = ReactDOM.createRoot(document.getElementById('root')); // Modern API
root.render(<App />);
