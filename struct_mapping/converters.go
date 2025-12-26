package struct_mapping

import (
	"github.com/jinzhu/copier"
)

// 1. ManualCopy 手动赋值
func ManualCopy(src *User, dst *UserDTO) {
	if src == nil {
		return
	}
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Email = src.Email
	dst.Age = src.Age
	dst.Address = src.Address
	dst.IsActive = src.IsActive
	dst.Score = src.Score
	dst.CreatedAt = src.CreatedAt

	if src.Tags != nil {
		dst.Tags = make([]string, len(src.Tags))
		copy(dst.Tags, src.Tags)
	}

	if src.Metadata != nil {
		dst.Metadata = make(map[string]string, len(src.Metadata))
		for k, v := range src.Metadata {
			dst.Metadata[k] = v
		}
	}
}

// 2. CopierCopy 使用 github.com/jinzhu/copier
func CopierCopy(src *User, dst *UserDTO) error {
	return copier.Copy(dst, src)
}

// 3. GoverterGenCopy 模拟 goverter 生成的代码
// goverter 会生成类似下面的代码，通常它会处理指针检查和切片/Map的深拷贝
// 这里的实现与 ManualCopy 类似，因为 goverter 本质上就是生成这样的代码
func GoverterGenCopy(src *User, dst *UserDTO) {
	if src == nil {
		return
	}
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Email = src.Email
	dst.Age = src.Age
	dst.Address = src.Address
	dst.IsActive = src.IsActive
	dst.Score = src.Score
	dst.CreatedAt = src.CreatedAt

	// goverter 生成的切片拷贝代码
	if src.Tags != nil {
		dst.Tags = make([]string, len(src.Tags))
		for i, v := range src.Tags {
			dst.Tags[i] = v
		}
	}

	// goverter 生成的 Map 拷贝代码
	if src.Metadata != nil {
		dst.Metadata = make(map[string]string, len(src.Metadata))
		for k, v := range src.Metadata {
			dst.Metadata[k] = v
		}
	}
}
