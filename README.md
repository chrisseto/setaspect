# Set Aspect

`setaspect` is a simple CLI to adjust the aspect ratio of PNG images by adding transparent padding.

Really, I just want to post screenshots on twitter and have them formattted appropriately.

For that reason, the default aspect ratio is 16:9.

## Installation

```
go install github.com/chrisseto/setaspect
```

## Usage

Write to an output file:
```
setaspect ./my-screenshot.png -o ./what-i-will-post-on-twitter.png
```

Or redirect to an output file
```
setaspect ./my-screenshot.png > ./what-i-will-post-on-twitter.png
```
