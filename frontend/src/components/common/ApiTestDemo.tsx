import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { userService, healthService, authService } from '../../services/api';

const ApiTestDemo: React.FC = () => {
  const { t } = useTranslation();
  const [testResults, setTestResults] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const addResult = (result: string) => {
    setTestResults(prev => [...prev, result]);
  };

  const clearResults = () => {
    setTestResults([]);
  };

  const testApiEndpoints = async () => {
    setIsLoading(true);
    clearResults();
    
    addResult('🧪 Testing API endpoints...');
    
    try {
      // Test health endpoint (no auth required)
      addResult('📡 Testing health endpoint...');
      const healthResponse = await healthService.check();
      addResult(`✅ Health: ${JSON.stringify(healthResponse.data)}`);
    } catch (error: any) {
      addResult(`❌ Health Error: ${error.message}`);
    }

    try {
      // Test auth login endpoint (versioned - should get 401 invalid credentials)
      addResult('🔐 Testing auth login endpoint (expects 401)...');
      await authService.login('test', 'test');
      addResult('✅ Auth: Success (unexpected!)');
    } catch (error: any) {
      if (error.response?.status === 401 && error.response?.data?.message === 'Invalid credentials') {
        addResult(`✅ Auth: Expected 401 'Invalid credentials' - versioned URL working!`);
      } else if (error.response?.status === 404) {
        addResult(`❌ Auth: 404 - URL construction failed!`);
      } else {
        addResult(`⚠️ Auth: Unexpected error ${error.response?.status}: ${error.response?.data?.message}`);
      }
    }

    try {
      // Test permissions endpoint (requires auth - will fail with 401)
      addResult('🔐 Testing permissions endpoint (expects 401)...');
      await userService.getPermissions();
      addResult('✅ Permissions: Success (unexpected!)');
    } catch (error: any) {
      if (error.response?.status === 401) {
        addResult(`✅ Permissions: Expected 401 - URL construction working!`);
      } else if (error.response?.status === 404) {
        addResult(`❌ Permissions: 404 - URL construction failed!`);
      } else {
        addResult(`⚠️ Permissions: Unexpected error ${error.response?.status}`);
      }
    }

    try {
      // Test profile endpoint (requires auth - will fail with 401)
      addResult('👤 Testing profile endpoint (expects 401)...');
      await userService.getProfile();
      addResult('✅ Profile: Success (unexpected!)');
    } catch (error: any) {
      if (error.response?.status === 401) {
        addResult(`✅ Profile: Expected 401 - URL construction working!`);
      } else if (error.response?.status === 404) {
        addResult(`❌ Profile: 404 - URL construction failed!`);
      } else {
        addResult(`⚠️ Profile: Unexpected error ${error.response?.status}`);
      }
    }

    setIsLoading(false);
    addResult('🏁 Testing complete!');
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 max-w-2xl">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
        API URL Construction Test
      </h3>
      
      <div className="space-y-4">
        <div className="text-sm text-gray-600 dark:text-gray-400">
          This test verifies that the frontend correctly constructs URLs:
          <ul className="mt-2 ml-4 list-disc text-xs">
            <li>All endpoints: versioned to v1 (e.g., /api/v1/auth/login, /api/v1/me/permissions)</li>
            <li>Health endpoint: unversioned (e.g., /api/health)</li>
          </ul>
        </div>
        
        <button
          onClick={testApiEndpoints}
          disabled={isLoading}
          className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? 'Testing...' : 'Test API Endpoints'}
        </button>
        
        <button
          onClick={clearResults}
          disabled={isLoading}
          className="ml-2 px-4 py-2 bg-gray-600 text-white text-sm font-medium rounded-md hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Clear Results
        </button>
        
        {testResults.length > 0 && (
          <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-md max-h-96 overflow-y-auto">
            <h4 className="text-sm font-semibold text-gray-800 dark:text-gray-200 mb-2">
              Test Results:
            </h4>
            <div className="space-y-1">
              {testResults.map((result, index) => (
                <div 
                  key={index}
                  className="text-xs font-mono text-gray-700 dark:text-gray-300"
                >
                  {result}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ApiTestDemo;