#!/bin/zsh
set -euo pipefail

cd "$(dirname "$0")/.."

echo "== cocos_attestation.spthy =="
tamarin-prover --prove verification/tamarin/cocos_attestation.spthy

echo
echo "== cocos_legacy_attestation.spthy =="
tamarin-prover --prove verification/tamarin/cocos_legacy_attestation.spthy

echo
echo "== cocos_same_endpoint_current.spthy =="
echo "[note] these two runs target one lemma at a time; the non-selected lemma may appear as 'analysis incomplete' in each individual run."
tamarin-prover --prove=received_attestation_has_server_origin verification/tamarin/cocos_same_endpoint_current.spthy
tamarin-prover --prove=same_endpoint_can_fail_under_leakage verification/tamarin/cocos_same_endpoint_current.spthy

echo
echo "== cocos_context_reuse_current.spthy =="
tamarin-prover --prove verification/tamarin/cocos_context_reuse_current.spthy
