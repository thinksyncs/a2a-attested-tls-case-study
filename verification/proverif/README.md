# ProVerif Models

This directory contains the `ProVerif` models used in the case study. They are
small correspondence-style models for the legacy/current Cocos AI aTLS designs.

The current-design models carry most of the analysis. The two legacy models are
kept as compact comparison points.

## Files

- `cocos_relay_attack.pv`
  current-design relay/diversion model under the explicit leakage assumption
- `cocos_same_endpoint_current.pv`
  current-design same-endpoint authenticity model
- `cocos_exporter_label_current.pv`
  current-design canonical exporter-label enforcement model
- `cocos_leaf_key_substitution_current.pv`
  current-design leaf-key substitution model
- `cocos_legacy_request_binding.pv`
  legacy request/report binding comparison model
- `cocos_legacy_relay_attack.pv`
  legacy relay/same-endpoint comparison model

## Run

```sh
eval "$(opam env --switch=default)"
proverif verification/proverif/cocos_relay_attack.pv
proverif verification/proverif/cocos_same_endpoint_current.pv
proverif verification/proverif/cocos_exporter_label_current.pv
proverif verification/proverif/cocos_leaf_key_substitution_current.pv
proverif verification/proverif/cocos_legacy_request_binding.pv
proverif verification/proverif/cocos_legacy_relay_attack.pv
```

Or run the wrapper:

```sh
./verification/run-proverif.sh
```

## Expected signals

Current-design relay / same-endpoint:

- `ClientAccepts ==> ClientSendsEARequest` is `true`
- `ClientAccepts ==> ServerIssuesAttestation` is `true`
- `ClientAccepts ==> ServerBuildsAuthenticator` is `false`
- `ClientAccepts ==> ServerBindsSameChannel` is `false`

Current-design canonical exporter-label enforcement:

- `ClientAccepts ==> ServerIssuesAttestation` is `true`
- `ClientAccepts ==> ServerUsesCanonicalLabel` is `false`

Current-design leaf-key substitution:

- `ClientAccepts ==> ServerAttestsLeafKey` is `true`

Legacy comparison:

- `ClientAcceptsLegacy ==> ClientRequestsEvidence` is `true`
- `ClientAcceptsLegacy ==> ServerIssuesLegacyAttestation` is `true`
- `ClientAcceptsLegacy ==> ServerCreatesLegacyReport` is `true`
- `ClientAcceptsLegacy ==> LegacyServerBindsSameChannel` is `false`

## How to read the results

The current models separate two ideas:

- accepted attestation still tracks genuine server-side attestation origin
- accepted attestation does not by itself establish same-endpoint authenticity

Read together with the compact legacy comparison models, they also show that
moving from the legacy pre-handshake design to the current post-handshake
design does not by itself resolve same-endpoint authenticity under the explicit
leakage assumption used by the relay models.

## Scope

These are abstract models. They are meant to clarify specific binding and
correspondence questions, not to serve as a full symbolic model of TLS or of
the complete Cocos AI implementation.
