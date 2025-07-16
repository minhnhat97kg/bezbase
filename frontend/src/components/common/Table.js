import React from 'react';

const Table = ({ columns, data, loading, renderActions }) => {
  if (loading) {
    return (
      <div className="flex justify-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-500 dark:text-gray-400">No data found</div>
      </div>
    );
  }

  return (
    <div className="overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-700">
            <tr>
              {columns.map((col) => (
                <th
                  key={col.key}
                  className={col.thClassName || "px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"}
                  style={col.style}
                >
                  {col.title}
                </th>
              ))}
              {renderActions && (
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
              )}
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {data.map((row, rowIndex) => (
              <tr key={row.id || rowIndex} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className={col.tdClassName || "px-6 py-4 whitespace-nowrap"}
                  >
                    {col.render ? col.render(row) : row[col.key]}
                  </td>
                ))}
                {renderActions && (
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    {renderActions(row)}
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Table;
