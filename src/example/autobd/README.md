# autobd - Automagic backup daemon
autobd is an automated backup daemon that uses inotify watches to act on events
when they happen. So basically it will watch a directory for events, then back
up all your stuff to other servers when you change, add or delete something.

Think of it as portal between multiple computers where you can put something in
and it appears in every other computer.

## Why
To learn some more Go, because I can, and because I have a use for this (mostly
because I don't have the patience for rsync/scp or any alternatives)

## Features
- Recursively watching directories: done
- Configuration file: done
- Setting custom event handles: in progress

## Future Features
- Server/client Code
    - server/client authentication using SSL certificates
- Hashing files on change
- Robust event logging
- Talk to influxdb with JSON for analytics
- ???

## Possible Features
(Not outside of the realm of possibility but less likely to become reality)
- Delta encoding for file transfers
- Daily snapshots (tarballs most likey, by host/date)


## Helping
As you might have noticed I'm new to Go but I'm enjoying it so far. If you see
something wrong or want to help out in general, drop me a line. I would really
appreciate any help you can give. You can reach me at <tyrell.wkeene@gmail> or
on twitter <@tywkeene>
