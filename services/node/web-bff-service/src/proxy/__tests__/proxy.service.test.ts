import { describe, it, expect } from 'vitest';
import { ProxyService } from '../proxy.service';
import { ConfigService } from '@nestjs/config';

describe('ProxyService', () => {
  it('returns list of service names', () => {
    const config = { get: (key: string, def: string) => def } as unknown as ConfigService;
    const service = new ProxyService(config);
    const names = service.getServiceNames();
    expect(names).toContain('identity');
    expect(names).toContain('tenant');
    expect(names).toContain('audit');
    expect(names.length).toBe(9);
  });
});
