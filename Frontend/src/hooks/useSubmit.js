import { useState } from 'react';
import { useToast } from '../context/ToastContext';
import { getErrorMessage } from '../utils/errors';

/**
 * Wraps a mutating request: blocks double submissions and reports errors as a toast.
 * Resolves with the result, or `undefined` when the request failed.
 */
export default function useSubmit(action) {
  const toast = useToast();
  const [submitting, setSubmitting] = useState(false);

  const submit = async (...args) => {
    if (submitting) return undefined;
    setSubmitting(true);
    try {
      return await action(...args);
    } catch (error) {
      toast.error(getErrorMessage(error));
      return undefined;
    } finally {
      setSubmitting(false);
    }
  };

  return { submit, submitting };
}
