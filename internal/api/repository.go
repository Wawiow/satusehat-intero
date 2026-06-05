package api

import (
	"database/sql"
	"fmt"
	"strings"
)

func savePatientLocal(db *sql.DB, p PersonResponse) error {
	query := `INSERT INTO patients (nik, id, ihs_number, name, gender, birth_date, phone, address) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(nik) DO UPDATE SET 
			  id=excluded.id, ihs_number=excluded.ihs_number, name=excluded.name, gender=excluded.gender, 
			  birth_date=excluded.birth_date, phone=excluded.phone, address=excluded.address`
	_, err := db.Exec(query, p.NIK, p.ID, p.IHSNumber, p.Name, p.Gender, p.BirthDate, p.Phone, p.Address)
	return err
}

func getPatientLocal(db *sql.DB, nik string) (*PersonResponse, error) {
	var p PersonResponse
	query := `SELECT nik, id, ihs_number, name, gender, birth_date, phone, address FROM patients WHERE nik = ?`
	err := db.QueryRow(query, nik).Scan(&p.NIK, &p.ID, &p.IHSNumber, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func getAllPatientsLocal(db *sql.DB) ([]PersonResponse, error) {
	query := `SELECT nik, id, ihs_number, name, gender, birth_date, phone, address FROM patients`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []PersonResponse
	for rows.Next() {
		var p PersonResponse
		if err := rows.Scan(&p.NIK, &p.ID, &p.IHSNumber, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.Address); err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func saveLocationLocal(db *sql.DB, loc LocationResponse) error {
	query := `INSERT INTO locations (id, identifier_value, name, description, phone) 
	          VALUES (?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET 
			  identifier_value=excluded.identifier_value, name=excluded.name, 
			  description=excluded.description, phone=excluded.phone`
	_, err := db.Exec(query, loc.ID, loc.IdentifierValue, loc.Name, loc.Description, loc.Phone)
	return err
}

func getLocationLocal(db *sql.DB, id string) (*LocationResponse, error) {
	var loc LocationResponse
	query := `SELECT id, identifier_value, name, description, phone FROM locations WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&loc.ID, &loc.IdentifierValue, &loc.Name, &loc.Description, &loc.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &loc, nil
}

func getAllLocationsLocal(db *sql.DB) ([]LocationResponse, error) {
	query := `SELECT id, identifier_value, name, description, phone FROM locations`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []LocationResponse
	for rows.Next() {
		var loc LocationResponse
		if err := rows.Scan(&loc.ID, &loc.IdentifierValue, &loc.Name, &loc.Description, &loc.Phone); err != nil {
			return nil, err
		}
		locs = append(locs, loc)
	}
	return locs, nil
}

func saveEncounterLocal(db *sql.DB, enc EncounterResponse) error {
	query := `INSERT INTO encounters (id, identifier_value, status, subject_id, location_id, start_time) 
	          VALUES (?, ?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET 
			  identifier_value=excluded.identifier_value, status=excluded.status, 
			  subject_id=excluded.subject_id, location_id=excluded.location_id, 
			  start_time=excluded.start_time`
	_, err := db.Exec(query, enc.ID, enc.IdentifierValue, enc.Status, enc.SubjectID, enc.LocationID, enc.StartTime)
	return err
}

func updateEncounterStatusLocal(db *sql.DB, id, status string) error {
	query := `UPDATE encounters SET status = ? WHERE id = ?`
	res, err := db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("encounter %s not found in local db", id)
	}
	return nil
}

func getEncounterLocal(db *sql.DB, id string) (*EncounterResponse, error) {
	var enc EncounterResponse
	query := `SELECT id, identifier_value, status, subject_id, location_id, start_time FROM encounters WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&enc.ID, &enc.IdentifierValue, &enc.Status, &enc.SubjectID, &enc.LocationID, &enc.StartTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &enc, nil
}

func getAllEncountersLocal(db *sql.DB) ([]EncounterResponse, error) {
	query := `SELECT id, identifier_value, status, subject_id, location_id, start_time FROM encounters`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var encs []EncounterResponse
	for rows.Next() {
		var enc EncounterResponse
		if err := rows.Scan(&enc.ID, &enc.IdentifierValue, &enc.Status, &enc.SubjectID, &enc.LocationID, &enc.StartTime); err != nil {
			return nil, err
		}
		encs = append(encs, enc)
	}
	return encs, nil
}

func searchLocationsLocal(db *sql.DB, id, identifier string, page, limit int) ([]LocationResponse, error) {
	offset := (page - 1) * limit
	query := `SELECT id, identifier_value, name, description, phone FROM locations`
	var args []interface{}
	var conditions []string

	if id != "" {
		conditions = append(conditions, "id = ?")
		args = append(args, id)
	}
	if identifier != "" {
		conditions = append(conditions, "identifier_value = ?")
		args = append(args, identifier)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []LocationResponse
	for rows.Next() {
		var loc LocationResponse
		if err := rows.Scan(&loc.ID, &loc.IdentifierValue, &loc.Name, &loc.Description, &loc.Phone); err != nil {
			return nil, err
		}
		locs = append(locs, loc)
	}
	return locs, nil
}

func savePractitionerLocal(db *sql.DB, p PersonResponse) error {
	query := `INSERT INTO practitioners (nik, id, ihs_number, name, gender, birth_date, phone, address) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(nik) DO UPDATE SET 
			  id=excluded.id, ihs_number=excluded.ihs_number, name=excluded.name, gender=excluded.gender, 
			  birth_date=excluded.birth_date, phone=excluded.phone, address=excluded.address`
	_, err := db.Exec(query, p.NIK, p.ID, p.IHSNumber, p.Name, p.Gender, p.BirthDate, p.Phone, p.Address)
	return err
}

func getPractitionerLocal(db *sql.DB, nik string) (*PersonResponse, error) {
	var p PersonResponse
	query := `SELECT nik, id, ihs_number, name, gender, birth_date, phone, address FROM practitioners WHERE nik = ?`
	err := db.QueryRow(query, nik).Scan(&p.NIK, &p.ID, &p.IHSNumber, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func getAllPractitionersLocal(db *sql.DB) ([]PersonResponse, error) {
	query := `SELECT nik, id, ihs_number, name, gender, birth_date, phone, address FROM practitioners`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var practitioners []PersonResponse
	for rows.Next() {
		var p PersonResponse
		if err := rows.Scan(&p.NIK, &p.ID, &p.IHSNumber, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.Address); err != nil {
			return nil, err
		}
		practitioners = append(practitioners, p)
	}
	return practitioners, nil
}

func searchPractitionersLocal(db *sql.DB, id, nik, name, gender, birthDate string, page, limit int) ([]PersonResponse, error) {
	offset := (page - 1) * limit
	query := `SELECT nik, id, ihs_number, name, gender, birth_date, phone, address FROM practitioners`
	var args []interface{}
	var conditions []string

	if id != "" {
		conditions = append(conditions, "id = ?")
		args = append(args, id)
	}
	if nik != "" {
		conditions = append(conditions, "nik = ?")
		args = append(args, nik)
	}
	if name != "" {
		cleanName := name
		if strings.HasPrefix(strings.ToLower(name), "dr. ") {
			cleanName = name[4:]
		} else if strings.HasPrefix(strings.ToLower(name), "dr ") {
			cleanName = name[3:]
		}
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+cleanName+"%")
	}
	if gender != "" {
		conditions = append(conditions, "gender = ?")
		args = append(args, gender)
	}
	if birthDate != "" {
		conditions = append(conditions, "birth_date = ?")
		args = append(args, birthDate)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pracs []PersonResponse
	for rows.Next() {
		var p PersonResponse
		if err := rows.Scan(&p.NIK, &p.ID, &p.IHSNumber, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.Address); err != nil {
			return nil, err
		}
		pracs = append(pracs, p)
	}
	return pracs, nil
}
