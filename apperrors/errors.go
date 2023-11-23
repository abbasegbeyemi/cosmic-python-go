package apperrors

import (
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type NonExistentProductError struct {
	Sku domain.Sku
}

func (n NonExistentProductError) Error() string {
	return fmt.Sprintf("%s does not exist", n.Sku)
}
