package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"
)

// 演示各种代码规范问题
func main() {
	fmt.Println("=== ch19: Go 代码规范检查和优化 ===")

	demoVet()
	demoGofmt()
	demoASTLint()
	demoCommonIssues()
	demoGolintConfig()

	// ============================================================
	// 逃逸分析演示
	// ============================================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== 逃逸分析演示 ===")
	fmt.Println(strings.Repeat("=", 60))

	noEscape()
	escapeToHeap()
	escapeInterface()
	escapeToCollection()
	escapeClosure()
	demoEscapeCommand()
	commonPitfalls()
	unsafeDemo()
	summary()
}

// ============================================================
// go vet 演示
// ============================================================
func demoVet() {
	fmt.Println("\n--- 1. go vet 静态分析 ---")

	// go vet 会检测以下问题：
	// - fmt.Printf 参数不匹配
	// - 错误的 mutex 使用
	// - 无效的结构体标签
	// - 空的 critical section

	badCode := `
package test

import "fmt"

func demo() {
	name := "Alice"
	fmt.Printf("Name: %d\n", name) // ❌ %d 用于数字，但传入 string
}
`
	// 将代码写入临时文件并运行 go vet
	tmpFile := "tmp_vet.go"
	os.WriteFile(tmpFile, []byte(badCode), 0644)
	defer os.Remove(tmpFile)

	cmd := exec.Command("go", "vet", tmpFile)
	output, _ := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Printf("go vet 发现问题:\n%s\n", output)
	}

	// 常见 vet 检查项
	fmt.Println("常用 vet 检查:")
	fmt.Println("  go vet ./...              # 检查所有包")
	fmt.Println("  go vet -printf ./...       # printf 风格检查")
	fmt.Println("  go vet -shadow ./...       # 变量遮蔽检查")
}

// ============================================================
// gofmt 演示
// ============================================================
func demoGofmt() {
	fmt.Println("\n--- 2. gofmt 代码格式化 ---")

	// 未格式化的代码
	unformatted := `package main

import ("fmt"; "strings")

func greet(name string)  {
      msg := "Hello, " + strings.TrimSpace(name)
fmt.Println(msg)
}

func main(){greet("  World  ")}
`

	fmt.Println("原始代码（有格式问题）:")
	fmt.Println(unformatted)

	// 解析并格式化
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", unformatted, parser.ParseComments)
	if err == nil {
		formatted, err := format.Source([]byte(unformatted))
		if err == nil {
			fmt.Println("格式化后:")
			fmt.Println(string(formatted))
		}
	}

	fmt.Println("gofmt 命令:")
	fmt.Println("  gofmt -w .              # 格式化并写入文件")
	fmt.Println("  gofmt -d .              # 显示差异")
	fmt.Println("  gofmt -l .              # 列出需格式化的文件")
}

// ============================================================
// AST 代码检查演示
// ============================================================
func demoASTLint() {
	fmt.Println("\n--- 3. 自定义 AST 代码检查 ---")

	code := `
package main

import "fmt"

const MaxRetries = 3 // 常量命名 OK
const maxretries = 3 // ❌ 应使用大写下划线格式

var GlobalVar = "test" // 全局变量应使用驼峰或帕斯卡

func calculate() {
	unused := "this variable is never used" // ❌ 未使用的变量
	_ = unused
}

func GetUserById(id int) string { // ❌ exported 函数应注释
	return "user"
}

func helper() { // ✅ 私有函数
}
`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	issues := lintFile(fset, f)
	fmt.Println("检测到的问题:")
	for _, issue := range issues {
		fmt.Printf("  [%s:%d] %s\n", issue.file, issue.line, issue.msg)
	}
}

type lintIssue struct {
	file string
	line int
	msg  string
}

