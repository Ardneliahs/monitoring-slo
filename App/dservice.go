package dservice

import (
	"net/http"
	"fmt"
	"encoding/json"
	"math/rand/v2"
	"time"
)

var startTime = time.Now()
func main() {
	http.HandleFunc("/health",healthHandler)
	http.HandleFunc("/work",workHandler)
	// nil means the default servemux will be used, we can also create our own.
	if err := http.ListenAndServe(":8080", nil); err != nil { //instanialisation of TCP connection is abstracted inside listen and serve
		fmt.Println("Server error:", err)
    }
}

func healthHandler(w http.ResponseWriter, r *http.Request){
	type HealthResponse struct {
		Status string `json:"status"`
		UptimeSec int `json:"uptime_sec"`
		Version string `json:"version"`
	}
	healthStatus:= HealthResponse{
		Status: "UP",
		UptimeSec: int(time.Since(startTime).Seconds()),
		Version: "1.0.0",
	}
	healthJson, _ := json.Marshal(healthStatus)
	w.Header().Set("Content-Type","application/json")
	w.Write(healthJson)
}

func workHandler(w http.ResponseWriter, r *http.Request){ //has varying latency, sometimes timeout, sometimes sends 5xx
	reqNumber := rand.IntN(100)
	if reqNumber == 1 {
		w.WriteHeader(http.StatusInternalServerError)
	} else if reqNumber == 2 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		time.Sleep(time.Duration(reqNumber) * time.Millisecond)
		w.Write([]byte("done!")) // because string not a single byte but a slice of bytes
	}
}