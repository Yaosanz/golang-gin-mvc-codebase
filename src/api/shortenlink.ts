import api from './config.ts';
import { getProfile } from './auth.ts';

export interface CreateShortenLinkRequest {
  url: string; // Backend expects 'url' field per Postman collection
}

export interface UpdateShortenLinkRequest {
  url?: string; // Backend expects 'url' field per Postman collection
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
  success?: boolean;
  code?: number;
  message: string;
  data:
    | ShortenLinkResponse[]
    | {
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
  code?: number;
  success?: boolean;
  message: string;
  data: ShortenLinkResponse;
}

// Transform Go struct field names to match frontend interface
const transformLink = (link: any): ShortenLinkResponse => ({
  id: link.id,
  code: link.short_code || link.code, // Backend returns 'short_code'
  url: link.original_url || link.url, // Backend returns 'original_url'
  created_at: link.created_at,
  updated_at: link.updated_at,
  user_id: link.user_id,
});

// Get all shortened links for current user - Task 5: Redis Caching
// Backend list route is registered with trailing slash
export const getAll = async (): Promise<ShortenLinkResponse[]> => {
  try {
    const response = await api.get<ShortenLinksListResponse>('/shorten-links/');
    const data = response.data.data;

    // Handle both array response and paginated response
    if (Array.isArray(data)) {
      return data.map(transformLink);
    } else {
      return (data?.contents || []).map(transformLink);
    }
  } catch (err: any) {
    // If permissions cache not warmed, fetch profile then retry once
    const msg: string | undefined = err?.response?.data?.message || err?.response?.data?.error;
    const status = err?.response?.status;
    if (status === 401 && msg && msg.toLowerCase().includes('permissions not cached')) {
      try {
        await getProfile();
        const retry = await api.get<ShortenLinksListResponse>('/shorten-links/');
        const data = retry.data.data;
        if (Array.isArray(data)) {
          return data.map(transformLink);
        }
        return (data?.contents || []).map(transformLink);
      } catch (retryErr) {
        throw retryErr;
      }
    }
    throw err;
  }
};

// Create shortened link - Task 4: Database Transactions & Task 5: Cache Invalidation
export const create = async (data: CreateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  if (!data.url || !data.url.trim()) {
    throw new Error('URL is required');
  }

  const payload = {
    url: data.url.trim(),
  };

  const response = await api.post<ShortenLinkDetailResponse>('/v1/shortenlinks', payload);
  return transformLink(response.data.data);
};

// Get shortened link by ID - Task 5: Redis Caching
export const getById = async (id: string): Promise<ShortenLinkResponse> => {
  const response = await api.get<ShortenLinkDetailResponse>(`/shorten-links/${id}`);
  return transformLink(response.data.data);
};

// Update shortened link - Task 4: Database Transactions & Task 5: Cache Invalidation
// Update shortened link - Task 4: Database Transactions & Task 5: Cache Invalidation
// Uses PUT /api/shortenlinks/:code with {url}
export const update = async (id: string, data: UpdateShortenLinkRequest): Promise<ShortenLinkResponse> => {
  const urlValue = data.url || data.original_url;
  if (!urlValue || !urlValue.trim()) {
    throw new Error('URL is required');
  }

  const payload = {
    url: urlValue.trim(),
  };

  const response = await api.put<ShortenLinkDetailResponse>(`/shortenlinks/${id}`, payload);
  return transformLink(response.data.data);
};

// Delete shortened link - Task 5: Cache Invalidation
export const remove = async (id: string): Promise<void> => {
  await api.delete(`/shorten-links/${id}`);
};

// Get shortened link by code (public redirect) - Task 5: Redis Caching (high-performance, 48-hour TTL)
export const getByCode = async (code: string): Promise<string> => {
  try {
    const response = await api.get(`/../shortenlinks/${code}`);
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
