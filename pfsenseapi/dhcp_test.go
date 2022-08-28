package pfsenseapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	dhcpLeasesTestResponse = `{"status":"ok","code":200,"return":0,"message":"Success",
"data":[{"ip":"192.168.60.7","type":"static","mac":"b4:5f:56:22:d4:33","if":"opt10",
"starts":"","ends":"","hostname":"host1","descr":"host1",
"online":true,"staticmap_array_index":1,"state":"static"}]}`
)

func TestDHCPService_ListLeases(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, dhcpLeasesTestResponse)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	newClient := NewClientWithNoAuth(server.URL)
	response, err := newClient.DHCP.ListLeases(context.Background())
	require.NoError(t, err)
	require.Len(t, response, 1)
}
