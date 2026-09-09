import React from 'react';
import ReactDOM from 'react-dom/client';
import { CrashBoundary } from './components/CrashBoundary';
import ManagementApp from './ManagementApp';
import AgentApp from './AgentApp';
import './theme/fonts.css';
import './theme/atmosphere.css';

const isAgent = window.location.pathname.startsWith('/agent');

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <CrashBoundary>{isAgent ? <AgentApp /> : <ManagementApp />}</CrashBoundary>
  </React.StrictMode>,
);
