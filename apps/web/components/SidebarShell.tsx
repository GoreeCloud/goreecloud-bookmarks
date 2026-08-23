import { ReactNode, useEffect, useState } from "react";
import { useRouter } from "next/router";
import { useTranslation } from "next-i18next";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import ProfileDropdown from "@/components/ProfileDropdown";
import { Button } from "./ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import Image from "next/image";
import Link from "next/link";

export function ActionIcon({
  icon,
  label,
  onClick,
  variant = "metal",
  tooltipSide,
}: {
  icon: string;
  label: string;
  onClick: () => void;
  variant?:
    | "metal"
    | "link"
    | "default"
    | "primary"
    | "accent"
    | "destructive"
    | "outline"
    | "secondary"
    | "ghost"
    | "simple";
  tooltipSide?: "top" | "right" | "bottom" | "left";
}) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant={variant}
            size="icon"
            className="h-11 w-11 shrink-0 rounded-xl lg:h-9 lg:w-9"
            onClick={onClick}
            aria-label={label}
          >
            <i className={cn(icon, "text-lg leading-none")} aria-hidden="true" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side={tooltipSide}>{label}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}

type SidebarShellContext = {
  collapsed?: boolean;
  toggle?: () => void;
  closeMobileSidebar: () => void;
};

export default function SidebarShell({
  className,
  toggleSidebar,
  sidebarIsCollapsed,
  children,
}: {
  className?: string;
  toggleSidebar?: () => void;
  sidebarIsCollapsed?: boolean;
  children: (ctx: SidebarShellContext) => ReactNode;
}) {
  const { t } = useTranslation();
  const router = useRouter();

  const { width } = useWindowDimensions();
  const isMobile = width < 1024;
  const [mobileExpanded, setMobileExpanded] = useState(false);

  useEffect(() => {
    setMobileExpanded(false);
  }, [router.asPath]);

  useEffect(() => {
    if (!isMobile || !mobileExpanded) return;

    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") setMobileExpanded(false);
    };

    document.addEventListener("keydown", closeOnEscape);
    return () => document.removeEventListener("keydown", closeOnEscape);
  }, [isMobile, mobileExpanded]);

  const collapsed = isMobile ? !mobileExpanded : sidebarIsCollapsed;
  const toggle = isMobile
    ? () => setMobileExpanded(!mobileExpanded)
    : toggleSidebar;

  return (
    <>
      {isMobile && mobileExpanded && (
        <>
          <div className="w-16 shrink-0" />
          <div
            className="fixed inset-0 z-30 bg-black/20 backdrop-blur-sm fade-in"
            onClick={() => setMobileExpanded(false)}
            aria-hidden="true"
          />
        </>
      )}
      <aside
        id="sidebar"
        className={cn(
          "z-20 flex h-screen flex-col border-r border-base-content/10 bg-base-200/80 backdrop-blur-xl",
          className,
          collapsed ? "w-16 p-2" : "w-72 p-2",
          isMobile &&
            mobileExpanded &&
            "fixed inset-y-0 left-0 z-40 shadow-2xl"
        )}
      >
        <Link
          href="/dashboard"
          className={cn(
            "mb-2 flex shrink-0 items-center rounded-xl transition-colors hover:bg-base-content/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary",
            collapsed ? "h-11 justify-center" : "gap-3 px-2 py-2"
          )}
          aria-label="GoreeCloud Bookmarks dashboard"
        >
          <Image
            src="/goreecloud-bookmarks.svg"
            width={36}
            height={36}
            alt=""
            aria-hidden="true"
            className="h-9 w-9 shrink-0 rounded-xl"
            unoptimized
            priority
          />
          {!collapsed && (
            <div className="min-w-0 leading-tight">
              <p className="truncate text-sm font-semibold text-base-content">
                GoreeCloud Bookmarks
              </p>
              <p className="truncate text-[11px] text-base-content/55">
                Private library
              </p>
            </div>
          )}
        </Link>

        <div className="flex min-h-0 flex-1 flex-col">
          {children({
            collapsed,
            toggle,
            closeMobileSidebar: () => setMobileExpanded(false),
          })}
        </div>

        <div
          className={cn(
            "mt-2 flex shrink-0 items-center gap-2 border-t border-base-content/10 pt-2",
            collapsed ? "flex-col" : "relative justify-between"
          )}
        >
          {collapsed ? (
            <ProfileDropdown
              tooltipSide="right"
              dropdownSide="right"
              dropdownAlign="end"
            />
          ) : (
            <div className="min-w-0 flex-1 px-1">
              <ProfileDropdown showName />
            </div>
          )}

          {toggle && (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    onClick={toggle}
                    size="icon"
                    className="h-11 w-11 rounded-xl lg:h-9 lg:w-9"
                    aria-label={
                      collapsed ? t("expand_sidebar") : t("shrink_sidebar")
                    }
                  >
                    <i className="bi-layout-sidebar" aria-hidden="true" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent side={collapsed ? "right" : "top"}>
                  {collapsed ? t("expand_sidebar") : t("shrink_sidebar")}
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}
        </div>
      </aside>
    </>
  );
}
