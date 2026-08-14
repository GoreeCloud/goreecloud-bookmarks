import React, { useLayoutEffect, useRef, useState } from "react";
import TextInput from "@/components/TextInput";
import { TokenExpiry } from "@linkwarden/types/global";
import toast from "react-hot-toast";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { useAddToken } from "@linkwarden/router/tokens";
import CopyButton from "../CopyButton";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

export default function NewTokenModal({ onClose }: Props) {
  const { t } = useTranslation();
  const [newToken, setNewToken] = useState("");
  const addToken = useAddToken();

  const initial = {
    name: "",
    expires: TokenExpiry.sevenDays,
  };

  const [token, setToken] = useState(initial);
  const [submitLoader, setSubmitLoader] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const normalizedName = token.name.trim();

  const submit = async () => {
    if (submitLoader || !normalizedName) return;

    setSubmitLoader(true);
    const load = toast.loading(t("creating_token"));

    await addToken.mutateAsync(
      { ...token, name: normalizedName },
      {
        onSettled: (data, error) => {
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            setNewToken(data.secretKey);
          }
        },
      }
    );
  };

  const getLabel = (expiry: TokenExpiry) => {
    switch (expiry) {
      case TokenExpiry.sevenDays:
        return t("7_days");
      case TokenExpiry.oneMonth:
        return t("30_days");
      case TokenExpiry.twoMonths:
        return t("60_days");
      case TokenExpiry.threeMonths:
        return t("90_days");
      case TokenExpiry.never:
        return t("no_expiration");
    }
  };

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <Modal toggleModal={onClose}>
      {newToken ? (
        <GlazeModalFrame
          title={t("access_token_created")}
          description={t("token_creation_notice")}
          icon="bi-key"
        >
          <div
            role="note"
            className="rounded-xl border border-warning/25 bg-warning/10 p-3"
          >
            <div className="flex items-center gap-2">
              <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap rounded-lg border border-base-300 bg-base-100 px-3 py-2 font-mono text-xs text-base-content">
                {newToken}
              </code>
              <CopyButton text={newToken} />
            </div>
          </div>
        </GlazeModalFrame>
      ) : (
        <GlazeModalFrame
          title={t("create_access_token")}
          icon="bi-key"
          footer={
            <>
              <Button
                type="button"
                variant="ghost"
                onClick={() => onClose()}
                disabled={submitLoader}
              >
                {t("cancel")}
              </Button>
              <Button
                type="button"
                variant="primary"
                onClick={submit}
                disabled={submitLoader || !normalizedName}
              >
                <i className="bi-key" aria-hidden="true" />
                {t("create_token")}
              </Button>
            </>
          }
        >
          <div className="grid gap-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-end">
            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <p className="mb-2 text-xs font-medium text-base-content/60">
                {t("name")}
              </p>
              <TextInput
                ref={inputRef}
                value={token.name}
                onChange={(e) => setToken({ ...token, name: e.target.value })}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    submit();
                  }
                }}
                placeholder={t("token_name_placeholder")}
                className="bg-base-100"
              />
            </div>

            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <p className="mb-2 text-xs font-medium text-base-content/60">
                {t("expires_in")}
              </p>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    type="button"
                    variant="metal"
                    className="w-full whitespace-nowrap sm:w-32"
                  >
                    {getLabel(token.expires)}
                    <i className="bi-chevron-down text-xs" aria-hidden="true" />
                  </Button>
                </DropdownMenuTrigger>

                <DropdownMenuContent align="end">
                  <DropdownMenuRadioGroup
                    value={token.expires.toString()}
                    onValueChange={(value) =>
                      setToken({
                        ...token,
                        expires: Number(value) as TokenExpiry,
                      })
                    }
                  >
                    <DropdownMenuRadioItem
                      value={TokenExpiry.sevenDays.toString()}
                    >
                      {t("7_days")}
                    </DropdownMenuRadioItem>
                    <DropdownMenuRadioItem
                      value={TokenExpiry.oneMonth.toString()}
                    >
                      {t("30_days")}
                    </DropdownMenuRadioItem>
                    <DropdownMenuRadioItem
                      value={TokenExpiry.twoMonths.toString()}
                    >
                      {t("60_days")}
                    </DropdownMenuRadioItem>
                    <DropdownMenuRadioItem
                      value={TokenExpiry.threeMonths.toString()}
                    >
                      {t("90_days")}
                    </DropdownMenuRadioItem>
                    <DropdownMenuRadioItem value={TokenExpiry.never.toString()}>
                      {t("no_expiration")}
                    </DropdownMenuRadioItem>
                  </DropdownMenuRadioGroup>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </GlazeModalFrame>
      )}
    </Modal>
  );
}
