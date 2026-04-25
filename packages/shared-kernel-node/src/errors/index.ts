// ---------------------------------------------------------------------------
// Error taxonomy (mirrors Go shared-kernel error types)
// ---------------------------------------------------------------------------

export class ComplaiError extends Error {
  readonly code: string;
  readonly statusCode: number;
  readonly details: Record<string, unknown>;

  constructor(
    message: string,
    code: string,
    statusCode: number,
    details: Record<string, unknown> = {},
  ) {
    super(message);
    this.name = 'ComplaiError';
    this.code = code;
    this.statusCode = statusCode;
    this.details = details;
    // Restore prototype chain for proper instanceof checks
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class NotFoundError extends ComplaiError {
  constructor(resource: string, id: string) {
    super(`${resource} not found: ${id}`, 'NOT_FOUND', 404, { resource, id });
    this.name = 'NotFoundError';
  }
}

export class ConflictError extends ComplaiError {
  constructor(message: string, details: Record<string, unknown> = {}) {
    super(message, 'CONFLICT', 409, details);
    this.name = 'ConflictError';
  }
}

export class ValidationError extends ComplaiError {
  constructor(message: string, details: Record<string, unknown> = {}) {
    super(message, 'VALIDATION_ERROR', 400, details);
    this.name = 'ValidationError';
  }
}

export class AuthorizationError extends ComplaiError {
  constructor(message = 'Unauthorized') {
    super(message, 'UNAUTHORIZED', 403);
    this.name = 'AuthorizationError';
  }
}

export class ProviderError extends ComplaiError {
  constructor(
    provider: string,
    message: string,
    details: Record<string, unknown> = {},
  ) {
    super(`Provider ${provider}: ${message}`, 'PROVIDER_ERROR', 502, {
      provider,
      ...details,
    });
    this.name = 'ProviderError';
  }
}

// ---------------------------------------------------------------------------
// Type guard
// ---------------------------------------------------------------------------

export function isComplaiError(err: unknown): err is ComplaiError {
  return err instanceof ComplaiError;
}
