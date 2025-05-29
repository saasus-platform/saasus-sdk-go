#!/bin/bash

# 使用法: ./split_client_gen.sh <入力ファイル.go>
# 例: ./split_client_gen.sh client.gen.go

input_file="$1"
output_dir="client_gen_split_output"

if [[ -z "$input_file" ]]; then
    echo "エラー: 入力ファイル名を指定してください。"
    echo "使用法: $0 <入力ファイル.go>"
    exit 1
fi

if [[ ! -f "$input_file" ]]; then
    echo "エラー: 入力ファイル '$input_file' が見つかりません。"
    exit 1
fi

mkdir -p "$output_dir"
rm -f "$output_dir"/*
echo "出力先ディレクトリ: $output_dir"

get_line_num() {
    local pattern="$1"
    grep -E -n -m 1 -- "$pattern" "$input_file" | cut -d: -f1
}

extract_chunk() {
    local start_line="$1"
    local end_line_exclusive="$2" # この行番号の直前まで
    local out_file="$3"
    local actual_end_line

    if [[ -z "$start_line" ]]; then
        echo "情報: セクション '$out_file' の開始点が見つかりません。空ファイルを作成します。"
        touch "$output_dir/$out_file"
        return
    fi

    # end_line_exclusive が空、または start_line が end_line_exclusive 以上の場合、EOFまで
    if [[ -z "$end_line_exclusive" || "$start_line" -ge "$end_line_exclusive" ]]; then
        actual_end_line="$L_EOF"
    else
        actual_end_line=$((end_line_exclusive - 1))
    fi

    # actual_end_line が start_line より小さい場合 (通常は上記の条件で回避されるはずだが念のため)
    if [[ "$actual_end_line" -lt "$start_line" ]]; then
         echo "情報: セクション '$out_file' の範囲が無効です (終了 $actual_end_line < 開始 $start_line)。空ファイルを作成します。"
         touch "$output_dir/$out_file"
         return
    fi

    echo "抽出中: $output_dir/$out_file (行: $start_line - $actual_end_line)"
    sed -n "${start_line},${actual_end_line}p" "$input_file" > "$output_dir/$out_file"
}

L_EOF=$(wc -l < "$input_file" | xargs)

# --- 各セクションの開始を定義するパターン (grep -E 向け) ---
P_CLIENT_STRUCT_START='^type Client struct'
P_CLIENT_INTERFACE_START='^type ClientInterface interface'
P_CLIENT_METHODS_START='^func \(c \*Client\) GetAuthInfo\('
P_REQUEST_BUILDERS_START='^func NewGetAuthInfoRequest\('
P_APPLY_EDITORS_START='^func \(c \*Client\) applyEditors\('
P_CLIENT_WITH_RESPONSES_STRUCT_START='^type ClientWithResponses struct'
P_CLIENT_WITH_RESPONSES_INTERFACE_START='^type ClientWithResponsesInterface interface'
P_CLIENT_WITH_RESPONSES_METHODS_START='^func \(c \*ClientWithResponses\) GetAuthInfoWithResponse\('
P_RESPONSE_STRUCTS_START='^type GetAuthInfoResponse struct'
P_PARSE_FUNCTIONS_START='^func ParseGetAuthInfoResponse\('

L_CLIENT_STRUCT_START=$(get_line_num "$P_CLIENT_STRUCT_START")
L_CLIENT_INTERFACE_START=$(get_line_num "$P_CLIENT_INTERFACE_START")
L_CLIENT_METHODS_START=$(get_line_num "$P_CLIENT_METHODS_START")
L_REQUEST_BUILDERS_START=$(get_line_num "$P_REQUEST_BUILDERS_START")
L_APPLY_EDITORS_START=$(get_line_num "$P_APPLY_EDITORS_START")
L_CLIENT_WITH_RESPONSES_STRUCT_START=$(get_line_num "$P_CLIENT_WITH_RESPONSES_STRUCT_START")
L_CLIENT_WITH_RESPONSES_INTERFACE_START=$(get_line_num "$P_CLIENT_WITH_RESPONSES_INTERFACE_START")
L_CLIENT_WITH_RESPONSES_METHODS_START=$(get_line_num "$P_CLIENT_WITH_RESPONSES_METHODS_START")
L_RESPONSE_STRUCTS_START=$(get_line_num "$P_RESPONSE_STRUCTS_START")
L_PARSE_FUNCTIONS_START=$(get_line_num "$P_PARSE_FUNCTIONS_START")

echo "--- 検出された行番号 ---"
echo "Client Struct Start: ${L_CLIENT_STRUCT_START:-未検出}"
echo "Client Interface Start: ${L_CLIENT_INTERFACE_START:-未検出}"
echo "Client Methods Start: ${L_CLIENT_METHODS_START:-未検出}"
echo "Request Builders Start: ${L_REQUEST_BUILDERS_START:-未検出}"
echo "Apply Editors Start: ${L_APPLY_EDITORS_START:-未検出}"
echo "ClientWithResponses Struct Start: ${L_CLIENT_WITH_RESPONSES_STRUCT_START:-未検出}"
echo "ClientWithResponses Interface Start: ${L_CLIENT_WITH_RESPONSES_INTERFACE_START:-未検出}"
echo "ClientWithResponses Methods Start: ${L_CLIENT_WITH_RESPONSES_METHODS_START:-未検出}"
echo "Response Structs Start: ${L_RESPONSE_STRUCTS_START:-未検出}"
echo "Parse Functions Start: ${L_PARSE_FUNCTIONS_START:-未検出}"
echo "ファイル終端 (EOF): $L_EOF"
echo "-------------------------"

# --- チャンク抽出 (検出された行番号の順序を考慮) ---
extract_chunk 1 "$L_CLIENT_STRUCT_START" "01_preamble.go"
extract_chunk "$L_CLIENT_STRUCT_START" "$L_CLIENT_INTERFACE_START" "02_client_definition.go"
extract_chunk "$L_CLIENT_INTERFACE_START" "$L_CLIENT_METHODS_START" "03_client_interface.go"
extract_chunk "$L_CLIENT_METHODS_START" "$L_REQUEST_BUILDERS_START" "04_client_methods.go"
extract_chunk "$L_REQUEST_BUILDERS_START" "$L_APPLY_EDITORS_START" "05_request_builders.go"
extract_chunk "$L_APPLY_EDITORS_START" "$L_CLIENT_WITH_RESPONSES_STRUCT_START" "06_apply_editors.go"
extract_chunk "$L_CLIENT_WITH_RESPONSES_STRUCT_START" "$L_CLIENT_WITH_RESPONSES_INTERFACE_START" "07_client_with_responses_definition.go"

# 実際のファイル構造に合わせてチャンクの区切りを調整
# ClientWithResponses Interface (L_CLIENT_WITH_RESPONSES_INTERFACE_START = 5237) の次は、
# Response Structs (L_RESPONSE_STRUCTS_START = 5576) が来る。
extract_chunk "$L_CLIENT_WITH_RESPONSES_INTERFACE_START" "$L_RESPONSE_STRUCTS_START" "08_client_with_responses_interface.go"

# 次に Response Structs (L_RESPONSE_STRUCTS_START = 5576)。
# これは ClientWithResponses Methods (L_CLIENT_WITH_RESPONSES_METHODS_START = 7525) の直前まで。
extract_chunk "$L_RESPONSE_STRUCTS_START" "$L_CLIENT_WITH_RESPONSES_METHODS_START" "09_response_structs.go"

# 次に ClientWithResponses Methods (L_CLIENT_WITH_RESPONSES_METHODS_START = 7525)。
# これは Parse Functions (L_PARSE_FUNCTIONS_START = 8618) の直前まで。
extract_chunk "$L_CLIENT_WITH_RESPONSES_METHODS_START" "$L_PARSE_FUNCTIONS_START" "10_client_with_responses_methods.go"

# 最後に Parse Functions (L_PARSE_FUNCTIONS_START = 8618) からファイル末尾まで。
extract_chunk "$L_PARSE_FUNCTIONS_START" "$((L_EOF + 1))" "11_parse_functions.go" # end_line_exclusive なので L_EOF+1 で L_EOF まで

echo "---"
echo "ファイルの分割処理が完了しました。'$output_dir' ディレクトリを確認してください。"
echo "注意: このスクリプトは提供されたファイルの構造とパターンに依存します。"
echo "他のファイルでは調整が必要になる場合があります。特に「未検出」と表示されたセクションのパターンを確認してください。"
echo "特定の関数名（例: GetAuthInfo）に依存しているパターンは、実際のファイルに合わせて変更する必要があるかもしれません。"
