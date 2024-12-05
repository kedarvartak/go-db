package query

type QueryProcessor struct {
    parser    *Parser
    optimizer *Optimizer
    executor  *Executor
}

func (qp *QueryProcessor) Execute(query string) (Result, error) {
    // 1. Parse the query
    parsed := qp.parser.Parse(query)
    
    // 2. Optimize the query 
    plan := qp.optimizer.Optimize(parsed)
    
    // 3. Execute the plan
    return qp.executor.Execute(plan)
} 