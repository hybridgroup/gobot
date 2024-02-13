# I2C

This document describes some basics for developers.

## Byte order

All common libraries (smbus, digispark, firmata, i2cget) read and write I2C data in the order LSByte, MSByte.  
Often the devices store its bytes in the reverse order and therefor needs to be swapped after reading.

## Linux syscall implementation

In general there are different ioctl features for I2C

* IOCTL I2C_RDWR, needs "I2C_FUNC_I2C"
* IOCTL SMBUS, needs "I2C_FUNC_SMBUS.."
* SYSFS I/O
* call of "i2c_smbus_* methods"

>The possible functions should be checked before by "I2C_FUNCS".

## SMBus ioctl workflow and functions

> Some calls are branched by kernels [i2c-dev.c:i2cdev_ioctl()](https://elixir.bootlin.com/linux/latest/source/drivers/i2c/i2c-dev.c#L392)
> to the next listed calls.

Set the device address `ioctl(file, I2C_TARGET, long addr)`. The call set the address directly to the character device.

Query the supported functions `ioctl(file, I2C_FUNCS, unsigned long *funcs)`. The call is converted to in-kernel function
[i2c.h:i2c_get_functionality()](https://elixir.bootlin.com/linux/latest/source/include/linux/i2c.h#L902)

Execute a function, if supported `ioctl(file, I2C_SMBUS, struct i2c_smbus_ioctl_data *args)`. The call is converted by
kernels [i2c-dev.c:i2cdev_ioctl_smbus](https://elixir.bootlin.com/linux/latest/source/drivers/i2c/i2c-dev.c#L311) to
[i2c.h:i2c_smbus_xfer()](https://elixir.bootlin.com/linux/latest/source/include/linux/i2c.h#L140). This leads to call of
[i2c-core-smbus.c:i2c_smbus_xfer_emulated()](https://elixir.bootlin.com/linux/latest/source/drivers/i2c/i2c-core-smbus.c#L607)
if adapter.algo->smbus_xfer() is not implemented (that means there is no implementation for the current platform/adapter).

```C
/* This is the structure as used in the I2C_SMBUS ioctl call */
struct i2c_smbus_ioctl_data {
  __u8 read_write;
  __u8 command;
  __u32 size;
  union i2c_smbus_data __user *data;
};

/*
 * Data for SMBus Messages
 */
#define I2C_SMBUS_BLOCK_MAX 32  /* As specified in SMBus standard */
union i2c_smbus_data {
  __u8 byte;
  __u16 word;
  __u8 block[I2C_SMBUS_BLOCK_MAX + 2]; /* block[0] is used for length */
             /* and one more for user-space compatibility */
};

```

```C
// default preparation for call of "status = __i2c_transfer(adapter, msg, nmsgs)"

msgbuf0[I2C_SMBUS_BLOCK_MAX+3];
msgbuf1[I2C_SMBUS_BLOCK_MAX+2];

msgbuf0[0] = command;

struct i2c_msg msg[2] = {
  {
    .addr = addr,
    .flags = flags,
    .len = 1,
    .buf = msgbuf0,
  }, {
    .addr = addr,
    .flags = flags | I2C_M_RD,
    .len = 0,
    .buf = msgbuf1,
  },
};
/* note: msg[0].buf == msgbuf0, means
         msg[0].buf[0] == msgbuf0[0]
         msg[0].buf[1] == msgbuf0[1]
         msg[0].buf[2] == msgbuf0[2]
*/
```

### Data flow for I2C_FUNC_SMBUS_WRITE_BYTE

```C
i2c_smbus_write_byte(const struct i2c_client *client, u8 value);
\\would call:
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_WRITE, value, I2C_SMBUS_BYTE, NULL);
```

gobot: d.smbusAccess(I2C_SMBUS_WRITE, val, I2C_SMBUS_BYTE, nil), calls without a platform driver  

```C
i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write=I2C_SMBUS_WRITE, command=val, size=I2C_SMBUS_BYTE, *data=NULL);
```

```C
// by default preparation (see above)
msg[0].len = 1;
msg[0].buf[0] = command; // corresponds to "val"
nmsg = 1;
```

### Data flow for I2C_FUNC_SMBUS_WRITE_BYTE_DATA

```C
i2c_smbus_write_byte_data(const struct i2c_client *client, u8 command, u8 value);
// would call:
data.byte = value;
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_WRITE, command, I2C_SMBUS_BYTE_DATA, &data);
```

gobot: d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&data)), calls without a platform driver

```C
datasize = sizeof(data->byte);
copy_from_user(&temp, data, datasize);

i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write=I2C_SMBUS_WRITE, command=reg, size=I2C_SMBUS_BYTE_DATA, *data=&temp);

// leads to:
msg[0].len = 2;
msg[0].buf[0] = command; // corresponds to "reg"
msg[0].buf[1] = data->byte;
nmsg = 1;
```

### Data flow for I2C_FUNC_SMBUS_WRITE_WORD_DATA

```C
i2c_smbus_write_word_data(const struct i2c_client *client, u8 command, u16 value);
// would call:
data.word = value;
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_WRITE, command, I2C_SMBUS_WORD_DATA, &data);
```

gobot: d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_WORD_DATA, unsafe.Pointer(&data)), calls without a platform driver

```C
datasize = sizeof(data->word);
copy_from_user(&temp, data, datasize);

i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write=I2C_SMBUS_WRITE, command=reg, size=I2C_SMBUS_WORD_DATA, *data=&temp);

// leads to:
msg[0].len = 3;
msg[0].buf[0] = command;
msg[0].buf[1] = data->word & 0xff;
msg[0].buf[2] = data->word >> 8;
nmsg = 1;
```

### Data flow for I2C_FUNC_SMBUS_WRITE_BLOCK_DATA

```C
i2c_smbus_write_block_data(const struct i2c_client *client, u8 command, u8 length, const u8 *values);
// would call:
data.block[0] = length; // set first data element to length
memcpy(&data.block[1], values, length); // ...followed by the real data values
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_WRITE, command, I2C_SMBUS_BLOCK_DATA, &data);
```

gobot: do not use this feature, but this would call without a platform driver

```C
datasize = sizeof(data->block); // it seems this just blocks 32+2
copy_from_user(&temp, data, datasize); // but possibly this only copy what is there

i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_WRITE, command=reg, size==I2C_SMBUS_BLOCK_DATA, *data==&temp);

// leads to:
// this reads the length from given user data, first element, add 2 for command and one more for user space compatibility:
msg[0].len = data->block[0] + 2;
msg[0].buf[0] = command; // by i2c_smbus_try_get_dmabuf(&msg[0], command)
msg[0].buf + 1 = data->block; // copy with size!, by memcpy(msg[0].buf + 1, data->block, msg[0].len - 1);
nmsg = 1;
```

> "i2c_smbus_write_block_data" adds the real data size as first value in data. Writing this to the device is intended
> for SMBus devices, see "section 6.5.7 Block Write/Read" of [smbus 3.0 specification](http://smbus.org/specs/SMBus_3_0_20141220.pdf).
> Implementing this behavior in gobot leads to writing the size as first data element, but this is normally not useful
> for i2c devices. When using ioctl calls there is no way to drop the size element by using this call.

### Data flow for I2C_FUNC_SMBUS_WRITE_I2C_BLOCK

```C
i2c_smbus_write_i2c_block_data(const struct i2c_client *client, u8 command, u8 length, const u8 *values);
// would call:
data.block[0] = length; // set first data element to length
memcpy(data.block + 1, values, length); // ...followed by the real data values
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_WRITE, command, I2C_SMBUS_I2C_BLOCK_DATA, &data);
```

gobot:

```go
// set the first element with the data size
dataLen := len(data)
buf := make([]byte, dataLen+1)
buf[0] = byte(dataLen)
copy(buf[1:], data)
d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_I2C_BLOCK_DATA, unsafe.Pointer(&buf[0]))
```

...calls without a platform driver

```C
datasize = sizeof(data->block); // it seems this just blocks 32+2
copy_from_user(&temp, data, datasize); // but possibly this only copy what is there

i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_WRITE, command=reg, size==I2C_SMBUS_I2C_BLOCK_DATA,
  *data==&temp);

// this reads the length from given user data, first element, add 1 for command:
msg[0].len = data->block[0] + 1;
msg[0].buf[0] = command; // by i2c_smbus_try_get_dmabuf(&msg[0], command)
msg[0].buf + 1 = data->block + 1; // copy real data without size, by memcpy(msg[0].buf + 1, data->block + 1, data->block[0]);
nmsg=1
```

### Data flow for I2C_FUNC_SMBUS_READ_BYTE

```C
i2c_smbus_read_byte(const struct i2c_client *client);
// would call:
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, &data);
return data.byte;
```

gobot: d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, unsafe.Pointer(&data)), calls without a platform driver

```C
datasize = sizeof(data->byte); // &temp keeps empty
i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_READ, command=0, size==I2C_SMBUS_BYTE, *data==&temp);

msg[0].len = 1; // per default
msg[0].flags = I2C_M_RD | flags;
msg[0].buf[0] = 0; // from command
nmsgs = 1;
```

afterwards copy content to target structure

```C
data->byte = msg[0].buf[0]; // copy read value back to given data pointer and one level above
copy_to_user(data, &temp, datasize); // copy one byte back to given data pointer
```

### Data flow for I2C_FUNC_SMBUS_READ_BYTE_DATA

```C
i2c_smbus_read_byte_data(const struct i2c_client *client, u8 command);
// would call:
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_READ, command, I2C_SMBUS_BYTE_DATA, &data);
return data.byte;
```

gobot: d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&data)), calls without a platform driver

