#!/bin/bash

# 輸出檔案名稱
OUTPUT="output.exe"
PROJECTNAME="bk_click_robot"

# 編譯函式
build_windows() {
    echo "正在編譯 Windows 平台的可執行檔案..."
    GOOS=windows GOARCH=amd64 go build -o $OUTPUT
    if [ $? -eq 0 ]; then
        echo "編譯完成！生成檔案：$OUTPUT"
    else
        echo "編譯失敗！請檢查程式碼或環境配置。"
    fi
}

# 清理函式
clean() {
    echo "正在清理..."
    if [ -f $OUTPUT ]; then
        rm -f $OUTPUT
        echo "已刪除 $OUTPUT"
    else
        echo "無檔案可清理。"
    fi
}

# 主選單
case $1 in
    build)
        build_windows
        ;;
    clean)
        clean
        ;;
    *)
        echo "使用方法："
        echo "  sh build.sh build   # 編譯 Windows 可執行檔"
        echo "  sh build.sh clean   # 清除生成的檔案"
        ;;
esac
