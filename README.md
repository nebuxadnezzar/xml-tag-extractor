# xml-tag-extractor
Extracts a XML document between a set of tags and transforms it into single line document

### Command line
```text

Usage: xte [OPTIONS] XML-file-path
OPTIONS:
  -boost
    	boost processing, helps big file processing
  -ca
  -convert.attributes
    	convert attributes to elements
  -h
  -help
    	show help
  -ol
  -one.liner
    	transform XML document into one-liner

  -rt string
  -root.tag string
    	add provided root tags to each document to make correct XML document

  -xp string
  -xml.paths string
    	CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story

```
### Examples
given this file
```xml
<root>
    <greetings>hello</greetings>
    <greetings>good bye <times>3</times>
    </greetings>
    <greetings id="123"/>
    <smiles>wide</smiles>
</root>
```
... to extract ```greetings``` tags together with its content use command
```bash
xte -ol -xp root:greetings my.xml
```
the output will be one-document-per-line records
```xml
<greetings>hello</greetings>
<greetings>good bye <times>3</times></greetings>
<greetings id="123"/>
```
to extract ```times``` tags together with its content use command
```bash
xte xte -ol -xp root:greetings:times my.xml
```

if not __tag-path__ argument provided __xte__ will print all tag paths and their count found in an XML file which can be useful, if you have huge file and don't know XML structure of the file.
```xml
xte my.xml
```
```
root		1
root:greetings		3
root:greetings:times		1
root:smiles		1

```

### Notes
>__xte__ doesn't parse out and convert attributes to elements. this will be next feature