func lintFile(fset *token.FileSet, f *ast.File) []lintIssue {
	var issues []lintIssue

	// 检查常量命名
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vs.Names {
						if name.IsExported() {
							continue // 导出的常量可以全大写
						}
						// 小写常量应该使用驼峰
						if strings.Contains(name.Name, "_") && !(strings.ToUpper(name.Name) == name.Name) {
							issues = append(issues, lintIssue{
								file: f.Name.Name,
								line: fset.Position(name.NamePos).Line,
								msg:  "常量命名应使用 CamelCase 或 SCREAMING_SNAKE_CASE",
							})
						}
					}
				}
			}
		}
	}

	// 检查函数注释
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.IsExported() && !hasDoc(fn) {
				issues = append(issues, lintIssue{
					file: f.Name.Name,
					line: fset.Position(fn.Pos()).Line,
					msg:  fmt.Sprintf("导出的函数 %s 缺少文档注释", fn.Name.Name),
				})
			}
		}
	}

	return issues
}

func hasDoc(fn *ast.FuncDecl) bool {
	return fn.Doc != nil && len(fn.Doc.Text()) > 0
}

// ============================================================
// 常见规范问题汇总
// ============================================================
func demoCommonIssues() {
	fmt.Println("\n--- 4. 常见规范问题及修复 ---")

	// 问题1: 错误处理
	fmt.Println("\n1. 错误处理:")
	badErrorHandle()
	goodErrorHandle()

	// 问题2: 上下文传递
	fmt.Println("\n2. Context 传递:")
	demoContextWrong()
	demoContextRight()

	// 问题3: 资源清理
	fmt.Println("\n3. defer 资源清理:")
	demoResourceCleanup()
}

// ❌ 错误: 忽略错误
func badErrorHandle() {
	fmt.Println("  ❌ 错误示例:")
	fmt.Println(`    file, _ := os.Open("test.txt")`)
	fmt.Println(`    // 错误被忽略，无法判断是否成功`)
}

// ✅ 正确: 处理错误
func goodErrorHandle() {
	fmt.Println("  ✅ 正确示例:")
	fmt.Println(`    file, err := os.Open("test.txt")
    if err != nil {
        return fmt.Errorf("open file: %w", err)
    }
    defer file.Close()`)
}

// ❌ 错误: 不传递 context
func demoContextWrong() {
	fmt.Println("  ❌ 错误示例:")
	fmt.Println(`    func BadFunc() {
        result := db.Query("SELECT ...") // ❌ 应传递 ctx
    }`)
}

// ✅ 正确: 使用 context
func demoContextRight() {
	fmt.Println("  ✅ 正确示例:")
	fmt.Println(`    func GoodFunc(ctx context.Context) error {
        rows, err := db.QueryContext(ctx, "SELECT ...")
        if err != nil {
            return err
        }
        defer rows.Close()
        // ...
    }`)
}

// 资源清理示例
func demoResourceCleanup() {
	// 正确使用 defer
	defer func() {
		// ❌ 错误: defer 在循环中
		// for _, file := range files {
		//     f, _ := os.Open(file)
		//     defer f.Close() // defer 不会在每次迭代后执行
		// }

		// ✅ 正确: 提取到函数或使用 defer 在循环内立即执行
		for _, file := range []string{"a.txt", "b.txt"} {
			func() {
				f, err := os.Open(file)
				if err != nil {
					return
				}
				defer f.Close()
			}()
		}
	}()
}

// ============================================================
// golangci-lint 配置示例
// ============================================================
func demoGolintConfig() {
	fmt.Println("\n--- 5. golangci-lint 配置示例 ---")

	config := `
# .golangci.yml 示例配置
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck      # 检查未处理的错误
    - gosimple      # 简化代码
    - govet         # 静态分析
    - ineffassign   # 检测未使用的赋值
    - staticcheck   # 静态检查
    - unused        # 检测未使用的代码

linters-settings:
  errcheck:
    check-type-assertions: true
  govet:
    enable-all: true
`
	fmt.Println("配置文件示例:")
	fmt.Println(config)
}

