package errorhandling

import (
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomErr .
type CustomErr struct {
	Code    int
	Message string
}

func (c CustomErr) Error() string {
	return fmt.Sprintf("%d:%s", c.Code, c.Message)
}

func New(code int, message string) CustomErr {
	return CustomErr{
		Code:    code,
		Message: message,
	}
}

func WrapErrorAsStatus(err CustomErr) *status.Status {
	return status.New(codes.Code(err.Code), err.Message)
}

func ParseCustomFromError(err error) *CustomErr {
	s, ok := status.FromError(err)
	if !ok {
		return &CustomErr{
			Code:    -1,
			Message: "未包装的错误：" + err.Error(),
		}
	}

	code, _ := strconv.Atoi(s.Code().String())
	return &CustomErr{
		Code:    code,
		Message: s.Message(),
	}
}
