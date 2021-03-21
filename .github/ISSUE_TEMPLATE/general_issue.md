---
name: Issue Template
about: Report a bug, feature, or enhancement
title: ''
labels: needs_triaging
assignees: davidji99

---

Hi there,

### Terraform Version
Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

### HerokuX Provider Version
Run `terraform -v` to show core and any provider versions. A sample output could be:

```
Terraform v0.14.8
+ provider.split v0.1.0
```

### Affected Resource(s)
Please list the resources as a list, for example:
- opc_instance
- opc_storage_volume

If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this.

### Terraform Configuration Files
```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

### Debug Output
Please provider a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist. Please MAKE SURE to mask any sensitive values.

### Panic Output
If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`.

### Expected Behavior
What should have happened?

### Actual Behavior
What actually happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### Important Factoids
Are there anything atypical about your accounts that we should know? For example: Running in EC2 Classic? Custom version of OpenStack? Tight ACLs?

### References
Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- GH-1234
