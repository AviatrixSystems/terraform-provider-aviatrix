ARM_APPLICATION_ID=$(terraform output ARM_APPLICATION_ID)
ARM_APPLICATION_KEY=$(terraform output ARM_APPLICATION_KEY)
ARM_DIRECTORY_ID=$(terraform output ARM_DIRECTORY_ID)
ARM_GW_SIZE=$(terraform output ARM_GW_SIZE)
ARM_REGION=$(terraform output ARM_REGION)
ARM_REGION2=$(terraform output ARM_REGION2)
ARM_SUBNET=$(terraform output ARM_SUBNET)
ARM_SUBSCRIPTION_ID=$(terraform output ARM_SUBSCRIPTION_ID)
ARM_VNET_ID=$(terraform output ARM_VNET_ID)
ARM_VNET_ID2=$(terraform output ARM_VNET_ID2)
AVIATRIX_CONTROLLER_IP=$(terraform output AVIATRIX_CONTROLLER_IP)
AVIATRIX_PASSWORD=$(terraform output AVIATRIX_PASSWORD)
AVIATRIX_USERNAME=$(terraform output AVIATRIX_USERNAME)
AWS_ACCESS_KEY=$(terraform output AWS_ACCESS_KEY)
AWS_ACCOUNT_NUMBER=$(terraform output AWS_ACCOUNT_NUMBER)
AWS_BGP_VGW_ID=$(terraform output AWS_BGP_VGW_ID)
AWS_REGION=$(terraform output AWS_REGION)
AWS_REGION2=$(terraform output AWS_REGION2)
AWS_SECRET_KEY=$(terraform output AWS_SECRET_KEY)
AWS_SUBNET=$(terraform output AWS_SUBNET)
AWS_SUBNET2=$(terraform output AWS_SUBNET2)
AWS_VPC_ID=$(terraform output AWS_VPC_ID)
AWS_VPC_ID2=$(terraform output AWS_VPC_ID2)
AWS_DX_GATEWAY_ID=$(terraform output AWS_DX_GATEWAY_ID)
DOMAIN_NAME=$(terraform output DOMAIN_NAME)
AWSGOV_ACCESS_KEY=$(terraform output AWSGOV_ACCESS_KEY)
AWSGOV_SECRET_KEY=$(terraform output AWSGOV_SECRET_KEY)
AWSGOV_ACCOUNT_NUMBER=$(terraform output AWSGOV_ACCOUNT_NUMBER)
AWSGOV_REGION=$(terraform output AWSGOV_REGION)
AWSGOV_SUBNET=$(terraform output AWSGOV_SUBNET)
AWSGOV_VPC_ID=$(terraform output AWSGOV_VPC_ID)
GCP_CREDENTIALS_FILEPATH=$(terraform output GCP_CREDENTIALS_FILEPATH)
GCP_ID=$(terraform output GCP_ID)
GCP_SUBNET=$(terraform output GCP_SUBNET)
GCP_VPC_ID=$(terraform output GCP_VPC_ID)
GCP_ZONE=$(terraform output GCP_ZONE)
IDP_METADATA=$(terraform output IDP_METADATA)
IDP_METADATA_TYPE=$(terraform output IDP_METADATA_TYPE)
OCI_API_KEY_FILEPATH=$(terraform output OCI_API_KEY_FILEPATH)
OCI_COMPARTMENT_ID=$(terraform output OCI_COMPARTMENT_ID)
OCI_REGION=$(terraform output OCI_REGION)
OCI_SUBNET=$(terraform output OCI_SUBNET)
OCI_TENANCY_ID=$(terraform output OCI_TENANCY_ID)
OCI_USER_ID=$(terraform output OCI_USER_ID)
OCI_VPC_ID=$(terraform output OCI_VPC_ID)
controller_private_ip=$(terraform output controller_private_ip)

