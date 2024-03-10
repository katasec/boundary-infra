# --------------------------------------------------------
# Azure Creds will be passed via environment variables:
# --------------------------------------------------------
#
# AZURE_TENANT_ID:                  Azure Tenant ID
# AZURE_CLIENT_ID:                  Azure App ID 
# AZURE_CLIENT_SECRET:              Azure App Password
# AZUREKEYVAULT_WRAPPER_VAULT_NAME: Key Vault Name
# BOUNDARY_POSTGRES_URL:            Postgres connection string


disable_mlock = true


controller {
  #name        = "env://HOSTNAME"
  name        = "boundary-controller.nprod.corp.internal"
  description = "A controller for a demo!"
  database {
    url = "env://BOUNDARY_POSTGRES_URL"
  }
  public_cluster_addr = "{{.Controller_PublicClusterAddress}}"
}

# API config
listener "tcp" {
  purpose              = "api"
  tls_disable          = false
  tls_cert_file        = "/tls/boundary-controller.crt"
  tls_key_file         = "/tls/boundary.key"
  cors_enabled         = true
  cors_allowed_origins = ["*"]
  address              = "0.0.0.0"
}

# Cluster config
listener "tcp" {
  purpose     = "cluster"
  tls_disable          = false
  tls_cert_file        = "/tls/boundary-controller.crt"
  tls_key_file         = "/tls/boundary.key"
  address              = "0.0.0.0"
}

# Root KMS configuration block: this is the root key for Boundary
# Using Azure Key Vault
kms "azurekeyvault" {
  purpose  = "root"
  key_name = "root"
}

# Worker authorization KMS
# Using Azure Key Vault
kms "azurekeyvault" {
  purpose  = "worker-auth"
  key_name = "worker"
}

# Recovery KMS block: configures the recovery key for Boundary
# Using Azure Key Vault
kms "azurekeyvault" {
  purpose  = "recovery"
  key_name = "recovery"

}
