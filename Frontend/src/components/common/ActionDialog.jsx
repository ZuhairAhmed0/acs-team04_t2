import { useRef } from 'react';
import ConfirmDialog from './ConfirmDialog';

/** Renders the confirmation for a `useConfirmedAction` controller. */
export default function ActionDialog({ controller }) {
  const { pending, cancel, confirm, submitting } = controller;
  const lastPending = useRef(null);
  if (pending) lastPending.current = pending; // keeps the text visible while the dialog fades out
  const { config, target } = pending ?? lastPending.current ?? {};

  return (
    <ConfirmDialog
      open={Boolean(pending)}
      title={config?.title}
      message={config?.message?.(target)}
      confirmLabel={config?.confirmLabel}
      tone={config?.tone}
      requireReason={config?.requireReason}
      loading={submitting}
      onConfirm={confirm}
      onClose={cancel}
    />
  );
}
