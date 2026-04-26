# a2a-attested-tls-case-study

This repository packages a small machine-checked case study on the Cocos AI
aTLS legacy/current designs for sharing in the context of recent discussion on
the IETF agent2agent list.

The bundle is intentionally scoped. It is meant to make the report and the
verification artifacts easy to inspect and rerun without asking readers to
navigate a larger development fork.

## What is here

- `reports/`
  the PDF report, its LaTeX source, and the local build script
- `verification/intra-handshake/`
  legacy intra-handshake Tamarin and ProVerif models
- `verification/post-handshake/`
  current post-handshake exported-authenticator Tamarin and ProVerif models
- `verification/go-regressions/`
  current post-handshake Go-language regression tests run against a Cocos
  source checkout

## What this repository is for

- sharing the machine-checked case study itself
- sharing the report as a stable PDF artifact
- rerunning the Tamarin, ProVerif, and Go-language checks included in the case
  study

## What this repository is not

- a full proof of the complete Cocos AI implementation
- a standalone clone of the full Cocos AI codebase
- a complete copy of the implementation-level Cocos test suite

The Go-language regressions included here are focused current-design tests.
They run against an external Cocos source checkout through `COCOS_SOURCE`.

## Requirements

- Tamarin Prover
- ProVerif
- Go toolchain
- Cocos source checkout for the Go-language regressions
- a LaTeX environment for rebuilding the PDF report

## Quick start

Build the PDF:

```sh
./reports/build-cocos-atls-report.sh
```

Run the Tamarin models:

```sh
./verification/run-tamarin.sh
```

Run the ProVerif models:

```sh
./verification/run-proverif.sh
```

Run the current-design Go-language regressions:

```sh
COCOS_SOURCE=/path/to/cocos ./verification/go-regressions/run-go-regressions.sh
```

When this repository is nested inside a Cocos checkout, the script will try the
parent directory automatically.

## Scope notes

- The report compares a legacy intra-handshake design anchor around commit
  `e372cfc` with a current post-handshake design anchor around commit
  `80bf813`, introduced through Cocos PR #582.
- The strongest concrete current-design finding is exporter-label
  binding-parameter confusion: the verifier should enforce the locally
  expected Cocos exporter label and fail closed on mismatch.
- The second implementation-facing result is weaker: attestation-required
  flows should distinguish a valid empty authenticator from attestation
  verified.
- The same-endpoint results are reduced-model proof obligations. They use the
  A1 leakage assumption from the report and do not model the full RFC 9261
  CertificateVerify/Finished path. They should not be read as a claim that the
  complete RFC 9261 or Cocos post-handshake path is vulnerable.
- The Go-language regressions also cover evidence/binder separation, leaf-key
  substitution, and session-scoped context reuse.
- The same-machine agent-identity result is paired with a narrow
  intended-agent-bound mitigation sketch in both Tamarin and ProVerif. It is a
  deployment-level identity obligation, not a relay or diversion model and not
  a Cocos production bug claim.
- The infrastructure-identity result is a narrow L2 diversion-style sketch. It
  separates valid shared measurement from intended infrastructure identity. In
  the report, L2 is split into L2a platform/VM appraisal and L2b intended
  service/deployment identity.
- The report separates relay, diversion, same-machine wrong-agent, and
  binding-parameter confusion as different threat-model categories.
- The L0--L5 notation is local to the report. It is not the same terminology as
  Cocos documentation that uses "Level 2 Binding" for TLS-session binding.
- The Tamarin and ProVerif models were developed with AI assistance and local
  human review, and some details still need deeper manual review. Their
  checker outputs apply to the bundled models, not to a complete model of the
  Cocos implementation.
