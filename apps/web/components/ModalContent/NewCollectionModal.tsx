import React, { useEffect, useLayoutEffect, useRef, useState } from "react";
import TextInput from "@/components/TextInput";
import { Collection } from "@linkwarden/prisma/client";
import Modal from "../Modal";
import { CollectionIncludingMembersAndLinkCount } from "@linkwarden/types/global";
import { useTranslation } from "next-i18next";
import { useCreateCollection } from "@linkwarden/router/collections";
import toast from "react-hot-toast";
import IconPicker from "../IconPicker";
import { IconWeight } from "@phosphor-icons/react";
import oklchVariableToHex from "@/lib/client/oklchVariableToHex";
import { Button } from "../ui/button";
import { Separator } from "../ui/separator";

type Props = {
  onClose: Function;
  parent?: CollectionIncludingMembersAndLinkCount;
};

export default function NewCollectionModal({ onClose, parent }: Props) {
  const { t } = useTranslation();

  const initial = {
    parentId: parent?.id,
    name: "",
    description: "",
    color: oklchVariableToHex("--p"),
  } as Partial<Collection>;

  const [collection, setCollection] = useState<Partial<Collection>>(initial);
  const [submitLoader, setSubmitLoader] = useState(false);
  const createCollection = useCreateCollection();
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setCollection(initial);
  }, [parent]);

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  const submit = async () => {
    if (submitLoader || !collection.name?.trim()) return;

    setSubmitLoader(true);
    const load = toast.loading(t("creating"));

    await createCollection.mutateAsync(collection, {
      onSettled: (data, error) => {
        setSubmitLoader(false);
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          onClose();
          toast.success(t("created"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <div className="flex items-start gap-3 pr-9">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200 text-primary">
          <i className={parent?.id ? "bi-folder-plus" : "bi-folder"} />
        </div>
        <div>
          <h2 className="text-lg font-semibold">
            {parent?.id ? t("new_sub_collection") : t("create_new_collection")}
          </h2>
          {parent?.id && (
            <p className="mt-1 text-xs text-neutral">
              {t("for_collection", { name: parent.name })}
            </p>
          )}
        </div>
      </div>

      <Separator className="my-4" />

      <div className="flex flex-col gap-5">
        <div className="grid gap-4 sm:grid-cols-[auto,1fr] sm:items-end">
          <div>
            <p className="mb-2 text-xs font-medium text-neutral">{t("icon")}</p>
            <div className="rounded-xl border border-neutral-content bg-base-200/50 p-1">
              <IconPicker
                color={collection.color || oklchVariableToHex("--p")}
                setColor={(color: string) =>
                  setCollection({ ...collection, color })
                }
                weight={(collection.iconWeight || "regular") as IconWeight}
                setWeight={(iconWeight: string) =>
                  setCollection({ ...collection, iconWeight })
                }
                iconName={collection.icon as string}
                setIconName={(icon: string) =>
                  setCollection({ ...collection, icon })
                }
                reset={() =>
                  setCollection({
                    ...collection,
                    color: oklchVariableToHex("--p"),
                    icon: "",
                    iconWeight: "",
                  })
                }
              />
            </div>
          </div>

          <div className="w-full">
            <p className="mb-2 text-xs font-medium text-neutral">{t("name")}</p>
            <TextInput
              ref={inputRef}
              className="h-11 bg-base-200"
              value={collection.name || ""}
              placeholder={t("collection_name_placeholder")}
              onChange={(e) =>
                setCollection({ ...collection, name: e.target.value })
              }
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  submit();
                }
              }}
            />
          </div>
        </div>

        <div className="w-full">
          <p className="mb-2 text-xs font-medium text-neutral">{t("description")}</p>
          <textarea
            className="h-28 w-full resize-none rounded-lg border border-neutral-content bg-base-200 p-3 text-sm outline-none duration-100 focus:border-primary"
            placeholder={t("collection_description_placeholder")}
            value={collection.description || ""}
            onChange={(e) =>
              setCollection({ ...collection, description: e.target.value })
            }
          />
        </div>

        <div className="flex items-center justify-end gap-2 border-t border-neutral-content pt-4">
          <Button type="button" variant="ghost" onClick={() => onClose()}>
            {t("cancel")}
          </Button>
          <Button
            type="button"
            variant="primary"
            onClick={submit}
            disabled={submitLoader || !collection.name?.trim()}
          >
            <i className="bi-plus-lg" />
            {t("create_collection_button")}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
