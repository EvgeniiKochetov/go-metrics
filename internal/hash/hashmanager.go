package hash

import (
	"bytes"
	"crypto/sha256"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/config"
	"net/http"
)

func estimateHash(value []byte, key string) [32]byte {

	for _, v := range []byte(key) {
		value = append(value, v)
	}

	return sha256.Sum256(value)

}

func CheckHash(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		key := config.GetInstance().GetFlag("k")
		if key != "" {
			b := make([]byte, 0)
			r.Body.Read(b)
			hash := r.Header.Get("HashSHA256")
			estimatedHash := estimateHash(b, key)

			if !bytes.Equal(estimatedHash[:], []byte(hash)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		h.ServeHTTP(w, r)
		//if key != "" {
		//
		//}
	}
	return http.HandlerFunc(hashFn)
}
