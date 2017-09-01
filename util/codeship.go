package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CodeShip struct {
	AuthString  string
	AccessToken string
	Builds      map[string]*Build
}
type Build struct {
	UUID             string `json:"uuid"`
	ProjectName      string `json:"name"`
	ProjectUUID      string `json:"project_uuid"`
	OrganizationUUID string `json:"organization_uuid"`
	Status           string `json:"status"`
	FinishedAt       time.Time
	LastCheck        time.Time
}

func (c *CodeShip) GetToken() error {
	url := "https://api.codeship.com/v2/auth"
	payload := strings.NewReader("{}")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "text/plain")
	req.Header.Add("authorization", "Basic "+c.AuthString)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	c.AccessToken = m["access_token"].(string)
	orgs := m["organizations"].([]interface{})
	for _, org := range orgs {
		orgArr := org.(map[string]interface{})
		err = c.getProjects(orgArr["uuid"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CodeShip) getProjects(uuid string) error {
	url := "https://api.codeship.com/v2/organizations/" + uuid + "/projects"
	payload := strings.NewReader("{}")
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", c.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	projects := m["projects"].([]interface{})
	for _, project := range projects {
		p := project.(map[string]interface{})
		err = c.getBuilds(p["organization_uuid"].(string), p["uuid"].(string), p["name"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *CodeShip) getBuilds(org, project, name string) error {
	url := "https://api.codeship.com/v2/organizations/" + org + "/projects/" + project + "/builds?per_page=1"
	payload := strings.NewReader("{}")
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", c.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		builds := m["builds"].([]interface{})
		for _, build := range builds {
			b := build.(map[string]interface{})
			uuid := b["uuid"].(string)
			status := b["status"].(string)
			project := b["project_uuid"].(string)
			org := b["organization_uuid"].(string)
			if c.Builds[name] == nil {
				c.Builds[name] = &Build{
					UUID:             uuid,
					ProjectName:      name,
					ProjectUUID:      project,
					OrganizationUUID: org,
					Status:           status,
					LastCheck:        time.Now(),
				}
				if b["finished_at"] != nil {
					finished_at := b["finished_at"].(string)
					finished, err := time.Parse(time.RFC3339, finished_at)
					if err != nil {
						return err
					}
					c.Builds[name].FinishedAt = finished
				}
			} else if c.Builds[name].UUID != uuid || (c.Builds[name].UUID == uuid && c.Builds[name].Status != status) {
				c.Builds[name].UUID = uuid
				c.Builds[name].Status = status
				c.Builds[name].LastCheck = time.Now()
				if b["finished_at"] != nil {
					finished_at := b["finished_at"].(string)
					finished, err := time.Parse(time.RFC3339, finished_at)
					if err != nil {
						return err
					}
					c.Builds[name].FinishedAt = finished
				}
				if status == "success" {
					sendNotify("", name, "Build Successful")
				} else if status == "error" {
					sendNotify("", name, "Build Failed")
				} else if status != "testing" {
					fmt.Println(status)
				}
			}
		}
	} else if res.StatusCode == 401 {
		err = c.GetToken()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Codeship Error: %s", m["error_message"])
	}
	return nil
}

func (c *CodeShip) PollChanges() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			week := time.Now().AddDate(0, 0, -7)
			month := time.Now().AddDate(0, 1, 0)
			months := time.Now().AddDate(0, 6, 0)
			for _, build := range c.Builds {
				if build.FinishedAt.After(week) {
					if build.LastCheck.After(time.Now().Add(-30 * time.Second)) {
						continue
					}
				} else if build.FinishedAt.After(month) {
					if build.LastCheck.After(time.Now().Add(-1 * time.Minute)) {
						continue
					}
				} else if build.FinishedAt.After(months) {
					if build.LastCheck.After(time.Now().Add(-90 * time.Second)) {
						continue
					}

				} else if !build.FinishedAt.IsZero() {
					if build.LastCheck.After(time.Now().Add(-2 * time.Minute)) {
						continue
					}
				}
				c.getBuilds(build.OrganizationUUID, build.ProjectUUID, build.ProjectName)
			}
		}
	}
}
