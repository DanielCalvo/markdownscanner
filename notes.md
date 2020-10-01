### Actionable items
- Replace all the fmt messages with log messages
- Add more log messages!
    - Add shorter log messages for the checks too!
- Missing production setup on orange pi
    - Set up AWS user with proper permissions (terraform this part too)
    - Cron set up
- Set up a working
- Is it possible to write a Prometheus Exporter?
- Stick to a coding standard pls (is it c or config? r or repo?)
- Order the functions in the packages (put close together what happens close together)
- Refine HTTP ERR (check the error and take it from there, it's more valuable to know that something is a timeout rather than a 404)
- Can you have a status page? You can have an independent routine running that updates a status page
- You are still missing tests!
```
Checked: {CHANGELOG-3.0.md /tmp/mdscanner/etcd-io/etcd/CHANGELOG-3.0.md https://github.com/etcd-io/etcd/tree/master/CHANGELOG-3.0.md [code changes] https://github.com/etcd-io/etcd/compare/v3.0.7...v3.0.8 HTTP 200}
```

### Code TODOs
- Redo the Dockerfile for the executable thing
- Careful with what you export (https://blog.golang.org/organizing-go-code)
- https://peter.bourgon.org/go-best-practices-2016/#repository-structure
- https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

- About Cobra: https://github.com/spf13/cobra
    - I think you don't need cobra for now as you'll only run with one parameter, but maybe investigate more functionality later
    - Like running the program with a parameter (--project) or to save results locally and don't upload them to S3 (--dont-upload)
    - Viper remains to be explored later as well: https://github.com/spf13/viper

### Things to do
- Scan projects that have not been scanned the longest first
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


#### Terraform
- Create the following:
    - S3 bucket (not public)
    - User to access the bucket
    - A policy for the user granting access only to that bucket

#### Golang
- Create a program that will upload files with S3 using user credentials


### md scanner ideas
- Create a cloudwatch event that launches an instance, does something and then exits?

- How do I stop and start Amazon EC2 instances at regular intervals using Lambda?
    - https://aws.amazon.com/premiumsupport/knowledge-center/start-stop-lambda-cloudwatch/

- How to connect to S3?
    - https://docs.aws.amazon.com/AmazonS3/latest/dev/security-best-practices.html

- Can you use Lambda pre signed URLs to upload files securely?
    - https://medium.com/@lakshmanLD/upload-file-to-s3-using-lambda-the-pre-signed-url-way-158f074cda6c

- How is Lambda local development done again? Can you launch a function that uploads a mock file to s3?

#### Maybe
- Just run things on an ec2 instance for now, worry about migrating to t2 micro later

#### Possible feature ideas
- Can you take the feature toggle approach as defined in the cloud ops book?
    - Application could save something either to S3 or local filesystem
    - "How to feature toggle golang"
    - https://featureflags.io/go-feature-flags/