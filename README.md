# TBC moonkin dps sim

To run, clone/download the repository locally.
Install golang https://golang.org/

Edit main.go with your character stats as you see on the character preview.
Enable any trinkets you have equipped, and change the talents based on your build.

Run with 
```
go run .
```

To run the webapp locally, do :
```
./web.sh && go run devserver/main.go 
```

# Known issues.
- Trinkets don't currently share cooldowns between them.
- Partial resists aren't modeled yet
- You can equip more than 2 trinkets.
- No haste support.

If you run into an issue or wish to contribute, please open an issue or a pull request!
