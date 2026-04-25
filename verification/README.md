# Verification Workspace

This directory contains the machine-checked verification artifacts bundled with
the public case-study repository.

The artifacts are separated by design generation first:

- `intra-handshake/`
  legacy pre-handshake / intra-handshake artifacts
- `post-handshake/`
  current post-handshake exported-authenticator artifacts

Each generation then contains `tamarin/` and/or `proverif/` subdirectories.

The implementation-level Go regressions discussed in the report are not bundled
here. This public repository keeps the verification models and the report itself
as a standalone artifact set for sharing and rerunning.

## Quick start

Run all Tamarin models:

```sh
verification/run-tamarin.sh
```

Run all ProVerif models:

```sh
verification/run-proverif.sh
```

## Intra-handshake artifacts

See `verification/intra-handshake/README.md`.

These artifacts correspond to the older legacy design where attestation is
carried inside the TLS handshake. They are retained as compact comparison
points for request gating, nonce/public-key binding, and legacy relay /
same-endpoint behavior under the same leakage-style abstraction.

## Post-handshake artifacts

See `verification/post-handshake/README.md`.

These artifacts correspond to the current design where the client sends an
exported-authenticator request after the TLS 1.3 handshake and the server
returns an exported authenticator carrying attestation material.

## Notes

- The scripts expect `tamarin-prover`, `opam`, and `proverif` to be available
  in `PATH`.
- In the Tamarin models, `verified` and `falsified` do not always line up with
  `good` and `bad`. For attack-trace lemmas, the important question is whether
  a trace is found. The report uses `trace found` and `no trace found` wording
  for those rows.
- The intended-agent-bound Tamarin model is a mitigation sketch, not a complete
  proof. Its direct wrong-agent trace is no longer found, but the helper lemma
  `acceptance_requires_intended_agent_response` is still falsified in this
  abstraction.
- If `opam` is not yet initialized in another environment, run:

```sh
opam init --disable-sandboxing -y
```
