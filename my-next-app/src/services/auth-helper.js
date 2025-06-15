// Helper functions for handling authentication and tokens

/**
 * Gets the authentication token from localStorage
 * @returns {string|null} The auth token or null if not found
 */
export const getToken = () => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
};

/**
 * Saves the authentication token to localStorage
 * @param {string} token - The token to save
 */
export const setToken = (token) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('token', token);
  }
};

/**
 * Removes the authentication token from localStorage
 */
export const removeToken = () => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
  }
};

/**
 * Checks if a user is authenticated based on token presence
 * @returns {boolean} True if authenticated, false otherwise
 */
export const isAuthenticated = () => {
  return !!getToken();
};

/**
 * Parse JWT token to get user data
 * @param {string} token - JWT token
 * @returns {object|null} Decoded token payload or null if invalid
 */
export const parseToken = (token) => {
  try {
    if (!token) return null;
    
    // Get the payload part of the JWT
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(window.atob(base64));
    
    // Check if token is expired
    if (payload.exp && Date.now() >= payload.exp * 1000) {
      removeToken();
      return null;
    }
    
    return payload;
  } catch (error) {
    console.error('Error parsing token:', error);
    return null;
  }
};

/**
 * Get user roles from JWT token
 * @returns {array} Array of user roles or empty array if not found
 */
export const getUserRoles = () => {
  const token = getToken();
  if (!token) return [];
  
  const payload = parseToken(token);
  return payload?.roles || [];
};

/**
 * Check if current user has admin role
 * @returns {boolean} True if user is admin, false otherwise
 */
export const isAdmin = () => {
  const roles = getUserRoles();
  return roles.includes('admin');
};

export default {
  getToken,
  setToken,
  removeToken,
  isAuthenticated,
  parseToken,
  getUserRoles,
  isAdmin
}; 