// ============================================================
// Benchmark 对比格式化前后性能
// ============================================================
func BenchmarkFormat(n int) {
	bad := strings.Repeat("x:=1;", n)
	start := time.Now()
	_, _ = format.Source([]byte(bad))
	fmt.Printf("格式化 %d 次操作耗时: %v\n", n, time.Since(start))
}

// ============================================================
// Go 逃逸分析（Escape Analysis）
// ============================================================
/*
什么是逃逸分析？
  Go 编译器通过逃逸分析决定变量分配在栈(Stack)还是堆(Heap)上。

栈分配：
  - 特点：速度快，分配和释放只需移动栈指针
  - 生命周期：函数返回后自动释放
  - 缺点：栈空间有限（通常 MB 级别）

堆分配：
  - 特点：空间大，但需要 GC（垃圾回收）管理
  - 生命周期：由 GC 决定何时释放
  - 缺点：有 GC 开销，可能导致内存碎片

逃逸规则：
  1. 变量逃逸到堆的情况：
     - 变量被返回（函数外部需要访问）
     - 变量被闭包引用
     - 变量类型不确定（如 interface{}）
     - 变量地址被写入 channel、slice、map、指针

  2. 变量留在栈的情况：
     - 变量作用域在函数内
     - 变量不逃逸到函数外

如何查看逃逸分析结果？
  go build -gcflags="-m" main.go
*/

// 场景1: 不逃逸 - 变量在栈上分配
func noEscape() {
	fmt.Println("\n--- 1. 不逃逸（栈分配）---")

	a := 10
	b := 20
	c := a + b

	fmt.Printf("a=%d, b=%d, c=%d (全部在栈上)\n", a, b, c)
	fmt.Println("&c 的地址:", &c)
	fmt.Println("→ 变量 c 不逃逸，分配在栈上")
}

// 场景2: 逃逸 - 返回指针
func escapeToHeap() {
	fmt.Println("--- 2. 逃逸到堆（返回指针）---")

	p := smallValue()
	fmt.Printf("返回的指针值: %d\n", *p)
	fmt.Println("→ smallValue 逃逸，分配在堆上")

	large := largeValue()
	fmt.Printf("largeValue: %s\n", large)
}

func smallValue() *int {
	x := 100
	fmt.Println("  smallValue: x 逃逸到堆")
	return &x
}

func largeValue() string {
	s := "Hello, 逃逸分析!"
	fmt.Println("  largeValue: s 逃逸到堆")
	return s
}

// 场景3: 逃逸 - interface 类型
func escapeInterface() {
	fmt.Println("\n--- 3. 逃逸（interface 类型）---")

	var i interface{} = 42
	fmt.Printf("interface{} 存储的值: %v\n", i)
	fmt.Println("→ 42 逃逸，因为赋值给 interface{}")

	n := 100
	fmt.Printf("fmt.Println: %d\n", n)
	fmt.Println("→ n 可能有轻微开销（视编译器版本而定）")
}

// 场景4: 逃逸 - 写入 map/slice/channel
func escapeToCollection() {
	fmt.Println("\n--- 4. 逃逸（写入集合）---")

	slice := make([]int, 0, 10)
	for i := 0; i < 5; i++ {
		slice = append(slice, i)
	}
	fmt.Printf("slice: %v\n", slice)
	fmt.Println("→ slice 可能逃逸（取决于编译器优化）")

	userMap := make(map[string]*int)
	val := 42
	userMap["age"] = &val
	fmt.Printf("map[\"age\"]: %d\n", *userMap["age"])
	fmt.Println("→ val 逃逸，存储在 map 中的指针必须逃逸")

	ch := make(chan *int, 1)
	num := 99
	ch <- &num
	fmt.Printf("channel 接收: %d\n", *<-ch)
	fmt.Println("→ num 逃逸，channel 持有指针必须逃逸")
}

