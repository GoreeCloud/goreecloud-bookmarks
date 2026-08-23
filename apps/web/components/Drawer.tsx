import React, { ReactNode } from "react";
import { Drawer as D } from "vaul";
import clsx from "clsx";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import useBodyScrollLock from "@/hooks/useBodyScrollLock";

type Props = {
  toggleDrawer: Function;
  children: ReactNode;
  className?: string;
  dismissible?: boolean;
  direction?: "left" | "right";
};

export default function Drawer({
  toggleDrawer,
  className,
  children,
  dismissible = true,
  direction,
}: Props) {
  const [drawerIsOpen, setDrawerIsOpen] = React.useState(true);
  const { width } = useWindowDimensions();

  useBodyScrollLock(width >= 640);

  if (width < 640) {
    return (
      <D.Root
        open={drawerIsOpen}
        onClose={() => dismissible && setDrawerIsOpen(false)}
        onAnimationEnd={(isOpen) => !isOpen && toggleDrawer()}
        dismissible={dismissible}
      >
        <D.Portal>
          <D.Overlay className="fixed inset-0 z-30 bg-black/45 backdrop-blur-[2px]" />
          <D.Content className="fixed bottom-0 left-0 right-0 z-40 mt-24 flex h-[90%] flex-col rounded-t-3xl !select-auto outline-none">
            <div
              className={clsx(
                "flex-1 overflow-y-auto rounded-t-3xl border-t border-base-300 bg-base-100 px-4 pb-5 pt-3 shadow-2xl",
                className
              )}
              data-testid="mobile-modal-container"
            >
              <div
                className="mx-auto mb-5 h-1.5 w-12 rounded-full bg-base-content/20"
                data-testid="mobile-modal-slider"
              />
              {children}
            </div>
          </D.Content>
        </D.Portal>
      </D.Root>
    );
  }

  const drawerDirection = direction || "right";

  return (
    <D.Root
      open={drawerIsOpen}
      onClose={() => dismissible && setDrawerIsOpen(false)}
      onAnimationEnd={(isOpen) => !isOpen && toggleDrawer()}
      dismissible={dismissible}
      direction={drawerDirection}
    >
      <D.Portal>
        <D.Overlay className="fixed inset-0 z-20 bg-black/30 backdrop-blur-[2px]" />
        <D.Content
          className={clsx(
            "fixed bottom-0 z-40 mt-24 flex h-full w-2/5 min-w-[30rem] max-w-6xl flex-col bg-base-100 !select-auto outline-none shadow-2xl",
            drawerDirection === "left"
              ? "left-0 rounded-r-2xl"
              : "right-0 rounded-l-2xl"
          )}
        >
          <div
            className={clsx(
              "flex-1 overflow-y-auto bg-base-100 p-4",
              drawerDirection === "left"
                ? "rounded-r-2xl border-r border-base-300"
                : "rounded-l-2xl border-l border-base-300",
              className
            )}
          >
            {children}
          </div>
        </D.Content>
      </D.Portal>
    </D.Root>
  );
}
