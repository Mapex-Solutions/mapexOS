import { connect, type JetStreamClient, type NatsConnection } from "nats";
import type { NatsClient, NatsConnectionOptions } from "./types";

/**
 * Creates a NATS client with JetStream capabilities.
 *
 * @param options - The connection options for establishing a NATS connection.
 * @returns A promise that resolves to a NatsClient object containing the NATS connection and JetStream client.
 */
export async function createNatsClient(options: NatsConnectionOptions): Promise<NatsClient> {
  const nc: NatsConnection = await connect(options);
  const js: JetStreamClient = nc.jetstream();
  return { nc, js };
}

/**
 * Closes the NATS client connection if it is not already closed.
 *
 * @param c - The NatsClient object containing the NATS connection to be closed.
 * @returns A promise that resolves when the connection is successfully closed.
 */
export async function closeNatsClient(c: NatsClient): Promise<void> {
  if (!c?.nc?.isClosed()) await c.nc.close();
}
