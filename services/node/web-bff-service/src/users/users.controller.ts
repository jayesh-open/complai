import {
  Controller,
  Get,
  Post,
  Put,
  Query,
  Param,
  Body,
  Inject,
  BadRequestException,
  ForbiddenException,
  NotFoundException,
  InternalServerErrorException,
  HttpCode,
  Logger,
} from '@nestjs/common';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import type { Cache } from 'cache-manager';
import { IdentityServiceClient } from './clients/identity-service-client';
import { UserRoleServiceClient } from './clients/user-role-service-client';
import { enrichUsersWithRoles } from './users.aggregator';
import type { Role, UserWithRole, PermissionPair } from './users.types';

const ROLE_CACHE_TTL_MS = 10 * 60 * 1000;
const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

@Controller('api/v1')
export class UsersController {
  private readonly logger = new Logger(UsersController.name);

  constructor(
    private readonly identityClient: IdentityServiceClient,
    private readonly userRoleClient: UserRoleServiceClient,
    @Inject(CACHE_MANAGER) private readonly cache: Cache,
  ) {}

  // ---- Users ----

  @Get('users')
  async listUsers(@Query('tenant_id') tenantId?: string) {
    this.requireTenantId(tenantId);

    const users = await this.identityClient.listUsers(tenantId!);
    const enriched = await enrichUsersWithRoles(users, tenantId!, this.userRoleClient);
    return { users: enriched, total: enriched.length };
  }

