package terraform

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	v := `{
  "format_version": "1.0",
  "terraform_version": "1.4.2",
  "values": {
    "outputs": {
      "a": { 
        "sensitive": false,
        "value": "a", "type": "string"
      },
      "b": {
        "sensitive": true,
        "value": "b",
        "type": "string"
      }
    },
    "root_module": {
    }
  }
}`
	s := Secret{}
	json.Unmarshal([]byte(v), &s)
}
