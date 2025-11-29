package mcp

const (
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternal       = -32603
)

func errorResponse(id any, code int, message string) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
}

func toolErrorResponse(id any, err error) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: ToolResult{
			Content: []Content{
				{
					Type: "text",
					Text: err.Error(),
				},
			},
			IsError: true,
		},
	}
}
