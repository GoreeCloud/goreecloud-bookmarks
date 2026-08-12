export default function DashboardItem({
  name,
  value,
  icon,
}: {
  name: string;
  value: number;
  icon: string;
}) {
  return (
    <div className="group flex min-h-24 items-center justify-between gap-4 w-full rounded-2xl border border-base-content/10 bg-base-100 p-4 shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-primary/30 hover:shadow-md">
      <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/10 transition-transform duration-200 group-hover:scale-105">
        <i className={`${icon} text-xl`} aria-hidden="true"></i>
      </div>
      <div className="min-w-0 text-right">
        <p className="truncate text-xs font-medium uppercase tracking-[0.12em] text-base-content/55">
          {name}
        </p>
        <p className="mt-1 text-3xl font-semibold tracking-tight text-base-content tabular-nums">
          {value || 0}
        </p>
      </div>
    </div>
  );
}
