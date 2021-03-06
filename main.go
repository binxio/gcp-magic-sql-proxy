package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    "google.golang.org/api/run/v1"
)

func main() {
    ctx := context.Background()
	stdout := log.New(os.Stdout, "", 1)
	// Get CloudRun service.
	runService, err := run.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Query for the right CloudRun instance.
	servicePtr, err := runService.Projects.Locations.Services.Get(
		fmt.Sprintf("projects/%s/locations/%s/services/%s",
			getEnv("GCP_PROJECT"),
			getEnv("REGION", "europe-west1"),
			getEnv("CR_SERVICE_NAME"))).Do()
	if err != nil {
		log.Fatal(err)
	}

	//Compose Argument to use with call of cloud_sql_proxy.
	sqlProxyCommandArg := fmt.Sprintf("-instances=%s=tcp:0.0.0.0:3306", getCloudRunVar(servicePtr, getEnv("CR_DB_ENV_NAME", "DB_SOCKET")))

	stdout.Print("Running: /cloud_sql_proxy " + sqlProxyCommandArg + "\n")
	cmd := exec.Command("/cloud_sql_proxy", sqlProxyCommandArg)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Create optional default value with ...string
func getEnv(key string, defaultValueOptional ...string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		log.Print(key + " not set.")
		if len(defaultValueOptional) > 0 {
			value = defaultValueOptional[0]
		} else {
			log.Fatal(key + " not set, and no default.")
		}
	}
	return value
}

// Get the environment variableValue for a variableKey.
// /cloudsql/ prefix if stripped to be able to use DB_SOCKET's 
func getCloudRunVar(servicePtr *run.Service, variableName string) string {
    variableVal := ""
	for _, con := range servicePtr.Spec.Template.Spec.Containers {
		for _, env := range con.Env {
			if variableName == env.Name {
                variableVal = strings.Replace(env.Value, "/cloudsql/", "", -1)
			}
		}
	}
	return variableVal
}

