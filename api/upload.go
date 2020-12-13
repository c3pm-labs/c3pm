package api

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func buildTarFile(files []string) bytes.Buffer {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, fn := range files {
		st, err := os.Lstat(fn)
		if err != nil {
			fmt.Printf("error handling %s: %s. ignoring...\n", fn, err.Error())
			continue
		}
		if st.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(fn)
			if err != nil {
				fmt.Printf("error handling %s: %s. ignoring...\n", fn, err.Error())
				continue
			}
			hdr := &tar.Header{
				Name:     fn,
				Linkname: target,
				Typeflag: tar.TypeSymlink,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				fmt.Printf("error writing %s: %s. ignoringa...", fn, err.Error())
				continue
			}
		} else {
			contents, err := ioutil.ReadFile(fn)
			if err != nil {
				fmt.Printf("error reading %s: %s. ignoringaa...\n", fn, err.Error())
				continue
			}
			hdr := &tar.Header{
				Name: fn,
				Mode: 0600,
				Size: st.Size(),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				fmt.Printf("error writing %s: %s. ignoringaaa...", fn, err.Error())
				continue
			}
			if _, err := tw.Write(contents); err != nil {
				fmt.Printf("error writing %s: %s. ignoringaaaa...", fn, err.Error())
				continue
			}
		}
	}
	tw.Close()
	return buf
}

// Upload is a wrapper function around C3PM's publish endpoint.
// It takes an array of file paths, creates a tar file containing
// them, then uploads it to the API.
func (c API) Upload(files []string) error {
	buf := buildTarFile(files)
	return c.send("POST", "/auth/publish", &buf)
}
