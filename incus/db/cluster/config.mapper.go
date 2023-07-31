//go:build linux && cgo && !agent

package cluster

// The code below was generated by incus-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cyphar/incus/incus/db/query"
	"github.com/cyphar/incus/shared/api"
)

var _ = api.ServerEnvironment{}

const configObjects = `SELECT %s_config.id, %s_config.%s_id, %s_config.key, %s_config.value
  FROM %s_config
  ORDER BY %s_config.id`

const configCreate = `INSERT INTO %s_config (%s_id, key, value)
  VALUES (?, ?, ?)`

const configDelete = `DELETE FROM %s_config WHERE %s_id = ?`

// configColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Config entity.
func configColumns() string {
	return "%s_config.id, %s_config.%s_id, %s_config.key, %s_config.value"
}

// getConfig can be used to run handwritten sql.Stmts to return a slice of objects.
func getConfig(ctx context.Context, stmt *sql.Stmt, parent string, args ...any) ([]Config, error) {
	objects := make([]Config, 0)

	dest := func(scan func(dest ...any) error) error {
		c := Config{}
		err := scan(&c.ID, &c.ReferenceID, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"%s_config\" table: %w", parent, err)
	}

	return objects, nil
}

// getConfigRaw can be used to run handwritten query strings to return a slice of objects.
func getConfigRaw(ctx context.Context, tx *sql.Tx, sql string, parent string, args ...any) ([]Config, error) {
	objects := make([]Config, 0)

	dest := func(scan func(dest ...any) error) error {
		c := Config{}
		err := scan(&c.ID, &c.ReferenceID, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"%s_config\" table: %w", parent, err)
	}

	return objects, nil
}

// GetConfig returns all available config.
// generator: config GetMany
func GetConfig(ctx context.Context, tx *sql.Tx, parent string, filters ...ConfigFilter) (map[int]map[string]string, error) {
	var err error

	// Result slice.
	objects := make([]Config, 0)

	configObjectsLocal := strings.Replace(configObjects, "%s_id", fmt.Sprintf("%s_id", parent), -1)
	fillParent := make([]any, strings.Count(configObjectsLocal, "%s"))
	for i := range fillParent {
		fillParent[i] = strings.Replace(parent, "_", "s_", -1) + "s"
	}

	queryStr := fmt.Sprintf(configObjectsLocal, fillParent...)
	queryParts := strings.SplitN(queryStr, "ORDER BY", 2)
	args := []any{}

	for i, filter := range filters {
		var cond string
		if i == 0 {
			cond = " WHERE ( %s )"
		} else {
			cond = " OR ( %s )"
		}

		entries := []string{}
		if filter.Key != nil {
			entries = append(entries, "key = ?")
			args = append(args, filter.Key)
		}

		if filter.Value != nil {
			entries = append(entries, "value = ?")
			args = append(args, filter.Value)
		}

		if len(entries) == 0 {
			return nil, fmt.Errorf("Cannot filter on empty ConfigFilter")
		}

		queryParts[0] += fmt.Sprintf(cond, strings.Join(entries, " AND "))
	}

	queryStr = strings.Join(queryParts, " ORDER BY")
	// Select.
	objects, err = getConfigRaw(ctx, tx, queryStr, parent, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"%s_config\" table: %w", parent, err)
	}

	resultMap := map[int]map[string]string{}
	for _, object := range objects {
		_, ok := resultMap[object.ReferenceID]
		if !ok {
			resultMap[object.ReferenceID] = map[string]string{}
		}

		resultMap[object.ReferenceID][object.Key] = object.Value
	}

	return resultMap, nil
}

// CreateConfig adds a new config to the database.
// generator: config Create
func CreateConfig(ctx context.Context, tx *sql.Tx, parent string, object Config) error {
	// An empty value means we are unsetting this key, so just return.
	if object.Value == "" {
		return nil
	}

	configCreateLocal := strings.Replace(configCreate, "%s_id", fmt.Sprintf("%s_id", parent), -1)
	fillParent := make([]any, strings.Count(configCreateLocal, "%s"))
	for i := range fillParent {
		fillParent[i] = strings.Replace(parent, "_", "s_", -1) + "s"
	}

	queryStr := fmt.Sprintf(configCreateLocal, fillParent...)
	_, err := tx.ExecContext(ctx, queryStr, object.ReferenceID, object.Key, object.Value)
	if err != nil {
		return fmt.Errorf("Insert failed for \"%s_config\" table: %w", parent, err)
	}

	return nil
}

// UpdateConfig updates the config matching the given key parameters.
// generator: config Update
func UpdateConfig(ctx context.Context, tx *sql.Tx, parent string, referenceID int, config map[string]string) error {
	// Delete current entry.
	err := DeleteConfig(ctx, tx, parent, referenceID)
	if err != nil {
		return err
	}

	// Insert new entries.
	for key, value := range config {
		object := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err = CreateConfig(ctx, tx, parent, object)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteConfig deletes the config matching the given key parameters.
// generator: config DeleteMany
func DeleteConfig(ctx context.Context, tx *sql.Tx, parent string, referenceID int) error {
	configDeleteLocal := strings.Replace(configDelete, "%s_id", fmt.Sprintf("%s_id", parent), -1)
	fillParent := make([]any, strings.Count(configDeleteLocal, "%s"))
	for i := range fillParent {
		fillParent[i] = strings.Replace(parent, "_", "s_", -1) + "s"
	}

	queryStr := fmt.Sprintf(configDeleteLocal, fillParent...)
	result, err := tx.ExecContext(ctx, queryStr, referenceID)
	if err != nil {
		return fmt.Errorf("Delete entry for \"%s_config\" failed: %w", parent, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}
