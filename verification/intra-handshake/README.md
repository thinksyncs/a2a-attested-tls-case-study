# Intra-Handshake Artifacts

This directory contains the legacy intra-handshake models.

In this generation, attestation is modeled as part of the TLS handshake rather
than as a post-handshake exported-authenticator exchange.

## Files

Tamarin:

- `tamarin/cocos_legacy_attestation.spthy`
  checks prior-request gating, `NO_TEE` fail-closed handling, and
  nonce/public-key binding

ProVerif:

- `proverif/cocos_legacy_request_binding.pv`
  checks legacy request/report binding
- `proverif/cocos_legacy_relay_attack.pv`
  checks the legacy-side relay / same-endpoint comparison under the same
  reduced A1 leakage-style abstraction used for comparison with the
  post-handshake models

These legacy models are retained as compact comparison material. They are not
the main subject of the report.

## Run

```sh
tamarin-prover --prove verification/intra-handshake/tamarin/cocos_legacy_attestation.spthy

eval "$(opam env --switch=default)"
proverif verification/intra-handshake/proverif/cocos_legacy_request_binding.pv
proverif verification/intra-handshake/proverif/cocos_legacy_relay_attack.pv
```

Or run the top-level wrappers:

```sh
./verification/run-tamarin.sh
./verification/run-proverif.sh
```
