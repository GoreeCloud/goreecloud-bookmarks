import React, { ReactNode, useState } from "react";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import Modal from "./Modal";
import { Separator } from "./ui/separator";

type Props = {
  toggleModal: () => void;
  className?: string;
  children: ReactNode;
  title: string;
  onConfirmed: () => void | Promise<void>;
  dismissible?: boolean;
};

export default function ConfirmationModal({
  toggleModal,
  className,
  children,
  title,
  onConfirmed,
  dismissible = true,
}: Props) {
  const { t } = useTranslation();
  const [isConfirming, setIsConfirming] = useState(false);

  const handleConfirm = async () => {
    if (isConfirming) return;

    setIsConfirming(true);

    try {
      await onConfirmed();
      toggleModal();
    } finally {
      setIsConfirming(false);
    }
  };

  return (
    <Modal
      toggleModal={toggleModal}
      className={className}
      dismissible={dismissible && !isConfirming}
    >
      <div aria-busy={isConfirming}>
        <h2 className="pr-12 text-xl font-semibold tracking-tight text-base-content">
          {title}
        </h2>
        <Separator className="mb-4 mt-2" />

        <div className="text-sm leading-6 text-base-content/80">{children}</div>

        <div className="mt-5 flex w-full flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-end">
          <Button
            variant="ghost"
            className="w-full hover:bg-base-200 sm:w-auto"
            onClick={toggleModal}
            disabled={isConfirming}
          >
            {t("cancel")}
          </Button>
          <Button
            variant="destructive"
            className="w-full sm:w-auto"
            onClick={handleConfirm}
            disabled={isConfirming}
          >
            {isConfirming && (
              <i className="bi-arrow-repeat animate-spin" aria-hidden="true" />
            )}
            {t("confirm")}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
