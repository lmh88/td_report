package curl

type AmazonauthInface interface {
	// GetProfiledList 获取profiledlist 列表
	GetProfiledList()
	GetProfiledByProfiled(profiled int64)
}
