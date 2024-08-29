package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"ghostprotocols.pk/internal/data"
	"ghostprotocols.pk/internal/validator"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		City     int64  `json:"city"`
		Address  string `json:"address"`
		Timings  string `json:"timings"`
		IsDealer bool   `json:"is_dealer"`
	}

	// time.Sleep(5 * time.Second)

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Phone: input.Phone,
		City:  input.City,
	}

	dealer := &data.Dealer{
		Address: input.Address,
		Timings: input.Timings,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.IsDealer {

		if data.ValidateDealer(v, user, dealer); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Users.InsertDealer(user, dealer)
	} else {

		if data.ValidateUser(v, user); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Users.InsertUser(user)
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicatePhone):
			v.AddError("phone", "a user with this phone already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user, "dealer": dealer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	var input struct {
		Name     string `json:"name"`
		City     int64  `json:"city"`
		Address  string `json:"address"`
		Timings  string `json:"timings"`
		IsDealer bool   `json:"is_dealer"`
	}

	user, err := app.models.Users.GetUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	user.City = input.City
	user.Name = input.Name

	dealer := &data.Dealer{
		Address: input.Address,
		Timings: input.Timings,
	}

	if input.IsDealer {
		dealer, err = app.models.Users.GetDealer(user.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		dealer.Address = input.Address
		dealer.Timings = input.Timings
	}

	v := validator.New()

	if input.IsDealer {

		if data.ValidateDealer(v, user, dealer); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Users.UpdateDealer(user, dealer)
	} else {

		if data.ValidateUser(v, user); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.Users.UpdateUser(user)
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicatePhone):
			v.AddError("phone", "a user with this phone already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user, "dealer": dealer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	user, err := app.models.Users.GetUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	dealer, err := app.models.Users.GetDealer(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			user.IsDealer = true
		}
	}

	// fmt.Printf(dealer.Address)

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user, "dealer": dealer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateProfilePicHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // Limit your file size to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the image file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine the format of the image
	format := header.Header.Get("Content-Type")

	var img image.Image

	switch format {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			http.Error(w, "Unable to decode JPEG image", http.StatusBadRequest)
			return
		}
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			http.Error(w, "Unable to decode PNG image", http.StatusBadRequest)
			return
		}
	case "image/webp":
		img, err = webp.Decode(file)
		if err != nil {
			http.Error(w, "Unable to decode WebP image", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unsupported image format", http.StatusBadRequest)
		return
	}

	// Calculate new dimensions while maintaining aspect ratio
	newHeight := uint(200)
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	newWidth := (newHeight * width) / height

	// Resize the image while maintaining aspect ratio
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Generate a UUID for the filename
	uuid := uuid.New().String()

	// Ensure the directory exists
	err = os.MkdirAll("./public/media/user-profile", os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create directory for images", http.StatusInternalServerError)
		fmt.Println("Error creating directory:", err)
		return
	}

	// Save the resized image as WebP
	err = saveImageAsWebP(resizedImg, "./public/media/user-profile/"+uuid+".webp")
	if err != nil {
		http.Error(w, "Unable to save resized image", http.StatusInternalServerError)
		fmt.Println("Error saving image:", err)
		return
	}

	err = app.models.Users.UpdateProfilePic(user.ID, (uuid + ".webp"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Return the UUID
	err = app.writeJSON(w, http.StatusCreated, envelope{"profile_url": uuid}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
