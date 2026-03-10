# DoorPiX

_DoorPiX_ is an application to run on a small computer like a raspberry pi and act as your doorbell. It allows the communication via SIP/VOIP protocol to integrate it with existing systems like a Fritz!Box like any other SIP/VOIP or video doorbell. The project is inspired by the [DoorPi](https://github.com/motom001/DoorPi), which unfortunatly is no longer under active development.

## Configuration

As the original project, DoorPiX is event-action based. There are some components which fire events and some which react on these. What to do for a specifc event is configured in the configuration file called `doorpix.yaml`. Each event is described via a path, wich contains some basic information about the event and some additional properties. You can filter for multiple events by using wildcards in the path. All events follow  a prefix like `<scope>/<node>/<type>/<event specfic path>`. Where the scope is `system` most of the time, there is also the `internal` scope, which is used for internal events, nothing will prevent you from using these, but typically the `system` should be enough. The `node` part is always `doorpix` for this application and can be used in conjuction with other systems. As an example other systems can fire events on the doorpi which than could start with the prefix `external/sensor/...` wich allows the integration of 3rd party systems. The `type` specifies some information about the  source of the event e.g. a `service` event or a `lifecycle` event. A path from a GPIO event would look like `system/doorpix/gpio/<pin number>/<edge>`. To listen for a rising edge on any pin you can use the wildcard path `system/doorpix/gpio/*/rising`. To listen on every edge of either pin 14 or 15 you could use `system/doorpix/gpio/[14,15]/*`. The wildcard syntax follows the [Match](https://pkg.go.dev/path#Match) syntax. 

For each event, there can be a number of actions e.g. call someone or fire a webhook. There is also an action to execute any arbitrary shell command, which increases the attac surface, but allows the implementation of any kind action, which can be done with shell commands. This is likely to be disabled per default in the future when the most of typical usecases are implemented. The configuration might look like:

```yaml
events:
  # call a number and execute a shell command if pin 14 is pressed
  - event: system/doorpix/gpio/14/rising
    steps:
      - type: invite
        with:
          uri: sip:0123456789@sip.example.org
      - type: shell
        with:
          cmd: echo "Hello World"
```

## Develop

### Dev-Shell

```sh
nix develop -c $SHELL
go run ./cmd/doorpix/main.go 
```

### Test

```sh
nix flake check
```
