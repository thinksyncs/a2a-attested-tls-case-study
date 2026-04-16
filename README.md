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
- `verification/tamarin/`
  Tamarin models for legacy/current gating, attestation presence,
  exporter-label enforcement, same-endpoint attack trace, same-machine
  intended-agent identity, and context reuse
- `verification/proverif/`
  ProVerif models for relay/diversion, same-endpoint authenticity,
  canonical exporter-label enforcement, leaf-key substitution, and compact
  legacy comparison models

## What this repository is for

- sharing the machine-checked case study itself
- sharing the report as a stable PDF artifact
- rerunning the Tamarin and ProVerif models included in the case study

## What this repository is not

- a full proof of the complete Cocos AI implementation
- a standalone clone of the full Cocos AI codebase
- a complete package of implementation-level Go regressions

The implementation-level Go regressions discussed in the report are not
repackaged here; this public repository is intentionally focused on the report
and the machine-checked models.

## Requirements

- Tamarin Prover
- ProVerif
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

## Scope notes

- The report compares a legacy pre-handshake design anchor around commit
  `e372cfc` with a current post-handshake design anchor around commit
  `80bf813`.
- The main current-design results are about same-endpoint authenticity,
  same-machine intended-agent identity, canonical exporter-label
  enforcement, missing attestation after offer, leaf-key substitution, and
  session-scoped context reuse.
- The relay/diversion results are design-level and rely on the explicit
  leakage assumption stated in the report.
