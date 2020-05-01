## i3-pomodoro: a pomodoro timer for i3blocks

This simple timer allows you to track Pomodoros with the greatest of ease!

### To use:

```shell
$ go get github.com/erikh/i3-pomodoro
```

Then put this in your i3blocks configuration

```ini
[pomodoro]
label=🕒
command=~/bin/i3-pomodoro # dir here should be $GOBIN or $GOPATH/bin
interval=1
```

And restart i3blocks or i3xrocks. Woo!

### Mouse Controls

- Mouse 1: Start/Pause the timer
- Mouse 2: Reset the timer
- Mouse 3: Cycle to the next timer (resetting the timer)

### Notes

Creates a file named `/tmp/pomodoro-state`. If anything bugs out, try removing
it!

Runs `i3-nagbar` to notify the user. You can change this to use `notify-send`
pretty easily near the top of the source.

Same with the timers and clock designations, they are also near the top.

### Author

Erik Hollensbe <erik+git@hollensbe.org>
