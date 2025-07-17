import React from 'react';

// Centralized icon configuration using Font Awesome free version
const ICON_CLASSES = {
  // Navigation
  dashboard: 'fas fa-tachometer-alt',
  user: 'fas fa-user',
  users: 'fas fa-users',
  shield: 'fas fa-shield-alt',
  key: 'fas fa-key',
  menu: 'fas fa-bars',
  
  // Actions
  edit: 'fas fa-edit',
  delete: 'fas fa-trash',
  trash: 'fas fa-trash',
  close: 'fas fa-times',
  search: 'fas fa-search',
  exit: 'fas fa-sign-out-alt',
  plus: 'fas fa-plus',
  eye: 'fas fa-eye',
  star: 'fas fa-star',
  
  // UI Elements
  chevronUp: 'fas fa-chevron-up',
  chevronDown: 'fas fa-chevron-down',
  sortAsc: 'fas fa-sort-amount-up',
  
  // Status & Feedback
  error: 'fas fa-times-circle',
  loading: 'fas fa-spinner',
  success: 'fas fa-check-circle',
  check: 'fas fa-check',
  warning: 'fas fa-exclamation-triangle',
  exclamation: 'fas fa-exclamation-triangle',
  
  // Theme
  moon: 'fas fa-moon',
  sun: 'fas fa-sun',
  
  // Misc
  rocket: 'fas fa-rocket',
  alarm: 'fas fa-bell',
  question: 'fas fa-question-circle',
  settings: 'fas fa-cog',
};

// Icon component
const Icon = ({ 
  name, 
  className = '', 
  size = '', 
  color = '', 
  ...props 
}) => {
  const iconClass = ICON_CLASSES[name];
  
  if (!iconClass) {
    console.warn(`Icon "${name}" not found in ICON_CLASSES`);
    return null;
  }
  
  const combinedClassName = `${iconClass} ${className} ${size} ${color}`.trim();
  
  return <i className={combinedClassName} {...props} />;
};

// Export both the Icon component and the classes for direct usage
export { ICON_CLASSES };
export default Icon;