SKIP_DATA_ACCOUNT="no"
SKIP_DATA_CALLER_IDENTITY="no"
SKIP_DATA_FIRENET="no"
SKIP_DATA_FIRENET_VENDOR_INTEGRATION="no"
SKIP_DATA_GATEWAY="no"
SKIP_DATA_SPOKE_GATEWAY="no"
SKIP_DATA_SPOKE_GATEWAY_AWS="no"
SKIP_DATA_SPOKE_GATEWAY_ARM="no"
SKIP_DATA_SPOKE_GATEWAY_GCP="no"
SKIP_DATA_TRANSIT_GATEWAY="no"
SKIP_DATA_TRANSIT_GATEWAY_AWS="no"
SKIP_DATA_TRANSIT_GATEWAY_ARM="no"
SKIP_DATA_TRANSIT_GATEWAY_GCP="no"
SKIP_ACCOUNT="no"
SKIP_ACCOUNT_AWS="no"
SKIP_ACCOUNT_ARM="no"
SKIP_ACCOUNT_GCP="no"
SKIP_ACCOUNT_OCI="yes"
SKIP_ACCOUNT_AWSGOV="no"
SKIP_ACCOUNT_USER="no"
SKIP_ARM_PEER="no"
SKIP_AWS_PEER="no"
SKIP_AWS_TGW="no"
SKIP_AWS_TGW_DIRECTCONNECT="no"
SKIP_AWS_TGW_VPC_ATTACHMENT="no"
SKIP_AWS_TGW_VPN_CONN="no"
SKIP_CONTROLLER_CONFIG="no"
SKIP_FIRENET="no"
SKIP_FIREWALL="no"
SKIP_FIREWALL_INSTANCE="no"
SKIP_FIREWALL_TAG="no"
SKIP_FQDN="no"
SKIP_GATEWAY="no"
SKIP_GATEWAY_AWS="no"
SKIP_GATEWAY_GCP="no"
SKIP_GATEWAY_ARM="no"
SKIP_GATEWAY_OCI="yes"
SKIP_GATEWAY_AWSGOV="yes"
SKIP_GATEWAY_DNAT="no"
SKIP_GATEWAY_DNAT_AWS="no"
SKIP_GATEWAY_DNAT_ARM="no"
SKIP_GATEWAY_SNAT="no"
SKIP_GATEWAY_SNAT_AWS="no"
SKIP_GATEWAY_SNAT_ARM="no"
SKIP_GEO_VPN="no"
SKIP_SAML_ENDPOINT="no"
SKIP_S2C="no"
SKIP_SPOKE="yes"
SKIP_SPOKE_ARM="yes"
SKIP_SPOKE_AWS="yes"
SKIP_SPOKE_GCP="yes"
SKIP_SPOKE_GATEWAY="no"
SKIP_SPOKE_GATEWAY_ARM="no"
SKIP_SPOKE_GATEWAY_AWS="no"
SKIP_SPOKE_GATEWAY_GCP="no"
SKIP_SPOKE_GATEWAY_OCI="yes"
SKIP_TRANSIT="yes"
SKIP_TRANSIT_AWS="yes"
SKIP_TRANSIT_ARM="yes"
SKIP_TRANSIT_GATEWAY="no"
SKIP_TRANSIT_GATEWAY_AWS="no"
SKIP_TRANSIT_GATEWAY_ARM="no"
SKIP_TRANSIT_GATEWAY_GCP="no"
SKIP_TRANSIT_GATEWAY_OCI="yes"
SKIP_TRANSIT_GATEWAY_PEERING="no"
SKIP_TRANS_PEER="no"
SKIP_TUNNEL="no"
SKIP_VPC="no"
SKIP_VGW_CONN="no"
SKIP_VPN_PROFILE="no"
SKIP_VPN_USER="no"
SKIP_VPN_USER_ACCELERATOR="no"

