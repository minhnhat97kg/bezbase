import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { authService, userService } from "../services/api";
import { rbacService } from "../services/rbacService";

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

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export { AuthContext };

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [permissions, setPermissions] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const initializeAuth = async () => {
      const token = localStorage.getItem("token");
      const savedUser = localStorage.getItem("user");

      if (token && savedUser) {
        try {
          const parsedUser = JSON.parse(savedUser);
          // Validate the token by fetching user profile
          try {
            const response = await userService.getProfile();
            setUser(response.data);
            localStorage.setItem("user", JSON.stringify(response.data));
            await fetchPermissions();
          } catch (error) {
            // Token is invalid, clear storage and redirect to login
            console.error("Token validation failed:", error);
            localStorage.removeItem("token");
            localStorage.removeItem("user");
            setUser(null);
            setPermissions([]);
          }
        } catch (error) {
          console.error("Error parsing saved user:", error);
          localStorage.removeItem("token");
          localStorage.removeItem("user");
        }
      }
      setLoading(false);
    };

    initializeAuth();
  }, []);

  // Fetch permissions if token exists
  const fetchPermissions = async () => {
    try {
      const res = await rbacService.getMyPermissions();
      setPermissions(res.data.permissions || []);
    } catch (err) {
      setPermissions([]);
    }
  };

  const login = async (username: string, password: string): Promise<{ success: boolean; error?: string }> => {
    try {
      const response = await authService.login(username, password);
      const { token, user } = response.data;

      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(user));
      setUser(user);
      fetchPermissions();

      return { success: true };
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.message || "Login failed",
      };
    }
  };

  const register = async (userData: any): Promise<{ success: boolean; error?: string }> => {
    try {
      const response = await authService.register(userData);
      const { token, user } = response.data;

      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(user));
      setUser(user);

      return { success: true };
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.message || "Registration failed",
      };
    }
  };

  const logout = (): void => {
    authService.logout();
    setUser(null);
  };

  const updateUser = (updatedUser: User): void => {
    setUser(updatedUser);
    localStorage.setItem("user", JSON.stringify(updatedUser));
  };

  const userProfile = async (): Promise<{ success: boolean; user?: User; error?: string }> => {
    try {
      const response = await userService.getProfile();
      const userData = response.data;
      setUser(userData);
      localStorage.setItem("user", JSON.stringify(userData));
      return { success: true, user: userData };
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.message || "Failed to fetch user profile",
      };
    }
  };

  const value = {
    user,
    permissions,
    login,
    register,
    logout,
    updateUser,
    userProfile,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
