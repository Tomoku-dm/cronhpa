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
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 1800
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: test1
  minReplicas: 10
  maxReplicas: 20
  cron:
    - schdedul: 0 7 0 0 1-5
      timezone: ZH 
      minReplicas : 10
      maxReplicas : 20
    - schdedul: 0 19 0 0 1-5
      timezone: ZH
      minReplicas : 5
      maxReplicas : 20
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 60
```

# Step

mkdir cronhpa && cd cronhpa
go  mod init cronhpa
kubebuilder init --domain tomoku.com 



# Link
https://github.com/tkestack/cron-hpa

