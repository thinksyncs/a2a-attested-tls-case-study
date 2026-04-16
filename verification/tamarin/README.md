# Tamarin Models

This directory contains the Tamarin models used in the case study.

## Files

- `cocos_attestation.spthy`
  current-design attestation offer/presence and canonical exporter-label checks
- `cocos_legacy_attestation.spthy`
  legacy-design request gating, `NO_TEE`, and nonce/public-key binding checks
- `cocos_same_endpoint_current.spthy`
  current-design same-endpoint authenticity model under leakage
- `cocos_agent_identity_current.spthy`
  current-design same-machine intended-agent identity failure model
- `cocos_agent_identity_bound_current.spthy`
  current-design intended-agent-bound mitigation sketch
- `cocos_context_reuse_current.spthy`
  current-design one-shot vs replayable request-context model

## Run

```sh
tamarin-prover --prove verification/tamarin/cocos_attestation.spthy
tamarin-prover --prove verification/tamarin/cocos_legacy_attestation.spthy
tamarin-prover --prove=received_attestation_has_server_origin verification/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove=same_endpoint_can_fail_under_leakage verification/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove verification/tamarin/cocos_agent_identity_current.spthy
tamarin-prover --prove verification/tamarin/cocos_agent_identity_bound_current.spthy
tamarin-prover --prove verification/tamarin/cocos_context_reuse_current.spthy
```

Or run the wrapper:

```sh
./verification/run-tamarin.sh
```

## Expected signals

Current-design attestation model:

- server-origin authenticity is `verified`
- offer-gated acceptance is `verified`
- no attested acceptance from plain mode is `verified`
- canonical exporter-label enforcement is `falsified`
- no success without attestation after offer is `falsified`

Current-design same-endpoint model:

- `received_attestation_has_server_origin` is `verified`
- `same_endpoint_can_fail_under_leakage` yields a verified attack trace

Current-design context reuse model:

- `session_acceptance_has_server_origin` is `verified`
- `session_context_is_one_shot` is `verified`
- `no_session_replay_exists` yields a verified replay trace

Current-design same-machine agent-identity model:

- `received_machine_attestation_has_machine_origin` is `verified`
- `intended_agent_identity_can_fail_on_same_machine` yields a verified attack trace

Current-design intended-agent-bound mitigation sketch:

- `received_bound_attestation_has_machine_origin` is `verified`
- `wrong_agent_identity_can_fail_on_same_machine` is `falsified` with no trace found

Legacy-design model:

- prior-request gating is `verified`
- `NO_TEE` fail-closed handling is `verified`
- nonce/public-key binding is `verified`

## How to read the results

The Tamarin side is strongest on gating, presence/absence conditions, and
session-style reuse behavior. In this workspace it also provides a same-endpoint
attack trace under the explicit leakage assumption used by the relay analysis.
In the current same-endpoint model, `ServerControlsLiveEndpoint` is emitted
only when the honest server receives and answers a live challenge on a specific
channel.

The current models therefore show a split:

- accepted attestation still tracks genuine server-side origin
- accepted attestation does not by itself imply same-endpoint authenticity
- machine-level attestation does not by itself imply intended-agent identity
- an intended-agent-bound attestation + live-response sketch removes the
  same-machine wrong-agent trace in this abstraction
- one-shot request-context behavior only holds for the explicit session-tracking
  path

## Scope

These models are intentionally abstract. They focus on protocol structure and
correspondence/gating questions rather than on full TLS transcript details,
full certificate validation, or concrete evidence parsing.
