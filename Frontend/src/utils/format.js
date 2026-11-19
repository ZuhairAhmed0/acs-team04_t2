const numberFormat = new Intl.NumberFormat('en-US');
const dateFormat = new Intl.DateTimeFormat('ar-u-nu-latn', { day: 'numeric', month: 'long', year: 'numeric' });
const timeFormat = new Intl.DateTimeFormat('ar-u-nu-latn', { hour: '2-digit', minute: '2-digit' });

export const formatNumber = (value) => numberFormat.format(Number(value) || 0);
export const formatMoney = (value) => `${formatNumber(value)} ج.س`;
export const formatDate = (iso) => dateFormat.format(new Date(iso));
export const formatDateTime = (iso) => `${formatDate(iso)}، ${timeFormat.format(new Date(iso))}`;
