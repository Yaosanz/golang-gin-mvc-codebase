import api from './config.ts';

export interface CreateUserRequest {
  name: string;
  username: string;
  email: string;
  phone?: string;
  is_active?: boolean;
  password: string;
}

export interface UpdateUserRequest {
  name?: string;
  username?: string;
  email?: string;
  phone?: string;
  is_active?: boolean;
  password?: string;
}

export interface UserResponse {
  id: number;
  name: string;
  username: string;
  email: string;
}

export const getAll = async (params?: { page?: number; limit?: number; search?: string }): Promise<{ data: UserResponse[]; total: number; page: number; limit: number }> => {
  const response = await api.get('/users', { params });
  return response.data;
};

export const create = async (data: CreateUserRequest): Promise<UserResponse> => {
  const response = await api.post('/users', data);
  return response.data.data || response.data;
};

export const getById = async (id: string): Promise<UserResponse> => {
  const response = await api.get(`/users/${id}`);
  return response.data.data || response.data;
};

export const update = async (id: string, data: UpdateUserRequest): Promise<UserResponse> => {
  const response = await api.patch(`/users/${id}`, data);
  return response.data.data || response.data;
};

export const remove = async (id: string): Promise<void> => {
  await api.delete(`/users/${id}`);
};
