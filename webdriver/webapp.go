package webdriver

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

type Webapp interface {
	// GetWebappUrl returns the URL for the given path on the running webapp.
	GetWebappUrl(path string) string
}

// GetStagingInstance gets the base url for a wptdashboard staging instance,
// assuming it was deployed via the git/git-deploy script.
func GetStagingInstance() (Webapp, error) {
	webapp := &staging_instance{}
	getVersion := exec.Command("bash", "-c", "echo $(source ../git/staging-version.sh; print_version)")
	if version, err := getVersion.Output(); err != nil {
		return nil, err
	} else {
		scanner := bufio.NewScanner(bytes.NewReader(version))
		if !scanner.Scan() {
			panic("Unable to read output of staging-version.sh")
		}
		webapp.version = scanner.Text()
	}
	return Webapp(webapp), nil
}

type staging_instance struct {
	version string
}

func (i *staging_instance) GetWebappUrl(path string) string {
	return fmt.Sprintf("http://%s-dot-wptdashboard.appspot.com", i.version)
}
