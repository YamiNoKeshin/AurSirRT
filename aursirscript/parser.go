package aursirscript

type AurSirScriptParser struct {
	statements []scriptStatement
	currentStatement scriptStatement
}

func NewParser() (Parser AurSirScriptParser) {
	return
}

func (p *AurSirScriptParser) ParseLine(line []byte) (err error){
	is, err := IsComment(line);
	if is || err != nil{
		return
	}

	for len(line)!= 0 {
		literal, _ := GetLiteral(line)
		if IsSymbol(literal){
			p.TerminateStatement(literal)
		} else {
			p.AddArgumentToStatement(literal)
		}
	}


	return
}

func (p *AurSirScriptParser) TerminateStatement(NextSymbol []byte){
	if p.statements == nil {
		p.statements = []scriptStatement{}
	}
	p.statements = append(p.statements,p.currentStatement)
	p.currentStatement.Symbol = NextSymbol
	p.currentStatement.Arguments = [][]byte{}
}

func (p *AurSirScriptParser)AddArgumentToStatement(Argument []byte){
	p.currentStatement.Arguments = append(p.currentStatement.Arguments,Argument)
}
