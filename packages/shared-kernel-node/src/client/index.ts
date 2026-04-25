import type { ApiResponse } from '../types/index.js';

// ---------------------------------------------------------------------------
// API client configuration
// ---------------------------------------------------------------------------

export interface ApiClientConfig {
  baseUrl: string;
  getToken: () => Promise<string>;
  tenantId: string;
}

// ---------------------------------------------------------------------------
// API client interface
// ---------------------------------------------------------------------------

export interface ApiClient {
  get<T>(path: string): Promise<ApiResponse<T>>;
  post<T>(path: string, body: unknown): Promise<ApiResponse<T>>;
  put<T>(path: string, body: unknown): Promise<ApiResponse<T>>;
  delete<T>(path: string): Promise<ApiResponse<T>>;
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

export function createApiClient(config: ApiClientConfig): ApiClient {
  const { baseUrl, getToken, tenantId } = config;

  async function request<T>(
    method: string,
    path: string,
    body?: unknown,
  ): Promise<ApiResponse<T>> {
    const token = await getToken();
    const requestId = crypto.randomUUID();
    const start = performance.now();

    const headers: Record<string, string> = {
      Authorization: `Bearer ${token}`,
      'X-Tenant-Id': tenantId,
      'X-Request-Id': requestId,
      'Content-Type': 'application/json',
      Accept: 'application/json',
    };

    const res = await fetch(`${baseUrl}${path}`, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    const latencyMs = Math.round(performance.now() - start);

    if (!res.ok) {
      const errorBody = await res.text();
      throw new Error(
        `API ${method} ${path} failed (${res.status}): ${errorBody}`,
      );
    }

    const data = (await res.json()) as T;

    return {
      data,
      meta: { requestId, latencyMs },
    };
  }

  return {
    get<T>(path: string) {
      return request<T>('GET', path);
    },
    post<T>(path: string, body: unknown) {
      return request<T>('POST', path, body);
    },
    put<T>(path: string, body: unknown) {
      return request<T>('PUT', path, body);
    },
    delete<T>(path: string) {
      return request<T>('DELETE', path);
    },
  };
}
