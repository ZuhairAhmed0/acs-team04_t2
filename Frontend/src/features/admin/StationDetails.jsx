import { Box, Stack, Typography } from '@mui/material';
import CongestionChip from '../../components/common/CongestionChip';
import DetailRow from '../../components/common/DetailRow';
import FuelStatusChip from '../../components/common/FuelStatusChip';
import StatusChip from '../../components/common/StatusChip';
import { bankLabel, cityLabel } from '../../config/constants';
import { formatDate } from '../../utils/format';

/** Registration and bank details of a station, plus its branches. */
export default function StationDetails({ station }) {
  const branches = station.branches ?? [];

  return (
    <Box>
      <DetailRow label="حالة المحطة" value={<StatusChip status={station.status} />} />
      <DetailRow label="الموقع" value={station.location} />
      <DetailRow label="البريد الإلكتروني" value={station.email} />
      <DetailRow label="البنك" value={bankLabel(station.bank_name)} />
      <DetailRow label="رقم الحساب" value={station.bank_account_number} />
      <DetailRow label="اسم صاحب الحساب" value={station.bank_account_name} />
      <DetailRow label="تاريخ التسجيل" value={formatDate(station.created_at)} />
      <Typography variant="subtitle2" sx={{ mt: 2, mb: 1 }}>الفروع ({branches.length})</Typography>
      <Stack spacing={1.5}>
        {branches.map((branch) => (
          <Stack key={branch.id} spacing={0.75}>
            <Typography variant="body2" sx={{ overflowWrap: 'break-word' }}>
              {cityLabel(branch.city)} — {branch.address}{branch.is_main ? ' (رئيسي)' : ''}
            </Typography>
            <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
              <FuelStatusChip status={branch.fuel_status} />
              <CongestionChip percentage={branch.congestion_percentage} />
            </Stack>
          </Stack>
        ))}
      </Stack>
    </Box>
  );
}
