package models

// SystemLog 系统操作日志
type SystemLog struct {
	Model
	UserID    uint   `gorm:"index" json:"user_id"`                         // 用户ID
	Username  string `gorm:"size:64" json:"username"`                       // 用户名
	Action    string `gorm:"size:128;index" json:"action"`                   // 操作类型
	Module    string `gorm:"size:64;index" json:"module"`                   // 模块
	Method    string `gorm:"size:32" json:"method"`                         // 方法
	Path      string `gorm:"size:256" json:"path"`                         // 请求路径
	IP        string `gorm:"size:45" json:"ip"`                            // IP地址
	Location  string `gorm:"size:128" json:"location"`                     // 地理位置
	UserAgent string `gorm:"size:512" json:"user_agent"`                   // 浏览器信息
	Status    int    `gorm:"default:1;index" json:"status"`                 // 状态 1:成功 0:失败
	Code      string `gorm:"size:32" json:"code"`                          // 错误码
	ErrorMsg  string `gorm:"size:512" json:"error_msg"`                     // 错误信息
	Message   string `gorm:"size:512" json:"message"`                      // 日志消息
	Duration  int64  `gorm:"default:0" json:"duration"`                     // 耗时(毫秒)
	Request   string `gorm:"type:text" json:"request"`                     // 请求参数
	Response  string `gorm:"type:text" json:"response"`                   // 响应结果
	Before    string `gorm:"type:text" json:"before"`                      // 操作前数据
	After     string `gorm:"type:text" json:"after"`                       // 操作后数据
}

// AccessLog 接口访问日志
type AccessLog struct {
	Model
	UserID     uint   `gorm:"index" json:"user_id"`                         // 用户ID (0=未登录)
	Username   string `gorm:"size:64" json:"username"`                       // 用户名
	Method     string `gorm:"size:10;index" json:"method"`                   // HTTP方法
	Path       string `gorm:"size:256;index" json:"path"`                    // 请求路径
	Query      string `gorm:"type:text" json:"query"`                       // 查询参数
	Body       string `gorm:"type:text" json:"body"`                        // 请求体
	IP         string `gorm:"size:45;index" json:"ip"`                      // IP地址
	UserAgent  string `gorm:"size:512" json:"user_agent"`                   // 浏览器信息
	StatusCode int    `gorm:"index" json:"status_code"`                    // 响应状态码
	Duration   int64  `gorm:"index" json:"duration"`                       // 耗时(毫秒)
	ErrorMsg   string `gorm:"size:512" json:"error_msg"`                   // 错误信息
}
