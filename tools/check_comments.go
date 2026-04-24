package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Go 文件注释检查工具
// 检查每个 Go 文件是否有完整的注释
func main() {
	dirs := []string{
		"server/internal/model",
		"server/internal/service",
		"server/internal/controller",
		"server/internal/repository",
	}

	for _, dir := range dirs {
		fmt.Printf("\n=== 检查 %s ===\n", dir)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			
			hasPackageComment, hasTypeComments, hasFuncComments, issues := checkFile(path)
			
			if !hasPackageComment || len(issues) > 0 {
				fmt.Printf("⚠️  %s\n", path)
				if !hasPackageComment {
					fmt.Printf("   ❌ 缺少包注释\n")
				}
				for _, issue := range issues {
					fmt.Printf("   ⚠️  %s\n", issue)
				}
			} else {
				fmt.Printf("✅ %s\n", path)
			}
			
			return nil
		})
	}
}

func checkFile(path string) (hasPackageComment bool, hasTypeComments bool, hasFuncComments bool, issues []string) {
	file, err := os.Open(path)
	if err != nil {
		issues = append(issues, fmt.Sprintf("无法打开文件：%v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	inType := false
	inFunc := false
	lastComment := ""
	
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// 检查包注释
		if lineNum <= 3 && strings.HasPrefix(line, "// Package") {
			hasPackageComment = true
		}
		
		// 检查类型定义
		if strings.HasPrefix(line, "type ") && strings.Contains(line, " struct") {
			inType = true
			if lastComment == "" {
				issues = append(issues, fmt.Sprintf("第%d行：类型定义缺少注释", lineNum))
			}
		}
		
		// 检查函数定义
		if strings.HasPrefix(line, "func ") {
			inFunc = true
			if lastComment == "" && !strings.Contains(line, "func (") {
				// 非方法函数必须有注释
				issues = append(issues, fmt.Sprintf("第%d行：函数缺少注释", lineNum))
			}
		}
		
		// 记录注释
		if strings.HasPrefix(line, "//") {
			lastComment = line
		} else if line != "" && !strings.HasPrefix(line, "package") && !strings.HasPrefix(line, "import") {
			if !inType && !inFunc {
				lastComment = ""
			}
		}
		
		// 重置状态
		if line == "}" {
			if inType {
				inType = false
			}
			if inFunc {
				inFunc = false
			}
		}
	}
	
	return hasPackageComment, hasTypeComments, hasFuncComments, issues
}
