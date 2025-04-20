# Docker README

This folder contains Dockerfiles for building images based on `alpine`, `busybox`, and `debian`.

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

## Notes
- Replace `myimage` with your desired image name.
- Ensure you are in the same directory as the Dockerfiles when running the commands.
