package avatar

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"strings"
)

func GenerateTextAvatar(username string) string {
	// 1. 取用户名首字母（大写）
	initial := strings.ToUpper(string(username[0]))
	if initial == "" {
		initial = "U"
	}

	// 2. 根据用户名生成固定颜色（避免每次生成颜色不一样）
	hash := fnv.New32a()
	hash.Write([]byte(username))
	colorIndex := hash.Sum32() % 8
	colors := []string{
		"#4F46E5", // 靛蓝（和你的项目主色一致）
		"#8B5CF6", // 紫
		"#EC4899", // 粉
		"#10B981", // 绿
		"#3B82F6", // 蓝
		"#F59E0B", // 橙
		"#EF4444", // 红
		"#6B7280", // 灰
	}
	bgColor := colors[colorIndex]

	// 3. 生成 SVG 字符串（圆形头像，居中文字）
	svg := fmt.Sprintf(`
		<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
			<circle cx="50" cy="50" r="50" fill="%s"/>
			<text x="50" y="60" font-family="Arial" font-size="40" font-weight="bold" text-anchor="middle" fill="white">%s</text>
		</svg>
	`, bgColor, initial)

	// 4. 转 Base64 编码（方便存储到数据库）
	base64Avatar := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
	return base64Avatar
}
