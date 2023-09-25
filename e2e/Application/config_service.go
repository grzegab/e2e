package Application

import (
	"e2e/Interface"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

type JWTdata struct {
	Type    string `json:"token_type"`
	Expire  int64  `json:"expires_in"`
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type Config struct {
	JWT             JWTdata
	DocUrl          string
	manualEndpoints []string
	TestingUrl      string
}

func GetConfig() Config {
	Interface.StartMsgLine()
	Interface.PrintSimpleText("Checking env file...")

	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		Interface.PrintSuccessText("Failed")
		Interface.EndMsgLine()

		Interface.ErrorMsg("Error reading ENV file: %s", []any{err.Error()})
		os.Exit(1)
	}

	Interface.PrintSuccessText("OK")
	Interface.EndMsgLine()

	config := Config{
		JWT:    JWTdata{},
		DocUrl: "",
	}

	Interface.PrintSimpleText("Checking for testing url...")
	config.TestingUrl = getTestingUrl("E2E_TEST_URL")
	Interface.PrintSuccessText("[" + config.TestingUrl + "]")
	Interface.EndMsgLine()
	Interface.PrintSimpleText("Checking for login url...")
	loginUrl := checkEnvUrlVariable("E2E_LOGIN_URL", "login")
	if loginUrl != "" {
		Interface.PrintSuccessText("[" + loginUrl + "]")
	} else {
		Interface.PrintFailText("NO LOGIN URL AVAILABLE")
	}
	Interface.EndMsgLine()
	if loginUrl != "" {
		Interface.PrintSimpleText("Obtaining JWT token...")
		user := os.Getenv("E2E_LOGIN_USERNAME")
		password := os.Getenv("E2E_LOGIN_PASSWORD")

		if user == "" || password == "" {
			Interface.PrintFailText("NO CREDENTIALS PROVIDED")
			Interface.EndMsgLine()
			Interface.PrintFailText("REMOVE LOGIN URL IF NOT NEEDED")
			os.Exit(1)
		}

		jwtData, loginErr := getJWTFromEnvVariables(loginUrl, user, password)

		if loginErr != nil {
			Interface.PrintFailText("INVALID CREDENTIALS PROVIDED")
			Interface.EndMsgLine()
			Interface.PrintFailText("CHECK PROVIDED CREDENTIALS")
			os.Exit(1)
		}

		config.JWT = jwtData

		Interface.PrintSuccessText("OK")
		Interface.EndMsgLine()
	}

	documentationUrl := checkEnvUrlVariable("E2E_DOC_URL", "documentation")
	Interface.PrintSimpleText("Checking for documentation url...")
	if documentationUrl == "" {
		Interface.PrintFailText("NO URL TO READ DOCUMENTATION")
	} else {
		config.DocUrl = checkEnvUrlVariable("E2E_DOC_URL", "documentation")
		Interface.PrintSuccessText("[" + config.DocUrl + "]")
	}

	config.manualEndpoints = loadManualEndpoints()

	Interface.EndMsgLine()
	return config
}

func checkEnvUrlVariable(envVariable string, variableType string) string {
	url := os.Getenv(envVariable)
	if url == "" {
		Interface.WarningMsg("Url for %s not provided", []any{variableType})
	} else {
		request, err := Interface.CreateRequest(url, http.MethodGet, "")

		if err != nil {
			Interface.ErrorMsg("Error while creating request: %s", []any{err.Error()})
		}

		success := Interface.PingRequest(request)
		if !success {
			Interface.WarningMsg("Url not reachable %s", []any{url})

			return ""
		}
	}

	return url
}

func getJWTFromEnvVariables(url string, login string, pass string) (JWTdata, error) {
	var err error
	var jwtData JWTdata
	loginBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, login, pass)
	request, err := Interface.CreateRequest(url, http.MethodPost, loginBody)

	if err != nil {
		Interface.ErrorMsg("Error while creating request: %s", []any{err.Error()})

		return jwtData, err
	}

	statusCode, body, _ := Interface.FullRequest(request)
	if statusCode != 200 {
		err = errors.New("bad response from JWT source")

		return jwtData, err
	}

	err = json.Unmarshal(body, &jwtData)

	if err != nil {
		Interface.ErrorMsg("Error decoding JWT: %s", []any{err.Error()})
	}

	return jwtData, err
}

func getTestingUrl(envVariable string) string {
	url := os.Getenv(envVariable)

	//@TODO: check if url is ok by looking at string not ping

	return url
}

func loadManualEndpoints() []string {
	manualEndpointsDir := "var/configs"
	// @TODO: check for pipeline, if not works replace with: "/go/src/e2e/var/configs"

	Interface.StartMsgLine()
	Interface.PrintSimpleText("Checking for manual endpoints files...")
	if _, err := os.Stat(manualEndpointsDir); os.IsNotExist(err) {
		Interface.PrintFailText("Nothing found")
		return []string{}
	}

	acceptedFiles, rejectedFiles := Interface.SearchForTestFiles(manualEndpointsDir)

	if acceptedFiles != nil {
		Interface.PrintSuccessText("Accepted: " + strconv.Itoa(len(acceptedFiles)))
	}

	if rejectedFiles != nil {
		Interface.PrintFailText("Rejected: " + strconv.Itoa(len(acceptedFiles)))
	}

	Interface.PrintSuccessText("OK")
	Interface.EndMsgLine()

	return acceptedFiles
}
