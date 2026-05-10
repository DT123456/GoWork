package service

import (
	"encoding/json"
	"gin-advanced/core"
	"gin-advanced/log"
	"gin-advanced/models"
	"time"

	"go.uber.org/zap"
)

// ============ Service 层统一日志 ============

// LogInfo 信息日志
func LogInfo(module, action string, fields ...interface{}) {
	msg := log.Sprintf("[%s] %s", module, action)
	log.MySQL(msg,
		buildFields(append([]interface{}{"action", action}, fields...))...,
	)
}

// buildFields 构建 zap 字段
func buildFields(fields []interface{}) []zap.Field {
	var zapFields []zap.Field
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			zapFields = append(zapFields, zap.Any(key, fields[i+1]))
		}
	}
	return zapFields
}

// LogError 错误日志
func LogError(module, action, errMsg string, fields ...interface{}) {
	log.E(log.Sprintf("[%s] %s failed: %s", module, action, errMsg))
}

// LogDebug 调试日志
func LogDebug(module, action string, fields ...interface{}) {
	log.D(log.Sprintf("[%s] %s", module, action))
}

// ============ 系统操作日志 ============

// CreateSystemLog 创建系统操作日志
func CreateSystemLog(systemLog *models.SystemLog) error {
	err := core.DB.Create(systemLog).Error
	if err != nil {
		LogError("system_log", "create", err.Error())
	} else {
		LogInfo("system_log", "create",
			"user_id", systemLog.UserID,
			"action", systemLog.Action,
		)
	}
	return err
}

// LogOperation 记录操作日志（快捷方法）
func LogOperation(userID uint, username, action, module, message string) error {
	systemLog := &models.SystemLog{
		UserID:   userID,
		Username: username,
		Action:   action,
		Module:   module,
		Status:   1,
		Message:  message,
	}
	return CreateSystemLog(systemLog)
}

// LogOperationError 记录操作失败日志
func LogOperationError(userID uint, username, action, module, message, errorMsg string) error {
	systemLog := &models.SystemLog{
		UserID:   userID,
		Username: username,
		Action:   action,
		Module:   module,
		Status:   0,
		Message:  message,
		Code:     "ERROR",
		ErrorMsg: errorMsg,
	}
	return CreateSystemLog(systemLog)
}

// GetSystemLogs 分页查询系统日志
func GetSystemLogs(page, pageSize int, filters map[string]interface{}) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := core.DB.Model(&models.SystemLog{})

	// 应用筛选条件
	if userID, ok := filters["user_id"]; ok && userID.(uint) > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if username, ok := filters["username"]; ok && username.(string) != "" {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if action, ok := filters["action"]; ok && action.(string) != "" {
		query = query.Where("action = ?", action)
	}
	if module, ok := filters["module"]; ok && module.(string) != "" {
		query = query.Where("module = ?", module)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

// GetSystemLogByID 根据ID获取日志详情
func GetSystemLogByID(id uint) (*models.SystemLog, error) {
	var log models.SystemLog
	err := core.DB.First(&log, id).Error
	return &log, err
}

// DeleteSystemLog 删除系统日志
func DeleteSystemLog(id uint) error {
	return core.DB.Delete(&models.SystemLog{}, id).Error
}

// DeleteSystemLogsBatch 批量删除系统日志
func DeleteSystemLogsBatch(ids []uint) error {
	return core.DB.Delete(&models.SystemLog{}, ids).Error
}

// ============ 接口访问日志 ============

// CreateAccessLog 创建访问日志
func CreateAccessLog(accessLog *models.AccessLog) error {
	err := core.DB.Create(accessLog).Error
	if err != nil {
		LogError("access_log", "create", err.Error())
	}
	return err
}

// GetAccessLogs 分页查询访问日志
func GetAccessLogs(page, pageSize int, filters map[string]interface{}) ([]models.AccessLog, int64, error) {
	var logs []models.AccessLog
	var total int64

	query := core.DB.Model(&models.AccessLog{})

	// 应用筛选条件
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if method, ok := filters["method"]; ok && method.(string) != "" {
		query = query.Where("method = ?", method)
	}
	if path, ok := filters["path"]; ok && path.(string) != "" {
		query = query.Where("path LIKE ?", "%"+path.(string)+"%")
	}
	if ip, ok := filters["ip"]; ok && ip.(string) != "" {
		query = query.Where("ip = ?", ip)
	}
	if statusCode, ok := filters["status_code"]; ok {
		query = query.Where("status_code = ?", statusCode)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

// GetAccessLogStatistics 获取访问统计
func GetAccessLogStatistics(startTime, endTime int64) (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// 总请求数
	var total int64
	core.DB.Model(&models.AccessLog{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).Count(&total)
	stats["total"] = total

	// 成功/失败数
	var successCount, failCount int64
	core.DB.Model(&models.AccessLog{}).Where("status_code >= 200 AND status_code < 400 AND created_at >= ? AND created_at <= ?", startTime, endTime).Count(&successCount)
	core.DB.Model(&models.AccessLog{}).Where("status_code >= 400 AND created_at >= ? AND created_at <= ?", startTime, endTime).Count(&failCount)
	stats["success"] = successCount
	stats["fail"] = failCount

	// 平均响应时间
	var avgDuration float64
	core.DB.Model(&models.AccessLog{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).Select("AVG(duration)").Scan(&avgDuration)
	stats["avg_duration"] = avgDuration

	// TOP 10 慢接口
	type SlowAPI struct {
		Path     string `json:"path"`
		AvgDuration float64 `json:"avg_duration"`
		Count    int64  `json:"count"`
	}
	var slowAPIs []SlowAPI
	core.DB.Model(&models.AccessLog{}).
		Select("path, AVG(duration) as avg_duration, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("path").
		Order("avg_duration DESC").
		Limit(10).
		Scan(&slowAPIs)
	stats["slow_apis"] = slowAPIs

	return stats, nil
}

// DeleteAccessLogsByDays 清理旧日志（按天数）
func DeleteAccessLogsByDays(days int) error {
	beforeTime := time.Now().AddDate(0, 0, -days).Unix()
	return core.DB.Where("created_at < ?", beforeTime).Delete(&models.AccessLog{}).Error
}

// ============ 日志统计 ============

// GetSystemLogStatistics 获取系统日志统计
func GetSystemLogStatistics(startTime, endTime int64) (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// 总操作数
	var total int64
	core.DB.Model(&models.SystemLog{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).Count(&total)
	stats["total"] = total

	// 成功/失败数
	var successCount, failCount int64
	core.DB.Model(&models.SystemLog{}).Where("status = 1 AND created_at >= ? AND created_at <= ?", startTime, endTime).Count(&successCount)
	core.DB.Model(&models.SystemLog{}).Where("status = 0 AND created_at >= ? AND created_at <= ?", startTime, endTime).Count(&failCount)
	stats["success"] = successCount
	stats["fail"] = failCount

	// TOP 操作类型
	type ActionStat struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	var topActions []ActionStat
	core.DB.Model(&models.SystemLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&topActions)
	stats["top_actions"] = topActions

	return stats, nil
}

// ObjectToJSON 对象转JSON字符串
func ObjectToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}
