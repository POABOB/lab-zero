// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: users.proto

package users

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on UserOrders with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserOrders) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserOrders with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserOrdersMultiError, or
// nil if none found.
func (m *UserOrders) ValidateAll() error {
	return m.validate(true)
}

func (m *UserOrders) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for BuyerId

	// no validation rules for SellerId

	// no validation rules for SellerAccount

	// no validation rules for GoodId

	// no validation rules for GoodName

	// no validation rules for GoodPrice

	// no validation rules for Status

	// no validation rules for Remark

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserOrdersValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserOrdersValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDeletedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserOrdersValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserOrdersValidationError{
				field:  "DeletedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UserOrdersMultiError(errors)
	}

	return nil
}

// UserOrdersMultiError is an error wrapping multiple validation errors
// returned by UserOrders.ValidateAll() if the designated constraints aren't met.
type UserOrdersMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserOrdersMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserOrdersMultiError) AllErrors() []error { return m }

// UserOrdersValidationError is the validation error returned by
// UserOrders.Validate if the designated constraints aren't met.
type UserOrdersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserOrdersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserOrdersValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserOrdersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserOrdersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserOrdersValidationError) ErrorName() string { return "UserOrdersValidationError" }

// Error satisfies the builtin error interface
func (e UserOrdersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserOrders.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserOrdersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserOrdersValidationError{}

// Validate checks the field values on GetUserInfoRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserInfoRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserInfoRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserInfoRequestMultiError, or nil if none found.
func (m *GetUserInfoRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserInfoRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := GetUserInfoRequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetUserInfoRequestMultiError(errors)
	}

	return nil
}

// GetUserInfoRequestMultiError is an error wrapping multiple validation errors
// returned by GetUserInfoRequest.ValidateAll() if the designated constraints
// aren't met.
type GetUserInfoRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserInfoRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserInfoRequestMultiError) AllErrors() []error { return m }

// GetUserInfoRequestValidationError is the validation error returned by
// GetUserInfoRequest.Validate if the designated constraints aren't met.
type GetUserInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserInfoRequestValidationError) ErrorName() string {
	return "GetUserInfoRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserInfoRequestValidationError{}

// Validate checks the field values on GetUserInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserInfoResponseMultiError, or nil if none found.
func (m *GetUserInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for Message

	for idx, item := range m.GetData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetUserInfoResponseValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetUserInfoResponseValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetUserInfoResponseValidationError{
					field:  fmt.Sprintf("Data[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetUserInfoResponseMultiError(errors)
	}

	return nil
}

// GetUserInfoResponseMultiError is an error wrapping multiple validation
// errors returned by GetUserInfoResponse.ValidateAll() if the designated
// constraints aren't met.
type GetUserInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserInfoResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserInfoResponseMultiError) AllErrors() []error { return m }

// GetUserInfoResponseValidationError is the validation error returned by
// GetUserInfoResponse.Validate if the designated constraints aren't met.
type GetUserInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserInfoResponseValidationError) ErrorName() string {
	return "GetUserInfoResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserInfoResponseValidationError{}
