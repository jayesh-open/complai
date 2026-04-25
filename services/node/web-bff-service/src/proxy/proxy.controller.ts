import { Controller, All, Req, Res, Param } from '@nestjs/common';
import { Request, Response } from 'express';
import { ProxyService } from './proxy.service';

@Controller('api/v1')
export class ProxyController {
  constructor(private readonly proxy: ProxyService) {}

  @All(':service/*path')
  async proxyRequest(
    @Param('service') service: string,
    @Param('path') path: string,
    @Req() req: Request,
    @Res() res: Response,
  ) {
    const result = await this.proxy.forward(
      service,
      `/v1/${path}`,
      req.method,
      req.headers as Record<string, string>,
      req.body,
    );
    res.status(result.status).json(result.data);
  }
}
