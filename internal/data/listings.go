package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ListingsModel struct {
	DB *sql.DB
}

type Seller struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerfied  bool      `json:"email_verified"`
	Phone         string    `json:"phone"`
	PhoneVerified bool      `json:"phone_verified"`
	IsDealer      bool      `json:"is_dealer"`
	City          int64     `json:"city,omitempty"`
	Address       string    `json:"address,omitempty"`
	DateJoined    time.Time `json:"date_joined"`
	Timings       string    `json:"timings,omitempty"`
}

type Listing struct {
	ID        int       `json:"id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"-"`

	Active   bool `json:"active"`
	Featured bool `json:"featured"`

	GpManaged   bool `json:"gp_managed"`
	GpCertified bool `json:"gp_certified"`
	GpYard      bool `json:"gp_yard"`

	Gallery []Image `json:"gallery"`

	MakeID    int32  `json:"-"`
	Make      string `json:"make,omitempty"`
	ModelID   int32  `json:"-"`
	Model     string `json:"model,omitempty"`
	VersionID int32  `json:"-"`
	Version   string `json:"version,omitempty"`
	Year      int32  `json:"year,omitempty"`
	Price     int64  `json:"price,omitempty"`

	RegistrationID int32  `json:"-"`
	Registration   string `json:"registration,omitempty"`
	CityID         int32  `json:"-"`
	City           string `json:"city,omitempty"`
	AreaID         int32  `json:"-"`
	Area           string `json:"area,omitempty"`

	Mileage        string `json:"mileage,omitempty"`
	TransmissionID int16  `json:"-"`
	Transmission   string `json:"transmission"`
	FuelTypeID     int16  `json:"-"`
	FuelType       string `json:"fuel_type,omitempty"`
	EngineCapacity int32  `json:"engine_capacity,omitempty"`
	BodyTypeID     int16  `json:"-"`
	BodyType       string `json:"body_type,omitempty"`

	ColorID int32  `json:"-"`
	Color   string `json:"color,omitempty"`
	Details string `json:"details,omitempty"`

	SellerID int32  `json:"-"`
	Seller   Seller `json:"seller,omitempty"`

	UpVersion int32 `json:"-"`
}

type Image struct {
	Url   string `json:"url"`
	Order int16  `json:"order"`
}

func (m *ListingsModel) Insert(listing *Listing) error {
	// Serialize the Gallery field to JSON
	galleryJSON, err := json.Marshal(listing.Gallery)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO listings (gallery, make, model, version, year, price, 
    registration, city, area, mileage, transmission, fuel_type, engine_capacity, body_type,
    color, details, seller)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	RETURNING id, updated_at, created_at;`

	args := []any{
		galleryJSON,
		listing.MakeID,
		listing.ModelID,
		sql.NullInt32{Int32: listing.VersionID, Valid: listing.VersionID != 0},
		listing.Year,
		listing.Price,
		listing.RegistrationID,
		listing.CityID,
		sql.NullInt32{Int32: listing.AreaID, Valid: listing.AreaID != 0},
		listing.Mileage,
		listing.TransmissionID,
		listing.FuelTypeID,
		listing.EngineCapacity,
		listing.BodyTypeID,
		listing.ColorID,
		listing.Details,
		listing.SellerID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID, &listing.UpdatedAt, &listing.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: User does not have enough listing limit.`:
			return ErrListingLimitReached
		default:
			return err
		}
	}

	return nil
}

