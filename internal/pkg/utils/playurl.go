package utils

import "strings"

const playURLMaskMarker = "/index.m3u8"

// MaskPlayURL 接口层播放地址脱敏：隐藏 /index.m3u8 前边四个字符。
// 例如 ...abcd/index.m3u8 -> ...****/index.m3u8；不足四位或不含该片段则原样返回。
func MaskPlayURL(raw string) string {
	if raw == "" {
		return raw
	}
	idx := indexPlayURLMarker(raw)
	if idx < 4 {
		return raw
	}
	return raw[:idx-4] + "****" + raw[idx:]
}

// IsMaskedPlayURL 判断地址是否已是脱敏形态，避免编辑回写时覆盖真实地址。
func IsMaskedPlayURL(raw string) bool {
	idx := indexPlayURLMarker(raw)
	if idx < 4 {
		return false
	}
	return raw[idx-4:idx] == "****"
}

func indexPlayURLMarker(raw string) int {
	if idx := strings.Index(raw, playURLMaskMarker); idx >= 0 {
		return idx
	}
	lower := strings.ToLower(raw)
	return strings.Index(lower, playURLMaskMarker)
}
