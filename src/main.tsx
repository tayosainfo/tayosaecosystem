import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';
import './utils/debugAuth'; // Load debug utilities
import './utils/apiInterceptor'; // Load API interceptor for 401 handling

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