echo "ARM_APPLICATION_ID=$ARM_APPLICATION_ID"
echo "ARM_APPLICATION_KEY=$ARM_APPLICATION_KEY"
echo "ARM_DIRECTORY_ID=$ARM_DIRECTORY_ID"
echo "ARM_GW_SIZE=$ARM_GW_SIZE"
echo "ARM_REGION=$ARM_REGION"
echo "ARM_REGION2=$ARM_REGION2"
echo "ARM_SUBNET=$ARM_SUBNET"
echo "ARM_SUBSCRIPTION_ID=$ARM_SUBSCRIPTION_ID"
echo "ARM_VNET_ID=$ARM_VNET_ID"
echo "ARM_VNET_ID2=$ARM_VNET_ID2"
echo "AVIATRIX_CONTROLLER_IP=$AVIATRIX_CONTROLLER_IP"
echo "AVIATRIX_PASSWORD=$AVIATRIX_PASSWORD"
echo "AVIATRIX_USERNAME=$AVIATRIX_USERNAME"
echo "AWS_ACCESS_KEY=$AWS_ACCESS_KEY"
echo "AWS_ACCOUNT_NUMBER=$AWS_ACCOUNT_NUMBER"
echo "AWS_BGP_VGW_ID=$AWS_BGP_VGW_ID"
echo "AWS_REGION=$AWS_REGION"
echo "AWS_REGION2=$AWS_REGION2"
echo "AWS_SECRET_KEY=$AWS_SECRET_KEY"
echo "AWS_SUBNET=$AWS_SUBNET"
echo "AWS_SUBNET2=$AWS_SUBNET2"
echo "AWS_VPC_ID=$AWS_VPC_ID"
echo "AWS_VPC_ID2=$AWS_VPC_ID2"
echo "AWS_DX_GATEWAY_ID=$AWS_DX_GATEWAY_ID"
echo "DOMAIN_NAME=$DOMAIN_NAME"
echo "AWSGOV_ACCESS_KEY=$AWSGOV_ACCESS_KEY"
echo "AWSGOV_ACCOUNT_NUMBER=$AWSGOV_ACCOUNT_NUMBER"
echo "AWSGOV_SECRET_KEY=$AWSGOV_SECRET_KEY"
echo "GCP_CREDENTIALS_FILEPATH=$GCP_CREDENTIALS_FILEPATH"
echo "GCP_ID=$GCP_ID"
echo "GCP_SUBNET=$GCP_SUBNET"
echo "GCP_VPC_ID=$GCP_VPC_ID"
echo "GCP_ZONE=$GCP_ZONE"
echo "IDP_METADATA=$IDP_METADATA"
echo "IDP_METADATA_TYPE=$IDP_METADATA_TYPE"
echo "OCI_API_KEY_FILEPATH=$OCI_API_KEY_FILEPATH"
echo "OCI_COMPARTMENT_ID=$OCI_COMPARTMENT_ID"
echo "OCI_REGION=$OCI_REGION"
echo "OCI_SUBNET=$OCI_SUBNET"
echo "OCI_TENANCY_ID=$OCI_TENANCY_ID"
echo "OCI_USER_ID=$OCI_USER_ID"
echo "OCI_VPC_ID=$OCI_VPC_ID"
echo "controller_private_ip=$controller_private_ip"

