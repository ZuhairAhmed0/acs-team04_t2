import { useState } from 'react';
import { useToast } from '../context/ToastContext';
import useSubmit from './useSubmit';

/**
 * "Ask for confirmation, run the request, toast, refresh" flow used by the admin screens.
 * `config` = { title, message(target), confirmLabel, tone, requireReason, success, run(target, reason) }.
 */
export default function useConfirmedAction(onDone) {
  const toast = useToast();
  const [pending, setPending] = useState(null);
  const { submit, submitting } = useSubmit(async (reason) => {
    await pending.config.run(pending.target, reason);
    return true;
  });

  const request = (config, target) => setPending({ config, target });
  const cancel = () => setPending(null);
  const confirm = async (reason) => {
    const succeeded = await submit(reason);
    if (!succeeded) return;
    toast.success(pending.config.success);
    setPending(null);
    onDone();
  };

  return { pending, request, cancel, confirm, submitting };
}
