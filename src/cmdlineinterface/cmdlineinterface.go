package cmdlineinterface

import (
	"fmt"
	"os"
	"bufio"
)

type CmdLineInterface struct {

}

func (cli CmdLineInterface) Run(){
	reader := bufio.NewReader(os.Stdin)
	for {
	text, _ := reader.ReadString('\n')
	switch text{
	case "quit\n", "quit":
		return
	default:
		printHelp()
	}
	}
}
    func printHelp(){
		fmt.Println(`
	Type quit to exit
		`)
	}
