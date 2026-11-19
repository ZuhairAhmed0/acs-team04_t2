import { useMemo, useState } from 'react';
import { Box, Button, Card, MenuItem, TextField } from '@mui/material';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import { TRANSACTION_STATUSES, TRANSACTION_TYPES, statusMeta } from '../../config/constants';
import useRequest from '../../hooks/useRequest';
import useSubmit from '../../hooks/useSubmit';
import { downloadBlob } from '../../utils/download';
import PageHeader from './PageHeader';
import Reveal from './Reveal';
import { ErrorState, PageSkeleton } from './States';
import TransactionList from './TransactionList';

const EMPTY_FILTERS = { type: '', status: '', from_date: '', to_date: '' };
const dateFieldProps = { slotProps: { inputLabel: { shrink: true } } };

/** Full transaction history with filters and CSV export; shared by the user and station accounts. */
export default function TransactionsView({ title, subtitle, fetchTransactions, exportTransactions }) {
  const [filters, setFilters] = useState(EMPTY_FILTERS);
  const params = useMemo(() => Object.fromEntries(Object.entries(filters).filter(([, value]) => value)), [filters]);
  const hasFilters = Object.keys(params).length > 0;

  const { data, error, reload } = useRequest(() => fetchTransactions(params), [params]);
  const { submit: downloadCsv, submitting: exporting } = useSubmit(async () => {
    downloadBlob(await exportTransactions(params), 'transactions.csv');
  });

  const handleChange = (event) => setFilters((prev) => ({ ...prev, [event.target.name]: event.target.value }));

  return (
    <>
      <PageHeader
        title={title}
        subtitle={subtitle}
        action={
          <Button variant="outlined" startIcon={<DownloadRounded />} onClick={() => downloadCsv()} disabled={exporting}>
            {exporting ? 'جارٍ التحميل...' : 'تنزيل CSV'}
          </Button>
        }
      />

      <Reveal>
        <Card sx={{ p: { xs: 1.5, sm: 2 }, mb: 2.5 }}>
          <Box sx={{ display: 'grid', gap: { xs: 1.5, sm: 2 }, gridTemplateColumns: { xs: 'repeat(2, minmax(0, 1fr))', md: 'repeat(4, minmax(0, 1fr)) auto' }, alignItems: 'center' }}>
            <TextField select size="small" name="type" label="نوع العملية" value={filters.type} onChange={handleChange}>
              <MenuItem value="">الكل</MenuItem>
              {Object.entries(TRANSACTION_TYPES).map(([value, label]) => (
                <MenuItem key={value} value={value}>{label}</MenuItem>
              ))}
            </TextField>
            <TextField select size="small" name="status" label="الحالة" value={filters.status} onChange={handleChange}>
              <MenuItem value="">الكل</MenuItem>
              {TRANSACTION_STATUSES.map((value) => (
                <MenuItem key={value} value={value}>{statusMeta(value).label}</MenuItem>
              ))}
            </TextField>
            <TextField type="date" size="small" name="from_date" label="من تاريخ" value={filters.from_date} onChange={handleChange} {...dateFieldProps} />
            <TextField type="date" size="small" name="to_date" label="إلى تاريخ" value={filters.to_date} onChange={handleChange} {...dateFieldProps} />
            <Button
              color="inherit"
              disabled={!hasFilters}
              onClick={() => setFilters(EMPTY_FILTERS)}
              sx={{ gridColumn: { xs: '1 / -1', md: 'auto' } }}
            >
              مسح الفلاتر
            </Button>
          </Box>
        </Card>
      </Reveal>

      {error && <ErrorState message={error} onRetry={reload} />}
      {!error && !data && <PageSkeleton blocks={2} />}
      {data && (
        <Reveal delay={50}>
          <Card>
            <TransactionList
              transactions={data}
              emptyTitle={hasFilters ? 'لا توجد نتائج مطابقة' : 'لا توجد معاملات بعد'}
              emptyDescription={hasFilters ? 'غيّر الفلاتر أو امسحها لعرض كل المعاملات.' : undefined}
            />
          </Card>
        </Reveal>
      )}
    </>
  );
}
