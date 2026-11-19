import { useState } from 'react';
import { Alert, Box, Button, MenuItem, Stack, TextField, ToggleButton, ToggleButtonGroup, Typography } from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import AuthLayout from '../../components/layout/AuthLayout';
import { BANKS } from '../../config/constants';
import useFormState from '../../hooks/useFormState';
import useSubmit from '../../hooks/useSubmit';
import { registerStation, registerUser } from '../../services/auth.service';

const ACCOUNT_TYPES = {
  user: { label: 'مستخدم', required: ['name', 'email', 'phone', 'password'] },
  station: {
    label: 'محطة وقود',
    required: ['name', 'location', 'bank_name', 'bank_account_number', 'bank_account_name', 'email', 'password'],
  },
};

const INITIAL_VALUES = {
  name: '', email: '', phone: '', password: '',
  location: '', bank_name: '', bank_account_number: '', bank_account_name: '',
};

const MIN_PASSWORD_LENGTH = 6;

const pickFields = (source, fields) => Object.fromEntries(fields.map((field) => [field, source[field]]));

export default function RegisterPage() {
  const navigate = useNavigate();
  const [accountType, setAccountType] = useState('user');
  const [stationSubmitted, setStationSubmitted] = useState(false);
  const { values, errors, setErrors, handleChange, validateRequired } = useFormState(INITIAL_VALUES);
  const isStation = accountType === 'station';

  const { submit, submitting } = useSubmit(async () => {
    if (isStation) {
      await registerStation(pickFields(values, ACCOUNT_TYPES.station.required));
      setStationSubmitted(true);
    } else {
      await registerUser(pickFields(values, ACCOUNT_TYPES.user.required));
      navigate('/login', { replace: true, state: { registered: true } });
    }
  });

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!validateRequired(ACCOUNT_TYPES[accountType].required)) return;
    if (!isStation && !/^\d{1,10}$/.test(values.phone.trim())) {
      setErrors({ phone: 'رقم الهاتف يجب أن يتكون من أرقام فقط وبحد أقصى 10 أرقام' });
      return;
    }
    if (values.password.length < MIN_PASSWORD_LENGTH) {
      setErrors({ password: `كلمة المرور ${MIN_PASSWORD_LENGTH} أحرف على الأقل` });
      return;
    }
    submit();
  };

  const field = (name, label, props = {}) => (
    <TextField name={name} label={label} value={values[name]} onChange={handleChange} error={Boolean(errors[name])} helperText={errors[name]} {...props} />
  );

  if (stationSubmitted) {
    return (
      <AuthLayout>
        <Typography variant="h4" component="h1">تم استلام طلبك</Typography>
        <Alert severity="info" sx={{ mt: 2, mb: 3, borderRadius: '14px' }}>
          طلب تسجيل المحطة قيد المراجعة. ستتمكن من الدخول بعد موافقة الإدارة، وتم إنشاء الفرع الرئيسي تلقائياً من الموقع الذي أدخلته.
        </Alert>
        <Button component={Link} to="/login" variant="contained" size="large" fullWidth>الرجوع لتسجيل الدخول</Button>
      </AuthLayout>
    );
  }

  return (
    <AuthLayout>
      <Typography variant="h4" component="h1">إنشاء حساب</Typography>
      <Typography color="text.secondary" sx={{ mt: 0.5, mb: 3 }}>اختار نوع الحساب ثم أكمل بياناتك.</Typography>

      <ToggleButtonGroup
        exclusive
        fullWidth
        value={accountType}
        onChange={(_, next) => next && setAccountType(next)}
        aria-label="نوع الحساب"
        sx={{ mb: 3, gap: 1.5, '& .MuiToggleButtonGroup-grouped': { border: '1px solid', borderColor: 'divider', borderRadius: '14px !important', margin: 0, py: 1.4 } }}
      >
        {Object.entries(ACCOUNT_TYPES).map(([value, { label }]) => (
          <ToggleButton key={value} value={value} color="primary">{label}</ToggleButton>
        ))}
      </ToggleButtonGroup>

      {!isStation && (
        <Alert severity="info" sx={{ mb: 2.5, borderRadius: '14px' }}>
          تحصل عند إنشاء حساب مستخدم على 10 نقاط ولاء مجانية (كل 10 نقاط = جالون وقود)، و+5 نقاط بعد كل شحن أو شراء وقود.
        </Alert>
      )}

      <Box component="form" onSubmit={handleSubmit} noValidate>
        <Stack spacing={2}>
          {field('name', isStation ? 'اسم المحطة' : 'الاسم الكامل')}
          {isStation && field('location', 'موقع المحطة (للفرع الرئيسي)')}
          {!isStation && field('phone', 'رقم الهاتف', { type: 'tel', placeholder: '09xxxxxxxx', inputProps: { maxLength: 10, inputMode: 'numeric' } })}
          {isStation && (
            <>
              {field('bank_name', 'البنك', { select: true, children: BANKS.map(({ value, label }) => <MenuItem key={value} value={value}>{label}</MenuItem>) })}
              {field('bank_account_number', 'رقم الحساب البنكي')}
              {field('bank_account_name', 'اسم صاحب الحساب')}
            </>
          )}
          {field('email', 'البريد الإلكتروني', { type: 'email', autoComplete: 'email' })}
          {field('password', 'كلمة المرور', { type: 'password', autoComplete: 'new-password' })}
          <Button type="submit" variant="contained" size="large" disabled={submitting}>
            {submitting ? 'جارٍ الإرسال...' : isStation ? 'إرسال طلب التسجيل' : 'إنشاء الحساب'}
          </Button>
        </Stack>
      </Box>

      <Typography variant="body2" color="text.secondary" sx={{ mt: 3, textAlign: 'center' }}>
        عندك حساب؟{' '}
        <Typography component={Link} to="/login" variant="body2" color="primary" fontWeight={700} sx={{ textDecoration: 'none' }}>
          تسجيل الدخول
        </Typography>
      </Typography>
    </AuthLayout>
  );
}
