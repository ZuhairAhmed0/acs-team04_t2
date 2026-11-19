export const ROLES = Object.freeze({ USER: 'user', STATION: 'station', ADMIN: 'admin' });

export const ROLE_HOME = {
  [ROLES.USER]: '/user',
  [ROLES.STATION]: '/station',
  [ROLES.ADMIN]: '/admin',
};

export const ROLE_LABELS = {
  [ROLES.USER]: 'حساب مستخدم',
  [ROLES.STATION]: 'حساب محطة',
  [ROLES.ADMIN]: 'مدير النظام',
};

export const SESSION_KEY = 'fuel_session';

export const BANKS = [
  { value: 'bank_of_khartoum', label: 'بنك الخرطوم' },
  { value: 'omdurman_national_bank', label: 'بنك أم درمان الوطني' },
  { value: 'faisal_islamic_bank', label: 'بنك فيصل الإسلامي' },
];
const BANK_ALIASES = {
  'Bank of Khartoum': 'بنك الخرطوم',
  'Omdurman National Bank': 'بنك أم درمان الوطني',
  'Faisal Islamic Bank': 'بنك فيصل الإسلامي',
};
export const bankLabel = (value) => BANKS.find((bank) => bank.value === value)?.label ?? BANK_ALIASES[value] ?? value;

export const CITIES = [
  { value: 'الخرطوم', label: 'الخرطوم' },
  { value: 'أم درمان', label: 'أم درمان' },
  { value: 'بحري', label: 'بحري' },
  { value: 'بورتسودان', label: 'بورتسودان' },
  { value: 'ود مدني', label: 'ود مدني' },
  { value: 'الأبيض', label: 'الأبيض' },
  { value: 'كسلا', label: 'كسلا' },
];
const CITY_ALIASES = {
  Khartoum: 'الخرطوم',
  Omdurman: 'أم درمان',
  'Khartoum North': 'بحري',
  Bahri: 'بحري',
  'Port Sudan': 'بورتسودان',
  'Wad Madani': 'ود مدني',
  'El Obeid': 'الأبيض',
  Kassala: 'كسلا',
};
export const cityLabel = (value) => CITIES.find((city) => city.value === value)?.label ?? CITY_ALIASES[value] ?? value;

export const FUEL_STATUS = {
  full: { label: 'ممتلئة', color: 'success' },
  low: { label: 'ممتلئة قليلاً', color: 'warning' },
  empty: { label: 'فارغة', color: 'error' },
};

export const TRANSACTION_TYPES = {
  topup: 'شحن المحفظة',
  allocation: 'تخصيص لمحطة',
  fuel_purchase: 'شراء وقود',
  station_transfer: 'تحويل بين المحطات',
  withdrawal: 'سحب كاش',
};

const STATUS_META = {
  completed: { label: 'مكتملة', color: 'success' },
  pending: { label: 'قيد المراجعة', color: 'warning' },
  failed: { label: 'فاشلة', color: 'error' },
  active: { label: 'نشط', color: 'success' },
  suspended: { label: 'موقوف', color: 'error' },
  approved: { label: 'تمت الموافقة', color: 'info' },
  rejected: { label: 'مرفوض', color: 'error' },
  revoked: { label: 'أُلغيت الموافقة', color: 'error' },
};
export const statusMeta = (status) => STATUS_META[status] ?? { label: status, color: 'default' };
export const TRANSACTION_STATUSES = ['completed', 'pending', 'failed'];
export const WITHDRAWAL_STATUSES = ['pending', 'approved', 'completed', 'rejected'];
