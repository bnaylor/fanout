# fanout example

Decided to try my hand at a kinda cool fanout model that I read about
[here](https://go.dev/blog/pipelines).

Basically it sets up a couple of "controllers" that each run a 3-stage
pipeline with a configurable number of goroutines, and uses channels to 
control the flow and thread termination/etc.

I went for a slightly more real-world example with the intent of actually 
adapting it for use @work..  There is a need for .. uh.. "bounded app discovery" :)

Copious console spam exists to track what's going on and to debug it.

- [another relevant link](https://dave.cheney.net/2013/04/30/curious-channels)
- [one more](https://dave.cheney.net/2014/03/19/channel-axioms)
- [original reason i was looking at it](https://www.reddit.com/r/golang/comments/48mnrp/go_channels_are_bad_and_you_should_feel_bad/)


My observation from getting it all to work is that it is really, really
easy to deadlock yourself with this stuff.  It is cool once it's *absolutely perfect* 
but until then it's a bear.
