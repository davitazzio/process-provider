#!/bin/bash

rm hello-process.sh 

kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: flights-service
  name: flights-service
spec:
  selector:
    matchLabels:
      run: flights-service
  template:
    metadata:
      labels:
        run: flights-service
    spec:
      containers:
      - image: liqo/flights-service
        name: flights-service
---
apiVersion: v1
kind: Service
metadata:
  labels:
    run: flights-service
  name: flights-service
spec:
  ports:
  - port: 7999
    targetPort: 7999
  selector:
    run: flights-service
EOF


while true; do
  echo "hello-world"
  sleep 1
done
