import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { HealthController } from './health/health.controller';
import { ProxyModule } from './proxy/proxy.module';
import { CalendarModule } from './calendar/calendar.module';
import { UsersModule } from './users/users.module';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    ProxyModule,
    CalendarModule,
    UsersModule,
  ],
  controllers: [HealthController],
})
export class AppModule {}
