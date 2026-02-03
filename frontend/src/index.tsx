import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App.tsx';
import reportWebVitals from './reportWebVitals.ts';

console.log('index.tsx loading...');

const rootElement = document.getElementById('root');
console.log('Root element:', rootElement);

if (!rootElement) {
  throw new Error('Root element not found');
}

// Ensure root is visible
rootElement.style.display = 'flex';
rootElement.style.flexDirection = 'column';
rootElement.style.width = '100%';
rootElement.style.height = '100%';

const root = ReactDOM.createRoot(rootElement);
console.log('ReactDOM root created');

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);

console.log('App rendered');

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
