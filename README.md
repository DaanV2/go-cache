# Cache

[![Pipeline](https://github.com/DaanV2/go-cache/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/DaanV2/go-cache/actions/workflows/pipeline.yaml)

[WIP] some sets / maps usable for large amount of storage of items concurrently.

```go
col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
require.NoError(t, err)

items := test_util.Generate(int(size))
test_util.Shuffle(items)

for _, item := range items {
	v, ok := col.GetOrAdd(item)
	require.True(t, ok)
	require.Equal(t, v, item)
}
```
