package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DataModel struct {
	DB *sql.DB
}

type Color struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	NameUr  string `json:"name_ur"`
	HexCode string `json:"hex_code"`
}

type Transmission struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type BodyType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type FuelType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Make struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	NameUr string `json:"name_ur"`
}

type Model struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	NameUr string `json:"name_ur"`
	MakeID int32  `json:"make_id"`
}

type Generation struct {
	ID        int32 `json:"id"`
	StartYear int32 `json:"start"`
	EndYear   int32 `json:"end"`
	ModelID   int32 `json:"model_id"`
}

type Version struct {
	ID             int32  `json:"id"`
	GenID          int32  `json:"gen_id"`
	ModelID        int32  `json:"model_id"`
	Name           string `json:"name"`
	NameUr         string `json:"name_ur"`
	Transmission   int16  `json:"transmission"`
	EngineCapacity int16  `json:"engine_capacity"`
	FuelType       int16  `json:"fuel_type"`
}

type Details struct {
	VersionID      int32 `json:"version_id"`
	Transmission   int32 `json:"transmission"`
	EngineCapacity int32 `json:"engine_capacity"`
	FuelType       int32 `json:"fuel_type"`
}

type City struct {
	ID      int32  `json:"version_id"`
	Name    string `json:"name"`
	NameUr  string `json:"name_ur"`
	Popular bool   `json:"popular"`
}

type Area struct {
	ID     int32  `json:"id"`
	CityID int32  `json:"city_id"`
	Name   string `json:"name"`
	NameUr string `json:"name_ur"`
}

type Registration struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	NameUr string `json:"name_ur"`
	Type   string `json:"type"`
}

func (m *DataModel) GetMakes() ([]*Make, error) {
	query := `
	SELECT id, name, name_ur  
	FROM data_makes; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	makes := []*Make{}

	for rows.Next() {
		var make Make
		err := rows.Scan(
			&make.ID,
			&make.Name,
			&make.NameUr,
		)
		if err != nil {
			return nil, err
		}

		makes = append(makes, &make)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return makes, nil
}

func (m *DataModel) GetModels() ([]*Model, error) {
	query := `
	SELECT id, name, name_ur, make_id  
	FROM data_models;  
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	models := []*Model{}

	for rows.Next() {
		var model Model
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.NameUr,
			&model.MakeID,
		)
		if err != nil {
			return nil, err
		}

		models = append(models, &model)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (m *DataModel) GetGenerations() ([]*Generation, error) {
	query := `
	SELECT id, start_year, end_year, model_id
	FROM data_generations;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	generations := []*Generation{}

	for rows.Next() {
		var generation Generation
		err := rows.Scan(
			&generation.ID,
			&generation.StartYear,
			&generation.EndYear,
			&generation.ModelID,
		)
		if err != nil {
			return nil, err
		}

		generations = append(generations, &generation)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return generations, nil
}

func (m *DataModel) GetVersions() ([]*Version, error) {
	query := `
	SELECT data_versions.id, data_versions.gen_id, data_versions.model_id,
	 data_versions.name, data_versions.name_ur, 
	 data_details.transmission_type, data_details.fuel_type,
	data_details.engine_capacity FROM data_versions Left
	 JOIN data_details ON data_details.version_id = data_versions.id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	versions := []*Version{}

	var tranmission sql.NullInt16
	var fueltype sql.NullInt16
	var engineCapacity sql.NullInt16

	for rows.Next() {
		var version Version
		err := rows.Scan(
			&version.ID,
			&version.GenID,
			&version.ModelID,
			&version.Name,
			&version.NameUr,
			&tranmission,
			&fueltype,
			&engineCapacity,
		)
		if err != nil {
			return nil, err
		}

		fmt.Println(tranmission, fueltype, engineCapacity)

		if tranmission.Valid {
			version.Transmission = tranmission.Int16
		} else {
			version.Transmission = 0
		}

		if fueltype.Valid {
			version.FuelType = fueltype.Int16
		} else {
			version.FuelType = 0
		}

		if engineCapacity.Valid {
			version.EngineCapacity = engineCapacity.Int16
		} else {
			version.EngineCapacity = 0
		}

		versions = append(versions, &version)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func (m *DataModel) GetColors() ([]*Color, error) {
	query := `
	SELECT id, name, name_ur, hex_code 
	FROM data_colors 
	WHERE version_id IS NULL;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	colors := []*Color{}

	for rows.Next() {
		var color Color
		err := rows.Scan(
			&color.ID,
			&color.Name,
			&color.NameUr,
			&color.HexCode,
		)
		if err != nil {
			return nil, err
		}

		colors = append(colors, &color)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return colors, nil
}

func (m *DataModel) GetTransmissions() ([]*Transmission, error) {
	query := `
	SELECT id, name
	FROM data_transmissions ;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	Transmissions := []*Transmission{}

	for rows.Next() {
		var transmission Transmission
		err := rows.Scan(
			&transmission.ID,
			&transmission.Name,
		)
		if err != nil {
			return nil, err
		}

		Transmissions = append(Transmissions, &transmission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Transmissions, nil
}

func (m *DataModel) GetBodyTypes() ([]*BodyType, error) {
	query := `
	SELECT id, name
	FROM data_body_types ; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	BodyTypes := []*BodyType{}

	for rows.Next() {
		var bodyType BodyType
		err := rows.Scan(
			&bodyType.ID,
			&bodyType.Name,
		)
		if err != nil {
			return nil, err
		}

		BodyTypes = append(BodyTypes, &bodyType)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return BodyTypes, nil
}

func (m *DataModel) GetFuelTypes() ([]*FuelType, error) {
	query := `
	SELECT id, name
	FROM fuel_types ; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	FuelTypes := []*FuelType{}

	for rows.Next() {
		var fuelType FuelType
		err := rows.Scan(
			&fuelType.ID,
			&fuelType.Name,
		)
		if err != nil {
			return nil, err
		}

		FuelTypes = append(FuelTypes, &fuelType)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return FuelTypes, nil
}

func (m *DataModel) GetCities() ([]*City, error) {
	query := `
	SELECT id, name, name_ur, popular  
	FROM cities; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	cities := []*City{}

	for rows.Next() {
		var city City
		err := rows.Scan(
			&city.ID,
			&city.Name,
			&city.NameUr,
			&city.Popular,
		)
		if err != nil {
			return nil, err
		}

		cities = append(cities, &city)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (m *DataModel) GetAreas() ([]*Area, error) {
	query := `
	SELECT id, name, name_ur, city  
	FROM areas; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	areas := []*Area{}

	for rows.Next() {
		var area Area
		err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.NameUr,
			&area.CityID,
		)
		if err != nil {
			return nil, err
		}

		areas = append(areas, &area)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return areas, nil
}

func (m *DataModel) GetRegistrations() ([]*Registration, error) {
	query := `
	SELECT id, name, name_ur, type  
	FROM registrations; 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	registrations := []*Registration{}

	for rows.Next() {
		var reg Registration
		err := rows.Scan(
			&reg.ID,
			&reg.Name,
			&reg.NameUr,
			&reg.Type,
		)
		if err != nil {
			return nil, err
		}

		registrations = append(registrations, &reg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}
