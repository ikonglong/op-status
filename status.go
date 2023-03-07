package opstatus

import (
	"fmt"
	"log"
	"strings"

	"github.com/ikonglong/op-status/http"
)

type any interface{}

// A pseudo-enum of Status instances mapped 1:1 with the Codes. This simplifies construction
// patterns for derived instances of Status.
var (
	// StatusOK means the operation completed successfully.
	//
	// HTTP Mapping: 200 OK
	StatusOK = CodeOK.toStatus()

	// StatusCancelled means the operation was cancelled (typically by the caller).
	//
	// HTTP Mapping: 499 Client Closed Request
	StatusCancelled = CodeCancelled.toStatus()

	// StatusUnknown may be returned when a `Status` value received
	// from another address space belongs to an error space that is not known
	// in this address space. Also errors raised by APIs that do not return
	// enough error information may be converted to this error.
	//
	// HTTP Mapping: 500 Internal Server OpError
	StatusUnknown = CodeUnknown.toStatus()

	// StatusInvalidArgument means the client specified an invalid argument.
	// Note that this differs from `FAILED_PRECONDITION`. `INVALID_ARGUMENT`
	// indicates arguments that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	//
	// HTTP Mapping: 400 Bad Request
	StatusInvalidArgument = CodeInvalidArgument.toStatus()

	// StatusDeadlineExceeded means the deadline expired before the operation
	// could complete. For operations that change the state of the system,
	// this error may be returned, even if the operation has completed successfully.
	// For example, a successful response from a server could have been delayed long
	// enough for the deadline to expire.
	//
	// HTTP Mapping: 504 Gateway Timeout
	StatusDeadlineExceeded = CodeDeadlineExceeded.toStatus()

	// StatusNotFound means some requested entity (e.g., file or directory) was not found.
	//
	// Note to server developers: if a request is denied for an entire class
	// of users, such as gradual feature rollout or undocumented allow list,
	// `NOT_FOUND` may be used. If a request is denied for some users within
	// a class of users, such as user-based access control, `PERMISSION_DENIED`
	// must be used.
	//
	// HTTP Mapping: 404 Not Found
	StatusNotFound = CodeNotFound.toStatus()

	// StatusAlreadyExists means the entity that a client attempted to create
	// (e.g., file or directory) already exists.
	//
	// HTTP Mapping: 409 Conflict
	StatusAlreadyExists = CodeAlreadyExists.toStatus()

	// StatusPermissionDenied means the caller does not have permission to execute the specified
	// operation. `PERMISSION_DENIED` must not be used for rejections
	// caused by exhausting some resource (use `RESOURCE_EXHAUSTED`
	// instead for those errors). `PERMISSION_DENIED` must not be
	// used if the caller can not be identified (use `UNAUTHENTICATED`
	// instead for those errors). This error code does not imply the
	// request is valid or the requested entity exists or satisfies
	// other pre-conditions.
	//
	// HTTP Mapping: 403 Forbidden
	StatusPermissionDenied = CodePermissionDenied.toStatus()

	// StatusUnauthenticated means the request does not have valid authentication
	// credentials for the operation.
	//
	// HTTP Mapping: 401 Unauthorized
	StatusUnauthenticated = CodeUnauthenticated.toStatus()

	// StatusResourceExhausted means some resource has been exhausted,
	// perhaps a per-user quota, or perhaps the entire file system is out of space.
	//
	// HTTP Mapping: 429 Too Many Requests
	StatusResourceExhausted = CodeResourceExhausted.toStatus()

	// StatusFailedPrecondition means the operation was rejected because the system is not in
	// a state required for the operation's execution.  For example, the directory
	// to be deleted is non-empty, an rmdir operation is applied to
	// a non-directory, etc.
	//
	// Service implementors can use the following guidelines to decide
	// between `FAILED_PRECONDITION`, `ABORTED`, and `UNAVAILABLE`:
	//  (a) Use `UNAVAILABLE` if the client can retry just the failing call.
	//  (b) Use `ABORTED` if the client should retry at a higher level. For
	//      example, when a client-specified test-and-set fails, indicating the
	//      client should restart a read-modify-write sequence.
	//  (c) Use `FAILED_PRECONDITION` if the client should not retry until
	//      the system state has been explicitly fixed. For example, if an "rmdir"
	//      fails because the directory is non-empty, `FAILED_PRECONDITION`
	//      should be returned since the client should not retry unless
	//      the files are deleted from the directory.
	//
	// HTTP Mapping: 400 Bad Request
	StatusFailedPrecondition = CodeFailedPrecondition.toStatus()

	// StatusAborted means the operation was aborted, typically due to
	// a concurrency issue such as a sequencer check failure or transaction abort.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// HTTP Mapping: 409 Conflict
	StatusAborted = CodeAborted.toStatus()

	// StatusOutOfRange means the operation was attempted past the valid range.
	// E.g., seeking or reading past end-of-file.
	//
	// Unlike `INVALID_ARGUMENT`, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate `INVALID_ARGUMENT` if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// `OUT_OF_RANGE` if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between `FAILED_PRECONDITION` and
	// `OUT_OF_RANGE`.  We recommend using `OUT_OF_RANGE` (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an `OUT_OF_RANGE` error to detect when
	// they are done.
	//
	// HTTP Mapping: 400 Bad Request
	StatusOutOfRange = CodeOutOfRange.toStatus()

	// StatusUnimplemented means the operation is not implemented or is
	// not supported/enabled in this service.
	//
	// HTTP Mapping: 501 Not Implemented
	StatusUnimplemented = CodeUnimplemented.toStatus()

	// StatusInternal means internal errors. This means that some invariants expected by the
	// underlying system have been broken. This error code is reserved for serious errors.
	//
	// HTTP Mapping: 500 Internal Server OpError
	StatusInternal = CodeInternal.toStatus()

	// StatusUnavailable means the service is currently unavailable. This is most likely a
	// transient condition, which can be corrected by retrying with
	// a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// HTTP Mapping: 503 Service Unavailable
	StatusUnavailable = CodeUnavailable.toStatus()

	// StatusDataLoss means unrecoverable data loss or corruption.
	//
	// HTTP Mapping: 500 Internal Server OpError
	StatusDataLoss = CodeDataLoss.toStatus()
)