// 场景5: 逃逸 - 闭包
func escapeClosure() {
	fmt.Println("\n--- 5. 逃逸（闭包）---")

	closure := outer()
	result := closure()
	fmt.Printf("闭包返回的值: %d\n", result)
	fmt.Println("→ x 逃逸，闭包捕获的变量必须分配在堆上")
}

func outer() func() int {
	x := 10
	fmt.Println("  outer: x 逃逸到堆")
	return func() int {
		return x + 1
	}
}

// 演示如何查看逃逸分析结果
func demoEscapeCommand() {
	fmt.Println("\n--- 逃逸分析命令 ---")
	fmt.Println("查看逃逸分析结果：")
	fmt.Println(`  go build -gcflags="-m" ./ch19`)
	fmt.Println(`  go build -gcflags="-m -m" ./ch19  # 更详细信息`)
	fmt.Println("")
	fmt.Println("输出示例说明：")
	fmt.Println("  ./main.go:10:6: moved to heap: x  # 变量 x 逃逸到堆")
	fmt.Println("  ./main.go:15:10: &y escapes to heap  # &y 逃逸")
	fmt.Println("  ./main.go:20:7: argument c not in memory  # c 不逃逸")
	fmt.Println("")
	fmt.Println("常见逃逸标记：")
	fmt.Println("  \"moved to heap\"   - 变量从栈移动到堆")
	fmt.Println("  \"escapes to heap\" - 变量逃逸到堆")
	fmt.Println("  \"not in memory\"   - 变量在寄存器/栈")
}

// 常见问题与优化
func commonPitfalls() {
	fmt.Println("\n--- 常见问题与优化 ---")

	badLoop := func() {
		var result []int
		for i := 0; i < 1000; i++ {
			tmp := make([]int, 1000)
			tmp[0] = i
			result = append(result, tmp[0])
		}
		_ = result
	}
	_ = badLoop

	goodLoop := func() {
		result := make([]int, 0, 1000)
		for i := 0; i < 1000; i++ {
			tmp := make([]int, 1000)
			tmp[0] = i
			result = append(result, tmp[0])
		}
		_ = result
	}
	_ = goodLoop

	fmt.Println("优化技巧:")
	fmt.Println("  1. 预分配 slice 容量，避免多次扩容")
	fmt.Println("  2. 小结构体可直接返回（栈分配）")
	fmt.Println("  3. 避免不必要的 interface{} 装箱")
	fmt.Println("  4. 使用 sync.Pool 复用大对象")
}

// 使用 unsafe.Pointer 分析内存
func unsafeDemo() {
	fmt.Println("\n--- unsafe.Pointer 内存分析 ---")

	str := "hello"
	strPtr := (*unsafe.Pointer)(unsafe.Pointer(&str))

	fmt.Printf("字符串: %s, 大小: %d 字节\n", str, len(str))
	fmt.Printf("字符串指针: %p\n", *strPtr)
	fmt.Println("→ 字符串内容在堆上（常量字符串可能有特殊优化）")
}

// 逃逸分析总结
func summary() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("逃逸分析总结")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(`
| 场景                 | 是否逃逸 | 原因                     |
|---------------------|---------|--------------------------|
| 返回指针              | ✅ 逃逸   | 外部需要访问               |
| 闭包捕获变量          | ✅ 逃逸   | 闭包生命周期超出函数        |
| interface{} 赋值      | ✅ 逃逸   | 编译器无法确定生命周期     |
| map/slice/channel   | ⚠️ 可能逃逸 | 集合持有引用需要堆分配     |
| 局部变量（不逃逸）     | ❌ 不逃逸  | 栈上分配，函数返回即释放    |

最佳实践：
  ✅ 优先返回结构体而非指针（小型结构体可栈分配）
  ✅ 预分配集合容量，减少 GC 压力
  ✅ 避免不必要的 interface{} 装箱
  ✅ 使用 go build -gcflags="-m" 检查逃逸情况
`)
}
