# kubernetes-schemas-store

A simple application to sync Kubernetes CRDs, and objects json schema to a remote storage with rclone

## Motivation

- Save time writing CRDs manifests
- No more wasting time to lookup for online json schema
- Can be self-hosted

## Usage

### Setup storage and kubernetes connection

1. Choose a compatible storage from https://rclone.org/overview/

2. Create rclone config https://rclone.org/docs/

> rclone config file is usually found at `~/.config/rclone/rclone.conf`

Assume we have

```conf
[my-s3]
type = S3
provider = Minio
access_key_id = my-access-key
secret_access_key = my-secret-key
endpoint = https://s3.example.com
no_check_bucket = true
```

3. (Optional) Create the bucket

This step can be skipped by removing `no_check_bucket`. Make sure your key has permission for creating bucket

4. (Optional) Create kubeconfig

> kubeconfig files are usually found at `~/.kube`

5. (Optional) Make sure your bucket does not cache the manifests, so they can be updated in the future

### Upload

#### In-cluster usage

- Usually run as a container with service account permissions

```bash
kss --auth-method in-cluster --destination my-s3:/my-bucket
```

#### Outside of cluster

```bash
kss --auth-method kubeconfig --kubeconfig ~/.kube/config --destination my-s3:/my-bucket
```
