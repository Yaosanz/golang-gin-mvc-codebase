import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface User {
  user_id: string;
  username: string;
  role: string;
  role_id: string;
}

interface AuthContextType {
  token: string | null;
  user: User | null;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};

const decodeToken = (token: string): User | null => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return {
      user_id: payload.user_id,
      username: payload.username,
      role: payload.role,
      role_id: payload.role_id,
    };
  } catch {
    return null;
  }
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    if (storedToken) {
      setToken(storedToken);
      setUser(decodeToken(storedToken));
    }
  }, []);

  const login = (newToken: string) => {
    setToken(newToken);
    setUser(decodeToken(newToken));
    localStorage.setItem('token', newToken);
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
  };

  return <AuthContext.Provider value={{ token, user, login, logout }}>{children}</AuthContext.Provider>;
};
