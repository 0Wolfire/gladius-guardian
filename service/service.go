package service

import "github.com/gladiusio/gladius-common/pkg/manager"

func SetupService(run func()) {
	// Define some variables
	name, displayName, description :=
		"GladiusGuardian",
		"Gladius Guardian",
		"Gladius Guardian"

	// Run the function "run" in guardian as a service
	manager.RunService(name, displayName, description, run)
}
