# Tamarin models for cocos attested TLS

This directory contains four Tamarin models:

- `verification/tamarin/cocos_attestation.spthy` for the current
  post-handshake exported-authenticator design
- `verification/tamarin/cocos_legacy_attestation.spthy` for the older
  pre-handshake OpenSSL custom-extension design
- `verification/tamarin/cocos_same_endpoint_current.spthy` for the narrower
  same-endpoint question in the current design
- `verification/tamarin/cocos_context_reuse_current.spthy` for one-shot
  request-context behavior in the current design

The current-design model was built around the following code paths:

- `pkg/atls/internal_transport/conn.go`
- `pkg/atls/provider.go`
- `pkg/atls/ea/authenticator.go`
- `pkg/atls/eaattestation/binding.go`
- `pkg/atls/eaattestation/verify.go`

## What the model captures

The model is intentionally abstract. It focuses on the security-relevant
message flow, not on byte-level TLS details:

1. A client sends an authenticator request with a fresh context.
2. The client either offers or does not offer the attestation extension.
3. The server only produces an attested authenticator when the offer is present.
4. The attested authenticator contains:
   - the request context,
   - the server public key,
   - an attestation payload,
   - a binder tied to the public key and the exported TLS value.
5. The client accepts only if the signed authenticator matches the offered
   context and binder structure.

This mirrors the current code-level design:

- `CMWAttestationOfferExtension()` from `pkg/atls/ea/cmw_attestation.go`
- `BuildLeafExtensions()` from `pkg/atls/provider.go`
- `CreateAuthenticator()` / `ValidateAuthenticatorWithAttestation()`
- `ComputeBinding()` / `VerifyPayload()`

The legacy model was built around:

- `pkg/atls/extensions.c`
- `pkg/atls/extensions.h`

and focuses on the older `ClientHello -> EncryptedExtensions -> Certificate`
attestation flow.

## What the model does not yet capture

- Full TLS 1.3 transcript hashing
- X.509 certificate chain validation
- Leaf-only certificate extension placement
- Concrete EAT or CoRIM parsing
- Real evidence verification logic

Those can be added later, but this model is already useful for replay,
agreement, and offer-gating properties.

## Current-design results

`cocos_attestation.spthy` verifies these properties:

- attested acceptance requires a prior attestation offer
- plain requests do not lead to attested acceptance
- accepted authenticators originate from a server creation event

`cocos_attestation.spthy` falsifies these properties:

- accepted attestation must use the default exporter label
- offered requests must not succeed without attestation

This falsified property corresponds to the current implementation accepting the
`binder.exporter_label` value from the payload instead of enforcing the default
attestation label during verification.

The second falsified property corresponds to the request/response logic allowing
an attestation offer to be ignored without forcing the client to abort.

`cocos_same_endpoint_current.spthy` asks a different question. It is a smaller
current-design model aimed at the correspondence:

- does client acceptance imply genuine attestation origin, and
- does it also imply the same live endpoint?

In that model:

- `received_attestation_has_server_origin` is verified
- `same_endpoint_can_fail_under_leakage` has a verified attack trace

This mirrors the same split seen in the `ProVerif` models, but in a form that
can later be extended with stronger order/state constraints.

`cocos_context_reuse_current.spthy` asks a narrower API/usage question:

- with explicit session-style context tracking, is a request context one-shot?
- without that tracking, does a replay trace still exist?

In that model:

- `session_acceptance_has_server_origin` is verified
- `session_context_is_one_shot` is verified
- `no_session_replay_exists` has a verified replay trace

This is intentionally scoped to the `ea.Session` semantics. It should not be
read as a claim that every current-design call path is automatically one-shot;
the implementation only gets that guarantee when a session object is used.

Relevant code:

- `pkg/atls/eaattestation/types.go`
- `pkg/atls/eaattestation/verify.go`
- `pkg/atls/internal_transport/conn.go`
- `pkg/atls/ea/authenticator.go`
- `pkg/atls/transport.go`

## Legacy-design results

`cocos_legacy_attestation.spthy` currently checks three narrower legacy
properties:

- attested legacy acceptance requires a prior evidence request
- a `NO_TEE` server response is fail-closed
- accepted legacy attestation binds the client nonce and server public key

This legacy model is intentionally smaller than the current-design model. It
does not yet attempt relay analysis or a full symbolic treatment of the old
OpenSSL handshake state.

## Run the models

Run:

```sh
tamarin-prover --prove verification/tamarin/cocos_attestation.spthy
tamarin-prover --prove verification/tamarin/cocos_legacy_attestation.spthy
tamarin-prover --prove=received_attestation_has_server_origin verification/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove=same_endpoint_can_fail_under_leakage verification/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove verification/tamarin/cocos_context_reuse_current.spthy
```

For `cocos_same_endpoint_current.spthy`, the two lemmas are run separately on
purpose. That means the non-selected lemma may appear as `analysis incomplete`
in each individual invocation; this is expected and does not mean the selected
lemma failed.

Or open the interactive UI:

```sh
tamarin-prover interactive verification/tamarin/cocos_attestation.spthy
tamarin-prover interactive verification/tamarin/cocos_legacy_attestation.spthy
tamarin-prover interactive verification/tamarin/cocos_same_endpoint_current.spthy
```

## Status in this workspace

`tamarin-prover` is installed in this workspace. Use
`verification/run-tamarin.sh` if you want the shortest
reproduction path for all current models.
