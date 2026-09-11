#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GENERATOR_SCRIPT="${SCRIPT_DIR}/run.sh"
REPORTS_DIR="${SCRIPT_DIR}/../postman-analyzer/reports"
OUTPUT_FILE="${SCRIPT_DIR}/performance_metrics_$(date +%Y%m%d_%H%M%S).txt"

echo "==================================================================="
echo "パフォーマンス測定"
echo "実行日時: $(date '+%Y-%m-%d %H:%M:%S')"
echo "==================================================================="
echo ""

# 測定対象モジュール
MODULES=("auth" "billing" "pricing" "communication" "awsmarketplace" "integration")

# 結果を保存
cat > "${OUTPUT_FILE}" << 'EOF'
=================================================================
パフォーマンス測定レポート
実行日時: $(date '+%Y-%m-%d %H:%M:%S')
=================================================================

EOF

total_time=0
module_count=0

for module in "${MODULES[@]}"; do
  report_path="${REPORTS_DIR}/${module}/postman_analysis.json"
  
  if [[ ! -f "${report_path}" ]]; then
    echo "[WARN] Report not found for module '${module}', skipping" >&2
    continue
  fi
  
  echo "[INFO] Measuring performance for module '${module}'"
  
  # 処理時間を測定（3回実行して平均を取る）
  times=()
  for i in {1..3}; do
    start_time=$(date +%s.%N)
    "${GENERATOR_SCRIPT}" -report "../postman-analyzer/reports/${module}/postman_analysis.json" -module "${module}" -overwrite > /dev/null 2>&1
    end_time=$(date +%s.%N)
    elapsed=$(echo "${end_time} - ${start_time}" | bc -l)
    times+=("${elapsed}")
    echo "  実行 ${i}/3: ${elapsed}秒"
  done
  
  # 平均時間を計算
  avg_time=$(echo "scale=3; (${times[0]} + ${times[1]} + ${times[2]}) / 3" | bc -l)
  
  echo "モジュール: ${module}" >> "${OUTPUT_FILE}"
  echo "  平均処理時間: ${avg_time}秒" >> "${OUTPUT_FILE}"
  echo "  実行1: ${times[0]}秒" >> "${OUTPUT_FILE}"
  echo "  実行2: ${times[1]}秒" >> "${OUTPUT_FILE}"
  echo "  実行3: ${times[2]}秒" >> "${OUTPUT_FILE}"
  
  # 目標値との比較（5秒以内）
  if (( $(echo "${avg_time} <= 5.0" | bc -l) )); then
    echo "  ステータス: ✓ 目標達成 (5秒以内)" >> "${OUTPUT_FILE}"
  else
    echo "  ステータス: ✗ 目標未達 (5秒超過)" >> "${OUTPUT_FILE}"
  fi
  echo "" >> "${OUTPUT_FILE}"
  
  total_time=$(echo "${total_time} + ${avg_time}" | bc -l)
  module_count=$((module_count + 1))
  
  echo "[SUCCESS] Module '${module}': 平均 ${avg_time}秒"
  echo ""
done

# サマリーを追加
if [[ ${module_count} -gt 0 ]]; then
  avg_total=$(echo "scale=3; ${total_time} / ${module_count}" | bc -l)
  
  cat >> "${OUTPUT_FILE}" << EOF
=================================================================
サマリー
=================================================================
測定モジュール数: ${module_count}
総処理時間: ${total_time}秒
平均処理時間: ${avg_total}秒

目標値: 5秒以内
EOF

  if (( $(echo "${avg_total} <= 5.0" | bc -l) )); then
    echo "全体ステータス: ✓ 目標達成" >> "${OUTPUT_FILE}"
  else
    echo "全体ステータス: ✗ 目標未達" >> "${OUTPUT_FILE}"
  fi
  
  echo "" >> "${OUTPUT_FILE}"
  echo "=================================================================" >> "${OUTPUT_FILE}"
fi

echo ""
echo "================================================================="
echo "測定完了！結果: ${OUTPUT_FILE}"
echo "================================================================="
cat "${OUTPUT_FILE}"
