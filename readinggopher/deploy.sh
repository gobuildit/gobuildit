#!/usr/bin/env bash

set -xe

docker push gcr.io/enocom-dev/readinggopher:0e77b51dba2d0c29dd4374e6937b62b2a01a6916

# Configure a firewall to allow ingress HTTP traffic
gcloud compute firewall-rules \
    create allow-http \
    --target-tags http-server \
    --allow tcp:80

# Create and deploy the container
gcloud beta compute instances \
    create-with-container readinggopher-vm \
    --tags http-server \
    --container-env CONSUMER_KEY="$CONSUMER_KEY",CONSUMER_SECRET="$CONSUMER_SECRET",TOKEN="$TOKEN",TOKEN_SECRET="$TOKEN_SECRET",FREQUENCY="$FREQUENCY" \
    --container-image \
    gcr.io/enocom-dev/readinggopher:0e77b51dba2d0c29dd4374e6937b62b2a01a6916
