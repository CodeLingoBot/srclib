package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/alexsaveliev/go-colorable-wrapper"

	"golang.org/x/tools/godoc/vfs"

	"sourcegraph.com/sourcegraph/srclib"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

type nopWriteCloser struct{}

func (w nopWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w nopWriteCloser) Close() error {
	return nil
}

func isDir(dir string) bool {
	di, err := os.Stat(dir)
	return err == nil && di.IsDir()
}

func execCmd(prog string, arg ...string) error {
	cmd := exec.Command(prog, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	log.Println("Running ", cmd.Args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: %s", cmd.Args, err)
	}
	return nil
}

func execSrcCmd(arg ...string) error {
	if len(arg) == 0 {
		log.Fatal("attempted to execute 'srclib' command with no arguments")
	}
	c := append(strings.Split(srclib.CommandName, " "), arg...)
	return execCmd(c[0], c[1:]...)
}

func SourceUnitMatchesArgs(specified []string, u *unit.SourceUnit) bool {
	var match bool
	if len(specified) == 0 {
		match = true
	} else {
		for _, unitSpec := range specified {
			if string(u.ID()) == unitSpec || u.Name == unitSpec {
				match = true
				break
			}
		}
	}

	return match
}

func PrintJSON(v interface{}, prefix string) {
	data, err := json.MarshalIndent(v, prefix, "  ")
	if err != nil {
		log.Fatal(err)
	}
	colorable.Println(string(data))
}

func OpenInputFiles(extraArgs []string) map[string]io.ReadCloser {
	inputs := make(map[string]io.ReadCloser)
	if len(extraArgs) == 0 {
		inputs["<stdin>"] = os.Stdin
	} else {
		for _, name := range extraArgs {
			f, err := os.Open(name)
			if err != nil {
				log.Fatal(err)
			}
			inputs[name] = f
		}
	}
	return inputs
}

func CloseAll(files map[string]io.ReadCloser) {
	for _, rc := range files {
		rc.Close()
	}
}

func readJSONFile(file string, v interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func readJSONFileFS(fs vfs.FileSystem, file string, v interface{}) (err error) {
	f, err := fs.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err == nil {
			err = err2
		}
	}()
	return json.NewDecoder(f).Decode(v)
}

func bytesString(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	logn := func(n, b float64) float64 {
		return math.Log(n) / math.Log(b)
	}
	e := math.Floor(logn(float64(s), 1000))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(1000, math.Floor(e))*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+"%s", val, suffix)
}

func percent(num, denom int) float64 {
	return 100 * float64(num) / float64(denom)
}

// parseRepoAndCommitID parses strings like "example.com/repo" and
// "example.com/repo@myrev".
func parseRepoAndCommitID(repoAndCommitID string) (uri, commitID string) {
	if i := strings.Index(repoAndCommitID, "@"); i != -1 {
		return repoAndCommitID[:i], repoAndCommitID[i+1:]
	}
	return repoAndCommitID, ""
}
