import { Fragment, useMemo, useState } from 'react';
import { Avatar, Box, Button, Card, Chip, Divider, InputAdornment, Stack, TextField, Typography } from '@mui/material';
import SearchRounded from '@mui/icons-material/SearchRounded';
import ActionDialog from '../../components/common/ActionDialog';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import { EmptyState, ErrorState, PageSkeleton } from '../../components/common/States';
import StatusChip from '../../components/common/StatusChip';
import { ROLES } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import useConfirmedAction from '../../hooks/useConfirmedAction';
import useRequest from '../../hooks/useRequest';
import { fetchUsers, updateUserRole, updateUserStatus } from '../../services/admin.service';
import { tokens } from '../../theme/tokens';

const ACTIONS = {
  suspend: {
    label: 'إيقاف', tone: 'error', title: 'إيقاف الحساب', confirmLabel: 'إيقاف الحساب', success: 'تم إيقاف الحساب',
    message: (user) => `لن يتمكن ${user.name} من الدخول حتى يُعاد تفعيل حسابه.`,
    run: (user) => updateUserStatus(user.id, 'suspended'),
  },
  activate: {
    label: 'تفعيل', tone: 'success', title: 'تفعيل الحساب', confirmLabel: 'تفعيل الحساب', success: 'تم تفعيل الحساب',
    message: (user) => `سيتمكن ${user.name} من الدخول مرة أخرى.`,
    run: (user) => updateUserStatus(user.id, 'active'),
  },
  promote: {
    label: 'ترقية إلى مدير', tone: 'primary', title: 'ترقية إلى مدير', confirmLabel: 'ترقية', success: 'تمت الترقية إلى مدير',
    message: (user) => `سيحصل ${user.name} على صلاحيات إدارة النظام كاملة.`,
    run: (user) => updateUserRole(user.id, ROLES.ADMIN),
  },
  revoke: {
    label: 'سحب صلاحية المدير', tone: 'error', title: 'سحب صلاحية المدير', confirmLabel: 'سحب الصلاحية', success: 'تم إرجاع الحساب إلى مستخدم',
    message: (user) => `سيعود ${user.name} إلى حساب مستخدم عادي.`,
    run: (user) => updateUserRole(user.id, ROLES.USER),
  },
};

const actionsFor = (user) => [
  user.status === 'active' ? 'suspend' : 'activate',
  user.role === ROLES.ADMIN ? 'revoke' : 'promote',
];

function UserRow({ user, isSelf, onAction }) {
  return (
    <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'flex-start', md: 'center' }} sx={{ px: 2.5, py: 2 }}>
      <Stack direction="row" spacing={1.5} alignItems="center" sx={{ flex: 1, minWidth: 0 }}>
        <Avatar sx={{ bgcolor: tokens.redSoft, color: 'primary.main', fontWeight: 800 }}>{user.name[0]}</Avatar>
        <Box sx={{ minWidth: 0 }}>
          <Typography fontWeight={700}>{user.name}{isSelf && ' (أنت)'}</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ overflowWrap: 'break-word' }}>{user.email} — {user.phone}</Typography>
        </Box>
      </Stack>
      <Stack direction="row" spacing={1}>
        {user.role === ROLES.ADMIN && <Chip size="small" color="primary" label="مدير" />}
        <StatusChip status={user.status} />
      </Stack>
      <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ width: { xs: '100%', md: 'auto' } }}>
        {actionsFor(user).map((key) => (
          <Button key={key} size="small" variant="outlined" color={ACTIONS[key].tone} disabled={isSelf} onClick={() => onAction(ACTIONS[key], user)}>
            {ACTIONS[key].label}
          </Button>
        ))}
      </Stack>
    </Stack>
  );
}

export default function AdminUsers() {
  const { user: me } = useAuth();
  const [query, setQuery] = useState('');
  const { data, error, reload } = useRequest(fetchUsers);
  const controller = useConfirmedAction(reload);

  const users = useMemo(() => {
    const needle = query.trim().toLowerCase();
    return (data ?? []).filter((user) => !needle || [user.name, user.email, user.phone].some((value) => value.toLowerCase().includes(needle)));
  }, [data, query]);

  return (
    <>
      <PageHeader title="المستخدمون" subtitle="فعّل الحسابات أو أوقفها، وأدر صلاحيات المديرين." />
      <TextField
        value={query}
        onChange={(event) => setQuery(event.target.value)}
        placeholder="ابحث بالاسم أو البريد أو الهاتف"
        inputProps={{ 'aria-label': 'بحث في المستخدمين' }}
        slotProps={{ input: { startAdornment: <InputAdornment position="start"><SearchRounded /></InputAdornment> } }}
        sx={{ mb: 2.5, maxWidth: 420 }}
      />

      {error && <ErrorState message={error} onRetry={reload} />}
      {!error && !data && <PageSkeleton blocks={2} />}
      {data && (
        <Reveal>
          <Card>
            {!users.length && <EmptyState title="لا يوجد مستخدمون مطابقون" description="جرّب كلمة بحث أخرى." />}
            {users.map((user, index) => (
              <Fragment key={user.id}>
                {index > 0 && <Divider />}
                <UserRow user={user} isSelf={user.id === me.id} onAction={controller.request} />
              </Fragment>
            ))}
          </Card>
        </Reveal>
      )}
      <ActionDialog controller={controller} />
    </>
  );
}
