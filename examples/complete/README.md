# Complete End-to-End Example

This example demonstrates a complete workflow using the Helmfile provider to
deploy multiple applications across different environments.

## Overview

This example shows:

- Provider configuration with plugin management
- Multiple helmfile releases with different configurations
- Environment-specific deployments
- Use of selectors and dependencies
- Integration with Terraform variables

## Files

- `main.tf` - Main Terraform configuration
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `helmfile.yaml` - Helmfile configuration
- `values/` - Environment-specific values files

## Architecture

```
┌─────────────────┐
│   Provider      │
│   Configuration │
└────────┬────────┘
         │
         ├─────────┬──────────┬──────────┐
         │         │          │          │
    ┌────▼───┐ ┌──▼────┐ ┌───▼────┐ ┌──▼─────┐
    │ Ingress│ │Database│ │Backend │ │Frontend│
    │ Release│ │Release │ │Release │ │Release │
    └────────┘ └────────┘ └────────┘ └────────┘
```

## Usage

1. **Initialize Terraform:**

   ```bash
   terraform init
   ```

2. **Review the plan:**

   ```bash
   terraform plan
   ```

3. **Apply the configuration:**

   ```bash
   terraform apply
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get pods -A
   helm list -A
   ```

## Customization

### Change Environment

Set the environment variable:

```bash
terraform apply -var="environment=staging"
```

### Override Values

Use variables to override default values:

```bash
terraform apply \
  -var="environment=production" \
  -var="replica_count=5" \
  -var="enable_monitoring=true"
```

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

## Notes

- This example assumes you have a Kubernetes cluster available
- Adjust the `helmfile.yaml` to match your actual chart requirements
- Review security settings before deploying to production
- Consider using remote state for team collaboration
