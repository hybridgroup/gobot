# Gobot CLI

Gobot has its own CLI to generate new platforms, adaptors, and drivers.

## Building the CLI

```
go build -o /path/to/dest/gobot .
```

## Running the CLI

```
/path/to/dest/gobot help
```

Should display help for the Gobot CLI:

```
CLI tool for generating new Gobot projects.

	NAME:
		 gobot - Command Line Utility for generating new Gobot adaptors, drivers, and platforms

	USAGE:
		 gobot [global options] command [command options] [arguments...]

...
```

## Installing from the snap

Gobot is also published in the [snap store](https://snapcraft.io/). It is not yet stable, so you can help testing it in any of the [supported Linux distributions](https://snapcraft.io/docs/core/install) with:

```
sudo snap install gobot --edge
```

## License
Copyright (c) 2013-2017 The Hybrid Group. Licensed under the Apache 2.0 license.
