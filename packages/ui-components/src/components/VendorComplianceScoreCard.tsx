import { cn } from '../lib/utils';

interface ScoreDimension {
  label: string;
  score: number;
  maxScore: number;
  status: 'pass' | 'warn' | 'fail';
  note?: string;
}

interface VendorComplianceScoreCardProps {
  vendorName: string;
  gstin: string;
  state: string;
  totalScore: number;
  maxScore?: number;
  riskLevel: string;
  category: string;
  lastReviewed: string;
  dimensions: ScoreDimension[];
  className?: string;
}

export function VendorComplianceScoreCard({
  vendorName, gstin, state, totalScore, maxScore = 100,
  riskLevel, category, lastReviewed, dimensions, className,
}: VendorComplianceScoreCardProps) {
  const scoreColor = totalScore >= 80 ? 'var(--success)' : totalScore >= 60 ? 'var(--warning)' : 'var(--danger)';

  return (
    <div className={cn('bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl p-5', className)} data-testid="vendor-scorecard">
      <div className="text-sm font-bold text-[var(--text-primary)]">{vendorName}</div>
      <div className="text-[11px] text-[var(--text-muted)] mt-0.5 font-mono">{gstin} &middot; {state}</div>
      <div className="mt-4 flex items-center gap-3">
        <div className="flex gap-0.5">
          {Array.from({ length: 10 }, (_, i) => (
            <div
              key={i}
              className="w-2.5 h-2.5 rounded-full"
              style={{
                backgroundColor: i < Math.round(totalScore / (maxScore / 10)) ? scoreColor : 'var(--border-default)',
              }}
            />
          ))}
        </div>
        <span className="text-sm font-bold" style={{ color: scoreColor }}>{totalScore} / {maxScore}</span>
        <span className="text-[11px] text-[var(--text-muted)]">&middot; Risk: {riskLevel}</span>
      </div>
      <div className="mt-3 bg-[var(--bg-tertiary)] border border-[var(--border-default)] rounded-lg p-3 space-y-2">
        {dimensions.map((d) => (
          <div key={d.label} className="flex items-center gap-2 text-xs">
            <span className="text-[var(--text-muted)] min-w-[140px]">{d.label}:</span>
            <span className="font-semibold text-[var(--text-primary)]">{d.score}/{d.maxScore}</span>
            <span>{d.status === 'pass' ? '✓' : '⚠'}</span>
            {d.note && <span className="text-[var(--text-disabled)]">&mdash; {d.note}</span>}
          </div>
        ))}
      </div>
      <div className="mt-3 text-[11px] text-[var(--text-muted)]">
        Category: {category} &middot; Last reviewed: {lastReviewed}
      </div>
    </div>
  );
}
