package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/fitbeard/radosgw-exporter/api"
)

type bucketStatsProm struct {
		Bucket       string
		Owner        string
		NumObjects   int64
		SizeKbActual int64
	    NumShards    int64
}

type bucketQuotasProm struct {
		Enabled    bool
		Owner      string
		MaxSizeKb  int64
		MaxObjects int64
}

type userStatsProm struct {
		Owner        string
		NumObjects   int64
		SizeKbActual int64
}

type userQuotasProm struct {
		Enabled    bool
		Owner      string
		MaxSizeKb  int64
		MaxObjects int64
}

// bucket num objects Gauge
var (
	numObjects = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_bucket_num_objects",
			Help: "Ceph RadosGW bucket num objects",
		},
		[]string{"bucket", "owner"},
	)
)

// bucket actual size in kb Gauge
var (
	sizeKbActual = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_bucket_size_kb",
			Help: "Ceph RadosGW bucket size kb",
		},
		[]string{"bucket", "owner"},
	)
)

// user num objects Gauge
var (
	userNumObjects = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_user_num_objects",
			Help: "Ceph RadosGW user num objects",
		},
		[]string{"owner"},
	)
)

// user actual size in kb Gauge
var (
	userSizeKbActual = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_user_size_kb",
			Help: "Ceph RadosGW user size kb",
		},
		[]string{"owner"},
	)
)

// bucket num shards
var (
	numShards = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_bucket_num_shards",
			Help: "Ceph RadosGW bucket num shards",
		},
		[]string{"bucket", "owner"},
	)
)

// user quota max size in kb Gauge
var (
	quotaUserMaxSizeKb = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_user_max_size_kb",
			Help: "Ceph RadosGW user quota max size kb",
		},
		[]string{"owner"},
	)
)

// user quota max objects Gauge
var (
	quotaUserMaxnumObjects = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_user_max_objects",
			Help: "Ceph RadosGW user quota max objects",
		},
		[]string{"owner"},
	)
)

// user quota enabled Gauge
var (
	quotaUserEnabled = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_user_enabled",
			Help: "Ceph RadosGW user quota enabled",
		},
		[]string{"owner"},
	)
)

// bucket quota enabled Gauge
var (
	quotaBucketEnabled = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_bucket_enabled",
			Help: "Ceph RadosGW bucket quota enabled",
		},
		[]string{"owner"},
	)
)

// bucket quota max size in kb Gauge
var (
	quotaBucketMaxSizeKb = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_bucket_max_size_kb",
			Help: "Ceph RadosGW bucket quota max size kb",
		},
		[]string{"owner"},
	)
)

// bucket quota max objects Gauge
var (
	quotaBucketMaxnumObjects = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ceph_radosgw_quota_bucket_max_objects",
			Help: "Ceph RadosGW bucket quota max objects",
		},
		[]string{"owner"},
	)
)

func boolToFloat(mybool bool) float64{
    if mybool {
        return 1
    }
    return 0
}

