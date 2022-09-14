import { useCallback, useEffect, useRef } from 'react';

const timerIntervalMSeconds = 1000 * 60;

type Callback = () => void;

export function useTimer(timeoutMinutes: number, callback: Callback): void {
  const isBlurNow = useRef(false);
  const dueTime = useRef(Date.now() + timeoutMinutes * 60 * 1000);

  const handleFocus = useCallback(() => {
    isBlurNow.current = false;
    // if focus, check immediately.
    checkTimeOut();
  }, [callback]);

  const handleBlur = () => {
    isBlurNow.current = true;
  };

  const handleInterval = () => {
    checkTimeOut();
  };

  const checkTimeOut = () => {
    const now = Date.now()
    if (now >= dueTime.current) {
        // check current is focused or not
        if (!isBlurNow.current) {
            callback();
            // reset
            dueTime.current = now + (timeoutMinutes * 60 * 1000)
        }
    }
  }

  useEffect(() => {
    const id = setInterval(handleInterval, timerIntervalMSeconds);
    return () => {
      clearInterval(id);
    };
  }, []);

  useEffect(() => {
    window.addEventListener('blur', handleBlur);
    return () => {
      window.removeEventListener('blur', handleBlur);
    };
  }, []);

  useEffect(() => {
    window.addEventListener('focus', handleFocus);
    return () => {
      window.removeEventListener('focus', handleFocus);
    };
  }, [handleFocus]);
}