 curl https://172.16.61.175:443/mutate  --cacert ../../deployment/ca.crt -k -X POST -H "Content-Type: application/json"  -d '{ "apiVersion": "admission.k8s.io/v1", "kind": "AdmissionReview", "request": { } }'