package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
)

// ebInterface represents the Elastic Beanstalk interface.
type ebInterface interface {
	filterEnvironments(res *elasticbeanstalk.DescribeEnvironmentsOutput, env string) ([]string, error)
	describeEnvironments() (*elasticbeanstalk.DescribeEnvironmentsOutput, error)
	saveConfig(appName, envName, envId string) (*elasticbeanstalk.CreateConfigurationTemplateOutput, error)
}

// saveConfig saves the Elastic Beanstalk configuration.
func (c *awsClient) saveConfig(appName, envName, envId string) (*elasticbeanstalk.CreateConfigurationTemplateOutput, error) {
	// Set the description and template name.
	description := "created by ebsc"
	templateName := fmt.Sprintf("ebsc-%v", envName)

	// Save the Elastic Beanstalk configuration.
	res, err := c.eb.CreateConfigurationTemplate(context.TODO(), &elasticbeanstalk.CreateConfigurationTemplateInput{
		ApplicationName: aws.String(appName),
		Description:     aws.String(description),
		EnvironmentId:   aws.String(envId),
		TemplateName:    aws.String(templateName),
	})

	if err != nil {
		return nil, err
	}

	// Return the configuration.
	return res, nil
}

// filterEnvironment filters the Elastic Beanstalk environments.
func (c *awsClient) filterEnvironments(res *elasticbeanstalk.DescribeEnvironmentsOutput, env string) ([]string, error) {
	// Create a slice to store the filtered environments.
	filtered := []string{}

	// Filter the environments.
	for _, e := range res.Environments {
		if strings.Contains(*e.EnvironmentName, env) {
			appEnvName := fmt.Sprintf("%v/%v/%v", *e.ApplicationName, *e.EnvironmentName, *e.EnvironmentId)
			filtered = append(filtered, appEnvName)
		}
	}

	// Check if there are any environments.
	if len(filtered) == 0 {
		return nil, fmt.Errorf("failed to filter environments")
	}

	// Return the filtered environments.
	return filtered, nil
}

// describeEnvironments describes the Elastic Beanstalk environments.
func (c *awsClient) describeEnvironments() (*elasticbeanstalk.DescribeEnvironmentsOutput, error) {
	// Describe the Elastic Beanstalk environments.
	res, err := c.eb.DescribeEnvironments(context.TODO(), &elasticbeanstalk.DescribeEnvironmentsInput{})

	if err != nil {
		return res, err
	}

	// Check if there are any environments.
	if len(res.Environments) == 0 {
		return nil, fmt.Errorf("no environments found")
	}

	// Return the environments.
	return res, nil
}
