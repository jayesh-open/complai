import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { BadRequestException, ForbiddenException, NotFoundException, InternalServerErrorException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { UsersController } from '../users.controller';
import { IdentityServiceClient } from '../clients/identity-service-client';
import { UserRoleServiceClient } from '../clients/user-role-service-client';

const originalFetch = globalThis.fetch;

function mockConfig(): ConfigService {
  return { get: (_k: string, d: string) => d } as unknown as ConfigService;
}

function mockCache() {
  const store = new Map<string, unknown>();
  return {
    get: vi.fn(async (key: string) => store.get(key) ?? undefined),
    set: vi.fn(async (key: string, value: unknown) => { store.set(key, value); }),
    _store: store,
  };
}

const TENANT = '00000000-0000-0000-0000-000000000001';
const USER_ID = '10000000-0000-0000-0000-000000000001';
const ROLE_ID = '20000000-0000-0000-0000-000000000001';

function createController(cacheOverride?: ReturnType<typeof mockCache>) {
  const config = mockConfig();
  const cache = cacheOverride ?? mockCache();
  const ctrl = new UsersController(
    new IdentityServiceClient(config),
    new UserRoleServiceClient(config),
    cache as any,
  );
  return { ctrl, cache };
}

afterEach(() => {
  globalThis.fetch = originalFetch;
  vi.restoreAllMocks();
});

beforeEach(() => {
  globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 503 });
});

// ---------- listUsers ----------

describe('UsersController — listUsers', () => {
  it('returns users array with enriched roles', async () => {
    const users = [{ id: USER_ID, email: 'a@b.com', first_name: 'A', last_name: 'B', status: 'active', created_at: '2026-01-01', updated_at: '2026-01-01' }];
    const roles = [{ id: ROLE_ID, name: 'admin', display_name: 'Admin' }];

    globalThis.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: users }) })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: { roles } }) });

    const { ctrl } = createController();
    const result = await ctrl.listUsers(TENANT);

    expect(result.users).toHaveLength(1);
    expect(result.users[0]!.role!.name).toBe('admin');
    expect(result.total).toBe(1);
  });

  it('throws 400 when tenant_id missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.listUsers()).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when tenant_id is invalid UUID', async () => {
    const { ctrl } = createController();
    await expect(ctrl.listUsers('not-a-uuid')).rejects.toThrow(BadRequestException);
  });

  it('returns empty users when identity service returns empty', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: [] }) });
    const { ctrl } = createController();
    const result = await ctrl.listUsers(TENANT);
    expect(result.users).toEqual([]);
    expect(result.total).toBe(0);
  });
});

// ---------- getUser ----------

describe('UsersController — getUser', () => {
  it('returns single user with role', async () => {
    const user = { id: USER_ID, email: 'a@b.com', first_name: 'A', last_name: 'B', status: 'active', created_at: '2026-01-01', updated_at: '2026-01-01' };
    const roles = [{ id: ROLE_ID, name: 'admin', display_name: 'Admin' }];

    globalThis.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: user }) })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: { roles } }) });

    const { ctrl } = createController();
    const result = await ctrl.getUser(USER_ID, TENANT);

    expect(result.id).toBe(USER_ID);
    expect(result.role!.name).toBe('admin');
  });

  it('throws 404 when user not found', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
    const { ctrl } = createController();
    await expect(ctrl.getUser(USER_ID, TENANT)).rejects.toThrow(NotFoundException);
  });

  it('throws 400 for invalid user_id', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getUser('not-uuid', TENANT)).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when tenant_id missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getUser(USER_ID)).rejects.toThrow(BadRequestException);
  });

  it('returns null role when user has no roles', async () => {
    const user = { id: USER_ID, email: 'a@b.com', first_name: 'A', last_name: 'B', status: 'active', created_at: '2026-01-01', updated_at: '2026-01-01' };

    globalThis.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: user }) })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: { roles: [] } }) });

    const { ctrl } = createController();
    const result = await ctrl.getUser(USER_ID, TENANT);
    expect(result.role).toBeNull();
  });
});

// ---------- createUser ----------

