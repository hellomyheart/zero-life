// Package response 定义API响应数据结构
package response

import "time"

// SystemInfoResp 系统信息响应
type SystemInfoResp struct {
	Version    string    `json:"version"`     // 系统版本
	GoVersion  string    `json:"go_version"`  // Go版本
	OS         string    `json:"os"`          // 操作系统
	Arch       string    `json:"arch"`        // 系统架构
	Uptime     string    `json:"uptime"`      // 运行时长
	Goroutines int       `json:"goroutines"`  // Goroutine数量
	StartTime  time.Time `json:"start_time"`  // 启动时间
}

// HealthCheckResp 健康检查响应
type HealthCheckResp struct {
	Status    string    `json:"status"`     // 状态（ok/error）
	Timestamp time.Time `json:"timestamp"`  // 时间戳
}
