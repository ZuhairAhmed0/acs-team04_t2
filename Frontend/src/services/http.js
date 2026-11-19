import axios from 'axios';

/** Configured via VITE_API_BASE_URL in Frontend/.env */
export const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/+$/, '');

export const SESSION_EXPIRED_EVENT = 'auth:session-expired';

export const http = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: 15000,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

let refreshPromise = null;

http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config;
    const status = error.response?.status;
    const url = original?.url ?? '';
    const isAuthMutation =
      url.includes('/auth/login') ||
      url.includes('/auth/refresh') ||
      url.includes('/auth/logout') ||
      url.includes('/auth/register');

    if (status === 401 && original && !original._retry && !isAuthMutation) {
      original._retry = true;
      try {
        refreshPromise ??= axios
          .post(`${API_BASE_URL}/api/v1/auth/refresh`, {}, { withCredentials: true })
          .finally(() => {
            refreshPromise = null;
          });

        await refreshPromise;
        return http(original);
      } catch {
        window.dispatchEvent(new Event(SESSION_EXPIRED_EVENT));
      }
    }
    return Promise.reject(error);
  },
);

/** Every JSON endpoint answers `{ success, data }` (or `{ success, message }` on some mutations). */
export const unwrap = (request) =>
  request.then((response) => (response.data.data !== undefined ? response.data.data : response.data));
