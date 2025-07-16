import React from 'react';
import Table from '../common/Table';
import Icon from '../common/Icons';

const RolesList = ({ 
  roles, 
  onEdit, 
  onDelete, 
  loading, 
  pagination, 
  onPageChange, 
  onPageSizeChange, 
  onSort, 
  sortField, 
  sortOrder 
}) => {
  const getStatusBadge = (isActive) => {
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
        isActive 
          ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100' 
          : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100'
      }`}>
        {isActive ? 'Active' : 'Inactive'}
      </span>
    );
  };

  const getSystemBadge = (isSystem) => {
    if (!isSystem) return null;
    return (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100">
        System
      </span>
    );
  };

  // Table configuration
  const columns = [
    {
      key: 'name',
      header: 'Role',
      nowrap: true,
      sortable: true,
      render: (value, row) => (
        <div>
          <div className="text-sm font-medium text-gray-900 dark:text-white">
            {row.display_name}
          </div>
          <div className="text-sm text-gray-500 dark:text-gray-400">
            {row.name}
          </div>
        </div>
      )
    },
    {
      key: 'description',
      header: 'Description',
      render: (value, row) => (
        <div className="text-sm text-gray-900 dark:text-white max-w-xs">
          {row.description || 'No description'}
        </div>
      )
    },
    {
      key: 'status',
      header: 'Status',
      nowrap: true,
      render: (value, row) => getStatusBadge(row.is_active)
    },
    {
      key: 'type',
      header: 'Type',
      nowrap: true,
      render: (value, row) => getSystemBadge(row.is_system)
    },
    {
      key: 'created_at',
      header: 'Created',
      nowrap: true,
      sortable: true,
      render: (value) => (
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {new Date(value).toLocaleDateString()}
        </span>
      )
    },
    {
      key: 'actions',
      header: 'Actions',
      align: 'right',
      nowrap: true,
      render: (value, row) => (
        <div className="flex justify-end space-x-2">
          <button
            onClick={() => onEdit(row)}
            className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
          >
            <Icon name="edit" />
          </button>
          {!row.is_system && (
            <button
              onClick={() => onDelete(row.name)}
              className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
            >
              <Icon name="delete" />
            </button>
          )}
        </div>
      )
    }
  ];

  return (
    <Table
      columns={columns}
      data={roles}
      loading={loading}
      emptyMessage="No roles found"
      emptyDescription="Create your first role to get started"
      pagination={pagination}
      onPageChange={onPageChange}
      onPageSizeChange={onPageSizeChange}
      onSort={onSort}
      sortField={sortField}
      sortOrder={sortOrder}
    />
  );
};

export default RolesList;