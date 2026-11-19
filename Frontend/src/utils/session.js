import { SESSION_KEY } from '../config/constants';

export const getStoredSession = () => {
  try {
    const raw = localStorage.getItem(SESSION_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
};

export const getStoredAccessToken = () => getStoredSession()?.accessToken ?? null;
export const getStoredRefreshToken = () => getStoredSession()?.refreshToken ?? null;

export const saveSession = (sessionData) => {
  try {
    const current = getStoredSession() || {};
    localStorage.setItem(SESSION_KEY, JSON.stringify({ ...current, ...sessionData }));
  } catch {
    // Ignore storage errors
  }
};

export const setStoredAccessToken = (accessToken) => {
  try {
    const current = getStoredSession() || {};
    localStorage.setItem(SESSION_KEY, JSON.stringify({ ...current, accessToken }));
  } catch {
    // Ignore storage errors
  }
};

export const updateStoredUser = (user) => {
  try {
    const current = getStoredSession() || {};
    localStorage.setItem(SESSION_KEY, JSON.stringify({ ...current, user }));
  } catch {
    // Ignore storage errors
  }
};

export const clearStoredSession = () => {
  try {
    localStorage.removeItem(SESSION_KEY);
  } catch {
    // Ignore storage errors
  }
};

export const clearLegacySession = clearStoredSession;

