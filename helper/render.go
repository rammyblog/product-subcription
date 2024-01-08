package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rammyblog/go-product-subscriptions/models"
)

type InputRequest struct {
	Data interface{}
}

func (u *InputRequest) Bind(r *http.Request) error {
	fieldErrMsg := "Field validation for '%s' failed on the '%s' tag"
	errMap := make(map[string]string)
	validate := validator.New()
	fmt.Println(u.Data.(models.User))
	err := validate.Struct(u.Data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMap[err.Field()] = fmt.Sprintf(fieldErrMsg, err.Field(), err.Tag())
		}
		jsonByte, jsonError := json.Marshal(errMap)
		if jsonError != nil {
			log.Fatal(jsonError)
		}
		return fmt.Errorf("%v", string(jsonByte))
	}
	return nil
}
