Detailed README file for the `gview` CLI tool, which explains how to use each command:

---

# gview CLI Tool

`gview` is a command-line interface (CLI) tool for managing AWS accounts and searching for EC2 instances and IAM users across multiple AWS accounts and regions. It uses a YAML configuration file to store AWS credentials and region information for different accounts.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Usage](#usage)
   - [Adding an Account](#adding-an-account)
   - [Removing an Account](#removing-an-account)
   - [Searching for an EC2 Instance](#searching-for-an-ec2-instance)
   - [Searching for an IAM User](#searching-for-an-iam-user)
4. [Options](#options)
5. [Example Commands](#example-commands)
6. [Help](#help)

## Installation

To install the `gview` CLI tool:

1. Build the tool using Go:
   ```bash
   go build -o gview
   ```

2. Move the binary to a directory included in your system’s `PATH`, e.g.:
   ```bash
   sudo mv gview /usr/local/bin/
   ```

## Configuration

Before using the tool, you need to set up an `aws-accounts.yaml` file that will store your AWS account credentials and regions.

### YAML File Structure  (WIP)
The Goal is to make `aws-accounts.yaml` encrypted. This is work in progress for now. 

```yaml
accounts:
  - name: "account1"
    access_key: "AKIA..."
    secret_key: "your-secret-key"
    regions:
      - "us-east-1"
      - "us-west-2"
  - name: "account2"
    access_key: "AKIA..."
    secret_key: "your-secret-key"
    regions:
      - "eu-central-1"
      - "eu-west-1"
```

The YAML file should contain a list of accounts, each with a name, access key, secret key, and associated regions.

## Usage

### Adding an Account

To add a new AWS account to the configuration file:

```bash
gview -add --account-name <ACCOUNTNAME> --access-key <ACCESSKEYS> --secret-key <SECRETACCESSKEYS> --regions <CommaSeparatedRegions>
```

#### Example:

```bash
gview -add --account-name test-account --access-key AKDGVODFKC72NHBVPE --secret-key VS67NFMCWH2+rAoTCIXHBsVSsNOjxk --regions us-east-1,us-west-2
```
/gview 

This command will add an account named `test-account` with the specified AWS credentials and regions.

### Removing an Account

To remove an AWS account from the configuration file:

```bash
gview -remove --account-name <AccountName>
```

#### Example:

```bash
gview -remove --account-name test-account
```

This command will remove the account named `test-account` from the configuration file.

### Searching for an EC2 Instance

To search for an EC2 instance by its ID or IP address across multiple accounts and regions:

```bash
gview -find --instance-id <InstanceID>
```

or

```bash
gview -find --ip <IPAddress>
```

#### Example:

```bash
gview -find --instance-id i-0abcd1234efgh5678
```

This command will search for the EC2 instance with the ID `i-0abcd1234efgh5678` across all configured accounts and regions.

```bash
gview -find --ip 192.168.1.1
```

This command will search for the EC2 instance with the IP address `192.168.1.1` across all configured accounts and regions.

### Searching for an IAM User

To search for IAM users across multiple accounts and regions:

```bash
gview -find --iam <ACCESSKEYS>
```

#### Example:

```bash
gview -find --iam AKIAW4D65KWCGKY5XROS 
```

This command will list all IAM users across all configured accounts and regions.

## Options

-  `-add`	:			  Add an AWS account
-  `-remove`		:		  Remove an AWS account
-  `-find`		:	  Find an AWS account
   -  `--instance-id <id>` :     Instance ID to search for
   -  `--ip <ip-address>`   :    IP address to search for
   - `--iam`            :		  Search for IAM user associated with the access key
-  `--config <file>`1:         Path to YAML file with AWS credentials (default: aws-accounts.yaml)
-  `--account-name <name>`:   AWS account name to add or remove
-  `--access-key <key>`  :   AWS access key to add
-  `--secret-key <key>` :   AWS secret key to add
-  `--regions <regions>` :Comma-separated AWS regions to add
-  `--help`: Show help information.

## Example Commands

1. **Add an Account**:
   ```bash
   gview -add --account-name test-account --access-key AKDGVODFKC72NHBVPE --secret-key VS67NFMCWH2+rAoTCIXHBsVSsNOjxk --regions us-east-1,us-west-2
   ```

2. **Remove an Account**:
   ```bash
   gview -remove --account-name test-account
   ```

3. **Search for an EC2 Instance by ID**:
   ```bash
   gview -find --instance-id i-0abcd1234efgh5678
   ```

4. **Search for an EC2 Instance by IP Address**:
   ```bash
   gview -find --ip 192.168.1.1
   ```

5. **Search for IAM Users**:
   ```bash
   gview -find --iam AKIAW4D3453GKY5XROS 
   ```

## Help

To display the help information:

```bash
gview --help
```

This will show the available options and usage examples.

---

This README should provide all the necessary information to effectively use the `gview` CLI tool. If you have any further questions or need additional functionality, feel free to ask!