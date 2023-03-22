package imgutil

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var NormalizedDateTime = time.Date(1980, time.January, 1, 0, 0, 1, 0, time.UTC)

type SaveDiagnostic struct {
	ImageName string
	Cause     error
}

type SaveError struct {
	Errors []SaveDiagnostic
}

func (e SaveError) Error() string {
	var errors []string
	for _, d := range e.Errors {
		errors = append(errors, fmt.Sprintf("[%s: %s]", d.ImageName, d.Cause.Error()))
	}
	return fmt.Sprintf("failed to write image to the following tags: %s", strings.Join(errors, ","))
}

// Platform represents the target arch/os/os_version for an image construction and querying.
type Platform struct {
	Architecture string
	OS           string
	OSVersion    string
}

type Image interface {
	// getters

	Architecture() (string, error)
	CreatedAt() (time.Time, error)
	Entrypoint() ([]string, error)
	Env(key string) (string, error)
	// Found tells whether the image exists in the repository by `Name()`.
	Found() bool
	GetAnnotateRefName() (string, error)
	// GetLayer retrieves layer by diff id. Returns a reader of the uncompressed contents of the layer.
	GetLayer(diffID string) (io.ReadCloser, error)
	Identifier() (Identifier, error)
	Label(string) (string, error)
	Labels() (map[string]string, error)
	// ManifestSize returns the size of the manifest. If a manifest doesn't exist, it returns 0.
	ManifestSize() (int64, error)
	Name() string
	OS() (string, error)
	OSVersion() (string, error)
	// TopLayer returns the diff id for the top layer
	TopLayer() (string, error)
	Variant() (string, error)
	WorkingDir() (string, error)

	// setters

	// AnnotateRefName set a value for the `org.opencontainers.image.ref.name` annotation
	AnnotateRefName(refName string) error
	Rename(name string)
	SetArchitecture(string) error
	SetCmd(...string) error
	SetEntrypoint(...string) error
	SetEnv(string, string) error
	SetLabel(string, string) error
	SetOS(string) error
	SetOSVersion(string) error
	SetVariant(string) error
	SetWorkingDir(string) error

	// modifiers

	AddLayer(path string) error
	AddLayerWithDiffID(path, diffID string) error
	Delete() error
	Rebase(string, Image) error
	RemoveLabel(string) error
	ReuseLayer(diffID string) error
	// Save saves the image as `Name()` and any additional names provided to this method.
	Save(additionalNames ...string) error
	// SaveAs ignores the image `Name()` method and saves the image according to name & additional names provided to this method
	SaveAs(name string, additionalNames ...string) error
	// SaveFile saves the image as a docker archive and provides the filesystem location
	SaveFile() (string, error)
}

type Identifier fmt.Stringer
