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
  checks a reduced same-endpoint proof obligation under the A1 leakage
  assumption
- `tamarin/cocos_agent_identity_current.spthy`
  checks same-machine wrong-agent behavior
- `tamarin/cocos_agent_identity_bound_current.spthy`
  sketches an intended-agent-bound mitigation
- `tamarin/cocos_infrastructure_identity_current.spthy`
  checks an L2 diversion-style intended-infrastructure identity gap
- `tamarin/cocos_infrastructure_identity_bound_current.spthy`
  sketches an intended-infrastructure-bound mitigation
- `tamarin/cocos_task_context_current.spthy`
  checks an L4 task/thread/delegation context gap
- `tamarin/cocos_task_context_bound_current.spthy`
  sketches an intended-task-bound mitigation
- `tamarin/cocos_context_reuse_current.spthy`
  checks session-scoped one-shot context use and replay without session tracking

ProVerif:

- `proverif/cocos_relay_attack.pv`
  checks reduced relay behavior under the A1 leakage assumption
- `proverif/cocos_same_endpoint_current.pv`
  checks reduced same-endpoint correspondence
- `proverif/cocos_exporter_label_current.pv`
  checks exporter-label binding-parameter behavior
- `proverif/cocos_leaf_key_substitution_current.pv`
  checks leaf-key substitution resistance
- `proverif/cocos_infrastructure_identity_current.pv`
  checks an L2 diversion-style intended-infrastructure identity gap
- `proverif/cocos_infrastructure_identity_bound_current.pv`
  sketches an intended-infrastructure-bound mitigation
- `proverif/cocos_agent_identity_current.pv`
  checks an L3 same-machine intended-agent identity gap
- `proverif/cocos_agent_identity_bound_current.pv`
  sketches an intended-agent-bound mitigation
- `proverif/cocos_task_context_current.pv`
  checks an L4 intended-task/context binding gap
- `proverif/cocos_task_context_bound_current.pv`
  sketches an intended-task-bound mitigation
- `proverif/CVE-2026-33697-evidence.md`
  notes how the post-handshake model maps to the inspected Cocos code path

## Run

```sh
tamarin-prover --prove verification/post-handshake/tamarin/cocos_attestation.spthy
tamarin-prover --prove=received_attestation_has_server_origin verification/post-handshake/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove=same_endpoint_can_fail_under_leakage verification/post-handshake/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_agent_identity_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_agent_identity_bound_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_infrastructure_identity_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_infrastructure_identity_bound_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_task_context_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_task_context_bound_current.spthy
tamarin-prover --prove verification/post-handshake/tamarin/cocos_context_reuse_current.spthy

eval "$(opam env --switch=default)"
proverif verification/post-handshake/proverif/cocos_relay_attack.pv
proverif verification/post-handshake/proverif/cocos_same_endpoint_current.pv
proverif verification/post-handshake/proverif/cocos_exporter_label_current.pv
proverif verification/post-handshake/proverif/cocos_leaf_key_substitution_current.pv
proverif verification/post-handshake/proverif/cocos_infrastructure_identity_current.pv
proverif verification/post-handshake/proverif/cocos_infrastructure_identity_bound_current.pv
proverif verification/post-handshake/proverif/cocos_agent_identity_current.pv
proverif verification/post-handshake/proverif/cocos_agent_identity_bound_current.pv
proverif verification/post-handshake/proverif/cocos_task_context_current.pv
proverif verification/post-handshake/proverif/cocos_task_context_bound_current.pv
```

Or run the top-level wrappers:

```sh
./verification/run-tamarin.sh
./verification/run-proverif.sh
```

## Expected highlights

- `cocos_attestation.spthy` verifies server-origin, offer-gating, and plain-mode
  separation, while the expected-label and missing-attestation lemmas are
  falsified.
- `cocos_same_endpoint_current.spthy` keeps server-origin but finds the
  same-endpoint trace under the reduced A1 leakage assumption. This is a
  proof-obligation signal, not a claim that the complete RFC 9261 / Cocos
  post-handshake path is vulnerable.
- `cocos_agent_identity_current.spthy` finds a same-machine wrong-agent trace.
- `cocos_agent_identity_bound_current.spthy` is only a mitigation sketch: the
  direct wrong-agent trace is no longer found, but
  `acceptance_requires_intended_agent_response` is still falsified.
- `cocos_infrastructure_identity_current.spthy` finds an L2 diversion-style
  trace where the accepted endpoint has a valid shared measurement but is not
  the intended infrastructure identity.
- `cocos_infrastructure_identity_bound_current.spthy` is a mitigation sketch:
  binding the attestation and live response to the intended infrastructure
  identity removes that direct diversion trace.
- `cocos_infrastructure_identity_current.pv` shows the same L2 split in
  correspondence form: acceptance still implies a valid shared-measurement
  attestation, but not response by the intended infrastructure endpoint.
- `cocos_infrastructure_identity_bound_current.pv` shows the corresponding
  symbolic mitigation shape when the intended infrastructure identity is part
  of the signed attestation consumed by the client.
- `cocos_agent_identity_current.pv` shows the same L3 split in correspondence
  form: acceptance still implies valid machine-level attestation, but not
  response by the intended agent.
- `cocos_agent_identity_bound_current.pv` shows the corresponding symbolic
  mitigation shape when the intended agent identity is part of the signed
  attestation consumed by the client.
- `cocos_task_context_current.spthy` finds an L4 trace where the intended agent
  is accepted for the intended task, but the live response can come from another
  task/thread/delegation context served by the same agent.
- `cocos_task_context_bound_current.spthy` is a mitigation sketch where the
  direct wrong-task trace is no longer found.
- `cocos_task_context_current.pv` shows the same L4 split in correspondence
  form: acceptance still implies valid agent-level attestation, but not response
  in the intended task context.
- `cocos_task_context_bound_current.pv` shows the corresponding symbolic
  mitigation shape when the intended task identity is part of the signed
  attestation consumed by the client and the live response.

## L2--L4 scope

The L2, L3, and L4 models are machine-checked design-obligation sketches. They
are useful because they separate lower-layer attestation success from the next
intended subject:

- L2: valid measurement is not automatically intended infrastructure identity.
- L3: valid machine attestation is not automatically intended agent identity.
- L4: valid agent-level attestation is not automatically intended task,
  thread, or delegation context.

The bound variants show a minimal symbolic shape that removes the corresponding
direct trace or correspondence failure. They do not claim that the current
Cocos production verifier implements those bindings. For that reason, this
repository keeps L2--L4 at the Tamarin/ProVerif level instead of adding
deployment-specific Go-language regressions.

## Claim-strength notes

The strongest implementation-facing finding in this directory is
exporter-label binding-parameter confusion. The missing-attestation result is
weaker and should be read as a policy-level fail-open risk for
attestation-required flows.

The same-endpoint models use the report's A1 leakage assumption: leakage of the
attested leaf or ephemeral private key after attestation issuance. They do not
assume TLS exporter-secret leakage or attestation-root leakage. They also do
not model the full RFC 9261 CertificateVerify/Finished path, so they are best
read as reduced-model proof obligations.

The L0--L5 notation used by the report is local terminology. It is not the same
as Cocos documentation that uses "Level 2 Binding" for TLS-session binding. In
the report, L2 is split into L2a platform/VM appraisal and L2b intended
service/deployment identity.

The Tamarin and ProVerif models were developed with AI assistance and local
human review. Their machine-checker outputs apply to these models, not to a
complete model of the Cocos implementation.
