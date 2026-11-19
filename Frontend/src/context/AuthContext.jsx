import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import * as authService from '../services/auth.service';
import { SESSION_EXPIRED_EVENT } from '../services/http';
import { clearStoredSession, getStoredSession, updateStoredUser } from '../utils/session';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => getStoredSession()?.user ?? null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;
    const session = getStoredSession();

    authService
      .fetchCurrentAccount()
      .then((account) => {
        if (active) {
          setUser(account);
          updateStoredUser(account);
        }
      })
      .catch(() => {
        if (active) {
          if (!session?.accessToken) {
            setUser(null);
          }
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    const handleExpired = () => {
      clearStoredSession();
      setUser(null);
    };
    window.addEventListener(SESSION_EXPIRED_EVENT, handleExpired);
    return () => {
      active = false;
      window.removeEventListener(SESSION_EXPIRED_EVENT, handleExpired);
    };
  }, []);

  const login = useCallback(async (credentials) => {
    const data = await authService.login(credentials);
    setUser(data.user);
    return data.user;
  }, []);

  const logout = useCallback(async () => {
    try {
      await authService.logout();
    } finally {
      clearStoredSession();
      setUser(null);
    }
  }, []);

  const value = useMemo(
    () => ({ user, isAuthenticated: Boolean(user), loading, login, logout }),
    [user, loading, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export const useAuth = () => useContext(AuthContext);
