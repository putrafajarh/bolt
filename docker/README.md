# Docker README

This folder contains Dockerfiles for building images based on `alpine`, `busybox`, and `debian`.

## Security
All container images in this project are configured to run as a non-root user. This ensures enhanced security by minimizing the privileges of the running processes within the containers. Running as a non-root user helps mitigate potential risks and limits the impact of vulnerabilities or malicious code execution.

## How to Build Specific Dockerfile

To build a specific Dockerfile, use the following commands:

### Build Alpine Image
```bash
docker build -f Dockerfile.alpine -t myimage:alpine ../
```

### Build BusyBox Image
```bash
docker build -f Dockerfile.busybox -t myimage:busybox ../
```

### Build Debian Image
```bash
docker build -f Dockerfile.debian -t myimage:debian ../
```

### Build Distroless Image
```bash
docker build -f Dockerfile.debian -t myimage:debian ../
```

## Notes
- Replace `myimage` with your desired image name.
- Ensure you are in the same directory as the Dockerfiles when running the commands. 
- For `distroless` image it lacks a libc (or other system libraries), your Go binary must be statically linked (`CGO_ENABLED=0`), Otherwise, it will fail with errors. If your package uses CGO, find alternative packages that are pure Go or use different base image.
