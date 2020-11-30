package client

import "fmt"

/** Custom error wrapper */
func wrapError(customMsg string, originalErr error) error {
	return fmt.Errorf("%s : %v", customMsg, originalErr)
}
