package types

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	// Azure specific
	TenantID     string
	ClientID     string
	ClientSecret string
	// GCP specific
	ProjectID   string
	KeyFilePath string
}
