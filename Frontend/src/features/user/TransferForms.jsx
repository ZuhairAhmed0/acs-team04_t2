import { useMemo, useRef, useState } from 'react';
import { Alert, Box, Button, Checkbox, Chip, FormControlLabel, MenuItem, Stack, TextField, Typography } from '@mui/material';
import AmountField from '../../components/common/AmountField';
import ConfirmDialog from '../../components/common/ConfirmDialog';
import useFormState from '../../hooks/useFormState';
import useSubmit from '../../hooks/useSubmit';
import { parseAmount, parseNonNegativeAmount } from '../../utils/amount';
import { formatMoney, formatNumber } from '../../utils/format';
import { generateIdempotencyKey } from '../../utils/idempotency';
import { allocateToStation, purchaseFuel, topUpWallet, transferBetweenStations } from '../../services/transfer.service';

const AMOUNT_ERROR = 'أدخل مبلغاً صحيحاً أكبر من صفر';
const DUPLICATE_WINDOW_MS = 10 * 60 * 1000;
const DEFAULT_LOYALTY_POINT_VALUE = 100;
const POINTS_PER_GALLON = 10;

function StationSelect({ stations, error, ...props }) {
  return (
    <TextField select error={Boolean(error)} helperText={error} {...props}>
      {stations.map(({ id, name }) => (
        <MenuItem key={id} value={id}>{name}</MenuItem>
      ))}
    </TextField>
  );
}

/** Shared behaviour for the forms: validate, attach stable per-attempt idempotency key, submit, reset and report the receipt. */
function useTransferForm(initialValues, validate, request, onSuccess) {
  const form = useFormState(initialValues);
  const { submit, submitting } = useSubmit(request);
  const idempotencyKeyRef = useRef(generateIdempotencyKey());

  const handleChange = (event) => {
    idempotencyKeyRef.current = generateIdempotencyKey();
    form.handleChange(event);
  };

  const submitValidated = async (payload) => {
    const receipt = await submit({ ...payload, idempotency_key: idempotencyKeyRef.current });
    if (receipt) {
      idempotencyKeyRef.current = generateIdempotencyKey();
      form.reset();
      onSuccess(receipt, payload);
    }
    return receipt;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const payload = validate(form);
    if (!payload) return;
    await submitValidated(payload);
  };

  const setFieldValue = (name, value) => {
    idempotencyKeyRef.current = generateIdempotencyKey();
    form.handleChange({ target: { name, value } });
  };

  return { ...form, handleChange, setFieldValue, submitting, handleSubmit, validatePayload: () => validate(form), submitValidated };
}

export function TopUpForm({ onSuccess }) {
  const form = useTransferForm(
    { amount: '' },
    ({ values, setErrors }) => {
      const amount = parseAmount(values.amount);
      if (!amount) return setErrors({ amount: AMOUNT_ERROR });
      return { amount };
    },
    topUpWallet,
    onSuccess,
  );

  return (
    <Stack component="form" spacing={2} onSubmit={form.handleSubmit} noValidate>
      <AmountField name="amount" value={form.values.amount} onChange={form.handleChange} error={Boolean(form.errors.amount)} helperText={form.errors.amount} />
      <Button type="submit" variant="contained" size="large" disabled={form.submitting}>شحن المحفظة</Button>
    </Stack>
  );
}

