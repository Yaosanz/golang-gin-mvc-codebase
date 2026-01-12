import api from './config.ts';

export interface CreateShortenLinkRequest {
  original_url: string;
}

export interface UpdateShortenLinkRequest {
  original_url: string;
}

export interface ShortenLinkResponse {
  ID: string;
  ShortCode: string;
  OriginalURL: string;
  CreatedAt: string;
}

export const getAll = async (): Promise<ShortenLinkResponse[]> => {
  const response = await api.get('/shorten-links');
  return response.data.data || [];
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
