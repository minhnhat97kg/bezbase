import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import { OrganizationProvider } from './context/OrganizationContext';
import { useAuth } from './hooks/useAuth';
import { useTranslation } from 'react-i18next';
import './i18n';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import UserManagement from './pages/UserManagement';
import UnifiedRBAC from './pages/UnifiedRBAC';
import OrganizationManagement from './pages/OrganizationManagement';
import CreateOrganization from './pages/CreateOrganization';
import OrganizationDebug from './components/debug/OrganizationDebug';

function App() {
  return (
    <AuthProvider>
      <ThemeProvider>
        <OrganizationProvider>
          <Router>
            <AppLayout />
          </Router>
        </OrganizationProvider>
      </ThemeProvider>
    </AuthProvider>
  );
}

function AppLayout() {
  const { user, loading } = useAuth();
  const { t } = useTranslation();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  if (loading) {
    return (
      <div className="min-h-screen flex justify-center items-center bg-gray-50 dark:bg-gray-900">
        <div className="text-lg text-gray-600 dark:text-gray-300">{t('common.loading')}</div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="*" element={<Navigate to="/login" />} />
        </Routes>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex bg-gray-50 dark:bg-gray-900">
      <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
      <div className={`flex-1 flex flex-col transition-all duration-300 ${sidebarOpen ? 'ml-64' : 'ml-0'}`}>
        <Header sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
        <main className="flex-1 p-6">
          <Routes>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/users" element={<UserManagement />} />
            <Route path="/rbac" element={<UnifiedRBAC />} />
            {/* Redirect old routes to unified RBAC */}
            <Route path="/roles" element={<Navigate to="/rbac" replace />} />
            <Route path="/advanced-rbac" element={<Navigate to="/rbac" replace />} />
            <Route path="/organizations/manage" element={<OrganizationManagement />} />
            <Route path="/organizations/create" element={<CreateOrganization />} />
            <Route path="/" element={<Navigate to="/dashboard" />} />
            <Route path="*" element={<Navigate to="/dashboard" />} />
          </Routes>
        </main>
      </div>
      {/* <OrganizationDebug /> */}
    </div>
  );
}

export default App;

