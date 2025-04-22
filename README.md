# 🌈 go-playground

Welcome to **Go Playground**! This repository serves as a collection of Go-related examples, demonstrations, and experiments. It's designed to showcase various features of the Go programming language and provide practical code snippets for learning and reference.

---

## 🌟 **About**
This repository contains:
- Hands-on examples of Go's core features.
- Code demonstrations for common Go patterns.
- Experimentation with Go libraries and frameworks.
- Best practices and idiomatic Go code.

Whether you're learning Go, refreshing your skills, or exploring advanced concepts, this repository is a great resource to get started.


## 논점 1
빈 슬라이스 생성할 때 nil initialization을 해야 하는가?

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go

var s []string


```
</td><td>

```go

s := []string{}
```
</td></tr>
</tbody></table>
