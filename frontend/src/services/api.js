import axios from 'axios';
import env from '../env/env';
import { getToken } from '../utils/cookies';

const api = axios.create({
  baseURL: env.API_BASE_URL,
});

api.interceptors.request.use(
  async (config) => {
    const token = await getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth APIs
export const register = (email, password, name) => {
  return api.post('/auth/register', { email, password, name });
};

export const login = (email, password) => {
  return api.post('/auth/login', { email, password });
};

// Task APIs
export const getTasks = () => {
  return api.get('/tasks');
};

export const createTask = (taskData) => {
  return api.post('/tasks', taskData);
};

export const updateTask = (id, taskData) => {
  return api.put(`/tasks/${id}`, taskData);
};

export const deleteTask = (id) => {
  return api.delete(`/tasks/${id}`);
};

export const searchTasks = (searchTerm) => {
  return api.get(`/tasks/search`, {
    params: { q: searchTerm }
  });
};

// User APIs
export const getCurrentUser = () => {
  return api.get('/users/me');
};

export const updateProfile = (userId, profileData) => {
  return api.put(`/users/${userId}/profile`, profileData);
};

// VULNERABILITY #2: Admin endpoint accessible without proper authorization check
// VULNERABILITY #4: Hardcoded API key sent in request
export const getAllUsers = () => {
  return api.get('/admin/users', {
    headers: {
      'X-Admin-Key': env.ADMIN_API_KEY
    }
  });
};

export default api;
