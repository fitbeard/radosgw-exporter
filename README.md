# Ceph RadosGW Exporter

[Prometheus](https://prometheus.io/) exporter that scrapes
[Ceph](http://ceph.com/) Rados Gateway buckets and users usage information.
This information is gathered from a Rados Gateway using the
[Admin Operations API](http://docs.ceph.com/docs/master/radosgw/adminops/).

## Requirements

* Configure admin entry point (default is 'admin'):

```text
rgw admin entry = "admin"
```

* Enable admin API (default is enabled):

```text
rgw enable apis = "s3, admin"
```

* This exporter requires a user that has the following capability, see the Admin Guide
[here](http://docs.ceph.com/docs/master/radosgw/admin/#add-remove-admin-capabilities)
for more details.

```text
"caps": [
    {
        "type": "buckets",
        "perm": "read"
    },
    {
        "type": "metadata",
        "perm": "read"
    },
    {
        "type": "usage",
        "perm": "read"
    },
    {
        "type": "users",
        "perm": "read"
    }
]
```

```shell
radosgw-admin user create --display-name="RadosGW Exporter" --uid="radosgw-exporter"
radosgw-admin caps add --uid="radosgw-exporter" --caps="buckets=read; metadata=read; usage=read; users=read"
```

## Metrics

```shell
# HELP ceph_radosgw_bucket_num_objects Ceph RadosGW bucket num objects
# TYPE ceph_radosgw_bucket_num_objects gauge
ceph_radosgw_bucket_num_objects{bucket="rook-ceph-bucket-checker-1219a9fe-1652-408e-b850-d761134a1743",owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_bucket_num_shards Ceph RadosGW bucket num shards
# TYPE ceph_radosgw_bucket_num_shards gauge
ceph_radosgw_bucket_num_shards{bucket="rook-ceph-bucket-checker-1219a9fe-1652-408e-b850-d761134a1743",owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 11
# HELP ceph_radosgw_bucket_size_kb Ceph RadosGW bucket size kb
# TYPE ceph_radosgw_bucket_size_kb gauge
ceph_radosgw_bucket_size_kb{bucket="rook-ceph-bucket-checker-1219a9fe-1652-408e-b850-d761134a1743",owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_quota_bucket_enabled Ceph RadosGW bucket quota enabled
# TYPE ceph_radosgw_quota_bucket_enabled gauge
ceph_radosgw_quota_bucket_enabled{owner="dashboard-admin"} 0
ceph_radosgw_quota_bucket_enabled{owner="exporter"} 0
ceph_radosgw_quota_bucket_enabled{owner="rgw-admin-ops-user"} 0
ceph_radosgw_quota_bucket_enabled{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_quota_bucket_max_objects Ceph RadosGW bucket quota max objects
# TYPE ceph_radosgw_quota_bucket_max_objects gauge
ceph_radosgw_quota_bucket_max_objects{owner="dashboard-admin"} -1
ceph_radosgw_quota_bucket_max_objects{owner="exporter"} -1
ceph_radosgw_quota_bucket_max_objects{owner="rgw-admin-ops-user"} -1
ceph_radosgw_quota_bucket_max_objects{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} -1
# HELP ceph_radosgw_quota_bucket_max_size_kb Ceph RadosGW bucket quota max size kb
# TYPE ceph_radosgw_quota_bucket_max_size_kb gauge
ceph_radosgw_quota_bucket_max_size_kb{owner="dashboard-admin"} -1
ceph_radosgw_quota_bucket_max_size_kb{owner="exporter"} -1
ceph_radosgw_quota_bucket_max_size_kb{owner="rgw-admin-ops-user"} -1
ceph_radosgw_quota_bucket_max_size_kb{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} -1
# HELP ceph_radosgw_quota_user_enabled Ceph RadosGW user quota enabled
# TYPE ceph_radosgw_quota_user_enabled gauge
ceph_radosgw_quota_user_enabled{owner="dashboard-admin"} 0
ceph_radosgw_quota_user_enabled{owner="exporter"} 0
ceph_radosgw_quota_user_enabled{owner="rgw-admin-ops-user"} 0
ceph_radosgw_quota_user_enabled{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_quota_user_max_objects Ceph RadosGW user quota max objects
# TYPE ceph_radosgw_quota_user_max_objects gauge
ceph_radosgw_quota_user_max_objects{owner="dashboard-admin"} -1
ceph_radosgw_quota_user_max_objects{owner="exporter"} -1
ceph_radosgw_quota_user_max_objects{owner="rgw-admin-ops-user"} -1
ceph_radosgw_quota_user_max_objects{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} -1
# HELP ceph_radosgw_quota_user_max_size_kb Ceph RadosGW user quota max size kb
# TYPE ceph_radosgw_quota_user_max_size_kb gauge
ceph_radosgw_quota_user_max_size_kb{owner="dashboard-admin"} 0
ceph_radosgw_quota_user_max_size_kb{owner="exporter"} 0
ceph_radosgw_quota_user_max_size_kb{owner="rgw-admin-ops-user"} 0
ceph_radosgw_quota_user_max_size_kb{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_user_num_objects Ceph RadosGW user num objects
# TYPE ceph_radosgw_user_num_objects gauge
ceph_radosgw_user_num_objects{owner="dashboard-admin"} 0
ceph_radosgw_user_num_objects{owner="exporter"} 0
ceph_radosgw_user_num_objects{owner="rgw-admin-ops-user"} 0
ceph_radosgw_user_num_objects{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
# HELP ceph_radosgw_user_size_kb Ceph RadosGW user size kb
# TYPE ceph_radosgw_user_size_kb gauge
ceph_radosgw_user_size_kb{owner="dashboard-admin"} 0
ceph_radosgw_user_size_kb{owner="exporter"} 0
ceph_radosgw_user_size_kb{owner="rgw-admin-ops-user"} 0
ceph_radosgw_user_size_kb{owner="rook-ceph-internal-s3-user-checker-1219a9fe-1652-408e-b850-d761134a1743"} 0
```

## Config

| Environment variable | Description | Default |
| --- | --- | ---|
| `RADOSGW_ENDPOINT` | URL for the RadosGW admin API (example: [https://rgw.deepswamp:8088]) | `NA` |
| `ACCESS_KEY` | RadosGW user access key | `NA` |
| `SECRET_KEY` | RadosGW user secret key | `NA` |
| `EXPORTER_PORT` | Exporter listen port | `9242` |
| `SKIP_SSL_VERIFICATION` | Skip RadosGW endpoint SSL verification | `false` |

### Docker

Container images are available for `amd64` and `arm64` platforms:

```shell
quay.io/tadas/radosgw-exporter:latest
quay.io/tadas/radosgw-exporter:<TAG>
```

```shell
docker run -d -p 9242:9242 \
-e "RADOSGW_ENDPOINT=<host>" \
-e "ACCESS_KEY=<access_key>" \
-e "SECRET_KEY=<secret_key>" \
quay.io/tadas/radosgw-exporter:latest
```

Resulting metrics can be then retrieved via your Prometheus server via the
`http://<exporter host>:9242/metrics` endpoint.

## Kubernetes

You can find an example of deployment using [Rook](https://rook.io/) operator in a K8s environment
in `examples/k8s` directory.

## Grafana

You can find an example of Grafana dashboard in `examples/grafana` directory.