describe('UsersController — createUser', () => {
  it('creates user and assigns role', async () => {
    const user = { id: 'u-new', email: 'new@b.com' };
    globalThis.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: user }) })
      .mockResolvedValueOnce({ ok: true });

    const { ctrl } = createController();
    const result = await ctrl.createUser(TENANT, { email: 'new@b.com', first_name: 'N', last_name: 'U', role_id: ROLE_ID });

    expect(result.user).toEqual(user);
    expect(result.role_assignment_failed).toBe(false);
  });

  it('returns role_assignment_failed when role assign fails', async () => {
    const user = { id: 'u-new', email: 'new@b.com' };
    globalThis.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ data: user }) })
      .mockResolvedValueOnce({ ok: false, status: 500 });

    const { ctrl } = createController();
    const result = await ctrl.createUser(TENANT, { email: 'new@b.com', first_name: 'N', last_name: 'U', role_id: ROLE_ID });

    expect(result.user).toEqual(user);
    expect(result.role_assignment_failed).toBe(true);
  });

  it('throws 500 when identity-service create fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const { ctrl } = createController();
    await expect(
      ctrl.createUser(TENANT, { email: 'new@b.com', first_name: 'N', last_name: 'U' }),
    ).rejects.toThrow(InternalServerErrorException);
  });

  it('throws 400 when email missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.createUser(TENANT, { first_name: 'N', last_name: 'U' } as any)).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when tenant_id missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.createUser(undefined, { email: 'a', first_name: 'N', last_name: 'U' })).rejects.toThrow(BadRequestException);
  });
});

// ---------- changeUserRole ----------

describe('UsersController — changeUserRole', () => {
  it('assigns new role successfully', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true });
    const { ctrl } = createController();
    const result = await ctrl.changeUserRole(USER_ID, TENANT, { role_id: ROLE_ID });
    expect(result.status).toBe('role_updated');
  });

  it('throws 400 when role_id missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.changeUserRole(USER_ID, TENANT, {} as any)).rejects.toThrow(BadRequestException);
  });

  it('throws 400 for invalid user_id', async () => {
    const { ctrl } = createController();
    await expect(ctrl.changeUserRole('bad', TENANT, { role_id: ROLE_ID })).rejects.toThrow(BadRequestException);
  });

  it('throws 500 when assign fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const { ctrl } = createController();
    await expect(ctrl.changeUserRole(USER_ID, TENANT, { role_id: ROLE_ID })).rejects.toThrow(InternalServerErrorException);
  });
});

// ---------- deactivateUser ----------

describe('UsersController — deactivateUser', () => {
  it('deactivates user successfully', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true });
    const { ctrl } = createController();
    const result = await ctrl.deactivateUser(USER_ID, TENANT);
    expect(result.status).toBe('deactivated');
  });

  it('throws 400 for invalid user_id', async () => {
    const { ctrl } = createController();
    await expect(ctrl.deactivateUser('bad', TENANT)).rejects.toThrow(BadRequestException);
  });

  it('throws 500 when deactivation fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const { ctrl } = createController();
    await expect(ctrl.deactivateUser(USER_ID, TENANT)).rejects.toThrow(InternalServerErrorException);
  });
});

// ---------- listRoles ----------

describe('UsersController — listRoles', () => {
  it('returns roles list', async () => {
    const roles = [{ id: ROLE_ID, name: 'admin' }];
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: roles }) });
    const { ctrl } = createController();
    const result = await ctrl.listRoles(TENANT);
    expect(result).toEqual(roles);
  });

  it('caches roles list', async () => {
    const roles = [{ id: ROLE_ID, name: 'admin' }];
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: roles }) });
    const cache = mockCache();
    const { ctrl } = createController(cache);

    await ctrl.listRoles(TENANT);
    expect(cache.set).toHaveBeenCalledOnce();
    const [key] = cache.set.mock.calls[0]!;
    expect(key).toBe(`roles:list:${TENANT}`);
  });

  it('returns cached roles on hit', async () => {
    const cached = [{ id: ROLE_ID, name: 'cached_admin' }];
    const cache = mockCache();
    cache._store.set(`roles:list:${TENANT}`, cached);
    const { ctrl } = createController(cache);

    const result = await ctrl.listRoles(TENANT);
    expect(result).toBe(cached);
  });

  it('throws 400 when tenant_id missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.listRoles()).rejects.toThrow(BadRequestException);
  });
});

