import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import { toast } from "react-hot-toast";
import { useTranslation } from "next-i18next";
import Modal from "@/components/Modal";
import { Button } from "@/components/ui/button";
import { Kbd } from "@/components/ui/kbd";
import { ADVANCED_SEARCH_OPERATORS } from "@/components/SearchBar";

export const isAppleDevice =
  typeof navigator !== "undefined" &&
  /Mac|iPhone|iPad|iPod/.test(navigator.platform);

const RECENT_SEARCHES_KEY = "recentSearches";
const MAX_RECENT_SEARCHES = 5;

type RecentSearch = {
  query: string;
  searchedAt: string;
};

const getRecentSearches = (): RecentSearch[] => {
  try {
    const stored = JSON.parse(
      localStorage.getItem(RECENT_SEARCHES_KEY) || "[]"
    );
    if (!Array.isArray(stored)) return [];
    return stored.filter((entry) => typeof entry?.query === "string");
  } catch {
    return [];
  }
};

const setRecentSearches = (recents: RecentSearch[]) => {
  localStorage.setItem(RECENT_SEARCHES_KEY, JSON.stringify(recents));
};

export default function SearchModal({ onClose }: { onClose: () => void }) {
  const { t } = useTranslation();
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [recents, setRecents] = useState<RecentSearch[]>(getRecentSearches);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", handler);
    return () => document.removeEventListener("keydown", handler);
  }, [onClose]);

  const submit = (query: string) => {
    const trimmed = query.trim();
    if (!trimmed) return;

    const updated = [
      { query: trimmed, searchedAt: new Date().toISOString() },
      ...recents.filter((entry) => entry.query !== trimmed),
    ].slice(0, MAX_RECENT_SEARCHES);

    setRecentSearches(updated);
    setRecents(updated);
    router.push("/search?q=" + encodeURIComponent(trimmed));
    onClose();
  };

  const removeRecent = (query: string) => {
    const updated = recents.filter((entry) => entry.query !== query);
    setRecentSearches(updated);
    setRecents(updated);
  };

  const appendOperator = (operator: string) => {
    setSearchQuery((prev) => {
      const needsSpace = prev.length > 0 && !prev.endsWith(" ");
      return `${prev}${needsSpace ? " " : ""}${operator}`;
    });
    requestAnimationFrame(() => inputRef.current?.focus());
  };

  return (
    <Modal toggleModal={onClose} hideCloseButton className="sm:!max-w-2xl">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          submit(searchQuery);
        }}
      >
        <div className="flex items-center gap-3 rounded-xl border border-base-content/10 bg-base-content/[0.025] px-3 py-2.5 focus-within:border-primary/35 focus-within:ring-2 focus-within:ring-primary/10">
          <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/[0.08] text-primary/80">
            <i className="bi-search text-sm leading-none" aria-hidden="true" />
          </span>

          <input
            type="search"
            ref={inputRef}
            autoFocus
            aria-label={t("search_for_links")}
            placeholder={t("search_for_links")}
            value={searchQuery}
            onChange={(e) => {
              e.target.value.includes("%") &&
                toast.error(t("search_query_invalid_symbol"));
              setSearchQuery(e.target.value.replace("%", ""));
            }}
            className="min-w-0 flex-1 bg-transparent text-base font-medium text-base-content outline-none placeholder:font-normal placeholder:text-base-content/35 sm:text-lg"
          />

          {searchQuery ? (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-10 w-10 shrink-0 rounded-lg text-base-content/40 hover:bg-base-content/[0.06] hover:text-base-content sm:h-8 sm:w-8"
              aria-label="Clear search"
              onClick={() => {
                setSearchQuery("");
                inputRef.current?.focus();
              }}
            >
              <i className="bi-x-lg text-xs" aria-hidden="true" />
            </Button>
          ) : (
            <div className="hidden shrink-0 items-center gap-1 sm:flex">
              <Kbd>{isAppleDevice ? "⌘" : "Ctrl"}</Kbd>
              <Kbd>K</Kbd>
            </div>
          )}
        </div>
      </form>

      <div className="mt-4 max-h-[28rem] overflow-y-auto pr-1">
        {recents.length > 0 && (
          <section>
            <div className="flex items-center justify-between px-1">
              <p className="text-[11px] font-semibold uppercase tracking-[0.08em] text-base-content/40">
                {t("recent")}
              </p>
            </div>

            <div className="mt-2 flex flex-col gap-1">
              {recents.map((entry) => (
                <div
                  key={entry.query}
                  className="group/recent flex min-h-11 cursor-pointer items-center gap-3 rounded-xl px-2.5 py-2.5 transition-colors hover:bg-base-content/[0.045] focus-within:bg-base-content/[0.045]"
                  onClick={() => submit(entry.query)}
                >
                  <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-base-content/[0.045] text-base-content/45">
                    <i className="bi-clock-history text-sm" aria-hidden="true" />
                  </span>
                  <span className="min-w-0 flex-1 truncate text-sm font-medium text-base-content/75">
                    {entry.query}
                  </span>
                  <span className="hidden text-[11px] text-base-content/30 sm:inline">
                    <Kbd>Enter</Kbd>
                  </span>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-9 w-9 shrink-0 rounded-lg text-base-content/35 opacity-80 hover:bg-base-content/[0.06] hover:text-base-content sm:h-7 sm:w-7 sm:opacity-70 group-hover/recent:opacity-100"
                    aria-label={t("delete")}
                    onClick={(e) => {
                      e.stopPropagation();
                      removeRecent(entry.query);
                    }}
                  >
                    <i className="bi-x-lg text-[11px]" aria-hidden="true" />
                  </Button>
                </div>
              ))}
            </div>
          </section>
        )}

        <section className={recents.length > 0 ? "mt-5" : ""}>
          <div className="px-1">
            <p className="text-[11px] font-semibold uppercase tracking-[0.08em] text-base-content/40">
              {t("suggested_search_operators")}
            </p>
          </div>

          <div className="mt-2 grid gap-1 sm:grid-cols-2">
            {ADVANCED_SEARCH_OPERATORS.map((entry) => (
              <button
                key={entry.operator}
                type="button"
                className="flex min-h-11 items-center justify-between gap-3 rounded-xl px-2.5 py-2.5 text-left transition-colors hover:bg-base-content/[0.045] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
                onClick={() => appendOperator(entry.operator)}
              >
                <div className="flex min-w-0 items-center gap-2.5">
                  <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/[0.07] text-primary/80">
                    <i className={`${entry.icon} text-sm`} aria-hidden="true" />
                  </span>
                  <span className="truncate text-sm font-medium text-base-content/70">
                    {t(entry.labelKey)}
                  </span>
                </div>
                <span className="shrink-0 rounded-md border border-base-content/10 bg-base-content/[0.035] px-1.5 py-0.5 font-mono text-[11px] text-base-content/50">
                  {entry.operator}
                </span>
              </button>
            ))}
          </div>
        </section>
      </div>

      <div className="mt-4 flex items-center justify-between border-t border-base-content/[0.07] pt-3 text-[11px] text-base-content/40">
        <div className="flex items-center gap-1.5">
          <Kbd>Enter</Kbd>
          <span>{t("search_for_links")}</span>
        </div>
        <div className="flex items-center gap-1.5">
          <Kbd>Esc</Kbd>
          <span>{t("cancel")}</span>
        </div>
      </div>
    </Modal>
  );
}
