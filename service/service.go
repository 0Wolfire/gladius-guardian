package service

import "github.com/gladiusio/gladius-utils/init/manager"

func SetupService(run func()) {
	// Define some variables
	name, displayName, description :=
		"GladiusGuardian",
		"Gladius Guardian",
		"Gladius Guardian"

	// Run the function "run" in newtworkd as a service
	manager.RunService(name, displayName, description, run)
}
