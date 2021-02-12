# s3sync
aws s3 sync CLI tool.

## Whai is this?
This is a CLI tool that synchronizes the state of your local directory with AWS S3.

## Usege
* Upload and sync local files in bulk.
```
% s3sync upload [-f]
```

* Download all at once from s3 and sync.
```
% s3sync download [-f]
```

* Upload and sync local a file.
```
% s3sync upload -m file
file list> aaa.txt
bbb.txt
ccc.txt
...
```

* Download and sync s3 a file.
```
% s3sync download -m file
file list> AAA.txt
BBB.txt
CCC.txt
...
```
