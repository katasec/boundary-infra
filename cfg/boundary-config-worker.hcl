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

listener "tcp" {
  address     = "0.0.0.0"
  purpose     = "proxy"
  #tls_disable = true

  tls_disable          = false
  tls_cert_file        = "/tls/boundary-worker.crt"
  tls_key_file         = "/tls/boundary.key"
}

worker {
  # Name attr must be unique
  public_addr  = "{{.Worker_PublicAddress}}"
  #name        = "env://HOSTNAME"
  name         = "{{.Worker_Name}}"
  description  = "A default worker created for demonstration"
  initial_upstreams  = ["{{.Worker_Controllers}}"]
  #controllers  = ["{{.Worker_Controllers}}"]
}

# Worker authorization KMS
# Using Azure Key Vault
kms "azurekeyvault" {
  purpose  = "worker-auth"
  key_name = "worker"
}