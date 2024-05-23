# Configuration

## SeaweedFS Master Node

[http://54.91.129.148:9333/](http://54.91.129.148:9333/ "http://54.91.129.148:9333/")

| Data Center | Location           |
| ----------- | ------------------ |
| dc1         | N. Virginia        |
| dc2         | California         |
| dc3         | Sao Paolo (Brazil) |
| dc4         | Paris (France)     |

### Master

```bash
weed master -port=9333 -mdir=/data/seaweedfs/master0 -ip="54.91.129.148" -ip.bind="172.31.28.30"
```

### Volume

```bash
weed volume -dataCenter=dc3 -rack=v1 -mserver="54.91.129.148:9333" -port=8080 -ip="18.230.226.109" -ip.bind="172.31.2.207" -preStopSeconds=1 --dir=/data/seaweedfs
```

-   [ ]
