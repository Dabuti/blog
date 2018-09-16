package models

import (
	"bufio"
	"bytes"
	"strings"
)

func getHTML(p *Post) string {
	body := string(p.Body)
	var buffer bytes.Buffer
	var tmp string

	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#+TITLE:") {
			tmp = strings.SplitAfter(line, "#+TITLE: ")[1]
			buffer.WriteString("<h1>" + tmp + "</h1>")
		} else if strings.HasPrefix(line, "#+AUTHOR:") {
			tmp = strings.SplitAfter(line, "#+AUTHOR: ")[1]
			buffer.WriteString("<h3>Author: " + tmp + "</h3>")
		} else if strings.HasPrefix(line, "* ") {
			tmp = strings.SplitAfter(line, "* ")[1]
			buffer.WriteString("<p>Section: " + tmp + "</p>")
		}
	}
	return buffer.String()
}
