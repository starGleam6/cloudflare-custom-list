#!/bin/bash

# 确保脚本出错时立即退出
set -e

# 项目名
PROJECT_NAME="cloudflare-custom-list"

# 不同平台和架构的列表
platforms=("windows/amd64" "windows/arm" "linux/amd64" "linux/arm" "darwin/amd64" "darwin/arm64")

for platform in "${platforms[@]}"
do 
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$PROJECT_NAME'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name

    echo 'Built for '$platform
done