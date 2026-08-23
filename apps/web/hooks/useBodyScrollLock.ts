import { useEffect } from "react";

export default function useBodyScrollLock(locked: boolean) {
  useEffect(() => {
    if (!locked) return;

    const previousOverflow = document.body.style.overflow;
    const previousPosition = document.body.style.position;

    document.body.style.overflow = "hidden";
    document.body.style.position = "relative";

    return () => {
      document.body.style.overflow = previousOverflow;
      document.body.style.position = previousPosition;
    };
  }, [locked]);
}
