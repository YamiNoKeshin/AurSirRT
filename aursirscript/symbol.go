package aursirscript

import (
	"unicode/utf8"
)

func IsSymbol(code []byte) (isSymbol bool){

	rune, _ := utf8.DecodeRune(code)
	switch rune {

	case '#', '-', '>', '<' :
		isSymbol=true
	}

	return
}

func GetSymbol(code []byte) (Symbol []byte){

	/*//rune, size := utf8.DecodeRune(code)
	switch rune {

	case '#', '-', '>', '<' :
		//Symbol=rune

	case '-':
		if IsSymbol(code[size:]){
		//	Symbol = rune + GetSymbol(code[size])
		}
	}
*/
	return
}
