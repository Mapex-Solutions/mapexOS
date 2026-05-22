# Bounded Context: PKI

**Service:** mapexVault
**Module path:** `src/modules/pki/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-12

## Purpose

The PKI bounded context owns the platform's CA storage (root + intermediate). It is stateless from a memory perspective — every request decrypts the CA from MongoDB on demand via the existing envelope encryption primitive (shared with the credentials module), uses the material, and discards it immediately. NEVER caches CA material in process memory beyond the lifetime of a single request. CA documents are inserted into the `pkiCertificateAuthorities` collection by the `mongodb-init` container at deploy time; the service itself does no in-process bootstrapping.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|---|---|---|
| CA | Certificate Authority record — one row per `kind` (root or intermediate) | Plugin marketplace authentication |
| Seed JSON | EJSON document per CA produced by `seed-encryptor`, inserted by `mongodb-init` | Credentials refresh schedule |
| Sign-server request | A CN + SANs payload from the broker provisioning script | OAuth token sign |
| Envelope encryption | Master Key wraps a per-record DEK; DEK encrypts payload | TLS transport encryption |

## Published Events (driven — outbound)

None. Cert state propagation rides on `mapexos.fanout.asset.invalidate` (owned by services/assets).

## Consumed Events (driving — inbound)

None. All driving inputs are HTTP internal endpoints.

## Driving Ports (what can call this module)

- HTTP internal API (API-key auth) under `/internal/pki`:
  - `GET /intermediate_ca_bundle` — returns intermediate cert + decrypted priv key (Assets MS on boot)
  - `GET /ca_chain` — root + intermediate cert PEM concatenated (public material)
  - `POST /sign_server` — signs a server cert with the intermediate CA (broker provisioning)

## Driven Ports (what this module requires)

- `repositories.CARepository` — Mongo persistence for `CertificateAuthority` records.
- `ports.EnvelopePort` — wraps the existing envelope primitive shared with the credentials module.
- `ports.X509SignerPort` — generates and signs certs with `crypto/x509`.

## Invariants and Business Rules

- Encrypted fields (`encryptedDEK`, `dekNonce`, `encryptedKey`, `keyNonce`) MUST never be serialised to API responses — the application-layer mapper reads them only to feed the envelope and discards plaintext immediately.
- `IsSystem` is always `true` for PKI rows — defense in depth; should the service ever expose a JWT route in the future, query filters can enforce `isSystem=false` for user listings (PKI rows never appear there).
- Endpoints hit before the operator runs `generate-pki.sh` + seeds Mongo return `ErrCANotBootstrapped`, mapped to HTTP 503 so the Assets MS retry loop keeps trying.
- Intermediate CA private key is decrypted on demand for `GetIntermediateCABundle` (Assets MS at boot) and for `SignServerCert` (broker provisioning). Plaintext key material exists in process memory ONLY for the duration of a single request.
- Seed JSON envelope-encrypts each CA private key with the same Master Key (`CREDENTIAL_MASTER_KEY`) the service uses to decrypt — keys MUST match between generation and runtime or `GET /internal/pki/intermediate_ca_bundle` fails decryption.

## Known Cross-Context Interactions

- **assets/mqttcerts** — fetches the intermediate CA bundle on boot, RAM-caches for high-throughput device cert signing.
- **mapex-mqtt-broker provisioning script** — calls `POST /sign_server` once per deploy to obtain a signed broker server cert.
- **credentials module (same service)** — shares the underlying envelope encryption primitive but has no direct domain coupling.
- **scripts/prebuild/pki/seed-encryptor** — Go helper invoked by `generate-pki.sh` that envelope-encrypts each CA private key and emits the per-CA EJSON files consumed by the `mongodb-init` container.
