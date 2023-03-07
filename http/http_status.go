package http

import (
	"fmt"
)

type Status int
type statusName string

const (
	StatusOK         = Status(200)
	StatusBadRequest = Status(400)

	StatusUnauthorized        = Status(401)
	StatusForbidden           = Status(403)
	StatusNotFound            = Status(404)
	StatusConflict            = Status(409)
	StatusTooManyRequests     = Status(429)
	StatusClientClosedRequest = Status(499)
	StatusInternalServerError = Status(500)
	StatusNotImplemented      = Status(501)
	StatusServiceUnavailable  = Status(503)
	StatusTimeout             = Status(504)
)

var statusToName = map[Status]statusName{
	StatusOK:                  "OK",
	StatusBadRequest:          "BadRequest",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "NotFound",
	StatusConflict:            "Conflict",
	StatusTooManyRequests:     "TooManyRequests",
	StatusClientClosedRequest: "ClientClosedRequest",
	StatusInternalServerError: "InternalServerError",
	StatusNotImplemented:      "NotImplemented",
	StatusServiceUnavailable:  "ServiceUnavailable",
	StatusTimeout:             "Timeout",
}

func (hs *Status) Code() int {
	return int(*hs)
}

func (hs *Status) name() string {
	return string(statusToName[*hs])
}

func (hs *Status) String() string {
	return fmt.Sprintf("%s(%v)", hs.name(), hs.Code())
}

func IsDefined(statusCode int) bool {
	status := Status(statusCode)
	_, found := statusToName[status]
	return found
}

func fromCode(statusCode int) (*Status, error) {
	if IsDefined(statusCode) {
		status := Status(statusCode)
		return &status, nil
	}
	return nil, fmt.Errorf("HTTP status for code %d is not defined", statusCode)
}
