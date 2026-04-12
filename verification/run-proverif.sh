#!/bin/zsh
set -euo pipefail

cd "$(dirname "$0")/.."

eval "$(opam env --switch=default)"

echo "== cocos_relay_attack.pv =="
proverif verification/proverif/cocos_relay_attack.pv

echo
echo "== cocos_same_endpoint_current.pv =="
proverif verification/proverif/cocos_same_endpoint_current.pv

echo
echo "== cocos_exporter_label_current.pv =="
proverif verification/proverif/cocos_exporter_label_current.pv

echo
echo "== cocos_leaf_key_substitution_current.pv =="
proverif verification/proverif/cocos_leaf_key_substitution_current.pv

echo
echo "== cocos_legacy_request_binding.pv =="
proverif verification/proverif/cocos_legacy_request_binding.pv

echo
echo "== cocos_legacy_relay_attack.pv =="
proverif verification/proverif/cocos_legacy_relay_attack.pv
