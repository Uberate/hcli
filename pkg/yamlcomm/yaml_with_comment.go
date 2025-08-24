package yamlcomm

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	commentTag = "comment" // Comment tag
	defaultTag = "default" // Default value tag
)

// MarshalWithComments converts the input struct to a YAML byte stream with comments
// Supports header comments, line-end type comments, default values, nested structures, ignored fields, and inline
func MarshalWithComments(input any) ([]byte, error) {
	root := &yaml.Node{}
	if err := structToYamlWithComments(reflect.ValueOf(input), root); err != nil {
		return nil, err
	}
	return yaml.Marshal(root)
}

// structToYamlWithComments recursively processes values, converting them to YAML nodes with comments
func structToYamlWithComments(v reflect.Value, node *yaml.Node) error {
	// Ensure the node always has a Kind property
	if node.Kind == 0 {
		// Default to scalar node
		node.Kind = yaml.ScalarNode
	}

	// Handle pointer types, get the pointed-to element
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// Encode nil pointer as null
			node.Kind = yaml.ScalarNode
			node.Tag = "!!null"
			return nil
		}
		return structToYamlWithComments(v.Elem(), node)
	}

	// If it's an interface type, get its actual value
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	// Handle struct types
	if v.Kind() == reflect.Struct {
		node.Kind = yaml.MappingNode
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := v.Field(i)

			// Parse yaml tag (name, ignore flag, inline flag)
			name, isIgnore, isInline := parseYamlTag(field.Tag.Get("yaml"))
			// Manually parse tag content
			tagStr := string(field.Tag)
			if isIgnore {
				continue // Skip ignored fields
			}

			// Handle inline tag: embed child struct fields into current node
			if isInline {
				// Recursively process child struct, result added directly to current node
				if err := structToYamlWithComments(fieldVal, node); err != nil {
					return fmt.Errorf("inline field %s: %w", field.Name, err)
				}
				continue // Already handled inline, no need to add current field
			}

			// Handle regular fields: create key-value pair nodes
			keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: name}
			valNode := &yaml.Node{}

			// Encode field value to value node
			if err := structToYamlWithComments(fieldVal, valNode); err != nil {
				return fmt.Errorf("encode field %s: %w", field.Name, err)
			}

			// Add header comment (comment tag)
			comment := parseCommentTag(tagStr)
			if comment != "" {
				keyNode.HeadComment = processComment(comment)
			}

			// Process line-end type and default value comments
			// Only add linecomment for basic types
			if isBasicType(field.Type) {
				fieldType := getSimpleTypeName(field.Type)
				lineComment := "# [" + fieldType + "]"
				if defVal := field.Tag.Get(defaultTag); defVal != "" && defVal != "nil" { // Ignore default value of "nil"
					lineComment = fmt.Sprintf("# [%s, {default=%s}]", fieldType, defVal)
				}
				valNode.LineComment = lineComment
			}

			// Add key-value pair to current node
			node.Content = append(node.Content, keyNode, valNode)
		}
		return nil
	}

	// Handle slices and arrays
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		node.Kind = yaml.SequenceNode
		for i := 0; i < v.Len(); i++ {
			itemNode := &yaml.Node{}
			if err := structToYamlWithComments(v.Index(i), itemNode); err != nil {
				return fmt.Errorf("encode slice item %d: %w", i, err)
			}

			node.Content = append(node.Content, itemNode)
		}
		return nil
	}

	// Handle maps
	if v.Kind() == reflect.Map {
		node.Kind = yaml.MappingNode
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()

			keyNode := &yaml.Node{}
			if err := structToYamlWithComments(key, keyNode); err != nil {
				return fmt.Errorf("encode map key: %w", err)
			}

			valNode := &yaml.Node{}
			if err := structToYamlWithComments(val, valNode); err != nil {
				return fmt.Errorf("encode map value: %w", err)
			}

			node.Content = append(node.Content, keyNode, valNode)
		}
		return nil
	}

	// Handle basic types
	// Create temporary node for encoding
	tempNode := &yaml.Node{}
	if err := tempNode.Encode(v.Interface()); err != nil {
		return err
	}
	// Copy temporary node properties to target node
	node.Kind = tempNode.Kind
	node.Tag = tempNode.Tag
	node.Value = tempNode.Value
	node.Content = tempNode.Content
	return nil
}

