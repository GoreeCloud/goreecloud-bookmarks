import React, { useEffect, useMemo, useState } from "react";
import Tree, {
  mutateTree,
  moveItemOnTree,
  RenderItemParams,
  TreeItem,
  TreeData,
  ItemId,
  TreeSourcePosition,
  TreeDestinationPosition,
} from "@atlaskit/tree";
import { Collection } from "@linkwarden/prisma/client";
import Link from "next/link";
import { CollectionIncludingMembersAndLinkCount } from "@linkwarden/types/global";
import { useRouter } from "next/router";
import toast from "react-hot-toast";
import { useTranslation } from "next-i18next";
import {
  useCollections,
  useUpdateCollection,
} from "@linkwarden/router/collections";
import { useUpdateUser, useUser } from "@linkwarden/router/user";
import Icon from "./Icon";
import { IconWeight } from "@phosphor-icons/react";
import Droppable from "./Droppable";
import { cn } from "@linkwarden/lib/utils";
import { Active, useDndContext } from "@dnd-kit/core";

interface ExtendedTreeItem extends TreeItem {
  data: Collection;
}

const CollectionListing = () => {
  const { active: droppableActive } = useDndContext();
  const { t } = useTranslation();
  const updateCollection = useUpdateCollection();
  const { data: collections = [], isLoading } = useCollections();
  const { data: user } = useUser();
  const updateUser = useUpdateUser();
  const router = useRouter();
  const [tree, setTree] = useState<TreeData | undefined>();

  const initialTree = useMemo(() => {
    if (collections.length > 0) {
      return buildTreeFromCollections(
        collections,
        router,
        tree,
        user?.collectionOrder
      );
    }
    return undefined;
  }, [collections, user]);

  useEffect(() => {
    setTree(initialTree);
  }, [initialTree]);

  useEffect(() => {
    if (user?.username) {
      if (
        (!user.collectionOrder || user.collectionOrder.length === 0) &&
        collections.length > 0
      ) {
        updateUser.mutate({
          ...user,
          collectionOrder: collections
            .filter((collection) => collection.parentId === null)
            .map((collection) => collection.id as number),
        });
      } else {
        const newCollectionOrder: number[] = [...(user.collectionOrder || [])];
        const existingCollectionIds = collections.map(
          (collection) => collection.id as number
        );
        const filteredCollectionOrder = user.collectionOrder.filter((id: any) =>
          existingCollectionIds.includes(id)
        );

        collections.forEach((collection) => {
          if (
            !filteredCollectionOrder.includes(collection.id as number) &&
            (!collection.parentId || collection.ownerId === user.id)
          ) {
            filteredCollectionOrder.push(collection.id as number);
          }
        });

        if (
          JSON.stringify(newCollectionOrder) !==
          JSON.stringify(user.collectionOrder)
        ) {
          updateUser.mutateAsync({
            ...user,
            collectionOrder: newCollectionOrder,
          });
        }
      }
    }
  }, [user, collections]);

  const onExpand = (movedCollectionId: ItemId) => {
    setTree((currentTree) =>
      mutateTree(currentTree!, movedCollectionId, { isExpanded: true })
    );
  };

  const onCollapse = (movedCollectionId: ItemId) => {
    setTree((currentTree) =>
      mutateTree(currentTree as TreeData, movedCollectionId, {
        isExpanded: false,
      })
    );
  };

  function reorderTreeItems(
    treeData: TreeData,
    movedCollectionId: ItemId,
    source: TreeSourcePosition,
    destination: TreeDestinationPosition
  ) {
    if (source.parentId === destination.parentId) {
      const parent = treeData.items[source.parentId];
      const children = [...parent.children];
      children.splice(source.index, 1);
      if (destination.index !== undefined) {
        children.splice(destination.index, 0, movedCollectionId);
      }
      parent.children = children;
      return treeData;
    }

    const sourceParent = treeData.items[source.parentId];
    const destinationParent = treeData.items[destination.parentId];

    sourceParent.children = sourceParent.children.filter(
      (id) => id !== movedCollectionId
    );

    if (!destinationParent.children) {
      destinationParent.children = [];
    }

    const destinationIndex =
      destination.index !== undefined
        ? destination.index
        : destinationParent.children.length;

    destinationParent.children.splice(destinationIndex, 0, movedCollectionId);
    destinationParent.hasChildren = true;
    destinationParent.isExpanded = true;
    treeData.items[movedCollectionId].data.parentId = destination.parentId;

    return treeData;
  }

  function flattenTreeIds(
    treeData: TreeData,
    nodeId: ItemId = "root",
    result: Array<ItemId> = []
  ) {
    const node = treeData.items[nodeId];

    if (nodeId !== "root") {
      result.push(node.id);
    }

    if (node.children && node.children.length > 0) {
      node.children.forEach((childId) => {
        flattenTreeIds(treeData, childId, result);
      });
    }

    return result;
  }

  const onDragEnd = async (
    source: TreeSourcePosition,
    destination: TreeDestinationPosition | undefined
  ) => {
    if (!destination || !tree) return;

    if (
      source.index === destination.index &&
      source.parentId === destination.parentId
    ) {
      return;
    }

    const movedCollectionId = Number(
      tree.items[source.parentId].children[source.index]
    );
    const movedCollection = collections.find(
      (collection) => collection.id === movedCollectionId
    );
    const destinationCollection = collections.find(
      (collection) => collection.id === Number(destination.parentId)
    );

    if (
      (movedCollection?.ownerId !== user?.id &&
        destination.parentId !== source.parentId) ||
      (destinationCollection?.ownerId !== user?.id &&
        destination.parentId !== "root")
    ) {
      return toast.error(t("cant_change_collection_you_dont_own"));
    }

    setTree((currentTree) => moveItemOnTree(currentTree!, source, destination));

    const newTree = reorderTreeItems(
      tree,
      movedCollectionId,
      source,
      destination
    );

    if (source.parentId !== destination.parentId) {
      await updateCollection.mutateAsync(
        {
          ...movedCollection,
          parentId:
            destination.parentId && destination.parentId !== "root"
              ? Number(destination.parentId)
              : destination.parentId === "root"
                ? "root"
                : null,
        },
        {
          onError: (error) => {
            toast.error(error.message);
          },
        }
      );
    }

    await updateUser.mutateAsync({
      ...user,
      collectionOrder: flattenTreeIds(newTree),
    });
  };

  if (isLoading) {
    return (
      <div className="flex flex-col gap-2 px-1">
        <div className="skeleton h-10 w-full rounded-lg" />
        <div className="skeleton h-10 w-11/12 rounded-lg" />
        <div className="skeleton h-10 w-full rounded-lg" />
      </div>
    );
  }

  if (!tree) {
    return (
      <p className="w-full truncate rounded-lg px-2 py-2 text-xs font-medium text-neutral">
        {t("you_have_no_collections")}
      </p>
    );
  }

  return (
    <Tree
      tree={tree}
      renderItem={(itemProps) =>
        renderItem({ ...itemProps }, router.asPath, droppableActive)
      }
      onExpand={onExpand}
      onCollapse={onCollapse}
      onDragEnd={onDragEnd}
      isDragEnabled
      isNestingEnabled
    />
  );
};

