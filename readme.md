# Argocd Repository Generator

This is an ArgoCD plugin generator that lets you template all repositories configured in ArgoCD by their name instead of their URL.

This simplifies ApplicationSet generators by allowing you to reference repositories by the name they have in ArgoCD, rather than hardcoding URLs.

For more information on ArgoCD plugin generators, please refer to the [official ArgoCD documentation](https://argo-cd.readthedocs.io/en/stable/user-guide/application-set/Generators-Plugin/).

## Installation

```bash
helm install argocd-repository-generator oci://ghcr.io/thorbenbelow/charts/argocd-repository-generator --version 0.1.0
```

## Example

Here is a basic example of how to use the generator. This `ApplicationSet` will find the repository named `example` in ArgoCD and use its URL in the template.

```yaml
apiVersion: v1
stringData:
  type: git
  url: https://github.com/argoproj/argocd-example-apps.git
  name: argocd-example-apps
kind: Secret
metadata:
  labels:
    argocd.argoproj.io/secret-type: repository
  name: example
  namespace: argocd
type: Opaque

---

apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: repositories-generator-example
  namespace: argocd
spec:
  goTemplate: true
  generators:
  - plugin:
      configMapRef:
        name: argocd-repository-generator # The name of the ConfigMap might differ based on the helm release name
  template:
    metadata:
      name: guestbook
      namespace: argocd
    spec:
      project: default
      source:
        repoURL: '{{ index .repositories "argocd-example-apps" }}'
        targetRevision: HEAD
        path: guestbook
      destination:
        server: 'https://kubernetes.default.svc'
        namespace: default
```
