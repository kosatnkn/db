package postgres_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/kosatnkn/db"
	"github.com/kosatnkn/db/postgres"
)

// NOTE: you will have to create a db named sample and in that a schema named sample and add the following table to it.
//
// | sample 					|
// | -------------------------- |
// | id (int, autoincrement)	|
// | name (varchar)				|
// | password (varchar) 		|
//

// newDBAdapter creates a new db adapter pointing to the test db.
func newDBAdapter(t *testing.T) db.AdapterInterface {

	cfg := postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "sample",
		User:     "postgres",
		Password: "admin",
		PoolSize: 10,
		Check:    true,
	}

	a, err := postgres.NewAdapter(cfg)
	if err != nil {
		t.Fatalf("Cannot create adapter. Error: %v", err)
	}

	return a
}

// clearTestTable clears all data from the test table.
func clearTestTable(t *testing.T) {

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	adapter.Query(context.Background(), `truncate sample.sample`, nil)
	adapter.Query(context.Background(), `alter sequence sample_id_seq restart with 1`, nil)

	t.Log("Table truncated")
}

// TestSelect tests select query.
func TestSelect(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	q := "select * from sample.sample"

	r, err := adapter.Query(context.Background(), q, nil)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	need := reflect.TypeOf(make([]map[string]interface{}, 0))
	got := reflect.TypeOf(r)

	if got != need {
		t.Errorf("Need %d, got %d", need, got)
	}
}

