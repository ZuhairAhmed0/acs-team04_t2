/** Amounts are integer SDG units. Returns the number, or null when invalid. */
export function parseAmount(value) {
  const amount = Number(value);
  return Number.isInteger(amount) && amount > 0 ? amount : null;
}

/** Parses an integer >= 0 (empty input returns 0). Returns null when negative or non-integer. */
export function parseNonNegativeAmount(value) {
  if (value === '' || value === null || value === undefined) return 0;
  const amount = Number(value);
  return Number.isInteger(amount) && amount >= 0 ? amount : null;
}