// statusList contains all the well-defined operation statuses indexed by their code values
var statusList = func() []Status {
	list := make([]Status, 0, len(codeToHTTPStatus))
	for _, code := range codeList {
		list = append(list, newStatus(code))
	}
	return list
}()

var codeToHTTPStatus = map[Code]http.Status{
	CodeOK:                 http.StatusOK,
	CodeInvalidArgument:    http.StatusBadRequest,
	CodeFailedPrecondition: http.StatusBadRequest,
	CodeOutOfRange:         http.StatusBadRequest,
	CodeUnauthenticated:    http.StatusUnauthorized,
	CodePermissionDenied:   http.StatusForbidden,
	CodeNotFound:           http.StatusNotFound,
	CodeAborted:            http.StatusConflict,
	CodeAlreadyExists:      http.StatusConflict,
	CodeResourceExhausted:  http.StatusTooManyRequests,
	CodeCancelled:          http.StatusClientClosedRequest,
	CodeDataLoss:           http.StatusInternalServerError,
	CodeUnknown:            http.StatusInternalServerError,
	CodeInternal:           http.StatusInternalServerError,
	CodeUnimplemented:      http.StatusNotImplemented,
	CodeUnavailable:        http.StatusServiceUnavailable,
	CodeDeadlineExceeded:   http.StatusTimeout,
}

var httpStatusToOpStatus = map[http.Status]Status{
	http.StatusOK:                  StatusOK,
	http.StatusBadRequest:          StatusInvalidArgument,
	http.StatusUnauthorized:        StatusUnauthenticated,
	http.StatusForbidden:           StatusPermissionDenied,
	http.StatusNotFound:            StatusNotFound,
	http.StatusConflict:            StatusAlreadyExists,
	http.StatusTooManyRequests:     StatusResourceExhausted,
	http.StatusClientClosedRequest: StatusCancelled,
	http.StatusInternalServerError: StatusInternal,
	http.StatusNotImplemented:      StatusUnimplemented,
	http.StatusServiceUnavailable:  StatusUnavailable,
	http.StatusTimeout:             StatusDeadlineExceeded,
}

// NewByHTTPStatus returns a copy of the status prototype mapped to given http status code.
func NewByHTTPStatus(statusCode int) *Status {
	if !http.IsDefined(statusCode) {
		unknownCopy := StatusUnknown
		return &unknownCopy
	}

	// Internally assure that there must be a unique op-status mapped to any defined https status
	// in order that the caller can take the fluid coding style.
	opStatus, found := httpStatusToOpStatus[http.Status(statusCode)]
	if found {
		log.Printf("[OpError] not found op-status mapped to given defined http status %v\n", statusCode)
	}
	return &opStatus
}

// NewWithCodeValue returns a copy of the status prototype mapped to given op status code.
func NewWithCodeValue(codeValue int) *Status {
	if codeValue < 0 || codeValue >= len(statusList) {
		return StatusUnknown.WithDescriptionf("Unknown op status code: %v", codeValue)
	}
	return &statusList[codeValue]
}

