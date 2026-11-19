import { http, unwrap } from './http';

export function fetchAdminDashboard() {
  return unwrap(http.get('/admin/dashboard')).then((data) => ({
    ...data,
    stations_count: data.stations_count ?? data.registered_stations_count,
  }));
}

/* ---- users ---- */

export function fetchUsers() {
  return unwrap(http.get('/admin/users'));
}

/** @param {'active'|'suspended'} status */
export function updateUserStatus(id, status) {
  return unwrap(http.patch(`/admin/users/${id}/status`, { status }));
}

/** @param {'user'|'admin'} role */
export function updateUserRole(id, role) {
  return unwrap(http.patch(`/admin/users/${id}/role`, { role }));
}

/* ---- stations ---- */

export function fetchPendingStations() {
  return unwrap(http.get('/admin/stations/pending')).then((list) =>
    (list ?? []).map((s) => ({ ...s, branches: s.branches ?? [] })),
  );
}

export function fetchStations() {
  return unwrap(http.get('/admin/stations')).then((list) =>
    (list ?? []).map((s) => ({ ...s, branches: s.branches ?? [] })),
  );
}

export function approveStation(id) {
  return unwrap(http.patch(`/admin/stations/${id}/approve`));
}

export function rejectStation(id, reason) {
  return unwrap(http.patch(`/admin/stations/${id}/reject`, { rejection_reason: reason }));
}

export function revokeStation(id) {
  return unwrap(http.patch(`/admin/stations/${id}/revoke`));
}

export function deleteStation(id) {
  return unwrap(http.delete(`/admin/stations/${id}`));
}

/* ---- cash withdrawals ---- */

export function fetchWithdrawals() {
  return unwrap(http.get('/admin/cash-withdrawals'));
}

export function approveWithdrawal(id) {
  return unwrap(http.patch(`/admin/cash-withdrawals/${id}/approve`));
}

export function rejectWithdrawal(id, reason) {
  return unwrap(http.patch(`/admin/cash-withdrawals/${id}/reject`, { rejection_reason: reason }));
}

export function completeWithdrawal(id, externalReference = '') {
  const body = externalReference ? { external_reference: externalReference } : {};
  return unwrap(http.patch(`/admin/cash-withdrawals/${id}/complete`, body));
}
