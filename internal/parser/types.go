package parser

const (
	Deployment                   = "Deployment"
	Service                      = "Service"
	ServiceAccount               = "ServiceAccount"
	Role                         = "Role"
	ClusterRole                  = "ClusterRole"
	ClusterRoleBinding           = "ClusterRoleBinding"
	RoleBinding                  = "RoleBinding"
	ConfigMap                    = "ConfigMap"
	MutatingWebhookConfiguration = "MutatingWebhookConfiguration"
	Certificate                  = "Certificate"
	Issuer                       = "Issuer"
)

const (
	DeploymentFile      = "deployment.yaml"
	ServiceFile         = "service.yaml"
	RbacFile            = "rbac.yaml"
	ServiceAccountFile  = "serviceaccount.yaml"
	ConfigMapFile       = "configmap.yaml"
	MutatingWebhookFile = "mutatingwebhook.yaml"
	CertificateFile     = "certificates.yaml"
)

var typeMap = map[string]string{
	Deployment:                   DeploymentFile,
	Service:                      ServiceFile,
	Role:                         RbacFile,
	ClusterRole:                  RbacFile,
	ServiceAccount:               ServiceAccountFile,
	RoleBinding:                  RbacFile,
	ClusterRoleBinding:           RbacFile,
	ConfigMap:                    ConfigMapFile,
	MutatingWebhookConfiguration: MutatingWebhookFile,
	Certificate:                  CertificateFile,
	Issuer:                       CertificateFile,
}

type data map[string]interface{}

func getType(t data) string {
	if v, ok := t["kind"]; ok {
		return v.(string)
	}
	return ""
}

func isType(t data, s string) bool {
	return s == getType(t)
}
