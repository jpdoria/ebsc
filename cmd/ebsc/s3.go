package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Interface interface {
	downloadConfig(dirName, appName, envName string) error
	searchConfigBucket() (string, error)
}

// downloadConfig downloads the configuration from the S3 bucket.
func (c *awsClient) downloadConfig(dirName, appName, envName string) error {
	// Retrieve the configuration bucket.
	bucketName, err := c.searchConfigBucket()

	if err != nil {
		return err
	}

	// Get the configuration file.
	keyPath := fmt.Sprintf("resources/templates/%v/ebsc-%v", appName, envName)
	res, err := c.s3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyPath),
	})

	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Create a new file to store the configuration.
	fileName := fmt.Sprintf("./%v/%v", dirName, envName)
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	// Write the configuration to the file.
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	return nil
}

// searchConfigBucket searches for the configuration bucket.
func (c *awsClient) searchConfigBucket() (string, error) {
	// Filter the bucket name using this format: elasticbeanstalk-<region>.
	filter := fmt.Sprintf("elasticbeanstalk-%v", *c.region)

	// List the S3 buckets.
	res, err := c.s3.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		return "", err
	}

	// Search for the bucket name.
	var bucketName string

	for _, b := range res.Buckets {
		if strings.Contains(*b.Name, filter) {
			bucketName = *b.Name
		}
	}

	return bucketName, nil
}
