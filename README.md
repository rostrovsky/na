# sodium

<p align="center">
    <img src="./na.png" alt="Sodium">
</p>

**N**ested **a**liases for your shell.

## TODO

* [ ] support for minimal form config
* [x] test on powershell
* [ ] readme
* [x] version
* [x] debug logs

## Example

1. Set aliases in config file.
2. Generate autocompletions for your shell:

    * bash

        ```shell
        source <(na completion bash)
        ```

    * zsh

        ```shell
        source <(na completion zsh)
        ```

    * powershell

        ```shell
        Register-ArgumentCompleter -CommandName na -ScriptBlock $__naCompleterBlock
        na completion powershell | Out-String | Invoke-Expression
        ```

    * fish

        ```shell
        na completion fish | source
        ```

## Config file schema

### Minimal form

```yaml
# minimal form config example
aliases:
  ssh:
    dev:
      host-x: ssh user@host-x
      host-y: ssh user@host-y
    prod:
      host-a: ssh user@host-x
      host-b: ssh user@host-y
  grep:
    heron: grep -hEron --with-filename --color=always
```

### Full form

In comparison to minimal form, in full format:

* you must add mandatory `_cmd` key which contains aliased command
* you can add optional `_info` key that enriches autocompletion hints with description.

```yaml
# full form config example
aliases:
  ssh:
    _info: aliases for SSH connections
    dev:
      _info: DEV environments aliases
      host-x:
        _info: makes ssh connection to host X on DEV env
        _cmd: ssh user@host-x
      host-y:
        _info: makes ssh connection to host Y on DEV env
        _cmd: ssh user@host-y
    prod:
      _info: PROD environments aliases
      host-a:
        _info: makes ssh connection to host A on PROD env
        _cmd: ssh user@host-x
      host-b:
        _info: makes ssh connection to host B on PROD env
        _cmd: ssh user@host-y
  grep: # _info key is completely optional though recommended
    heron:
      _cmd: grep -hEron --with-filename --color=always
```

## Config file location

By default, `na` expects configuration file placed in `~/.config/sodium/.narc.yaml`.

You can override default config file location by setting `SODIUM_CONFIG` environment variable.

## Supported shells

`na` supports all [cobra](https://github.com/spf13/cobra)-generated autocompletions:

* bash
* fish
* powershell
* zsh

## Reserved aliases

Below aliases cannot be used due to being [cobra](https://github.com/spf13/cobra) bultins:

* `completion`
* `help`
* `--help`
* `-h`

## Debug logs

If need to to see debug logs, set `SODIUM_LOG_LEVEL` env variable to `debug` (case insensitive).

## License

MIT
