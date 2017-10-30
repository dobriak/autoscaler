package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Client struct
type Client struct {
	BaseURL    string
	UserAgent  string
	httpClient *http.Client
	Token      string
}

//client - reference to a http client
var client Client

func init() {
	//get baseurl from env
	baseURL := os.Getenv("AS_BASEURL")
	if len(baseURL) == 0 {
		log.Panicln("Please supply AS_BASEURL")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr, Timeout: time.Second * 10}
	client = Client{baseURL, "tester", httpClient, ""}
	client.auth()
}

//GetAllMarathonApps func
func (c *Client) GetAllMarathonApps() MarathonApps {
	req, err := c.newRequest("GET", "/service/marathon/v2/apps", nil)
	if err != nil {
		return MarathonApps{}
	}

	body, _ := c.do(req)
	var result MarathonApps
	err = json.Unmarshal(body, &result)
	if err != nil {
		return MarathonApps{}
	}
	return result
}

//AppExists for an App
func (c *Client) AppExists(a *App) bool {
	var mApps = c.GetAllMarathonApps()
	for _, mApp := range mApps.Apps {
		if a.AppID == mApp.ID {
			return true
		}
	}
	return false
}

//GetMarathonApp func
func (c *Client) GetMarathonApp(appID string) MarathonApp {
	req, err := c.newRequest("GET", fmt.Sprintf("/service/marathon/v2/apps/%s", appID), nil)
	if err != nil {
		return MarathonApp{}
	}

	body, _ := c.do(req)

	var result MarathonApp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return MarathonApp{}
	}
	return result
}

//ScaleMarathonApp scales to target number of instances
func (c *Client) ScaleMarathonApp(appID string, instances int) {
	data := MarathonAppInstances{instances}
	req, err := c.newRequest("PUT", fmt.Sprintf("/service/marathon/v2/apps/%s", appID), data)
	if err != nil {
		log.Errorln(err)
	}

	body, err := c.do(req)
	if err != nil {
		log.Errorln(err)
	}

	var resp MarathonScaleResult
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Infof("Successfully scaled app %s: version %s, deploymentId %s",
			appID, resp.Version, resp.DeploymentID)
	}

	fmt.Println(resp)
}

//GetTaskStats func
func (c *Client) GetTaskStats(taskID string, slaveID string) TaskStats {
	req, err := c.newRequest("GET", fmt.Sprintf("/slave/%s/monitor/statistics.json", slaveID), nil)
	if err != nil {
		log.Errorln("Error querying statistics.json")
		return TaskStats{}
	}
	body, _ := c.do(req)
	var result []TaskStats
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Errorln("Error unmarshalling TasksStats")
		return TaskStats{}
	}
	for _, ts := range result {
		if ts.ExecutorID == taskID {
			return ts
		}
	}
	return TaskStats{}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, c.BaseURL+path, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "token="+c.Token)
	} else {
		req.Header.Del("Authorization")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, err
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		if 401 == resp.StatusCode {
			log.Infoln("Authentication expired. Re-authorizing account")
			log.Panicln("Not implemented")
		} else {
			return nil, fmt.Errorf("%s", body)
		}
	}
	return body, err

}

//isJSON: false if not a json string
func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func downloadFile(filepath string, path string) (err error) {
	// Check if exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) auth() {
	// Do we have username/password?
	user := os.Getenv("AS_USERID")
	pass := os.Getenv("AS_PASSWORD")
	if len(user) == 0 || len(pass) == 0 {
		log.Panicln("Set AS_USERID and AS_PASSWORD env vars")
	}
	usrPass := DcosBasicAuth{user, pass}

	req, err := client.newRequest("POST", "/acs/api/v1/auth/login", usrPass)
	if err != nil {
		fmt.Println(err)
		log.Panicln("Error trying to auth")
	}

	body, _ := c.do(req)
	var result DcosAuthResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(body)
		fmt.Println(err)
		log.Panicln("Couldn't convert to dcosAuthResponse")
	}

	log.Infof("Token obtained: %s", result.Token)
	c.Token = result.Token

}