```C
datasize = sizeof(data->byte); // &temp keeps empty
i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_READ, command=reg, size==I2C_SMBUS_BYTE_DATA, *data==&temp);

msg[0].len = 1; // per default
msg[0].buf[0] = command;
msg[1].len = 1;
nmsgs = 2;
```

afterwards copy content to target structure

```C
data->byte = msg[1].buf[0]; // copy read value back to given data pointer and one level above
copy_to_user(data, &temp, datasize); // copy one byte back to given data pointer
```

### Data flow for I2C_FUNC_SMBUS_READ_WORD_DATA

```C
i2c_smbus_read_word_data(const struct i2c_client *client, u8 command);
// would call:
status = i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_READ, command, I2C_SMBUS_WORD_DATA, &data);
return data.word;
```

gobot: d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_WORD_DATA, unsafe.Pointer(&data)), calls without a platform driver

```C
datasize = sizeof(data->word); // &temp keeps empty
i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_READ, command=reg, size==I2C_SMBUS_BYTE_DATA, *data==&temp);

msg[0].len = 1; // per default
msg[0].buf[0] = command;
msg[1].len = 2;
nmsgs = 2;
```

afterwards copy content to target structure

```C
data->word = msgbuf1[0] | (msgbuf1[1] << 8); // copy read value back to given data pointer and one level above
copy_to_user(data, &temp, datasize); // copy 2 bytes back to given data pointer
```

