package service

import (
	"gin-advanced/core"
	"gin-advanced/models"
	"testing"
)

func init() {
	// 初始化测试环境
	core.InitViper()
	core.InitMySQL()
}

// ============ User Service Tests ============

func TestCreateUser(t *testing.T) {
	user := &models.User{
		Username: "testuser123",
		Password: "password123",
		Nickname: "测试用户",
	}

	err := CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser failed: %v", err)
	}

	// 清理
	DeleteUser(int(user.ID))
}

func TestGetUserByUsername(t *testing.T) {
	// 先创建用户
	user := &models.User{
		Username: "test_get_user",
		Password: "password123",
		Nickname: "测试",
	}
	CreateUser(user)

	// 测试查询
	result, err := GetUserByUsername("test_get_user")
	if err != nil {
		t.Errorf("GetUserByUsername failed: %v", err)
	}
	if result.Username != "test_get_user" {
		t.Errorf("Expected username 'test_get_user', got '%s'", result.Username)
	}

	// 清理
	DeleteUser(int(user.ID))
}

func TestGetUserByID(t *testing.T) {
	// 先创建用户
	user := &models.User{
		Username: "test_get_id",
		Password: "password123",
		Nickname: "测试ID",
	}
	CreateUser(user)

	// 测试查询
	result, err := GetUserByID(int(user.ID))
	if err != nil {
		t.Errorf("GetUserByID failed: %v", err)
	}
	if result.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, result.ID)
	}

	// 清理
	DeleteUser(int(user.ID))
}

// ============ Role Service Tests ============

func TestCreateRole(t *testing.T) {
	role := &models.Role{
		Name:        "test_role",
		Description: "测试角色",
	}

	err := CreateRole(role)
	if err != nil {
		t.Errorf("CreateRole failed: %v", err)
	}

	// 验证
	result, _ := GetRoleByID(role.ID)
	if result.Name != "test_role" {
		t.Errorf("Expected role name 'test_role', got '%s'", result.Name)
	}
}

func TestGetRoles(t *testing.T) {
	roles, err := GetRoles()
	if err != nil {
		t.Errorf("GetRoles failed: %v", err)
	}
	if len(roles) == 0 {
		t.Log("No roles found, this is OK for empty database")
	}
}

func TestAssignRoleToUser(t *testing.T) {
	// 创建用户和角色
	user := &models.User{
		Username: "test_assign_user",
		Password: "password123",
		Nickname: "测试分配",
	}
	CreateUser(user)

	role := &models.Role{
		Name:        "test_assign_role",
		Description: "测试分配角色",
	}
	CreateRole(role)

	// 分配角色
	err := AssignRoleToUser(user.ID, role.ID)
	if err != nil {
		t.Errorf("AssignRoleToUser failed: %v", err)
	}

	// 验证
	userRoles, _ := GetUserRoles(user.ID)
	if len(userRoles) == 0 {
		t.Errorf("Expected user to have roles, got none")
	}
}

// ============ Log Service Tests ============

func TestCreateSystemLog(t *testing.T) {
	log := &models.SystemLog{
		UserID:   1,
		Username: "test_user",
		Action:   "test_action",
		Module:   "test_module",
		IP:       "127.0.0.1",
		Status:   1,
		Message:  "测试日志",
	}

	err := CreateSystemLog(log)
	if err != nil {
		t.Errorf("CreateSystemLog failed: %v", err)
	}
}

func TestGetSystemLogs(t *testing.T) {
	filters := map[string]interface{}{}
	_, total, err := GetSystemLogs(1, 10, filters)
	if err != nil {
		t.Errorf("GetSystemLogs failed: %v", err)
	}
	t.Logf("Found %d system logs", total)
}

func TestLogOperation(t *testing.T) {
	err := LogOperation(1, "test_user", "TEST_ACTION", "test_module", "测试操作日志")
	if err != nil {
		t.Errorf("LogOperation failed: %v", err)
	}
}

// ============ Benchmark Tests ============

func BenchmarkGetUserByUsername(b *testing.B) {
	// 创建测试用户
	user := &models.User{
		Username: "bench_user",
		Password: "password123",
		Nickname: "Benchmark",
	}
	CreateUser(user)
	defer DeleteUser(int(user.ID))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserByUsername("bench_user")
	}
}

func BenchmarkGetRoles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRoles()
	}
}
