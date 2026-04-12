# ProVerif Models

This directory contains small `ProVerif` models for the Cocos AI aTLS designs.
The main focus is still the current post-handshake design, with two compact
legacy comparison models for the older pre-handshake generation.

## Why this exists

The earlier Tamarin model in `verification/tamarin/cocos_attestation.spthy`
primarily checked attestation offer handling and exporter-label handling. It
did not model:

- ephemeral key compromise,
- two distinct channels/sessions,
- or the distinction between "proof of possession of the attested key" and
  "same-channel communication with the genuine attested endpoint".

That is why it did not reproduce the published CVE.

## Models

`cocos_relay_attack.pv` is intentionally minimal. It models:

- a client-side exported-authenticator request with a fresh request context,
- a server creating an authenticator that contains an attestation payload over
  `(service identity, ephemeral public key, request context)`,
- a client accepting a peer that can present that authenticator and answer a
  live challenge with the corresponding ephemeral private key,
- and leakage of that ephemeral private key after attestation.

The key point is that the attestation payload does not contain any explicit
representation of the client's channel/session identity.

`cocos_same_endpoint_current.pv` is a companion model with the same narrow
purpose, but framed directly around the question:

- does client acceptance imply genuine attestation origin, and
- does it also imply the same live endpoint?

This file is meant to be the practical first step before adding order/state
constraints in Tamarin.

`cocos_exporter_label_current.pv` is a separate minimal model for the
exporter-label question. It checks whether:

- client acceptance still implies genuine server-side attestation origin, and
- client acceptance also implies use of the canonical exporter label.

This model is intentionally narrow and does not attempt to cover relay or
same-endpoint behavior.

`cocos_leaf_key_substitution_current.pv` is a separate minimal model for
leaf-key substitution resistance. It checks whether client acceptance still
implies that the attested leaf public key is the same key used for the live
authenticator proof.

`cocos_legacy_request_binding.pv` is a compact legacy comparison model. It
checks whether client acceptance implies:

- a prior evidence request carrying the client nonce, and
- a server-issued report over that same nonce and the server leaf public key.

It is intentionally smaller than the current-design ProVerif models and is
meant only to give the legacy generation one matching `ProVerif` anchor in the
workspace.

`cocos_legacy_relay_attack.pv` is a second legacy comparison model. It keeps
the same pre-handshake nonce/report shape but adds:

- a live challenge on the client side,
- proof of possession of the attested leaf key,
- and an explicit leakage assumption for that leaf private key.

Its purpose is to ask the legacy-side version of the same relay / same-endpoint
question used in the current-design models.

## Mapping to CoCoS code

The abstraction is meant to track the concrete flow in CoCoS:

- client sends an exported-authenticator request:
  `pkg/atls/internal_transport/conn.go`
- server builds the CMW attestation leaf extension from the TLS connection
  state, request context, and leaf certificate:
  `pkg/atls/provider.go`
- the binding uses `ExportKeyingMaterial(label, context)` plus the leaf public
  key:
  `pkg/atls/eaattestation/binding.go`
- client verifies the authenticator and then verifies the attestation binder:
  `pkg/atls/ea/authenticator.go`
  `pkg/atls/eaattestation/verify.go`

## Expected results

Run:

```sh
eval "$(opam env --switch=default)"
proverif verification/proverif/cocos_relay_attack.pv
proverif verification/proverif/cocos_same_endpoint_current.pv
proverif verification/proverif/cocos_exporter_label_current.pv
proverif verification/proverif/cocos_leaf_key_substitution_current.pv
proverif verification/proverif/cocos_legacy_request_binding.pv
proverif verification/proverif/cocos_legacy_relay_attack.pv
```

For the relay / same-endpoint models, you should see:

- acceptance still implies a prior exported-authenticator request:
  `ClientAccepts ==> ClientSendsEARequest`
- the attestation-origin correspondence succeed:
  `ClientAccepts ==> ServerIssuesAttestation`
- the server-side authenticator-build correspondence fail under relay:
  `ClientAccepts ==> ServerBuildsAuthenticator`
- the same-endpoint / same-channel correspondence fail:
  `ClientAccepts ==> ServerBindsSameChannel`
  or
  `ClientAccepts ==> ServerControlsLiveEndpoint`

That is the minimal signal we wanted: a client can accept a peer based on a
genuine attestation and key possession, without a proof that the client's live
channel terminates at the genuine attested service. It also shows a useful
split: client acceptance still implies a prior client-side request, but does
not imply that the same channel identifier was used on the server-side
authenticator build event under the relay/leakage scenario.

For the exporter-label model, you should see:

- the attestation-origin correspondence succeed:
  `ClientAccepts ==> ServerIssuesAttestation`
- the canonical-label correspondence fail:
  `ClientAccepts ==> ServerUsesCanonicalLabel`

That is the matching `ProVerif` signal for the exporter-label question already
seen in Tamarin and in the Go regression test.

For the leaf-key substitution model, you should see:

- `ClientAccepts ==> ServerAttestsLeafKey` succeed

That is the minimal positive signal we want here: within this abstraction, a
client-accepted authenticator stays tied to the attested leaf key rather than
allowing a substituted key.

For the legacy comparison model, you should see:

- `ClientAcceptsLegacy ==> ClientRequestsEvidence` succeed
- `ClientAcceptsLegacy ==> ServerIssuesLegacyAttestation` succeed
- `ClientAcceptsLegacy ==> ServerCreatesLegacyReport` succeed

That is the legacy-side signal we want here: acceptance remains gated by a
prior request and tied to the nonce/public-key report used in the older
pre-handshake design, while still giving the legacy side a more directly
comparable attestation-origin query.

For the legacy relay comparison model, you should see:

- `ClientAcceptsLegacy ==> ClientRequestsEvidence` succeed
- `ClientAcceptsLegacy ==> ServerIssuesLegacyAttestation` succeed
- `ClientAcceptsLegacy ==> LegacyServerBindsSameChannel` fail

That gives the legacy generation a CVE-aligned relay/same-endpoint comparison
under an explicit leaf-key leakage assumption: acceptance still tracks genuine
legacy attestation origin, but it does not establish same-endpoint
authenticity.

## Why ProVerif first

For this specific question, `ProVerif` is a good first step because the core
issue is a correspondence question: whether a client acceptance event implies
the same live endpoint, rather than only genuine attestation origin plus key
possession.

If we later want to tighten ordering, reuse prevention, or one-shot context
constraints, that is the point where a follow-up Tamarin model becomes useful.

## Scope

This is still an abstract model. It does not claim to be a full symbolic model
of TLS or of the complete CoCoS implementation. Its purpose is narrower:
capture the design gap described in the CVE more directly than the earlier
Tamarin model did, while keeping the names and data flow close to the CoCoS
code path.
