import { Box, Button, Skeleton, Stack, Typography } from '@mui/material';
import InboxRounded from '@mui/icons-material/InboxRounded';
import ErrorOutlineRounded from '@mui/icons-material/ErrorOutlineRounded';
import { tokens } from '../../theme/tokens';

export function EmptyState({ title, description, action, icon: Icon = InboxRounded }) {
  return (
    <Stack alignItems="center" spacing={1} sx={{ py: 6, px: 2, textAlign: 'center' }}>
      <Box sx={{ width: 56, height: 56, borderRadius: '50%', display: 'grid', placeItems: 'center', bgcolor: tokens.redSoft, color: 'primary.main' }}>
        <Icon />
      </Box>
      <Typography variant="subtitle1">{title}</Typography>
      {description && (
        <Typography color="text.secondary" variant="body2" sx={{ maxWidth: 360 }}>
          {description}
        </Typography>
      )}
      {action}
    </Stack>
  );
}

export function ErrorState({ message, onRetry }) {
  return (
    <EmptyState
      icon={ErrorOutlineRounded}
      title="تعذّر تحميل البيانات"
      description={message}
      action={onRetry && <Button onClick={onRetry} variant="outlined" sx={{ mt: 1 }}>إعادة المحاولة</Button>}
    />
  );
}

export function PageSkeleton({ blocks = 3 }) {
  return (
    <Stack spacing={2.5}>
      <Skeleton variant="text" width={220} height={44} />
      {Array.from({ length: blocks }, (_, index) => (
        <Skeleton key={index} variant="rounded" height={index === 0 ? 180 : 120} sx={{ borderRadius: '20px' }} />
      ))}
    </Stack>
  );
}
