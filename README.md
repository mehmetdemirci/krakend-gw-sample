# krakend-gw-sample

This is a sample krakend api gateway that shows how to use KrakenD Flexible configuration and writing custom plugin.

###Running

**Build docker image:**
```sh
docker build . -t krakend-gw:1.0.0
```

**Deploy to k8s:**
```sh
kubectl apply -f deployment.yaml

kubectl apply -f service.yaml
```