/** No branch is selected: the allocation goes to the station's shared balance. */
export function AllocateForm({ stations, availableBalance, onSuccess }) {
  const form = useTransferForm(
    { station_id: '', amount: '' },
    ({ values, setErrors }) => {
      const amount = parseAmount(values.amount);
      const errors = {
        ...(!values.station_id && { station_id: 'اختر المحطة' }),
        ...(!amount && { amount: AMOUNT_ERROR }),
        ...(amount && typeof availableBalance === 'number' && amount > availableBalance && {
          amount: `المبلغ يتجاوز الرصيد المتاح (${formatMoney(availableBalance)})`,
        }),
      };
      if (Object.keys(errors).length) return setErrors(errors);
      return { station_id: values.station_id, amount };
    },
    allocateToStation,
    onSuccess,
  );

  return (
    <Stack component="form" spacing={2} onSubmit={form.handleSubmit} noValidate>
      {typeof availableBalance === 'number' && (
        <Typography variant="body2" color="text.secondary">
          الرصيد المتاح: <strong>{formatMoney(availableBalance)}</strong>
        </Typography>
      )}
      <StationSelect name="station_id" label="المحطة" stations={stations} value={form.values.station_id} onChange={form.handleChange} error={form.errors.station_id} />
      <AmountField name="amount" value={form.values.amount} onChange={form.handleChange} error={Boolean(form.errors.amount)} helperText={form.errors.amount} />
      <Button type="submit" variant="contained" size="large" disabled={form.submitting}>تأكيد التخصيص</Button>
    </Stack>
  );
}

/** Source and destination are stations only; branches are never part of the request. */
export function StationTransferForm({ stations, allocatedBalance, onSuccess }) {
  const form = useTransferForm(
    { from_station_id: '', to_station_id: '', amount: '' },
    ({ values, setErrors }) => {
      const amount = parseAmount(values.amount);
      const errors = {
        ...(!values.from_station_id && { from_station_id: 'اختر المحطة المصدر' }),
        ...(!values.to_station_id && { to_station_id: 'اختر المحطة الوجهة' }),
        ...(values.from_station_id && values.from_station_id === values.to_station_id && { to_station_id: 'اختر محطتين مختلفتين' }),
        ...(!amount && { amount: AMOUNT_ERROR }),
        ...(amount && typeof allocatedBalance === 'number' && amount > allocatedBalance && {
          amount: `المبلغ يتجاوز الرصيد المخصص (${formatMoney(allocatedBalance)})`,
        }),
      };
      if (Object.keys(errors).length) return setErrors(errors);
      return { from_station_id: values.from_station_id, to_station_id: values.to_station_id, amount };
    },
    transferBetweenStations,
    onSuccess,
  );

  return (
    <Stack component="form" spacing={2} onSubmit={form.handleSubmit} noValidate>
      {typeof allocatedBalance === 'number' && (
        <Typography variant="body2" color="text.secondary">
          الرصيد المخصص المتاح للتحويل: <strong>{formatMoney(allocatedBalance)}</strong>
        </Typography>
      )}
      <StationSelect name="from_station_id" label="من المحطة" stations={stations} value={form.values.from_station_id} onChange={form.handleChange} error={form.errors.from_station_id} />
      <StationSelect name="to_station_id" label="إلى المحطة" stations={stations} value={form.values.to_station_id} onChange={form.handleChange} error={form.errors.to_station_id} />
      <AmountField name="amount" value={form.values.amount} onChange={form.handleChange} error={Boolean(form.errors.amount)} helperText={form.errors.amount} />
      <Button type="submit" variant="contained" size="large" disabled={form.submitting}>تأكيد التحويل</Button>
    </Stack>
  );
}

