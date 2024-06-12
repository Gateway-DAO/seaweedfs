.PHONY: image kompose

image:
	docker build -t gateway/fs:

kompose:
	rm ../k8/*
	kompose -f docker-compose.local.yml convert -o ../k8
