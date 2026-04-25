package server

import "errors"

// MathService 提供数学计算服务
type MathService struct{}

// Args RPC 调用参数
type Args struct {
	A, B int
}

// Add 计算两数之和
func (m *MathService) Add(args Args, reply *int) error {
	if args.A < 0 || args.B < 0 {
		return errors.New("参数必须为非负数")
	}
	*reply = args.A + args.B
	return nil
}

// Multiply 计算两数之积
func (m *MathService) Multiply(args Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
