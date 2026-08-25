import { ReactNode } from "react";
import { Separator } from "../ui/separator";

type Tone = "default" | "destructive";

type Props = {
  title: ReactNode;
  description?: ReactNode;
  icon?: string;
  tone?: Tone;
  children: ReactNode;
  footer?: ReactNode;
};

export default function GlazeModalFrame({
  title,
  description,
  icon = "bi-window",
  tone = "default",
  children,
  footer,
}: Props) {
  const iconTone =
    tone === "destructive"
      ? "border-error/20 bg-error/10 text-error"
      : "border-base-300 bg-base-200 text-primary";

  return (
    <>
      <div className="flex items-start gap-3 pr-10">
        <div
          className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border ${iconTone}`}
          aria-hidden="true"
        >
          <i className={`${icon} text-base`} />
        </div>
        <div className="min-w-0 pt-0.5">
          <h2 className="text-lg font-semibold leading-6 text-base-content">
            {title}
          </h2>
          {description && (
            <p className="mt-1 text-sm leading-5 text-base-content/60">
              {description}
            </p>
          )}
        </div>
      </div>

      <Separator className="my-4" />

      <div className="space-y-4 text-sm leading-6 text-base-content/80">
        {children}
      </div>

      {footer && (
        <div className="mt-5 flex items-center justify-end gap-2 border-t border-base-300 pt-4">
          {footer}
        </div>
      )}
    </>
  );
}
