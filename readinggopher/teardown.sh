#!/usr/bin/env bash

set -xe

gcloud container images delete gcr.io/enocom-dev/readinggopher:95b649d9f5c51b9be4da457906d6515b5c5999f2

gcloud compute instances delete readinggopher-vm

gcloud compute firewall-rules delete allow-http
