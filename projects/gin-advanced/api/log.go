package api

import (
	"gin-advanced/log"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemLogsResponse 系统日志列表响应
type SystemLogsResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// GetSystemLogs
// @Summary 获取系统日志列表
// @Description 分页获取系统日志，支持按用户、模块、操作等条件筛选
// @Tags 系统日志
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param user_id query int false "用户ID"
// @Param username query string false "用户名"
// @Param action query string false "操作类型"
// @Param module query string false "模块"
// @Param status query int false "状态码"
// @Param start_time query int false "开始时间戳"
// @Param end_time query int false "结束时间戳"
// @Success 200 {object} utils.Response{data=SystemLogsResponse}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/system [get]
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

	utils.Success(c, SystemLogsResponse{
		List:  logs,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}

// GetSystemLogDetail
// @Summary 获取系统日志详情
// @Description 根据ID获取单条系统日志的详细信息
// @Tags 系统日志
// @Produce json
// @Security BearerAuth
// @Param id path int true "日志ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /admin/logs/system/{id} [get]
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

// GetSystemLogStatistics
// @Summary 获取系统日志统计
// @Description 统计指定时间范围内的系统日志数量
// @Tags 系统日志
// @Produce json
// @Security BearerAuth
// @Param start_time query int false "开始时间戳"
// @Param end_time query int false "结束时间戳"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/system/stats [get]
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

// DeleteSystemLogsRequest 删除系统日志请求
type DeleteSystemLogsRequest struct {
	IDs []uint `json:"ids"`
}

// DeleteSystemLogs
// @Summary 删除系统日志
// @Description 批量删除指定的系统日志
// @Tags 系统日志
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body DeleteSystemLogsRequest true "日志ID列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/system [delete]
func DeleteSystemLogs(c *gin.Context) {
	var req DeleteSystemLogsRequest

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

// AccessLogsResponse 访问日志列表响应
type AccessLogsResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// GetAccessLogs
// @Summary 获取访问日志列表
// @Description 分页获取访问日志，支持按用户、请求方法、路径等条件筛选
// @Tags 访问日志
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param user_id query int false "用户ID"
// @Param method query string false "请求方法"
// @Param path query string false "请求路径"
// @Param ip query string false "IP地址"
// @Param status_code query int false "状态码"
// @Param start_time query int false "开始时间戳"
// @Param end_time query int false "结束时间戳"
// @Success 200 {object} utils.Response{data=AccessLogsResponse}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/access [get]
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

	utils.Success(c, AccessLogsResponse{
		List:  logs,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}

// GetAccessLogStatistics
// @Summary 获取访问统计
// @Description 统计指定时间范围内的访问日志数量和趋势
// @Tags 访问日志
// @Produce json
// @Security BearerAuth
// @Param start_time query int false "开始时间戳"
// @Param end_time query int false "结束时间戳"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/access/stats [get]
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

// CleanAccessLogsRequest 清理访问日志请求
type CleanAccessLogsRequest struct {
	Days int `json:"days"`
}

// CleanAccessLogs
// @Summary 清理访问日志
// @Description 删除指定天数之前的访问日志
// @Tags 访问日志
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CleanAccessLogsRequest false "保留天数" default(30)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/logs/access [delete]
func CleanAccessLogs(c *gin.Context) {
	var req CleanAccessLogsRequest

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
