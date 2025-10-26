#!/bin/bash

# Git自动上传脚本
# 使用方法: ./git_push.sh "commit message"
# 或者: ./git_push.sh

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 检查是否在git仓库中
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_message $RED "错误: 当前目录不是git仓库"
        exit 1
    fi
}

# 检查是否有未提交的更改
check_changes() {
    if git diff --quiet && git diff --cached --quiet; then
        print_message $YELLOW "没有检测到任何更改"
        return 1
    fi
    return 0
}

# 获取commit message
get_commit_message() {
    local message="$1"
    
    if [ -z "$message" ]; then
        # 如果没有提供commit message，尝试从git status生成
        local status_output=$(git status --porcelain)
        local modified_files=$(echo "$status_output" | grep -E '^[AM]' | wc -l)
        local new_files=$(echo "$status_output" | grep -E '^A' | wc -l)
        local deleted_files=$(echo "$status_output" | grep -E '^D' | wc -l)
        
        if [ $modified_files -gt 0 ] || [ $new_files -gt 0 ] || [ $deleted_files -gt 0 ]; then
            message="Update: "
            if [ $modified_files -gt 0 ]; then
                message+="modified $modified_files files"
            fi
            if [ $new_files -gt 0 ]; then
                if [ $modified_files -gt 0 ]; then message+=", "; fi
                message+="added $new_files files"
            fi
            if [ $deleted_files -gt 0 ]; then
                if [ $modified_files -gt 0 ] || [ $new_files -gt 0 ]; then message+=", "; fi
                message+="deleted $deleted_files files"
            fi
        else
            message="Auto commit $(date '+%Y-%m-%d %H:%M:%S')"
        fi
    fi
    
    echo "$message"
}

# 主函数
main() {
    print_message $BLUE "=== Git自动上传脚本 ==="
    
    # 检查git仓库
    check_git_repo
    
    # 检查是否有更改
    if ! check_changes; then
        print_message $YELLOW "没有检测到任何更改，退出"
        exit 0
    fi
    
    # 获取commit message
    local commit_message=$(get_commit_message "$1")
    print_message $GREEN "Commit message: $commit_message"
    
    # 显示当前状态
    print_message $BLUE "\n当前git状态:"
    git status --short
    
    # 添加所有更改
    print_message $BLUE "\n添加所有更改..."
    git add .
    
    # 提交更改
    print_message $BLUE "提交更改..."
    git commit -m "$commit_message"
    
    # 获取当前分支
    local current_branch=$(git branch --show-current)
    print_message $BLUE "当前分支: $current_branch"
    
    # 推送到远程仓库
    print_message $BLUE "推送到远程仓库..."
    if git push -u origin "$current_branch"; then
        print_message $GREEN "✅ 成功推送到远程仓库!"
    else
        print_message $RED "❌ 推送失败，请检查网络连接或权限"
        exit 1
    fi
    
    print_message $GREEN "\n🎉 Git上传完成!"
}

# 显示帮助信息
show_help() {
    echo "Git自动上传脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [commit_message]"
    echo ""
    echo "参数:"
    echo "  commit_message  可选的提交信息"
    echo ""
    echo "示例:"
    echo "  $0 \"添加新功能\""
    echo "  $0 \"修复bug\""
    echo "  $0              # 使用自动生成的提交信息"
    echo ""
    echo "功能:"
    echo "  - 自动检测并添加所有更改"
    echo "  - 支持自定义提交信息"
    echo "  - 自动推送到当前分支"
    echo "  - 彩色输出和错误处理"
}

# 检查参数
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# 执行主函数
main "$1"
