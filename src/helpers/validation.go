package helpers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func FieldRequiredValidation(fieldValue interface{}, require string) *ErrorResponse {
	var errorResponse *ErrorResponse

	v := validator.New()
	err := v.Var(fieldValue, require)

	if err != nil {
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			return &ErrorResponse{ErrorMessage: err.Error()}
		}

		for _, err := range validationErr {
			errorResponse = &ErrorResponse{
				ErrorMessage: fmt.Sprintf("this field must contain %s", err.ActualTag()),
			}
			break
		}
	}

	return errorResponse
}

func StructValidation(param any) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validator.New().Struct(param)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			if err.Param() == "" {
				message = fmt.Sprintf("%s must contain %s", err.Field(), err.ActualTag())
			} else {
				message = fmt.Sprintf("%s must contain %s=%s", err.Field(), err.ActualTag(), err.Param())
			}

			errors = append(errors, &ErrorResponse{
				ErrorMessage: message,
			})
		}
	}
	return errors
}

func PasswordValidation(password string, errors []*ErrorResponse) []*ErrorResponse {
	for _, err := range errors {
		if strings.Contains(fmt.Sprintf("%s", err), "password") {
			return errors
		}
	}

	uppercasePassword := regexp.MustCompile(`[A-Z]`)
	spaceProhibitedPassword := regexp.MustCompile(`[\s]`)
	numberPassword := regexp.MustCompile(`[0-9]`)
	specialPassword := regexp.MustCompile(`[\W_]`)

	if !uppercasePassword.MatchString(password) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "password must contain at least one uppercase letter",
		})
	} else if spaceProhibitedPassword.MatchString(password) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "password must contain no space",
		})
	} else if !numberPassword.MatchString(password) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "password must contain at least one digit number",
		})
	} else if !specialPassword.MatchString(password) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "password must contain at least one special letter",
		})
	}

	return errors
}

func ImageValidation(file *multipart.FileHeader) []*ErrorResponse {
	var errors []*ErrorResponse

	if file.Size > int64(2<<20) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "image must contain less than 2 MB",
		})
	}

	fileHeader, err := file.Open()
	if err != nil {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "failed to open image",
		})
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	if _, err := fileHeader.Read(buffer); err != nil {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "image failed to read image",
		})
	}

	fileType := http.DetectContentType(buffer)
	validFileTypes := []string{"image/png", "image/jpeg", "image/jpg"}
	if !isValidFileType(validFileTypes, fileType) {
		errors = append(errors, &ErrorResponse{
			ErrorMessage: "image must contain .png|.jpg|.jpeg",
		})
	}

	return errors
}

func isValidFileType(validFileTypes []string, fileType string) bool {
	for _, validType := range validFileTypes {
		if validType == fileType {
			return true
		}
	}
	return false
}
