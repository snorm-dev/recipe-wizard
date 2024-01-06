package domerr

import "fmt"

type DomainErrorType int

const (
	NotFound DomainErrorType = iota
	UserNotFound
	Internal
	RecipeScraperFailure
	Forbidden
	DecodeJsonFailure
)

type DomainError struct {
	errorType DomainErrorType
	code      string
	message   string
}

func newDomainError(errorType DomainErrorType, code string, message string) *DomainError {
	return &DomainError{
		errorType,
		code,
		message,
	}
}

func (e *DomainError) Type() DomainErrorType {
	return e.errorType
}

func (e *DomainError) Code() string {
	return e.code
}

func (e *DomainError) Message() string {
	return e.message
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

var ErrNotFound *DomainError = newDomainError(NotFound, "not_found", "the requested resource does not exist")
var ErrUserNotFound *DomainError = newDomainError(UserNotFound, "user_not_found", "the user with the given id does not exist")
var ErrInternal *DomainError = newDomainError(Internal, "internal_error", "something went wrong")
var ErrRecipeScraperFailure *DomainError = newDomainError(RecipeScraperFailure, "recipe_scraper_failure", "the recipe scraper could not parse the given url")
var ErrForbidden *DomainError = newDomainError(Forbidden, "forbidden_access", "you do not have authorization to access that resource")
