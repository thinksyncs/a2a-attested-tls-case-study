#!/bin/zsh
set -euo pipefail

cd "$(dirname "$0")/.."

eval "$(opam env --switch=default)"

echo "== cocos_relay_attack.pv =="
proverif verification/post-handshake/proverif/cocos_relay_attack.pv

echo
echo "== cocos_same_endpoint_current.pv =="
proverif verification/post-handshake/proverif/cocos_same_endpoint_current.pv

echo
echo "== cocos_exporter_label_current.pv =="
proverif verification/post-handshake/proverif/cocos_exporter_label_current.pv

echo
echo "== cocos_leaf_key_substitution_current.pv =="
proverif verification/post-handshake/proverif/cocos_leaf_key_substitution_current.pv

echo
echo "== cocos_infrastructure_identity_current.pv =="
proverif verification/post-handshake/proverif/cocos_infrastructure_identity_current.pv

echo
echo "== cocos_infrastructure_identity_bound_current.pv =="
proverif verification/post-handshake/proverif/cocos_infrastructure_identity_bound_current.pv

echo
echo "== cocos_agent_identity_current.pv =="
proverif verification/post-handshake/proverif/cocos_agent_identity_current.pv

echo
echo "== cocos_agent_identity_bound_current.pv =="
proverif verification/post-handshake/proverif/cocos_agent_identity_bound_current.pv

echo
echo "== cocos_task_context_current.pv =="
proverif verification/post-handshake/proverif/cocos_task_context_current.pv

echo
echo "== cocos_task_context_bound_current.pv =="
proverif verification/post-handshake/proverif/cocos_task_context_bound_current.pv

echo
echo "== cocos_legacy_request_binding.pv =="
proverif verification/intra-handshake/proverif/cocos_legacy_request_binding.pv

echo
echo "== cocos_legacy_relay_attack.pv =="
proverif verification/intra-handshake/proverif/cocos_legacy_relay_attack.pv
