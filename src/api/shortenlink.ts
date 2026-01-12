import api from './config.ts';

export interface CreateShortenLinkRequest {
  original_url: string;
}

export interface UpdateShortenLinkRequest {
  original_url: string;
}

export interface ShortenLinkResponse {
  id: string;
  short_code: string;
  original_url: string;
  created_at: string;
}

export const getAll = async (params?: { page?: number; limit?: number; search?: string }): Promise<{ data: ShortenLinkResponse[]; total: number; page: number; limit: number }> => {
  const response = await api.get('/shorten-links', { params });
  return response.data;
};

export const create = async (data: CreateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  const response = await api.post('/shorten-links', data);
  return response.data.data || response.data;
};

export const getById = async (id: string): Promise<ShortenLinkResponse> => {
  const response = await api.get(`/shorten-links/${id}`);
  return response.data.data || response.data;
};

export const update = async (id: string, data: UpdateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  const response = await api.patch(`/shorten-links/${id}`, data);
  return response.data.data || response.data;
};

export const remove = async (id: string): Promise<void> => {
  await api.delete(`/shorten-links/${id}`);
};

export const redirect = async (code: string): Promise<string> => {
  const response = await api.get(`/r/${code}`, { baseURL: api.defaults.baseURL?.replace('/api/v1', '') });
  return response.data;
};
