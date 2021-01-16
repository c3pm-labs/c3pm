package registry

import (
	"encoding/xml"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type Options struct {
	RegistryURL string
}

type ListRegistryResponse struct {
	Name     string
	Contents []struct {
		Key string `xml:"Key"`
	} `xml:"Contents"`
}

func GetLastVersion(dependency string, options Options) (*semver.Version, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", options.RegistryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating version query: %w", err)
	}
	q := req.URL.Query()
	q.Add("typeList", "2")
	q.Add("prefix", dependency)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching versions: %w", err)
	}
	defer resp.Body.Close()
	var registryResponse ListRegistryResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	err = xml.Unmarshal(body, &registryResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling versions: %w", err)
	}
	if len(registryResponse.Contents) == 0 {
		fmt.Printf("%s: package not found\n", dependency)
		os.Exit(1)
	}
	vs := make([]*semver.Version, len(registryResponse.Contents))
	for i, r := range registryResponse.Contents {
		version := filepath.Base(r.Key)
		v, err := semver.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("error parsing version %s: %w", r, err)
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))
	return vs[len(vs)-1], nil
}

func FetchPackage(dependency string, version *semver.Version, options Options) (*os.File, error) {
	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s/%s/%s", options.RegistryURL, dependency, version.String()))
	if err != nil {
		return nil, fmt.Errorf("error fetching package %s: %w", dependency, err)
	}
	defer resp.Body.Close()
	file, err := ioutil.TempFile("", fmt.Sprintf("%s.%s.*.tar", dependency, version.String()))
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
