//+build !dev

package env

const (
	//API_ENDPOINT is the URL of the C3PM API
	API_ENDPOINT = "https://c3pm.herokuapp.com/v1"
	//REGISTRY_ENDPOINT is the URL to an S3-compatible bucket that will be used as the package registry.
	REGISTRY_ENDPOINT = "https://s3.fr-par.scw.cloud"
	//REGISTRY_BUCKET_NAME is the name of the s3 bucket
	REGISTRY_BUCKET_NAME = "registry-c3pm-io"
)
