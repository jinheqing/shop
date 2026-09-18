// Package util — 通用工具函数
// 本文件提供地理数据脱敏工具，用于在 public API 出参层
// 自动 strip 精确到乡/镇/村/街道级别的茶园地址，
// 输出最多到"市/州"级别的模糊描述。
//
// 设计原因（合规红线）：
//   - 普洱茶核心产区（云南临沧/版纳/普洱、福建武夷山、四川雅安等）
//     多位于自然保护区、军事禁区、或航空管制山区，
//     对外发布精确坐标可能触碰《测绘法》第 26 条
//     "不得以任何形式擅自公开重要地理信息数据"
//     以及《国家安全法》第 77 条中关于"泄露国家秘密"的条款。
//   - 精确到村/乡的地址，即使没有坐标，也足以让懂行的人
//     通过地名反推具体山头位置（例如"班章村""曼岗村"），
//     属于敏感商业信息。
//   - 本函数作为**最后一道防线**：即使运营后台误录精确地址，
//     public API 出参也会自动脱敏，永不外泄。
package util

import (
	"regexp"
	"strings"
)

// SanitizeGardenLocation — 脱敏茶园位置字符串
// 输入可以是任何格式的地址文本，输出最多保留到"市/州"级别。
// 如果已经是模糊化格式（含 "·" 且不含精确行政层级词），则原样返回。
//
// 脱敏后格式示例：
//
//	"云南省临沧市临翔区邦东乡曼岗村茶园"        → "云南省 · 临沧市"
//	"云南省西双版纳州勐海县布朗山乡班章村茶园"  → "云南省 · 西双版纳州"
//	"云南省 · 临沧市 · 云雾茶区"                 → "云南省 · 临沧市 · 云雾茶区"  (已是模糊,原样)
//	"Yunnan · Lincang · Manzhang Village"       → "Yunnan · Lincang"
//	"Master Selection"                          → "Master Selection"  (不触发)
func SanitizeGardenLocation(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	// ====== 已模糊化检测 ======
	// 含 "·" 分隔符 + 不含精确层级词(乡/镇/村/街道) + 不含英文精确后缀
	// → 运营已经手动处理过，直接返回。
	if strings.Contains(raw, "·") {
		if !preciseLevelRe.MatchString(raw) && !preciseLevelEnRe.MatchString(raw) {
			return raw
		}
	}

	// ====== 中文地址 ======
	if chineseRe.MatchString(raw) {
		return sanitizeCN(raw)
	}

	// ====== 英文地址 ======
	return sanitizeEN(raw)
}

// --- 编译期正则 ---

// chineseRe — 是否含中文
var chineseRe = regexp.MustCompile(`\p{Han}`)

// preciseLevelRe — 精确行政层级词（比"市/州"更细）
// 注意：故意不含"区"，因为"云雾茶区""产区""功能区"里的"区"不是行政区。
// 临翔区这种市辖区属于中等粒度，保留不影响安全。
var preciseLevelRe = regexp.MustCompile(`县|乡|镇|村|街道|社区|行政村|行政村公所`)

// preciseLevelEnRe — 英文精确后缀
var preciseLevelEnRe = regexp.MustCompile(`(?i)village|township|county|hamlet|street|road|block|precinct`)

// prefectureRe — 州/市级关键词（需要保留的最细粒度）
// 顺序很重要：长词优先（"自治州"要比"州"先匹配）
var prefectureSuffixes = []string{"自治州", "地区", "特别行政区", "盟", "省辖市", "市", "州"}

// --- 脱敏实现 ---

func sanitizeCN(raw string) string {
	// 1. 带 "·" 分隔符的半模糊格式
	if strings.Contains(raw, "·") {
		return sanitizeCN_Dotted(raw)
	}

	// 2. 无分隔符的扁平格式
	return sanitizeCN_Flat(raw)
}

// sanitizeCN_Dotted — "云南省 · 临沧市 · 临翔区 · 邦东乡 · 曼岗村"
// 策略：逐段保留，遇到含州/市关键词的段就停
func sanitizeCN_Dotted(raw string) string {
	parts := strings.Split(raw, "·")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
		if containsAny(p, prefectureSuffixes) {
			break
		}
		// 如果当前段包含精确层级词（乡/镇/村），那也停（虽然我们前面已经过了hasDot检测，
		// 但这里再保险一次）
		if preciseLevelRe.MatchString(p) {
			out = out[:len(out)-1] // 把这个精确段去掉
			break
		}
	}
	return strings.Join(out, " · ")
}

// sanitizeCN_Flat — "云南省临沧市临翔区邦东乡曼岗村茶园"
// 策略：从左向右找第一个州/市/地区/盟，截断在它结束的位置
func sanitizeCN_Flat(raw string) string {
	// 先切掉末尾非行政后缀（茶园/茶区/茶山/古寨/寨子）
	trimmed := trailingNonAdminRe.ReplaceAllString(raw, "")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return raw
	}

	// 在整个字符串里找所有 prefecture suffixes 的第一个出现位置
	// 注意：是"第一个出现"而不是按优先级顺序——因为可能有"德宏州...潞西市"这种情况
	bestIdx := -1
	bestLen := 0
	for _, suf := range prefectureSuffixes {
		i := strings.Index(trimmed, suf)
		if i >= 0 && (bestIdx == -1 || i < bestIdx) {
			bestIdx = i
			bestLen = len(suf)
		}
	}

	if bestIdx < 0 {
		// 整个字符串没出现州/市/地区/盟——可能是"云南省临翔区..."不完整地址
		// 保留省 + 定性描述
		if i := strings.Index(trimmed, "省"); i >= 0 {
			return trimmed[:i+len("省")] + " · 茶山产区"
		}
		return raw
	}

	prefix := trimmed[:bestIdx+bestLen]
	// 确保省和市之间有 " · "
	prefix = addSeparatorBetweenProvinceAndCity(prefix)
	return prefix
}

// trailingNonAdminRe — 末尾的非行政后缀（茶园/茶区/茶山/古寨/寨子）
var trailingNonAdminRe = regexp.MustCompile(`(茶园|茶区|茶山|古寨|寨子|茶庄|茶林)$`)

// addSeparatorBetweenProvinceAndCity — 确保 "省" 之后的市/州 之间有 " · "
func addSeparatorBetweenProvinceAndCity(s string) string {
	// 已经有 "·" 就不动
	if strings.Contains(s, "·") {
		return s
	}
	// 在 "省" 之后插入 " · "
	if i := strings.Index(s, "省"); i >= 0 && i+1 < len(s) {
		return s[:i+len("省")] + " · " + s[i+len("省"):]
	}
	return s
}

// containsAny — s 是否包含 patterns 中的任意一个
func containsAny(s string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

// sanitizeEN — 英文地址脱敏，保留前两段（province + prefecture/city）
func sanitizeEN(raw string) string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == '·' || r == '|' || r == ','
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
		if len(out) >= 2 {
			break
		}
	}
	if len(out) == 0 {
		return raw
	}
	return strings.Join(out, " · ")
}
