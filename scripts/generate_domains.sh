#!/bin/bash

# 対象となるフォルダ
TARGET_DIR="internal/domains/models"

# findコマンドを使用してフォルダ配下のファイルを取得し、ループで処理
find "$TARGET_DIR" -type f | while read filepath; do
    # 例: echoコマンドを実行
    echo $filepath
    gomodifytags -file $filepath -w -all -add-tags json
done
