## wsoffcli

Collect data from https://ws-tcg.com/

### Synopsis

Collect data from https://ws-tcg.com/.

Create a json file for each card with most information.

Example:
'wsoffcli fetch -n IMC' will fetch all cards with a code starting with 'IMC'

If you want more than one use '##' as seperator like 'wsoffcli fetch -n BD##IM'

'--serie' use a hidden number in the official site, this number is increment for each new set (e.g Kadokawa is number 259, Goblin 260 ...).

To use environ variable, use the prefix 'WSOFF'.
	 

### Options

```
  -h, --help           help for wsoffcli
  -n, --neo string     Neo standar by set
  -s, --serie string   serie number
```

### SEE ALSO

* [wsoffcli completion](doc/wsoffcli_completion.md)	 - Generate the autocompletion script for the specified shell
* [wsoffcli fetch](doc/wsoffcli_fetch.md)	 - Fetch cards
* [wsoffcli gendoc](doc/wsoffcli_gendoc.md)	 - Generate doc with Cobra
* [wsoffcli products](doc/wsoffcli_products.md)	 - Get products information

###### Auto generated by spf13/cobra on 16-Nov-2023
