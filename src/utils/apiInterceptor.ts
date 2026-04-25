/**
 * API Interceptor for handling 401 errors globally
 * 
 * This wraps fetch to automatically handle unauthorized responses
 * by dispatching a custom event that the AuthContext listens to.
 */

const originalFetch = window.fetch;

window.fetch = async (...args) => {
  const response = await originalFetch(...args);
  
  // If we get a 401, dispatch an event for the AuthContext to handle
  if (response.status === 401) {
    console.warn('401 Unauthorized response detected');
    window.dispatchEvent(
      new CustomEvent('auth:unauthorized', {
        detail: { status: 401, url: args[0] },
      })
    );
  }
  
  return response;
};

export {}; // Make this a module
