#!/bin/bash

# 현재 디렉토리부터 모든 go.mod 파일 찾기
find . -name "go.mod" | while read modfile; do
    dir=$(dirname "$modfile")  # go.mod가 있는 디렉토리 추출
    echo "Running 'go mod tidy' in $dir"
    
    # go.mod가 있는 디렉토리에서 실행
    (cd "$dir" && go mod tidy)
done

echo "✅ All go.mod directories have been tidied!"
