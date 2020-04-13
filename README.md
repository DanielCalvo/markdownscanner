### Markdownscanner (TM)
Are your links broken? Let's find out: https://mdscanner.dcalvo.dev/  
Please note that this project is under development and I would not consider it ready for general use.  

### Okay but now for real
While signing up to contribute to k8s, I found a broken link on the sign up process. [This was then my my first contribution.](https://github.com/kubernetes/community/pull/4304)  
I then wondered: How many other markdown links are broken in open source projects? As it turns out, a lot of them.  
This tool will hopefully help me find and fix these links.

### Next actionable items
- Order output on the report, 404s first!
- Create usage instructions
- Dockerfile is missing!
- Set a timeout / handle timeouts appropriately (the velodrome link is an example)
- Don't check github issues!
- Also ignore checks to github users
- See how the 429 re-check logic would fit in together with the rest of the stuff
- Create a flag for the port number and set it as 8080 by default
- How is the allProjects flag going to work moving forward? `//allProjects := flag.Bool("allprojects", false, "Do you want to scan everything?")`
- Find a way to get github groups by API so that you automatically get all the github projects!
- It looks like your "IGNORED" logic is quite expensive. Any way to make it faster? Can you benchmark it?
- Also don't check release notes. Check some example on Kops
- How do you handle an invalid repository or a typo on repositories.yaml?
- How do you handle a timeout on a repository clone?
- Don't forget to document which links get ignored
- Don't forget to somehow implement header checks for markdown files? Those might be tough...
- To ignore: If file is named changelog
- Recheck go.mod and general install (create instructions)
- Instructions on how to use: https://mdscanner.dcalvo.dev/
- To ignore: If the destination is a pull request on github
- Review error handling (handle errors better than just printing stuff (maybe see where it may be appropriate to panic, and handle timeoouts too!))
- Prioritize repositories that have not been scanned yet. Then scan from oldest to newest!
- Unit tests are missing
- Delete repo once scan is completed. Make this an option later
- Add the Scan results for PROJECT part on the project page.
- Use some kind of Logger class/module (depending on language) that allows you to log at different levels (DEBUG, INFO, etc.)
- Don't forge to handle timeouts appropriately (both on repositories and on http links)
- It would be nice to have the metadata on the scan page as well (Put an indication of which repo the results page is for on the results page. You might need to send more data over)
- Repo clone sometimes times out and acts weird
- Organize this in packages. That big 400 line single file isn't nice to work with.
- Go over all the comments in the code and take action on them. And clean the code up, too.
- Maybe check later: https://www.datadoghq.com/blog/go-logging/
- Remove the excessive newlines on the html if you can, more info here: https://github.com/golang/go/issues/9969
- Reorganize code a bit (You need to be a bit more clear with this one)

Error during git clone when running only for the first time:
```text
2020/04/11 01:06:14 Beginning main loop: Iterating over Git Repositories from repositories.yaml
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x93b78b]

goroutine 1 [running]:
main.GetMarkdownFiles.func1(0x0, 0x0, 0x0, 0x0, 0xb54e00, 0xc0008bcb10, 0x0, 0x0)
        /home/daniel/tmp/markdownscanner/main.go:99 +0x5b
path/filepath.Walk(0x0, 0x0, 0xc0000fbb70, 0x0, 0x0)
        /usr/lib/go-1.13/src/path/filepath/path.go:402 +0x6a
main.GetMarkdownFiles(0x0, 0x0, 0xc0000287e0, 0x25, 0x0, 0x0, 0x0)
        /home/daniel/tmp/markdownscanner/main.go:98 +0xd4
main.main()
        /home/daniel/tmp/markdownscanner/main.go:341 +0x122d

Process finished with exit code 2
```

Revise the explanatory text, see where it can fit:
```text
In more detail
- The program finds all the Markdown files in a given project. A markdown file is any file that ends in ".md"
- The program then ops all of those markdown files and find all of the links inside them. Links in markdown look like
- The program will then check the destination of all these links and generate a report. Links can be working (200), pointing to resources that cannot be found (404) or in a variety of other states.
- Click on HTML report below to see that status of the links in markdown files for a given project.
```
