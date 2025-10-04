package yamlutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// FieldInfo represents a parsed struct field with metadata
type FieldInfo struct {
	Name      string
	Type      string
	Default   string
	Describe  string
	Value     interface{}
	YAMLName  string
	IsPointer bool
}

// ParseStruct parses a struct and returns field information
func ParseStruct(obj interface{}) ([]FieldInfo, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", val.Kind())
	}

	typ := val.Type()
	var fields []FieldInfo

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get YAML tag name or use field name
		yamlName := field.Name
		if yamlTag := field.Tag.Get("yaml"); yamlTag != "" {
			if yamlTag == "-" {
				continue
			}
			yamlName = yamlTag
		}

		desc := getTagValue(string(field.Tag), "describe")

		// Parse field info
		fieldInfo := FieldInfo{
			Name:      field.Name,
			YAMLName:  yamlName,
			Type:      getTypeName(fieldValue.Type()),
			Default:   field.Tag.Get("default"),
			Describe:  strings.ReplaceAll(desc, "\n", "\n# "),
			Value:     getFieldValue(fieldValue),
			IsPointer: fieldValue.Kind() == reflect.Ptr,
		}

		fields = append(fields, fieldInfo)
	}

	return fields, nil
}

// GenerateYAMLFromStruct generates YAML output from a struct
func GenerateYAMLFromStruct(obj interface{}) (string, error) {
	fields, err := ParseStruct(obj)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, field := range fields {
		// Add description comment if available
		if field.Describe != "" {
			builder.WriteString(fmt.Sprintf("# %s\n", field.Describe))
		}

		// Write field name and value
		value := getFormattedValue(field)
		builder.WriteString(fmt.Sprintf("%s: %s", field.YAMLName, value))

		// Add type comment for basic types
		if isBasicType(field.Type) {
			typeComment := fmt.Sprintf(" # [%s", field.Type)
			if field.Default != "" {
				typeComment += fmt.Sprintf(", default=%s", formatDefaultForComment(field.Default, field.Type))
			}
			typeComment += "]"
			builder.WriteString(typeComment)
		}

		builder.WriteString("\n\n")
	}

	return builder.String(), nil
}

// Render is the main function that returns YAML as bytes
func Render(obj interface{}) ([]byte, error) {
	yaml, err := GenerateYAMLFromStruct(obj)
	if err != nil {
		return nil, err
	}
	return []byte(yaml), nil
}

// Helper functions
func getTypeName(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		return getTypeName(typ.Elem())
	}

	switch typ.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct:
		return "struct"
	case reflect.Map:
		return "map"
	default:
		return typ.Name()
	}
}

func getFieldValue(value reflect.Value) interface{} {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	}
	return value.Interface()
}

func getFormattedValue(field FieldInfo) string {
	// Handle nil pointers by using default value
	if field.IsPointer && field.Value == nil && field.Default != "" {
		return formatDefaultValue(field.Default, field.Type)
	}

	// Handle zero values for basic types
	if field.Value == nil || (isBasicType(field.Type) && isZeroValue(field.Value)) {
		if field.Default != "" {
			return formatDefaultValue(field.Default, field.Type)
		}
		return getZeroValue(field.Type)
	}

	// For nested structs, generate YAML recursively
	if reflect.TypeOf(field.Value).Kind() == reflect.Struct {
		nestedYAML, err := GenerateYAMLFromStruct(field.Value)
		if err == nil {
			// Indent nested content
			lines := strings.Split(nestedYAML, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = "  " + line
				}
			}
			return "\n" + strings.Join(lines, "\n")
		}
	}

	// Format the actual value
	return formatValue(field.Value, field.Type)
}

func formatValue(value interface{}, typeName string) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case []string:
		return formatStringSlice(v)
	default:
		// Handle slices/arrays
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			return formatSlice(val)
		}
		// For other complex types, use simple representation
		return fmt.Sprintf("%v", v)
	}
}

func formatStringSlice(slice []string) string {
	if len(slice) == 0 {
		return "[]"
	}
	var builder strings.Builder
	for _, item := range slice {
		builder.WriteString(fmt.Sprintf("\n  - %q", item))
	}
	return builder.String()
}

