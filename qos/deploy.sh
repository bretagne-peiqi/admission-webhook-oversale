kubectl delete MutatingWebhookConfiguration qos 
kubectl delete deployment qos-server -n webhook

basedir=.
ca_pem_b64="$(openssl base64 -A <"${basedir}/ca.crt")"
sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"${basedir}/deployment.yaml" \
    | kubectl apply -f -

