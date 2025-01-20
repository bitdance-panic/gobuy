#!/bin/bash

root_path=${PWD}

# 1. 生成proto
echo "------ start generate proto ------"
idl_dir="$root_path/app/rpc/idl"
# 固定不变
MODULE="github.com/bitdance-panic/gobuy/app/rpc"

cd "$idl_dir/.."

# 遍历 IDL 文件夹中的所有 .proto 文件
for proto_file in "$idl_dir"/*.proto; do
    # 获取文件名（不带路径和后缀）
    name=$(basename "$proto_file" .proto)

    kitex -module="$MODULE" -I=./ ./idl/"$name".proto

    # 检查 kitex 是否执行成功
    if [ $? -eq 0 ]; then
        echo "Successfully generated code for $name.proto"
    else
        echo "Failed to generate code for $name.proto"
    fi
done
echo "------ end generate proto ------"


echo "------ start generate swagger ------"
# 2. 生成swagger
paths=(
    "./app/services/gateway"
)
# # 遍历每个路径并执行操作
for path in "$root_path/${paths[@]}"; do
    echo "Processing path: $path"
    # 检查路径是否存在
    if [ -d "$path" ]; then
        cd $path
        swag init
        echo "Swagger doc init success in $path"
    else
        echo "Path $path does not exist, skipping..."
    fi
done
echo "------ end generate swagger ------"