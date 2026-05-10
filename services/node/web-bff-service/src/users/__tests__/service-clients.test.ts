import { describe, it, expect, vi, afterEach } from 'vitest';
import { ConfigService } from '@nestjs/config';
import { IdentityServiceClient } from '../clients/identity-service-client';
import { UserRoleServiceClient } from '../clients/user-role-service-client';

const originalFetch = globalThis.fetch;

function mockConfig(overrides: Record<string, string> = {}): ConfigService {
  return { get: (key: string, def: string) => overrides[key] ?? def } as unknown as ConfigService;
}

afterEach(() => {
  globalThis.fetch = originalFetch;
  vi.restoreAllMocks();
});

// ---------- IdentityServiceClient ----------

describe('IdentityServiceClient', () => {
  describe('listUsers', () => {
    it('returns users array from data envelope', async () => {
      const users = [{ id: 'u1', email: 'a@b.com', first_name: 'A', last_name: 'B', status: 'active' }];
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: users }),
      });
      const client = new IdentityServiceClient(mockConfig());
      const result = await client.listUsers('tenant-1');

      expect(result).toEqual(users);
      const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
      expect(url).toContain('/v1/users');
    });

    it('sends X-Tenant-Id header', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: [] }) });
      const client = new IdentityServiceClient(mockConfig());
      await client.listUsers('tid-abc');

      const opts = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![1] as RequestInit;
      expect((opts.headers as Record<string, string>)['X-Tenant-Id']).toBe('tid-abc');
    });

    it('returns empty array on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.listUsers('t')).toEqual([]);
    });

    it('returns empty array on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.listUsers('t')).toEqual([]);
    });

    it('returns empty array when body has no data', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({}) });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.listUsers('t')).toEqual([]);
    });

    it('uses base URL from config', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: [] }) });
      const client = new IdentityServiceClient(mockConfig({ IDENTITY_SERVICE_BASE_URL: 'http://id:9000' }));
      await client.listUsers('t');
      const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
      expect(url).toMatch(/^http:\/\/id:9000/);
    });
  });

  describe('getUser', () => {
    it('returns user from data envelope', async () => {
      const user = { id: 'u1', email: 'a@b.com' };
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: user }) });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.getUser('t', 'u1')).toEqual(user);
    });

    it('returns null on 404', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.getUser('t', 'u1')).toBeNull();
    });

    it('returns null on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('timeout'));
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.getUser('t', 'u1')).toBeNull();
    });
  });

  describe('createUser', () => {
    it('returns created user from data envelope', async () => {
      const user = { id: 'u2', email: 'c@d.com' };
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: user }) });
      const client = new IdentityServiceClient(mockConfig());
      const result = await client.createUser('t', { email: 'c@d.com', first_name: 'C', last_name: 'D' });
      expect(result).toEqual(user);

      const opts = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![1] as RequestInit;
      expect(opts.method).toBe('POST');
    });

    it('returns null on failure', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.createUser('t', { email: 'x', first_name: 'X', last_name: 'Y' })).toBeNull();
    });

    it('returns null on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.createUser('t', { email: 'x', first_name: 'X', last_name: 'Y' })).toBeNull();
    });
  });

  describe('deactivateUser', () => {
    it('returns true on success', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.deactivateUser('t', 'u1')).toBe(true);

      const opts = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![1] as RequestInit;
      expect(opts.method).toBe('PATCH');
    });

    it('returns false on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.deactivateUser('t', 'u1')).toBe(false);
    });

    it('returns false on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('timeout'));
      const client = new IdentityServiceClient(mockConfig());
      expect(await client.deactivateUser('t', 'u1')).toBe(false);
    });
  });
});

// ---------- UserRoleServiceClient ----------

describe('UserRoleServiceClient', () => {
  describe('listRoles', () => {
    it('returns roles from data envelope', async () => {
      const roles = [{ id: 'r1', name: 'admin', display_name: 'Admin' }];
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: roles }) });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.listRoles('t')).toEqual(roles);
    });

    it('returns empty array on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.listRoles('t')).toEqual([]);
    });

    it('returns empty array on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.listRoles('t')).toEqual([]);
    });

    it('uses base URL from config', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: [] }) });
      const client = new UserRoleServiceClient(mockConfig({ USER_ROLE_SERVICE_BASE_URL: 'http://ur:7000' }));
      await client.listRoles('t');
      const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
      expect(url).toMatch(/^http:\/\/ur:7000/);
    });
  });

  describe('getRoleDetail', () => {
    it('returns role with permissions from data envelope', async () => {
      const detail = { role: { id: 'r1', name: 'admin' }, permissions: [{ resource: 'gst', action: 'view' }] };
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: detail }) });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getRoleDetail('t', 'r1')).toEqual(detail);
    });

    it('returns null on 404', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getRoleDetail('t', 'r1')).toBeNull();
    });

    it('returns null on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('timeout'));
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getRoleDetail('t', 'r1')).toBeNull();
    });
  });

  describe('getUserRoles', () => {
    it('returns roles from data.roles envelope', async () => {
      const roles = [{ id: 'r1', name: 'admin' }];
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: { roles } }),
      });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getUserRoles('t', 'u1')).toEqual(roles);
    });

    it('returns empty array on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getUserRoles('t', 'u1')).toEqual([]);
    });

    it('returns empty array on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.getUserRoles('t', 'u1')).toEqual([]);
    });
  });

  describe('assignRole', () => {
    it('returns true on success', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.assignRole('t', 'u1', 'r1')).toBe(true);

      const opts = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![1] as RequestInit;
      expect(opts.method).toBe('POST');
      expect(JSON.parse(opts.body as string)).toEqual({ role_id: 'r1' });
    });

    it('returns false on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.assignRole('t', 'u1', 'r1')).toBe(false);
    });

    it('returns false on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('timeout'));
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.assignRole('t', 'u1', 'r1')).toBe(false);
    });
  });

  describe('createRole', () => {
    it('returns created role from data envelope', async () => {
      const role = { id: 'r2', name: 'custom' };
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: role }) });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.createRole('t', { name: 'custom', display_name: 'Custom' })).toEqual(role);
    });

    it('returns null on failure', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.createRole('t', { name: 'x', display_name: 'X' })).toBeNull();
    });
  });

  describe('updateRolePermissions', () => {
    it('returns true on success', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.updateRolePermissions('t', 'r1', ['p1', 'p2'])).toBe(true);

      const opts = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![1] as RequestInit;
      expect(opts.method).toBe('PUT');
      expect(JSON.parse(opts.body as string)).toEqual({ permission_ids: ['p1', 'p2'] });
    });

    it('returns false on non-2xx', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 403 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.updateRolePermissions('t', 'r1', ['p1'])).toBe(false);
    });
  });

  describe('seedRoles', () => {
    it('returns seeded:true on 200', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, status: 200 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.seedRoles('t1')).toEqual({ seeded: true, conflict: false });
    });

    it('returns conflict:true on 409', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 409 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.seedRoles('t1')).toEqual({ seeded: false, conflict: true });
    });

    it('returns seeded:false on other errors', async () => {
      globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.seedRoles('t1')).toEqual({ seeded: false, conflict: false });
    });

    it('returns seeded:false on network error', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
      const client = new UserRoleServiceClient(mockConfig());
      expect(await client.seedRoles('t1')).toEqual({ seeded: false, conflict: false });
    });
  });
});
