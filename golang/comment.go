package golang

import "strings"

// formatComment replaces a '...' prefix with the ellipsisName and prefixes all lines
// with a '// '.
func formatComment(ellipsisName, doc string) string {
	doc = strings.TrimSpace(doc)
	if len(doc) > 0 {
		tmp := &strings.Builder{}
		if strings.HasPrefix(doc, "...") {
			tmp.WriteString(ellipsisName)
			tmp.WriteString(" ")
			tmp.WriteString(strings.TrimSpace(doc[3:]))
		} else {
			tmp.WriteString(doc)
		}
		str := tmp.String()
		tmp.Reset()
		for _, line := range strings.Split(str, "\n") {
			tmp.WriteString("// ")
			tmp.WriteString(line)
			tmp.WriteString("\n")
		}

		return tmp.String()
	}

	return ""
}

func deEllipsis(ellipsisName, doc string) string {
	tmp := &strings.Builder{}
	if strings.HasPrefix(doc, "...") {
		tmp.WriteString(ellipsisName)
		tmp.WriteString(" ")
		tmp.WriteString(strings.TrimSpace(doc[3:]))
	} else {
		tmp.WriteString(doc)
	}

	return tmp.String()
}
