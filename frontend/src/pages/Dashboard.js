import React, { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { healthService } from '../services/api';

const Dashboard = () => {
  const { user } = useAuth();
  const [healthStatus, setHealthStatus] = useState(null);

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const response = await healthService.check();
        setHealthStatus(response.data);
      } catch (error) {
        console.error('Health check failed:', error);
        setHealthStatus({ status: 'error', message: 'Health check failed' });
      }
    };

    checkHealth();
  }, []);

  return (
    <div className="max-w-6xl mx-auto">
      <div className="mb-6">
        <p className="text-gray-600">Welcome to your BezBase dashboard</p>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Welcome Card */}
        <div className="card lg:col-span-2">
          <h3 className="text-xl font-semibold text-gray-900 mb-4">
            Welcome back, {user?.first_name} {user?.last_name}!
          </h3>
          <div className="space-y-2 text-gray-600">
            <p><span className="font-medium">Email:</span> {user?.email}</p>
            <p><span className="font-medium">Member since:</span> {new Date(user?.created_at).toLocaleDateString()}</p>
          </div>
        </div>

        {/* System Status Card */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">System Status</h3>
          {healthStatus ? (
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <span className="text-sm font-medium text-gray-700">Status:</span>
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                  healthStatus.status === 'healthy' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-red-100 text-red-800'
                }`}>
                  {healthStatus.status}
                </span>
              </div>
              <p className="text-sm text-gray-600">{healthStatus.message}</p>
            </div>
          ) : (
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600"></div>
              <span className="text-sm text-gray-600">Checking system status...</span>
            </div>
          )}
        </div>

        {/* Quick Actions Card */}
        <div className="card lg:col-span-3">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <p className="text-gray-600 mb-4">
            This is your main dashboard. You can add more features here such as:
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900 mb-1">User Statistics</h4>
              <p className="text-sm text-gray-600">View your usage metrics</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900 mb-1">Recent Activities</h4>
              <p className="text-sm text-gray-600">Check your latest actions</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900 mb-1">Quick Settings</h4>
              <p className="text-sm text-gray-600">Manage your preferences</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900 mb-1">Notifications</h4>
              <p className="text-sm text-gray-600">Stay updated with alerts</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

