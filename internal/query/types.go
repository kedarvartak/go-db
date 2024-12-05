package query

type Parser struct{}
type Optimizer struct{}
type Executor struct{}
type Result struct{}

func (p *Parser) Parse(query string) interface{} {
	return nil
}

func (o *Optimizer) Optimize(parsed interface{}) interface{} {
	return nil
}

func (e *Executor) Execute(plan interface{}) (Result, error) {
	return Result{}, nil
}
