package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var name string = "testacc"
var ranges string = "172.16.0.1-172.16.0.8,172.16.0.10"
var comment string = "terraform-acc-test-pool"
var updatedRanges string = "172.16.0.1-172.16.0.8,172.16.0.16"
var updatedComment string = "terraform acc test pool updated"

func TestAddUpdateAndDeletePool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedPool := &Pool{
		Name:    name,
		Ranges:  ranges,
		Comment: comment,
	}
	pool, err := c.AddPool(expectedPool)

	if err != nil {
		t.Fatalf("Error creating a pool with: %v", err)
	}

	expectedPool.Id = pool.Id
	if !reflect.DeepEqual(pool, expectedPool) {
		t.Errorf("The pool does not match what we expected. actual: %v expected: %v", pool, expectedPool)
	}

	expectedPool.Comment = updatedComment
	expectedPool.Ranges = updatedRanges
	pool, err = c.UpdatePool(expectedPool)

	if err != nil {
		t.Errorf("Error updating pool with: %v", err)
	}

	if !reflect.DeepEqual(pool, expectedPool) {
		t.Errorf("Updated pool does not match the expected: %v expected: %v", expectedPool, pool)
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

	require.Truef(t, IsNotFoundError(err), "client should have NotFound error error but instead received")
}

func TestFindPoolByName_forExistingPool(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	p := &Pool{
		Name:    name,
		Ranges:  ranges,
		Comment: comment,
	}
	pool, err := c.AddPool(p)

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

	require.True(t, IsNotFoundError(err),
		"client should have NotFound error")
}
