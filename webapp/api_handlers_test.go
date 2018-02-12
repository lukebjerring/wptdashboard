package webapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"github.com/stretchr/testify/assert"
	base "github.com/w3c/wptdashboard/shared"
)

var (
	staticDataTime, _   = time.Parse(time.RFC3339, "2017-10-18T00:00:00Z")
	staticRunSHA        = "b952881825"
	summaryUrlFmtString = "/static/wptd/" + staticRunSHA + "/%s"
	chromeRun           = base.TestRun{
		BrowserName:    "chrome",
		BrowserVersion: "63.0",
		OSName:         "linux",
		OSVersion:      "3.16",
		Revision:       staticRunSHA,
		ResultsURL:     fmt.Sprintf(summaryUrlFmtString, "chrome-63.0-linux-summary.json.gz"),
		CreatedAt:      staticDataTime,
	}
	staticTestRuns = []base.TestRun{
		chromeRun,
		{
			BrowserName:    "edge",
			BrowserVersion: "15",
			OSName:         "windows",
			OSVersion:      "10",
			Revision:       staticRunSHA,
			ResultsURL:     fmt.Sprintf(summaryUrlFmtString, "edge-15-windows-10-sauce-summary.json.gz"),
			CreatedAt:      staticDataTime,
		},
		{
			BrowserName:    "firefox",
			BrowserVersion: "57.0",
			OSName:         "linux",
			OSVersion:      "*",
			Revision:       staticRunSHA,
			ResultsURL:     fmt.Sprintf(summaryUrlFmtString, "firefox-57.0-linux-summary.json.gz"),
			CreatedAt:      staticDataTime,
		},
		{
			BrowserName:    "safari",
			BrowserVersion: "10",
			OSName:         "macos",
			OSVersion:      "10.12",
			Revision:       staticRunSHA,
			ResultsURL:     fmt.Sprintf(summaryUrlFmtString, "safari-10-macos-10.12-sauce-summary.json.gz"),
			CreatedAt:      staticDataTime,
		},
	}
)

func TestApiTestRunsHandler_Browser(t *testing.T) {
	req, instance := AETestRequest("http://wpt.fyi/api/runs?browser=chrome")
	defer instance.Close()
	if err := addStaticRuns(req); err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	apiTestRunsHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Just the chrome run.
	expected, _ := json.Marshal([]base.TestRun{chromeRun})
	actual, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expected, actual)
}

func TestApiTestRunsHandler_MaxCount(t *testing.T) {
	req, instance := AETestRequest("http://wpt.fyi/api/runs?max-count=1")
	defer instance.Close()
	// Add the runs twice
	if err := addStaticRuns(req); err != nil {
		panic(err)
	} else if err := addStaticRuns(req); err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	apiTestRunsHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// But only expect one of each.
	expected, _ := json.Marshal(staticTestRuns)
	actual, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expected, actual)
}

func AETestRequest(url string) (*http.Request, aetest.Instance) {
	var aetestInstance aetest.Instance
	var err error
	// Need strongly consistent so GetAll actually returns results (without waiting).
	opts := &aetest.Options{StronglyConsistentDatastore: true}
	if aetestInstance, err = aetest.NewInstance(opts); err != nil {
		panic(err)
	}

	var req *http.Request
	req, err = aetestInstance.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return req, aetestInstance
}

func addStaticRuns(req *http.Request) error {
	ctx := appengine.NewContext(req)
	keys := make([]*datastore.Key, len(staticTestRuns))
	for i := range staticTestRuns {
		keys[i] = datastore.NewIncompleteKey(ctx, "TestRun", nil)
	}
	if _, err := datastore.PutMulti(ctx, keys, staticTestRuns); err != nil {
		panic(err)
	}
	log.Printf("Added %v %s entities", len(staticTestRuns), "TestRun")
	return nil
}