// TestInsert tests insert query.
func TestInsert(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	q := `insert into sample.sample(name, password) values (?name, ?password) returning id`
	params := map[string]interface{}{
		"name":     "Success Data 1",
		"password": "pwd1",
	}

	r, err := adapter.Query(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 1
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 1
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is inserted
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	cNeed := "1, Success Data 1, pwd1"
	cGot := fmt.Sprintf("%d, %s, %s", int(cr[0]["id"].(int64)), cr[0]["name"], cr[0]["password"])
	if cGot != cNeed {
		t.Errorf("Need `%s`, got `%s`", cNeed, cGot)
	}
}

// TestUpdate tests update query.
func TestUpdate(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	// insert
	q := `insert into sample(name, password) values (?name, ?password)`
	params := map[string]interface{}{
		"name":     "Success Data 1",
		"password": "pwd1",
	}

	_, err := adapter.Query(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error inserting: %v", err)
	}

	//update
	q = `update sample set name = ?name, password = ?password where id = ?id`
	params = map[string]interface{}{
		"id":       1,
		"name":     "Success Data 2",
		"password": "pwd2",
	}

	r, err := adapter.Query(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error updating: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 1
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 0
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is updated
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	cNeed := "1, Success Data 2, pwd2"
	cGot := fmt.Sprintf("%d, %s, %s", int(cr[0]["id"].(int64)), cr[0]["name"], cr[0]["password"])
	if cGot != cNeed {
		t.Errorf("Need `%s`, got `%s`", cNeed, cGot)
	}
}

// TestDelete tests delete query.
func TestDelete(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	// insert
	q := `insert into sample(name, password) values (?name, ?password)`
	params := map[string]interface{}{
		"name":     "Success Data 1",
		"password": "pwd1",
	}

	_, err := adapter.Query(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error inserting: %v", err)
	}

	// delete
	q = `delete from sample where id = ?id`
	params = map[string]interface{}{
		"id": 1,
	}

	r, err := adapter.Query(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error deleting: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 1
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 0
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is inserted
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) > 0 {
		t.Errorf("Need 0 record, got %d records", len(r))
	}
}

// TestSelectBulk tests bulk select query.
func TestSelectBulk(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	q := "select * from sample"

	_, err := adapter.QueryBulk(context.Background(), q, nil)
	if err == nil {
		t.Errorf("Need error, got nil")
	}

	need := "Select queries are not allowed. Use Query() instead"
	got := err.Error()
	if got != need {
		t.Errorf("Need %s, got %s", need, got)
	}
}

// TestInsertBulk tests bulk insert query.
func TestInsertBulk(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	q := `insert into sample(name, password) values (?name, ?password)`

	params := make([]map[string]interface{}, 0)
	params = append(params, map[string]interface{}{
		"name":     "Name 1",
		"password": "pwd1",
	})
	params = append(params, map[string]interface{}{
		"name":     "Name 2",
		"password": "pwd2",
	})

	r, err := adapter.QueryBulk(context.Background(), q, params)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 2
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 2
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is inserted
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	cNeed := "1, Name 1, pwd1"
	cGot := fmt.Sprintf("%d, %s, %s", int(cr[0]["id"].(int64)), cr[0]["name"], cr[0]["password"])
	if cGot != cNeed {
		t.Errorf("Record 1: need `%s`, got `%s`", cNeed, cGot)
	}

	cNeed = "2, Name 2, pwd2"
	cGot = fmt.Sprintf("%d, %s, %s", int(cr[1]["id"].(int64)), cr[1]["name"], cr[1]["password"])
	if cGot != cNeed {
		t.Errorf("Record 2: need `%s`, got `%s`", cNeed, cGot)
	}
}

// TestUpdateBulk tests bulk update query.
func TestUpdateBulk(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	// insert
	q := `insert into sample(name, password) values (?name, ?password)`

	ips := make([]map[string]interface{}, 0)
	ips = append(ips, map[string]interface{}{
		"name":     "Name 1",
		"password": "pwd1",
	})
	ips = append(ips, map[string]interface{}{
		"name":     "Name 2",
		"password": "pwd2",
	})

	_, err := adapter.QueryBulk(context.Background(), q, ips)
	if err != nil {
		t.Fatalf("Error inserting: %v", err)
	}

	// update
	q = `update sample set name = ?name, password = ?password where id = ?id`

	ups := make([]map[string]interface{}, 0)
	ups = append(ups, map[string]interface{}{
		"id":       1,
		"name":     "Name 1 Updated",
		"password": "pwd1 Updated",
	})
	ups = append(ups, map[string]interface{}{
		"id":       2,
		"name":     "Name 2 Updated",
		"password": "pwd2 Updated",
	})

	r, err := adapter.QueryBulk(context.Background(), q, ups)
	if err != nil {
		t.Fatalf("Error updating: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 2
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 0
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is updated
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) != 2 {
		t.Errorf("Need 2 records, got %d records", len(r))
	}

	cNeed := "1, Name 1 Updated, pwd1 Updated"
	cGot := fmt.Sprintf("%d, %s, %s", int(cr[0]["id"].(int64)), cr[0]["name"], cr[0]["password"])
	if cGot != cNeed {
		t.Errorf("Record 1: need `%s`, got `%s`", cNeed, cGot)
	}

	cNeed = "2, Name 2 Updated, pwd2 Updated"
	cGot = fmt.Sprintf("%d, %s, %s", int(cr[1]["id"].(int64)), cr[1]["name"], cr[1]["password"])
	if cGot != cNeed {
		t.Errorf("Record 2: need `%s`, got `%s`", cNeed, cGot)
	}
}

// TestDeleteBulk tests bulk delete query.
func TestDeleteBulk(t *testing.T) {

	clearTestTable(t)

	adapter := newDBAdapter(t)
	defer adapter.Destruct()

	// insert
	q := `insert into sample(name, password) values (?name, ?password)`

	ips := make([]map[string]interface{}, 0)
	ips = append(ips, map[string]interface{}{
		"name":     "Name 1",
		"password": "pwd1",
	})
	ips = append(ips, map[string]interface{}{
		"name":     "Name 2",
		"password": "pwd2",
	})

	_, err := adapter.QueryBulk(context.Background(), q, ips)
	if err != nil {
		t.Fatalf("Error inserting: %v", err)
	}

	// delete
	q = `delete from sample where id = ?id`

	dps := make([]map[string]interface{}, 0)
	dps = append(dps, map[string]interface{}{
		"id": 1,
	})
	dps = append(dps, map[string]interface{}{
		"id": 2,
	})

	r, err := adapter.QueryBulk(context.Background(), q, dps)
	if err != nil {
		t.Fatalf("Error deleting: %v", err)
	}
	if len(r) == 0 {
		t.Errorf("Need 1 record, got %d records", len(r))
	}

	need := 2
	got := int(r[0]["affected_rows"].(int64))
	if got != need {
		t.Errorf("Affected rows: need `%d`, got `%d`", need, got)
	}

	need = 0
	got = int(r[0]["last_insert_id"].(int64))
	if got != need {
		t.Errorf("Last insert id: need `%d`, got `%d`", need, got)
	}

	// check whether all data is updated
	cr, _ := adapter.Query(context.Background(), `select * from sample`, nil)
	if len(cr) != 0 {
		t.Errorf("Need 0 records, got %d records", len(r))
	}
}