package util

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnwrapAllErrors(t *testing.T) {
	err1 := errors.New("error1")
	err2 := fmt.Errorf("error2: %w", err1)
	err3 := fmt.Errorf("error3: %w", err2)
	err4 := fmt.Errorf("error4: %w", err3)

	require.ErrorIs(t, UnwrapAllErrors(err4), err1)
}
