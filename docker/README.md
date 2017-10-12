# Containerized builds

In order to improve debugging cycle time, there are two containers:

- `wptd-base` - includes all necessary files
- `wptd-testrun` - the final container for running the tests

First build the base container.

```sh
cd docker/wptd-base
docker build -t wptd-base .
```

Then you can build the test run container.

```sh
docker build -t wptd-testrun .

SAUCE_USER=rutabaga
SAUCE_KEY=rutabaga
docker run \
    -e "PLATFORM_ID=edge-15-windows-10-sauce" \
    -e "RUN_PATH=gamepad" \
    -e "SAUCE_USER=$SAUCE_USER" \
    -e "SAUCE_KEY=$SAUCE_KEY" \
    -p 4445:4445 \
    wptd-testrun
```

Remove the `RUN_PATH` variable if you wish to run the whole suite.

## Push a new version to the registry

Be advised this is dangerous as all builds use this container.

```sh
IMAGE_NAME=gcr.io/wptdashboard/wptd-testrun
docker tag wptd-testrun $IMAGE_NAME
gcloud docker -- push $IMAGE_NAME
```

## Start a test run VM

This starts a VM that runs a containerized test run and uploads the results.

```sh
VM_NAME=test-vm-docker-run
gcloud compute instances create $VM_NAME \
    --metadata-from-file startup-script=vm-startup.sh \
    --zone us-central1-c \
    --image-project cos-cloud \
    --image cos-stable-55-8872-76-0
```

To view the logs:

```sh
gcloud compute ssh $VM_NAME

# In the VM
tail -f /var/log/syslog
```