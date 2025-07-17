import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Icon from './common/Icons';

const Sidebar = ({ sidebarOpen, setSidebarOpen }) => {
  const { user, logout } = useAuth();
  const location = useLocation();

  const handleLogout = () => {
    logout();
  };

  const navigation = [
    {
      name: 'Dashboard',
      href: '/dashboard',
      icon: <Icon name="dashboard" />,
    },
    {
      name: 'Profile',
      href: '/profile',
      icon: <Icon name="user" />,
    },
    {
      name: 'User Management',
      href: '/users',
      icon: <Icon name="users" />,
    },
    {
      name: 'Role Management',
      href: '/roles',
      icon: <Icon name="shield" />,
    },
  ];

  const isActive = (href) => location.pathname === href;

  if (!user) {
    return null; // Don't show sidebar if user is not logged in
  }

  return (
    <div
      className={`fixed inset-y-0 left-0 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 min-h-screen transform ${
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      } transition-transform duration-300 ease-in-out z-30`}
    >
      {/* Logo/Brand */}
      <div className="flex items-center h-16 px-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center">
          <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mr-3">
            <Icon name="rocket" className="text-white" />
          </div>
          <h1 className="text-xl font-semibold text-gray-900 dark:text-white">BezBase</h1>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 py-6 space-y-1">
        {navigation.map((item) => (
          <Link
            key={item.name}
            to={item.href}
            className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${isActive(item.href)
                ? 'bg-primary-50 text-primary-700 border-r-2 border-primary-600 dark:bg-primary-800 dark:text-primary-200 dark:border-primary-500'
                : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white'
              }`}
          >
            <span className={`mr-3 ${isActive(item.href) ? 'text-primary-600 dark:text-primary-400' : 'text-gray-400 dark:text-gray-500'}`}>
              {item.icon}
            </span>
            {item.name}
          </Link>
        ))}
      </nav>

      {/* User Actions */}
      <div className="px-4 py-4 border-t border-gray-200 dark:border-gray-700">
        <div className="flex items-center mb-4 px-3">
          <div className="w-8 h-8 bg-gray-300 dark:bg-gray-700 rounded-full flex items-center justify-center mr-3">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {user?.first_name?.[0]}{user?.last_name?.[0]}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
              {user?.first_name} {user?.last_name}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{user?.email}</p>
          </div>
        </div>
        <button
          onClick={handleLogout}
          className="flex items-center w-full px-3 py-2 text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white transition-colors"
        >
          <Icon name="exit" className="mr-3 text-gray-400 dark:text-gray-500" />
          Sign out
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
