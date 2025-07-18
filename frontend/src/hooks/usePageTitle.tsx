import { useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const usePageTitle = (): string => {
  const location = useLocation();
  const { t } = useTranslation();
  
  const getPageTitle = (): string => {
    switch (location.pathname) {
      case '/dashboard':
        return t('pageTitle.dashboard');
      case '/profile':
        return t('pageTitle.profile');
      default:
        return t('pageTitle.default');
    }
  };

  return getPageTitle();
};

export default usePageTitle;