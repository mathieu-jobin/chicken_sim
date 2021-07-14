# TBC moonkin dps sim

The sim is still in development and has bugs.
You can use it at https://thatcodingguy.github.io/chicken_sim/

# Development

To run, clone/download the repository locally.
Install golang https://golang.org/

Edit main.go with your character stats as you see on the character preview.
Enable any trinkets you have equipped, and change the talents based on your build.

Run with 
```
go run .
```

# Known issues.
- Trinkets don't currently share cooldowns between them.
- Partial resists aren't modeled yet
- You can equip more than 2 trinkets.
- No haste support.

If you run into an issue or wish to contribute, please open an issue or a pull request!
