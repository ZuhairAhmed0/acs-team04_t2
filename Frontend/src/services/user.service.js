import { http, unwrap } from './http';

export function fetchUserDashboard() {
  return unwrap(http.get('/user/dashboard'));
}

export function fetchUserWallet() {
  return unwrap(http.get('/user/wallet'));
}

/** @param {string} [city] e.g. "Omdurman"; omit for all cities. */
export function fetchUserStations(city) {
  return unwrap(http.get('/user/stations', { params: city ? { city } : {} }));
}

/** @param {{ type?: string, status?: string, from_date?: string, to_date?: string }} [filters] */
export function fetchUserTransactions(filters = {}) {
  return unwrap(http.get('/user/transactions', { params: filters }));
}

/** Resolves with a CSV `Blob`. */
export function exportUserTransactions(filters = {}) {
  return http.get('/user/transactions/export', { params: filters, responseType: 'blob' }).then((response) => response.data);
}
