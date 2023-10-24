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

```yaml
ble ble
```

## Config file location

TBD

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