package sources

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	lxd "github.com/lxc/lxd/shared"

	"github.com/lxc/distrobuilder/shared"
)

// AlpineLinuxHTTP represents the Alpine Linux downloader.
type AlpineLinuxHTTP struct{}

// NewAlpineLinuxHTTP creates a new AlpineLinuxHTTP instance.
func NewAlpineLinuxHTTP() *AlpineLinuxHTTP {
	return &AlpineLinuxHTTP{}
}

// Run downloads an Alpine Linux mini root filesystem.
func (s *AlpineLinuxHTTP) Run(definition shared.Definition, rootfsDir string) error {
	fname := fmt.Sprintf("alpine-minirootfs-%s-%s.tar.gz", definition.Image.Release,
		definition.Image.ArchitectureMapped)
	tarball := fmt.Sprintf("%s/v%s/releases/%s/%s", definition.Source.URL,
		strings.Join(strings.Split(definition.Image.Release, ".")[0:2], "."),
		definition.Image.ArchitectureMapped, fname)

	url, err := url.Parse(tarball)
	if err != nil {
		return err
	}

	if url.Scheme != "https" && len(definition.Source.Keys) == 0 {
		return errors.New("GPG keys are required if downloading from HTTP")
	}

	err = shared.DownloadSha256(tarball, tarball+".sha256")
	if err != nil {
		return err
	}

	// Force gpg checks when using http
	if url.Scheme != "https" {
		shared.DownloadSha256(tarball+".asc", "")
		valid, err := shared.VerifyFile(
			filepath.Join(os.TempDir(), fname),
			filepath.Join(os.TempDir(), fname+".asc"),
			definition.Source.Keys,
			definition.Source.Keyserver)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Failed to verify tarball")
		}
	}

	// Unpack
	err = lxd.Unpack(filepath.Join(os.TempDir(), fname), rootfsDir, false, false)
	if err != nil {
		return err
	}

	return nil
}
