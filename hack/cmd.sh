 curl https://172.16.61.175:443/mutate  --cacert ../../deployment/ca.crt -k -X POST -H "Content-Type: application/json"  -d '{ "apiVersion": "admission.k8s.io/v1", "kind": "AdmissionReview", "request": {"uid": "705ab4f5-6393-11e8-b7cc-42010a800002", "kind": {"group": "", "version": "v1", "kind": "Node"}, "resource": {"version": "v1", "resource":"nodes"}, "operation": "UPDATE", "userInfo": {"username": "1234"} } }'

# echo "W3sib3AiOiJyZXBsYWNlIiwicGF0aCI6Ii9zdGF0dXMvYWxsb2NhdGFibGUvY3B1IiwidmFsdWUiOjB9XQ=="|base64 -d
