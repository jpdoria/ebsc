package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// awsClient represents the AWS client.
type awsClient struct {
	region *string
	eb     *elasticbeanstalk.Client
	s3     *s3.Client
}

// backup represents the backup struct.
type backup struct {
	eb  ebInterface
	dir dirInterface
	s3  s3Interface
}

func main() {
	// Define the flags.
	env := flag.String("env", "dev", "the environment to filter e.g., dev, qa, prod")
	region := flag.String("region", "us-east-1", "the region code of the environment")

	flag.Parse()

	// Validate the flags.
	if *env != "dev" && *env != "qa" && *env != "prod" {
		fmt.Println("invalid environment")
		os.Exit(1)
	}

	if *region == "" {
		fmt.Println("invalid region")
		os.Exit(1)
	}

	fmt.Println("starting backup")

	// Load the SDK's default configuration.
	sdkConfig, _ := config.LoadDefaultConfig(context.TODO())

	// Create a new client with the Elastic Beanstalk client and environment.
	b := &backup{
		dir: &dirManager{},
		eb:  &awsClient{eb: elasticbeanstalk.NewFromConfig(sdkConfig)},
		s3: &awsClient{
			s3:     s3.NewFromConfig(sdkConfig),
			region: region,
		},
	}

	// Get the current date and time and add the value to the dirManager struct.
	b.dir.getDateTime()

	// Describe the Elastic Beanstalk environments.
	res, err := b.eb.describeEnvironments()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Filter the Elastic Beanstalk environments.
	out, err := b.eb.filterEnvironments(res, *env)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a backup of each environment. Ignore the error of creating a directory
	// because we are going to reuse the directory for other environments.
	for _, o := range out {
		appName := strings.Split(o, "/")[0]
		envName := strings.Split(o, "/")[1]
		envId := strings.Split(o, "/")[2]
		var path string

		fmt.Printf("creating a backup for %q\n", envName)

		// Check if there's a directory that contains the application name.
		// If there is, skip creating a directory.
		exists, path := b.dir.directoryExists(strings.ToLower(appName))

		// Create a directory using the application name.
		if !exists && path == "" {
			path, err = b.dir.createDirectory(strings.ToLower(appName))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// Create a backup of the environment.
		_, err = b.eb.saveConfig(appName, envName, envId)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Download the config from S3.
		err = b.s3.downloadConfig(path, appName, envName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("backup has been created for %q\n", envName)
	}

	// Compress the backup directory.
	err = b.dir.compressDirectory()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("backup has been completed")
}
