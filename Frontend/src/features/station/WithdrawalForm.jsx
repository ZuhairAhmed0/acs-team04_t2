import { useRef, useState } from 'react';
import { Button, Card, Stack, Typography } from '@mui/material';
import AmountField from '../../components/common/AmountField';
import { useToast } from '../../context/ToastContext';
import useSubmit from '../../hooks/useSubmit';
import { submitCashWithdrawal } from '../../services/station.service';
import { parseAmount } from '../../utils/amount';
import { formatMoney } from '../../utils/format';
import { generateIdempotencyKey } from '../../utils/idempotency';

export default function WithdrawalForm({ balance, onSubmitted }) {
  const toast = useToast();
  const [amount, setAmount] = useState('');
  const [error, setError] = useState('');
  const { submit, submitting } = useSubmit(submitCashWithdrawal);
  const idempotencyKeyRef = useRef(generateIdempotencyKey());

  const handleSubmit = async (event) => {
    event.preventDefault();
    const parsed = parseAmount(amount);
    if (!parsed) return setError('أدخل مبلغاً صحيحاً أكبر من صفر');
    if (parsed > balance) return setError(`المبلغ يتجاوز رصيد المحطة (${formatMoney(balance)})`);

    const request = await submit({ amount: parsed, idempotency_key: idempotencyKeyRef.current });
    if (request) {
      idempotencyKeyRef.current = generateIdempotencyKey();
      setAmount('');
      toast.success('تم إرسال طلب السحب للإدارة');
      onSubmitted();
    }
    return undefined;
  };

  return (
    <Card sx={{ p: 2.5, height: '100%' }}>
      <Typography variant="h6">طلب سحب كاش</Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2.5 }}>يصل الطلب للإدارة للمراجعة، وبعد التحويل يُحدَّث وضعه هنا.</Typography>
      <Stack component="form" spacing={2} onSubmit={handleSubmit} noValidate>
        <AmountField
          value={amount}
          onChange={(event) => {
            idempotencyKeyRef.current = generateIdempotencyKey();
            setAmount(event.target.value);
            setError('');
          }}
          error={Boolean(error)}
          helperText={error}
        />
        <Button type="submit" variant="contained" size="large" disabled={submitting}>إرسال طلب السحب</Button>
      </Stack>
    </Card>
  );
}
