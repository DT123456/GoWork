package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/register": {
            "post": {
                "description": "用户注册",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["用户"],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "username": {"type": "string", "example": "testuser"},
                                "password": {"type": "string", "example": "password123"},
                                "nickname": {"type": "string", "example": "测试用户"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "成功"}
                }
            }
        },
        "/login": {
            "post": {
                "description": "用户登录",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["用户"],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "username": {"type": "string"},
                                "password": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "code": {"type": "integer"},
                                "msg": {"type": "string"},
                                "data": {
                                    "type": "object",
                                    "properties": {
                                        "token": {"type": "string"},
                                        "user_id": {"type": "integer"},
                                        "username": {"type": "string"},
                                        "roles": {"type": "array", "items": {"type": "string"}}
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/captcha": {
            "get": {
                "description": "获取图形验证码",
                "produces": ["application/json"],
                "tags": ["验证码"],
                "summary": "获取验证码",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "captcha_id": {"type": "string"},
                                "image": {"type": "string"}
                            }
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "description": "上传文件",
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "tags": ["文件"],
                "summary": "上传文件",
                "parameters": [
                    {
                        "type": "file",
                        "description": "文件",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {"description": "成功"}
                }
            }
        },
        "/admin/roles": {
            "get": {
                "security": [{"Bearer": []}],
                "tags": ["角色管理"],
                "summary": "获取角色列表",
                "responses": {
                    "200": {"description": "成功"}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "tags": ["角色管理"],
                "summary": "创建角色",
                "parameters": [
                    {
                        "description": "角色信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string"},
                                "description": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "成功"}
                }
            }
        },
        "/admin/logs/system": {
            "get": {
                "security": [{"Bearer": []}],
                "tags": ["系统日志"],
                "summary": "获取系统日志列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {"description": "成功"}
                }
            }
        },
        "/admin/logs/access": {
            "get": {
                "security": [{"Bearer": []}],
                "tags": ["访问日志"],
                "summary": "获取访问日志列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {"description": "成功"}
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "token",
            "in": "header",
            "description": "JWT认证Token"
        }
    }
}`

// SwaggerInfo Swagger配置
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Gin-Advanced API 文档",
	Description:      "Gin-Advanced RESTful API 文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
