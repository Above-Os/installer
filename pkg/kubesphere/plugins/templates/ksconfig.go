package templates

import (
	"text/template"

	"github.com/lithammer/dedent"
)

var (
	KsConfigTempl = template.Must(template.New("KsConfig").Parse(
		dedent.Dedent(`apiVersion: v1
data:
  kubesphere.yaml: |
    authentication:
      authenticateRateLimiterMaxTries: 10
      authenticateRateLimiterDuration: 10m0s
      loginHistoryRetentionPeriod: 168h
      maximumClockSkew: 10s
      multipleLogin: True
      kubectlImage: kubesphere/kubectl:v1.22.0
      jwtSecret: "{{ .JwtSecret }}"
      oauthOptions:
        accessTokenMaxAge: 1209600000000000
        clients:
        - name: kubesphere
          secret: kubesphere
          redirectURIs:
          - '*'
    redis:
      host: redis.kubesphere-system.svc
      port: 6379
      password: KUBESPHERE_REDIS_PASSWORD
      db: 0
    network:
      enableNetworkPolicy: true
      ippoolType: none
    multicluster:
      clusterRole: none
    monitoring:
      endpoint: http://prometheus-operated.kubesphere-monitoring-system.svc:9090
      enableGPUMonitoring: false
    gpu:
      kinds:
      - resourceName: nvidia.com/gpu
        resourceType: GPU
        default: True
    notification:
      endpoint: http://notification-manager-svc.kubesphere-monitoring-system.svc:19093

    terminal:
      image: alpine:3.14
      timeout: 600

    gateway:
      watchesPath: /var/helm-charts/watches.yaml
      repository: kubesphere/nginx-ingress-controller
      tag: v1.1.0
      namespace: kubesphere-controls-system
kind: ConfigMap
metadata:
  name: kubesphere-config
  namespace: kubesphere-system
		`)))
)