export default CollectionListing;

const renderItem = (
  { item, onExpand, onCollapse, provided }: RenderItemParams,
  currentPath: string,
  droppableActive: Active | null
) => {
  const collection = item.data;
  const isActive = currentPath === `/collections/${collection.id}`;

  return (
    <Droppable
      id={`side-bar-collection-${collection.id}`}
      data={{
        name: collection.name,
        id: collection.id,
        ownerId: collection.ownerId,
      }}
      className="group"
    >
      <div ref={provided.innerRef} {...provided.draggableProps} className="mb-0.5">
        <div
          className={cn(
            isActive
              ? "bg-primary/15 is-active"
              : droppableActive
                ? "select-none"
                : "hover:bg-base-content/5",
            "relative flex min-h-10 items-center rounded-lg pr-1.5 transition-colors duration-100"
          )}
        >
          {Dropdown(item as ExtendedTreeItem, onExpand, onCollapse)}

          <Link
            href={`/collections/${collection.id}`}
            className="min-w-0 flex-1 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
            {...provided.dragHandleProps}
          >
            <div className="flex min-w-0 items-center gap-2 py-2 pl-8">
              <span className="flex h-5 w-5 shrink-0 items-center justify-center">
                {collection.icon ? (
                  <Icon
                    icon={collection.icon}
                    size={18}
                    weight={(collection.iconWeight || "regular") as IconWeight}
                    color={collection.color}
                  />
                ) : (
                  <i
                    className="bi-folder-fill text-sm"
                    style={{ color: collection.color }}
                    aria-hidden="true"
                  />
                )}
              </span>

              <p
                className={cn(
                  "min-w-0 flex-1 truncate text-xs",
                  isActive ? "font-semibold text-base-content" : "font-medium"
                )}
              >
                {collection.name}
              </p>

              {collection.isPublic && (
                <i
                  className="bi-globe2 shrink-0 text-[11px] text-neutral"
                  title="This collection is being shared publicly."
                />
              )}

              <span className="min-w-5 shrink-0 rounded-full bg-base-200 px-1.5 py-0.5 text-center text-[10px] tabular-nums text-neutral">
                {collection._count?.links || 0}
              </span>
            </div>
          </Link>
        </div>
      </div>
    </Droppable>
  );
};

