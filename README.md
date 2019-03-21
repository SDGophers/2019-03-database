# Database

Today's challange is a simlple disk backed data base which should be almost
[ACID](https://en.wikipedia.org/wiki/ACID_\(computer_science\)). The database
will contain data about the world's countries.

The end goal is to read off data from several json files to fill the database
and persist the data in a binary file. Code should be written to also query
this database quickly; in O(log n) via presorting and binary search or O(1) via
hashing (note: this is the hard route but a good exercise).

## Serialization

Serialization is the process of turning data structures into a (machine or
computer) readable, persistant format. This is a lot of jargon, but some
examples include JSON, YAML, XML, and Protocol Buffers. These file formats
are used to share data in a manner independent of the programming language used
to generate them. Although, these formats can sometimes be considered their own
programming languages.

The first thing we have to do is find a way to serialize a single country.
Here's what we currently have:
```go
// Country is a single database row
type Country struct {
	Name string // change me
	Population uint64
	North, South, East, West float64
}
```

There is one problem with the current structure, the name field is variable
length. One way to get around this is to use a null terminated string as in
the C/C++. This means, that the string ends once the first byte of value `0`
is reached.

Run the following command to get an example of this:
```text
go run scratch/serialize.go | hexdump -C
variable sized
00000000  6d 65 20 74 6f 77 6e 00  01 00 00 00 00 00 00 00  |me town.........|
00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00000030
```

Looking at the code for this program you'll notice an if statement. This checks
if the `Country` structure if it's statically sized (it is currently not).

This there a few problems with this serialization method but most important is
speed. If we know the size of our row struture, we can quickly iterate and
index our table by using multiples of the row size.

To achive this, we'll change the Name field to a static array. This is not
conventional `go`, but will simplify the code later. Here's what it should look
like.

```go
// Country is a single database row.
type Country struct {
	//Name string // change me
	Name [KeySize]byte
	Population uint64
	North, South, East, West float64
}
```

Due to go's strict typing, the `serialize.go` program will not work but the
the necessary changes are commented in the file.

```text
go run serialize.go | hexdump -C
static size:  96
00000000  6d 65 20 74 6f 77 6e 00  00 00 00 00 00 00 00 00  |me town.........|
00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00000030  00 00 00 00 00 00 00 00  01 00 00 00 00 00 00 00  |................|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00000060
```

## Database interface

Let's look at the country database interface. This is used by some prewritten
programs to test, benchmark, and populate the database.

```go
type CoutntryDB interface {
	Get(name string) (Country, error)
	Set(name string, data Country) error
	Del(name string) (Country, error)
}
```

This defines an arbitrary key-value store with `string` keys and `Country`
values. Note that name doesn't have to be a separate argument, because the
`Country` struct has a `Name` field. I chose to leave this because it quickly
communicates what we are indexing by.

Another possible function could be
```go
ForEach(cb func(string, Country) error) error
```
would call cb `for each` key-value pair in the database. One difficulty with
this is that it can be complicated to implement with "atomicity".

## Database implementation

So, the task to create an implemntation of CountryDB.

Note that while the country name is redundant as a key and field in the
interface, this redundancy is not necessary in the binary file.

### List Implementation

Okay, now to the fun stuff. Our database will be essentially a contiguous array
living on disk. As previously stated, the elements are statically sized and can
be quickly indexed by their position in the array. We'll use this to our
advantage. Here's a

* **get** - binary search through the array
* **set** - search for the element, if it exists overwrite, else shift the back
	elements and write the new one
* **del** - search for th element, if exists shift the back elements towards
	front & truncate file, else do nothing.

### Hash Map implementation

A hash map implementation similarly takes advantage of the array indexing.
First, the key needs a hash function, this a function that takes in the key
and returns an integer. For strings, one could sum every byte of the string.
This hash value is then modded (%) with the table size to get a valid index
within the hash table. This index is a starting point to look for an element.
From this index one could search linearly until a match is found or a
tombstone/null element is reached.
[wikipedia has a detailed explanation](https://en.wikipedia.org/wiki/Hash_table)