// NewWithCode returns a copy of the status prototype mapped to given op status code.
func NewWithCode(code Code) *Status {
	return &statusList[code.value]
}

// Status defines the status of an operation by providing a standard Code in conjunction with an
// optional Case and an optional description. Instances of Status are created by starting with the
// template for the appropriate Code and supplementing it with additional information:
//  StatusNotFound.WithDescription("Could not find 'important_file.txt'")
//
// The logical error model that Status defines is suitable for different programming environments,
// including REST APIs and RPC APIs.
type Status struct {
	code        Code
	theCase     Case
	description string
	details     map[string]any
}

func newStatus(code Code) Status {
	return Status{
		code: code,
	}
}

// WithDescription returns a derived instance of this Status with the given description. Leading and
// trailing whitespace is removed.
func (s *Status) WithDescription(description string) *Status {
	description = strings.TrimSpace(description)
	if s.description == description {
		copy := *s
		return &copy // return a copy of this Status
	}
	return &Status{
		code:        s.code,
		theCase:     s.theCase,
		description: description,
		details:     copyDetails(s.details),
	}
}

// WithDescriptionf returns a derived instance of this Status with the formatted description. Leading and
// trailing whitespace is removed.
func (s *Status) WithDescriptionf(descFmt string, fmtArgs ...any) *Status {
	return s.WithDescription(fmt.Sprintf(descFmt, fmtArgs))
}

// AugmentDescription returns a derived instance of this Status augmenting the current description
// with additional detail.
func (s *Status) AugmentDescription(additionalDetail string) *Status {
	if additionalDetail == "" {
		copy := *s
		return &copy // return a copy of this Status
	}

	newMsg := ""
	if s.description == "" {
		newMsg = additionalDetail
	} else {
		newMsg = s.description + "\n" + additionalDetail
	}
	return s.WithDescription(newMsg)
}

// WithCase returns a derived instance of this Status with the given case.
func (s *Status) WithCase(theCase Case) *Status {
	if s.theCase == theCase { // todo 深度比较 case
		copy := *s
		return &copy // return a copy of this Status
	}
	return &Status{
		code:        s.code,
		theCase:     theCase,
		description: s.description,
		details:     s.details,
	}
}

// WithCaseAndDesc returns a derived instance of this Status with the given case and description.
func (s *Status) WithCaseAndDesc(theCase Case, description string) *Status {
	description = strings.TrimSpace(description)
	if s.theCase == theCase && s.description == description { // todo 深度比较 case
		copy := *s
		return &copy
	}
	return &Status{
		code:        s.code,
		theCase:     theCase,
		description: description,
		details:     copyDetails(s.details),
	}
}

// WithCaseAndDescf returns a derived instance of this Status with the given case and formatted description.
func (s *Status) WithCaseAndDescf(theCase Case, descFmt string, fmtArgs ...any) *Status {
	desc := fmt.Sprintf(descFmt, fmtArgs)
	return s.WithCaseAndDesc(theCase, desc)
}

// AddDetail adds a detail about the failure.
func (s *Status) AddDetail(key string, value any) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	s.details[key] = value
}

// AddDetails adds details about the failure.
func (s *Status) AddDetails(details map[string]any) {
	for key, val := range details {
		s.AddDetail(key, val)
	}
}

func (s *Status) Code() Code {
	return s.code
}

func (s *Status) Description() string {
	return s.description
}

func (s *Status) TheCase() Case {
	return s.theCase
}

func (s *Status) Details() map[string]any {
	return s.details
}

// IsOK tells if this status is OK, i.e., not an error
func (s *Status) IsOK() bool {
	return s.code == CodeOK
}

// ToErrorCondition creates a string from this Status that describe current error condition
func (s *Status) ToErrorCondition() string {
	if s.description == "" {
		return s.code.String()
	}
	return s.code.String() + ": " + s.description
}

// RetryAdvice provides advice on retry for this status.
func (s *Status) RetryAdvice() RetryAdvice {
	advice := NoAdvice
	if s.code == CodeUnavailable {
		advice = JustRetryFailingCall
	} else if s.code == CodeFailedPrecondition {
		advice = NotRetryUntilStateFixed
	} else if s.code == CodeAborted || s.code == CodeResourceExhausted {
		advice = RetryAtHigherLevel
	}
	return advice
}

func copyDetails(details map[string]any) map[string]any {
	if details == nil {
		return map[string]any{}
	}
	copy := make(map[string]any, len(details)+2)
	for k, v := range details {
		copy[k] = v
	}
	return copy
}
