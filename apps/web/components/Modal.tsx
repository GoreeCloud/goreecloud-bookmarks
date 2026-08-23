import React, { MouseEventHandler, ReactNode } from "react";
import ReactDOM from "react-dom";
import ClickAwayHandler from "@/components/ClickAwayHandler";
import { Drawer } from "vaul";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import useBodyScrollLock from "@/hooks/useBodyScrollLock";
import { Button } from "@/components/ui/button";

type Props = {
  toggleModal: Function;
  children: ReactNode;
  className?: string;
  dismissible?: boolean;
  hideCloseButton?: boolean;
};

export default function Modal({
  toggleModal,
  className,
  children,
  dismissible = true,
  hideCloseButton,
}: Props) {
  const [drawerIsOpen, setDrawerIsOpen] = React.useState(true);
  const { width } = useWindowDimensions();

  useBodyScrollLock(width >= 640);

  if (width < 640) {
    return (
      <Drawer.Root
        open={drawerIsOpen}
        onClose={() => dismissible && setDrawerIsOpen(false)}
        onAnimationEnd={(isOpen) => !isOpen && toggleModal()}
        dismissible={dismissible}
      >
        <Drawer.Portal>
          <Drawer.Overlay className="fixed inset-0 z-40 bg-black/45 backdrop-blur-[2px]" />
          <Drawer.Content className="fixed bottom-0 left-0 right-0 z-50 mt-24 flex h-[90%] flex-col rounded-t-3xl outline-none">
            <div
              className="flex-1 overflow-y-auto rounded-t-3xl border-t border-base-300 bg-base-100 px-4 pb-5 pt-3 shadow-2xl"
              data-testid="mobile-modal-container"
            >
              <div
                className="mx-auto mb-5 h-1.5 w-12 flex-shrink-0 rounded-full bg-base-content/20"
                data-testid="mobile-modal-slider"
              />

              {children}
            </div>
          </Drawer.Content>
        </Drawer.Portal>
      </Drawer.Root>
    );
  }

  return ReactDOM.createPortal(
    <div
      className="fade-in fixed inset-0 z-40 flex items-center justify-center overflow-y-auto bg-black/35 px-4 py-6 backdrop-blur-sm"
      data-testid="modal-outer"
    >
      <ClickAwayHandler
        onClickOutside={() => dismissible && toggleModal()}
        className={`m-auto w-full sm:w-11/12 sm:max-w-xl ${className || ""}`}
      >
        <div
          className="slide-up relative m-auto overflow-y-auto rounded-2xl border border-base-300 bg-base-100 p-5 shadow-2xl sm:overflow-y-visible"
          data-testid="modal-container"
        >
          {dismissible && !hideCloseButton && (
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleModal as MouseEventHandler<HTMLButtonElement>}
              className="absolute right-3 top-3 z-10 h-10 w-10 rounded-full text-base-content/60 hover:bg-base-200 hover:text-base-content"
              aria-label="Close dialog"
            >
              <i
                className="bi-x text-xl"
                data-testid="close-modal-button"
                aria-hidden="true"
              />
            </Button>
          )}
          {children}
        </div>
      </ClickAwayHandler>
    </div>,
    document.body
  );
}
