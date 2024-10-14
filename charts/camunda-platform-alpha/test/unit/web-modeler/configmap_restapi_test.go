package web_modeler

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

type configmapRestAPITemplateTest struct {
	suite.Suite
	chartPath string
	release   string
	namespace string
	templates []string
}

func TestRestAPIConfigmapTemplate(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../../../")
	require.NoError(t, err)

	suite.Run(t, &configmapRestAPITemplateTest{
		chartPath: chartPath,
		release:   "camunda-platform-test",
		namespace: "camunda-platform-" + strings.ToLower(random.UniqueId()),
		templates: []string{"templates/web-modeler/configmap-restapi.yaml"},
	})
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectAuthClientApiAudience() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                                "true",
			"webModeler.restapi.mail.fromAddress":               "example@example.com",
			"global.identity.auth.webModeler.clientApiAudience": "custom-audience",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("custom-audience", configmapApplication.Camunda.Modeler.Security.JWT.Audience.InternalAPI)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectAuthPublicApiAudience() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                                "true",
			"webModeler.restapi.mail.fromAddress":               "example@example.com",
			"global.identity.auth.webModeler.publicApiAudience": "custom-audience",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("custom-audience", configmapApplication.Camunda.Modeler.Security.JWT.Audience.PublicAPI)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectIdentityServiceUrlWithFullnameOverride() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                  "true",
			"webModeler.restapi.mail.fromAddress": "example@example.com",
			"identity.fullnameOverride":           "custom-identity-fullname",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("http://custom-identity-fullname:80", configmapApplication.Camunda.Identity.BaseURL)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectIdentityServiceUrlWithNameOverride() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                  "true",
			"webModeler.restapi.mail.fromAddress": "example@example.com",
			"identity.nameOverride":               "custom-identity",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("http://camunda-platform-test-custom-identity:80", configmapApplication.Camunda.Identity.BaseURL)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectIdentityType() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                    "true",
			"webModeler.restapi.mail.fromAddress":   "example@example.com",
			"global.identity.auth.type":             "MICROSOFT",
			"global.identity.auth.issuerBackendUrl": "https://example.com",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("MICROSOFT", configmapApplication.Camunda.Identity.Type)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectKeycloakServiceUrl() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                    "true",
			"webModeler.restapi.mail.fromAddress":   "example@example.com",
			"global.identity.keycloak.url.protocol": "http",
			"global.identity.keycloak.url.host":     "keycloak",
			"global.identity.keycloak.url.port":     "80",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("http://keycloak:80/auth/realms/camunda-platform", configmapApplication.Camunda.Modeler.Security.JWT.Issuer.BackendUrl)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetCorrectKeycloakServiceUrlWithCustomPort() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                    "true",
			"webModeler.restapi.mail.fromAddress":   "example@example.com",
			"global.identity.keycloak.url.protocol": "http",
			"global.identity.keycloak.url.host":     "keycloak",
			"global.identity.keycloak.url.port":     "8888",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("http://keycloak:8888/auth/realms/camunda-platform", configmapApplication.Camunda.Modeler.Security.JWT.Issuer.BackendUrl)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetSmtpCredentials() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                   "true",
			"webModeler.restapi.mail.fromAddress":  "example@example.com",
			"webModeler.restapi.mail.smtpUser":     "modeler-user",
			"webModeler.restapi.mail.smtpPassword": "modeler-password",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("modeler-user", configmapApplication.Spring.Mail.Username)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetExternalDatabaseConfiguration() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                           "true",
			"webModeler.restapi.mail.fromAddress":          "example@example.com",
			"postgresql.enabled":                           "false",
			"webModeler.restapi.externalDatabase.url":      "jdbc:postgresql://postgres.example.com:65432/modeler-database",
			"webModeler.restapi.externalDatabase.user":     "modeler-user",
			"webModeler.restapi.externalDatabase.password": "modeler-password",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal("jdbc:postgresql://postgres.example.com:65432/modeler-database", configmapApplication.Spring.Datasource.Url)
	s.Require().Equal("modeler-user", configmapApplication.Spring.Datasource.Username)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetClusterConfigurationWithDefaultValuesIfZeebeEnabled() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                  "true",
			"webModeler.restapi.mail.fromAddress": "example@example.com",
			"postgresql.enabled":                  "false",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal(1, len(configmapApplication.Camunda.Modeler.Clusters))
	s.Require().Equal("camunda-platform-test-zeebe", configmapApplication.Camunda.Modeler.Clusters[0].Name)
	s.Require().Equal("8.6.0", configmapApplication.Camunda.Modeler.Clusters[0].Version)
	s.Require().Equal("OAUTH", configmapApplication.Camunda.Modeler.Clusters[0].Authentication)
	s.Require().Equal("grpc://camunda-platform-test-zeebe-gateway:26500", configmapApplication.Camunda.Modeler.Clusters[0].Url.Zeebe.Grpc)
	s.Require().Equal("http://camunda-platform-test-zeebe-gateway:8080", configmapApplication.Camunda.Modeler.Clusters[0].Url.Zeebe.Rest)
	s.Require().Equal("http://camunda-platform-test-operate:80", configmapApplication.Camunda.Modeler.Clusters[0].Url.Operate)
	s.Require().Equal("http://camunda-platform-test-tasklist:80", configmapApplication.Camunda.Modeler.Clusters[0].Url.Tasklist)
	s.Require().Equal("http://camunda-platform-test-keycloak:80/auth/realms/camunda-platform/protocol/openid-connect/token", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Url)
	s.Require().Equal("zeebe-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Zeebe)
	s.Require().Equal("operate-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Operate)
	s.Require().Equal("tasklist-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Tasklist)
	s.Require().Equal("", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Scope)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldSetClusterConfigurationWithCustomValuesIfZeebeEnabled() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                     "true",
			"webModeler.restapi.mail.fromAddress":    "example@example.com",
			"postgresql.enabled":                     "false",
			"global.zeebeClusterName":                "test-cluster",
			"global.identity.auth.zeebe.tokenScope":  "test-scope",
			"global.identity.auth.zeebe.audience":    "test-zeebe-api",
			"global.identity.auth.operate.audience":  "test-operate-api",
			"global.identity.auth.tasklist.audience": "test-tasklist-api",
			"global.identity.auth.tokenUrl":          "https://example.com/auth/realms/test/protocol/openid-connect/token",
			"zeebe.image.tag":                        "8.7.0-alpha1",
			"zeebeGateway.contextPath":               "/zeebe",
			"zeebeGateway.service.grpcPort":          "26600",
			"zeebeGateway.service.restPort":          "8090",
			"operate.contextPath":                    "/operate",
			"operate.service.port":                   "8080",
			"tasklist.contextPath":                   "/tasklist",
			"tasklist.service.port":                  "8080",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Equal(1, len(configmapApplication.Camunda.Modeler.Clusters))
	s.Require().Equal("test-cluster", configmapApplication.Camunda.Modeler.Clusters[0].Name)
	s.Require().Equal("8.7.0-alpha1", configmapApplication.Camunda.Modeler.Clusters[0].Version)
	s.Require().Equal("OAUTH", configmapApplication.Camunda.Modeler.Clusters[0].Authentication)
	s.Require().Equal("grpc://test-cluster-gateway:26600", configmapApplication.Camunda.Modeler.Clusters[0].Url.Zeebe.Grpc)
	s.Require().Equal("http://test-cluster-gateway:8090/zeebe", configmapApplication.Camunda.Modeler.Clusters[0].Url.Zeebe.Rest)
	s.Require().Equal("http://camunda-platform-test-operate:8080/operate", configmapApplication.Camunda.Modeler.Clusters[0].Url.Operate)
	s.Require().Equal("http://camunda-platform-test-tasklist:8080/tasklist", configmapApplication.Camunda.Modeler.Clusters[0].Url.Tasklist)
	s.Require().Equal("https://example.com/auth/realms/test/protocol/openid-connect/token", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Url)
	s.Require().Equal("test-zeebe-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Zeebe)
	s.Require().Equal("test-operate-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Operate)
	s.Require().Equal("test-tasklist-api", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Audience.Tasklist)
	s.Require().Equal("test-scope", configmapApplication.Camunda.Modeler.Clusters[0].Oauth.Scope)
}

func (s *configmapRestAPITemplateTest) TestContainerShouldNotSetClusterConfigurationIfZeebeDisabled() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"webModeler.enabled":                  "true",
			"webModeler.restapi.mail.fromAddress": "example@example.com",
			"postgresql.enabled":                  "false",
			"zeebe.enabled":                       "false",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var configmap corev1.ConfigMap
	var configmapApplication WebModelerRestAPIApplicationYAML
	helm.UnmarshalK8SYaml(s.T(), output, &configmap)

	err := yaml.Unmarshal([]byte(configmap.Data["application.yaml"]), &configmapApplication)
	if err != nil {
		s.Fail("Failed to unmarshal yaml. error=", err)
	}

	// then
	s.Require().Empty(configmapApplication.Camunda.Modeler.Clusters)
}
