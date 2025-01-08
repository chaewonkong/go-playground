package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

// cel go로 age가 minAge 이상인지 판단
func main() {
	env, err := cel.NewEnv(
		cel.Variable("age", cel.IntType),
		cel.Variable("minAge", cel.IntType),
	)

	if err != nil {
		log.Fatal(err)
	}

	ast, issues := env.Compile(`age >= minAge`)
	if issues != nil && issues.Err() != nil {
		log.Fatal(issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatal(err)
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"age":    27,
		"minAge": 19,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