func main() {

	// go get rid of any additional metrics 
    // we have to expose our metrics with a custom registry
	registry := prometheus.NewRegistry()

	registry.MustRegister(numObjects)
	registry.MustRegister(sizeKbActual)
	registry.MustRegister(numShards)
	registry.MustRegister(userNumObjects)
	registry.MustRegister(userSizeKbActual)
	registry.MustRegister(quotaUserMaxSizeKb)
	registry.MustRegister(quotaUserMaxnumObjects)
	registry.MustRegister(quotaUserEnabled)
	registry.MustRegister(quotaBucketEnabled)
	registry.MustRegister(quotaBucketMaxSizeKb)
	registry.MustRegister(quotaBucketMaxnumObjects)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	endpoint := os.Getenv("RADOSGW_ENDPOINT")

	port := ":"+os.Getenv("EXPORTER_PORT")
	if port == ":" {
		port = ":9242"
	}

	if os.Getenv("RADOSGW_ENDPOINT") == "" {
		fmt.Println("RadosGW endpoint is not provided. Set environment variable RADOSGW_ENDPOINT.")
		os.Exit(1)
	}

	if endpoint := os.Getenv("RADOSGW_ENDPOINT"); endpoint != "" {
		if valid, _ := regexp.MatchString("https?://[a-zA-Z.0-9]+", strings.ToLower(endpoint)); valid {
		} else {
			fmt.Println("RadosGW endpoint URL must start with http:// or https://")
			os.Exit(1)
		}
	}

	if port := os.Getenv("EXPORTER_PORT"); port != "" {
		if _, err := strconv.Atoi(port); err == nil {
		} else {
			fmt.Println("Configured port is not a valid number:", port)
			os.Exit(1)
		}
	}

	if os.Getenv("ACCESS_KEY") == "" {
		fmt.Println("RadosGW credentials are not provided. Set environment variable ACCESS_KEY.")
		os.Exit(1)
	}

	if os.Getenv("SECRET_KEY") == "" {
		fmt.Println("RadosGW credentials are not provided. Set environment variable SECRET_KEY.")
		os.Exit(1)
	}

	var UsersJSON api.Users
	var BucketsJSON api.Buckets

	go func() {
		for {
			// get list of users from ceph admin api
			respu := api.ListUsersJSON(endpoint)
			if err := json.Unmarshal([]byte(respu), &UsersJSON); err != nil {
				fmt.Println("Error unmarshaling users:", err)
				continue
			}

			// get user quotas for each user
			uquotas := getUserQuotas(endpoint, UsersJSON)

			// get bucket quotas for each user
			bquotas := getBucketQuotas(endpoint, UsersJSON)

			// get buckets form ceph admin api
			respb := api.ListBucketsJSON(endpoint)
			if err := json.Unmarshal([]byte(respb), &BucketsJSON); err != nil {
				fmt.Println("Error unmarshaling buckets:", err)
				continue
			}

			// get buckets stats
			bstats := getBucketsStats(endpoint, BucketsJSON)

			// get users stats
			ustats := getUsersStats(endpoint, UsersJSON)

			// We need to reset counters before each update to clean up obsolete values
			quotaUserMaxSizeKb.Reset()
			quotaUserMaxnumObjects.Reset()
			quotaUserEnabled.Reset()
			quotaBucketEnabled.Reset()
			quotaBucketMaxSizeKb.Reset()
			quotaBucketMaxnumObjects.Reset()

			numObjects.Reset()
			sizeKbActual.Reset()
			numShards.Reset()
			userNumObjects.Reset()
			userSizeKbActual.Reset()


			for i := range uquotas {
				quotaUserMaxSizeKb.WithLabelValues(uquotas[i].Owner).Set(float64(uquotas[i].MaxSizeKb))
				quotaUserMaxnumObjects.WithLabelValues(uquotas[i].Owner).Set(float64(uquotas[i].MaxObjects))
				quotaUserEnabled.WithLabelValues(uquotas[i].Owner).Set(boolToFloat(uquotas[i].Enabled))
			}
			for i := range bquotas {
				if bquotas[i].MaxSizeKb == 0 {
					quotaBucketMaxSizeKb.WithLabelValues(bquotas[i].Owner).Set(float64(-1))
				} else {
					quotaBucketMaxSizeKb.WithLabelValues(bquotas[i].Owner).Set(float64(bquotas[i].MaxSizeKb))
				}
				quotaBucketMaxnumObjects.WithLabelValues(bquotas[i].Owner).Set(float64(bquotas[i].MaxObjects))
				quotaBucketEnabled.WithLabelValues(bquotas[i].Owner).Set(boolToFloat(bquotas[i].Enabled))
			}
			for i := range bstats {
				numObjects.WithLabelValues(bstats[i].Bucket, bstats[i].Owner).Set(float64(bstats[i].NumObjects))
				sizeKbActual.WithLabelValues(bstats[i].Bucket, bstats[i].Owner).Set(float64(bstats[i].SizeKbActual))
				numShards.WithLabelValues(bstats[i].Bucket, bstats[i].Owner).Set(float64(bstats[i].NumShards))
			}
			for i := range ustats {
				userNumObjects.WithLabelValues(ustats[i].Owner).Set(float64(ustats[i].NumObjects))
				userSizeKbActual.WithLabelValues(ustats[i].Owner).Set(float64(ustats[i].SizeKbActual))
			}
			time.Sleep(2 * time.Minute)
		}
	}()


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
             <head><title>Ceph RadosGW Exporter</title></head>
             <body>
             <h1>Ceph RadosGW Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
		if err != nil {
			log.Fatalln(err)
		}
	})

	// metrics endpoint
	http.Handle("/metrics", handler)

	log.Println("Starting HTTP server on", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

// getBucketsStats call ceph admin api for each bucket stats and return an array
// containing info for prometheus exporter
func getBucketsStats(endpoint string, BucketsJSON api.Buckets) []bucketStatsProm {
	bucketsstats := make([]bucketStatsProm, len(BucketsJSON))
	for i := range BucketsJSON {
		var bsprom bucketStatsProm
		var bstats api.BucketStats

		// get stats for each bucket
		respbs := api.GetBucketStatsJSON(endpoint, BucketsJSON[i])
		if err := json.Unmarshal([]byte(respbs), &bstats); err != nil {
			fmt.Println("Error unmarshaling bucket stats:", err)
			continue
		}

		bsprom.Bucket = BucketsJSON[i]
		bsprom.Owner = bstats.Owner
		bsprom.NumShards = bstats.NumShards
		bsprom.NumObjects = bstats.Usage.RgwMain.NumObjects
		bsprom.SizeKbActual = bstats.Usage.RgwMain.SizeKbActual
		bucketsstats[i] = bsprom
	}
	return bucketsstats
}

// getUsersStats call ceph admin api for each user stats and return an array
// containing info for prometheus exporter
func getUsersStats(endpoint string, UsersJSON api.Users) []userStatsProm {
	usersstats := make([]userStatsProm, len(UsersJSON))
	for i := range UsersJSON {
		var usprom userStatsProm
		var ustats api.UserStats

		// get stats for each user
		respbs := api.GetUserStatsJSON(endpoint, UsersJSON[i])
		if err := json.Unmarshal([]byte(respbs), &ustats); err != nil {
			fmt.Println("Error unmarshaling user stats:", err)
			continue
		}

		usprom.Owner = UsersJSON[i]
		usprom.NumObjects = ustats.Stats.NumObjects
		usprom.SizeKbActual = ustats.Stats.SizeKbActual
		usersstats[i] = usprom
	}
	return usersstats
}

// getUserQuotas call ceph admin api for each user quota stats and return an array
// containing info for prometheus exporter
func getUserQuotas(endpoint string, UsersJSON api.Users) []userQuotasProm {
	userquotas := make([]userQuotasProm, len(UsersJSON))
	for i := range UsersJSON {
		var uqprom userQuotasProm
		var uquotas api.UserQuotas

		// get quota for each user
		respuq := api.GetUserQuotasJSON(endpoint, UsersJSON[i])
		if err := json.Unmarshal([]byte(respuq), &uquotas); err != nil {
			fmt.Println("Error unmarshaling user quotas:", err)
			continue
		}

		uqprom.Enabled = uquotas.Enabled
		uqprom.Owner = UsersJSON[i]
		uqprom.MaxSizeKb = uquotas.MaxSizeKb
		uqprom.MaxObjects = uquotas.MaxObjects
		userquotas[i] = uqprom
	}
	return userquotas
}

// getBucketQuotas call ceph admin api for each user bucket quota stats and return an array
// containing info for prometheus exporter
func getBucketQuotas(endpoint string, UsersJSON api.Users) []bucketQuotasProm {
	bucketquotas := make([]bucketQuotasProm, len(UsersJSON))
	for i := range UsersJSON {
		var bqprom bucketQuotasProm
		var bquotas api.BucketQuotas

		// get quota for each bucket
		respuq := api.GetBucketQuotasJSON(endpoint, UsersJSON[i])
		if err := json.Unmarshal([]byte(respuq), &bquotas); err != nil {
			fmt.Println("Error unmarshaling bucket quotas:", err)
			continue
		}

		bqprom.Enabled = bquotas.Enabled
		bqprom.Owner = UsersJSON[i]
		bqprom.MaxSizeKb = bquotas.MaxSizeKb
		bqprom.MaxObjects = bquotas.MaxObjects
		bucketquotas[i] = bqprom
	}
	return bucketquotas
}
