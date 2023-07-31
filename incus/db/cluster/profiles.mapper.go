//go:build linux && cgo && !agent

package cluster

// The code below was generated by incus-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyphar/incus/incus/db/query"
	"github.com/cyphar/incus/shared/api"
)

var _ = api.ServerEnvironment{}

var profileObjects = RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  ORDER BY projects.id, profiles.name
`)

var profileObjectsByID = RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  WHERE ( profiles.id = ? )
  ORDER BY projects.id, profiles.name
`)

var profileObjectsByName = RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  WHERE ( profiles.name = ? )
  ORDER BY projects.id, profiles.name
`)

var profileObjectsByProject = RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  WHERE ( project = ? )
  ORDER BY projects.id, profiles.name
`)

var profileObjectsByProjectAndName = RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  WHERE ( project = ? AND profiles.name = ? )
  ORDER BY projects.id, profiles.name
`)

var profileID = RegisterStmt(`
SELECT profiles.id FROM profiles
  JOIN projects ON profiles.project_id = projects.id
  WHERE projects.name = ? AND profiles.name = ?
`)

var profileCreate = RegisterStmt(`
INSERT INTO profiles (project_id, name, description)
  VALUES ((SELECT projects.id FROM projects WHERE projects.name = ?), ?, ?)
`)

var profileRename = RegisterStmt(`
UPDATE profiles SET name = ? WHERE project_id = (SELECT projects.id FROM projects WHERE projects.name = ?) AND name = ?
`)

var profileUpdate = RegisterStmt(`
UPDATE profiles
  SET project_id = (SELECT projects.id FROM projects WHERE projects.name = ?), name = ?, description = ?
 WHERE id = ?
