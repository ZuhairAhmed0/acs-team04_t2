import { Box, Card, CardActionArea, Stack, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import { tokens } from '../../theme/tokens';

export default function StatCard({ label, value, icon: Icon, hint, to }) {
  const content = (
    <Box sx={{ p: 2.5 }}>
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
        <Box>
          <Typography color="text.secondary" variant="body2">
            {label}
          </Typography>
          <Typography variant="h4" sx={{ mt: 0.5 }}>
            {value}
          </Typography>
          {hint && (
            <Typography variant="caption" color="text.secondary">
              {hint}
            </Typography>
          )}
        </Box>
        {Icon && (
          <Box sx={{ width: 44, height: 44, borderRadius: '14px', display: 'grid', placeItems: 'center', bgcolor: tokens.redSoft, color: 'primary.main' }}>
            <Icon />
          </Box>
        )}
      </Stack>
    </Box>
  );

  return (
    <Card sx={{ height: '100%' }}>
      {to ? (
        <CardActionArea component={Link} to={to} sx={{ height: '100%' }}>
          {content}
        </CardActionArea>
      ) : (
        content
      )}
    </Card>
  );
}
