#!/usr/bin/env node

/**
 * Workflow Dispatcher — Simulates a Router event by publishing directly to the WORKFLOW-EXECUTION stream.
 *
 * Usage:
 *   node index.js                    # Send payload.json with auto-generated workflowUUID
 *   node index.js --uuid MY-UUID     # Send with custom workflowUUID
 *   node index.js --file other.json  # Send a different payload file
 *
 * Environment:
 *   NATS_URL      — NATS server URL (default: nats://localhost:4222)
 *   NATS_USER     — NATS username (default: service)
 *   NATS_PASSWORD  — NATS password (default: service_secret)
 */

const { connect, StringCodec } = require('nats');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

async function main() {
  const args = process.argv.slice(2);

  // Parse args
  let payloadFile = path.join(__dirname, 'payload.json');
  let customUUID = null;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--file' && args[i + 1]) {
      payloadFile = path.resolve(args[++i]);
    } else if (args[i] === '--uuid' && args[i + 1]) {
      customUUID = args[++i];
    }
  }

  // Load payload
  if (!fs.existsSync(payloadFile)) {
    console.error(`Payload file not found: ${payloadFile}`);
    process.exit(1);
  }

  const payload = JSON.parse(fs.readFileSync(payloadFile, 'utf-8'));

  // Generate or use custom workflowUUID
  const workflowUUID = customUUID || crypto.randomUUID();
  if (payload.data) {
    payload.data.workflowUUID = workflowUUID;
  }

  // Generate unique eventTrackerId and executionId
  payload.eventTrackerId = crypto.randomUUID();
  payload.executionId = crypto.randomUUID();
  payload.created = new Date().toISOString();

  // Connect to NATS
  const natsUrl = process.env.NATS_URL || 'nats://localhost:4222';
  const natsUser = process.env.NATS_USER || 'service';
  const natsPass = process.env.NATS_PASSWORD || 'service_secret';

  console.log(`Connecting to ${natsUrl}...`);

  const nc = await connect({
    servers: natsUrl,
    user: natsUser,
    pass: natsPass,
  });

  const sc = StringCodec();
  const subject = 'workflow.execution.router';
  const data = JSON.stringify(payload);

  // Publish
  nc.publish(subject, sc.encode(data));
  await nc.flush();

  console.log(`Published to: ${subject}`);
  console.log(`  mode:         ${payload.mode}`);
  console.log(`  instanceId:   ${payload.data?.instanceId || 'N/A'}`);
  console.log(`  workflowUUID: ${workflowUUID}`);
  console.log(`  trackerId:    ${payload.eventTrackerId}`);
  console.log('');
  console.log('Payload:');
  console.log(JSON.stringify(payload, null, 2));

  await nc.close();
  console.log('\nDone.');
}

main().catch((err) => {
  console.error('Error:', err.message);
  process.exit(1);
});
