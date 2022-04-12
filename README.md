# Terraform IAM generator
Generate the needed IAM Policy for your terraform code to run!

## Usage
```
$ terraform-iam-generator -help
  -dir string
    	terraform directory to use (default "current working directory")
  -tf-var value
    		Terraform variables to use. Specify like this: "KEY=VALUE". Can be used multiple times
  -tf-var-file value
    	Path to a terraform variables file. Must be relative to the passed directory. Can be used multiple times
```

## How are the policies generated?
AWS publishes metrics when calling their API via the SDKs. It's called **Client Side Monitoring**, short **CSM**. Terraform IAM Generator spins up a UDP server while running the terraform apply and destroy commands to get the information of the api calls. Unfortunatly, the requested resources are not included in these metrics.