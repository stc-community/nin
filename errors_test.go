package nin

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	baseError := errors.New("test error")
	err := &Error{
		Err:  baseError,
		Type: ErrorTypePrivate,
	}
	assert.Equal(t, err.Error(), baseError.Error())
	assert.Equal(t, map[string]any{"error": baseError.Error()}, err.JSON())

	assert.Equal(t, err.SetType(ErrorTypePublic), err)
	assert.Equal(t, ErrorTypePublic, err.Type)

	assert.Equal(t, err.SetMeta("some data"), err)
	assert.Equal(t, "some data", err.Meta)
	assert.Equal(t, map[string]any{
		"error": baseError.Error(),
		"meta":  "some data",
	}, err.JSON())

	jsonBytes, _ := json.Marshal(err)
	assert.Equal(t, "{\"error\":\"test error\",\"meta\":\"some data\"}", string(jsonBytes))

	err.SetMeta(map[string]any{ //nolint: errcheck
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, map[string]any{
		"error":  baseError.Error(),
		"status": "200",
		"data":   "some data",
	}, err.JSON())

	err.SetMeta(map[string]any{ //nolint: errcheck
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, map[string]any{
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	}, err.JSON())

	type customError struct {
		status string
		data   string
	}
	err.SetMeta(customError{status: "200", data: "other data"}) //nolint: errcheck
	assert.Equal(t, customError{status: "200", data: "other data"}, err.JSON())
}

func TestErrorSlice(t *testing.T) {
	errs := errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
		{Err: errors.New("second"), Type: ErrorTypePrivate, Meta: "some data"},
		{Err: errors.New("third"), Type: ErrorTypePublic, Meta: map[string]any{"status": "400"}},
	}

	assert.Equal(t, errs, errs.ByType(ErrorTypeAny))
	assert.Equal(t, "third", errs.Last().Error())
	assert.Equal(t, []string{"first", "second", "third"}, errs.Errors())
	assert.Equal(t, []string{"third"}, errs.ByType(ErrorTypePublic).Errors())
	assert.Equal(t, []string{"first", "second"}, errs.ByType(ErrorTypePrivate).Errors())
	assert.Equal(t, []string{"first", "second", "third"}, errs.ByType(ErrorTypePublic|ErrorTypePrivate).Errors())
	assert.Empty(t, errs.ByType(ErrorTypeBind))
	assert.Empty(t, errs.ByType(ErrorTypeBind).String())

	assert.Equal(t, []any{
		map[string]any{"error": "first"},
		map[string]any{"error": "second", "meta": "some data"},
		map[string]any{"error": "third", "status": "400"},
	}, errs.JSON())
	jsonBytes, _ := json.Marshal(errs)
	assert.Equal(t, "[{\"error\":\"first\"},{\"error\":\"second\",\"meta\":\"some data\"},{\"error\":\"third\",\"status\":\"400\"}]", string(jsonBytes))
	errs = errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
	}
	assert.Equal(t, map[string]any{"error": "first"}, errs.JSON())
	jsonBytes, _ = json.Marshal(errs)
	assert.Equal(t, "{\"error\":\"first\"}", string(jsonBytes))

	errs = errorMsgs{}
	assert.Nil(t, errs.Last())
	assert.Nil(t, errs.JSON())
	assert.Empty(t, errs.String())
}

type TestErr string

func (e TestErr) Error() string { return string(e) }

// TestErrorUnwrap tests the behavior of nin.Error with "errors.Is()" and "errors.As()".
// "errors.Is()" and "errors.As()" have been added to the standard library in go 1.13.
func TestErrorUnwrap(t *testing.T) {
	innerErr := TestErr("some error")

	// 2 layers of wrapping : use 'fmt.Errorf("%w")' to wrap a nin.Error{}, which itself wraps innerErr
	err := fmt.Errorf("wrapped: %w", &Error{
		Err:  innerErr,
		Type: ErrorTypeAny,
	})

	// check that 'errors.Is()' and 'errors.As()' behave as expected :
	assert.True(t, errors.Is(err, innerErr))
	var testErr TestErr
	assert.True(t, errors.As(err, &testErr))
}
