import React, { useState } from "react";
import TextInput from "@/components/TextInput";
import { CollectionIncludingMembersAndLinkCount } from "@linkwarden/types/global";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useUpdateCollection } from "@linkwarden/router/collections";
import toast from "react-hot-toast";
import IconPicker from "../IconPicker";
import { IconWeight } from "@phosphor-icons/react";
import oklchVariableToHex from "@/lib/client/oklchVariableToHex";
import { Button } from "../ui/button";
import { Separator } from "../ui/separator";

type Props = {
  onClose: Function;
  activeCollection: CollectionIncludingMembersAndLinkCount;
};

export default function EditCollectionModal({
  onClose,
  activeCollection,
}: Props) {
  const { t } = useTranslation();
  const [collection, setCollection] =
    useState<CollectionIncludingMembersAndLinkCount>(activeCollection);
  const [submitLoader, setSubmitLoader] = useState(false);
  const updateCollection = useUpdateCollection();

  const submit = async () => {
    if (submitLoader || !collection.name.trim()) return;

    setSubmitLoader(true);
    const load = toast.loading(t("updating_collection"));

    await updateCollection.mutateAsync(collection, {
      onSettled: (data, error) => {
        setSubmitLoader(false);
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          onClose();
          toast.success(t("updated"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <div className="flex items-start gap-3 pr-9">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200 text-primary">
          <i className="bi-folder-gear" />
        </div>
        <div>
          <h2 className="text-lg font-semibold">{t("edit_collection_info")}</h2>
          <p className="mt-1 text-xs text-neutral">{collection.name}</p>
        </div>
      </div>

      <Separator className="my-4" />

      <div className="flex flex-col gap-5">
        <div className="grid gap-4 sm:grid-cols-[auto,1fr] sm:items-end">
          <div>
            <p className="mb-2 text-xs font-medium text-neutral">{t("icon")}</p>
            <div className="rounded-xl border border-neutral-content bg-base-200/50 p-1">
              <IconPicker
                color={collection.color}
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
              className="h-11 bg-base-200"
              value={collection.name}
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
            disabled={submitLoader || !collection.name.trim()}
          >
            <i className="bi-check2" />
            {t("save_changes")}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
