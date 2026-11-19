import { generateIdempotencyKey } from '../utils/idempotency';
import { http, unwrap } from './http';

function postIdempotent(url, payload) {
  const { idempotency_key, ...body } = payload ?? {};
  const key = idempotency_key || generateIdempotencyKey();
  return unwrap(
    http.post(url, { ...body, idempotency_key: key }, { headers: { 'Idempotency-Key': key } }),
  );
}

export function topUpWallet(payload) {
  return postIdempotent('/transfers/topup', payload);
}

/** Allocation targets a station only; no branch is ever sent. */
export function allocateToStation(payload) {
  return postIdempotent('/transfers/allocate', payload);
}

/** Station-to-station transfers select stations only; no branches are sent. */
export function transferBetweenStations(payload) {
  return postIdempotent('/transfers/station-to-station', payload);
}

export function purchaseFuel(payload) {
  return postIdempotent('/transfers/purchase', payload);
}

export function fetchTransfers() {
  return unwrap(http.get('/transfers'));
}

export function fetchTransferById(id) {
  return unwrap(http.get(`/transfers/${id}`));
}
