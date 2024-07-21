
Kube Admission Webhook
===========



This [validating & mutating webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers) was developed to enforce Kubernetes policies across Licious Clusters.

Policies to be Supported by the Webhook:
* Disallow Public Repository on Image Pull URI.
* Mutate Resource Requests / Limit for specification not defined.
* Mutate Inject Custom SideCar / Annotations.
* Notify for any Validating policy failure on Slack.





