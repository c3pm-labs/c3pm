//+build !dev

package env

const (
	//API_ENDPOINT is the URL of the C3PM API
	API_ENDPOINT = "https://c3pm.herokuapp.com/v1"
	//REGISTRY_ENDPOINT is the URL to an S3-compatible bucket that will be used as the package registry.
	REGISTRY_ENDPOINT = "https://registry-c3pm-io.s3.fr-par.scw.cloud"
)
