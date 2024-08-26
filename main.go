package main

import (
    "context"
    "flag"
    "fmt"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
    "os"
    "strings"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    "github.com/aws/aws-sdk-go-v2/service/iam"
    "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Account struct {
    Name       string   `yaml:"name"`
    AccessKey  string   `yaml:"access_key"`
    SecretKey  string   `yaml:"secret_key"`
    Regions    []string `yaml:"regions"`
}

type Config struct {
    Accounts []Account `yaml:"accounts"`
}

var (
	add   = flag.Bool("add", false, "Add an AWS account")
	remove   = flag.Bool("remove", false, "Remove an AWS account")
	find   = flag.Bool("find", false, "find an AWS account")

	instanceID   = flag.String("instance-id", "", "Instance ID to search for")
    ipAddress    = flag.String("ip", "", "IP address to search for")
    searchIAM    = flag.String("iam", "", "Search for IAM user associated with the access key")
    yamlFile     = flag.String("config", "aws-accounts.yaml", "Path to YAML file with AWS credentials")
    accountName  = flag.String("account-name", "", "AWS account name to add or remove")
    accessKey    = flag.String("access-key", "", "AWS access key to add")
    secretKey    = flag.String("secret-key", "", "AWS secret key to add")
    regions      = flag.String("regions", "us-east-1", "Comma-separated AWS regions to add")
    showHelp     = flag.Bool("help", false, "Show help information")
)

func main() {
    flag.Parse()

    if *showHelp {
        showHelpMessage()
        return
    }

    if _, err := os.Stat(*yamlFile); os.IsNotExist(err) {
        createConfig(*yamlFile)
    }

    switch {
    case *add:
        if *accountName == "" || *accessKey == "" || *secretKey == "" {
            log.Fatalf("All of --account-name, --access-key, and --secret-key must be provided for 'add' operation")
        }
        addAccount(*yamlFile, *accountName, *accessKey, *secretKey, *regions)
    case *remove:
        if *accountName == "" {
            log.Fatalf("The --account-name must be provided for 'remove' operation")
        }
        removeAccount(*yamlFile, *accountName)
    case *find:
        if *searchIAM!= "" {
            search(*yamlFile, *instanceID, *ipAddress, *searchIAM)
        } else if *instanceID != "" || *ipAddress != "" {
            search(*yamlFile, *instanceID, *ipAddress, *searchIAM)
        } else {
            log.Fatalf("Specify either --instance-id or --ip to perform search")
        }
    default:
        log.Fatalf("No valid operation flag provided. Use --add, --remove, or --search.")
    }
}

func showHelpMessage() {
    fmt.Println(`Usage: gview [options]

Options:
  --add 				  Add an AWS account
  --remove				  Remove an AWS account
  --find				  Find an AWS account
  --instance-id <id>      Instance ID to search for
  --ip <ip-address>       IP address to search for
  --iam            		  Search for IAM user associated with the access key
  --config <file>         Path to YAML file with AWS credentials (default: aws-accounts.yaml)
  --account-name <name>   AWS account name to add or remove
  --access-key <key>      AWS access key to add
  --secret-key <key>      AWS secret key to add
  --regions <regions>     Comma-separated AWS regions to add
  --help                  Show this help message
`)
}

func createConfig(yamlFile string) {
    config := Config{Accounts: []Account{}}
    data, err := yaml.Marshal(&config)
    if err != nil {
        log.Fatalf("Error marshaling YAML: %v", err)
    }

    if err := ioutil.WriteFile(yamlFile, data, 0600); err != nil {
        log.Fatalf("Error writing YAML file: %v", err)
    }
}

func readConfig(yamlFile string) (Config, error) {
    data, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return Config{}, fmt.Errorf("error reading YAML file: %v", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return Config{}, fmt.Errorf("error parsing YAML file: %v", err)
    }
    return config, nil
}

func writeConfig(yamlFile string, config Config) error {
    data, err := yaml.Marshal(&config)
    if err != nil {
        return fmt.Errorf("error marshaling YAML: %v", err)
    }

    if err := ioutil.WriteFile(yamlFile, data, 0600); err != nil {
        return fmt.Errorf("error writing YAML file: %v", err)
    }
    return nil
}

func addAccount(yamlFile, name, accessKey, secretKey, regions string) {
    config, err := readConfig(yamlFile)
    if err != nil {
        log.Fatal(err)
    }

    for _, account := range config.Accounts {
        if account.Name == name {
            log.Fatalf("Account %s already exists", name)
        }
    }

    newAccount := Account{
        Name:      name,
        AccessKey: accessKey,
        SecretKey: secretKey,
        Regions:   strings.Split(regions, ","),
    }
    config.Accounts = append(config.Accounts, newAccount)

    if err := writeConfig(yamlFile, config); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Account %s added successfully\n", name)
}

