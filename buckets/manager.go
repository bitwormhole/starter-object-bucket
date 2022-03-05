package buckets

// Manager 用来统一管理各种驱动的管理器
// 【inject:"#buckets.Manager"】
type Manager interface {
	FindDriver(name string) (Driver, error)
}
