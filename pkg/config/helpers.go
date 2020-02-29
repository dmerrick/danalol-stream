package config

func IsProduction() bool {
	return Environment == "production"
}

func IsStaging() bool {
	return Environment == "staging"
}
