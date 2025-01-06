package status

import (
	"fmt"
	"github.com/LatticeBCLab/go-lattice/common/codes"
)

type Status struct {
	Code    codes.Code
	Message string
}

func (s *Status) code() codes.Code {
	if s == nil {
		return codes.OK
	}
	return s.Code
}

// Message returns the message contained in s.
func (s *Status) message() string {
	if s == nil {
		return ""
	}
	return s.Message
}

func (s *Status) Error() string {
	return s.String()
}

func (s *Status) String() string {
	return fmt.Sprintf("internal error: code = %d desc = %s", s.code(), s.message())
}

func (s *Status) Err() error {
	if s.code() == codes.OK {
		return nil
	}
	return s
}

func New(c codes.Code, msg string) *Status {
	return &Status{
		Code:    c,
		Message: msg,
	}
}

func Newf(code codes.Code, format string, args ...any) *Status {
	return New(code, fmt.Sprintf(format, args...))
}

func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}

func Errorf(c codes.Code, format string, a ...any) error {
	return Error(c, fmt.Sprintf(format, a...))
}
