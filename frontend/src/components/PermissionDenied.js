import React from 'react';
import { Link } from 'react-router-dom';

const PermissionDenied = () => {
    return (
        <div style={{ textAlign: 'center', marginTop: '50px' }}>
            <h2>403 - Permission Denied</h2>
            <p>У вас нет прав доступа к этой странице.</p>
            <Link to="/">Вернуться на главную</Link>
        </div>
    );
};

export default PermissionDenied;
