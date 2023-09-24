Start server:

```bash
go run ./lesson2
```

GET operation:

```bash
curl localhost:13377/get
```

PUT operation:

```bash
curl -d deadbeef localhost:13377/replace
```

