package qos

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
	tlsDir         = `/run/secrets/tls`
	tlsCertFile    = `crt.pem`
	tlsKeyFile     = `key.pem`
)

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
	if pod.Labels["oversale.huifu.com"] == "false" {
		return patches
	}
	podName := pod.Name

	for _, container := range pod.Spec.Containers {
		origin := container.Resources.Requests.Cpu().Value()
		log.Printf("%s: original cpu value is %f", podName, origin)
		coff := float64(0.8)
		fixed := float64(origin)*coff

		log.Printf("%s: changing cpu value to %f", podName, fixed)
		patches = append(patches, getPatchItem("replace", "/spec/containers/resources/requests/cpu", fixed))
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
