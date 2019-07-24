# go-asciitable
Draw ASCII tables in console. Bloody alpha at the edge of being unusable, your PRs are welcome.

## Usage

As simple as that (as of now):

```go
table := asciitable.NewSimpleTable()
table.Header("Name", "Description")
table.AddRow("Strawberry", "Fruit in yours granma's garten")
table.AddRow("Raspberry", "Useless hardware you anyway will buy")
table.AddRow("Perl", "Also language")
fmt.Println(table.Render())
```

Result:

```
+----------+------------------------------------+
|Name      |Description                         |
+----------+------------------------------------+
|Strawberry|Fruit in yours granma's garten      |
+----------+------------------------------------+
|Raspberry |Useless hardware you anyway will buy|
+----------+------------------------------------+
|Perl      |Also language                       |
+----------+------------------------------------+
```
