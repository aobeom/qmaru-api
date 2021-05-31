package configs

var mode int = 1

// Deployment 部署模式 0: debug 1: release
func Deployment() bool {
	return mode == 0
}
