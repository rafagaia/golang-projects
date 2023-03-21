package main

import (
	"net/url"
	"strings"
)

/*
*	ERRORS:
**/

/* store all errors associated with form validations
*	- map with index string and entries as slice of string
*		since we might have more than one kind of error for a particular field
*	- made a type for errors so we can have functions associated with this type.
**/
type errors map[string][]string

// get first error message for a particular field
func (e errors) Get(field string) string {
	errorSlice := e[field]
	if len(errorSlice) == 0 {
		return ""
	}

	return errorSlice[0]
}

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

/*
*	FORMS:
**/

type Form struct {
	Data   url.Values //values of the URL in a form POST
	Errors errors
}

// Create a new variable of type Form
func NewForm(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: map[string][]string{},
	}
}

// determine if a particular field exists in a Form POST
func (f *Form) Has(field string) bool {
	x := f.Data.Get(field)
	if x == "" {
		return false
	}
	return true
}

/* make certain fields required
* can take more than one field in a single call to this function
**/
func (f *Form) Required(fields ...string) {
	// don't need index, just field itself
	for _, field := range fields {
		value := f.Data.Get(field)
		// get rid of leading empty space
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// generic check function
func (f *Form) Check(ok bool, key, message string) {
	if !ok {
		f.Errors.Add(key, message)
	}
}

// return true if there are no Errors, false if there's at least one error
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
