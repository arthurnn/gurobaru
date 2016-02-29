# gurobaru

Is a unique identifier server based on mysql sequence generator.

### How to use it.

Build:
```
go build
```

Run the server:
```
./gurobaru
```

Test it out
```
echo "FETCH ID" | nc localhost 8000
```

Or go to: http://localhost:8000