const Dropdown = (
  item: ExtendedTreeItem,
  onExpand: (id: ItemId) => void,
  onCollapse: (id: ItemId) => void
) => {
  if (!item.children || item.children.length === 0) return null;

  const expanded = item.isExpanded;

  return (
    <button
      type="button"
      onClick={() => (expanded ? onCollapse(item.id) : onExpand(item.id))}
      className="absolute left-0 top-1/2 z-10 flex h-9 w-8 -translate-y-1/2 items-center justify-center rounded-lg text-neutral opacity-70 transition hover:bg-base-content/10 hover:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
      aria-label={expanded ? "Collapse collection" : "Expand collection"}
    >
      <i
        className={
          expanded ? "bi-chevron-down text-[10px]" : "bi-chevron-right text-[10px]"
        }
        aria-hidden="true"
      />
    </button>
  );
};

const buildTreeFromCollections = (
  collections: CollectionIncludingMembersAndLinkCount[],
  router: ReturnType<typeof useRouter>,
  tree?: TreeData,
  order?: number[]
): TreeData => {
  if (order) {
    collections.sort((a: any, b: any) => {
      return order.indexOf(a.id) - order.indexOf(b.id);
    });
  }

  function getTotalLinkCount(collectionId: number): number {
    const collection = items[collectionId];
    if (!collection) return 0;

    let totalLinkCount = (collection.data as any)._count?.links || 0;

    if (collection.hasChildren) {
      collection.children.forEach((childId) => {
        totalLinkCount += getTotalLinkCount(childId as number);
      });
    }

    return totalLinkCount;
  }

  const items: { [key: string]: ExtendedTreeItem } = collections.reduce(
    (acc: any, collection) => {
      acc[collection.id as number] = {
        id: collection.id,
        children: [],
        hasChildren: false,
        isExpanded: tree?.items[collection.id as number]?.isExpanded || false,
        data: {
          id: collection.id,
          parentId: collection.parentId,
          name: collection.name,
          description: collection.description,
          color: collection.color,
          icon: collection.icon,
          iconWeight: collection.iconWeight,
          isPublic: collection.isPublic,
          ownerId: collection.ownerId,
          createdAt: collection.createdAt,
          updatedAt: collection.updatedAt,
          _count: {
            links: collection._count?.links,
          },
        },
      };
      return acc;
    },
    {}
  );

  const activeCollectionId = Number(router.asPath.split("/collections/")[1]);

  if (activeCollectionId) {
    for (const item in items) {
      const collection = items[item];
      if (Number(item) === activeCollectionId && collection.data.parentId) {
        let parentId = collection.data.parentId || null;
        while (parentId && items[parentId]) {
          items[parentId].isExpanded = true;
          parentId = items[parentId].data.parentId;
        }
      }
    }
  }

  collections.forEach((collection) => {
    const parentId = collection.parentId;
    if (parentId && items[parentId] && collection.id) {
      items[parentId].children.push(collection.id);
      items[parentId].hasChildren = true;
    }
  });

  collections.forEach((collection) => {
    const collectionId = collection.id;
    if (items[collectionId as number] && collection.id) {
      const linkCount = getTotalLinkCount(collectionId as number);
      (items[collectionId as number].data as any)._count.links = linkCount;
    }
  });

  const rootId = "root";
  items[rootId] = {
    id: rootId,
    children: (collections
      .filter(
        (collection) =>
          collection.parentId === null ||
          !collections.find((item) => item.id === collection.parentId)
      )
      .map((collection) => collection.id) || "") as unknown as string[],
    hasChildren: true,
    isExpanded: true,
    data: { name: "Root" } as Collection,
  };

  return { rootId, items };
};
