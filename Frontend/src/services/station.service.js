import { generateIdempotencyKey } from '../utils/idempotency';
import { http, unwrap } from './http';

export function fetchStationDashboard() {
  return unwrap(http.get('/station/dashboard'));
}

export function fetchStationWallet() {
  return unwrap(http.get('/station/wallet'));
}

export function fetchStationTransactions(filters = {}) {
  return unwrap(http.get('/station/transactions', { params: filters }));
}

/** Resolves with a CSV `Blob`. */
export function exportStationTransactions(filters = {}) {
  return http.get('/station/transactions/export', { params: filters, responseType: 'blob' }).then((response) => response.data);
}

export function fetchBranches() {
  return unwrap(http.get('/station/branches'));
}

/** @param {{ city: string, address: string }} payload */
export function addBranch(payload) {
  return unwrap(http.post('/station/branches', payload));
}

/** @param {'full'|'low'|'empty'} fuelStatus */
export function updateBranchFuelStatus(branchId, fuelStatus) {
  return unwrap(http.patch(`/station/branches/${branchId}/fuel-status`, { fuel_status: fuelStatus }));
}

/** @param {number} congestionPercentage 0..100 */
export function updateBranchCongestion(branchId, congestionPercentage) {
  return unwrap(http.patch(`/station/branches/${branchId}/congestion`, { congestion_percentage: congestionPercentage }));
}

export function submitCashWithdrawal(payload) {
  const { idempotency_key, ...body } = payload ?? {};
  const key = idempotency_key || generateIdempotencyKey();
  return unwrap(
    http.post('/station/cash-withdrawals', { ...body, idempotency_key: key }, { headers: { 'Idempotency-Key': key } }),
  );
}

export function fetchStationWithdrawals() {
  return unwrap(http.get('/station/cash-withdrawals'));
}
