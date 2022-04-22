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

const coff float32 = 1.2

var (
	nodeResource = metav1.GroupVersionResource{Version: "v1", Resource: "nodes"}
)

func getPatchItem(op string, path string, val interface{}) patchOperation {
	return patchOperation{
		Op:    op,
		Path:  path,
		Value: val,
	}
}

func initPatch(node corev1.Node) []patchOperation {
	var patches []patchOperation
	nodeName := node.Name
	origin := float32(node.Status.Allocatable.Cpu().Value())
	log.Printf("%s: original cpu value is %f", nodeName, origin)
	fixed := origin * coff

	log.Printf("%s: changing cpu value to %f", nodeName, fixed)
	patches = append(patches, getPatchItem("replace", "/status/allocatable/cpu", fixed))
	return patches
}

func applyNodeConfig(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {
	if req.Resource != nodeResource {
		log.Printf("expect resource to be %s", nodeResource)
		return nil, nil
	}
	// Parse the Node object.
	raw := req.Object.Raw
	node := corev1.Node{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &node); err != nil {
		return nil, fmt.Errorf("could not deserialize node object: %v", err)
	}
	var patches []patchOperation
	patches = initPatch(node)
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
