import { Injectable, HttpException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

interface ServiceEndpoints {
  [key: string]: string;
}

@Injectable()
export class ProxyService {
  private readonly endpoints: ServiceEndpoints;

  constructor(private config: ConfigService) {
    this.endpoints = {
      'identity': config.get('IDENTITY_SERVICE_URL', 'http://localhost:8081'),
      'tenant': config.get('TENANT_SERVICE_URL', 'http://localhost:8082'),
      'user-role': config.get('USER_ROLE_SERVICE_URL', 'http://localhost:8083'),
      'master-data': config.get('MASTER_DATA_SERVICE_URL', 'http://localhost:8084'),
      'document': config.get('DOCUMENT_SERVICE_URL', 'http://localhost:8085'),
      'notification': config.get('NOTIFICATION_SERVICE_URL', 'http://localhost:8086'),
      'audit': config.get('AUDIT_SERVICE_URL', 'http://localhost:8087'),
      'workflow': config.get('WORKFLOW_SERVICE_URL', 'http://localhost:8089'),
      'rules-engine': config.get('RULES_ENGINE_SERVICE_URL', 'http://localhost:8090'),
    };
  }

  async forward(
    service: string,
    path: string,
    method: string,
    headers: Record<string, string>,
    body?: unknown,
  ): Promise<{ status: number; data: unknown }> {
    const baseUrl = this.endpoints[service];
    if (!baseUrl) {
      throw new HttpException(`Unknown service: ${service}`, 400);
    }

    const url = `${baseUrl}${path}`;
    const fetchHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (headers['x-tenant-id']) fetchHeaders['X-Tenant-Id'] = headers['x-tenant-id'];
    if (headers['x-request-id']) fetchHeaders['X-Request-Id'] = headers['x-request-id'];
    if (headers['authorization']) fetchHeaders['Authorization'] = headers['authorization'];

    const init: RequestInit = { method, headers: fetchHeaders };
    if (body && method !== 'GET') {
      init.body = JSON.stringify(body);
    }

    const res = await fetch(url, init);
    const data = await res.json().catch(() => null);

    return { status: res.status, data };
  }

  getServiceNames(): string[] {
    return Object.keys(this.endpoints);
  }
}
