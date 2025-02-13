import React, { useEffect, useState } from 'react';

const Reports = () => {
    const [reports, setReports] = useState([]);

    useEffect(() => {
        // Simulating API call
        const fetchReports = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/reports'); // Replace with your API URL
                if (!response.ok) throw new Error('Failed to fetch reports');
                const data = await response.json();
                setReports(data);
            } catch (error) {
                console.error('Error fetching reports:', error);
            }
        };

        fetchReports();
    }, []);

    return (
        <div className="container">
            <h1>View Reports</h1>
            <ul>
                {reports.length > 0 ? (
                    reports.map(report => (
                        <li key={report.id}>
                            <strong>{report.title}</strong>: {report.description}
                        </li>
                    ))
                ) : (
                    <li>No reports available.</li>
                )}
            </ul>
        </div>
    );
};

export default Reports;
