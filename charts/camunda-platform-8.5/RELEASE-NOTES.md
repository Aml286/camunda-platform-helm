The changelog is automatically generated and it follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.
<!-- generated by git-cliff -->
### Release Info

Supported versions:

- Camunda applications: [8.5](https://github.com/camunda/camunda-platform/releases?q=tag%3A8.5&expanded=true)
- Helm values: [10.4.2](https://artifacthub.io/packages/helm/camunda/camunda-platform/10.4.2#parameters)
- Helm CLI: [3.16.1](https://github.com/helm/helm/releases/tag/v3.16.1)

Camunda images:

- docker.io/camunda/connectors-bundle:8.5.9
- docker.io/camunda/identity:8.5.6
- docker.io/camunda/operate:8.5.7
- docker.io/camunda/optimize:8.5.6
- docker.io/camunda/tasklist:8.5.8
- docker.io/camunda/zeebe:8.5.8
- registry.camunda.cloud/console/console-sm:8.5.112
- registry.camunda.cloud/web-modeler-ee/modeler-restapi:8.5.10
- registry.camunda.cloud/web-modeler-ee/modeler-webapp:8.5.10
- registry.camunda.cloud/web-modeler-ee/modeler-websockets:8.5.10

Non-Camunda images:

- docker.io/bitnami/elasticsearch:8.12.2
- docker.io/bitnami/keycloak:23.0.7
- docker.io/bitnami/os-shell:12-debian-12-r18
- docker.io/bitnami/postgresql:14.13.0
- docker.io/bitnami/postgresql:15.8.0

### Verification

To verify the integrity of the Helm chart using [Cosign](https://docs.sigstore.dev/signing/quickstart/):

```shell
cosign verify-blob camunda-platform-10.4.2.tgz \
  --bundle camunda-platform-10.4.2.cosign.bundle \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --certificate-identity "https://github.com/camunda/camunda-platform-helm/.github/workflows/chart-release-chores.yml@refs/pull/2454/merge"
```