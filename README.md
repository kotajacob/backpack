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
There are four different operations: `buy`, `add`, `remove`, and `set` which
take a string indicating an item with an optional count and price. If the count
is given it comes first and if the price is given it comes last.

# buy
If no count is given it will be 1. Buy does not accept a price option in the
request. The `owner` option must always be used with `buy`. The owner is the
buyer and the seller is the room from which the command is called. If the seller
has enough in stock and the buyer has enough coins both are removed and the item
is given to the buyer.
```
/inv owner[#finn] buy[10 apples]
/inv owner[#gordon] buy[10 regular arrows]
/inv owner[#aurora] buy[mighty sword]
```

## add
If no count is given it will be 1. If no price is given the price will simply
not be changed. The default price is "not for sale".
```
/inv add[bow]
/inv add[2 apple]
/inv add[4 apple 10]
```

## remove
If no count is given ALL of that item will be removed. The price works just like
in add.
```
/inv remove[bow]
/inv remove[2 apple]
/inv remove[4 apple 10]
```

# set
Set is the simplest of the operations. It will just change the given inventory
to match the given number and value of items.
```
/inv set[25 regular arrows 2]
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
