import { describe, it, expect, vi } from 'vitest';
import { enrichUsersWithRoles } from '../users.aggregator';
import type { UserRoleServiceClient } from '../clients/user-role-service-client';
import type { User, Role } from '../users.types';

function makeUser(overrides: Partial<User> = {}): User {
  return {
    id: 'u1',
    tenant_id: 't1',
    email: 'test@example.com',
    email_verified: true,
    first_name: 'Test',
    last_name: 'User',
    status: 'active',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

function makeRole(overrides: Partial<Role> = {}): Role {
  return {
    id: 'r1',
    tenant_id: 't1',
    name: 'admin',
    display_name: 'Admin',
    is_system: true,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

function mockUserRoleClient(getUserRolesFn?: (tid: string, uid: string) => Promise<Role[]>): UserRoleServiceClient {
  return {
    getUserRoles: getUserRolesFn ?? vi.fn().mockResolvedValue([]),
  } as unknown as UserRoleServiceClient;
}

describe('enrichUsersWithRoles', () => {
  it('enriches single user with their primary role', async () => {
    const user = makeUser({ id: 'u1' });
    const role = makeRole({ id: 'r1', name: 'admin', display_name: 'Admin' });
    const client = mockUserRoleClient(async () => [role]);

    const result = await enrichUsersWithRoles([user], 't1', client);

    expect(result).toHaveLength(1);
    expect(result[0]!.id).toBe('u1');
    expect(result[0]!.role).toEqual({ id: 'r1', name: 'admin', display_name: 'Admin' });
  });

  it('enriches multiple users in parallel', async () => {
    const users = [
      makeUser({ id: 'u1', email: 'a@b.com' }),
      makeUser({ id: 'u2', email: 'c@d.com' }),
    ];
    const fn = vi.fn()
      .mockResolvedValueOnce([makeRole({ id: 'r1', name: 'admin', display_name: 'Admin' })])
      .mockResolvedValueOnce([makeRole({ id: 'r2', name: 'auditor', display_name: 'Auditor' })]);
    const client = mockUserRoleClient(fn);

    const result = await enrichUsersWithRoles(users, 't1', client);

    expect(result).toHaveLength(2);
    expect(result[0]!.role!.name).toBe('admin');
    expect(result[1]!.role!.name).toBe('auditor');
    expect(fn).toHaveBeenCalledTimes(2);
  });

  it('sets role to null when user has no roles', async () => {
    const user = makeUser();
    const client = mockUserRoleClient(async () => []);

    const result = await enrichUsersWithRoles([user], 't1', client);

    expect(result[0]!.role).toBeNull();
  });

  it('uses first role when user has multiple roles', async () => {
    const user = makeUser();
    const roles = [
      makeRole({ id: 'r1', name: 'admin', display_name: 'Admin' }),
      makeRole({ id: 'r2', name: 'auditor', display_name: 'Auditor' }),
    ];
    const client = mockUserRoleClient(async () => roles);

    const result = await enrichUsersWithRoles([user], 't1', client);

    expect(result[0]!.role!.id).toBe('r1');
  });

  it('handles role lookup error gracefully — sets role to null', async () => {
    const users = [
      makeUser({ id: 'u1' }),
      makeUser({ id: 'u2' }),
    ];
    const fn = vi.fn()
      .mockResolvedValueOnce([makeRole()])
      .mockRejectedValueOnce(new Error('network error'));
    const client = mockUserRoleClient(fn);

    const result = await enrichUsersWithRoles(users, 't1', client);

    expect(result[0]!.role).not.toBeNull();
    expect(result[1]!.role).toBeNull();
  });

  it('returns empty array for empty input', async () => {
    const client = mockUserRoleClient();
    const result = await enrichUsersWithRoles([], 't1', client);
    expect(result).toEqual([]);
  });

  it('preserves all user fields in output', async () => {
    const user = makeUser({
      id: 'u1',
      email: 'test@x.com',
      first_name: 'Jane',
      last_name: 'Doe',
      status: 'active',
      created_at: '2026-03-01T00:00:00Z',
      updated_at: '2026-04-01T00:00:00Z',
    });
    const client = mockUserRoleClient(async () => []);

    const result = await enrichUsersWithRoles([user], 't1', client);

    expect(result[0]).toEqual({
      id: 'u1',
      email: 'test@x.com',
      first_name: 'Jane',
      last_name: 'Doe',
      status: 'active',
      created_at: '2026-03-01T00:00:00Z',
      updated_at: '2026-04-01T00:00:00Z',
      role: null,
    });
  });
});
