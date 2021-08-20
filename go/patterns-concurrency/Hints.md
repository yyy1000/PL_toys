Convert data state into code state when it makes programs clearer.</br>
Convert mutexes into goroutines when it makes programs clearer.</br>
Use additional goroutines to hold additional code state.</br>
Use goroutines to let independent concerns run independently.</br>
Consider the effect of slow goroutines.</br>
Know why and when each communication will proceed.</br>
Know why and when each goroutine will exit.</br>
Type Ctrl-\ to kill a program and dump all its goroutine stacks.</br>
Use the HTTP server’s /debug/pprof/goroutine to inspect live goroutine stacks.</br>
Use a buffered channel as a concurrent blocking queue.</br>
Think carefully before introducing unbounded queuing.</br>
Close a channel to signal that no more values will be sent.</br>
Stop timers you don’t need.</br>
Prefer defer for unlocking mutexes.</br>
Use a mutex if that is the clearest way to write the code.</br>
Use a goto if that is the clearest way to write the code.</br>
Use goroutines, channels, and mutexes together
if that is the clearest way to write the code.</br>