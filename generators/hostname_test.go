package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lxc/distrobuilder/image"
	"github.com/lxc/distrobuilder/shared"
)

func TestHostnameGeneratorCreateLXCData(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "distrobuilder-test")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator := Get("hostname")
	if generator == nil {
		t.Fatal("Expected hostname generator, got nil")
	}

	definition := shared.DefinitionImage{
		Distribution: "ubuntu",
		Release:      "artful",
	}

	image := image.NewLXCImage(cacheDir, definition, shared.DefinitionTargetLXC{})

	err := os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	createTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hostname"), "hostname")

	err = generator.CreateLXCData(cacheDir, "/etc/hostname", image)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	validateTestFile(t, filepath.Join(cacheDir, "tmp", "etc", "hostname"), "hostname")
	validateTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hostname"), "LXC_NAME\n")

	err = RestoreFiles(cacheDir)
	if err != nil {
		t.Fatalf("Failed to restore files: %s", err)
	}

	validateTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hostname"), "hostname")
}

func TestHostnameGeneratorCreateLXDData(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "distrobuilder-test")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator := Get("hostname")
	if generator == nil {
		t.Fatal("Expected hostname generator, got nil")
	}

	definition := shared.DefinitionImage{
		Distribution: "ubuntu",
		Release:      "artful",
	}

	image := image.NewLXDImage(cacheDir, definition)

	err := os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	err = generator.CreateLXDData(cacheDir, "/etc/hostname", image)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	validateTestFile(t, filepath.Join(cacheDir, "templates", "hostname.tpl"), "{{ container.name }}\n")
}