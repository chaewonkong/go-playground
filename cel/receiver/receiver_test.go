package receiver_test

import (
	"cel/receiver"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
)

func TestReceiver(t *testing.T) {
	t.Run("TestReceiver", func(t *testing.T) {
		party := &receiver.Party{
			Attributes: []receiver.Attribute{
				{Name: "AvgScore", Type: "float", Value: 9.0, Default: 0},
				{Name: "AvgAge", Type: "float", Value: 16.0, Default: 0},
			},
		}

		attrs := []receiver.AttributeSchema{
			{Name: "AvgScore", Type: "float", Default: "0", Optional: false},
			{Name: "AvgAge", Type: "float", Default: "0", Optional: false},
		}
		expr := "party.AvgScore > 8.0 && party.AvgAge > 15.0 ? true : false"

		// CEL 환경 설정
		env, err := cel.NewEnv(
			receiver.NewPartyProvider(attrs),
		)
		assert.NoError(t, err)

		// 표현식 컴파일
		ast, issues := env.Compile(expr)
		assert.NoError(t, issues.Err())

		// 프로그램 실행 준비
		prg, err := env.Program(ast)
		assert.NoError(t, err)

		// 프로그램 실행
		out, _, err := prg.Eval(map[string]any{"party": party})
		assert.NoError(t, err)
		assert.Equal(t, true, out.Value())
	})
}
