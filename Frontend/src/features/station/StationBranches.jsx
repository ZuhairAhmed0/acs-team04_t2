import { useEffect, useState } from 'react';
import { Alert, Box, Button, Card, Chip, Slider, Stack, Typography } from '@mui/material';
import AddRounded from '@mui/icons-material/AddRounded';
import FuelStatusToggle from '../../components/common/FuelStatusToggle';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import { ErrorState, PageSkeleton } from '../../components/common/States';
import { cityLabel } from '../../config/constants';
import { useToast } from '../../context/ToastContext';
import useRequest from '../../hooks/useRequest';
import useSubmit from '../../hooks/useSubmit';
import { fetchBranches, updateBranchCongestion, updateBranchFuelStatus } from '../../services/station.service';
import AddBranchDialog from './AddBranchDialog';

function BranchCard({ branch, onChanged }) {
  const toast = useToast();
  const [congestion, setCongestion] = useState(branch.congestion_percentage ?? 0);
  const { submit, submitting } = useSubmit(updateBranchFuelStatus);
  const { submit: saveCongestion, submitting: savingCongestion } = useSubmit(updateBranchCongestion);

  useEffect(() => {
    setCongestion(branch.congestion_percentage ?? 0);
  }, [branch.congestion_percentage]);

  const handleStatusChange = async (status) => {
    const updated = await submit(branch.id, status);
    if (updated) {
      toast.success('تم تحديث حالة الوقود');
      onChanged();
    }
  };

  const handleCongestionCommit = async (_, value) => {
    const updated = await saveCongestion(branch.id, value);
    if (updated) {
      toast.success('تم تحديث نسبة الازدحام');
      onChanged();
    }
  };

  return (
    <Card sx={{ p: { xs: 2, sm: 2.5 }, height: '100%' }}>
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start" spacing={1} sx={{ mb: 2 }}>
        <Box sx={{ minWidth: 0 }}>
          <Typography variant="h6">{cityLabel(branch.city)}</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ overflowWrap: 'break-word' }}>{branch.address}</Typography>
        </Box>
        {branch.is_main && <Chip size="small" color="primary" label="الفرع الرئيسي" sx={{ flexShrink: 0 }} />}
      </Stack>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>حالة الوقود</Typography>
      <FuelStatusToggle value={branch.fuel_status} onChange={handleStatusChange} disabled={submitting} />
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mt: 2 }}>
        <Typography variant="body2" color="text.secondary">نسبة الازدحام</Typography>
        <Typography variant="body2" fontWeight={700}>{congestion}%</Typography>
      </Stack>
      <Slider
        value={congestion}
        min={0}
        max={100}
        step={5}
        valueLabelDisplay="auto"
        disabled={savingCongestion}
        aria-label={`نسبة الازدحام لفرع ${cityLabel(branch.city)}`}
        sx={{ py: 1.5 }}
        onChange={(_, value) => setCongestion(value)}
        onChangeCommitted={handleCongestionCommit}
      />
    </Card>
  );
}

export default function StationBranches() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const { data, error, reload } = useRequest(fetchBranches);

  return (
    <>
      <PageHeader
        title="الفروع"
        subtitle="حدّث حالة الوقود في كل فرع وأضف فروعاً جديدة."
        action={<Button variant="contained" startIcon={<AddRounded />} onClick={() => setDialogOpen(true)}>إضافة فرع</Button>}
      />
      <Alert severity="info" sx={{ mb: 3, borderRadius: '14px' }}>كل الفروع تشترك في رصيد المحطة نفسه، ولا يُقسَّم الرصيد بينها.</Alert>

      {error && <ErrorState message={error} onRetry={reload} />}
      {!error && !data && <PageSkeleton blocks={2} />}
      {data && (
        <Box sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: 'repeat(auto-fill, minmax(min(100%, 320px), 1fr))' }}>
          {data.map((branch, index) => (
            <Reveal key={branch.id} delay={Math.min(index, 5) * 50}>
              <BranchCard branch={branch} onChanged={reload} />
            </Reveal>
          ))}
        </Box>
      )}

      <AddBranchDialog open={dialogOpen} onClose={() => setDialogOpen(false)} onAdded={reload} />
    </>
  );
}
