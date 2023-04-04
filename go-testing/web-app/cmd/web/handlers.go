package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
	"webapp/pkg/data"
)

var pathToTemplates = "./templates/"

/*
*	- method with app receiver, as we'll need to share application data with handler
*		- as a handler, it needs something to write to
**/
func (app *application) Home(write http.ResponseWriter, req *http.Request) {
	var td = make(map[string]any)

	if app.Session.Exists(req.Context(), "test") {
		msg := app.Session.GetString(req.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(
			req.Context(),
			"test",
			"Hit this page at "+time.Now().UTC().String())
	}
	_ = app.renderPage(write, req, "home.page.gohtml", &TemplateData{Data: td})
}

func (app *application) Profile(write http.ResponseWriter, req *http.Request) {
	_ = app.renderPage(write, req, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP    string
	Data  map[string]any // any = alias to interface
	Error string
	Flash string
	User  data.User
}

func (app *application) renderPage(write http.ResponseWriter, request *http.Request, tmplt string, td *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(
		path.Join(pathToTemplates, tmplt),
		path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(write, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(request.Context())

	td.Error = app.Session.PopString(request.Context(), "error")
	td.Flash = app.Session.PopString(request.Context(), "flash")

	// if user exists in session, cast it to data.User and add to template's data
	if app.Session.Exists(request.Context(), "user") {
		td.User = app.Session.Get(request.Context(), "user").(data.User)
	}

	// execute the template, passing it data if any
	err = parsedTemplate.Execute(write, td)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) Login(write http.ResponseWriter, req *http.Request) {
	// parse form data
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(write, "bad request", http.StatusBadRequest)
		return
	}

	// validate data
	form := NewForm(req.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		// redirect to the login page with error message
		app.Session.Put(req.Context(), "error", "Invalid login credentials")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}

	email := req.Form.Get("email")
	password := req.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		// redirect to the login page with error message
		app.Session.Put(req.Context(), "error", "Invalid login!")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}

	// got the User. Authenticate the user.
	// if not authenticated then redirect with error
	if !app.authenticate(req, user, password) {
		app.Session.Put(req.Context(), "error", "Invalid login!")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}
	// prevent fixation attack
	_ = app.Session.RenewToken(req.Context())

	//store success message in session
	app.Session.Put(req.Context(), "flash", "Successfully logged in!")

	// redirect to some other page
	http.Redirect(write, req, "/user/profile", http.StatusSeeOther) // 303
}

func (app *application) authenticate(req *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}

	// user added to Session, that's how we know if logged in or not
	app.Session.Put(req.Context(), "user", user)

	return true
}

func (app *application) UploadProfilePic(w http.ResponseWriter, r *http.Request) {
	// call a function that extracts a file from an upload (request)
	files, err := app.UploadFiles(r, "./static/img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// get the user from the session - must exist, or we shouldn't be able to get to this handler function
	user := app.Session.Get(r.Context(), "user").(data.User)
	// create a var of type data.UserImage
	var i = data.UserImage{
		UserID:   user.ID,
		FileName: files[0].OriginalFileName, // yes, if different users upload files with same name, they'll clash. We're more interested in Tests for now.
	}

	// insert the user image into user_images
	_, err = app.DB.InsertUserImage(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
	* at this point, the user's info has changed.
	*	refresh the sessional variable "user"
	*		(to have the right, updated user info in the session)
	**/
	updatedUser, err := app.DB.GetUser(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "user", updatedUser)

	// redirect back to profile page
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

/*
* handle upload(s) of file(s)
* 	considers situations where there is more than one file in request.
* Obs: just a simple implementation.
*		- no test of mime type of uploaded files
*		- no generation of a random file name to be stored along with original file name
* 	of course, above would be required for a production system.
**/
func (app *application) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile // slice of pointers to uploaded files

	/* parse the form to have access to files:
	*  	- maximum determined file size: 5 megabytes. If bigger, return error.
	**/

	err := r.ParseMultipartForm(int64(1024 * 1024 * 5))
	if err != nil {
		return nil, fmt.Errorf("[file_size_limit] Uploaded file exceeds size of %d bytes.", 1024*1024*5)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, header := range fileHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := header.Open()
				if err != nil {
					return nil, err
				}
				// this is why we're using an inline function - because calling defer inside a for loop can cause resource leak.
				defer infile.Close()

				uploadedFile.OriginalFileName = header.Filename

				// this is where we're going to write bytes from request to save the image to our File System.
				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}
