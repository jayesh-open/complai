import type { EmployeeITRDetail } from "../../compliance/itr/types";
import { getEmployeeDetail } from "../../compliance/itr/mock-detail-data";

export type TokenStatus = "valid" | "expired" | "used";
export type MagicLinkStatus = "SENT" | "VIEWED" | "APPROVED" | "EXPIRED" | "USED";

export interface TokenVerificationResult {
  status: TokenStatus;
  employeeDetail?: EmployeeITRDetail;
  employerName?: string;
  expiresAt?: string;
  expiredAt?: string;
  usedAt?: string;
  outcome?: "APPROVED" | "FILED" | "ACKNOWLEDGED" | "REJECTED";
  generatedAt?: string;
}

const VALID_TOKENS = ["mlk-valid-001", "mlk-valid-002", "mlk-valid-003"];
const EXPIRED_TOKENS = ["mlk-expired-001", "mlk-expired-002"];
const USED_TOKENS = ["mlk-used-001", "mlk-used-002"];

export function verifyMagicLinkToken(token: string): Promise<TokenVerificationResult> {
  return new Promise((resolve) => {
    setTimeout(() => {
      if (VALID_TOKENS.includes(token)) {
        const detail = getEmployeeDetail("emp-itr-0001");
        resolve({
          status: "valid",
          employeeDetail: detail ?? undefined,
          employerName: "Acme Technologies Pvt. Ltd.",
          expiresAt: "2026-08-01T23:59:59Z",
          generatedAt: "2026-07-25T10:00:00Z",
        });
      } else if (EXPIRED_TOKENS.includes(token)) {
        resolve({
          status: "expired",
          expiredAt: "2026-07-20T23:59:59Z",
          generatedAt: "2026-07-13T10:00:00Z",
        });
      } else if (USED_TOKENS.includes(token)) {
        resolve({
          status: "used",
          usedAt: "2026-07-22T14:30:00Z",
          outcome: "APPROVED",
          generatedAt: "2026-07-18T10:00:00Z",
        });
      } else {
        resolve({
          status: "expired",
          expiredAt: "2026-07-01T00:00:00Z",
          generatedAt: "2026-06-24T10:00:00Z",
        });
      }
    }, 800);
  });
}
