## KV

Another kv thing,

The only reaon this exists is because I need a simple DB for alfred/keyboard maestro tasks.

## Install

```
  go install github.com/fmeyer/kv@latest
```

## Usage

### In a bash script

```
V=$(kv g -k "$KEY" 2>/dev/null)
ES=$?

if [ $ES -eq 0 ]; then
    # do whatver you need with $VALUE

else
    # exec task store result in $V

    kv s -k "$KEY" -v $V
fi

```

### HELP

```
Usage:
  kv [command]

Available Commands:
  g           Get a value
  l           List all keys
  s           Set a value

Flags:
  -h, --help   help for kv

Use "kv [command] --help" for more information about a command.
```


### GET
```
Usage:
  kv g [flags]

Flags:
  -k, --key string   Key
```

### SET
```
Usage:
  kv s [flags]

Flags:
  -k, --key string     Key
  -v, --value string   Value
```

### LIST

```
Usage:
  kv l [flags]

```
