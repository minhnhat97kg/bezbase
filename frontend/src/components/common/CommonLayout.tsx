import React, { ReactNode } from 'react';
import Icon from './Icons';

interface CommonLayoutProps {
  title: string;
  subtitle?: string;
  children: ReactNode;
  headerActions?: ReactNode;
  message?: string;
  error?: string;
  onMessageDismiss?: () => void;
  onErrorDismiss?: () => void;
  className?: string;
}

/**
 * CommonLayout - A reusable layout component for standard pages without tabs
 * 
 * This component provides a consistent layout structure for pages that need:
 * - A header with title and subtitle
 * - Header actions (buttons, etc.)
 * - Success/error message display
 * - A clean content area
 * 
 * Perfect for list pages, forms, dashboards, and other single-content pages.
 */
const CommonLayout: React.FC<CommonLayoutProps> = ({
  title,
  subtitle,
  children,
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

            {/* Main Content */}
            <div>
              {children}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CommonLayout;