echo "SKIP_DATA_ACCOUNT=$SKIP_DATA_ACCOUNT"
echo "SKIP_DATA_CALLER_IDENTITY=$SKIP_DATA_CALLER_IDENTITY"
echo "SKIP_DATA_FIRENET=$SKIP_DATA_FIRENET"
echo "SKIP_DATA_FIRENET_VENDOR_INTEGRATION=$SKIP_DATA_FIRENET_VENDOR_INTEGRATION"
echo "SKIP_DATA_GATEWAY=$SKIP_DATA_GATEWAY"
echo "SKIP_DATA_SPOKE_GATEWAY=$SKIP_DATA_SPOKE_GATEWAY"
echo "SKIP_DATA_SPOKE_GATEWAY_AWS=$SKIP_DATA_SPOKE_GATEWAY_AWS"
echo "SKIP_DATA_SPOKE_GATEWAY_ARM=$SKIP_DATA_SPOKE_GATEWAY_ARM"
echo "SKIP_DATA_SPOKE_GATEWAY_GCP=$SKIP_DATA_SPOKE_GATEWAY_GCP"
echo "SKIP_DATA_TRANSIT_GATEWAY=$SKIP_DATA_TRANSIT_GATEWAY"
echo "SKIP_DATA_TRANSIT_GATEWAY_AWS=$SKIP_DATA_TRANSIT_GATEWAY_AWS"
echo "SKIP_DATA_TRANSIT_GATEWAY_ARM=$SKIP_DATA_TRANSIT_GATEWAY_ARM"
echo "SKIP_DATA_TRANSIT_GATEWAY_GCP=$SKIP_DATA_TRANSIT_GATEWAY_GCP"
echo "SKIP_ACCOUNT=$SKIP_ACCOUNT"
echo "SKIP_ACCOUNT_ARM=$SKIP_ACCOUNT_ARM"
echo "SKIP_ACCOUNT_AWS=$SKIP_ACCOUNT_AWS"
echo "SKIP_ACCOUNT_GCP=$SKIP_ACCOUNT_GCP"
echo "SKIP_ACCOUNT_OCI=$SKIP_ACCOUNT_OCI"
echo "SKIP_ACCOUNT_AWSGOV=$SKIP_ACCOUNT_AWSGOV"
echo "SKIP_ACCOUNT_USER=$SKIP_ACCOUNT_USER"
echo "SKIP_ARM_PEER=$SKIP_ARM_PEER"
echo "SKIP_AWS_PEER=$SKIP_AWS_PEER"
echo "SKIP_AWS_TGW=$SKIP_AWS_TGW"
echo "SKIP_AWS_TGW_DIRECTCONNECT=$SKIP_AWS_TGW_DIRECTCONNECT"
echo "SKIP_AWS_TGW_VPC_ATTACHMENT=$SKIP_AWS_TGW_VPC_ATTACHMENT"
echo "SKIP_AWS_TGW_VPN_CONN=$SKIP_AWS_TGW_VPN_CONN"
echo "SKIP_CONTROLLER_CONFIG=$SKIP_CONTROLLER_CONFIG"
echo "SKIP_FIRENET=$SKIP_FIRENET"
echo "SKIP_FIREWALL=$SKIP_FIREWALL"
echo "SKIP_FIREWALL_INSTANCE=$SKIP_FIREWALL_INSTANCE"
echo "SKIP_FIREWALL_TAG=$SKIP_FIREWALL_TAG"
echo "SKIP_FQDN=$SKIP_FQDN"
echo "SKIP_GATEWAY=$SKIP_GATEWAY"
echo "SKIP_GATEWAY_ARM=$SKIP_GATEWAY_ARM"
echo "SKIP_GATEWAY_AWS=$SKIP_GATEWAY_AWS"
echo "SKIP_GATEWAY_GCP=$SKIP_GATEWAY_GCP"
echo "SKIP_GATEWAY_OCI=$SKIP_GATEWAY_OCI"
echo "SKIP_GATEWAY_AWSGOV=$SKIP_GATEWAY_AWSGOV"
echo "SKIP_GATEWAY_DNAT=$SKIP_GATEWAY_DNAT"
echo "SKIP_GATEWAY_DNAT_ARM=$SKIP_GATEWAY_DNAT_ARM"
echo "SKIP_GATEWAY_DNAT_AWS=$SKIP_GATEWAY_DNAT_AWS"
echo "SKIP_GATEWAY_SNAT=$SKIP_GATEWAY_SNAT"
echo "SKIP_GATEWAY_SNAT_ARM=$SKIP_GATEWAY_SNAT_ARM"
echo "SKIP_GATEWAY_SNAT_AWS=$SKIP_GATEWAY_SNAT_AWS"
echo "SKIP_GEO_VPN=$SKIP_GEO_VPN"
echo "SKIP_SAML_ENDPOINT=$SKIP_SAML_ENDPOINT"
echo "SKIP_S2C=$SKIP_S2C"
echo "SKIP_SPOKE=$SKIP_SPOKE"
echo "SKIP_SPOKE_ARM=$SKIP_SPOKE_ARM"
echo "SKIP_SPOKE_AWS=$SKIP_SPOKE_AWS"
echo "SKIP_SPOKE_GCP=$SKIP_SPOKE_GCP"
echo "SKIP_SPOKE_GATEWAY=$SKIP_SPOKE_GATEWAY"
echo "SKIP_SPOKE_GATEWAY_ARM=$SKIP_SPOKE_GATEWAY_ARM"
echo "SKIP_SPOKE_GATEWAY_AWS=$SKIP_SPOKE_GATEWAY_AWS"
echo "SKIP_SPOKE_GATEWAY_GCP=$SKIP_SPOKE_GATEWAY_GCP"
echo "SKIP_SPOKE_GATEWAY_OCI=$SKIP_SPOKE_GATEWAY_OCI"
echo "SKIP_TRANSIT=$SKIP_TRANSIT"
echo "SKIP_TRANSIT_AWS=$SKIP_TRANSIT_AWS"
echo "SKIP_TRANSIT_ARM=$SKIP_TRANSIT_ARM"
echo "SKIP_TRANSIT_GATEWAY=$SKIP_TRANSIT_GATEWAY"
echo "SKIP_TRANSIT_GATEWAY_AWS=$SKIP_TRANSIT_GATEWAY_AWS"
echo "SKIP_TRANSIT_GATEWAY_ARM=$SKIP_TRANSIT_GATEWAY_ARM"
echo "SKIP_TRANSIT_GATEWAY_GCP=$SKIP_TRANSIT_GATEWAY_GCP"
echo "SKIP_TRANSIT_GATEWAY_OCI=$SKIP_TRANSIT_GATEWAY_OCI"
echo "SKIP_TRANSIT_GATEWAY_PEERING=$SKIP_TRANSIT_GATEWAY_PEERING"
echo "SKIP_TRANS_PEER=$SKIP_TRANS_PEER"
echo "SKIP_TUNNEL=$SKIP_TUNNEL"
echo "SKIP_VPC=$SKIP_VPC"
echo "SKIP_VGW_CONN=$SKIP_VGW_CONN"
echo "SKIP_VPN_PROFILE=$SKIP_VPN_PROFILE"
echo "SKIP_VPN_USER=$SKIP_VPN_USER"
echo "SKIP_VPN_USER_ACCELERATOR=$SKIP_VPN_USER_ACCELERATOR"

