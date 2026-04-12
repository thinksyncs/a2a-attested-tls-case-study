#!/bin/zsh
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd "$script_dir"

if ! command -v pdflatex >/dev/null 2>&1; then
  echo "pdflatex not found in PATH" >&2
  exit 1
fi

rm -f \
  cocos_atls_verification_report.aux \
  cocos_atls_verification_report.out \
  cocos_atls_verification_report.toc \
  cocos_atls_verification_report.lof \
  cocos_atls_verification_report.lot \
  cocos_atls_verification_report.log
pdflatex -interaction=nonstopmode -halt-on-error cocos_atls_verification_report.tex
pdflatex -interaction=nonstopmode -halt-on-error cocos_atls_verification_report.tex
pdflatex -interaction=nonstopmode -halt-on-error cocos_atls_verification_report.tex
