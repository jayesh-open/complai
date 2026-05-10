import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { User, CreateUserPayload } from '../users.types';

@Injectable()
export class IdentityServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(IdentityServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('IDENTITY_SERVICE_BASE_URL', 'http://localhost:8081');
  }

  async listUsers(tenantId: string): Promise<User[]> {
    const url = `${this.baseUrl}/v1/users`;
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
      this.logger.warn(`listUsers failed: ${msg}`);
      return [];
    }
  }

  async getUser(tenantId: string, userId: string): Promise<User | null> {
    const url = `${this.baseUrl}/v1/users/${userId}`;
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
      this.logger.warn(`getUser failed: ${msg}`);
      return null;
    }
  }

  async createUser(tenantId: string, payload: CreateUserPayload): Promise<User | null> {
    const url = `${this.baseUrl}/v1/users`;
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
        this.logger.warn(`createUser returned ${res.status}`);
        return null;
      }
      const body = await res.json().catch(() => null);
      return body?.data ?? null;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`createUser failed: ${msg}`);
      return null;
    }
  }

  async deactivateUser(tenantId: string, userId: string): Promise<boolean> {
    const url = `${this.baseUrl}/v1/users/${userId}`;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        method: 'PATCH',
        headers: { 'X-Tenant-Id': tenantId, 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'inactive' }),
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) {
        this.logger.warn(`deactivateUser returned ${res.status}`);
        return false;
      }
      return true;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      this.logger.warn(`deactivateUser failed: ${msg}`);
      return false;
    }
  }
}
