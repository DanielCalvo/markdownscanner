### Markdownscanner (TM)
Are your links broken? Let's find out: https://mdscanner.dcalvo.dev/ 
Please note that this project is under development and I would not consider it ready for general use.  

### Okay but now for real
While signing up to contribute to k8s, I found a broken link on the sign up process. [This was then my my first contribution.](https://github.com/kubernetes/community/pull/4304)  
I then wondered: How many other markdown links are broken in open source projects? As it turns out, a lot of them.  
This tool will hopefully help me find and fix these links.

### Next actionable items
- Scan from oldest to newest!
- Create a flag for the port number and set it as 8080 by default
- Make sure repository folder is clean at program start up
- Review error handling (handle errors better than just printing stuff (maybe see where it may be appropriate to panic, and handle timeoouts too!))
- How does kubernetes does logging? Or kubelet maybe? Or some other thing? Is there logging in packages?
- Set a timeout / handle timeouts appropriately (the velodrome link is an example). Further info here: http://networkbit.ch/golang-http-client/#minimal
- Create usage instructions (with docker!)
- Also ignore checks to github users (but how?)
- Don't forget to document which links get ignored
- Don't forget to somehow implement header checks for markdown files? Those might be tough...
- Recheck go.mod and general install (create instructions)
- Unit tests are missing (How do you handle an invalid repository or a typo on repositories.yaml?)
- Add the <h3>Scan results for PROJECT</h3> part on the project page.
- Use some kind of Logger class/module (depending on language) that allows you to log at different levels (DEBUG, INFO, etc.)
- It would be nice to have the metadata on the scan page as well (Put an indication of which repo the results page is for on the results page. You might need to send more data over)
- Organize this in packages. That big 400 line single file isn't nice to work with.
- Go over all the comments in the code and take action on them. And clean the code up, too.
- Maybe check later: https://www.datadoghq.com/blog/go-logging/
- Remove the excessive newlines on the html if you can, more info here: https://github.com/golang/go/issues/9969
- Reorganize code a bit (You need to be a bit more clear with this one)

#### Docker shenanigans
- All commands were ran from the project root:
```
docker build . -t danitest
docker run danitest -it /bin/bash
docker run -d -v $(pwd)/html_static:/app/server/html_static -p 8080:8080 danitest
```