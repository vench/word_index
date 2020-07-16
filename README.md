# Word index

Quick search for keywords in documents.


Example find key word in document.

```
documents = []string{
		`Sorry. We’re having trouble getting your pages back.`,
		`Still not able to restore your session?`,
		`Our Dockerfile will have two section.`,
		`Quick search for keywords in documents.`,
}

index := NewIndexBin()
index.Add(documents...)
i := index.Find(`key*`)
if i != -1 {
    println(documents[i])
}
```

### TODO

[ ] bin operations