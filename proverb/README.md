# Proverb

Proverbs provides a single JSON endpoint `/poverbs/random` which returns one of
the nineteen [Go Proverbs](https://go-proverbs.github.io), a collection of
design principles distilled into pity statements. The application is written in
Go.

## Building the Docker Image

To build your own copy of the Docker image, use:

```
docker build --force-rm -t gcr.io/[GCP_PROJECT_ID_HERE]/proverb:latest .
```
