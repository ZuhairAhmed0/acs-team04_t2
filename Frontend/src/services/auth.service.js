import { http, unwrap } from './http';

export function login(credentials) {
  return http.post('/auth/login', credentials).then((response) => {
    const payload = response.data.data ?? response.data;
    return { user: payload.user };
  });
}

export function registerUser(payload) {
  return unwrap(http.post('/auth/register/user', payload));
}

export function registerStation(payload) {
  return unwrap(http.post('/auth/register/station', payload));
}

export function fetchCurrentAccount() {
  return unwrap(http.get('/auth/me'));
}

export function logout() {
  return http.post('/auth/logout', {}).catch(() => null);
}
