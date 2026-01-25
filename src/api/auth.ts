import api from './config.ts';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
    user: {
      id: string;
      username: string;
      email: string;
      name: string;
      phone: string;
      roles: string[];
    };
  };
}

export interface ProfileResponse {
  code: number;
  message: string;
  data: {
    id: string;
    username: string;
    email: string;
    name: string;
    phone: string;
    created_at: string;
    updated_at: string;
    roles: string[];
  };
}

export interface RegisterRequest {
  name: string;
  username: string;
  email: string;
  phone?: string;
  password: string;
}

export interface UpdateProfileRequest {
  name?: string;
  email?: string;
  phone?: string;
}

// Login - Task 1: JWT Multi-Role Auth
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/login', data);
  return response.data;
};

// Register - Task 1: JWT Multi-Role Auth
export const register = async (data: RegisterRequest): Promise<{ code: number; message: string }> => {
  const response = await api.post('/auth/register', data);
  return response.data;
};

// Get Profile - Task 2: RBAC Token Payload (returns roles array)
export const getProfile = async (): Promise<ProfileResponse> => {
  const response = await api.get<ProfileResponse>('/auth/profile');
  return response.data;
};

// Update Profile - Self-service profile update
export const updateProfile = async (data: UpdateProfileRequest): Promise<ProfileResponse> => {
  const response = await api.put<ProfileResponse>('/users', data);
  return response.data;
};

// Logout
export const logout = () => {
  localStorage.removeItem('token');
};