func (m *ListingsModel) GetById(id int64) (*Listing, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT 
    l.id, l.created_at, l.updated_at,
    l.active, l.featured, 
    l.gp_managed, l.gp_certified, l.gp_yard,
    l.gallery, 
    m.name AS make_name,
    mo.name AS model_name,
    COALESCE(v.name, '') AS version_name,
    l.year, l.price, 
    r.name AS registration_name,
    ci.name AS city_name,
    a.name AS area_name,
    l.mileage,
    t.name AS transmission_name,
    f.name AS fuel_type_name,
    l.engine_capacity, 
    b.name AS body_type_name,
    col.name AS color_name,
    l.details, 
    u.id AS seller_id,
    u.name AS seller_name,
    u.phone AS seller_phone,
    u.phone_verified AS phone_verified,
    u.email AS seller_phone,
    u.email_verified AS email_verified,
    CASE WHEN d.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_dealer
FROM listings l
LEFT JOIN data_makes m ON l.make = m.id
LEFT JOIN data_models mo ON l.model = mo.id
LEFT JOIN data_versions v ON l.version = v.id
LEFT JOIN registrations r ON l.registration = r.id
LEFT JOIN cities ci ON l.city = ci.id
LEFT JOIN areas a ON l.area = a.id
LEFT JOIN data_transmissions t ON l.transmission = t.id
LEFT JOIN fuel_types f ON l.fuel_type = f.id
LEFT JOIN data_body_types b ON l.body_type = b.id
LEFT JOIN data_colors col ON l.color = col.id
LEFT JOIN users u ON l.seller = u.id
LEFT JOIN dealers d ON u.id = d.user_id
WHERE l.id = $1;`

	var listing Listing

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	row := m.DB.QueryRowContext(ctx, query, id)

	var (
		galleryBytes []byte         // Declare galleryBytes as []byte to store JSON data
		versionName  sql.NullString // Use sql.NullString for nullable version_name
		areaName     sql.NullString // Use sql.NullString for nullable area_name
	)

	// Scan the row into Listing struct fields
	err := row.Scan(
		&listing.ID,
		&listing.CreatedAt,
		&listing.UpdatedAt,

		&listing.Active,
		&listing.Featured,

		&listing.GpManaged,
		&listing.GpCertified,
		&listing.GpYard,

		&galleryBytes,

		&listing.Make,
		&listing.Model,
		&versionName, // Scan version_name into sql.NullString
		&listing.Year,
		&listing.Price,

		&listing.Registration,
		&listing.City,
		&areaName, // Scan area_name into sql.NullString

		&listing.Mileage,
		&listing.Transmission,
		&listing.FuelType,
		&listing.EngineCapacity,

		&listing.BodyType,
		&listing.Color,
		&listing.Details,

		&listing.Seller.ID,
		&listing.Seller.Name,
		&listing.Seller.Phone,
		&listing.Seller.PhoneVerified,
		&listing.Seller.Email,
		&listing.Seller.EmailVerfied,
		&listing.Seller.IsDealer,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Handle version_name which is nullable
	if versionName.Valid {
		listing.Version = versionName.String
	} else {
		listing.Version = "" // Set to empty string if NULL
	}

	// Handle area_name which is nullable
	if areaName.Valid {
		listing.Area = areaName.String
	} else {
		listing.Area = "" // Set to empty string if NULL
	}

	// Unmarshal the JSON array bytes into []Image
	var gallery []Image
	if err := json.Unmarshal(galleryBytes, &gallery); err != nil {
		return nil, err
	}

	// Assign unmarshaled gallery to the listing
	listing.Gallery = gallery

	// Print the listing for debugging
	// fmt.Printf("Listing: %+v\n", listing)

	return &listing, nil
}

func (m *ListingsModel) GetForUpdate(id int64) (*Listing, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		id, created_at, updated_at,
		active,  featured,
		gp_managed, gp_certified, gp_yard,
		gallery,
		make, model, version, year, price,
		registration, city, area,
		mileage, transmission, fuel_type, engine_capacity, body_type,
		color, details, 
		seller, upversion
	FROM listings 

	WHERE id = $1;`

	var listing Listing

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	row := m.DB.QueryRowContext(ctx, query, id)

	var (
		galleryBytes []byte        // Declare galleryBytes as []byte to store JSON data
		versionName  sql.NullInt32 // Use sql.NullString for nullable version_name
		areaName     sql.NullInt32 // Use sql.NullString for nullable area_name
	)

	// Scan the row into Listing struct fields
	err := row.Scan(
		&listing.ID,
		&listing.CreatedAt,
		&listing.UpdatedAt,

		&listing.Active,
		&listing.Featured,

		&listing.GpManaged,
		&listing.GpCertified,
		&listing.GpYard,

		&galleryBytes,

		&listing.MakeID,
		&listing.ModelID,
		&versionName, // Scan version_name into sql.NullString
		&listing.Year,
		&listing.Price,

		&listing.RegistrationID,
		&listing.CityID,
		&areaName, // Scan area_name into sql.NullString

		&listing.Mileage,
		&listing.TransmissionID,
		&listing.FuelTypeID,
		&listing.EngineCapacity,

		&listing.BodyTypeID,
		&listing.ColorID,
		&listing.Details,

		&listing.Seller.ID,
		&listing.UpVersion,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Handle version_name which is nullable
	if versionName.Valid {
		listing.VersionID = versionName.Int32
	} else {
		listing.VersionID = 0 // Set to empty string if NULL
	}

	// Handle area_name which is nullable
	if areaName.Valid {
		listing.AreaID = areaName.Int32
	} else {
		listing.AreaID = 0 // Set to empty string if NULL
	}

	// Unmarshal the JSON array bytes into []Image
	var gallery []Image
	if err := json.Unmarshal(galleryBytes, &gallery); err != nil {
		return nil, err
	}

	// Assign unmarshaled gallery to the listing
	listing.Gallery = gallery

	// Print the listing for debugging
	// fmt.Printf("Listing: %+v\n", listing)

	return &listing, nil
}

func (m ListingsModel) Update(listing *Listing) error {
	query := `
	UPDATE listings 
	SET active = $1, featured = $2,
	gp_managed = $3, gp_certified = $4, gp_yard = $5,
	gallery = $6,
	make = $7, model = $8, version = $9, year = $10, price = $11,
	registration = $12, city = $13, area = $14,
	mileage = $15, transmission = $16, fuel_type = $17, engine_capacity = $18, body_type = $19,
	color = $20, details = $21, upversion = upversion + 1
	WHERE id = $22 AND upversion = $23
	RETURNING upversion;
	`

	args := []any{
		listing.Active, listing.Featured,
		listing.GpManaged, listing.GpCertified, listing.GpYard,
		listing.Gallery,
		listing.MakeID, listing.ModelID,
		sql.NullInt32{Int32: listing.VersionID, Valid: listing.VersionID != 0},
		listing.Year, listing.Price,
		listing.RegistrationID, listing.CityID,
		sql.NullInt32{Int32: listing.AreaID, Valid: listing.AreaID != 0},
		listing.Mileage, listing.TransmissionID, listing.FuelTypeID, listing.EngineCapacity, listing.BodyTypeID,
		listing.ColorID, listing.Details, listing.ID, listing.UpVersion,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&listing.UpVersion)
	if err != nil {
		switch {
		default:
			return err

		}
	}

	return nil
}

type ListingFilter struct {
	Make             int32        `json:"make,omitempty"`
	Model            int32        `json:"model,omitempty"`
	Version          int32        `json:"version,omitempty"`
	Year             NumberFilter `json:"year,omitempty"`
	City             int32        `json:"city,omitempty"`
	Area             int32        `json:"area,omitempty"`
	FuelType         int32        `json:"fuel_type,omitempty"`
	TransmissionAuto *bool        `json:"transmission_is_auto,omitempty"`
	Active           *bool        `json:"active,omitempty"`
	Featured         *bool        `json:"featured,omitempty"`
	GpManaged        *bool        `json:"gp_managed,omitempty"`
}

type NumberFilter struct {
	Start int32 `json:"start,omitempty"`
	End   int32 `json:"end,omitempty"`
}

func (m *ListingsModel) GetAll(f ListingFilter, s Sorting) ([]*Listing, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT 
			l.id, l.updated_at,
			l.active, l.featured, 
			l.gp_managed,l.gp_certified, l.gp_yard,
			l.gallery, 
			m.name AS make,
			mo.name AS model,
			COALESCE(v.name, '') AS version,
			l.year, l.price, 
			ci.name AS city,
			COALESCE(a.name, '') AS area,
			t.name AS transmission,
			l.mileage,
    		f.name AS fuel_type,
			COUNT(*) OVER()
		FROM listings l
		LEFT JOIN data_makes m ON l.make = m.id
		LEFT JOIN data_models mo ON l.model = mo.id
		LEFT JOIN data_versions v ON l.version = v.id
		LEFT JOIN cities ci ON l.city = ci.id
		LEFT JOIN areas a ON l.area = a.id
		LEFT JOIN data_transmissions t ON l.transmission = t.id
		LEFT JOIN fuel_types f ON l.fuel_type = f.id
		WHERE
			($1::INT IS NULL OR l.make = $1)
			AND ($2::INT IS NULL OR l.model = $2)
			AND ($3::INT IS NULL OR l.version = $3)
			AND ($4::INT IS NULL OR l.year >= $4)
			AND ($5::INT IS NULL OR l.year <= $5)
			AND ($6::INT IS NULL OR l.city = $6)
			AND ($7::INT IS NULL OR l.area = $7)
			AND ($8::INT IS NULL OR l.fuel_type = $8)
			AND ($9::INT IS NULL OR l.transmission = $9)
			AND ($10::BOOL IS NULL OR l.active = $10)
			AND ($11::BOOL IS NULL OR l.featured = $11)
			AND ($12::BOOL IS NULL OR l.gp_managed = $12)
		
		ORDER BY %s %s, l.id DESC
		LIMIT $13 OFFSET $14;

	`, s.sortColumn(), s.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		sql.NullInt32{Int32: f.Make, Valid: f.Make != 0},
		sql.NullInt32{Int32: f.Model, Valid: f.Model != 0},
		sql.NullInt32{Int32: f.Version, Valid: f.Version != 0},
		sql.NullInt32{Int32: f.Year.Start, Valid: f.Year.Start != 0},
		sql.NullInt32{Int32: f.Year.End, Valid: f.Year.End != 0},
		sql.NullInt32{Int32: f.City, Valid: f.City != 0},
		sql.NullInt32{Int32: f.Area, Valid: f.Area != 0},
		sql.NullInt32{Int32: f.FuelType, Valid: f.FuelType != 0},
		sql.NullBool{Bool: f.TransmissionAuto != nil && *f.TransmissionAuto, Valid: f.TransmissionAuto != nil},
		sql.NullBool{Bool: f.Active != nil && *f.Active, Valid: f.Active != nil},
		sql.NullBool{Bool: f.Featured != nil && *f.Featured, Valid: f.Featured != nil},
		sql.NullBool{Bool: f.GpManaged != nil && *f.GpManaged, Valid: f.GpManaged != nil},
		s.limit(),
		s.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	listings := []*Listing{}
	totalRecords := 0
	for rows.Next() {
		var (
			galleryBytes []byte // Declare galleryBytes as []byte to store JSON data
			// versionName  sql.NullString // Use sql.NullString for nullable version_name
			areaName sql.NullString // Use sql.NullString for nullable area_name
		)

		var listing Listing

		err := rows.Scan(
			&listing.ID,
			&listing.UpdatedAt,
			&listing.Active,
			&listing.Featured,
			&listing.GpManaged,
			&listing.GpCertified,
			&listing.GpYard,
			&galleryBytes,
			&listing.Make,
			&listing.Model,
			&listing.Version,
			&listing.Year,
			&listing.Price,
			&listing.City,
			&listing.Area,
			&listing.Transmission,
			&listing.Mileage,
			&listing.FuelType,
			&totalRecords,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// if versionName.Valid {
		// 	listing.Version = versionName.String
		// } else {
		// 	listing.Version = ""
		// }

		if areaName.Valid {
			listing.Area = areaName.String
		} else {
			listing.Area = ""
		}

		var gallery []Image
		if err := json.Unmarshal(galleryBytes, &gallery); err != nil {
			return nil, Metadata{}, err
		}

		listing.Gallery = gallery

		listings = append(listings, &listing)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, s.Page, s.PageSize)

	return listings, metadata, nil
}

// func (m *ListingsModel) GetAll(f ListingFilter, s Sorting) ([]*Listing, Metadata, error) {
// 	query := fmt.Sprintf(`
// 		SELECT
// 			l.id, l.updated_at, l.active, l.featured, l.gp_managed, l.gallery,
// 			m.name AS make_name,
// 			mo.name AS model_name,
// 			COALESCE(v.name, '') AS version_name,
// 			l.year, l.price,
// 			ci.name AS city_name,
// 			COALESCE(a.name, '') AS area_name,
// 			l.mileage,
// 			t.name AS transmission,
// 			f.name AS fuel_type_name,
// 			COUNT(*) OVER()
// 		FROM listings l
// 		LEFT JOIN data_makes m ON l.make = m.id
// 		LEFT JOIN data_models mo ON l.model = mo.id
// 		LEFT JOIN data_versions v ON l.version = v.id
// 		LEFT JOIN cities ci ON l.city = ci.id
// 		LEFT JOIN areas a ON l.area = a.id
// 		LEFT JOIN fuel_types f ON l.fuel_type = f.id
// 		LEFT JOIN data_transmissions t ON l.transmission = t.id
// 		WHERE
// 			($1::INT IS NULL OR l.make = $1)
// 			AND ($2::INT IS NULL OR l.model = $2)
// 			AND ($3::INT IS NULL OR l.version = $3)
// 			AND ($4::INT IS NULL OR l.year >= $4)
// 			AND ($5::INT IS NULL OR l.year <= $5)
// 			AND ($6::INT IS NULL OR l.city = $6)
// 			AND ($7::INT IS NULL OR l.area = $7)
// 			AND ($8::INT IS NULL OR l.fuel_type = $8)
// 			AND ($9::INT IS NULL OR l.transmission = $9)
// 			AND ($10::BOOL IS NULL OR l.active = $10)
// 			AND ($11::BOOL IS NULL OR l.featured = $11)
// 			AND ($12::BOOL IS NULL OR l.gp_managed = $12)
// 		ORDER BY %s %s, l.id DESC
// 		LIMIT $13 OFFSET $14;
// 	`, s.sortColumn(), s.sortDirection())

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	args := []interface{}{
// 		sql.NullInt32{Int32: f.Make, Valid: f.Make != 0},
// 		sql.NullInt32{Int32: f.Model, Valid: f.Model != 0},
// 		sql.NullInt32{Int32: f.Version, Valid: f.Version != 0},
// 		sql.NullInt32{Int32: f.Year.Start, Valid: f.Year.Start != 0},
// 		sql.NullInt32{Int32: f.Year.End, Valid: f.Year.End != 0},
// 		sql.NullInt32{Int32: f.City, Valid: f.City != 0},
// 		sql.NullInt32{Int32: f.Area, Valid: f.Area != 0},
// 		sql.NullInt32{Int32: f.FuelType, Valid: f.FuelType != 0},
// 		sql.NullBool{Bool: f.TransmissionAuto != nil && *f.TransmissionAuto, Valid: f.TransmissionAuto != nil},
// 		sql.NullBool{Bool: f.Active != nil && *f.Active, Valid: f.Active != nil},
// 		sql.NullBool{Bool: f.Featured != nil && *f.Featured, Valid: f.Featured != nil},
// 		sql.NullBool{Bool: f.GpManaged != nil && *f.GpManaged, Valid: f.GpManaged != nil},
// 		s.limit(),
// 		s.offset(),
// 	}

// 	rows, err := m.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, Metadata{}, fmt.Errorf("error querying database: %w", err)
// 	}
// 	defer rows.Close()

// 	listings := []*Listing{}
// 	totalRecords := 0
// 	for rows.Next() {
// 		var (
// 			galleryBytes []byte         // Declare galleryBytes as []byte to store JSON data
// 			versionName  sql.NullString // Use sql.NullString for nullable version_name
// 			areaName     sql.NullString // Use sql.NullString for nullable area_name
// 		)

// 		var listing Listing

// 		err := rows.Scan(
// 			&listing.ID,
// 			&listing.UpdatedAt,
// 			&listing.Active,
// 			&listing.Featured,
// 			&listing.GpManaged,
// 			&galleryBytes,
// 			&listing.Make,
// 			&listing.Model,
// 			&versionName,
// 			&listing.Year,
// 			&listing.Price,
// 			&listing.City,
// 			&areaName,
// 			&listing.Mileage,
// 			&listing.Transmission,
// 			&listing.FuelType,
// 			&totalRecords,
// 		)
// 		if err != nil {
// 			return nil, Metadata{}, err
// 		}

// 		if versionName.Valid {
// 			listing.Version = versionName.String
// 		} else {
// 			listing.Version = ""
// 		}

// 		if areaName.Valid {
// 			listing.Area = areaName.String
// 		} else {
// 			listing.Area = ""
// 		}

// 		var gallery []Image
// 		if err := json.Unmarshal(galleryBytes, &gallery); err != nil {
// 			return nil, Metadata{}, err
// 		}

// 		listing.Gallery = gallery

// 		listings = append(listings, &listing)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, Metadata{}, err
// 	}
// 	metadata := calculateMetadata(totalRecords, s.Page, s.PageSize)

// 	return listings, metadata, nil
// }
