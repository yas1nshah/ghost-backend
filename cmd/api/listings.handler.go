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
	"github.com/patrickmn/go-cache"
)

func (app *application) getListingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	listing, err := app.models.Listings.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"listing": listing}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createListingHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	var input struct {
		Gallery        []data.Image `json:"gallery"`
		Make           int32        `json:"make"`
		Model          int32        `json:"model"`
		Version        *int32       `json:"version"`
		Year           int32        `json:"year"`
		Price          int64        `json:"price"`
		Registration   int32        `json:"registration"`
		City           int32        `json:"city"`
		Area           *int32       `json:"area"`
		Mileage        string       `json:"mileage"`
		Transmission   int16        `json:"transmission"`
		FuelType       int16        `json:"fueltype"`
		EngineCapacity int32        `json:"engine_capacity"`
		BodyType       int16        `json:"body_type"`
		Color          int32        `json:"color"`
		Details        string       `json:"details"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Assign optional fields with default values if not provided
	versionID := int32(0)
	if input.Version != nil {
		versionID = *input.Version
	}

	areaID := int32(0)
	if input.Area != nil {
		areaID = *input.Area
	}

	listing := &data.Listing{
		Gallery:        input.Gallery,
		MakeID:         input.Make,
		ModelID:        input.Model,
		VersionID:      versionID,
		Year:           input.Year,
		Price:          input.Price,
		RegistrationID: input.Registration,
		CityID:         input.City,
		AreaID:         areaID,
		Mileage:        input.Mileage,
		TransmissionID: input.Transmission,
		FuelTypeID:     input.FuelType,
		EngineCapacity: input.EngineCapacity,
		BodyTypeID:     input.BodyType,
		ColorID:        input.Color,
		Details:        input.Details,
		SellerID:       int32(user.ID),
	}
	v := validator.New()

	err = app.models.Listings.Insert(listing)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrListingLimitReached):
			v.AddError("limit", "You have reached your Listing Limit")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) saveGalleryHandler(w http.ResponseWriter, r *http.Request) {
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
	newHeight := uint(750)
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	newWidth := (newHeight * width) / height

	// Resize the image while maintaining aspect ratio
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Generate a UUID for the filename
	uuid := uuid.New().String()

	// Ensure the directory exists
	err = os.MkdirAll("./listings/images", os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create directory for images", http.StatusInternalServerError)
		fmt.Println("Error creating directory:", err)
		return
	}

	// Save the resized image as WebP
	err = saveImageAsWebP(resizedImg, "./listings/images/"+uuid+".webp")
	if err != nil {
		http.Error(w, "Unable to save resized image", http.StatusInternalServerError)
		fmt.Println("Error saving image:", err)
		return
	}

	// Return the UUID
	fmt.Fprintln(w, uuid)
}

func saveImageAsWebP(img image.Image, filepath string) error {
	// Create the output file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	// Encode the image as WebP
	err = webp.Encode(out, img, &webp.Options{Lossless: true})
	if err != nil {
		return fmt.Errorf("error encoding webp: %w", err)
	}

	return nil
}

func (app *application) getListingsByFilter(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.ListingFilter
		data.Sorting
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Sorting.Page = app.readInt(qs, "page", 1, v)
	input.Sorting.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sorting.Sort = app.readString(qs, "sort", "updated")
	input.Sorting.SortSafelist = []string{"updated", "-updated"}
	if data.ValidateFilters(v, input.Sorting); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	input.ListingFilter.Make = int32(app.readInt(qs, "make", 0, v))
	input.ListingFilter.Model = int32(app.readInt(qs, "model", 0, v))
	input.ListingFilter.Version = int32(app.readInt(qs, "version", 0, v))

	input.ListingFilter.Year.Start, input.ListingFilter.Year.End = app.readRange(qs, "year", 0, v)
	input.ListingFilter.City = int32(app.readInt(qs, "city", 0, v))
	input.ListingFilter.Area = int32(app.readInt(qs, "area", 0, v))
	input.ListingFilter.FuelType = int32(app.readInt(qs, "fuel_type", 0, v))
	input.ListingFilter.TransmissionAuto = app.readBool(qs, "transmission_is_auto", v)
	input.ListingFilter.Active = app.readBool(qs, "active", v)
	input.ListingFilter.Featured = app.readBool(qs, "featured", v)
	input.ListingFilter.GpManaged = app.readBool(qs, "gp_managed", v)

	listings, metadata, err := app.models.Listings.GetAll(input.ListingFilter, input.Sorting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "listings": listings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getHomeFeed(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.ListingFilter
		data.Sorting
	}

	v := validator.New()
	input.Sorting.Page = 1
	input.Sorting.PageSize = 8
	input.Sorting.Sort = "updated"
	input.Sorting.SortSafelist = []string{"updated", "-updated"}
	if data.ValidateFilters(v, input.Sorting); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cacheKey := "home_feed"
	if cachedData, found := app.cache.Get(cacheKey); found {
		// Type assertion to convert cachedData back to envelope type
		if cachedEnvelope, ok := cachedData.(envelope); ok {
			err := app.writeJSON(w, http.StatusOK, cachedEnvelope, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}
	}

	setTrue := true
	setFalse := false
	input.ListingFilter.Active = &setTrue
	input.ListingFilter.Featured = &setTrue
	input.ListingFilter.GpManaged = &setFalse

	featuredListings, _, err := app.models.Listings.GetAll(input.ListingFilter, input.Sorting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	input.ListingFilter.Featured = &setFalse
	input.ListingFilter.GpManaged = &setTrue
	gpManagedListings, _, err := app.models.Listings.GetAll(input.ListingFilter, input.Sorting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	input.ListingFilter.GpManaged = &setFalse
	recentListings, _, err := app.models.Listings.GetAll(input.ListingFilter, input.Sorting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"featured_listings":   featuredListings,
		"gp_managed_listings": gpManagedListings,
		"recent_listings":     recentListings,
	}

	// Cache the response
	app.cache.Set(cacheKey, response, cache.DefaultExpiration)

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
