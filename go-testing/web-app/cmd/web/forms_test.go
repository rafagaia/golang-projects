package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Form_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("anything")
	if has {
		t.Error("Form shows Has field when it should not.")
	}

	postedData := url.Values{}
	// "a" as 'field' name, "b" as 'message' value
	postedData.Add("a", "b")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows Form does not have field when it should.")
	}
}

func Test_Form_Required(t *testing.T) {
	req := httptest.NewRequest("POST", "/anyurl", nil)
	form := NewForm(req.PostForm)

	form.Required("a", "ab", "ac") // requiring 3 forms

	if form.Valid() {
		t.Error("Form shows valid when required fields are missing.")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("ab", "ab")
	postedData.Add("ac", "ac")

	// using a "real request" now. Doesn't really make any difference here as opposed to httptest.NewRequest
	req, _ = http.NewRequest("POST", "/anyurl", nil)
	req.PostForm = postedData

	form = NewForm(req.PostForm)
	form.Required("a", "ab", "ac")
	if !form.Valid() {
		t.Error("Shows post does not have required fields, when it indeed does.")
	}
}
