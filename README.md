# FileCrypt

This command line utility will:
- encrypt a file into a gzip file.
- decrypt the encrypted gzip file back into the original file


Here is the help message:
```
$ FileCrypt -help
Usage:
  -d    decrypt flag
  -e    encrypt flag
  -help
        show usage message
  -i string
        input file to be encrypted; required
  -key string
        encryption key; required
  -o string
        output file to be decrypted; required
This code will encrypt and write to a gzip file; or decrypt from same
Usage: FileCrypt -key keystring -e -i filename_to_encrypt -o encypted_filename
      or
Usage: FileCrypt -key keystring -d -i filename_to_decrypt -o decripted_filename
where -e means encrypt and -d means decrypt
```

Here is a sample use:
```
$ FileCrypt -e -key BingoWasHisNameOHisNameO -i xx.go -o xx.gz
$ FileCrypt -d -key BingoWasHisNameOHisNameO -i xx.gz -o yy.go
$ cksum xx.* yy.go
371583179 3029 xx.go
3795034373 959 xx.gz
371583179 3029 yy.go
$ 
```

Note: the key is adjusted to be 32 bytes which invokes the AES-256 algorithm.
If the key provided is longer than 32 bytes, then only the first 32 are used.