/** Records a branch-level fuel_purchase transaction for a tapped station branch, with optional loyalty redemption. */
export function FuelPurchaseForm({
  branch,
  recentPurchases = [],
  availableBalance,
  allocatedBalance,
  loyaltyPoints = 0,
  loyaltyPointValue = DEFAULT_LOYALTY_POINT_VALUE,
  isRegularUser = false,
  onSuccess,
  onCancel,
}) {
  const [duplicateConfirm, setDuplicateConfirm] = useState(null);
  const pointValue = Number(loyaltyPointValue) || DEFAULT_LOYALTY_POINT_VALUE;
  const gallonValue = POINTS_PER_GALLON * pointValue;
  const safeLoyaltyPoints = Math.max(0, Number(loyaltyPoints) || 0);
  const redeemableGallons = Math.floor(safeLoyaltyPoints / POINTS_PER_GALLON);
  const redeemablePoints = redeemableGallons * POINTS_PER_GALLON;
  const redeemableValue = redeemablePoints * pointValue;

  const form = useTransferForm(
    {
      payment_source: 'available',
      amount: '',
      use_loyalty: false,
      loyalty_points_used: '',
    },
    ({ values, setErrors }) => {
      const source = values.payment_source;
      const errors = {
        ...(!branch?.branch_id && { amount: 'اختر فرعاً متاحاً أولاً' }),
        ...(branch?.fuel_status === 'empty' && { amount: 'الفرع المختار فارغ من الوقود' }),
      };

      let amount = 0;
      let pointsUsed = 0;
      if (source === 'loyalty') {
        const parsedPoints = parseAmount(values.loyalty_points_used);
        if (!isRegularUser) {
          errors.payment_source = 'نقاط الولاء متاحة لحسابات المستخدمين فقط';
        } else if (redeemablePoints < POINTS_PER_GALLON) {
          errors.payment_source = 'تحتاج إلى 10 نقاط على الأقل لاستبدال الوقود';
        } else if (!parsedPoints) {
          errors.loyalty_points_used = 'أدخل عدد نقاط صحيحاً (10، 20، 30...)';
        } else if (parsedPoints % POINTS_PER_GALLON !== 0) {
          errors.loyalty_points_used = 'يمكنك استخدام 10 نقاط أو مضاعفاتها فقط (10، 20، 30...)';
        } else if (parsedPoints > redeemablePoints) {
          errors.loyalty_points_used = `النقاط المطلوبة تتجاوز رصيدك القابل للاستبدال (${formatNumber(redeemablePoints)} نقطة)`;
        } else {
          pointsUsed = parsedPoints;
        }
        amount = 0;
      } else {
        const sourceBalance = source === 'allocated' ? allocatedBalance : availableBalance;
        const usingPoints = isRegularUser && values.use_loyalty;
        const parsedAmount = usingPoints ? parseNonNegativeAmount(values.amount) : parseAmount(values.amount);

        if (parsedAmount === null) {
          errors.amount = usingPoints ? 'أدخل مبلغاً صحيحاً (0 أو أكثر)' : AMOUNT_ERROR;
        } else {
          amount = parsedAmount;
        }

        if (usingPoints) {
          const parsedPoints = parseAmount(values.loyalty_points_used);
          if (!parsedPoints) {
            errors.loyalty_points_used = 'أدخل عدد نقاط صحيحاً (10، 20، 30...)';
          } else if (parsedPoints % POINTS_PER_GALLON !== 0) {
            errors.loyalty_points_used = 'يمكنك استخدام 10 نقاط أو مضاعفاتها فقط (10، 20، 30...)';
          } else if (parsedPoints > redeemablePoints) {
            errors.loyalty_points_used = `النقاط المطلوبة تتجاوز رصيدك القابل للاستبدال (${formatNumber(redeemablePoints)} نقطة)`;
          } else if (amount > 0 && parsedPoints * pointValue > amount) {
            errors.loyalty_points_used = 'خصم النقاط لا يمكن أن يتجاوز قيمة الشراء';
          } else {
            pointsUsed = parsedPoints;
          }
        }

        const walletCharge = amount === 0 ? 0 : Math.max(0, amount - pointsUsed * pointValue);
        if (parsedAmount !== null && typeof sourceBalance === 'number' && walletCharge > sourceBalance) {
          errors.amount = `المبلغ المخصوم من المحفظة يتجاوز الرصيد المختار (${formatMoney(sourceBalance)})`;
        }
      }

      if (Object.keys(errors).length) return setErrors(errors);
      return {
        branch_id: branch.branch_id,
        amount,
        payment_source: source,
        ...(pointsUsed > 0 && { loyalty_points_used: pointsUsed }),
      };
    },
    purchaseFuel,
    onSuccess,
  );

  const isLoyaltySource = form.values.payment_source === 'loyalty';
  const usingPoints = isLoyaltySource || (isRegularUser && form.values.use_loyalty);
  const parsedAmount = usingPoints ? parseNonNegativeAmount(form.values.amount) : parseAmount(form.values.amount);
  const parsedCustomPoints = parseAmount(form.values.loyalty_points_used);
  const activePointsUsed = usingPoints && parsedCustomPoints && parsedCustomPoints % POINTS_PER_GALLON === 0
    ? parsedCustomPoints
    : 0;
  const loyaltyDiscount = activePointsUsed * pointValue;
  const estimatedWalletCharge = isLoyaltySource || parsedAmount === 0 || parsedAmount === null
    ? 0
    : Math.max(0, parsedAmount - loyaltyDiscount);
  const maxRedeemablePoints = !isLoyaltySource && parsedAmount && parsedAmount > 0
    ? Math.min(redeemablePoints, Math.floor(parsedAmount / gallonValue) * POINTS_PER_GALLON)
    : redeemablePoints;

  const matchingRecentPurchase = useMemo(() => {
    if (!branch?.branch_id) return null;
    if (!isLoyaltySource && parsedAmount === null) return null;
    if (!estimatedWalletCharge && !activePointsUsed) return null;
    const cutoff = Date.now() - DUPLICATE_WINDOW_MS;
    return recentPurchases.find((item) => {
      if (item.type !== 'fuel_purchase' || item.status !== 'completed') return false;
      if (
        Number(item.branch_id) !== Number(branch.branch_id) ||
        Number(item.amount) !== estimatedWalletCharge ||
        Number(item.loyalty_points_used ?? 0) !== activePointsUsed
      ) {
        return false;
      }
      const createdMs = item.created_at ? new Date(item.created_at).getTime() : Date.now();
      return Number.isNaN(createdMs) || createdMs >= cutoff;
    }) ?? null;
  }, [activePointsUsed, branch?.branch_id, estimatedWalletCharge, isLoyaltySource, parsedAmount, recentPurchases]);

  const handlePurchaseSubmit = async (event) => {
    event.preventDefault();
    const payload = form.validatePayload();
    if (!payload) return;
    if (matchingRecentPurchase) {
      setDuplicateConfirm({ payload, previous: matchingRecentPurchase });
      return;
    }
    await form.submitValidated(payload);
  };

  const handleConfirmDuplicate = async () => {
    if (!duplicateConfirm) return;
    await form.submitValidated(duplicateConfirm.payload);
    setDuplicateConfirm(null);
  };

  return (
    <>
      <Stack component="form" spacing={2} onSubmit={handlePurchaseSubmit} noValidate>
        <TextField
          select
          name="payment_source"
          label="مصدر الدفع"
          value={form.values.payment_source}
          onChange={(event) => {
            const nextSource = event.target.value;
            form.handleChange(event);
            if (nextSource === 'loyalty') {
              form.setFieldValue('use_loyalty', false);
              form.setFieldValue('amount', '0');
              if (redeemablePoints >= POINTS_PER_GALLON) {
                form.setFieldValue('loyalty_points_used', String(POINTS_PER_GALLON));
              }
            } else if (form.values.amount === '0') {
              form.setFieldValue('amount', '');
            }
          }}
          error={Boolean(form.errors.payment_source)}
          helperText={form.errors.payment_source}
        >
          <MenuItem value="available">
            الرصيد المتاح{typeof availableBalance === 'number' ? ` (${formatMoney(availableBalance)})` : ''}
          </MenuItem>
          <MenuItem value="allocated">
            الرصيد المخصص{typeof allocatedBalance === 'number' ? ` (${formatMoney(allocatedBalance)})` : ''}
          </MenuItem>
          {isRegularUser && (
            <MenuItem value="loyalty" disabled={redeemablePoints < POINTS_PER_GALLON}>
              نقاط الولاء فقط — 0 ج.س ({formatNumber(redeemablePoints)} نقطة قابلة للاستبدال)
            </MenuItem>
          )}
        </TextField>

        {isLoyaltySource ? (
          <Stack spacing={1.5}>
            {redeemablePoints < POINTS_PER_GALLON ? (
              <Alert severity="info" sx={{ borderRadius: '14px' }}>
                تحتاج إلى 10 نقاط على الأقل لاستبدال جالون وقود (كل 10 نقاط = {formatMoney(gallonValue)}).
              </Alert>
            ) : (
              <>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                  <Chip
                    label={`10 نقاط / 1 جالون`}
                    size="small"
                    color="primary"
                    variant="outlined"
                    clickable
                    onClick={() => form.setFieldValue('loyalty_points_used', String(POINTS_PER_GALLON))}
                  />
                  {redeemablePoints > POINTS_PER_GALLON && (
                    <Chip
                      label={`الكل: ${formatNumber(redeemablePoints)} نقطة (${formatNumber(redeemableGallons)} جالون)`}
                      size="small"
                      color="primary"
                      variant="outlined"
                      clickable
                      onClick={() => form.setFieldValue('loyalty_points_used', String(redeemablePoints))}
                    />
                  )}
                </Stack>
                <TextField
                  fullWidth
                  type="number"
                  name="loyalty_points_used"
                  label="عدد النقاط للاستبدال (مضاعفات 10: 10، 20، 30...)"
                  value={form.values.loyalty_points_used}
                  onChange={form.handleChange}
                  error={Boolean(form.errors.loyalty_points_used)}
                  helperText={form.errors.loyalty_points_used || `0 ج.س من المحفظة • كل 10 نقاط = 1 جالون (${formatMoney(gallonValue)})`}
                  slotProps={{ htmlInput: { min: POINTS_PER_GALLON, max: redeemablePoints, step: POINTS_PER_GALLON, inputMode: 'numeric' } }}
                />
              </>
            )}
          </Stack>
        ) : (
          <AmountField
            name="amount"
            min={form.values.use_loyalty ? 0 : 1}
            label={form.values.use_loyalty ? 'المبلغ (يمكن إدخال 0 ج.س لاستخدام النقاط فقط)' : 'المبلغ'}
            value={form.values.amount}
            onChange={form.handleChange}
            error={Boolean(form.errors.amount)}
            helperText={form.errors.amount}
          />
        )}

        {!isLoyaltySource && isRegularUser && redeemablePoints >= POINTS_PER_GALLON && (
          <Box sx={{ p: 1.5, borderRadius: '14px', border: '1px solid', borderColor: 'divider' }}>
            <FormControlLabel
              control={
                <Checkbox
                  checked={Boolean(form.values.use_loyalty)}
                  onChange={(event) => {
                    const checked = event.target.checked;
                    form.setFieldValue('use_loyalty', checked);
                    if (!checked) {
                      form.setFieldValue('loyalty_points_used', '');
                      if (form.values.amount === '0') form.setFieldValue('amount', '');
                    } else if (!form.values.loyalty_points_used && maxRedeemablePoints >= POINTS_PER_GALLON) {
                      form.setFieldValue('loyalty_points_used', String(maxRedeemablePoints));
                    }
                  }}
                />
              }
              label={`استخدام نقاط الولاء (القابل للاستبدال: ${formatNumber(redeemablePoints)} نقطة = ${formatMoney(redeemableValue)})`}
            />
            {form.values.use_loyalty && (
              <Stack spacing={1} sx={{ mt: 1 }}>
                <Stack direction="row" spacing={1} alignItems="flex-start">
                  <TextField
                    fullWidth
                    size="small"
                    type="number"
                    name="loyalty_points_used"
                    label="عدد النقاط (مضاعفات 10: 10، 20، 30...)"
                    value={form.values.loyalty_points_used}
                    onChange={form.handleChange}
                    error={Boolean(form.errors.loyalty_points_used)}
                    helperText={form.errors.loyalty_points_used || `كل 10 نقاط = ${formatMoney(gallonValue)} (1 جالون)`}
                    slotProps={{ htmlInput: { min: POINTS_PER_GALLON, max: redeemablePoints, step: POINTS_PER_GALLON, inputMode: 'numeric' } }}
                  />
                  {maxRedeemablePoints >= POINTS_PER_GALLON && (
                    <Button
                      size="small"
                      variant="outlined"
                      sx={{ mt: 0.5, flexShrink: 0 }}
                      onClick={() => form.setFieldValue('loyalty_points_used', String(maxRedeemablePoints))}
                    >
                      الحد الأقصى
                    </Button>
                  )}
                </Stack>
              </Stack>
            )}
          </Box>
        )}

        {activePointsUsed > 0 && !form.errors.amount && !form.errors.loyalty_points_used && (
          <Alert severity="info" sx={{ borderRadius: '14px' }}>
            {isLoyaltySource || parsedAmount === 0
              ? `سيتم شراء وقود بقيمة ${formatMoney(0)} من محفظتك باستخدام ${formatNumber(activePointsUsed)} نقطة ولاء (${formatMoney(loyaltyDiscount)}).`
              : `قيمة الوقود ${formatMoney(parsedAmount)} — خصم نقاط الولاء ${formatMoney(loyaltyDiscount)} (${formatNumber(activePointsUsed)} نقطة) — يُخصم من المحفظة ${formatMoney(estimatedWalletCharge)}.`}
          </Alert>
        )}

        {matchingRecentPurchase && (
          <Alert severity="warning" sx={{ borderRadius: '14px' }}>
            تنبيه: قمت بنفس عملية شراء الوقود (
            {matchingRecentPurchase.amount === 0 && matchingRecentPurchase.loyalty_points_used > 0
              ? `${formatNumber(matchingRecentPurchase.loyalty_points_used)} نقطة ولاء`
              : formatMoney(matchingRecentPurchase.amount)}
            ) في هذا الفرع مؤخراً
            {matchingRecentPurchase.reference ? ` (مرجع: ${matchingRecentPurchase.reference})` : ''}.
          </Alert>
        )}
        <Stack direction="row" spacing={1.5} justifyContent="flex-end">
          {onCancel && (
            <Button color="inherit" onClick={onCancel} disabled={form.submitting}>
              إلغاء
            </Button>
          )}
          <Button type="submit" variant="contained" size="large" disabled={form.submitting}>
            تأكيد شراء الوقود
          </Button>
        </Stack>
      </Stack>

      <ConfirmDialog
        open={Boolean(duplicateConfirm)}
        title="تنبيه تكرار عملية الشراء"
        message={
          duplicateConfirm
            ? `لقد نفذت نفس عملية شراء الوقود (${
              duplicateConfirm.payload.amount === 0 && duplicateConfirm.payload.loyalty_points_used > 0
                ? `${formatNumber(duplicateConfirm.payload.loyalty_points_used)} نقطة ولاء`
                : formatMoney(duplicateConfirm.payload.amount)
            }) في هذا الفرع مؤخراً${duplicateConfirm.previous?.reference ? ` (مرجع: ${duplicateConfirm.previous.reference})` : ''}. هل تريد بالتأكيد تنفيذ العملية مرة أخرى؟`
            : ''
        }
        confirmLabel="نعم، كرر الشراء"
        tone="warning"
        loading={form.submitting}
        onConfirm={handleConfirmDuplicate}
        onClose={() => setDuplicateConfirm(null)}
      />
    </>
  );
}
