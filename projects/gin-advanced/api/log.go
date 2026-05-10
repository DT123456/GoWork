package api

import (
	"gin-advanced/log"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ============ 系统日志接口 ============

// GetSystemLogs 获取系统日志列表
func GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := map[string]interface{}{}

	// 筛选条件
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"], _ = strconv.ParseUint(userID, 10, 32)
	}
	if username := c.Query("username"); username != "" {
		filters["username"] = username
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if module := c.Query("module"); module != "" {
		filters["module"] = module
	}
	if status := c.Query("status"); status != "" {
		filters["status"], _ = strconv.Atoi(status)
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := strconv.ParseInt(startTime, 10, 64); err == nil {
			filters["start_time"] = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := strconv.ParseInt(endTime, 10, 64); err == nil {
			filters["end_time"] = t
		}
	}

	logs, total, err := service.GetSystemLogs(page, pageSize, filters)
	if err != nil {
		log.E(log.Sprintf("查询系统日志失败: %v", err))
		utils.Error(c, "查询失败")
		return
	}

	utils.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetSystemLogDetail 获取系统日志详情
func GetSystemLogDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	systemLog, err := service.GetSystemLogByID(uint(id))
	if err != nil {
		utils.NotFound(c, "日志不存在")
		return
	}

	utils.Success(c, systemLog)
}

// GetSystemLogStatistics 获取系统日志统计
func GetSystemLogStatistics(c *gin.Context) {
	// 默认查询最近7天
	endTime := time.Now().Unix()
	startTime := time.Now().AddDate(0, 0, -7).Unix()

	if start := c.Query("start_time"); start != "" {
		if t, err := strconv.ParseInt(start, 10, 64); err == nil {
			startTime = t
		}
	}
	if end := c.Query("end_time"); end != "" {
		if t, err := strconv.ParseInt(end, 10, 64); err == nil {
			endTime = t
		}
	}

	stats, err := service.GetSystemLogStatistics(startTime, endTime)
	if err != nil {
		log.E(log.Sprintf("查询系统日志统计失败: %v", err))
		utils.Error(c, "查询失败")
		return
	}

	utils.Success(c, stats)
}

// DeleteSystemLogs 删除系统日志
func DeleteSystemLogs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := service.DeleteSystemLogsBatch(req.IDs); err != nil {
		log.E(log.Sprintf("删除系统日志失败: %v", err))
		utils.Error(c, "删除失败")
		return
	}

	utils.Success(c, "删除成功")
}

// ============ 访问日志接口 ============

// GetAccessLogs 获取访问日志列表
func GetAccessLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := map[string]interface{}{}

	// 筛选条件
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"], _ = strconv.ParseUint(userID, 10, 32)
	}
	if method := c.Query("method"); method != "" {
		filters["method"] = method
	}
	if path := c.Query("path"); path != "" {
		filters["path"] = path
	}
	if ip := c.Query("ip"); ip != "" {
		filters["ip"] = ip
	}
	if statusCode := c.Query("status_code"); statusCode != "" {
		filters["status_code"], _ = strconv.Atoi(statusCode)
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := strconv.ParseInt(startTime, 10, 64); err == nil {
			filters["start_time"] = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := strconv.ParseInt(endTime, 10, 64); err == nil {
			filters["end_time"] = t
		}
	}

	logs, total, err := service.GetAccessLogs(page, pageSize, filters)
	if err != nil {
		log.E(log.Sprintf("查询访问日志失败: %v", err))
		utils.Error(c, "查询失败")
		return
	}

	utils.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetAccessLogStatistics 获取访问统计
func GetAccessLogStatistics(c *gin.Context) {
	// 默认查询最近7天
	endTime := time.Now().Unix()
	startTime := time.Now().AddDate(0, 0, -7).Unix()

	if start := c.Query("start_time"); start != "" {
		if t, err := strconv.ParseInt(start, 10, 64); err == nil {
			startTime = t
		}
	}
	if end := c.Query("end_time"); end != "" {
		if t, err := strconv.ParseInt(end, 10, 64); err == nil {
			endTime = t
		}
	}

	stats, err := service.GetAccessLogStatistics(startTime, endTime)
	if err != nil {
		log.E(log.Sprintf("查询访问统计失败: %v", err))
		utils.Error(c, "查询失败")
		return
	}

	utils.Success(c, stats)
}

// CleanAccessLogs 清理访问日志
func CleanAccessLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days"` // 保留多少天的日志
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Days = 30 // 默认保留30天
	}

	if err := service.DeleteAccessLogsByDays(req.Days); err != nil {
		log.E(log.Sprintf("清理访问日志失败: %v", err))
		utils.Error(c, "清理失败")
		return
	}

	utils.Success(c, "清理成功")
}
