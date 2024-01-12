package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"os"
	"strconv"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

// Buckets array
type Buckets []string

// BucketStats - bucket stats json structure
type BucketStats struct {
		Bucket    string `json:"bucket"`
		NumShards int64 `json:"num_shards"`
		ID        string `json:"id"`
		Marker    string `json:"marker"`
		Owner     string `json:"owner"`
		Ver       string `json:"ver"`
		Mtime     string `json:"mtime"`
		MaxMarker string `json:"max_marker"`
		Usage     struct {
			RgwMain struct {
				SizeKb       int64 `json:"size_kb"`
				SizeKbActual int64 `json:"size_kb_actual"`
				NumObjects   int64 `json:"num_objects"`
			} `json:"rgw.main"`
			RgwNone struct {
				SizeKb       int64 `json:"size_kb"`
				SizeKbActual int64 `json:"size_kb_actual"`
				NumObjects   int64 `json:"num_objects"`
			} `json:"rgw.none"`
			RgwMultimeta struct {
				SizeKb       int64 `json:"size_kb"`
				SizeKbActual int64 `json:"size_kb_actual"`
				NumObjects   int64 `json:"num_objects"`
			} `json:"rgw.multimeta"`
		} `json:"usage"`
		BucketQuota struct {
			Enabled    bool
			MaxSizeKb  int64 `json:"max_size_kb"`
			MaxObjects int64 `json:"max_objects"`
		} `json:"bucket_quota"`
}

// BucketQuotas - bucket quotas json structure
type BucketQuotas struct {
		Enabled    bool  `json:"enabled"`
		MaxSizeKb  int64 `json:"max_size_kb"`
		MaxObjects int64 `json:"max_objects"`
}

// Users array
type Users []string

// UserStats - user stats json structure
type UserStats struct {
		Stats struct {
			SizeKb       int64 `json:"size_kb"`
			SizeKbActual int64 `json:"size_kb_actual"`
			NumObjects   int64 `json:"num_objects"`
		} `json:"stats"`
}

// UserQuotas - user quotas json structure
type UserQuotas struct {
		Enabled    bool  `json:"enabled"`
		MaxSizeKb  int64 `json:"max_size_kb"`
		MaxObjects int64 `json:"max_objects"`
}

// ListBuckets returns buckets list array
func ListBuckets(endpoint string) (Buckets, error) {
    var b Buckets

    bjson := ListBucketsJSON(endpoint)
    if bjson == "" {
        return b, errors.New("empty JSON response")
    }

    if err := json.Unmarshal([]byte(bjson), &b); err != nil {
        return b, fmt.Errorf("failed to unmarshal JSON data for bucket list: %v", err)
    }

    return b, nil
}

// GetBucketStats return bucket stats
func GetBucketStats(endpoint string, bucket string) (BucketStats, error) {
	var bs BucketStats

	bsjson := GetBucketStatsJSON(endpoint, bucket)
    if bsjson == "" {
        return bs, errors.New("empty JSON response")
    }

	if err := json.Unmarshal([]byte(bsjson), &bs); err != nil {
		return bs, fmt.Errorf("failed to unmarshal JSON data for bucket stats: %v", err)
	}

	return bs, nil
}

// GetUserStats return bucket stats
func GetUserStats(endpoint string, user string) (UserStats, error) {
	var us UserStats

	usjson := GetUserStatsJSON(endpoint, user)
    if usjson == "" {
        return us, errors.New("empty JSON response")
    }

	if err := json.Unmarshal([]byte(usjson), &us); err != nil {
		return us, fmt.Errorf("failed to unmarshal JSON data for user stats: %v", err)
	}

	return us, nil
}

// ListBucketsJSON list all buckets in a zonegroup
func ListBucketsJSON(endpoint string) string {
		url := endpoint + "/admin/bucket"
		buckets := adminAPI(url)
		return buckets
}

// GetBucketStatsJSON return bucket stats in json
func GetBucketStatsJSON(endpoint string, bucket string) string {
		url := endpoint + "/admin/bucket?bucket=" + bucket
		bstats := adminAPI(url)
		return bstats
}

// GetUserStatsJSON return user stats in json
func GetUserStatsJSON(endpoint string, user string) string {
		url := endpoint + "/admin/user?uid=" + user  + "&stats=True"
		ustats := adminAPI(url)
		return ustats
}

// ListUsers returns users list array
func ListUsers(endpoint string) (Users, error) {
    var u Users
	ujson := ListUsersJSON(endpoint)
    if ujson == "" {
        return u, errors.New("empty JSON response")
    }
    if err := json.Unmarshal([]byte(ujson), &u); err != nil {
        return u, fmt.Errorf("failed to unmarshal JSON data for user list: %v", err)
    }

    return u, nil
}

// ListUsersJSON list all users in a zonegroup
func ListUsersJSON(endpoint string) string {
		url := endpoint + "/admin/metadata/user"
		users := adminAPI(url)
		return users
}

// GetUserQuotasJSON return user quotas
func GetUserQuotasJSON(endpoint string, user string) string {
		url := endpoint + "/admin/user?quota&quota-type=user&uid=" + user
		uquotas := adminAPI(url)
		return uquotas
}

// GetBucketQuotasJSON return bucket quotas
func GetBucketQuotasJSON(endpoint string, user string) string {
		url := endpoint + "/admin/user?quota&quota-type=bucket&uid=" + user
		bquotas := adminAPI(url)
		return bquotas
}

func createHTTPClient(skipSSLVerification bool) *http.Client {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: skipSSLVerification,
        },
    }

    return &http.Client{
        Timeout:   time.Second * 600,
        Transport: tr,
    }
}

// admin Api calls specific admin api url
func adminAPI(url string) string {

	accessKey := os.Getenv("ACCESS_KEY")
	os.Setenv("AWS_ACCESS_KEY_ID", accessKey)

	secretKey := os.Getenv("SECRET_KEY")
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)

	signer := v4.NewSigner(credentials.NewEnvCredentials())

	skipSSL := os.Getenv("SKIP_SSL_VERIFICATION")
    skipSSLVerification, _ := strconv.ParseBool(skipSSL)

	client := createHTTPClient(skipSSLVerification)

	req, _ := http.NewRequest("GET", url, nil)
	_, _ = signer.Sign(req, nil, "s3", "us-east-1", time.Now())

	var resp *http.Response
	var err error
	
	for {
		resp, err = client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Request timed out, retrying...")
			} else {
				log.Println("Network error:", err)
			}

			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatal("Failed even after retrying:", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Received non-OK status code: ", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
