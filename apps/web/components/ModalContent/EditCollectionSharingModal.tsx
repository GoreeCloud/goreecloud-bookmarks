import React, { useEffect, useState } from "react";
import TextInput from "@/components/TextInput";
import toast from "react-hot-toast";
import {
  AccountSettings,
  CollectionIncludingMembersAndLinkCount,
  Member,
} from "@linkwarden/types/global";
import getPublicUserData from "@/lib/client/getPublicUserData";
import usePermissions from "@/hooks/usePermissions";
import ProfilePhoto from "../ProfilePhoto";
import addMemberToCollection from "@/lib/client/addMemberToCollection";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useUpdateCollection } from "@linkwarden/router/collections";
import { useUser } from "@linkwarden/router/user";
import CopyButton from "../CopyButton";
import { useRouter } from "next/router";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  activeCollection: CollectionIncludingMembersAndLinkCount;
};

export default function EditCollectionSharingModal({
  onClose,
  activeCollection,
}: Props) {
  const { t } = useTranslation();
  const [collection, setCollection] =
    useState<CollectionIncludingMembersAndLinkCount>(activeCollection);
  const [propagateToSubcollections, setPropagateToSubcollections] =
    useState(false);
  const [submitLoader, setSubmitLoader] = useState(false);
  const [memberIdentifier, setMemberIdentifier] = useState("");
  const [collectionOwner, setCollectionOwner] = useState<
    Partial<AccountSettings>
  >({});

  const updateCollection = useUpdateCollection();
  const { data: user } = useUser();
  const permissions = usePermissions(collection.id as number);
  const router = useRouter();
  const isPublicRoute = router.pathname.startsWith("/public");
  const canManageSharing = permissions === true && !isPublicRoute;

  const currentURL = new URL(document.URL);
  const publicCollectionURL = `${currentURL.origin}/public/collections/${collection.id}`;

  useEffect(() => {
    const fetchOwner = async () => {
      const owner = await getPublicUserData(collection.ownerId as number);
      setCollectionOwner(owner);
    };

    fetchOwner();
    setCollection(activeCollection);
  }, []);

  const submit = async () => {
    if (submitLoader || !collection) return;

    setSubmitLoader(true);
    const load = toast.loading(t("updating_collection"));

    await updateCollection.mutateAsync(
      { ...collection, propagateToSubcollections },
      {
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
      }
    );
  };

  const setMemberState = (newMember: Member) => {
    if (!collection) return null;

    setCollection({
      ...collection,
      members: [...collection.members, newMember],
    });
    setMemberIdentifier("");
  };

  const addMember = () => {
    const identifier = memberIdentifier.trim().replace(/^@/, "");
    if (!identifier) return;

    addMemberToCollection(user as any, identifier, collection, setMemberState, t);
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={canManageSharing ? t("share_and_collaborate") : t("team")}
        icon="bi-people"
        footer={
          canManageSharing ? (
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
                disabled={submitLoader}
              >
                <i className="bi-check-lg" aria-hidden="true" />
                {t("save_changes")}
              </Button>
            </>
          ) : undefined
        }
      >
        {canManageSharing && (
          <section className="rounded-xl border border-base-300 bg-base-200/40 p-3">
            <div className="flex items-start gap-3">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-base-300 bg-base-100 text-primary">
                <i className="bi-globe2" aria-hidden="true" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="font-medium text-base-content">
                  {t("make_collection_public")}
                </p>
                <p className="mt-0.5 text-xs leading-5 text-base-content/60">
                  {t("make_collection_public_desc")}
                </p>
              </div>
              <input
                type="checkbox"
                checked={collection.isPublic}
                onChange={() =>
                  setCollection({
                    ...collection,
                    isPublic: !collection.isPublic,
                  })
                }
                className="checkbox checkbox-primary mt-1"
                aria-label={t("make_collection_public_checkbox")}
              />
            </div>
          </section>
        )}

        {collection.isPublic && (
          <section>
            <p className="mb-2 text-xs font-medium text-base-content/60">
              {t("sharable_link")}
            </p>
            <div className="flex w-full items-center justify-between gap-2 overflow-x-auto whitespace-nowrap rounded-xl border border-base-300 bg-base-200/50 p-2 pl-3">
              <span className="min-w-0 overflow-hidden text-ellipsis text-xs">
                {publicCollectionURL}
              </span>
              <CopyButton text={publicCollectionURL} />
            </div>
          </section>
        )}

        {canManageSharing && (
          <section className="rounded-xl border border-base-300 p-3">
            <div className="mb-3">
              <p className="font-medium text-base-content">{t("members")}</p>
              <p className="mt-0.5 text-xs text-base-content/60">
                {t("add_member_placeholder")}
              </p>
            </div>

            <div className="flex items-center gap-2">
              <TextInput
                value={memberIdentifier}
                className="bg-base-200"
                placeholder={t("add_member_placeholder")}
                onChange={(e) => setMemberIdentifier(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    addMember();
                  }
                }}
              />

              <Button
                type="button"
                variant="primary"
                size="icon"
                className="h-10 w-10 shrink-0"
                onClick={addMember}
                disabled={!memberIdentifier.trim()}
                aria-label={t("members")}
              >
                <i className="bi-person-add text-lg" aria-hidden="true" />
              </Button>
            </div>
          </section>
        )}

        {collection?.members[0]?.user && (
          <section>
            <p className="mb-2 text-xs font-medium text-base-content/60">
              {t("team")}
            </p>
            <div className="overflow-hidden rounded-xl border border-base-300 bg-base-100">
              <div
                className="flex items-center justify-between gap-3 bg-base-200/50 p-3"
                title={`@${collectionOwner.username} is the owner of this collection`}
              >
                <div className="flex min-w-0 items-center">
                  <div className="shrink-0">
                    <ProfilePhoto
                      src={collectionOwner.image || undefined}
                      name={collectionOwner.name}
                    />
                  </div>
                  <div className="ml-2 min-w-0">
                    <p className="truncate text-sm font-semibold text-base-content">
                      {collectionOwner.name}
                    </p>
                    <p className="truncate text-xs text-base-content/60">
                      @{collectionOwner.username}
                    </p>
                  </div>
                </div>
                <span className="rounded-full border border-base-300 bg-base-100 px-2.5 py-1 text-xs font-medium">
                  {t("owner")}
                </span>
              </div>

              <div className="divide-y divide-base-300">
                {[...collection.members]
                  .sort((a, b) => (a.userId as number) - (b.userId as number))
                  .map((member) => {
                    const roleKey: "viewer" | "contributor" | "admin" =
                      !member.canCreate &&
                      !member.canUpdate &&
                      !member.canDelete
                        ? "viewer"
                        : member.canCreate &&
                            !member.canUpdate &&
                            !member.canDelete
                          ? "contributor"
                          : "admin";

                    const handleRoleChange = (newRole: string) => {
                      const updatedMember = {
                        ...member,
                        canCreate: newRole !== "viewer",
                        canUpdate: newRole === "admin",
                        canDelete: newRole === "admin",
                      };

                      setCollection({
                        ...collection,
                        members: collection.members.map((candidate) =>
                          candidate.userId === member.userId
                            ? updatedMember
                            : candidate
                        ),
                      });
                    };

                    return (
                      <div
                        key={member.userId}
                        className="flex items-center justify-between gap-3 p-3"
                      >
                        <div className="flex min-w-0 items-center">
                          <ProfilePhoto
                            src={member.user.image || undefined}
                            name={member.user.name}
                          />
                          <div className="ml-2 min-w-0">
                            <p className="truncate text-sm font-semibold text-base-content">
                              {member.user.name}
                            </p>
                            <p className="truncate text-xs text-base-content/60">
                              @{member.user.username}
                            </p>
                          </div>
                        </div>

                        <div className="flex shrink-0 items-center gap-1">
                          {canManageSharing ? (
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button
                                  type="button"
                                  variant="ghost"
                                  className="h-9 rounded-lg px-2.5 text-xs"
                                >
                                  {t(roleKey)}
                                  <i
                                    className="bi-chevron-down text-xs"
                                    aria-hidden="true"
                                  />
                                </Button>
                              </DropdownMenuTrigger>

                              <DropdownMenuContent sideOffset={4} align="end">
                                <DropdownMenuRadioGroup
                                  value={roleKey}
                                  onValueChange={handleRoleChange}
                                >
                                  <DropdownMenuRadioItem value="viewer">
                                    <div>
                                      <p className="whitespace-nowrap font-bold">
                                        {t("viewer")}
                                      </p>
                                      <p className="whitespace-nowrap">
                                        {t("viewer_desc")}
                                      </p>
                                    </div>
                                  </DropdownMenuRadioItem>

                                  <DropdownMenuRadioItem value="contributor">
                                    <div>
                                      <p className="whitespace-nowrap font-bold">
                                        {t("contributor")}
                                      </p>
                                      <p className="whitespace-nowrap">
                                        {t("contributor_desc")}
                                      </p>
                                    </div>
                                  </DropdownMenuRadioItem>

                                  <DropdownMenuRadioItem value="admin">
                                    <div>
                                      <p className="whitespace-nowrap font-bold">
                                        {t("admin")}
                                      </p>
                                      <p className="whitespace-nowrap">
                                        {t("admin_desc")}
                                      </p>
                                    </div>
                                  </DropdownMenuRadioItem>
                                </DropdownMenuRadioGroup>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          ) : (
                            <span className="px-2 text-xs text-base-content/60">
                              {t(roleKey)}
                            </span>
                          )}

                          {canManageSharing && (
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="h-9 w-9 rounded-lg text-base-content/50 hover:bg-error/10 hover:text-error"
                              onClick={() => {
                                setCollection({
                                  ...collection,
                                  members: collection.members.filter(
                                    (candidate) =>
                                      candidate.userId !== member.userId
                                  ),
                                });
                              }}
                              aria-label={t("remove_member")}
                            >
                              <i className="bi-x text-lg" aria-hidden="true" />
                            </Button>
                          )}
                        </div>
                      </div>
                    );
                  })}
              </div>
            </div>
          </section>
        )}

        {canManageSharing && (
          <label className="flex cursor-pointer items-start gap-3 rounded-xl border border-base-300 bg-base-200/40 p-3">
            <input
              type="checkbox"
              checked={propagateToSubcollections}
              onChange={() =>
                setPropagateToSubcollections(!propagateToSubcollections)
              }
              className="checkbox checkbox-primary mt-0.5"
            />
            <span className="min-w-0">
              <span className="block font-medium text-base-content">
                {t("apply_members_roles_to_subcollections")}
              </span>
              <span className="mt-0.5 block text-xs leading-5 text-base-content/60">
                {t("apply_members_roles_to_subcollections_desc")}
              </span>
            </span>
          </label>
        )}
      </GlazeModalFrame>
    </Modal>
  );
}
