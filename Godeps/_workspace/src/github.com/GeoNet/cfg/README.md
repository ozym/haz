# cfg

Golang library for application configuration.

* Config is defined using JSON and can optionally be overridden using environment variables.
* The configer tool can be used to automatically generate a file `docker-run.sh` which documents which configuration properties can be overridden from environment variable.

Refer to the documentation on `cfg.Load()` for further details.

configer is suitable for use with `go generate` e.g., include it as a comment in your code:

```
...
//go:generate configer fits.json
var (
	config = cfg.Load()
)
...
```

Then run `go generate` to create the `docker-run.sh` file from your JSON config file.

The log prefix is set to the executable name.  The an additional string can optionally be appended to the prefix by setting cfg.Build
at compile time e.g.,

```
 ... -ldflags "-X github.com/GeoNet/cfg.Build `git rev-parse --short HEAD`" ...
 ```
