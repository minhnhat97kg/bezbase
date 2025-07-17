import React, { useState, useEffect } from 'react';
import { rbacService } from '../../services/api';
import Table from '../common/Table';
import Icon from '../common/Icons';

const ActionsList = () => {
  const [actions, setActions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({
    currentPage: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0,
    pageSizeOptions: [5, 10, 25, 50]
  });
  const [filters, setFilters] = useState({
    search: ''
  });

  const fetchActions = async (page = pagination.currentPage) => {
    try {
      setLoading(true);
      const params = {
        page,
        page_size: pagination.pageSize,
        ...filters
      };

      // Remove empty filters
      Object.keys(params).forEach(key => {
        if (params[key] === '' || params[key] === null || params[key] === undefined) {
          delete params[key];
        }
      });

      const response = await rbacService.getActions(params);
      
      setActions(response.data.data);
      setPagination(prev => ({
        ...prev,
        currentPage: response.data.page,
        pageSize: response.data.page_size,
        total: response.data.total_items,
        totalPages: response.data.total_pages
      }));
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch actions');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchActions();
  }, []);

  useEffect(() => {
    fetchActions(1);
  }, [filters]);

  const handlePageChange = (page) => {
    fetchActions(page);
  };

  const handlePageSizeChange = (pageSize) => {
    setPagination(prev => ({ ...prev, pageSize, currentPage: 1 }));
    fetchActions(1);
  };

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  };


  const getActionIcon = (actionId) => {
    const iconMap = {
      'create': 'plus',
      'read': 'eye',
      'update': 'edit',
      'delete': 'delete',
      '*': 'star'
    };

    return iconMap[actionId] || 'question';
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString();
  };

  const columns = [
    {
      key: 'id',
      header: 'ID',
      nowrap: true,
      render: (value, row) => (
        <div className="flex items-center space-x-2">
          <Icon name={getActionIcon(value)} className="w-4 h-4 text-gray-500" />
          <code className="text-sm bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">
            {value}
          </code>
        </div>
      )
    },
    {
      key: 'name',
      header: 'Name',
      sortable: true,
      render: (value, row) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">
            {value}
          </div>
          <div className="text-sm text-gray-500 dark:text-gray-400">
            {row.description}
          </div>
        </div>
      )
    },
    {
      key: 'created_at',
      header: 'Created At',
      sortable: true,
      render: (value) => formatDate(value)
    }
  ];

  return (
    <div>
      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Search Controls */}
      <div className="mb-6">
        <div className="flex-1 max-w-md">
          <input
            type="text"
            placeholder="Search actions..."
            value={filters.search}
            onChange={(e) => handleFilterChange('search', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
          />
        </div>
      </div>

      <Table
        columns={columns}
        data={actions}
        loading={loading}
        emptyMessage="No actions found"
        emptyDescription="No actions match your search criteria"
        pagination={pagination}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
      />
    </div>
  );
};

export default ActionsList;