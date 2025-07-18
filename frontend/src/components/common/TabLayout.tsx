import React, { ReactNode } from 'react';
import Icon from './Icons';

interface Tab {
  id: string;
  label: string;
  icon?: string;
  disabled?: boolean;
}

interface TabLayoutProps {
  title?: string;
  subtitle?: string;
  children: ReactNode;
  tabs: Tab[];
  activeTab: string;
  onTabChange: (tab: string) => void;
  showTabs?: boolean;
  headerActions?: ReactNode;
  message?: string;
  error?: string;
  onMessageDismiss?: () => void;
  onErrorDismiss?: () => void;
  className?: string;
}

/**
 * TabLayout - A reusable layout component for pages with tabbed navigation
 * 
 * This component provides a consistent layout structure for pages that need:
 * - A header with title and subtitle
 * - Tabbed navigation interface
 * - Success/error message display
 * - Header actions (buttons, etc.)
 * 
 * Perfect for admin panels, settings pages, management interfaces, etc.
 */
const TabLayout: React.FC<TabLayoutProps> = ({
  title,
  subtitle,
  children,
  tabs,
  activeTab,
  onTabChange,
  showTabs = true,
  headerActions = null,
  message = '',
  error = '',
  onMessageDismiss,
  onErrorDismiss,
  className = '',
}) => {
  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      {(title || headerActions) && (
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
            <div className="flex justify-between items-start">
              <div>
                {title && <h1 className="text-2xl font-semibold text-gray-900">{title}</h1>}
                {subtitle && <p className="text-gray-600 mt-1">{subtitle}</p>}
              </div>
              {headerActions && (
                <div className="flex items-center space-x-2">
                  {headerActions}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Tabs Navigation */}
      {showTabs && tabs.length > 0 && (
        <div className="bg-white shadow rounded-lg">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8 px-6" aria-label="Tabs">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => !tab.disabled && onTabChange && onTabChange(tab.id)}
                  disabled={tab.disabled}
                  className={`py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap flex items-center ${
                    activeTab === tab.id
                      ? 'border-indigo-500 text-indigo-600'
                      : tab.disabled
                      ? 'border-transparent text-gray-400 cursor-not-allowed'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                  role="tab"
                  aria-selected={activeTab === tab.id}
                  aria-controls={`tabpanel-${tab.id}`}
                >
                  {tab.icon && (
                    <span className="mr-2">{tab.icon}</span>
                  )}
                  {tab.label}
                </button>
              ))}
            </nav>
          </div>
        </div>
      )}

      {/* Messages */}
      {message && (
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <span className="text-sm text-green-700">{message}</span>
            </div>
            {onMessageDismiss && (
              <button
                onClick={onMessageDismiss}
                className="text-green-400 hover:text-green-600 transition-colors"
                aria-label="Dismiss success message"
              >
                ×
              </button>
            )}
          </div>
        </div>
      )}

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <span className="text-sm text-red-700">{error}</span>
            </div>
            {onErrorDismiss && (
              <button
                onClick={onErrorDismiss}
                className="text-red-400 hover:text-red-600 transition-colors"
                aria-label="Dismiss error message"
              >
                ×
              </button>
            )}
          </div>
        </div>
      )}

      {/* Tab Content */}
      <div role="tabpanel" id={`tabpanel-${activeTab}`}>
        {children}
      </div>
    </div>
  );
};

export default TabLayout;
