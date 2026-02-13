#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="/workspace"
PROTO_DIR="$ROOT_DIR/proto"

if [[ ! -d "$PROTO_DIR" ]]; then
	echo "Proto directory not found: $PROTO_DIR" >&2
	exit 1
fi

mapfile -t PROTO_FILES < <(find "$PROTO_DIR" -type f -name "*.proto" | sort)

if [[ ${#PROTO_FILES[@]} -eq 0 ]]; then
	echo "No .proto files found under: $PROTO_DIR" >&2
	exit 1
fi

echo "Generating protos (${#PROTO_FILES[@]} files)"

rm -f "$ROOT_DIR/go-service/proto"/*.pb.go || true
rm -f "$ROOT_DIR/python-service/proto"/*_pb2*.py || true

protoc -I"$ROOT_DIR" \
	--go_out="$ROOT_DIR/go-service" \
	--go_opt=paths=source_relative \
	--go-grpc_out="$ROOT_DIR/go-service" \
	--go-grpc_opt=paths=source_relative \
	"${PROTO_FILES[@]/$ROOT_DIR\//}"

python3 -m grpc_tools.protoc -I"$ROOT_DIR" \
	--python_out="$ROOT_DIR/python-service" \
	--grpc_python_out="$ROOT_DIR/python-service" \
	"${PROTO_FILES[@]/$ROOT_DIR\//}"

if [[ ! -f "$ROOT_DIR/python-service/proto/__init__.py" ]]; then
	touch "$ROOT_DIR/python-service/proto/__init__.py"
fi

echo "Done"
