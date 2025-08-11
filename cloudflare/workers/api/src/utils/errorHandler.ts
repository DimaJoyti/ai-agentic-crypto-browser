/**
 * Error handling utilities for Cloudflare Worker
 */

import { corsHeaders } from './cors';

export interface APIError {
  error: string;
  message: string;
  code?: string;
  details?: any;
  timestamp: string;
}

export function errorHandler(error: any): Response {
  console.error('Worker error:', error);

  let status = 500;
  let errorResponse: APIError = {
    error: 'Internal Server Error',
    message: 'An unexpected error occurred',
    timestamp: new Date().toISOString(),
  };

  // Handle specific error types
  if (error instanceof Error) {
    errorResponse.message = error.message;
    
    // Handle specific error types
    if (error.name === 'ValidationError') {
      status = 400;
      errorResponse.error = 'Bad Request';
    } else if (error.name === 'UnauthorizedError') {
      status = 401;
      errorResponse.error = 'Unauthorized';
    } else if (error.name === 'ForbiddenError') {
      status = 403;
      errorResponse.error = 'Forbidden';
    } else if (error.name === 'NotFoundError') {
      status = 404;
      errorResponse.error = 'Not Found';
    } else if (error.name === 'RateLimitError') {
      status = 429;
      errorResponse.error = 'Too Many Requests';
    }
  }

  return new Response(JSON.stringify(errorResponse), {
    status,
    headers: {
      'Content-Type': 'application/json',
      ...corsHeaders,
    },
  });
}

export class ValidationError extends Error {
  constructor(message: string, public details?: any) {
    super(message);
    this.name = 'ValidationError';
  }
}

export class UnauthorizedError extends Error {
  constructor(message: string = 'Unauthorized') {
    super(message);
    this.name = 'UnauthorizedError';
  }
}

export class ForbiddenError extends Error {
  constructor(message: string = 'Forbidden') {
    super(message);
    this.name = 'ForbiddenError';
  }
}

export class NotFoundError extends Error {
  constructor(message: string = 'Not Found') {
    super(message);
    this.name = 'NotFoundError';
  }
}

export class RateLimitError extends Error {
  constructor(message: string = 'Rate limit exceeded') {
    super(message);
    this.name = 'RateLimitError';
  }
}
