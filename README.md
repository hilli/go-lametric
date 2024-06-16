# go-lametric

Go lib for controlling the LaMetric devices (SKY and TIME (2022 or later). 

Also: `lametric-homekit-hub` - Bridge LaMetric devices with Apple Home.

So far only display on/off and brightness controls, as that is what make the most sense to support.

## Using with Apple Home

Check the [releases](https://github.com/hilli/go-lametric/releases) for something suitable. Setup the environment to point to the LaMetric device and set the key

- `LAMETRIC_HOSTNAME` to the IP or a hostname on your network pointing to your device.
- `LAMETRIC_API_KEY` to an API key. Get it from the settings on the device in the LaMetric app.
- `LAMETRIC_DIY_PUSH_URL` (optional) to the DIY app's push URL.

## Developing

This repo uses [Taskfile](https://taskfile.dev/). If you haven't already go read Taskfile.dev's [installation instructions](https://taskfile.dev/installation/).

### Run tests

#### Unit tests

```shell
task test
```

#### Integration tests

For code examples and running against a device use

```shell
task test-integration
```

## Notes

LaMetric Developer Documentation:
https://lametric-documentation.readthedocs.io/en/latest/

