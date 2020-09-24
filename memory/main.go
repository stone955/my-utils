package main

type smallStruct struct {
	a, b int64
	c, d float64
}

//go:noinline
func smallAlloc() *smallStruct {
	// go:noinline 会禁用内联 不能加空格

	// go tool compile "-m" main.go 运行逃逸分析命令
	// 	.\main.go:12:9: &smallStruct literal escapes to heap

	// go tool compile -S main.go 得到汇编代码，展示分配细节
	return &smallStruct{}
}
