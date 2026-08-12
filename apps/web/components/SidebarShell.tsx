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
            className="h-9 w-9 shrink-0 rounded-xl"
            onClick={onClick}
            aria-label={label}
          >
            <i className={cn(icon, "text-lg leading-none")} />
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

  const collapsed = isMobile ? !mobileExpanded : sidebarIsCollapsed;
  const toggle = isMobile
    ? () => setMobileExpanded(!mobileExpanded)
    : toggleSidebar;

  return (
    <>
      {isMobile && mobileExpanded && (
        <>
          <div className="w-14 shrink-0" />
          <div
            className="fixed inset-0 bg-black/20 backdrop-blur-sm z-30 fade-in"
            onClick={() => setMobileExpanded(false)}
          />
        </>
      )}
      <aside
        id="sidebar"
        className={cn(
          "h-screen flex flex-col z-20 bg-base-200/80 backdrop-blur-xl border-r border-base-content/10",
          className,
          collapsed ? "p-2 w-14" : "p-2 w-72",
          isMobile &&
            mobileExpanded &&
            "fixed inset-y-0 left-0 z-40 shadow-2xl"
        )}
      >
        <Link
          href="/dashboard"
          className={cn(
            "shrink-0 flex items-center rounded-xl mb-2 transition-colors hover:bg-base-content/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary",
            collapsed ? "justify-center h-10" : "gap-3 px-2 py-2"
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

        <div className="min-h-0 flex-1 flex flex-col">
          {children({
            collapsed,
            toggle,
            closeMobileSidebar: () => setMobileExpanded(false),
          })}
        </div>

        <div
          className={cn(
            "shrink-0 flex items-center gap-2 mt-2 pt-2 border-t border-base-content/10",
            collapsed ? "flex-col" : "justify-between relative"
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
                    className="rounded-xl"
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
