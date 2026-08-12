import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import { toast } from "react-hot-toast";
import { useTranslation } from "next-i18next";
import { Button } from "@/components/ui/button";

type Props = {
  placeholder?: string;
  fullWidth?: boolean;
};

export const ADVANCED_SEARCH_OPERATORS = [
  {
    operator: "url:",
    labelKey: "url",
    icon: "bi-link-45deg",
  },
  {
    operator: "tag:",
    labelKey: "tag",
    icon: "bi-tag",
  },
  {
    operator: "pinned:true",
    labelKey: "pinned",
    icon: "bi-pin-angle",
  },
  {
    operator: "before:",
    labelKey: "before",
    icon: "bi-calendar-minus",
  },
  {
    operator: "after:",
    labelKey: "after",
    icon: "bi-calendar-plus",
  },
] as const;

export default function SearchBar({ placeholder, fullWidth }: Props) {
  const router = useRouter();
  const { t } = useTranslation();
  const [searchQuery, setSearchQuery] = useState("");
  const [showSuggestions, setShowSuggestions] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    router.query.q
      ? setSearchQuery(decodeURIComponent(router.query.q as string))
      : setSearchQuery("");
  }, [router.query.q]);

  const handleSuggestionClick = (operator: string) => {
    setSearchQuery((prev) => {
      const needsSpace = prev.length > 0 && !prev.endsWith(" ");
      return `${prev}${needsSpace ? " " : ""}${operator}`;
    });
    requestAnimationFrame(() => inputRef.current?.focus());
  };

  const submitSearch = () => {
    if (router.pathname.startsWith("/public")) {
      if (!searchQuery) {
        return router.push("/public/collections/" + router.query.id);
      }

      return router.push(
        "/public/collections/" +
          router.query.id +
          "?q=" +
          encodeURIComponent(searchQuery || "")
      );
    }

    return router.push("/search?q=" + encodeURIComponent(searchQuery));
  };

  return (
    <div
      className={`relative flex items-center ${fullWidth ? "w-full" : ""}`}
    >
      <div className="pointer-events-none absolute left-3 z-10 flex h-8 w-8 items-center justify-center text-base-content/45">
        <i className="bi-search text-sm leading-none" aria-hidden="true" />
      </div>

      <input
        id="search-box"
        type="search"
        ref={inputRef}
        aria-label={placeholder || t("search_for_links")}
        placeholder={placeholder || t("search_for_links")}
        value={searchQuery}
        onFocus={() => setShowSuggestions(true)}
        onBlur={() => setShowSuggestions(false)}
        onChange={(e) => {
          e.target.value.includes("%") &&
            toast.error(t("search_query_invalid_symbol"));
          setSearchQuery(e.target.value.replace("%", ""));
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            e.preventDefault();
            submitSearch();
          }
          if (e.key === "Escape") {
            setShowSuggestions(false);
            inputRef.current?.blur();
          }
        }}
        className={`h-11 rounded-xl border border-base-content/10 bg-base-100 pl-11 pr-12 text-sm text-base-content shadow-sm outline-none transition-all duration-150 placeholder:text-base-content/35 hover:border-base-content/20 focus:border-primary/40 focus:ring-2 focus:ring-primary/10 sm:h-10 sm:pr-11 ${
          fullWidth ? "w-full" : "w-full max-w-[15rem] md:w-80 md:max-w-full"
        }`}
      />

      {searchQuery && (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="absolute right-0.5 h-10 w-10 rounded-lg text-base-content/40 hover:bg-base-content/[0.06] hover:text-base-content sm:right-1.5 sm:h-7 sm:w-7"
          aria-label="Clear search"
          onMouseDown={(e) => e.preventDefault()}
          onClick={() => {
            setSearchQuery("");
            requestAnimationFrame(() => inputRef.current?.focus());
          }}
        >
          <i className="bi-x-lg text-xs" aria-hidden="true" />
        </Button>
      )}

      {showSuggestions && (
        <div className="absolute left-0 top-full z-50 mt-2 w-full min-w-[18rem]">
          <div
            className="overflow-hidden rounded-xl border border-base-content/10 bg-base-100 p-2 shadow-xl ring-1 ring-black/[0.02]"
            onMouseDown={(e) => e.preventDefault()}
          >
            <div className="flex items-center justify-between px-2 pb-1.5 pt-1">
              <p className="text-[11px] font-semibold uppercase tracking-[0.08em] text-base-content/40">
                {t("suggested_search_operators")}
              </p>
            </div>

            <div className="flex flex-col gap-0.5">
              {ADVANCED_SEARCH_OPERATORS.map((entry) => (
                <button
                  key={entry.operator}
                  type="button"
                  className="flex min-h-11 items-center justify-between gap-3 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-base-content/[0.045] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
                  onClick={() => handleSuggestionClick(entry.operator)}
                >
                  <div className="flex min-w-0 items-center gap-2.5">
                    <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/[0.07] text-primary/80">
                      <i className={`${entry.icon} text-sm`} aria-hidden="true" />
                    </span>
                    <span className="truncate text-xs font-medium text-base-content/70">
                      {t(entry.labelKey)}
                    </span>
                  </div>
                  <span className="shrink-0 rounded-md border border-base-content/10 bg-base-content/[0.035] px-1.5 py-0.5 font-mono text-[11px] text-base-content/50">
                    {entry.operator}
                  </span>
                </button>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
