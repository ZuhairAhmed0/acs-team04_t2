import { SESSION_KEY } from '../config/constants';

/** Removes any legacy client-side session entry from localStorage (tokens now live in HttpOnly cookies). */
export const clearLegacySession = () => {
  try {
    localStorage.removeItem(SESSION_KEY);
  } catch {
    // Ignore storage access errors
  }
};
