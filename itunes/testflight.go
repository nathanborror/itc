package itunes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const testflightPath = "https://itunesconnect.apple.com/testflight/v2"

// Tester contains properties about a tester.
type Tester struct {
	ID                string `json:"id"`
	ProviderID        int    `json:"providerId"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Email             string `json:"email"`
	LatestInstallInfo struct {
		LatestInstalledAppAdamID    string `json:"latestInstalledAppAdamId"`
		LatestInstalledBuildID      string `json:"latestInstalledBuildId"`
		LatestInstalledDate         string `json:"latestInstalledDate"`
		LatestInstalledShortVersion string `json:"latestInstalledShortVersion"`
		LatestInstalledVersion      string `json:"latestInstalledVersion"`
	} `json:"latestInstallInfo"`
	AppAdamID              int      `json:"appAdamId"`
	AccountID              string   `json:"accountId"`
	InviteToken            string   `json:"inviteToken"`
	Status                 string   `json:"status"`
	StatusModTime          string   `json:"statusModTime"`
	LatestInstalledTrain   string   `json:"latestInstalledTrain"`
	LatestInstalledVersion string   `json:"latestInstalledVersion"`
	Groups                 []string `json:"groups"`
	InstallCount           int      `json:"installCount"`
	SessionCount           int      `json:"sessionCount"`
	CrashCount             int      `json:"crashCount"`
}

// TesterListResponse contains a list of Testers.
type TesterListResponse struct {
	Data []Tester `json:"data"`
}

// TesterList returns a TesterListResponse that contains all the testers for
// a given app identifier.
func (c *Client) TesterList(providerID int, appID int, paging *Paging) (*TesterListResponse, error) {
	host := fmt.Sprintf("%s/providers/%d/apps/%d/testers", testflightPath, providerID, appID)
	req, err := c.NewRequest("GET", host, nil)
	req.URL.RawQuery = paging.Encode(req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, c.ErrorForServiceErrors(body)
	}
	var testers TesterListResponse
	if err := json.Unmarshal(body, &testers); err != nil {
		return nil, err
	}
	return &testers, nil
}

// TesterGroup contains information about a test group.
type TesterGroup struct {
	ID                     string `json:"id"`
	ProviderID             int    `json:"providerId"`
	AppAdamID              int    `json:"appAdamId"`
	Name                   string `json:"name"`
	Created                string `json:"created"`
	Active                 bool   `json:"active"`
	IsInternalGroup        bool   `json:"isInternalGroup"`
	IsDefaultExternalGroup bool   `json:"isDefaultExternalGroup"`
}

// TesterGroupResponse contains a list of TesterGroups.
type TesterGroupResponse struct {
	Data []TesterGroup `json:"data"`
}

// TesterGroups returns a list of tester groups.
func (c *Client) TesterGroups(providerID int, appID int, paging *Paging) (*TesterGroupResponse, error) {
	host := fmt.Sprintf("%s/providers/%d/apps/%d/groups", testflightPath, providerID, appID)
	req, err := c.NewRequest("GET", host, nil)
	req.URL.RawQuery = paging.Encode(req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, c.ErrorForServiceErrors(body)
	}
	var testers TesterGroupResponse
	if err := json.Unmarshal(body, &testers); err != nil {
		return nil, err
	}
	return &testers, nil
}

// CreateTester contains properties for adding a test user.
type CreateTester struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// TesterCreate adds a new tester to the app.
func (c *Client) TesterCreate(testers []CreateTester, providerID int, appID int, groupID string) error {
	data, err := json.Marshal(testers)
	if err != nil {
		return err
	}
	host := fmt.Sprintf("%s/providers/%d/apps/%d/groups/%s/testers", testflightPath, providerID, appID, groupID)
	req, err := c.NewRequest("POST", host, bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error (%d): %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// TesterDelete deletes a tester.
func (c *Client) TesterDelete(testerID string, providerID int, appID int) error {
	host := fmt.Sprintf("%s/providers/%d/apps/%d/testers/%s", testflightPath, providerID, appID, testerID)
	req, err := c.NewRequest("DELETE", host, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error (%d): %s", resp.StatusCode, resp.Status)
	}
	return nil
}
