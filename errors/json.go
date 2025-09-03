package errors

import "fmt"

type JsonRPCError struct {
	Code    ExitStatus `json:"code,string"`         // error code
	Message string     `json:"message"`       // The human readable error message associated to Code

}

func (e *JsonRPCError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}