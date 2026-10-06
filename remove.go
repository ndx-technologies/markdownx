package markdownx

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	inlineLink    = regexp.MustCompile(`!?\[([^\]]*)\]\([^()]*\)`)
	referenceLink = regexp.MustCompile(`\[([^\]]*)\]\[[^\]]*\]`)
	delimiters    = regexp.MustCompile("[`*]+")
	strike        = regexp.MustCompile(`~~+`)
)

// Remove Markdown markers, keep words
func Remove(text string) string {
	var (
		lines   []string
		inFence bool
	)

	for line := range strings.SplitSeq(text, "\n") {
		line = strings.TrimSpace(line)

		if isFence(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			lines = append(lines, line)
			continue
		}
		if isDefinition(line) {
			continue
		}
		if isTableDelimiter(line) {
			continue
		}
		if isRule(line) {
			lines = append(lines, "")
			continue
		}

		lines = append(lines, cleanseLine(line))
	}

	return strings.Join(lines, "\n")
}

func cleanseLine(line string) string {
	for strings.HasPrefix(line, ">") {
		line = strings.TrimSpace(line[1:])
	}

	line = cleanseListItem(line)
	line = cleanseHeading(line)

	return strings.Join(strings.Fields(cleanseInline(line)), " ")
}

func cleanseListItem(line string) string {
	for _, bullet := range []string{"-", "*", "+", "•"} {
		rest, ok := strings.CutPrefix(line, bullet)
		if !ok || (rest != "" && !strings.HasPrefix(rest, " ") && !strings.HasPrefix(rest, "\t")) {
			continue
		}
		line = strings.TrimSpace(rest)
		break
	}

	if n := orderedListLen(line); n > 0 {
		line = strings.TrimSpace(line[n:])
	}

	return cleanseTaskBox(line)
}

func cleanseTaskBox(line string) string {
	runes := []rune(line)
	if len(runes) < 3 || runes[0] != '[' || runes[2] != ']' || runes[1] == '[' || runes[1] == ']' {
		return line
	}
	rest := runes[3:]
	if len(rest) > 0 && rest[0] != ' ' && rest[0] != '\t' {
		return line
	}
	return strings.TrimSpace(string(rest))
}

func orderedListLen(line string) int {
	i := 0
	for i < len(line) && i < 3 && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i == 0 || i == len(line) {
		return 0
	}
	if line[i] != '.' && line[i] != ')' {
		return 0
	}
	if i+1 < len(line) && line[i+1] != ' ' && line[i+1] != '\t' {
		return 0
	}
	return i + 1
}

func cleanseHeading(line string) string {
	n := 0
	for n < len(line) && n < 6 && line[n] == '#' {
		n++
	}
	if n == 0 {
		return line
	}
	if n < len(line) && line[n] != ' ' && line[n] != '\t' {
		return line
	}
	line = strings.TrimSpace(line[n:])
	if trimmed := strings.TrimRight(line, "#"); trimmed != line && strings.HasSuffix(trimmed, " ") {
		line = strings.TrimSpace(trimmed)
	}
	return line
}

func cleanseInline(line string) string {
	line = inlineLink.ReplaceAllString(line, "$1")
	line = referenceLink.ReplaceAllString(line, "$1")
	line = delimiters.ReplaceAllString(line, "")
	line = strike.ReplaceAllString(line, "")
	line = cleanseFootnoteRef(line)
	line = cleanseUnderscores(line)

	return strings.ReplaceAll(line, "|", "")
}

func cleanseFootnoteRef(line string) string {
	for {
		i := strings.Index(line, "[^")
		if i < 0 {
			return line
		}
		j := strings.IndexByte(line[i:], ']')
		if j < 0 {
			return line
		}
		line = line[:i] + line[i+j+1:]
	}
}

func cleanseUnderscores(line string) string {
	var (
		b     strings.Builder
		runes = []rune(line)
	)

	for i := 0; i < len(runes); {
		if runes[i] != '_' {
			b.WriteRune(runes[i])
			i++
			continue
		}

		j := i
		for j < len(runes) && runes[j] == '_' {
			j++
		}

		if isWordRune(runes, i-1) && isWordRune(runes, j) {
			b.WriteString(strings.Repeat("_", j-i))
		}

		i = j
	}

	return b.String()
}

func isWordRune(runes []rune, i int) bool {
	if i < 0 || i >= len(runes) {
		return false
	}
	return unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])
}

func isFence(line string) bool {
	return strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~")
}

func isDefinition(line string) bool {
	if !strings.HasPrefix(line, "[") {
		return false
	}
	i := strings.IndexAny(line[1:], "[]")
	if i < 0 {
		return false
	}
	return strings.HasPrefix(line[1+i:], "]:")
}

func isTableDelimiter(line string) bool {
	if !strings.Contains(line, "|") || !strings.Contains(line, "-") {
		return false
	}
	for cell := range strings.SplitSeq(line, "|") {
		if strings.Trim(cell, "-: ") != "" {
			return false
		}
	}
	return true
}

func isRule(line string) bool {
	line = strings.ReplaceAll(line, " ", "")
	if len(line) < 3 {
		return false
	}
	for _, r := range line {
		if r != '-' && r != '*' && r != '_' && r != '=' {
			return false
		}
	}
	return true
}
