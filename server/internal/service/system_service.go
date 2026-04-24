// Package service 业务逻辑层，实现核心业务逻辑
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"runtime"
	"time"
)

// SystemService 系统管理业务服务
// 负责处理系统管理相关的业务逻辑
type SystemService struct {
	startTime time.Time // 系统启动时间
}

// NewSystemService 创建系统管理服务实例
// 返回：
//   - *SystemService: 系统管理服务实例
func NewSystemService() *SystemService {
	return &SystemService{
		startTime: time.Now(),
	}
}

// GetSystemInfo 获取系统信息
// 返回：
//   - *response.SystemInfoResp: 系统信息
//   - error: 错误信息
func (s *SystemService) GetSystemInfo() (*response.SystemInfoResp, error) {
	return &response.SystemInfoResp{
		Version:     "1.0.0",
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		Uptime:      time.Since(s.startTime).String(),
		Goroutines:  runtime.NumGoroutine(),
		StartTime:   s.startTime,
	}, nil
}

// HealthCheck 健康检查
// 返回：
//   - *response.HealthCheckResp: 健康检查结果
//   - error: 错误信息
func (s *SystemService) HealthCheck() (*response.HealthCheckResp, error) {
	return &response.HealthCheckResp{
		Status:    "ok",
		Timestamp: time.Now(),
	}, nil
}