func removeAccount(yamlFile, name string) {
    config, err := readConfig(yamlFile)
    if err != nil {
        log.Fatal(err)
    }

    var updatedAccounts []Account
    for _, account := range config.Accounts {
        if account.Name != name {
            updatedAccounts = append(updatedAccounts, account)
        }
    }

    if len(updatedAccounts) == len(config.Accounts) {
        log.Fatalf("Account %s not found", name)
    }

    config.Accounts = updatedAccounts

    if err := writeConfig(yamlFile, config); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Account %s removed successfully\n", name)
}

func search(yamlFile, instanceID, ipAddress string, searchIAM string) {
    cfg, err := readConfig(yamlFile)
    if err != nil {
        log.Fatalf("Failed to read configuration: %v", err)
    }

    for _, account := range cfg.Accounts {
        for _, region := range account.Regions {
            awsCfg, err := config.LoadDefaultConfig(context.TODO(),
                config.WithRegion(region),
                config.WithCredentialsProvider(
                    aws.NewCredentialsCache(
                        credentials.NewStaticCredentialsProvider(account.AccessKey, account.SecretKey, ""),
                    ),
                ),
            )
            if err != nil {
                log.Printf("Unable to load SDK config for account %s, region %s: %v", account.Name, region, err)
                continue
            }

			if searchIAM != "" {
                searchIAMUser(awsCfg, account.Name,searchIAM)
            } else {
                searchEC2Instances(awsCfg, account.Name, region, instanceID, ipAddress)
            }
        }
    }
}

func searchEC2Instances(cfg aws.Config, accountName, region, instanceID, ipAddress string) {
    svc := ec2.NewFromConfig(cfg)

    var results []*ec2.DescribeInstancesOutput

    if instanceID != "" {
        result, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
            InstanceIds: []string{instanceID},
        })
        if err != nil {
            log.Printf("Error describing instance %s in account %s, region %s: %v\n", instanceID, accountName, region, err)
        } else {
            results = append(results, result)
        }
    }

    if ipAddress != "" {
        privateIPFilter := &ec2.DescribeInstancesInput{
            Filters: []types.Filter{
                {
                    Name:   aws.String("private-ip-address"),
                    Values: []string{ipAddress},
                },
            },
        }
        publicIPFilter := &ec2.DescribeInstancesInput{
            Filters: []types.Filter{
                {
                    Name:   aws.String("ip-address"),
                    Values: []string{ipAddress},
                },
            },
        }

        privateResult, err := svc.DescribeInstances(context.TODO(), privateIPFilter)
        if err != nil {
            log.Printf("Error describing instances by private IP %s in account %s, region %s: %v\n", ipAddress, accountName, region, err)
        } else {
            results = append(results, privateResult)
        }

        publicResult, err := svc.DescribeInstances(context.TODO(), publicIPFilter)
        if err != nil {
            log.Printf("Error describing instances by public IP %s in account %s, region %s: %v\n", ipAddress, accountName, region, err)
        } else {
            results = append(results, publicResult)
        }
    }

    for _, result := range results {
        if len(result.Reservations) > 0 {
            fmt.Printf("Instances found in account: %s, Region: %s\n", accountName, region)
            for _, reservation := range result.Reservations {
                for _, instance := range reservation.Instances {
                    fmt.Printf("Instance ID: %s\n", *instance.InstanceId)
                    fmt.Printf("Instance Type: %s\n", instance.InstanceType)
                    fmt.Printf("Region: %s\n", region)
                    fmt.Printf("Private IP Address: %s\n", *instance.PrivateIpAddress)
                    if instance.PublicIpAddress != nil {
                        fmt.Printf("Public IP Address: %s\n", *instance.PublicIpAddress)
                    }
                    fmt.Println("-----")
					os.Exit(0)

                }
            }
        }
    }
}

func searchIAMUser(cfg aws.Config, accountName, searchAccessKey string) {
    svc := iam.NewFromConfig(cfg)
    
    // First, list all IAM users
    listUsersOutput, err := svc.ListUsers(context.TODO(), &iam.ListUsersInput{})
    if err != nil {
        log.Printf("Error listing IAM users in account %s: %v\n", accountName, err)
        return
    }

    for _, user := range listUsersOutput.Users {
        // For each user, list their access keys
        listAccessKeysOutput, err := svc.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
            UserName: user.UserName,
        })
        if err != nil {
            log.Printf("Error listing access keys for user %s in account %s: %v\n", *user.UserName, accountName, err)
            continue
        }

        // Check each access key for a match
        for _, accessKeyMetadata := range listAccessKeysOutput.AccessKeyMetadata {
            if *accessKeyMetadata.AccessKeyId == searchAccessKey {
                fmt.Printf("Found IAM user with access key %s in account %s\n", searchAccessKey, accountName)
                fmt.Printf("User Name: %s\n", *user.UserName)
                fmt.Printf("User ID: %s\n", *user.UserId)
                fmt.Printf("ARN: %s\n", *user.Arn)
				// fmt.Printf("Path : %s\n",*user.Path )
                fmt.Println("-----")
                return // Exit as soon as a match is found
            }
        }
    }

    fmt.Printf("No IAM user found with access key %s in account %s\n", searchAccessKey, accountName)
}

