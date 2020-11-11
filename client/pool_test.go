package client

import (
	"reflect"
	"testing"
)

var name string = "testacc"
var ranges string = "172.16.0.1-172.16.0.8,172.16.0.10"
var comment string = "terraform-acc-test-pool"

func TestAddPoolAndDeletePool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedPool := &Pool{
		Name:    name,
		Ranges:  ranges,
		Comment: comment,
	}
	pool, err := c.AddPool(
		name,
		ranges,
		comment,
	)

	if err != nil {
		t.Fatalf("Error creating a pool with: %v", err)
	}

	if len(pool.Id) < 1 {
		t.Errorf("The created pool does not have an Id: %v", pool)
	}

	if pool.Name != expectedPool.Name {
		t.Errorf("The pool Name fields do not match. actual: %v expected: %v", pool.Name, expectedPool.Name)
	}

	if pool.Ranges != expectedPool.Ranges {
		t.Errorf("The pool Ranges fields do not match. actual: %v expected: %v", pool.Ranges, expectedPool.Ranges)
	}

	if pool.Comment != expectedPool.Comment {
		t.Errorf("The pool Comment fields do not match. actual: %v expected: %v", pool.Comment, expectedPool.Comment)
	}

	foundPool, err := c.FindPool(pool.Id)

	if err != nil {
		t.Errorf("Error getting pool with: %v", err)
	}

	if !reflect.DeepEqual(pool, foundPool) {
		t.Errorf("Created pool and found pool do not match. actual: %v expected: %v", foundPool, pool)
	}

	err = c.DeletePool(pool.Id)

	if err != nil {
		t.Errorf("Error deleting pool with: %v", err)
	}
}

func TestFindPool_forNonExistingPool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	poolId := "Invalid id"
	_, err := c.FindPool(poolId)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("client should have NotFound error error but instead received '%v'", err)
	}
}

func TestFindPoolByName_forExistingPool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())
	pool, err := c.AddPool(
		name,
		ranges,
		comment,
	)

	expectedPool, err := c.FindPoolByName(pool.Name)
	if err != nil {
		t.Fatalf("Error finding pool by name with: %v", err)
	}
	if pool.Name != expectedPool.Name {
		t.Errorf("The pool Name fields do not match. actual: %v expected: %v", pool.Name, expectedPool.Name)
	}
	c.DeletePool(pool.Id)
}

func TestFindPoolByName_forNonExistingPool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	poolName := "Invalid name"
	_, err := c.FindPoolByName(poolName)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("client should have NotFound error error but instead received '%v'", err)
	}
}