### Data flow for I2C_FUNC_SMBUS_READ_BLOCK_DATA

```C
i2c_smbus_read_block_data(const struct i2c_client *client, u8 command, u8 *values);
// would call:
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_READ, command, I2C_SMBUS_BLOCK_DATA, &data);
// data.block[0] is the read length (N)
memcpy(values, &data.block[1], data.block[0]); // copy starting from data.block[1] to "values", N elements
return data.block[0]; // number of read bytes (N)
```

gobot: do not use this feature, but this would call without a platform driver

```C
datasize = sizeof(data->block); // &temp keeps empty
i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_READ, command=reg, size==I2C_SMBUS_BYTE_DATA, *data==&temp);

msg[0].len = 1; // per default
msg[0].buf[0] = command;
msg[1].flags |= I2C_M_RECV_LEN;
msg[1].len = 1; // block length will be added by the underlying bus driver
msg[1].buf[0] = 0; // by i2c_smbus_try_get_dmabuf(&msg[1], 0);
nmsgs = 2;
```

afterwards copy content to target structure

```C
memcpy(data->block, msg[1].buf, msg[1].buf[0] + 1); // copy read value back to given data pointer and one level above
copy_to_user(data, &temp, datasize); // copy all bytes (according to given size) back to given data pointer
```

