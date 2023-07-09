# go-tb-dedup
Go-based Thunderbird Mailbox Deduper

# NOTES

- Some code is purposefully not covered; and while the Go coverage tooling has chosen
  to not support any tagging to mark it ignored (see https://github.com/golang/go/issues/53271)
  this code does mark it by placing `//go: cover ignore` before the block to be ignored.

- This code demonstrates one of the worse performance cases of Go.
