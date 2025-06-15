import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Auth API
export const authAPI = {
  register: (userData) => api.post('/auth/register', userData),
  login: (credentials) => api.post('/auth/login', credentials),
  getCurrentUser: () => api.get('/auth/me'),
};

// Products API
export const productsAPI = {
  getAll: (params) => api.get('/products', { params }),
  getById: (id) => api.get(`/products/${id}`),
  create: (productData) => api.post('/products', productData),
  update: (id, productData) => api.put(`/products/${id}`, productData),
  delete: (id) => api.delete(`/products/${id}`),
  
  // Product attributes
  addColor: (id, colorData) => api.post(`/products/${id}/colors`, colorData),
  addSize: (id, sizeData) => api.post(`/products/${id}/sizes`, sizeData),
  updateInventory: (id, inventoryData) => api.post(`/products/${id}/inventory`, inventoryData),
  addImage: (id, imageData) => api.post(`/products/${id}/images`, imageData),
  deleteColor: (id, colorId) => api.delete(`/products/${id}/colors/${colorId}`),
  deleteSize: (id, sizeId) => api.delete(`/products/${id}/sizes/${sizeId}`),
  deleteImage: (id, imageId) => api.delete(`/products/${id}/images/${imageId}`),
};

// Categories API
export const categoriesAPI = {
  getAll: () => api.get('/categories'),
  getById: (id) => api.get(`/categories/${id}`),
  create: (categoryData) => api.post('/categories', categoryData),
  update: (id, categoryData) => api.put(`/categories/${id}`, categoryData),
  delete: (id) => api.delete(`/categories/${id}`),
};

// Cart API
export const cartAPI = {
  getCart: () => api.get('/cart'),
  addToCart: (item) => api.post('/cart', item),
  updateCartItem: (id, updates) => api.put(`/cart/${id}`, updates),
  removeFromCart: (id) => api.delete(`/cart/${id}`),
  clearCart: () => api.delete('/cart'),
};

// Wishlist API
export const wishlistAPI = {
  getWishlist: () => api.get('/wishlist'),
  addToWishlist: (item) => api.post('/wishlist', item),
  removeFromWishlist: (id) => api.delete(`/wishlist/${id}`),
  clearWishlist: () => api.delete('/wishlist'),
};

// Address API
export const addressAPI = {
  getAll: () => api.get('/addresses'),
  getById: (id) => api.get(`/addresses/${id}`),
  create: (addressData) => api.post('/addresses', addressData),
  update: (id, addressData) => api.put(`/addresses/${id}`, addressData),
  delete: (id) => api.delete(`/addresses/${id}`),
};

// Order API
export const orderAPI = {
  placeOrder: (orderData) => api.post('/orders', orderData),
  getAll: () => api.get('/orders'),
  getById: (id) => api.get(`/orders/${id}`),
  updateStatus: (id, status) => api.put(`/orders/${id}/status`, { status }),
};

export default api; 