### Data flow for I2C_FUNC_SMBUS_READ_I2C_BLOCK

```C
i2c_smbus_read_i2c_block_data(const struct i2c_client *client, u8 command, u8 length, u8 *values);
// would call:
data.block[0] = length; // set first data element to length
i2c_smbus_xfer(client->adapter, client->addr, client->flags, I2C_SMBUS_READ, command, I2C_SMBUS_I2C_BLOCK_DATA, &data);
// data.block[0] is the read length (N)
memcpy(values, &data.block[1], data.block[0]); // copy starting from data.block[1] to "values", N elements
return data.block[0]; // number of read bytes (N)
```

gobot:

```go
dataLen := len(data)
// set the first element with the data size
buf := make([]byte, dataLen+1)
buf[0] = byte(dataLen)
copy(buf[1:], data)
log.Printf("buffer: %v", buf)
d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_I2C_BLOCK_DATA, unsafe.Pointer(&buf[0])); err != nil {
// get data from buffer without first size element
copy(data, buf[1:])
```

...calls without a platform driver

```C
datasize = sizeof(data->block); // &temp keeps empty, datasize includes the "size" byte (means real data + 1)
copy_from_user(&temp, data, datasize); // but possibly this only copy what is there

i2c_smbus_xfer_emulated(*adapter, addr, flags, read_write==I2C_SMBUS_READ, command=reg, size==I2C_SMBUS_I2C_BLOCK_DATA, *data==&temp);

msg[0].len = 1; // per default
msg[0].buf[0] = command;
msg[1].len = data->block[0]; // this reads the length from given user data, first element:
msg[1].buf[0] = 0; // by i2c_smbus_try_get_dmabuf(&msg[1], 0);
nmsgs = 2;
```

afterwards copy content to target structure

```C
// copy read values back, starting from second value (contains length)
memcpy(data->block + 1, msg[1].buf, data->block[0]);
// and one level above:
copy_to_user(data, &temp, datasize); // copy all bytes (according to given size) back to given data pointer
// this means, the caller needs to strip the real data starting from second byte (like "i2c_smbus_read_i2c_block_data")
```

## Links

* <https://www.kernel.org/doc/Documentation/i2c/dev-interface>
* <https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/i2c-dev.h#L42>
* <https://stackoverflow.com/questions/9974592/i2c-slave-ioctl-purpose> (good for understanding, but there are some
  small errors in the provided example)
* <https://stackoverflow.com/questions/55976683/read-a-block-of-data-from-a-specific-registerfifo-using-c-c-and-i2c-in-raspb>

> Qotation from kernel.org: "If possible, use the provided i2c_smbus_* methods described below instead of issuing direct
> ioctls."  
> We do not do this at the moment, instead we using the "IOCTL SMBUS".

Because the syscall needs uintptr in Go, there are some known pitfalls with that. Following documents could be helpful:

* <https://go101.org/article/unsafe.html>
* <https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang>
* <https://stackoverflow.com/questions/59042646/whats-the-difference-between-uint-and-uintptr-in-golang>
* <https://go.dev/play/p/Wd7hWn9Zsu>
* for go vet false positives, see: <https://github.com/golang/go/issues/41205>

Basically by convert to an uintptr, which is than just a number to an object existing at the moment of creation without
any other reference, the garbage collector will possible destroy the original object. Therefor uintptr should be avoided
as long as possible.
