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
	mv := strings.Split(myVersion, ".")
	ov := strings.Split(officialVersion, ".")

	for i, num := range mv {
		currNum, err := strconv.Atoi(num)
		if err != nil {
			return -1, err
		}

		expeNum, err := strconv.Atoi(ov[i])
		if err != nil {
			return -1, err
		}

		if currNum < expeNum {
			return -1, nil
		} else if currNum > expeNum {
			return 1, nil
		}
	}
	return 0, nil
}