func formatSlice(slice reflect.Value) string {
	if slice.Len() == 0 {
		return "[]"
	}

	var builder strings.Builder

	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface()

		// Handle nested structs in arrays
		if reflect.TypeOf(item).Kind() == reflect.Struct {
			nestedYAML, err := GenerateYAMLFromStruct(item)
			if err == nil {
				// Indent each line and replace first character with '-'
				lines := strings.Split(nestedYAML, "\n")
				for j, line := range lines {
					if line != "" {
						if j == 0 {
							// First line: replace with '-'
							if strings.HasPrefix(line, "#") {
								// Comment line: keep as is with proper indentation
								builder.WriteString(fmt.Sprintf("\n  - %s", line))
							} else {
								builder.WriteString(fmt.Sprintf("\n  - %s", line))
							}
						} else {
							// Other lines: add two spaces indentation
							if strings.HasPrefix(line, "#") {
								// Comment line: indent with 4 spaces
								builder.WriteString(fmt.Sprintf("\n    %s", line))
							} else {
								// Regular line: indent with 4 spaces
								builder.WriteString(fmt.Sprintf("\n    %s", line))
							}
						}
					}
				}
				builder.WriteString("\n")
			} else {
				builder.WriteString(fmt.Sprintf("\n  - %v", item))
			}
		} else if str, ok := item.(string); ok {
			builder.WriteString(fmt.Sprintf("\n  - %q", str))
		} else {
			builder.WriteString(fmt.Sprintf("\n  - %v", item))
		}
	}

	return builder.String()
}

func formatDefaultValue(defaultStr, typeName string) string {
	switch typeName {
	case "string":
		return fmt.Sprintf("%q", defaultStr)
	case "boolean":
		return defaultStr
	case "integer", "float":
		return defaultStr
	default:
		return defaultStr
	}
}

func formatDefaultForComment(defaultStr, typeName string) string {
	switch typeName {
	case "string":
		return fmt.Sprintf("%q", defaultStr)
	default:
		return defaultStr
	}
}

func getZeroValue(typeName string) string {
	switch typeName {
	case "string":
		return `""`
	case "boolean":
		return "false"
	case "integer", "float":
		return "0"
	default:
		return ""
	}
}

func isBasicType(typeName string) bool {
	return typeName == "string" || typeName == "boolean" || typeName == "integer" || typeName == "float"
}

func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case bool:
		return !v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v == 0
	case float32, float64:
		return v == 0.0
	default:
		return false
	}
}

// Utility function to convert first character to lowercase
func firstToLower(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func getTagValue(tagStr string, key string) string {
	// 将tag转换为字符串进行处理
	length := len(tagStr)
	i := 0

	for i < length {
		// 跳过空格
		for i < length && tagStr[i] == ' ' {
			i++
		}
		if i >= length {
			break
		}

		// 提取key
		keyStart := i
		for i < length && tagStr[i] != ':' && tagStr[i] != ' ' {
			i++
		}
		currentKey := tagStr[keyStart:i]

		// 检查是否找到目标key
		if currentKey != key {
			i = skipValue(tagStr, i)
			continue
		}

		// 找到目标key，提取其值
		if i >= length || tagStr[i] != ':' {
			return "" // 没有值
		}
		i++ // 跳过冒号

		// 解析值
		return parseValue(tagStr, &i)
	}

	return "" // 未找到key
}

// skipValue 跳过值部分，用于跳过非目标key的值
func skipValue(tag string, i int) int {
	length := len(tag)

	// 跳过空格
	for i < length && tag[i] == ' ' {
		i++
	}

	if i >= length {
		return i
	}

	// 如果是引号包裹的值，找到对应的结束引号
	if tag[i] == '"' || tag[i] == '\'' {
		quote := tag[i]
		i++
		for i < length && tag[i] != quote {
			// 处理转义的引号
			if tag[i] == '\\' && i+1 < length && tag[i+1] == quote {
				i += 2 // 跳过转义和引号
			} else {
				i++
			}
		}
		if i < length {
			i++ // 跳过结束引号
		}
	} else {
		// 非引号包裹的值，直到遇到空格或结束
		for i < length && tag[i] != ' ' {
			i++
		}
	}

	return i
}

// parseValue 解析值，处理引号和转义字符
func parseValue(tag string, i *int) string {
	length := len(tag)
	var result strings.Builder

	// 跳过空格
	for *i < length && tag[*i] == ' ' {
		*i++
	}

	if *i >= length {
		return ""
	}

	// 处理引号包裹的值
	if tag[*i] == '"' || tag[*i] == '\'' {
		quote := tag[*i]
		*i++ // 跳过开始引号

		for *i < length && tag[*i] != quote {
			// 处理转义字符
			if tag[*i] == '\\' && *i+1 < length {
				// 处理常见的转义序列
				switch tag[*i+1] {
				case 'n':
					result.WriteRune('\n')
				case 'r':
					result.WriteRune('\r')
				case 't':
					result.WriteRune('\t')
				case '"', '\'', '\\':
					result.WriteRune(rune(tag[*i+1]))
				default:
					// 保留未知的转义序列
					result.WriteRune('\\')
					result.WriteRune(rune(tag[*i+1]))
				}
				*i += 2
			} else {
				result.WriteRune(rune(tag[*i]))
				*i++
			}
		}

		if *i < length {
			*i++ // 跳过结束引号
		}
	} else {
		// 非引号包裹的值，直到遇到空格
		start := *i
		for *i < length && tag[*i] != ' ' {
			*i++
		}
		result.WriteString(tag[start:*i])
	}

	return result.String()
}
