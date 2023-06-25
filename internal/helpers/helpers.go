package helpers

import (
	"fmt"
)

func IsNotNil(nullables ...any) error {
	for _, nullable := range nullables {
		if nullable == nil {
			return fmt.Errorf("err something is nil")
		}
	}

	return nil
}
