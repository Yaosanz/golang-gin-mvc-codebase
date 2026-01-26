import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';

const BASE_URL = process.env.REACT_APP_API_URL || '/api';

const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
});

// Request interceptor to add JWT token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      } as any;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  },
);

// Response interceptor to handle 401
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      const data: any = error.response?.data;
      const msg: string | undefined = data?.message || data?.error;
      const url = (error.config as any)?.url || '';
      const recentTsStr = localStorage.getItem('recent_login_ts') || '0';
      const recentTs = Number(recentTsStr);
      const withinWarmup = Date.now() - recentTs < 10000; // 10s warm-up window

      // During warm-up window after login, never auto-redirect on 401
      if (withinWarmup) {
        console.warn('[401 skipped - warmup]', url, msg);
        return Promise.reject(error);
      }

      // Skip redirect when backend is telling us to warm permissions
      if (msg && msg.toLowerCase().includes('permissions not cached')) {
        console.warn('[401 skipped - permissions cache]', url, msg);
        return Promise.reject(error);
      }

      // Only redirect for truly invalid/expired tokens
      const lower = (msg || '').toLowerCase();
      const shouldLogout = lower.includes('invalid token') || lower.includes('token expired') || lower.includes('unauthorized') || lower.includes('session expired') || lower.includes('please login again');
      if (shouldLogout) {
        localStorage.removeItem('token');
        localStorage.removeItem('recent_login_ts');
        window.location.href = '/login';
        return Promise.reject(error);
      }

      // Default: do not auto-logout; surface error to caller
      console.warn('[401 surfaced]', url, msg);
      return Promise.reject(error);
    }
    return Promise.reject(error);
  },
);

export default api;
