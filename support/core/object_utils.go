package core

import "strings"

// ComputeDownloadURL 根据参数计算下载地址
// p1: 基本的URL
// p2: 对象的名称(路径)
func ComputeDownloadURL(p1, p2 string) string {
	count := 0
	if strings.HasSuffix(p1, "/") {
		count++
	}
	if strings.HasPrefix(p2, "/") {
		count++
	}
	switch count {
	case 0:
		return p1 + "/" + p2
	case 1:
		return p1 + p2
	default:
		return p1 + p2[1:]
	}
}
