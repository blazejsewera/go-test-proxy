# Go Test Proxy

A Go proxy to plug between the frontend
and backend services for last-resort testing.

## Why?

Here's the story:

You are dropped inside a project.
However, it doesn't have any tests.
Your onboarding process is nowhere to be seen,
but somehow you manage to get the application running —
both the backend and frontend.

You've got a task:
add a feature to the frontend application,
but there's a catch:
the backend endpoint is not yet implemented.
Your manager told you that
"we want only the frontend part for now,"
and you have to find a way to move forward somehow.

Alright, you've been in such situations before.
You tell your fellow backend expert
that you want to negotiate how this endpoint should look like,
but all you hear is that they are not yet assigned to the backend part,
and they have more important things to do for now.

So you think you will simply write this endpoint yourself,
but as you wander through the backend code,
you are astonished how untested the code can be.
You finally drop this idea,
instead opting for a simple prototype of this endpoint
in some kind of Express.js, FastAPI, or some other quick-to-write framework.

Your frontend teammate sees you do this
and wants to get it running on their machine,
so you spend some totally avoidable amount of time
helping them set up the right version of Node.js, Yarn, or Python.

After all this, it turns out that the application blows up
without other backend endpoints,
so you have to forward _all of them_,
except the one that you want to test.
You want to just give one binary to your fellow colleagues,
so they get unblocked.
You want to be able to just cross-compile it.
Just do it and after you're done,
just have another tool in your software engineer's toolbox.

This is it. This is a tool.
