package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
)

func TestFindDnsRecord_nonExistantRecordReturnsError(t *testing.T) {

	c := client.NewClient(client.GetConfigFromEnv())

	recordName := "record.does.not.exist"
	_, err := c.FindDnsRecord(recordName)

	expectedErrStr := fmt.Sprintf("dns record `%s` not found", recordName)
	if err == nil || err.Error() != expectedErrStr {
		t.Errorf("client should have received error indicating the following dns record `%s` was not found. Instead error was nil", recordName)
	}
}