// ---------- getRoleDetail ----------

describe('UsersController — getRoleDetail', () => {
  it('returns role with permissions', async () => {
    const detail = { role: { id: ROLE_ID, name: 'admin', is_system: true }, permissions: [] };
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: detail }) });
    const { ctrl } = createController();
    const result = await ctrl.getRoleDetail(ROLE_ID, TENANT);
    expect(result).toEqual(detail);
  });

  it('throws 404 when role not found', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
    const { ctrl } = createController();
    await expect(ctrl.getRoleDetail(ROLE_ID, TENANT)).rejects.toThrow(NotFoundException);
  });

  it('throws 400 for invalid role_id', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getRoleDetail('bad', TENANT)).rejects.toThrow(BadRequestException);
  });
});

// ---------- updateRolePermissions ----------

describe('UsersController — updateRolePermissions', () => {
  it('throws 403 for system role', async () => {
    const detail = { role: { id: ROLE_ID, name: 'admin', is_system: true }, permissions: [] };
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: detail }) });
    const { ctrl } = createController();

    await expect(
      ctrl.updateRolePermissions(ROLE_ID, TENANT, { permission_pairs: [{ resource: 'gst', action: 'view' }] }),
    ).rejects.toThrow(ForbiddenException);
  });

  it('throws 404 when role not found', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
    const { ctrl } = createController();

    await expect(
      ctrl.updateRolePermissions(ROLE_ID, TENANT, { permission_pairs: [] }),
    ).rejects.toThrow(NotFoundException);
  });

  it('throws 400 when permission_pairs missing', async () => {
    const { ctrl } = createController();
    await expect(
      ctrl.updateRolePermissions(ROLE_ID, TENANT, {} as any),
    ).rejects.toThrow(BadRequestException);
  });

  it('throws 400 for invalid role_id', async () => {
    const { ctrl } = createController();
    await expect(
      ctrl.updateRolePermissions('bad', TENANT, { permission_pairs: [] }),
    ).rejects.toThrow(BadRequestException);
  });

  it('throws 400 (pending resolution) for non-system role', async () => {
    const detail = { role: { id: ROLE_ID, name: 'custom', is_system: false }, permissions: [] };
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: detail }) });
    const { ctrl } = createController();

    await expect(
      ctrl.updateRolePermissions(ROLE_ID, TENANT, { permission_pairs: [{ resource: 'gst', action: 'view' }] }),
    ).rejects.toThrow(BadRequestException);
  });
});

// ---------- createRole ----------

describe('UsersController — createRole', () => {
  it('creates role successfully', async () => {
    const role = { id: 'r-new', name: 'custom' };
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: role }) });
    const { ctrl } = createController();
    const result = await ctrl.createRole(TENANT, { name: 'custom', display_name: 'Custom' });
    expect(result).toEqual(role);
  });

  it('throws 400 when name missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.createRole(TENANT, { display_name: 'X' } as any)).rejects.toThrow(BadRequestException);
  });

  it('throws 500 when create fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const { ctrl } = createController();
    await expect(
      ctrl.createRole(TENANT, { name: 'x', display_name: 'X' }),
    ).rejects.toThrow(InternalServerErrorException);
  });
});

// ---------- seedRoles ----------

describe('UsersController — seedRoles', () => {
  it('seeds successfully', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: true, status: 200 });
    const { ctrl } = createController();
    const result = await ctrl.seedRoles(TENANT);
    expect(result.status).toBe('seeded');
  });

  it('returns already_seeded on conflict', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 409 });
    const { ctrl } = createController();
    const result = await ctrl.seedRoles(TENANT);
    expect(result.status).toBe('already_seeded');
  });

  it('throws 400 for invalid tenant_id', async () => {
    const { ctrl } = createController();
    await expect(ctrl.seedRoles('bad')).rejects.toThrow(BadRequestException);
  });

  it('throws 500 when seed fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const { ctrl } = createController();
    await expect(ctrl.seedRoles(TENANT)).rejects.toThrow(InternalServerErrorException);
  });
});
