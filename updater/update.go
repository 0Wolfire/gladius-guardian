package updater

import (
	"github.com/gladiusio/gladius-common/pkg/requests"
)

// GetVersion - get version number from droplet
func GetVersion() (string, error) {
	res, err := requests.SendRequest("GET", "https://gladius-node-hackathon.nyc3.digitaloceanspaces.com", nil)
	if err != nil {
		return "", err
	}
	return res, nil
}

func CompareVersion(myVersion, officialVersion string) {
}
