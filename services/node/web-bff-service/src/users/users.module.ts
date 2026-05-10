import { Module } from '@nestjs/common';
import { CacheModule } from '@nestjs/cache-manager';
import { UsersController } from './users.controller';
import { IdentityServiceClient } from './clients/identity-service-client';
import { UserRoleServiceClient } from './clients/user-role-service-client';

@Module({
  imports: [
    CacheModule.register({
      ttl: 10 * 60 * 1000,
      max: 100,
    }),
  ],
  controllers: [UsersController],
  providers: [IdentityServiceClient, UserRoleServiceClient],
})
export class UsersModule {}
