# cSense

![cSense logo](https://raw.githubusercontent.com/danielkrainas/csense/master/docs/logo/csense-logo.png)

cSense (Container Sense) allows you to subscribe to container events with web hooks. Hooks are registered with cSense via an HTTP API and may contain selectors, like image tag or container name, to limit the containers or events that the hook should be notified about.

For more information, see the [project site.](https://github.com/danielkrainas/csense)

## Usage

```
$ docker run dakr/csense agent
```
