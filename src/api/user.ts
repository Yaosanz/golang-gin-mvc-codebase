import api from './config.ts';

export interface CreateUserRequest {
  name: string;
  username: string;
  email: string;
  phone?: string;
  password: string;
  role_id?: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
}

export interface UserResponse {
  id: string;
  name: string;
  username: string;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
  roles?: Array<{ id: string; name: string }>;
}

export interface UsersListResponse {
  code: number;
  message: string;
  data: {
    contents: UserResponse[];
    pagination: {
      total_data: number;
      current_page: number;
      per_page: number;
      total_pages: number;
    };
  };
}

export interface UserDetailResponse {
  code: number;
  message: string;
  data: UserResponse;
}

// Get all users (Admin only) - Task 2: RBAC
export const getAll = async (params?: { page?: number; limit?: number; search?: string }): Promise<{ data: UserResponse[]; total: number; page: number; limit: number }> => {
  const response = await api.get<UsersListResponse>('/users', { params });
  const resData = response.data.data;
  return {
    data: resData.contents || [],
    total: resData.pagination?.total_data || 0,
    page: resData.pagination?.current_page || 1,
    limit: resData.pagination?.per_page || 10,
  };
};

// Create user (Admin only) - Task 4: Database Transactions
export const create = async (data: CreateUserRequest): Promise<UserResponse> => {
  const response = await api.post<UserDetailResponse>('/users', data);
  return response.data.data || response.data;
};

// Get user by ID (Admin only) - Task 3: Middleware Token Identification
export const getById = async (id: string): Promise<UserResponse> => {
  const response = await api.get<UserDetailResponse>(`/users/${id}`);
  return response.data.data || response.data;
};

// Update user (Admin only) - Task 4: Database Transactions
export const update = async (id: string, data: UpdateUserRequest): Promise<UserResponse> => {
  const response = await api.patch<UserDetailResponse>(`/users/${id}`, data);
  return response.data.data || response.data;
};

// Delete user (Admin only) - Task 3: Middleware Token Identification
export const remove = async (id: string): Promise<void> => {
  await api.delete(`/users/${id}`);
};
