import React, { useState, useEffect } from 'react'; // Import useEffect
import Axios from 'axios';

const App = () => {
  const [data, setData] = useState(""); // Use square brackets for useState

  const getData = async () => {
    try {
      const response = await Axios.get("http://localhost:5000/getData");
      setData(response.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    getData();
  }, []); // Dependency array is fine here

  return (
    <div>{data}</div> // Render the data
  );
};

export default App;