// parseCommentTag manually parses the comment tag
// parseCommentTag manually parses the comment tag
func parseCommentTag(tagStr string) string {
	// Find the start position of the comment tag
	commentStart := strings.Index(tagStr, "comment:")
	if commentStart == -1 {
		return ""
	}

	// Find the start position of the comment tag value (skip "comment:" prefix)
	valueStart := commentStart + len("comment:")

	// Check if there are quotes
	if valueStart < len(tagStr) && tagStr[valueStart] == '"' {
		// Handle quoted values
		valueStart++ // Skip opening quote
		valueEnd := valueStart

		// Find closing quote, taking care to handle escape characters
		for valueEnd < len(tagStr) {
			if tagStr[valueEnd] == '"' {
				// Check if it's an escaped quote
				if valueEnd > 0 && tagStr[valueEnd-1] == '\\' {
					valueEnd++
					continue
				}
				break
			}
			valueEnd++
		}

		if valueEnd < len(tagStr) {
			// Get the raw string value
			value := tagStr[valueStart:valueEnd]
			// Handle escaped newlines
			value = strings.ReplaceAll(value, "\\n", "\n")
			return value
		}
		return ""
	} else {
		// Handle unquoted values
		valueEnd := valueStart
		for valueEnd < len(tagStr) && tagStr[valueEnd] != ' ' {
			valueEnd++
		}
		// Get the raw string value
		value := tagStr[valueStart:valueEnd]
		// Handle escaped newlines
		value = strings.ReplaceAll(value, "\\n", "\n")
		return value
	}
}

// parseYamlTag parses the yaml tag, returning name, ignore flag, and inline flag
// yaml tag format: "name[,inline]" or "-" (means ignore)
func parseYamlTag(input string) (name string, isIgnore bool, isInline bool) {
	if input == "" {
		return "", false, false
	}

	parts := strings.Split(input, ",")
	// First part is the field name
	name = strings.TrimSpace(parts[0])
	if name == "-" {
		return "", true, false // Ignore this field
	}

	// Check for inline tag (case insensitive)
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if strings.EqualFold(strings.TrimSpace(part), "inline") {
				return name, false, true
			}
		}
	}

	return name, false, false
}

// getSimpleTypeName gets the simplified name of a type (for comment display)
func getSimpleTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return getSimpleTypeName(t.Elem()) // Hide pointer type '*' symbol
	case reflect.Slice, reflect.Array:
		return getSimpleTypeName(t.Elem()) + "[]"
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", getSimpleTypeName(t.Key()), getSimpleTypeName(t.Elem()))
	case reflect.Struct:
		return t.Name()
	default:
		return t.String()
	}
}

// processComment processes multi-line comments and handles auto-wrapping
func processComment(comment string) string {
	// Process multi-line comments
	commentLines := strings.Split(comment, "\n")
	var cleanedComments []string
	for _, line := range commentLines {
		line = strings.TrimSpace(line)
		if line != "" {
			if !strings.HasPrefix(line, "#") {
				line = "# " + line
			}
			// If line length exceeds 100 characters, auto-wrap
			// Need to consider the 2 characters occupied by '# ' prefix
			if len(line) > 100 {
				// Remove '# ' prefix before wrapping
				content := line[2:]
				wrappedLines := wrapText(content, 98) // 100 - 2('# ')
				// Re-add '# ' prefix to each wrapped line
				for i, wrappedLine := range wrappedLines {
					wrappedLines[i] = "# " + wrappedLine
				}
				cleanedComments = append(cleanedComments, wrappedLines...)
			} else {
				cleanedComments = append(cleanedComments, line)
			}
		}
	}
	// Join all comment lines
	return strings.Join(cleanedComments, "\n")
}

// wrapText wraps text at the specified width
func wrapText(text string, width int) []string {
	var lines []string
	for len(text) > width {
		// Find the nearest space position to break the line
		index := strings.LastIndex(text[:width], " ")
		if index == -1 {
			// If no space is found, force a line break at the specified width
			index = width
		}
		lines = append(lines, text[:index])
		text = text[index+1:]
	}
	if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}

// Example struct to demonstrate how to use comment tags
// type SomeStruct struct {
//     AStr string `comment:"Describe of next line key, may be description was mutiple.\nOther string of description." yaml:"AStr" default:"xxx"`
// }

// isBasicType determines whether a type is a basic type
func isBasicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return true
	default:
		return false
	}
}
