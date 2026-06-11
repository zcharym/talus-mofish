import React from 'react'
import ReactDOM from 'react-dom/client'
import ManagementApp from './ManagementApp'
import AgentApp from './AgentApp'

const isAgent = window.location.pathname.startsWith('/agent')

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    {isAgent ? <AgentApp /> : <ManagementApp />}
  </React.StrictMode>,
)
