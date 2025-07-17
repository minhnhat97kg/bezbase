import React, { ReactNode } from 'react';
import Icon from './Icons';

interface Tab {
  id: string;
  name: string;
  icon?: string;
}

interface TabLayoutProps {
  title: string;
  subtitle?: string;
  children: ReactNode;
  tabs?: Tab[];
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
  tabs = [],
  activeTab,
  onTabChange,
  showTabs = false,
  headerActions = null,
  message = '',
  error = '',
  onMessageDismiss,
  onErrorDismiss,
  className = '',
}) => {
  return (
    <div className={`min-h-screen bg-gray-50 py-4 ${className}`}>
      <div className="mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 mb-4">
          <div className="px-4 py-3 border-b border-gray-200">
            <div className="flex justify-between items-start">
              <div>
                <h1 className="text-2xl font-semibold text-gray-900">{title}</h1>
                {subtitle && <p className="text-gray-600 mt-1">{subtitle}</p>}
              </div>
              {headerActions && (
                <div className="flex items-center space-x-2">
                  {headerActions}
                </div>
              )}
            </div>
          </div>

          {/* Tabs Navigation */}
          {showTabs && tabs.length > 0 && (
            <div className="flex border-b border-gray-200">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => onTabChange && onTabChange(tab.id)}
                  className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors flex items-center ${activeTab === tab.id
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    }`}
                  role="tab"
                  aria-selected={activeTab === tab.id}
                  aria-controls={`tabpanel-${tab.id}`}
                >
                  {tab.icon && (
                    <Icon name={tab.icon} className="mr-2" />
                  )}
                  {tab.name}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Content Area */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="p-4">
            {/* Success Message */}
            {message && (
              <div className="mb-4 bg-green-50 border border-green-200 rounded-lg p-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Icon name="check" className="text-green-500 mr-2" />
                    <span className="text-sm text-green-700">{message}</span>
                  </div>
                  {onMessageDismiss && (
                    <button
                      onClick={onMessageDismiss}
                      className="text-green-400 hover:text-green-600 transition-colors"
                      aria-label="Dismiss success message"
                    >
                      <Icon name="close" className="w-4 h-4" />
                    </button>
                  )}
                </div>
              </div>
            )}

            {/* Error Message */}
            {error && (
              <div className="mb-4 bg-red-50 border border-red-200 rounded-lg p-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Icon name="error" className="text-red-500 mr-2" />
                    <span className="text-sm text-red-700">{error}</span>
                  </div>
                  {onErrorDismiss && (
                    <button
                      onClick={onErrorDismiss}
                      className="text-red-400 hover:text-red-600 transition-colors"
                      aria-label="Dismiss error message"
                    >
                      <Icon name="close" className="w-4 h-4" />
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
        </div>
      </div>
    </div>
  );
};

export default TabLayout;
