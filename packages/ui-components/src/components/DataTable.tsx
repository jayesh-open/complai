"use client";
import { cn } from '../lib/utils';
import { useThemeStore, type DensityMode } from '../lib/theme-store';
import { DENSITY_ROW_HEIGHT } from '../lib/utils';

interface Column<T> {
  key: string;
  header: string;
  render?: (row: T) => React.ReactNode;
  className?: string;
  align?: 'left' | 'right' | 'center';
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  density?: DensityMode;
  onRowClick?: (row: T) => void;
  emptyMessage?: string;
  className?: string;
}

export function DataTable<T extends Record<string, unknown>>({
  columns,
  data,
  density: densityProp,
  onRowClick,
  emptyMessage = 'No data found',
  className,
}: DataTableProps<T>) {
  const storeDensity = useThemeStore((s) => s.density);
  const density = densityProp ?? storeDensity;
  const rowHeight = DENSITY_ROW_HEIGHT[density];

  return (
    <div className={cn('bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden', className)}>
      <table className="w-full" data-density={density} data-testid="data-table">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            {columns.map((col) => (
              <th
                key={col.key}
                className={cn(
                  'px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)]',
                  col.align === 'right' && 'text-right',
                  col.align === 'center' && 'text-center',
                  col.className,
                )}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} className="text-center py-12 text-[var(--text-muted)] text-sm">
                {emptyMessage}
              </td>
            </tr>
          ) : (
            data.map((row, i) => (
              <tr
                key={i}
                onClick={() => onRowClick?.(row)}
                className={cn(
                  'border-b border-[var(--border-default)] last:border-b-0',
                  'hover:bg-[var(--bg-tertiary)] transition-colors duration-150',
                  onRowClick && 'cursor-pointer',
                )}
                style={{ height: `${rowHeight}px` }}
              >
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className={cn(
                      'px-[18px] text-xs text-[var(--text-secondary)]',
                      col.align === 'right' && 'text-right tabular-nums font-semibold',
                      col.align === 'center' && 'text-center',
                      col.className,
                    )}
                  >
                    {col.render ? col.render(row) : String(row[col.key] ?? '')}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
