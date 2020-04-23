### Markdownscanner (TM)
Are your links broken? Let's find out: https://mdscanner.dcalvo.dev/ 
Please note that this project is under (sporadic) development and it's not finished. I just needed an excuse to code something...  

### Okay but now for real
While signing up to contribute to k8s, I found a broken link on the sign up process. [This was then my my first contribution.](https://github.com/kubernetes/community/pull/4304)  
I then wondered: How many other markdown links are broken in open source projects? As it turns out, a lot of them.  
This tool will hopefully help me find and fix these links.

### Things to do
- Scan projects that have not been scanned the longest first
- Create a flag for the port number and set it as 8080 by default
- Make sure repository folder is clean at program start up
- Review error handling (handle errors better than just printing stuff (maybe see where it may be appropriate to panic, and handle timeoouts too!))
- Check how having using a more sophisticated logging method could improve the program
- Set a timeout / handle timeouts on HTTP checks properly. Further info here: http://networkbit.ch/golang-http-client/#minimal
- Create usage instructions
- Is there a way to ignore checks on github users?
- Don't forget to document which links get ignored somewhere
- Find a way to implement header checks for markdown files. That might be tough...
- Recheck go.mod and general install (create instructions)
- There are no unit tests. You should create some!
- Add the `Scan results for PROJECT` part on the project page.
- Do some research on "how to organize your go project" and apply it here.
- Remove the excessive newlines on the html if you can, more info here: https://github.com/golang/go/issues/9969

#### Docker shenanigans
- All commands were ran from the project root:
```
docker build . -t danitest
docker run danitest -it /bin/bash
docker run -d -v $(pwd)/html_static:/app/server/html_static -p 8080:8080 danitest
```