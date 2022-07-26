package client

import "testing"

func TestAddInterfaceListUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	list, err := c.AddInterfaceList(&InterfaceList{
		Name:    "mylist",
		Comment: "Created by terraform",
	})
	if err != nil {
		t.Fatal(err)
	}

	found, err := c.FindInterfaceList(list.Name)
	if err != nil {
		t.Fatal(err)
	}

	if found.Name != list.Name {
		t.Errorf("expected name to be %q, got %q", list.Name, found.Name)
	}

	list.Comment = "updated list"
	updated, err := c.UpdateInterfaceList(list)
	if err != nil {
		t.Error(err)
	}

	if updated.Comment != "updated list" {
		t.Errorf("expected comment to be %q, got %q", list.Comment, updated.Comment)
	}

	// cleanup
	if err := c.DeleteInterfaceList(list.Name); err != nil {
		t.Error(err)
	}

	_, err = c.FindInterfaceList(list.Name)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
