import { useContext } from 'react';
import { AuthContext } from '../context/AuthContext';

interface User {
  id: string;
  username: string;
  email: string;
  [key: string]: any;
}

interface AuthContextType {
  user: User | null;
  permissions: string[];
  login: (username: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (userData: any) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  updateUser: (updatedUser: User) => void;
  userProfile: () => Promise<{ success: boolean; user?: User; error?: string }>;
  loading: boolean;
}

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

