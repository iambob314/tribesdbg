# tribesdbg

Runs Tribes.exe normally and waits for it to exit. Then, it highlights any error logs that occurred while running.
Only works if console.log is active (`$Console::logMode > 0`).

Log patterns highlighted:
```
something.cs Line: 1234 - Syntax error.
someFunc: Unknown command.
```

Usage:
```
tribesdbg.exe Tribes.exe [args to pass to Tribes.exe here]
```