export ARM_APPLICATION_ID
export ARM_APPLICATION_KEY
export ARM_DIRECTORY_ID
export ARM_GW_SIZE
export ARM_REGION
export ARM_REGION2
export ARM_SUBNET
export ARM_SUBSCRIPTION_ID
export ARM_VNET_ID
export ARM_VNET_ID2
export AVIATRIX_CONTROLLER_IP
export AVIATRIX_PASSWORD
export AVIATRIX_USERNAME
export AWS_ACCESS_KEY
export AWS_ACCOUNT_NUMBER
export AWS_BGP_VGW_ID
export AWS_REGION
export AWS_REGION2
export AWS_SECRET_KEY
export AWS_SUBNET
export AWS_SUBNET2
export AWS_VPC_ID
export AWS_VPC_ID2
export AWS_DX_GATEWAY_ID
export DOMAIN_NAME
export AWSGOV_ACCESS_KEY
export AWSGOV_ACCOUNT_NUMBER
export AWSGOV_SECRET_KEY
export GCP_CREDENTIALS_FILEPATH
export GCP_ID
export GCP_SUBNET
export GCP_VPC_ID
export GCP_ZONE
export IDP_METADATA
export IDP_METADATA_TYPE
export OCI_API_KEY_FILEPATH
export OCI_COMPARTMENT_ID
export OCI_REGION
export OCI_SUBNET
export OCI_TENANCY_ID
export OCI_USER_ID
export OCI_VPC_ID
export controller_private_ip

