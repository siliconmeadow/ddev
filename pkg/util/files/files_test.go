package files

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/drud/ddev/pkg/testcommon"
	"github.com/drud/drud-go/utils/system"
	"github.com/stretchr/testify/assert"
)

var (
	temp            = os.TempDir()
	cwd             string
	testArchiveURL  = "https://github.com/drud/wordpress/releases/download/v0.1.0/files.tar.gz"
	testArchivePath = path.Join(testcommon.CreateTmpDir("filetest"), "files.tar.gz")
)

func TestMain(m *testing.M) {
	err := system.DownloadFile(testArchivePath, testArchiveURL)
	if err != nil {
		log.Fatalf("archive download failed: %s", err)
	}

	cwd, err = os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}

	testRun := m.Run()

	os.Exit(testRun)
}

func TestUntar(t *testing.T) {
	assert := assert.New(t)
	exDir := path.Join(temp, "extract")

	err := Untar(testArchivePath, exDir)
	assert.NoError(err)

	err = os.RemoveAll(exDir)
	assert.NoError(err)
}

// TestCopyFile tests copying a file.
func TestCopyFile(t *testing.T) {
	assert := assert.New(t)
	dest := path.Join(temp, "testfile2")

	err := os.Chmod(testArchivePath, 0644)
	assert.NoError(err)

	err = CopyFile(testArchivePath, dest)
	assert.NoError(err)

	file, err := os.Stat(dest)
	assert.NoError(err)
	assert.Equal(int(file.Mode()), 0644)

	err = os.RemoveAll(dest)
	assert.NoError(err)
}

// TestCopyDir tests copying a directory.
func TestCopyDir(t *testing.T) {
	assert := assert.New(t)
	dest := path.Join(temp, "copy")
	err := os.Mkdir(dest, 0755)
	assert.NoError(err)

	// test source not a directory
	err = CopyDir(testArchivePath, temp)
	assert.Error(err)
	assert.Contains(err.Error(), "source is not a directory")

	// test destination exists
	err = CopyDir(temp, cwd)
	assert.Error(err)
	assert.Contains(err.Error(), "destination already exists")
	err = os.RemoveAll(dest)
	assert.NoError(err)

	// copy a directory.
	err = CopyDir(cwd, dest)
	assert.NoError(err)
	assert.True(system.FileExists(path.Join(dest, "files.go")))
	assert.True(system.FileExists(path.Join(dest, "files_test.go")))

	err = os.RemoveAll(dest)
	assert.NoError(err)
}