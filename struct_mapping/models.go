package struct_mapping

// User 源结构体
type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	Address   string
	IsActive  bool
	Score     float64
	Tags      []string
	Metadata  map[string]string
	CreatedAt int64
}

// UserDTO 目标结构体
type UserDTO struct {
	ID        int
	Name      string
	Email     string
	Age       int
	Address   string
	IsActive  bool
	Score     float64
	Tags      []string
	Metadata  map[string]string
	CreatedAt int64
}
