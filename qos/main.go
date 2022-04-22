package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `crt.pem`
	tlsKeyFile  = `key.pem`
)

const coff float64 = 0.8

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

func getPatchItem(op string, path string, val interface{}) patchOperation {
	return patchOperation{
		Op:    op,
		Path:  path,
		Value: val,
	}
}

func initPatch(pod corev1.Pod) []patchOperation {
	var patches []patchOperation
	if pod.Labels["oversale"] == "disabled" {
		return patches
	}
	podName := pod.Name
	for i, container := range pod.Spec.Containers {
		origin := container.Resources.Requests.Cpu().AsApproximateFloat64()
		log.Printf("%s, %d: original cpu value is %f", podName, i, origin)
		fixed := origin * coff

		log.Printf("%s,%d: changing cpu value to %f", podName, i, fixed)
		path := fmt.Sprintf("/spec/containers/%d/resources/requests/cpu", i)
		patches = append(patches, getPatchItem("replace", path, fixed))
	}
	return patches
}

func applyNodeConfig(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {
	if req.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		return nil, nil
	}
	// Parse the Node object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize node object: %v", err)
	}
	var patches []patchOperation
	patches = initPatch(pod)
 	log.Printf("testing pod struct podName %s, podRequestCPU %f, podLimitCPU %d", pod.Name, pod.Spec.Containers[0].Resources.Limits.Cpu().AsApproximateFloat64(), pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue() )
	return patches, nil
}
func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)
	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applyNodeConfig))
	log.Printf("listen on port 8443")
	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}
	server.ListenAndServeTLS(certPath, keyPath)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
