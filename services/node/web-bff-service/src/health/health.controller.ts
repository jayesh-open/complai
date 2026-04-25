import { Controller, Get } from '@nestjs/common';

@Controller()
export class HealthController {
  @Get('health')
  health() {
    return { status: 'ok', service: 'web-bff-service' };
  }

  @Get('ping')
  ping() {
    return 'pong';
  }
}
