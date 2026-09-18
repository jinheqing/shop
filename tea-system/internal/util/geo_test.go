package util

import "testing"

func TestSanitizeGardenLocation(t *testing.T) {
	cases := []struct{ in, want string }{
		// 典型精确地址 → 市/州级
		{"云南省临沧市临翔区邦东乡曼岗村茶园", "云南省 · 临沧市"},
		{"云南省西双版纳州勐海县布朗山乡班章村茶园", "云南省 · 西双版纳州"},
		{"云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园", "云南省 · 普洱市"},
		{"四川省雅安市天全县茶区", "四川省 · 雅安市"},
		// 自治州 / 地区
		{"云南省德宏州潞西市芒市镇", "云南省 · 德宏州"},
		// 已是模糊格式（含 · 且无精确层级词）→ 原样返回
		{"云南省 · 临沧市 · 云雾茶区", "云南省 · 临沧市 · 云雾茶区"},
		// 英文 → 保留前 2 段
		{"Yunnan · Lincang · Manzhang Village", "Yunnan · Lincang"},
		{"Yunnan · Lincang · Tea Mountain", "Yunnan · Lincang · Tea Mountain"},
		// 空 / 特殊值 → 原样
		{"", ""},
		{"Master Selection", "Master Selection"},
	}
	for _, c := range cases {
		got := SanitizeGardenLocation(c.in)
		if got != c.want {
			t.Errorf("Sanitize(%q) = %q, want %q", c.in, got, c.want)
		} else {
			t.Logf("OK  %q → %q", c.in, got)
		}
	}
}
