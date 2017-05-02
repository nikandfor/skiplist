#!/bin/sh

if [ "$#" -lt 1 ] ; then
	echo usage: $0 raw_type
	exit 1
fi

tp=$1

dir=cg

mkdir -p $dir
rm -f $dir/*_test.go

sed "s/interface{} \/\* val \*\//$tp/g" skiplist.go > cg/skiplist.go

if [ "$tp" = "int" ] ; then
	sed "s/, ok := \(.*\).Value().(int); !ok ||/ := \1.Value();/g" skiplist_test.go | sed "s/ \w*.Value() == nil ||//" | sed "s/.Value().(int)/.Value()/" > cg/skiplist_test.go
fi

cat >cg/less.go <<EOF
package skiplist

var IntAsc LessFunc = func(a, b $tp) bool { return a < b }
EOF
