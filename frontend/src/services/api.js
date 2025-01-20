import axios from 'axios';

const API = axios.create({ baseURL: process.env.REACT_APP_API_URL });

export const fetchPosts = () => API.get('/posts');
export const login = (credentials) => API.post('/auth/login', credentials);
export const signup = (userData) => API.post('/auth/signup', userData);
