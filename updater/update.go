package updater

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gladiusio/gladius-common/pkg/requests"
)

// GetOfficialVersions - get version numbers from official droplet
func GetOfficialVersions() (map[string]string, error) {
	res, err := requests.SendRequest("GET", "https://gladius-version.nyc3.digitaloceanspaces.com/version.json", nil)
	if err != nil {
		return nil, err
	}

	var versions = make(map[string]string)
	err = json.Unmarshal([]byte(res), &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// GetVersion - get individual version number from module
func GetVersion(module string) (string, error) {
	var port int
	switch module {
	case "guardian":
		port = 7791
	case "edged":
		port = 7946
	case "network-gateway":
		port = 3001
	default:
		port = 0
	}

	if port == 0 {
		return "", fmt.Errorf("Module %s not found", module)
	}
	res, err := requests.SendRequest("GET", fmt.Sprintf("http://localhost:%d/version", port), nil)
	if err != nil {
		return "", err
	}

	var response = make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		return "", err
	}

	res1 := response["response"].(map[string]interface{})
	version := res1["version"].(string)

	return version, nil
}

// CompareVersion - are you on the right version?
// -1 = you are on a older version
//  0 = you are up-to-date
//  1 = you are on a newer version
func CompareVersion(myVersion, officialVersion string) (int, error) {
	mv := strings.Replace(myVersion, ".", "", -1)
	ov := strings.Replace(officialVersion, ".", "", -1)

	thisVersion, err := strconv.Atoi(mv)
	if err != nil {
		return -1, err
	}

	realVersion, err := strconv.Atoi(ov)
	if err != nil {
		return -1, err
	}

	if thisVersion < realVersion {
		return -1, nil
	} else if thisVersion > realVersion {
		return 1, nil
	} else {
		return 0, nil
	}
}

// VersionHandler - handler for the version
func VersionHandler(_res []byte) (map[string]string, error) {

	var response = make(map[string]string)

	err := json.Unmarshal(_res, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
