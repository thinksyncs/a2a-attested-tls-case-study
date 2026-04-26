# Verification Workspace

This directory contains the machine-checked verification artifacts bundled with
the public case-study repository.

The artifacts are separated by design generation first:

- `intra-handshake/`
  legacy intra-handshake artifacts
- `post-handshake/`
  current post-handshake exported-authenticator artifacts
- `go-regressions/`
  current post-handshake Go-language regression tests

Each design-generation directory contains `tamarin/` and/or `proverif/`
subdirectories.

The Go-language regressions are bundled as a small harness under
`verification/go-regressions/`. They are current-design implementation anchors
for exporter-label binding-parameter behavior, attestation-required policy,
evidence/binder separation, leaf-key substitution, and context reuse. They run
against a Cocos source checkout selected with `COCOS_SOURCE`.

## Reading the results

The report separates claim strengths:

- The strongest implementation-facing finding is exporter-label
  binding-parameter confusion.
- The missing-attestation result is a narrower policy-level fail-open risk for
  attestation-required flows.
- The same-endpoint models are reduced-model proof obligations under the A1
  leakage assumption. They do not model the full RFC 9261
  CertificateVerify/Finished path and should not be read as a production-path
  break.
- The L2--L4 identity and task/context models are design-obligation sketches,
  not Cocos production implementation tests.

The L0--L5 notation used in the report is local terminology. It is not the
same as Cocos documentation that uses "Level 2 Binding" for TLS-session
binding.

## Quick start

Run all Tamarin models:

```sh
verification/run-tamarin.sh
```

Run all ProVerif models:

```sh
verification/run-proverif.sh
```

Run the current-design Go-language regressions:

```sh
COCOS_SOURCE=/path/to/cocos verification/go-regressions/run-go-regressions.sh
```

If this public repository is nested inside a Cocos checkout, the wrapper tries
the parent directory automatically.

## Intra-handshake artifacts

See `verification/intra-handshake/README.md`.

These artifacts correspond to the older legacy design where attestation is
carried inside the TLS handshake. They are retained as compact comparison
points for request gating, nonce/public-key binding, and legacy relay /
same-endpoint behavior under the same reduced A1 leakage-style abstraction.

## Post-handshake artifacts

See `verification/post-handshake/README.md`.

These artifacts correspond to the current design where the client sends an
exported-authenticator request after the TLS 1.3 handshake and the server
returns an exported authenticator carrying attestation material.

## Notes

- The scripts expect `tamarin-prover`, `opam`, and `proverif` to be available
  in `PATH`.
- The Go regression wrapper expects `go` in `PATH` and a Cocos source checkout.
- In the Tamarin models, `verified` and `falsified` do not always line up with
  `good` and `bad`. For attack-trace lemmas, the important question is whether
  a trace is found. The report uses `trace found` and `no trace found` wording
  for those rows.
- The intended-agent-bound Tamarin model is a mitigation sketch, not a complete
  proof. Its direct wrong-agent trace is no longer found, but the helper lemma
  `acceptance_requires_intended_agent_response` is still falsified in this
  abstraction.
- The intended-infrastructure-bound Tamarin model is also a mitigation sketch.
  It checks the abstract L2 diversion-style distinction between a valid shared
  measurement and the intended infrastructure identity. The report separates
  L2a platform/VM appraisal from L2b intended service/deployment identity.
- The matching ProVerif L2 models express the same point as correspondence
  queries: valid shared-measurement attestation is not the same as intended
  infrastructure identity unless that identity is explicitly bound.
- The matching ProVerif L3 models express the agent-identity point as
  correspondence queries: valid machine-level attestation is not the same as
  intended-agent response unless that agent identity is explicitly bound.
- The matching Tamarin and ProVerif L4 models express the task/context point:
  valid agent-level attestation is not the same as intended task, thread, or
  delegation context unless that context is explicitly bound.
- The L2--L4 models are design-obligation checks, not Cocos production
  implementation tests. They are intentionally kept at the Tamarin/ProVerif
  level because concrete tests would depend on deployment policy, evidence
  fields, agent routing, task identifiers, and authorization semantics outside
  this small verification workspace.
- The Tamarin and ProVerif models were developed with AI assistance and local
  human review, and some details still need deeper manual review. Their
  machine-checker outputs apply to these models, not to a complete model of the
  Cocos implementation.
- If `opam` is not yet initialized in another environment, run:

```sh
opam init --disable-sandboxing -y
```
