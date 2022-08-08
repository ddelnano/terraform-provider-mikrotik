package client

import "testing"

func TestAddInterfaceListMemberUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	listName := "test_list"

	list, err := c.AddInterfaceList(&InterfaceList{
		Name: listName,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.DeleteInterfaceList(list.Id); err != nil {
			t.Error(err)
		}
	}()

	listMember, err := c.AddInterfaceListMember(&InterfaceListMember{
		List:      list.Name,
		Interface: "*0",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.DeleteInterfaceListMember(listMember.Id); err != nil {
			t.Error(err)
		}
		if m, err := c.FindInterfaceListMember(listMember.Id); err == nil || m != nil {
			t.Errorf("expected error to be present and list member to be nil")
		}
	}()

	found, err := c.FindInterfaceListMember(listMember.Id)
	if err != nil {
		t.Fatal(err)
	}

	if found.List != list.Name {
		t.Errorf("expected name to be %q, got %q", list.Name, found.List)
	}

	listMember.Interface = "ether1"
	updated, err := c.UpdateInterfaceListMember(listMember)
	if err != nil {
		t.Error(err)
	}

	if updated.Interface != "ether1" {
		t.Errorf("expected updated interface to be %q, got %q", listMember.Interface, updated.Interface)
	}
}
