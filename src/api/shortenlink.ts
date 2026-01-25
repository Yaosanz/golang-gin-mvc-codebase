import api from './config.ts';

export interface CreateShortenLinkRequest {
  url: string;
}

export interface UpdateShortenLinkRequest {
  url: string;
}

export interface ShortenLinkResponse {
  id: string;
  code: string;
  url: string;
  created_at: string;
  updated_at: string;
  user_id: string;
}

export interface ShortenLinksListResponse {
  code: number;
  message: string;
  data: {
    contents: ShortenLinkResponse[];
    pagination?: {
      total_data: number;
      current_page: number;
      per_page: number;
      total_pages: number;
    };
  };
}

export interface ShortenLinkDetailResponse {
  code: number;
  message: string;
  data: ShortenLinkResponse;
}

// Get all shortened links for current user - Task 5: Redis Caching
export const getAll = async (): Promise<ShortenLinkResponse[]> => {
  const response = await api.get<ShortenLinksListResponse>('/shorten-links');
  return response.data.data?.contents || [];
};

// Create shortened link - Task 4: Database Transactions & Task 5: Cache Invalidation
export const create = async (data: CreateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  const response = await api.post<ShortenLinkDetailResponse>('/shorten-links', data);
  return response.data.data;
};

// Get shortened link by ID - Task 5: Redis Caching
export const getById = async (id: string): Promise<ShortenLinkResponse> => {
  const response = await api.get<ShortenLinkDetailResponse>(`/shorten-links/${id}`);
  return response.data.data;
};

// Update shortened link - Task 4: Database Transactions & Task 5: Cache Invalidation
export const update = async (id: string, data: UpdateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  const response = await api.patch<ShortenLinkDetailResponse>(`/shorten-links/${id}`, data);
  return response.data.data;
};

// Delete shortened link - Task 5: Cache Invalidation
export const remove = async (id: string): Promise<void> => {
  await api.delete(`/shorten-links/${id}`);
};

// Get shortened link by code (public redirect) - Task 5: Redis Caching (high-performance, 48-hour TTL)
export const getByCode = async (code: string): Promise<string> => {
  try {
    const response = await api.get(`/shortenlinks/${code}`);
    return response.data.data?.url || response.data.url;
  } catch (error) {
    throw error;
  }
};

// Redirect to original URL - Task 5: Cache Hit Performance
export const redirect = async (code: string): Promise<void> => {
  try {
    const response = await api.get(`/r/${code}`, {
      validateStatus: (status) => status < 400,
    });
    if (response.status === 301 || response.status === 302) {
      window.location.href = response.headers.location || response.data;
    }
  } catch (error) {
    console.error('Redirect failed:', error);
  }
};