`)

var profileDeleteByProjectAndName = RegisterStmt(`
DELETE FROM profiles WHERE project_id = (SELECT projects.id FROM projects WHERE projects.name = ?) AND name = ?
`)

// GetProfileID return the ID of the profile with the given key.
// generator: profile ID
func GetProfileID(ctx context.Context, tx *sql.Tx, project string, name string) (int64, error) {
	stmt, err := Stmt(tx, profileID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"profileID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, project, name)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "Profile not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"profiles\" ID: %w", err)
	}

	return id, nil
}

// ProfileExists checks if a profile with the given key exists.
// generator: profile Exists
func ProfileExists(ctx context.Context, tx *sql.Tx, project string, name string) (bool, error) {
	_, err := GetProfileID(ctx, tx, project, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// profileColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Profile entity.
func profileColumns() string {
	return "profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')"
}

// getProfiles can be used to run handwritten sql.Stmts to return a slice of objects.
func getProfiles(ctx context.Context, stmt *sql.Stmt, args ...any) ([]Profile, error) {
	objects := make([]Profile, 0)

	dest := func(scan func(dest ...any) error) error {
		p := Profile{}
		err := scan(&p.ID, &p.ProjectID, &p.Project, &p.Name, &p.Description)
		if err != nil {
			return err
		}

		objects = append(objects, p)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	return objects, nil
}

// getProfilesRaw can be used to run handwritten query strings to return a slice of objects.
func getProfilesRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]Profile, error) {
	objects := make([]Profile, 0)

	dest := func(scan func(dest ...any) error) error {
		p := Profile{}
		err := scan(&p.ID, &p.ProjectID, &p.Project, &p.Name, &p.Description)
		if err != nil {
			return err
		}

		objects = append(objects, p)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	return objects, nil
}

// GetProfiles returns all available profiles.
// generator: profile GetMany
func GetProfiles(ctx context.Context, tx *sql.Tx, filters ...ProfileFilter) ([]Profile, error) {
	var err error

	// Result slice.
	objects := make([]Profile, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(tx, profileObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"profileObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Project != nil && filter.Name != nil && filter.ID == nil {
			args = append(args, []any{filter.Project, filter.Name}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, profileObjectsByProjectAndName)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"profileObjectsByProjectAndName\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(profileObjectsByProjectAndName)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"profileObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Project != nil && filter.ID == nil && filter.Name == nil {
			args = append(args, []any{filter.Project}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, profileObjectsByProject)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"profileObjectsByProject\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(profileObjectsByProject)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"profileObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Name != nil && filter.ID == nil && filter.Project == nil {
			args = append(args, []any{filter.Name}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, profileObjectsByName)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"profileObjectsByName\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(profileObjectsByName)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"profileObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID != nil && filter.Project == nil && filter.Name == nil {
			args = append(args, []any{filter.ID}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, profileObjectsByID)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"profileObjectsByID\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(profileObjectsByID)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"profileObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID == nil && filter.Project == nil && filter.Name == nil {
			return nil, fmt.Errorf("Cannot filter on empty ProfileFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getProfiles(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getProfilesRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	return objects, nil
}

// GetProfileDevices returns all available Profile Devices
// generator: profile GetMany
func GetProfileDevices(ctx context.Context, tx *sql.Tx, profileID int, filters ...DeviceFilter) (map[string]Device, error) {
	profileDevices, err := GetDevices(ctx, tx, "profile", filters...)
	if err != nil {
		return nil, err
	}

	devices := map[string]Device{}
	for _, ref := range profileDevices[profileID] {
		_, ok := devices[ref.Name]
		if !ok {
			devices[ref.Name] = ref
		} else {
			return nil, fmt.Errorf("Found duplicate Device with name %q", ref.Name)
		}
	}

	return devices, nil
}

// GetProfileConfig returns all available Profile Config
// generator: profile GetMany
func GetProfileConfig(ctx context.Context, tx *sql.Tx, profileID int, filters ...ConfigFilter) (map[string]string, error) {
	profileConfig, err := GetConfig(ctx, tx, "profile", filters...)
	if err != nil {
		return nil, err
	}

	config, ok := profileConfig[profileID]
	if !ok {
		config = map[string]string{}
	}

	return config, nil
}

// GetProfile returns the profile with the given key.
// generator: profile GetOne
func GetProfile(ctx context.Context, tx *sql.Tx, project string, name string) (*Profile, error) {
	filter := ProfileFilter{}
	filter.Project = &project
	filter.Name = &name

	objects, err := GetProfiles(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "Profile not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"profiles\" entry matches")
	}
}

// CreateProfile adds a new profile to the database.
// generator: profile Create
func CreateProfile(ctx context.Context, tx *sql.Tx, object Profile) (int64, error) {
	// Check if a profile with the same key exists.
	exists, err := ProfileExists(ctx, tx, object.Project, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"profiles\" entry already exists")
	}

	args := make([]any, 3)

	// Populate the statement arguments.
	args[0] = object.Project
	args[1] = object.Name
	args[2] = object.Description

	// Prepared statement to use.
	stmt, err := Stmt(tx, profileCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"profileCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"profiles\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"profiles\" entry ID: %w", err)
	}

	return id, nil
}

// CreateProfileDevices adds new profile Devices to the database.
// generator: profile Create
func CreateProfileDevices(ctx context.Context, tx *sql.Tx, profileID int64, devices map[string]Device) error {
	for key, device := range devices {
		device.ReferenceID = int(profileID)
		devices[key] = device
	}

	err := CreateDevices(ctx, tx, "profile", devices)
	if err != nil {
		return fmt.Errorf("Insert Device failed for Profile: %w", err)
	}

	return nil
}

// CreateProfileConfig adds new profile Config to the database.
// generator: profile Create
func CreateProfileConfig(ctx context.Context, tx *sql.Tx, profileID int64, config map[string]string) error {
	referenceID := int(profileID)
	for key, value := range config {
		insert := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err := CreateConfig(ctx, tx, "profile", insert)
		if err != nil {
			return fmt.Errorf("Insert Config failed for Profile: %w", err)
		}

	}

	return nil
}

// RenameProfile renames the profile matching the given key parameters.
// generator: profile Rename
func RenameProfile(ctx context.Context, tx *sql.Tx, project string, name string, to string) error {
	stmt, err := Stmt(tx, profileRename)
	if err != nil {
		return fmt.Errorf("Failed to get \"profileRename\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(to, project, name)
	if err != nil {
		return fmt.Errorf("Rename Profile failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}

	return nil
}

// UpdateProfile updates the profile matching the given key parameters.
// generator: profile Update
func UpdateProfile(ctx context.Context, tx *sql.Tx, project string, name string, object Profile) error {
	id, err := GetProfileID(ctx, tx, project, name)
	if err != nil {
		return err
	}

	stmt, err := Stmt(tx, profileUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"profileUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Project, object.Name, object.Description, id)
	if err != nil {
		return fmt.Errorf("Update \"profiles\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}

// UpdateProfileDevices updates the profile Device matching the given key parameters.
// generator: profile Update
func UpdateProfileDevices(ctx context.Context, tx *sql.Tx, profileID int64, devices map[string]Device) error {
	err := UpdateDevices(ctx, tx, "profile", int(profileID), devices)
	if err != nil {
		return fmt.Errorf("Replace Device for Profile failed: %w", err)
	}

	return nil
}

// UpdateProfileConfig updates the profile Config matching the given key parameters.
// generator: profile Update
func UpdateProfileConfig(ctx context.Context, tx *sql.Tx, profileID int64, config map[string]string) error {
	err := UpdateConfig(ctx, tx, "profile", int(profileID), config)
	if err != nil {
		return fmt.Errorf("Replace Config for Profile failed: %w", err)
	}

	return nil
}

// DeleteProfile deletes the profile matching the given key parameters.
// generator: profile DeleteOne-by-Project-and-Name
func DeleteProfile(ctx context.Context, tx *sql.Tx, project string, name string) error {
	stmt, err := Stmt(tx, profileDeleteByProjectAndName)
	if err != nil {
		return fmt.Errorf("Failed to get \"profileDeleteByProjectAndName\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(project, name)
	if err != nil {
		return fmt.Errorf("Delete \"profiles\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "Profile not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d Profile rows instead of 1", n)
	}

	return nil
}
