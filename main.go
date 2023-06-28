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

	logger.Infof("Apple ID: %s", connection.AppleIDConnection.AppleID)

	//	fmt.Println("This is the value specified for the input 'example_step_input':", os.Getenv("example_step_input"))
	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "EXAMPLE_STEP_OUTPUT", "--value", "the value you want to share").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}
