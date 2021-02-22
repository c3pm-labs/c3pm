// Package registry handles interaction with the C3PM registry.
// It handles file downloading and version querying.
package registry

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/c3pm-labs/c3pm/env"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

//Options holds the options to pass to every function interacting with the registry
type Options struct {
	//RegistryURL is the URL to call to reach the registry.
	RegistryURL string
}

type Client struct {
	s3Client *s3.Client
}

func NewClient(opts Options) *Client {
	s3Client := s3.New(s3.Options{
		UsePathStyle:     true,
		Region: "fr-par",
		EndpointResolver: s3.EndpointResolverFromURL(opts.RegistryURL),
	})

	return &Client{s3Client: s3Client}
}

//GetLastVersion calls the registry to find the latest version published to the API.
//The version found can be different to the version that has been published to the API in case of support of ancient versions.
//For example, if a package is currently at version 3.3.0, but the maintainer last pushed version 2.7.3, a patch for version 2.7.
//The version returned by GetLastVersion will be 3.3.0, because it is the highest SemVer version number.
func (c *Client) GetLastVersion(ctx context.Context, pkgName string) (*semver.Version, error) {
	resp, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(env.REGISTRY_BUCKET_NAME),
		Prefix: aws.String(pkgName),
	})
	if err != nil {
		fmt.Println("HERE")
		return nil, err
	}
	vs := make([]*semver.Version, len(resp.Contents))
	for i, r := range resp.Contents {
		version := filepath.Base(*r.Key)
		v, err := semver.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("error parsing version %s: %w", r, err)
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))
	return vs[len(vs)-1], nil
}

//FetchPackage downloads a package given it's name and version number.
func (c *Client) FetchPackage(ctx context.Context, pkgName string, version *semver.Version) (*os.File, error) {
	key := fmt.Sprintf("/%s/%s", pkgName, version.String())

	resp, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(env.REGISTRY_BUCKET_NAME),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	file, err := ioutil.TempFile("", fmt.Sprintf("%s.%s.*.tar", pkgName, version.String()))
	if err != nil {
		return nil, fmt.Errorf("error creating temporary package file: %w", err)
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error downloading package: %w", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error returning to file beginning: %w", err)
	}
	return file, nil
}
