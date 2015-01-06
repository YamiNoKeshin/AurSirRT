package test

import (
	"bufio"
	"os"
	"testing"

	"github.com/joernweissenborn/aursirrt/aursirscript"
)

func readLines(path string) ([]string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

var code []string = readLines("testcode.aursirscript")



func TestParseComment(t *testing.T){
	comment := []byte(code[0])
	noncomment := []byte(code[2])
	is, err := aursirscript.IsComment(comment)

	if !is || err != nil {
		t.Error("Could not parse Comment")
	}

	is, err = aursirscript.IsComment(noncomment)

	if is || err != nil {
		t.Error("Could not parse non comment")
	}
	_, err = aursirscript.IsComment([]byte("/ you shall not parse"))
	_,isCommentError := err.(aursirscript.CommentError)
	if !isCommentError {
		t.Error("Could not detect Malformed Comment")
	}
}
