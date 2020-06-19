package integration

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudfoundry/dagger"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var uri string

func TestIntegration(t *testing.T) {
	var Expect = NewWithT(t).Expect

	root, err := filepath.Abs("./..")
	Expect(err).ToNot(HaveOccurred())

	uri, err = dagger.PackageBuildpack(root)
	Expect(err).NotTo(HaveOccurred())

	// HACK: we need to fix dagger and the package.sh scripts so that this isn't required
	uri = fmt.Sprintf("%s.tgz", uri)

	suite := spec.New("Integration", spec.Report(report.Terminal{}))
	suite("SimpleApp", testSimpleApp)
	suite("Caching", testCaching)
	suite("Logging", testLogging)
	suite("NoConfApp", testNoConfApp)

	defer AfterSuite(t)
	suite.Run(t)
}

func AfterSuite(t *testing.T) {
	var Expect = NewWithT(t).Expect

	Expect(dagger.DeleteBuildpack(uri)).To(Succeed())
}

func GetGitVersion() (string, error) {
	gitExec := pexec.NewExecutable("git")
	revListOut := bytes.NewBuffer(nil)

	err := gitExec.Execute(pexec.Execution{
		Args:   []string{"rev-list", "--tags", "--max-count=1"},
		Stdout: revListOut,
	})
	if err != nil {
		return "", err
	}

	stdout := bytes.NewBuffer(nil)
	err = gitExec.Execute(pexec.Execution{
		Args:   []string{"describe", "--tags", strings.TrimSpace(revListOut.String())},
		Stdout: stdout,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(stdout.String(), "v")), nil
}
