# sflicense

"Cracking Cisco’s Sourcefire license system" writeup code

## Project structure

```
.
├── LICENSE
├── README.md
├── checklic
│   └── checklic.go
├── crypto
│   └── rsaGen.go
├── licgen
│   └── sflicgen.go
└── tools
    └── interceptor.sh
```

## Cracking guide

**Do this at your own risk!**

**Cracking software is illegal and unmoral, please use this only for testing and educational
purposes.**

### RSA key generation

Execute:

```bash
go run rsaGen.go
```

This will generate a public and a private RSA key of 4096 bits.

Format the public key for the patchline with:

```bash
hexdump -ve '1/1 "_x%.2x"' public.der |sed 's/_/\\/g'
```

### Calculate the SHA1 of the public key

On macos:

```bash
shasum -a1 crypto/public.der |cut -d" " -f1 |sed -E 's/(.{2})/\1\\x/g' |rev |cut -d"\\" -f2- |rev
```

### Prepare the patchlines (x86-64 version) 


```bash
printf '${your_formated_public_key}' | dd seek=$((0x10b48)) conv=notrunc bs=1 of=${target}

printf '${your_formated_sha1}' | dd seek=$((0x10f48)) conv=notrunc bs=1 of=${target}
```

### Patch!

SSH into the FSM and get sudo, execute the two patchlines changing the `${target}` for
the _checklic_ binary.

Repeat the same process on all your Sensors.

### Generate your license

For a FSM license run:

```bash
go run sflicgen.go -l 66:00:11:22:33:44:55 -k ../crypto/private.pem -fsm
```

and for a sensor license run:

```bash
go run sflicgen.go -l 66:00:11:22:33:44:55 -k ../crypto/private.pem -n 6 -mid 63E -mod 3D8120
```

Upload it to your FSM and enjoy


## Known models:

| Model ID | Model Name |
|----------|------------|
| 63E      | 3D8120     |
| 63F      | 3D7120     |
| 63G      | 3d7110     |
| 63H      | 3D7030     |
| 63J      | 3D7010     |
| 63L      | 3D7125     |
| 63P      | 3D7150     |


