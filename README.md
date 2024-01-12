# Ceph RadosGW Usage Exporter

[Prometheus](https://prometheus.io/) exporter that scrapes
[Ceph](http://ceph.com/) Rados Gateway buckets and users usage information.
This information is gathered from a Rados Gateway using the
[Admin Operations API](http://docs.ceph.com/docs/master/radosgw/adminops/).

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
