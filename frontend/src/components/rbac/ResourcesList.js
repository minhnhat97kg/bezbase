import React, { useState, useEffect } from 'react';
import { rbacService } from '../../services/api';
import Table from '../common/Table';
import Icon from '../common/Icons';

const ResourcesList = () => {
  const [resources, setResources] = useState([]);
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

  const fetchResources = async (page = pagination.currentPage) => {
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

      const response = await rbacService.getResources(params);
      
      setResources(response.data.data);
      setPagination(prev => ({
        ...prev,
        currentPage: response.data.page,
        pageSize: response.data.page_size,
        total: response.data.total_items,
        totalPages: response.data.total_pages
      }));
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch resources');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchResources();
  }, []);

  useEffect(() => {
    fetchResources(1);
  }, [filters]);

  const handlePageChange = (page) => {
    fetchResources(page);
  };

  const handlePageSizeChange = (pageSize) => {
    setPagination(prev => ({ ...prev, pageSize, currentPage: 1 }));
    fetchResources(1);
  };

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  };


  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString();
  };

  const columns = [
    {
      key: 'id',
      header: 'ID',
      nowrap: true,
      render: (value) => (
        <code className="text-sm bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">
          {value}
        </code>
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
            placeholder="Search resources..."
            value={filters.search}
            onChange={(e) => handleFilterChange('search', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
          />
        </div>
      </div>

      <Table
        columns={columns}
        data={resources}
        loading={loading}
        emptyMessage="No resources found"
        emptyDescription="No resources match your search criteria"
        pagination={pagination}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
      />
    </div>
  );
};

export default ResourcesList;