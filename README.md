# 🚀 Coupon Service API

This is a simple coupon service API built with Go. It allows you to create, apply, and retrieve coupons.

## 🐳 Running with Docker

To run the service with Docker, follow these steps:

1. Build the Docker image:

```sh
docker build -t coupon-service .
```

2. Run the Docker container:

```sh
docker run -p 8080:8080 coupon-service
```

The service will be available at `http://localhost:8080`.

## 🚀 Deploying with Kubernetes

To deploy the service with Kubernetes, follow these steps:

1. Apply the Deployment configuration:

```sh
kubectl apply -f k8s/Deployment.yaml
```

2. Apply the HPA configuration:

```sh
kubectl apply -f k8s/HPA.yaml
```

3. Apply the Service configuration:

```sh
kubectl apply -f k8s/Service.yaml
```

4. Get the service URL:

```sh
minikube service list
```

The service will be available at the outputted `URL`

## Docker Hub
I have also uploaded the images to my [Docker Hub Profile](https://hub.docker.com/u/zvdy) Both for 32 and non 32 core CPU's.

## 📬 Making Requests

Here are some example `curl` commands to interact with the API:

- Create a coupon:

```sh
curl -X POST http://localhost:8080/api/create -d '{"discount": 10, "code": "Superdiscount", "minBasketValue": 50}' -H "Content-Type: application/json"
```

- Apply a coupon:

```sh
curl -X POST http://localhost:8080/api/apply -d '{"basket": {"value": 100}, "code": "Superdiscount"}' -H "Content-Type: application/json"
```

- Retrieve coupons: 

```sh
curl -X GET http://localhost:8080/api/coupons -d '{"codes": ["Superdiscount"]}' -H "Content-Type: application/json"
```

- Retrieve many coupons: 

```sh
curl -X GET http://localhost:8080/api/coupons -d '{"codes": ["Superdiscount1", "Superdiscount2", "Superdiscount3"]}' -H "Content-Type: application/json"
```

- Check HTTP response count and metrics. Uses _`prometheus`_ 

```sh
curl -X GET http://localhost:8080/metrics
```

## 🏋️‍♀️ Load Testing

To perform load testing on the service, we have a Go script named `loadtest.go`. This script sends multiple concurrent requests to the service and measures the response times.

Before running the load test, make sure the service is running either locally or in a Docker container or Kubernetes Cluster.

---

To run the load test, follow these steps:

1. Navigate to the directory containing `loadtest.go`:

```sh
cd path/to/directory
```

2. Run the `loadtest.go` script:

```sh
go run loadtest.go
```

This will start the load testing process. The script will print the response times for each request and a summary at the end.

Please note that you might need to adjust the parameters in the `loadtest.go` script to match your testing requirements, such as the number of concurrent requests or the total number of requests to send.

---

> You can use _[jq](https://jqlang.github.io/jq/_)_ in order to get formatted/prettier outputs just execute your curl command as usual, then add:  | jq and it will be formated 