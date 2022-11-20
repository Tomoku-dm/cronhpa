# cronhpa
This k8s operator project,support timezone HPA with cronjob time

# CRD
```
apiVersion: autoscaling.tomoku.com/v1beta1
kind: CronHPA
metadata:
  name: cronhpa-test1
  namespace: default
spec:
  template:
    spec:
      behavior:
        scaleDown:
          stabilizationWindowSeconds: 1800
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: test1
      minReplicas: 10
      maxReplicas: 20
      metrics:
      - type: Resource
        resource:
          name: cpu
          target:
            type: Utilization
            averageUtilization: 60
  cron:
  - name: "daytime"
    schdedul: "0 7 0 0 1-5"
    timezone: "Asia/China"
    minReplicas : 10
    maxReplicas : 20
  - name: "nighttime"
    timezone: "Asia/China"
    "schdedul: "0 19 0 0 1-5"
    minReplicas : 5
    maxReplicas : 20
```

# Step

mkdir cronhpa && cd cronhpa
go  mod init cronhpa
kubebuilder init --domain tomoku.com 
kubebuilder create api --group cronhpa --version v1 --kind CronHPA


# Link
https://github.com/tkestack/cron-hpa
https://sqbu-github.cisco.com/WebexPlatform/aws-iam-controller/commit/862fb638b57b4bacb6831ab669291e3c224d7711
https://github.com/dtaniwaki/cron-hpa