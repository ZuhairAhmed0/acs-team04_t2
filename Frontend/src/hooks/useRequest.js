import { useCallback, useEffect, useRef, useState } from 'react';
import { getErrorMessage } from '../utils/errors';

/**
 * Runs `requestFn` on mount and whenever `deps` change.
 * Keeps the previous data while reloading so lists do not flash.
 */
export default function useRequest(requestFn, deps = []) {
  const [state, setState] = useState({ data: null, loading: true, error: null });
  const latestCall = useRef(0);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const reload = useCallback(async () => {
    const call = ++latestCall.current;
    setState((prev) => ({ ...prev, loading: true, error: null }));
    try {
      const data = await requestFn();
      if (call === latestCall.current) setState({ data, loading: false, error: null });
    } catch (error) {
      if (call === latestCall.current) {
        setState({ data: null, loading: false, error: getErrorMessage(error) });
      }
    }
  }, deps);

  useEffect(() => {
    reload();
  }, [reload]);

  return { ...state, reload };
}