  @Get('users/:user_id')
  async getUser(
    @Param('user_id') userId: string,
    @Query('tenant_id') tenantId?: string,
  ) {
    this.requireTenantId(tenantId);
    if (!UUID_RE.test(userId)) throw new BadRequestException('Invalid user_id');

    const user = await this.identityClient.getUser(tenantId!, userId);
    if (!user) throw new NotFoundException('User not found');

    const roles = await this.userRoleClient.getUserRoles(tenantId!, userId);
    const primary = roles.length > 0 ? roles[0]! : null;
    const result: UserWithRole = {
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
    return result;
  }

  @Post('users')
  async createUser(
    @Query('tenant_id') tenantId?: string,
    @Body() body?: { email?: string; first_name?: string; last_name?: string; role_id?: string },
  ) {
    this.requireTenantId(tenantId);
    if (!body?.email || !body.first_name || !body.last_name) {
      throw new BadRequestException('email, first_name, last_name are required');
    }

    const user = await this.identityClient.createUser(tenantId!, {
      email: body.email,
      first_name: body.first_name,
      last_name: body.last_name,
    });
    if (!user) {
      throw new InternalServerErrorException('Failed to create user in identity service');
    }

    let roleAssignmentFailed = false;
    if (body.role_id) {
      const ok = await this.userRoleClient.assignRole(tenantId!, user.id, body.role_id);
      if (!ok) {
        this.logger.warn(`Role assignment failed for user ${user.id}, role ${body.role_id}`);
        roleAssignmentFailed = true;
      }
    }

    return { user, role_assignment_failed: roleAssignmentFailed };
  }

  @Put('users/:user_id/role')
  async changeUserRole(
    @Param('user_id') userId: string,
    @Query('tenant_id') tenantId?: string,
    @Body() body?: { role_id?: string },
  ) {
    this.requireTenantId(tenantId);
    if (!UUID_RE.test(userId)) throw new BadRequestException('Invalid user_id');
    if (!body?.role_id) throw new BadRequestException('role_id is required');

    const ok = await this.userRoleClient.assignRole(tenantId!, userId, body.role_id);
    if (!ok) {
      throw new InternalServerErrorException('Failed to assign role');
    }
    return { status: 'role_updated' };
  }

  @Post('users/:user_id/deactivate')
  @HttpCode(200)
  async deactivateUser(
    @Param('user_id') userId: string,
    @Query('tenant_id') tenantId?: string,
  ) {
    this.requireTenantId(tenantId);
    if (!UUID_RE.test(userId)) throw new BadRequestException('Invalid user_id');

    const ok = await this.identityClient.deactivateUser(tenantId!, userId);
    if (!ok) {
      throw new InternalServerErrorException('Failed to deactivate user');
    }
    return { status: 'deactivated' };
  }

  // ---- Roles ----

  @Get('roles')
  async listRoles(@Query('tenant_id') tenantId?: string) {
    this.requireTenantId(tenantId);

    const cacheKey = `roles:list:${tenantId}`;
    const cached = await this.cache.get<Role[]>(cacheKey);
    if (cached) return cached;

    const roles = await this.userRoleClient.listRoles(tenantId!);
    await this.cache.set(cacheKey, roles, ROLE_CACHE_TTL_MS);
    return roles;
  }

  @Get('roles/:role_id')
  async getRoleDetail(
    @Param('role_id') roleId: string,
    @Query('tenant_id') tenantId?: string,
  ) {
    this.requireTenantId(tenantId);
    if (!UUID_RE.test(roleId)) throw new BadRequestException('Invalid role_id');

    const detail = await this.userRoleClient.getRoleDetail(tenantId!, roleId);
    if (!detail) throw new NotFoundException('Role not found');
    return detail;
  }

  @Put('roles/:role_id/permissions')
  async updateRolePermissions(
    @Param('role_id') roleId: string,
    @Query('tenant_id') tenantId?: string,
    @Body() body?: { permission_pairs?: PermissionPair[] },
  ) {
    this.requireTenantId(tenantId);
    if (!UUID_RE.test(roleId)) throw new BadRequestException('Invalid role_id');
    if (!body?.permission_pairs || !Array.isArray(body.permission_pairs)) {
      throw new BadRequestException('permission_pairs array is required');
    }

    const detail = await this.userRoleClient.getRoleDetail(tenantId!, roleId);
    if (!detail) throw new NotFoundException('Role not found');
    if (detail.role.is_system) {
      throw new ForbiddenException('Cannot modify system role permissions');
    }

    // TODO (Part 14): resolve permission_pairs → permission_ids via a permissions lookup/create endpoint
    // For now, permission pair → ID resolution requires direct DB or a new user-role-service endpoint
    throw new BadRequestException(
      'Permission pair resolution not yet available — use direct permission_ids via user-role-service',
    );
  }

  @Post('roles')
  async createRole(
    @Query('tenant_id') tenantId?: string,
    @Body() body?: { name?: string; display_name?: string; description?: string },
  ) {
    this.requireTenantId(tenantId);
    if (!body?.name || !body?.display_name) {
      throw new BadRequestException('name and display_name are required');
    }

    const role = await this.userRoleClient.createRole(tenantId!, {
      name: body.name,
      display_name: body.display_name,
      description: body.description,
    });
    if (!role) {
      throw new InternalServerErrorException('Failed to create role');
    }

    // TODO (Part 14): if template provided, copy permissions from template
    // TODO (Part 14): if permission_pairs provided, resolve and assign
    return role;
  }

  @Post('tenants/:tenant_id/seed-roles')
  @HttpCode(200)
  async seedRoles(@Param('tenant_id') tenantId: string) {
    if (!UUID_RE.test(tenantId)) throw new BadRequestException('Invalid tenant_id');

    const result = await this.userRoleClient.seedRoles(tenantId);
    if (result.conflict) {
      return { status: 'already_seeded', message: 'Tenant roles already seeded' };
    }
    if (!result.seeded) {
      throw new InternalServerErrorException('Failed to seed roles');
    }
    return { status: 'seeded' };
  }

  private requireTenantId(tenantId?: string): asserts tenantId is string {
    if (!tenantId) throw new BadRequestException('tenant_id query parameter is required');
    if (!UUID_RE.test(tenantId)) throw new BadRequestException('Invalid tenant_id format');
  }
}
