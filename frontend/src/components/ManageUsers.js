import React, { useEffect, useState } from 'react';

const ManageUsers = () => {
    const [users, setUsers] = useState([]);

    useEffect(() => {
        // Simulating API call
        const fetchUsers = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/users'); // Replace with your API URL
                if (!response.ok) throw new Error('Failed to fetch users');
                const data = await response.json();
                setUsers(data);
            } catch (error) {
                console.error('Error fetching users:', error);
            }
        };

        fetchUsers();
    }, []);

    return (
        <div className="container">
            <h1>Manage Users</h1>
            <table border="1">
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                </tr>
                </thead>
                <tbody>
                {users.length > 0 ? (
                    users.map(user => (
                        <tr key={user.id}>
                            <td>{user.id}</td>
                            <td>{user.name}</td>
                            <td>{user.email}</td>
                            <td>{user.role}</td>
                        </tr>
                    ))
                ) : (
                    <tr><td colSpan="4">No users found.</td></tr>
                )}
                </tbody>
            </table>
        </div>
    );
};

export default ManageUsers;
