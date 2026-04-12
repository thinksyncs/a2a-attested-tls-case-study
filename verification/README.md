# Verification Workspace

This directory collects the formal verification artifacts bundled with the
public case-study repository.

## What is here

- `verification/tamarin/`
  Tamarin models for both the legacy pre-handshake and current
  post-handshake designs
- `verification/proverif/`
  ProVerif models for relay/diversion, same-endpoint authenticity,
  canonical exporter-label enforcement, leaf-key substitution, and compact
  legacy/current comparison points

The implementation-level Go regressions discussed in the report are not bundled
here. This public repository keeps the formal models and the report itself as a
standalone artifact set for sharing and rerunning.

## Quick start

### 1. Tamarin

```sh
verification/run-tamarin.sh
```

Expected highlights:

- legacy model:
  - `legacy_acceptance_requires_prior_request` is verified
  - `legacy_no_tee_is_fail_closed` is verified
  - `legacy_attestation_binds_nonce_and_public_key` is verified
- `attested_authenticator_has_server_origin` is verified
- `attested_acceptance_requires_prior_offer` is verified
- `plain_requests_do_not_produce_attested_acceptance` is verified
- `accepted_attestation_must_use_default_exporter_label` is falsified
- `offered_requests_must_not_succeed_without_attestation` is falsified

### 2. ProVerif

```sh
verification/run-proverif.sh
```

Expected highlights:

- `ClientAccepts ==> ServerIssuesAttestation` is `true`
- `ClientAccepts ==> ServerBindsSameChannel` is `false`
- `ClientAccepts ==> ServerUsesCanonicalLabel` is `false`
- `ClientAccepts ==> ServerAttestsLeafKey` is `true`
- `ClientAcceptsLegacy ==> ClientRequestsEvidence` is `true`
- `ClientAcceptsLegacy ==> ServerIssuesLegacyAttestation` is `true`
- `ClientAcceptsLegacy ==> ServerCreatesLegacyReport` is `true`

## Notes

- The scripts expect `tamarin-prover`, `opam`, and `proverif` to be available
  in `PATH`.
- If `opam` is not yet initialized in another environment, run:

```sh
opam init --disable-sandboxing -y
```
