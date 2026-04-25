# Quality Review Report

Date: 2026-04-25

Reviewer: Akira Okutomi

Repository: `a2a-attested-tls-case-study`

## Scope

This review covers the public case-study repository layout, the bundled
Tamarin and ProVerif models, the top-level verification documentation, and the
PDF report source after the intra-handshake / post-handshake repository split.

The goal of this review was not to claim a full verification of Cocos AI aTLS.
It was to check whether the public artifact repository accurately describes the
models it contains and whether the load-bearing verification claims can be
rerun without obvious semantic or documentation mismatches.

## Review Findings Checked

### 1. ProVerif challenge-channel shadowing

Files checked:

- `verification/post-handshake/proverif/cocos_relay_attack.pv`
- `verification/post-handshake/proverif/cocos_same_endpoint_current.pv`

Issue checked:

The challenge input previously rebound `chan_id_from_client`, which allowed
the honest-server same-channel witness to drift away from the channel learned
from the exported-authenticator request.

Resolution:

Both post-handshake ProVerif models now match the challenge channel against the
previously learned request channel with `=chan_id_from_client`.

Result:

The ProVerif runner completes without the previous identifier-rebound warning.

### 2. Tamarin mitigation-helper lemma documentation

Files checked:

- `verification/README.md`
- `verification/post-handshake/README.md`
- `verification/post-handshake/tamarin/cocos_agent_identity_bound_current.spthy`

Issue checked:

The mitigation sketch reports that the direct same-machine wrong-agent trace is
not found, but the helper correspondence lemma
`acceptance_requires_intended_agent_response` is still falsified.

Resolution:

The top-level and post-handshake verification READMEs now explicitly document
that split result.

Result:

Readers rerunning `verification/run-tamarin.sh` should no longer see a mismatch
between the documented expected highlights and the actual Tamarin summary.

### 3. Intra-handshake / post-handshake repository split

Files checked:

- `verification/intra-handshake/`
- `verification/post-handshake/`
- `verification/run-tamarin.sh`
- `verification/run-proverif.sh`
- `README.md`
- `verification/README.md`

Issue checked:

External review requested a clearer separation between the older
intra-handshake design artifacts and the current post-handshake
exported-authenticator artifacts.

Resolution:

Verification artifacts are now grouped first by design generation:

- `verification/intra-handshake/`
- `verification/post-handshake/`

The wrapper scripts and README files were updated to use the new paths.

### 4. Terminology cleanup

Files checked:

- `README.md`
- `verification/`
- `reports/cocos_atls_verification_report.tex`

Issue checked:

The previous wording used terms that were easy to overread: one exporter-label
term sounded RFC-defined, one attack label mixed two different attack classes,
and one proof-strength term suggested a result that the report does not prove.

Resolution:

The public repository now uses narrower wording:

- `expected/default Cocos exporter label`
- `relay-style same-endpoint`

The searched repository text now avoids those overbroad terms.

## Commands Rerun

The following commands were rerun from the repository root:

```sh
./verification/run-proverif.sh
./verification/run-tamarin.sh
./reports/build-cocos-atls-report.sh
```

Observed result:

- ProVerif completed successfully.
- Tamarin completed successfully.
- The PDF report rebuilt successfully.

Remaining LaTeX messages are typography-only overfull/underfull box messages.
No stale old-path references, undefined references, or structural LaTeX
warnings remained after the final build.

## Residual Caution

The models remain intentionally small. In particular:

- The relay-style same-endpoint results depend on the explicit leakage
  assumption stated in the report.
- The intended-agent-bound Tamarin model is still a mitigation sketch, not a
  complete protocol design.
- Implementation-level Go regressions are discussed in the report but are not
  repackaged in this public artifact repository.

This review confirms that the public artifact repository is internally more
consistent after the cleanup, not that it provides a full end-to-end proof of
the complete implementation.
