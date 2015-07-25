

## Docker Builds

There is a script to build all projects that produce a binary executable and have a Dockerfile (the db is excluded).  Remove the build dir first to ensure a clean build:

```
rm -rf docker-build-tmp
./docker.sh
```
