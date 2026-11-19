export const getErrorMessage = (error, fallback = 'حدث خطأ غير متوقع، حاول مرة أخرى') =>
  error?.response?.data?.message ?? (error?.code === 'ERR_NETWORK' ? 'تعذّر الاتصال بالخادم' : fallback);
