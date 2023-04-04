package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
)

// helper function
func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

// helper function
func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}

func Test_application_handlers(t *testing.T) {
	var tests = []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/cat", http.StatusNotFound, "/cat", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	// when we're done with the test function, close the server:
	defer ts.Close()

	// ----> Custom HTTP Client - so we're not limited to the 10 redirects of default http client used in the for loop in this test function:
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // accept invalid ssl certificates (not signed) so we can use our custom client without throwing unknown certficate error
	}

	// create the client and override some default configs
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // so we only return the first response code
		},
	}

	// Custom HTTP Client <-----

	// range through test data
	for _, e := range tests {
		// call built-in test server, request to specified url
		response, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		// check for final status code getting back, as there's a status code before redirect, and one after
		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, response.StatusCode)
		}

		if response.Request.URL.Path != e.expectedURL {
			t.Errorf("%s: expected final url of %s; but got %s.", e.name, e.expectedURL, response.Request.URL.Path)
		}

		response2, _ := client.Get(ts.URL + e.url)
		if response2.StatusCode != e.expectedFirstStatusCode {
			t.Errorf("%s: expected first response status code to be %d; but got %d.", e.name, e.expectedFirstStatusCode, response2.StatusCode)
		}
	}
}

func Test_app_Home(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello, world!", "<small>From Session: hello, world!"},
	}

	for _, e := range tests {
		// create a request
		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context()) // make sure there's nothing in the Session before proceeding

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession) // "test" because its the expected key in the Home() Handler
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(res, req)

		// check status code
		if res.Code != http.StatusOK {
			t.Errorf("TestApp_Home returned wrong status code; expected 200 but got %d.", res.Code)
		}

		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: did not find %s in response body.", e.name, e.expectedHTML)
		}
	}
}

func Test_app_renderWithBadTemplate(t *testing.T) {
	// set pathToTemplates to a location with a bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	res := httptest.NewRecorder()

	err := app.renderPage(res, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("Expected error from bad template, but did not get one.")
	}

	pathToTemplates = "./../../templates/"
}

// table test
func Test_app_Login(t *testing.T) {
	var tests = []struct {
		name               string
		postedData         url.Values // what GO expects from a Form Post. Wrapper for some strings.
		expectedStatusCode int
		expectedLoc        string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"}, // can be sending a checkbox for example, not just a string. Thus the {}
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
		{
			name: "missing form data",
			postedData: url.Values{
				"email":    {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "user not found",
			postedData: url.Values{
				"email":    {"invalid@email.com"}, // can be sending a checkbox for example, not just a string. Thus the {}
				"password": {"anypassword"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "bad credentials",
			postedData: url.Values{
				"email":    {"admin@example.com"}, // can be sending a checkbox for example, not just a string. Thus the {}
				"password": {"wrongpassword"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, e := range tests {
		req, _ := http.NewRequest(
			"POST",
			"/login",
			strings.NewReader(e.postedData.Encode()), //Encode convers postedData to an io.Reader
		)
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // content-type that go expects to find from an html form post
		// ResponseRecorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d.", e.name, e.expectedStatusCode, rr.Code)
		}

		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != e.expectedLoc {
				t.Errorf("%s: expected location %s but got %s.", e.name, e.expectedLoc, actualLoc.String())
			}
		} else {
			t.Errorf("%s: no location header set.", e.name)
		}
	}
}

func Test_app_UploadFiles(t *testing.T) {
	// set up pipes (pipeRead, pipeWrite)
	pr, pw := io.Pipe()

	// create a new writer, of type *io.Writer
	writer := multipart.NewWriter(pw)

	// create a waitgroup, and add 1 to its counter
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// simulate uploading a file using a goroutine and our writer
	go simulatePNGUpload("./testdata/test_img.png", writer, wg, t)

	// read from the pipe, which receives data
	request := httptest.NewRequest("POST", "/anyurl", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	// call app.UploadFiles
	uploadedFiles, err := app.UploadFiles(request, "./testdata/uploads/")
	if err != nil {
		t.Error("Testing app UploadFiles error: ", err)
	}
	// perform our test cases
	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName)); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s.", err.Error()) // this is the only case we need to test, because if file does not exist, UploadFiles failed
	}

	// clean up any files created during test
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName))
}

func simulatePNGUpload(fileToUpload string, writer *multipart.Writer, wg *sync.WaitGroup, t *testing.T) {
	// so we don't have any resource leaks:
	defer writer.Close()
	// subtract 1 from our waitgroup's counter
	defer wg.Done()

	// create the multipart-form data field 'file' with value being filename
	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Error(err) // if we get the path name right, it should work fine without errors.
	}

	// open the actual file that we're simulating an upload
	f, err := os.Open(fileToUpload)
	if err != nil {
		t.Error(err) // if we get the path name right, it should work fine without errors.
	}
	// just to avoid resource leaks:
	defer f.Close()

	// decode the image (if we can upload a PNG, it should work for jpeg, gif, etc)
	img, _, err := image.Decode(f)
	if err != nil {
		t.Error("Error decoding image: ", err)
	}

	// write the PNG to our io.Writer
	err = png.Encode(part, img)
	if err != nil {
		t.Error(err)
	}
}
