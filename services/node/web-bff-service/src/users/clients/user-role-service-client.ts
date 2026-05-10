import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { Role, RoleWithPermissions, CreateRolePayload } from '../users.types';

@Injectable()
export class UserRoleServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(UserRoleServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('USER_ROLE_SERVICE_BASE_URL', 'http://localhost:8083');
  }

  async listRoles(tenantId: string): Promise<Role[]> {
    const url = `${this.baseUrl}/v1/roles`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) return [];
      const body = await res.json().catch(() => null);
      return body?.data ?? [];
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`listRoles failed: ${msg}`);
      return [];
    }
  }

  async getRoleDetail(tenantId: string, roleId: string): Promise<RoleWithPermissions | null> {
    const url = `${this.baseUrl}/v1/roles/${roleId}`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) return null;
      const body = await res.json().catch(() => null);
      return body?.data ?? null;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`getRoleDetail failed: ${msg}`);
      return null;
    }
  }

  async getUserRoles(tenantId: string, userId: string): Promise<Role[]> {
    const url = `${this.baseUrl}/v1/users/${userId}/roles`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) return [];
      const body = await res.json().catch(() => null);
      return body?.data?.roles ?? [];
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`getUserRoles failed: ${msg}`);
      return [];
    }
  }

  async assignRole(tenantId: string, userId: string, roleId: string): Promise<boolean> {
    const url = `${this.baseUrl}/v1/users/${userId}/roles`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        method: 'POST',
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        body: JSON.stringify({ role_id: roleId }),
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) {
        this.logger.warn(`assignRole returned ${res.status}`);
        return false;
      }
      return true;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`assignRole failed: ${msg}`);
      return false;
    }
  }

  async createRole(tenantId: string, payload: CreateRolePayload): Promise<Role | null> {
    const url = `${this.baseUrl}/v1/roles`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        method: 'POST',
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) {
        this.logger.warn(`createRole returned ${res.status}`);
        return null;
      }
      const body = await res.json().catch(() => null);
      return body?.data ?? null;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`createRole failed: ${msg}`);
      return null;
    }
  }

  async updateRolePermissions(tenantId: string, roleId: string, permissionIds: string[]): Promise<boolean> {
    const url = `${this.baseUrl}/v1/roles/${roleId}/permissions`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        method: 'PUT',
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        body: JSON.stringify({ permission_ids: permissionIds }),
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) {
        this.logger.warn(`updateRolePermissions returned ${res.status}`);
        return false;
      }
      return true;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`updateRolePermissions failed: ${msg}`);
      return false;
    }
  }

  async seedRoles(tenantId: string): Promise<{ seeded: boolean; conflict: boolean }> {
    const url = `${this.baseUrl}/v1/tenants/${tenantId}/seed-roles`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (res.status === 409) return { seeded: false, conflict: true };
      if (!res.ok) {
        this.logger.warn(`seedRoles returned ${res.status}`);
        return { seeded: false, conflict: false };
      }
      return { seeded: true, conflict: false };
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`seedRoles failed: ${msg}`);
      return { seeded: false, conflict: false };
    }
  }
}
