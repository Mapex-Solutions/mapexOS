# Bounded Context: KEK (Key Encryption Keys)

**Last reviewed:** 2026-06-05

## Purpose

Hand out platform key-encryption keys (KEKs) to authorized internal services so
they can envelope-encrypt their own secrets locally, without ever holding the
Vault master key. A KEK is a random AES-256 key stored envelope-encrypted (with
`credential_master_key`) in Mongo and decrypted on demand for the caller. This
mirrors the `pki` bounded context, which hands out the intermediate CA material
the same way.

## Ubiquitous Language

- **KEK**: a key-encryption key. Random 32-byte AES-256 key, stored as its
  64-hex string, envelope-encrypted at rest.
- **Context**: the secret category a KEK protects (e.g. `lorawan_device_keys`).
  One KEK per context, unique. Extensible: new categories add a new context.
- **IsSystem**: always `true` for these rows; a defense-in-depth guard so
  system KEKs never appear in any future user-facing listing.

## Driven Ports

- `EncryptionKeyRepository`: reads `EncryptionKey` records from Mongo by context.
- `EnvelopePort`: decrypts the stored envelope (decrypt only; the seed encrypts).

## Driving Ports

- `KekServicePort.GetKEK(ctx, context)`: load + decrypt the KEK for a context.

## Published / Consumed

- Exposes `GET /internal/kek/:context` (API-key gated). Returns
  `{context, kek}` where `kek` is the decrypted 64-hex key. Consumed by the
  assets MS and the LNS at boot to seed their in-RAM envelope service.
- Status: 200 (ok), 503 (not seeded yet, caller retries), 500 (real failure).

## Invariants

- KEK rows are NEVER authored by this service: the mongodb-init container
  (`kek-bootstrap.sh`) seeds them, encrypting with the same `credential_master_key`
  the service decrypts with.
- The plaintext KEK exists in process memory only for the duration of a request
  and is never logged or persisted.

## Cross-Context Interactions

- Shares the Vault's `*envelope.EnvelopeService` (master key from
  `credential_master_key`) via the shared `bootstrap/encryption.go` provider.
