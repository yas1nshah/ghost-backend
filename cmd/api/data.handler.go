package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func (app *application) updateMakes(w http.ResponseWriter, r *http.Request) {
	makes, err := app.models.Data.GetMakes()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize makes data to JSON format
	jsonData, err := json.MarshalIndent(makes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/makes.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data updated successfully"))
}

func (app *application) updateModels(w http.ResponseWriter, r *http.Request) {
	models, err := app.models.Data.GetModels()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize makes data to JSON format
	jsonData, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/models.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data updated successfully"))
}

func (app *application) updateGenerations(w http.ResponseWriter, r *http.Request) {
	generations, err := app.models.Data.GetGenerations()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize generations data to JSON format
	jsonData, err := json.MarshalIndent(generations, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/generations.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Generations data updated successfully"))
}

func (app *application) updateVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := app.models.Data.GetVersions()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize versions data to JSON format
	jsonData, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/versions.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Versions data updated successfully"))
}

func (app *application) updateColors(w http.ResponseWriter, r *http.Request) {
	colors, err := app.models.Data.GetColors()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize colors data to JSON format
	jsonData, err := json.MarshalIndent(colors, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/colors.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Colors data updated successfully"))
}

func (app *application) updateTransmissions(w http.ResponseWriter, r *http.Request) {
	transmissions, err := app.models.Data.GetTransmissions()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize transmissions data to JSON format
	jsonData, err := json.MarshalIndent(transmissions, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/transmissions.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transmissions data updated successfully"))
}

func (app *application) updateBodyTypes(w http.ResponseWriter, r *http.Request) {
	bodyTypes, err := app.models.Data.GetBodyTypes()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize body types data to JSON format
	jsonData, err := json.MarshalIndent(bodyTypes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/body_types.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Body types data updated successfully"))
}

func (app *application) updateFuelTypes(w http.ResponseWriter, r *http.Request) {
	fuelTypes, err := app.models.Data.GetFuelTypes()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize fuel types data to JSON format
	jsonData, err := json.MarshalIndent(fuelTypes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/fuel_types.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fuel types data updated successfully"))
}

func (app *application) updateCities(w http.ResponseWriter, r *http.Request) {
	fuelTypes, err := app.models.Data.GetCities()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize fuel types data to JSON format
	jsonData, err := json.MarshalIndent(fuelTypes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/cities.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cities data updated successfully"))
}

func (app *application) updateAreas(w http.ResponseWriter, r *http.Request) {
	fuelTypes, err := app.models.Data.GetAreas()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize fuel types data to JSON format
	jsonData, err := json.MarshalIndent(fuelTypes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/areas.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Areas data updated successfully"))
}

func (app *application) updateRegistration(w http.ResponseWriter, r *http.Request) {
	fuelTypes, err := app.models.Data.GetRegistrations()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Serialize fuel types data to JSON format
	jsonData, err := json.MarshalIndent(fuelTypes, "", "  ")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the JSON data to a file
	filePath := "./data/registrations.json"
	file, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Optionally, send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registrations data updated successfully"))
}
