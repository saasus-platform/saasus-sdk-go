#!/bin/bash
# Test Param Generator 実行スクリプト
# ビルド不要で直接実行

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# go runで直接実行（ビルド不要）
cd "$SCRIPT_DIR"
exec go run . "$@"
