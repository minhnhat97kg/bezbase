import { useLocation } from 'react-router-dom';

const usePageTitle = () => {
  const location = useLocation();
  
  const getPageTitle = () => {
    switch (location.pathname) {
      case '/dashboard':
        return 'Dashboard';
      case '/profile':
        return 'Profile Settings';
      default:
        return 'Dashboard';
    }
  };

  return getPageTitle();
};

export default usePageTitle;