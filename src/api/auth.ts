import api from './config.ts';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  message: string;
  data: {
    access_token: string;
    token_type: string;
  };
}

export interface RegisterRequest {
  name: string;
  username: string;
  email: string;
  phone?: string;
  password: string;
}

export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/login', data);
  return response.data;
};

export const register = async (data: RegisterRequest): Promise<{ success: boolean; message: string }> => {
  const response = await api.post('/auth/register', data);
  return response.data;
};

export const logout = () => {
  localStorage.removeItem('token');
};
