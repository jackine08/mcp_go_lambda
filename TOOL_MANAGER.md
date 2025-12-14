# Tool Manager 사용 가이드

## 개요
ToolManager는 MCP tool을 자동으로 관리하는 시스템입니다. 새로운 tool provider를 추가하면 자동으로 등록됩니다.

## 아키텍처

```
ToolManager (registry.go)
    ├── ToolProvider interface
    │   └── RegisterTools(server)
    │
    ├── Calculator (calculator.go)
    │   ├── Add()
    │   ├── Multiply()
    │   ├── Subtract()
    │   └── Divide()
    │
    └── StringTools (string.go)
        ├── ToUpper()
        ├── ToLower()
        └── Reverse()
```

## 새로운 Tool Provider 추가하기

### 1단계: Tool Provider 구조체 생성

```go
package tools

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// 입력 타입 정의
type YourInput struct {
    Field string `json:"field" jsonschema:"field description"`
}

// Tool Provider 구조체
type YourTools struct{}

func NewYourTools() *YourTools {
    return &YourTools{}
}
```

### 2단계: Tool 메소드 구현

```go
// Tool 메소드 - MCP tool handler 시그니처 준수
func (t *YourTools) YourMethod(ctx context.Context, req *mcp.CallToolRequest, input YourInput) (
    *mcp.CallToolResult,
    map[string]interface{},
    error,
) {
    // 로직 구현
    result := "your result"
    
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{
                Text: result,
            },
        },
    }, map[string]interface{}{"result": result}, nil
}
```

### 3단계: RegisterTools 메소드 구현

```go
// ToolProvider interface 구현
func (t *YourTools) RegisterTools(server *mcp.Server) {
    // 메소드를 추가하면 여기에 등록
    mcp.AddTool(server, &mcp.Tool{
        Name:        "your_tool",
        Description: "설명",
    }, t.YourMethod)
    
    // 더 많은 tool 추가 가능
}
```

### 4단계: server.go에 등록

```go
// internal/server/server.go
toolManager := tools.NewToolManager()

toolManager.
    Register(tools.NewCalculator()).
    Register(tools.NewStringTools()).
    Register(tools.NewYourTools())  // ← 여기에 추가!

toolManager.RegisterAll(server)
```

## 완료! 🎉

이제 새로운 tool이 자동으로 MCP 서버에 등록됩니다.

## 예시: Math Tools 추가하기

```go
// internal/tools/math.go
package tools

import (
    "context"
    "math"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type NumberInput struct {
    Value float64 `json:"value" jsonschema:"the number"`
}

type MathTools struct{}

func NewMathTools() *MathTools {
    return &MathTools{}
}

func (m *MathTools) Sqrt(ctx context.Context, req *mcp.CallToolRequest, input NumberInput) (
    *mcp.CallToolResult,
    map[string]interface{},
    error,
) {
    result := math.Sqrt(input.Value)
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{
                Text: fmt.Sprintf("√%f = %f", input.Value, result),
            },
        },
    }, map[string]interface{}{"result": result}, nil
}

func (m *MathTools) RegisterTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "sqrt",
        Description: "제곱근을 계산합니다",
    }, m.Sqrt)
}
```

그리고 server.go에 한 줄만 추가:
```go
.Register(tools.NewMathTools())
```

## 장점

✅ **자동 등록**: Register()만 호출하면 모든 tool이 자동 등록  
✅ **체이닝**: Register() 여러 개를 체이닝으로 연결  
✅ **확장성**: 새 provider 추가가 간단함  
✅ **타입 안전성**: 컴파일 타임에 타입 체크  
✅ **명확한 구조**: 각 tool provider가 독립적인 파일로 분리
