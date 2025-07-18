import React from 'react';
import { useTranslation } from 'react-i18next';
import Icon from './Icons';

interface TableColumn {
  key: string;
  header?: string;
  label?: string;
  sortable?: boolean;
  nowrap?: boolean;
  align?: string;
  className?: string;
  render?: (value: any, row: any) => React.ReactNode;
}

interface TablePagination {
  currentPage?: number;
  totalPages?: number;
  pageSize?: number;
  total?: number;
  pageSizeOptions?: number[];
}

interface TableProps {
  columns: TableColumn[];
  data: any[];
  loading?: boolean;
  emptyMessage?: string;
  emptyDescription?: string;
  onSort?: (field: string) => void;
  sortField?: string;
  sortOrder?: 'asc' | 'desc';
  pagination?: TablePagination;
  onPageChange?: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
  className?: string;
  rowClassName?: string;
  onRowClick?: (row: any) => void;
}

const Table: React.FC<TableProps> = ({
  columns,
  data,
  loading = false,
  emptyMessage,
  emptyDescription,
  onSort,
  sortField,
  sortOrder,
  pagination,
  onPageChange,
  onPageSizeChange,
  className = '',
  rowClassName = '',
  onRowClick,
}) => {
  const { t } = useTranslation();
  
  const defaultEmptyMessage = emptyMessage || t('table.noData');
  // Sort icon component
  const getSortIcon = (field: string) => {
    if (sortField !== field) {
      return (
        <Icon name="sortAsc" className="text-gray-400" />
      );
    }
    return sortOrder === 'asc' ? (
      <Icon name="chevronUp" className="text-blue-600" />
    ) : (
      <Icon name="chevronDown" className="text-blue-600" />
    );
  };

  // Pagination component
  const renderPagination = () => {
    if (!pagination) return null;

    const {
      currentPage = 1,
      totalPages = 1,
      pageSize = 10,
      total = 0,
      pageSizeOptions = [1, 5, 10, 25, 50]
    } = pagination;

    // Ensure totalPages is at least 1
    const safeTotalPages = Math.max(1, totalPages);

    const pages = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    let end = Math.min(safeTotalPages, start + maxVisible - 1);

    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1);
    }

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    return (
      <div className="flex items-center justify-between px-4 py-3 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 sm:px-6">
        <div className="flex items-center justify-between w-full">
          <div className="flex items-center">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              {t('table.showingResults', {
                from: ((currentPage - 1) * pageSize) + 1,
                to: Math.min(currentPage * pageSize, total),
                total: total
              })}
            </p>
            {onPageSizeChange && (
              <select
                value={pageSize}
                onChange={(e) => onPageSizeChange(Number(e.target.value))}
                className="ml-4 px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
              >
                {pageSizeOptions.map(size => (
                  <option key={size} value={size}>{size} {t('pagination.itemsPerPage')}</option>
                ))}
              </select>
            )}
          </div>

          <div className="flex items-center space-x-1">
            <button
              onClick={() => onPageChange && onPageChange(1)}
              disabled={currentPage === 1}
              className="px-2 py-1 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 disabled:opacity-50"
            >
              {t('pagination.first')}
            </button>
            <button
              onClick={() => onPageChange && onPageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="px-2 py-1 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 disabled:opacity-50"
            >
              {t('pagination.previous')}
            </button>

            {pages.map(page => (
              <button
                key={page}
                onClick={() => onPageChange && onPageChange(page)}
                className={`px-3 py-1 text-sm rounded-md ${page === currentPage
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  }`}
              >
                {page}
              </button>
            ))}

            <button
              onClick={() => onPageChange && onPageChange(currentPage + 1)}
              disabled={currentPage === safeTotalPages}
              className="px-2 py-1 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 disabled:opacity-50"
            >
              {t('pagination.next')}
            </button>
            <button
              onClick={() => onPageChange && onPageChange(safeTotalPages)}
              disabled={currentPage === safeTotalPages}
              className="px-2 py-1 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 disabled:opacity-50"
            >
              {t('pagination.last')}
            </button>
          </div>
        </div>
      </div>
    );
  };

  // Loading state
  if (loading) {
    return (
      <div className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className={`min-w-full divide-y divide-gray-200 dark:divide-gray-700 ${className}`}>
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                {columns.map((column, index) => (
                  <th
                    key={index}
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                  >
                    {column.header || column.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center">
                  <div className="flex justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        {pagination && renderPagination()}
      </div>
    );
  }

  // Empty state
  if (!data || data.length === 0) {
    return (
      <div className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className={`min-w-full divide-y divide-gray-200 dark:divide-gray-700 ${className}`}>
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                {columns.map((column, index) => (
                  <th
                    key={index}
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                  >
                    {column.header || column.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center">
                  <div className="text-gray-500 dark:text-gray-400">{defaultEmptyMessage}</div>
                  {emptyDescription && (
                    <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                      {emptyDescription}
                    </p>
                  )}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        {pagination && renderPagination()}
      </div>
    );
  }

  // Main table
  return (
    <div className="overflow-hidden">
      <div className="overflow-x-auto">
        <table className={`min-w-full divide-y divide-gray-200 dark:divide-gray-700 ${className}`}>
          <thead className="bg-gray-50 dark:bg-gray-700">
            <tr>
              {columns.map((column, index) => (
                <th
                  key={index}
                  className={`px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${column.sortable ? 'cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600' : ''
                    } ${column.align === 'right' ? 'text-right' : ''}`}
                  onClick={column.sortable && onSort ? () => onSort(column.key) : undefined}
                >
                  <div className={`flex items-center ${column.align === 'right' ? 'justify-end' : ''} space-x-1`}>
                    <span>{column.header || column.label}</span>
                    {column.sortable && onSort && getSortIcon(column.key)}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {data.map((row, rowIndex) => (
              <tr
                key={row.id || rowIndex}
                className={`hover:bg-gray-50 dark:hover:bg-gray-700 ${rowClassName} ${onRowClick ? 'cursor-pointer' : ''
                  }`}
                onClick={onRowClick ? () => onRowClick(row) : undefined}
              >
                {columns.map((column, colIndex) => (
                  <td
                    key={colIndex}
                    className={`px-6 py-4 ${column.nowrap ? 'whitespace-nowrap' : ''} ${column.align === 'right' ? 'text-right' : ''
                      } ${column.className || ''}`}
                  >
                    {column.render ? column.render(row[column.key], row) : row[column.key]}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {pagination && renderPagination()}
    </div>
  );
};

export default Table;
