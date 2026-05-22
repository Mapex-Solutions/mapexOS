export interface Publisher {
  publish(subject: string, payload: unknown, headers?: Record<string, string>): Promise<void>;
}

export interface Subscriber {
  /**
   * Subscribe to a subject. Returns a stop() function to cleanly drain/unsubscribe.
   */
  subscribe(options: {
    stream?: string;
    subject: string;
    durable?: string;
    queueGroup?: string;
    pull?: boolean;
    handler: (data: Uint8Array) => Promise<void> | void;
  }): Promise<() => Promise<void>>;
}
