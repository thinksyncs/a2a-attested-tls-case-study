# Go-Language Regression Tests

This directory contains focused Go-language regression tests for the current
post-handshake exported-authenticator path.

The tests are not a copy of the full Cocos test suite. They are the concrete
implementation anchors used by the report for:

- self-declared non-default exporter-label acceptance,
- missing attestation after an explicit offer,
- API-level separation between evidence verification and binder verification,
- leaf-key substitution resistance,
- session-scoped one-shot request contexts, and
- replayable API-level validation without session tracking.

## Run

From the repository root:

```sh
COCOS_SOURCE=/path/to/cocos ./verification/go-regressions/run-go-regressions.sh
```

If this public repository is nested inside a Cocos checkout, the script tries
the parent directory automatically.

The wrapper creates a temporary Go module, points
`github.com/ultravioletrs/cocos` at `COCOS_SOURCE`, and runs the tests in
`current/authenticator_regression_test.go`.

## Scope

These are current-design implementation regressions. They are not Tamarin or
ProVerif models, and they do not cover the legacy intra-handshake design.
