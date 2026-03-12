package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0o755)
}

func ReadJSON(filePath string, dest any) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, dest)
}

func WriteJSON(filePath string, payload any) error {
	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return err
	}
	content, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, append(content, '\n'), 0o644)
}

func SHA256(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func TimestampUTC() string {
	return strings.ReplaceAll(strings.TrimSuffix(time.Now().UTC().Format(time.RFC3339), "Z"), ":", "") + "Z"
}

func Slugify(value string) string {
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.ToLower(value)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" {
		return "task"
	}
	return slug
}

func RenderJSON(payload any) string {
	content, _ := json.MarshalIndent(payload, "", "  ")
	return string(content)
}

func RenderText(payload any) string {
	value := reflect.ValueOf(payload)
	if !value.IsValid() {
		return ""
	}
	if value.Kind() == reflect.Struct {
		typ := value.Type()
		lines := make([]string, 0, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			name := field.Tag.Get("json")
			name = strings.TrimSuffix(name, ",omitempty")
			if name == "" {
				name = strings.ToLower(field.Name)
			}
			lines = append(lines, fmt.Sprintf("%s: %v", name, value.Field(i).Interface()))
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprint(payload)
}

func ReadSimpleYAML(filePath string) map[string]any {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return map[string]any{}
	}
	root := map[string]any{}
	type frame struct {
		indent int
		target map[string]any
	}
	stack := []frame{{indent: -1, target: root}}

	for _, rawLine := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(rawLine) == "" || strings.HasPrefix(strings.TrimSpace(rawLine), "#") {
			continue
		}
		indent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))
		trimmed := strings.TrimSpace(rawLine)
		parts := strings.SplitN(trimmed, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := ""
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}

		for len(stack) > 1 && indent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}

		current := stack[len(stack)-1].target
		if value == "" {
			child := map[string]any{}
			current[key] = child
			stack = append(stack, frame{indent: indent, target: child})
			continue
		}
		current[key] = normalizeScalar(value)
	}

	return root
}

func normalizeScalar(value string) any {
	unquoted := strings.Trim(value, `"'`)
	switch unquoted {
	case "true":
		return true
	case "false":
		return false
	}
	if number, err := strconv.Atoi(unquoted); err == nil {
		return number
	}
	return unquoted
}
