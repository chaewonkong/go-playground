package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
)

// 로그 구조체 정의
type LogEntry struct {
	Level   string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
}

// 구조체를 map[string]interface{}로 변환하는 함수
func structToMap(log LogEntry) map[string]interface{} {
	return map[string]interface{}{
		"level":   log.Level,
		"service": log.Service,
		"message": log.Message,
	}
}

func main() {
	// CEL 환경 설정 (log를 map 타입으로 설정)
	env, err := cel.NewEnv(
		cel.Variable("log", cel.MapType(cel.StringType, cel.AnyType)), // log는 Map 형태로 변환하여 사용
	)
	if err != nil {
		fmt.Println("환경 설정 오류:", err)
		return
	}

	// 로그 필터 표현식 정의
	expression := `(log.level in ["ERROR", "WARN"]) && (log.service == "auth-service")`

	// 표현식 컴파일
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		fmt.Println("컴파일 오류:", issues.Err())
		return
	}

	// 프로그램 실행 준비
	prg, err := env.Program(ast)
	if err != nil {
		fmt.Println("프로그램 생성 오류:", err)
		return
	}

	// JSON 형태의 로그 데이터
	logData := []byte(`
	[
		{"level": "INFO", "service": "auth-service", "message": "User logged in"},
		{"level": "ERROR", "service": "auth-service", "message": "Authentication failed"},
		{"level": "WARN", "service": "payment-service", "message": "Payment delayed"},
		{"level": "WARN", "service": "auth-service", "message": "Token expired"}
	]`)

	// JSON을 []LogEntry 구조체로 변환
	var logs []LogEntry
	if err := json.Unmarshal(logData, &logs); err != nil {
		fmt.Println("JSON Unmarshal 오류:", err)
		return
	}

	// 필터링된 로그 출력
	fmt.Println("필터링된 로그:")
	for _, log := range logs {
		logMap := structToMap(log) // 구조체를 map 형태로 변환
		out, _, err := prg.Eval(map[string]interface{}{"log": logMap})
		if err != nil {
			fmt.Println("평가 오류:", err)
			continue
		}

		// 필터 조건을 만족하면 출력
		if out.Value().(bool) {
			fmt.Println(log)
		}
	}
}