export SKIP_DATA_ACCOUNT
export SKIP_DATA_CALLER_IDENTITY
export SKIP_DATA_FIRENET
export SKIP_DATA_FIRENET_VENDOR_INTEGRATION
export SKIP_DATA_GATEWAY
export SKIP_DATA_SPOKE_GATEWAY
export SKIP_DATA_SPOKE_GATEWAY_AWS
export SKIP_DATA_SPOKE_GATEWAY_ARM
export SKIP_DATA_SPOKE_GATEWAY_GCP
export SKIP_DATA_TRANSIT_GATEWAY
export SKIP_DATA_TRANSIT_GATEWAY_AWS
export SKIP_DATA_TRANSIT_GATEWAY_ARM
export SKIP_DATA_TRANSIT_GATEWAY_GCP
export SKIP_ACCOUNT
export SKIP_ACCOUNT_AWS
export SKIP_ACCOUNT_ARM
export SKIP_ACCOUNT_GCP
export SKIP_ACCOUNT_OCI
export SKIP_ACCOUNT_AWSGOV
export SKIP_ACCOUNT_USER
export SKIP_ARM_PEER
export SKIP_AWS_PEER
export SKIP_AWS_TGW
export SKIP_AWS_TGW_DIRECTCONNECT
export SKIP_AWS_TGW_VPC_ATTACHMENT
export SKIP_AWS_TGW_VPN_CONN
export SKIP_CONTROLLER_CONFIG
export SKIP_FIRENET
export SKIP_FIREWALL
export SKIP_FIREWALL_INSTANCE
export SKIP_FIREWALL_TAG
export SKIP_FQDN
export SKIP_GATEWAY
export SKIP_GATEWAY_AWS
export SKIP_GATEWAY_GCP
export SKIP_GATEWAY_ARM
export SKIP_GATEWAY_OCI
export SKIP_GATEWAY_DNAT
export SKIP_GATEWAY_DNAT_AWS
export SKIP_GATEWAY_DNAT_ARM
export SKIP_GATEWAY_SNAT
export SKIP_GATEWAY_SNAT_AWS
export SKIP_GATEWAY_SNAT_ARM
export SKIP_GATEWAY_AWSGOV
export SKIP_GEO_VPN
export SKIP_SAML_ENDPOINT
export SKIP_S2C
export SKIP_SPOKE
export SKIP_SPOKE_ARM
export SKIP_SPOKE_AWS
export SKIP_SPOKE_GCP
export SKIP_SPOKE_GATEWAY
export SKIP_SPOKE_GATEWAY_ARM
export SKIP_SPOKE_GATEWAY_AWS
export SKIP_SPOKE_GATEWAY_GCP
export SKIP_SPOKE_GATEWAY_OCI
export SKIP_TRANSIT
export SKIP_TRANSIT_AWS
export SKIP_TRANSIT_ARM
export SKIP_TRANSIT_GATEWAY
export SKIP_TRANSIT_GATEWAY_AWS
export SKIP_TRANSIT_GATEWAY_ARM
export SKIP_TRANSIT_GATEWAY_GCP
export SKIP_TRANSIT_GATEWAY_OCI
export SKIP_TRANSIT_GATEWAY_PEERING
export SKIP_TRANS_PEER
export SKIP_TUNNEL
export SKIP_VPC
export SKIP_VGW_CONN
export SKIP_VPN_PROFILE
export SKIP_VPN_USER
export SKIP_VPN_USER_ACCELERATOR
