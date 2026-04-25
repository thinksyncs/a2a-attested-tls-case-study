# Post-Handshake Artifacts

This directory contains the current post-handshake exported-authenticator
models.

In this generation, the client first completes the TLS 1.3 handshake, then
sends an exported-authenticator request with a request context. The server
returns a post-handshake exported authenticator carrying certificate and
attestation material.

## Files

Tamarin:

- `tamarin/cocos_attestation.spthy`
  checks offer gating, attestation presence, and exporter-label behavior
- `tamarin/cocos_same_endpoint_current.spthy`
  checks same-endpoint behavior under the explicit leakage assumption
- `tamarin/cocos_agent_identity_current.spthy`
  checks same-machine wrong-agent behavior
- `tamarin/cocos_agent_identity_bound_current.spthy`
  sketches an intended-agent-bound mitigation
- `tamarin/cocos_context_reuse_current.spthy`
  checks session-scoped one-shot context use and replay without session tracking

ProVerif:

- `proverif/cocos_relay_attack.pv`
  checks relay behavior under the explicit leakage assumption
- `proverif/cocos_same_endpoint_current.pv`
  checks same-endpoint correspondence
- `proverif/cocos_exporter_label_current.pv`
  checks expected/default Cocos attestation exporter-label behavior
- `proverif/cocos_leaf_key_substitution_current.pv`
  checks leaf-key substitution resistance
- `proverif/CVE-2026-33697-evidence.md`
  notes how the post-handshake model maps to the inspected Cocos code path

## Run

```sh
tamarin-prover --prove verification/post-handshake/tamarin/cocos_attestation.spthy
tamarin-prover --prove=received_attestation_has_server_origin verification/post-handshake/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove=same_endpoint_can_fail_under_leakage verification/post-handshake/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_agent_identity_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_agent_identity_bound_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_context_reuse_current.spthy

eval "$(opam env --switch=default)"
proverif verification/post-handshake/proverif/cocos_relay_attack.pv
proverif verification/post-handshake/proverif/cocos_same_endpoint_current.pv
proverif verification/post-handshake/proverif/cocos_exporter_label_current.pv
proverif verification/post-handshake/proverif/cocos_leaf_key_substitution_current.pv
```

Or run the top-level wrappers:

```sh
./verification/run-tamarin.sh
./verification/run-proverif.sh
```

## Expected highlights

- `cocos_attestation.spthy` verifies server-origin, offer-gating, and plain-mode
  separation, while the default-label and missing-attestation lemmas are
  falsified.
- `cocos_same_endpoint_current.spthy` keeps server-origin but finds the
  same-endpoint attack trace under the explicit leakage assumption.
- `cocos_agent_identity_current.spthy` finds a same-machine wrong-agent trace.
- `cocos_agent_identity_bound_current.spthy` is only a mitigation sketch: the
  direct wrong-agent trace is no longer found, but
  `acceptance_requires_intended_agent_response` is still falsified.
