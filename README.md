# backpack

An inventory management discord bot.

# build and run
```
go build .
export BACKPACK_TOKEN=GET_ONE_FROM_DISCORD
export BACKPACK_DATA=/home/backpack/data
./backpack
```

# usage
There are two different operations: `add` and `remove` which take a string
indicating an item to add with an optional count and price.

## add
The count is specified first. If it is missing the count will be 1. The price is
specified last and if it is missing the price will simply not be changed. The
default price is simply "not for sale".
```
/inv add[bow]
/inv add[2 apple]
/inv add[4 apple 10]
```

## remove
The count is specified first. If it is missing ALL of that item will be removed.
The price works just like in add.
```
/inv remove[bow]
/inv remove[2 apple]
/inv remove[4 apple 10]
```

## owner
An owner may be specified and will be used instead of the current channel's
name. For example, if a channel named `#finn` exists this will give 1 apple to
that inventory instead of the current channel:
```
/inv owner[#finn] add[1 apple]
```

# author
Written and maintained by Dakota Walsh.
Up-to-date sources can be found at https://git.sr.ht/~kota/backpack/

# license
GNU AGPL version 3 or later, see LICENSE.
Copyright 2022 Dakota Walsh
