// ---------------------------------------------------------------------------
// Indian-locale formatting helpers (mirrors Go shared-kernel)
// ---------------------------------------------------------------------------

/**
 * Format a number as INR with Indian grouping (12,34,567.89).
 * Prefix: ₹
 */
export function formatINR(amount: number): string {
  const negative = amount < 0;
  const abs = Math.abs(amount);
  const [intPart, decPart] = abs.toFixed(2).split('.');

  // Indian grouping: last 3 digits, then groups of 2
  let formatted: string;
  if (intPart!.length <= 3) {
    formatted = intPart!;
  } else {
    const last3 = intPart!.slice(-3);
    const rest = intPart!.slice(0, -3);
    const groups: string[] = [];
    for (let i = rest.length; i > 0; i -= 2) {
      groups.unshift(rest.slice(Math.max(0, i - 2), i));
    }
    formatted = groups.join(',') + ',' + last3;
  }

  return (negative ? '-' : '') + '₹' + formatted + '.' + decPart;
}

/**
 * Format a Date as DD/MM/YYYY.
 */
export function formatDate(date: Date): string {
  const dd = String(date.getDate()).padStart(2, '0');
  const mm = String(date.getMonth() + 1).padStart(2, '0');
  const yyyy = String(date.getFullYear());
  return `${dd}/${mm}/${yyyy}`;
}

/**
 * Parse a DD/MM/YYYY string into a Date.
 */
export function parseDate(str: string): Date {
  const parts = str.split('/');
  if (parts.length !== 3) {
    throw new Error(`Invalid date format: ${str}. Expected DD/MM/YYYY`);
  }
  const [dd, mm, yyyy] = parts as [string, string, string];
  const date = new Date(Number(yyyy), Number(mm) - 1, Number(dd));
  if (isNaN(date.getTime())) {
    throw new Error(`Invalid date: ${str}`);
  }
  return date;
}

/**
 * Validate GSTIN format (15-character alphanumeric with checksum pattern).
 * Pattern: 2-digit state code + 10-char PAN + 1 entity code + Z + 1 check digit
 */
export function validateGSTIN(gstin: string): boolean {
  const pattern = /^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$/;
  return pattern.test(gstin);
}

/**
 * Validate PAN format (AAAAA9999A).
 */
export function validatePAN(pan: string): boolean {
  const pattern = /^[A-Z]{5}[0-9]{4}[A-Z]{1}$/;
  return pattern.test(pan);
}

/**
 * Validate TAN format (AAAA99999A).
 */
export function validateTAN(tan: string): boolean {
  const pattern = /^[A-Z]{4}[0-9]{5}[A-Z]{1}$/;
  return pattern.test(tan);
}

/**
 * Format a number in Indian compact notation:
 * - >= 1 crore (10M)  -> "1.5Cr"
 * - >= 1 lakh (100K)  -> "45L"
 * - otherwise          -> plain number
 */
export function formatCompact(num: number): string {
  const abs = Math.abs(num);
  const sign = num < 0 ? '-' : '';

  if (abs >= 1_00_00_000) {
    // crore
    const crores = abs / 1_00_00_000;
    const display = crores % 1 === 0 ? crores.toFixed(0) : crores.toFixed(1);
    return sign + display + 'Cr';
  }
  if (abs >= 1_00_000) {
    // lakh
    const lakhs = abs / 1_00_000;
    const display = lakhs % 1 === 0 ? lakhs.toFixed(0) : lakhs.toFixed(1);
    return sign + display + 'L';
  }
  return sign + abs.toLocaleString('en-IN');
}
