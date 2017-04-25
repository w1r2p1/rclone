package webdav

import (
	"os"
	"os/exec"
	"testing"

	"github.com/ncw/rclone/fstest"
	_ "github.com/ncw/rclone/local"
	"github.com/stretchr/testify/assert"
)

// TestWebDav runs the webdav server then runs the unit tests for the
// webdav remote against it.
func TestWebDav(t *testing.T) {
	fstest.Initialise()

	fremote, _, clean, err := fstest.RandomRemote(*fstest.RemoteName, *fstest.SubDir)
	assert.NoError(t, err)
	defer clean()

	err = fremote.Mkdir("")
	assert.NoError(t, err)

	// Start the server
	go func() {
		err := serveWebDav(fremote)
		assert.NoError(t, err)
	}()
	// FIXME shut it down somehow?

	// Change directory to run the tests
	err = os.Chdir("../../../webdav")
	assert.NoError(t, err, "failed to cd to webdav remote")

	// Run the webdav tests with an on the fly remote
	args := []string{"test"}
	if testing.Verbose() {
		args = append(args, "-v")
	}
	if *fstest.Verbose {
		args = append(args, "-verbose")
	}
	args = append(args, "-remote", "webdavtest:")
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"RCLONE_CONFIG_WEBDAVTEST_TYPE=webdav",
		"RCLONE_CONFIG_WEBDAVTEST_URL=http://localhost:8081/",
		"RCLONE_CONFIG_WEBDAVTEST_VENDOR=other",
	)
	out, err := cmd.CombinedOutput()
	if len(out) != 0 {
		t.Logf("\n----------\n%s----------\n", string(out))
	}
	assert.NoError(t, err, "Running webdav integration tests")
}