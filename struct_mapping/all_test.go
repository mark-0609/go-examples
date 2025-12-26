package struct_mapping

import (
	"testing"
)

// 初始化测试数据
var srcUser = &User{
	ID:        1001,
	Name:      "John Doe",
	Email:     "john@example.com",
	Age:       30,
	Address:   "123 Main St, New York, NY",
	IsActive:  true,
	Score:     99.5,
	Tags:      []string{"admin", "editor", "viewer"},
	Metadata:  map[string]string{"role": "admin", "dept": "it"},
	CreatedAt: 1678888888,
}

func BenchmarkManualCopy(b *testing.B) {
	var dst UserDTO
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ManualCopy(srcUser, &dst)
	}
}

func BenchmarkGoverterCopy(b *testing.B) {
	var dst UserDTO
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GoverterGenCopy(srcUser, &dst)
	}
}

func BenchmarkCopierCopy(b *testing.B) {
	var dst UserDTO
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CopierCopy(srcUser, &dst)
	}
}

func BenchmarkSimpleInFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i
	}
}

func TestPlaceholder(t *testing.T) {
	// 占位符，确保 go test 运行
}
