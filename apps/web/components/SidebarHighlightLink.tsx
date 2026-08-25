import Link from "next/link";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

export default function SidebarHighlightLink({
  title,
  href,
  icon,
  active,
  external,
  sidebarIsCollapsed,
  expanded,
  onToggleExpand,
}: {
  title: string;
  href: string;
  icon: string;
  active?: boolean;
  external?: boolean;
  sidebarIsCollapsed?: boolean;
  expanded?: boolean;
  onToggleExpand?: () => void;
}) {
  const showToggle = Boolean(onToggleExpand) && !sidebarIsCollapsed;

  return (
    <TooltipProvider>
      <Tooltip>
        <div
          className={cn(
            "group flex items-center transition-colors",
            active ? "bg-primary/20" : "hover:bg-neutral/20",
            sidebarIsCollapsed
              ? "h-11 w-11 justify-center rounded-xl"
              : "min-h-10 rounded-lg"
          )}
        >
          <TooltipTrigger asChild>
            <Link
              href={href}
              title={title}
              target={external ? "_blank" : undefined}
              rel={external ? "noreferrer" : undefined}
              className={cn(
                "flex min-w-0 items-center gap-2 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50",
                sidebarIsCollapsed
                  ? "h-11 w-11 justify-center"
                  : "min-h-10 flex-1 px-3"
              )}
            >
              <i
                className={cn(
                  icon,
                  "shrink-0 text-lg text-primary drop-shadow"
                )}
                aria-hidden="true"
              />
              {!sidebarIsCollapsed && (
                <p className="min-w-0 flex-1 truncate text-xs font-bold capitalize">
                  {title}
                </p>
              )}
            </Link>
          </TooltipTrigger>

          {showToggle && (
            <button
              type="button"
              onClick={onToggleExpand}
              aria-expanded={expanded}
              aria-label={`${expanded ? "Collapse" : "Expand"} ${title}`}
              className="mr-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-base-content/50 transition-colors hover:bg-base-content/10 hover:text-base-content focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50"
            >
              <i
                className={cn(
                  expanded ? "bi-chevron-down" : "bi-chevron-right",
                  "text-[10px]"
                )}
                aria-hidden="true"
              />
            </button>
          )}
        </div>

        {sidebarIsCollapsed && (
          <TooltipContent side="right">{title}</TooltipContent>
        )}
      </Tooltip>
    </TooltipProvider>
  );
}
