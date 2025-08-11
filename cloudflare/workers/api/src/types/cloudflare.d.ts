/**
 * Cloudflare Workers type declarations
 * These types are provided by the Cloudflare Workers runtime
 */

declare global {
  // Cloudflare Workers global types
  interface KVNamespace {
    get(key: string, options?: { type?: 'text' | 'json' | 'arrayBuffer' | 'stream' }): Promise<string | null>;
    put(key: string, value: string | ArrayBuffer | ArrayBufferView | ReadableStream, options?: {
      expirationTtl?: number;
      expiration?: number;
      metadata?: any;
    }): Promise<void>;
    delete(key: string): Promise<void>;
    list(options?: {
      prefix?: string;
      limit?: number;
      cursor?: string;
    }): Promise<{
      keys: Array<{ name: string; expiration?: number; metadata?: any }>;
      list_complete: boolean;
      cursor?: string;
    }>;
  }

  interface D1Database {
    prepare(query: string): D1PreparedStatement;
    dump(): Promise<ArrayBuffer>;
    batch<T = unknown>(statements: D1PreparedStatement[]): Promise<D1Result<T>[]>;
    exec(query: string): Promise<D1ExecResult>;
  }

  interface D1PreparedStatement {
    bind(...values: any[]): D1PreparedStatement;
    first<T = unknown>(colName?: string): Promise<T | null>;
    run(): Promise<D1Result>;
    all<T = unknown>(): Promise<D1Result<T>>;
    raw<T = unknown>(): Promise<T[]>;
  }

  interface D1Result<T = unknown> {
    results?: T[];
    success: boolean;
    error?: string;
    meta: {
      duration: number;
      size_after: number;
      rows_read: number;
      rows_written: number;
    };
  }

  interface D1ExecResult {
    count: number;
    duration: number;
  }

  interface R2Bucket {
    head(key: string): Promise<R2Object | null>;
    get(key: string, options?: R2GetOptions): Promise<R2ObjectBody | null>;
    put(key: string, value: ReadableStream | ArrayBuffer | ArrayBufferView | string | null | Blob, options?: R2PutOptions): Promise<R2Object>;
    delete(key: string | string[]): Promise<void>;
    list(options?: R2ListOptions): Promise<R2Objects>;
  }

  interface R2Object {
    key: string;
    version: string;
    size: number;
    etag: string;
    httpEtag: string;
    uploaded: Date;
    httpMetadata?: R2HTTPMetadata;
    customMetadata?: Record<string, string>;
    range?: R2Range;
  }

  interface R2ObjectBody extends R2Object {
    body: ReadableStream;
    bodyUsed: boolean;
    arrayBuffer(): Promise<ArrayBuffer>;
    text(): Promise<string>;
    json<T = unknown>(): Promise<T>;
    blob(): Promise<Blob>;
  }

  interface R2GetOptions {
    onlyIf?: R2Conditional;
    range?: R2Range;
  }

  interface R2PutOptions {
    onlyIf?: R2Conditional;
    httpMetadata?: R2HTTPMetadata;
    customMetadata?: Record<string, string>;
    md5?: ArrayBuffer | string;
    sha1?: ArrayBuffer | string;
    sha256?: ArrayBuffer | string;
    sha384?: ArrayBuffer | string;
    sha512?: ArrayBuffer | string;
  }

  interface R2ListOptions {
    limit?: number;
    prefix?: string;
    cursor?: string;
    delimiter?: string;
    startAfter?: string;
    include?: ('httpMetadata' | 'customMetadata')[];
  }

  interface R2Objects {
    objects: R2Object[];
    truncated: boolean;
    cursor?: string;
    delimitedPrefixes: string[];
  }

  interface R2HTTPMetadata {
    contentType?: string;
    contentLanguage?: string;
    contentDisposition?: string;
    contentEncoding?: string;
    cacheControl?: string;
    cacheExpiry?: Date;
  }

  interface R2Range {
    offset?: number;
    length?: number;
    suffix?: number;
  }

  interface R2Conditional {
    etagMatches?: string;
    etagDoesNotMatch?: string;
    uploadedBefore?: Date;
    uploadedAfter?: Date;
  }

  interface DurableObjectNamespace {
    newUniqueId(options?: { jurisdiction?: string }): DurableObjectId;
    idFromName(name: string): DurableObjectId;
    idFromString(id: string): DurableObjectId;
    get(id: DurableObjectId): DurableObjectStub;
  }

  interface DurableObjectId {
    toString(): string;
    equals(other: DurableObjectId): boolean;
  }

  interface DurableObjectStub {
    id: DurableObjectId;
    fetch(input: RequestInfo, init?: RequestInit): Promise<Response>;
  }

  interface DurableObjectState {
    waitUntil(promise: Promise<any>): void;
    id: DurableObjectId;
    storage: DurableObjectStorage;
    blockConcurrencyWhile<T>(callback: () => Promise<T>): Promise<T>;
    acceptWebSocket(ws: WebSocket, tags?: string[]): void;
    getWebSockets(tag?: string): WebSocket[];
    setWebSocketAutoResponse(maybeReqResp?: WebSocketRequestResponsePair): void;
    getWebSocketAutoResponse(): WebSocketRequestResponsePair | null;
    getWebSocketAutoResponseTimestamp(ws: WebSocket): Date | null;
  }

  interface DurableObjectStorage {
    get<T = unknown>(key: string, options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<T | undefined>;
    get<T = unknown>(keys: string[], options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<Map<string, T>>;
    put<T>(key: string, value: T, options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<void>;
    put<T>(entries: Record<string, T>, options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<void>;
    delete(key: string, options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<boolean>;
    delete(keys: string[], options?: { allowConcurrency?: boolean; noCache?: boolean }): Promise<number>;
    list<T = unknown>(options?: {
      start?: string;
      startAfter?: string;
      end?: string;
      prefix?: string;
      reverse?: boolean;
      limit?: number;
      allowConcurrency?: boolean;
      noCache?: boolean;
    }): Promise<Map<string, T>>;
    transaction<T>(closure: (txn: DurableObjectTransaction) => Promise<T>): Promise<T>;
    getAlarm(options?: { allowConcurrency?: boolean }): Promise<number | null>;
    setAlarm(scheduledTime: number | Date, options?: { allowConcurrency?: boolean }): Promise<void>;
    deleteAlarm(options?: { allowConcurrency?: boolean }): Promise<void>;
    sync(): Promise<void>;
  }

  interface DurableObjectTransaction {
    get<T = unknown>(key: string): Promise<T | undefined>;
    get<T = unknown>(keys: string[]): Promise<Map<string, T>>;
    put<T>(key: string, value: T): Promise<void>;
    put<T>(entries: Record<string, T>): Promise<void>;
    delete(key: string): Promise<boolean>;
    delete(keys: string[]): Promise<number>;
    list<T = unknown>(options?: {
      start?: string;
      startAfter?: string;
      end?: string;
      prefix?: string;
      reverse?: boolean;
      limit?: number;
    }): Promise<Map<string, T>>;
    rollback(): void;
  }

  interface WebSocketRequestResponsePair {
    request: string;
    response: string;
  }

  interface ExecutionContext {
    waitUntil(promise: Promise<any>): void;
    passThroughOnException(): void;
  }

  // WebSocket Pair
  class WebSocketPair {
    0: WebSocket;
    1: WebSocket;
  }

  // Cloudflare-specific WebSocket extensions
  interface WebSocket {
    accept(): void;
    send(message: string | ArrayBuffer): void;
    close(code?: number, reason?: string): void;
    addEventListener(type: 'message', listener: (event: MessageEvent) => void): void;
    addEventListener(type: 'close', listener: (event: CloseEvent) => void): void;
    addEventListener(type: 'error', listener: (event: ErrorEvent) => void): void;
    addEventListener(type: 'open', listener: (event: Event) => void): void;
    removeEventListener(type: string, listener: EventListener): void;
    readyState: number;
    url: string;
    protocol: string;
    extensions: string;
    bufferedAmount: number;
    binaryType: string;
  }
}

export {};
