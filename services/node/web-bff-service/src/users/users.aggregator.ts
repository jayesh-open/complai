import { Logger } from '@nestjs/common';
import type { UserRoleServiceClient } from './clients/user-role-service-client';
import type { User, UserWithRole } from './users.types';

const logger = new Logger('UsersAggregator');

export async function enrichUsersWithRoles(
  users: User[],
  tenantId: string,
  userRoleClient: UserRoleServiceClient,
): Promise<UserWithRole[]> {
  // N+1 pattern — batch optimization deferred to Part 14
  return Promise.all(
    users.map(async (user) => {
      try {
        const roles = await userRoleClient.getUserRoles(tenantId, user.id);
        const primary = roles.length > 0 ? roles[0]! : null;
        return {
          id: user.id,
          email: user.email,
          first_name: user.first_name,
          last_name: user.last_name,
          status: user.status,
          created_at: user.created_at,
          updated_at: user.updated_at,
          role: primary
            ? { id: primary.id, name: primary.name, display_name: primary.display_name }
            : null,
        };
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        logger.warn(`Failed to fetch roles for user ${user.id}: ${msg}`);
        return {
          id: user.id,
          email: user.email,
          first_name: user.first_name,
          last_name: user.last_name,
          status: user.status,
          created_at: user.created_at,
          updated_at: user.updated_at,
          role: null,
        };
      }
    }),
  );
}
