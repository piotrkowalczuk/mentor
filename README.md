# Zordon
Zordon is a tool for defining and running multi-services Go applications.
With Zordon, you use a Alphasfile to configure your application’s services.
Then, using a single command, you create and start all the services from your configuration.
Any application that is go gettable can be used as one.

## Installation
To build Zordon from the source code yourself you need to have a working Go environment.
You can directly use the go tool to download and install the Zordon.

```bash
$ go get github.com/piotrkowalczuk/zordon
```


## Commands

* **morphintime** - starts all services in specified order. It performs restart on exit code 1.
* **powerup** - is trying to update each service. If any change is found by git, service will be skipped
* **recruit** - install all services

## Alphasfile
Zordon can run services automatically by reading the instructions from a Alphasfile.
A Alphasfile is a [hlc](https://github.com/hashicorp/hcl) document that contains all the definitions.

### Example

```hlc
service "gnatsd" {
  import = "github.com/nats-io/gnatsd"

  arguments {
    D = false
    V = false
    T = false
    p = 9010
    m = 9011
  }
}

service "prometheus" {
  import = "github.com/prometheus/prometheus/cmd/prometheus"
  log = "json"

  arguments {
    log.format = "logger:stdout?json=true"
    web.listen-address = ":9020"
  }
}

```

## License

[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl.txt), see [LICENSE](LICENSE).