package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/retryhttp"
	"github.com/bitrise-io/go-xcode/devportalservice"
)

func main() {
	logger := log.NewLogger()
	http := retryhttp.NewClient(logger).StandardClient()
	buildUrl := os.Getenv("BITRISE_BUILD_URL")
	buildToken := os.Getenv("BITRISE_BUILD_API_TOKEN")
	provider := devportalservice.NewBitriseClient(http, buildUrl, buildToken)

	if provider == nil {
		logger.Errorf("Failed to connect to Bitrise.")
		os.Exit(1)
	}

	logger.Infof("%s", "Successfully connected to Bitrise.")

	connection, err := provider.GetAppleDeveloperConnection()
	if err != nil {
		logger.Errorf("Failed to connect to Bitrise: %s", err)
		os.Exit(1)
	}

	session, err := connection.AppleIDConnection.FastlaneLoginSession()
	if err != nil {
		logger.Errorf("Failed to get Fastlane session: %s", err)
		os.Exit(1)
	}

	environmentVariables := map[string]string{
		"FASTLANE_USER":     connection.AppleIDConnection.AppleID,
		"FASTLANE_PASSWORD": connection.AppleIDConnection.Password,
		"FASTLANE_SESSION":  session,
	}
	setEnv(environmentVariables)
	os.Exit(0)
}

func setEnv(vars map[string]string) {
	for key, value := range vars {
		// You can find more usage examples on envman's GitHub page
		//  at: https://github.com/bitrise-io/envman
		cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", key, "--value", value).CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
			os.Exit(1)
		}
	}
}
