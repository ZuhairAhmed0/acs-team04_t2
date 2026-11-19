import { useState } from 'react';
import { Alert, Box, Button, IconButton, InputAdornment, Stack, TextField, Typography } from '@mui/material';
import EmailRounded from '@mui/icons-material/EmailRounded';
import LockRounded from '@mui/icons-material/LockRounded';
import VisibilityOffRounded from '@mui/icons-material/VisibilityOffRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import AuthLayout from '../../components/layout/AuthLayout';
import { ROLE_HOME } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import useFormState from '../../hooks/useFormState';
import useSubmit from '../../hooks/useSubmit';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [showPassword, setShowPassword] = useState(false);
  const { values, errors, handleChange, validateRequired } = useFormState({ email: '', password: '' });
  const [formError, setFormError] = useState('');
  const justRegistered = location.state?.registered;

  const { submit, submitting } = useSubmit(async () => {
    try {
      const user = await login(values);
      navigate(location.state?.from?.pathname ?? ROLE_HOME[user.role], { replace: true });
    } catch (error) {
      setFormError(error?.response?.data?.message ?? 'تعذّر تسجيل الدخول');
    }
  });

  const handleSubmit = (event) => {
    event.preventDefault();
    setFormError('');
    if (validateRequired(['email', 'password'])) submit();
  };

  return (
    <AuthLayout>
      <Typography variant="h4" component="h1">تسجيل الدخول</Typography>
      <Typography color="text.secondary" sx={{ mt: 0.5, mb: 3 }}>
        أهلاً بيك، ادخل بحسابك للمتابعة.
      </Typography>

      {justRegistered && <Alert severity="success" sx={{ mb: 2, borderRadius: '14px' }}>تم إنشاء الحساب. سجّل دخولك للمتابعة.</Alert>}
      {formError && <Alert severity="error" sx={{ mb: 2, borderRadius: '14px' }}>{formError}</Alert>}

      <Box component="form" onSubmit={handleSubmit} noValidate>
        <Stack spacing={2}>
          <TextField
            name="email"
            type="email"
            label="البريد الإلكتروني"
            value={values.email}
            onChange={handleChange}
            error={Boolean(errors.email)}
            helperText={errors.email}
            autoComplete="email"
            slotProps={{ input: { startAdornment: <InputAdornment position="start"><EmailRounded fontSize="small" /></InputAdornment> } }}
          />
          <TextField
            name="password"
            type={showPassword ? 'text' : 'password'}
            label="كلمة المرور"
            value={values.password}
            onChange={handleChange}
            error={Boolean(errors.password)}
            helperText={errors.password}
            autoComplete="current-password"
            slotProps={{
              input: {
                startAdornment: <InputAdornment position="start"><LockRounded fontSize="small" /></InputAdornment>,
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton aria-label={showPassword ? 'إخفاء كلمة المرور' : 'إظهار كلمة المرور'} onClick={() => setShowPassword((prev) => !prev)} edge="end">
                      {showPassword ? <VisibilityOffRounded /> : <VisibilityRounded />}
                    </IconButton>
                  </InputAdornment>
                ),
              },
            }}
          />
          <Button type="submit" variant="contained" size="large" disabled={submitting}>
            {submitting ? 'جارٍ الدخول...' : 'دخول'}
          </Button>
        </Stack>
      </Box>

      <Typography variant="body2" color="text.secondary" sx={{ mt: 3, textAlign: 'center' }}>
        ما عندك حساب؟{' '}
        <Typography component={Link} to="/register" variant="body2" color="primary" fontWeight={700} sx={{ textDecoration: 'none' }}>
          إنشاء حساب
        </Typography>
      </Typography>
    </AuthLayout>
  );
}
