package dsql

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientTransportError(t *testing.T) {
	// some tcp error
}

func TestClientBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `
    {
      "__type":"com.amazonaws.dynamodb.v20111205#ProvisionedThroughputExceededException",
      "message":"The level of configured provisioned throughput for the table was exceeded. Consider increasing your provisioning level with the UpdateTable API"
    }`)
	}))
	defer ts.Close()

	client := NewDynamoClient("dyanmodb://access:secret@us-east-1")
	client.urlStr = ts.URL
	_, actual := client.Post(&Query{})

	expected := DynamoError{
		Type:    "com.amazonaws.dynamodb.v20111205#ProvisionedThroughputExceededException",
		Message: "The level of configured provisioned throughput for the table was exceeded. Consider increasing your provisioning level with the UpdateTable API",
	}

	if expected.Error() != actual.Error() {
		t.Error("expected ", expected)
		t.Error("actual   ", actual)
	}
}

func TestClientOK(t *testing.T) {

}
