#!/bin/bash
# 编译当前目录下所有 .proto 文件，生成 Go 代码到当前目录

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 检查 protoc 是否安装
if ! command -v protoc &>/dev/null; then
  echo "❌ protoc 未安装，请先安装："
  echo "   brew install protobuf"
  exit 1
fi

# 检查 Go 插件是否安装
if ! command -v protoc-gen-go &>/dev/null; then
  echo "❌ protoc-gen-go 未安装，请执行："
  echo "   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  exit 1
fi

# 用临时目录接收，再把生成的 .go 文件移回当前目录
TMP_DIR=$(mktemp -d)
CS_OUT="./cs_out"
mkdir -p "$CS_OUT"

echo "🔧 开始编译 proto 文件..."

for proto_file in *.proto; do
  echo "  📄 编译: $proto_file"

  # 生成 Go 代码到临时目录
  protoc \
    --go_out="$TMP_DIR" \
    "$proto_file"

  # 生成 C# 代码（可选，无插件时跳过）
  protoc \
    --csharp_out="$CS_OUT" \
    "$proto_file" 2>/dev/null || echo "  ⚠️  C# 插件未找到，跳过生成 C# 代码"
done

# 把生成的 .go 文件全部移回当前目录
find "$TMP_DIR" -name "*.go" -exec mv {} "$SCRIPT_DIR/" \;
rm -rf "$TMP_DIR"

echo "✅ 编译完成！"
echo "   Go 代码: $SCRIPT_DIR"
echo "   C# 代码: $CS_OUT"