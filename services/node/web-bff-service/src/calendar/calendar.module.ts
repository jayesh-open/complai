import { Module } from '@nestjs/common';
import { CacheModule } from '@nestjs/cache-manager';
import { CalendarController } from './calendar.controller';
import { GstServiceClient } from './clients/gst-service-client';
import { Gstr9ServiceClient } from './clients/gstr9-service-client';
import { TdsServiceClient } from './clients/tds-service-client';
import { ItrServiceClient } from './clients/itr-service-client';

@Module({
  imports: [
    CacheModule.register({
      ttl: 5 * 60 * 1000,
      max: 200,
    }),
  ],
  controllers: [CalendarController],
  providers: [GstServiceClient, Gstr9ServiceClient, TdsServiceClient, ItrServiceClient],
})
export class CalendarModule {}
