// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/icza/backscanner"
	"golang.org/x/exp/slices"
)

const (
	docCommentTitleArgTag = "title"
	docCommentDescArgTag  = "description"
	docCommentEnumArgTag  = "enum"
)

func SetListFromDocIfFound(values *[]interface{}, valueTagName string, fieldDoc map[string]string) {
	var enumCSV string
	if ok := SetValueFromDocIfFound(&enumCSV, valueTagName, fieldDoc); ok {
		enumValues := strings.Split(enumCSV, ",")
		*values = make([]interface{}, 0, len(enumValues))
		for _, value := range enumValues {
			*values = append(*values, strings.TrimSpace(value))
		}
	}
}

func SetValueFromDocIfFound(value *string, valueTagName string, fieldDoc map[string]string) (ok bool) {
	if newValue, exists := fieldDoc[valueTagName]; exists && newValue != "" {
		*value = newValue
		ok = true
	}
	return
}

// GetDoc reads a given file and returns documentation
// when reverse is false the file is read from the beginning until the end byte position or the EOF (endBytePos < 0 reads the entire file)
// when reverse is true the file is read from the end byte position to the beginning.
func GetDoc(filename string, endBytePos int, reverse bool) (args map[string]string, err error) {
	head, args, _, err := parseCommentLines(filename, endBytePos, reverse)
	if err != nil {
		return
	}

	if _, ok := args[docCommentDescArgTag]; !ok && head != nil { // if there is no description tag use the head lines instead
		args[docCommentDescArgTag] = strings.Join(head, "\n")
	}

	return
}

// parseCommentLines parses docummentation comment to sections:
//
// # Quo in officia nobis autem pariatur sit tenetur ut dolores.		<--- head line
// # Deleniti asperiores quaerat.										<--- head line
// # @title: Incidunt aperiam sit facilis.								<--- argument tag
// # Voluptatem officiis aperiam.										<--- reminder
//
// It returns head (lines before the first argument tag with comment symbol removed), args (argument tags grouped by tag name), lines (list of all lines).
func parseCommentLines(filename string, endBytePos int, reverse bool) (head []string, args map[string]string, lines []string, err error) {
	if filename == "" {
		return
	}
	lines, err = readCommentLines(filename, endBytePos, reverse)
	if err != nil {
		return
	}

	head = make([]string, 0, len(lines))
	args = make(map[string]string, len(lines))
	for _, line := range lines {
		if groups := docCommentArgMatcher.FindStringSubmatch(line); groups != nil {
			args[groups[1]] = groups[2]
		}
		if len(args) < 1 { // all lines before the first argument tag form the comment head
			head = append(head, strings.TrimSpace(strings.TrimPrefix(line, "#")))
		}
	}

	return
}

// read all comment lines above a given end byte until the first non-comment or the end of file
func readCommentLines(filename string, endBytePos int, reverse bool) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	reader := NewLineReader(file, endBytePos, reverse)
	if err := reader.Read(&lines); err != nil {
		return lines, err
	}

	comments := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" && len(comments) < 1 {
			continue // ignore empty lines until the first comment above the end position
		} else if docCommentMatcher.MatchString(trimmedLine) {
			comments = append(comments, trimmedLine)
			continue
		}
		break
	}

	if reverse {
		slices.Reverse(comments) // lines are read bottom-up
	}
	return comments, nil
}

type LineReader interface {
	Read(lines *[]string) (err error)
}

func NewLineReader(fp *os.File, endBytePos int, reverse bool) LineReader {
	if reverse {
		return &backwardReader{
			reader:     fp,
			endBytePos: endBytePos,
		}
	}
	return &forwardReader{
		reader:     fp,
		endBytePos: int64(endBytePos),
	}
}

type forwardReader struct {
	reader     *os.File
	endBytePos int64
}

func (r *forwardReader) Read(lines *[]string) (err error) {
	var reader io.Reader = r.reader
	if r.endBytePos >= 0 {
		reader = io.LimitReader(r.reader, r.endBytePos)
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
	return scanner.Err()
}

type backwardReader struct {
	reader     *os.File
	endBytePos int
}

func (r *backwardReader) Read(lines *[]string) (err error) {
	scanner := backscanner.New(r.reader, r.endBytePos)
	for {
		line, _, err := scanner.Line()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*lines = append(*lines, line)
	}

	return nil
}
