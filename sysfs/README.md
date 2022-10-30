# sysfs

## I2C

### Byte order

All common libraries (smbus, digispark, firmata, i2cget) read and write I2C data in the order LSByte, MSByte.  
Often the devices store its bytes in the reverse order and therefor needs to be swapped after reading.

### Linux syscall implementation

In general there are different ioctl features for I2C

* IOCTL I2C_RDWR, needs "I2C_FUNC_I2C"
* IOCTL SMBUS, needs "I2C_FUNC_SMBUS.."
* SYSFS I/O
* call of "i2c_smbus_* methods"

>The possible functions should be checked before by "I2C_FUNCS".

for further reading see:

* https://www.kernel.org/doc/Documentation/i2c/dev-interface
* https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/i2c-dev.h#L42
* https://stackoverflow.com/questions/9974592/i2c-slave-ioctl-purpose (good for understanding, but there are some small errors in the provided example)

>Qotation from kernel.org: "If possible, use the provided i2c_smbus_* methods described below instead of issuing direct ioctls." We do not do this at the moment, instead we using the "IOCTL SMBUS".

Because the syscall needs uintptr in Go, there are some known pitfalls with that. Following documents could be helpful:

* https://go101.org/article/unsafe.html
* https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
* https://stackoverflow.com/questions/59042646/whats-the-difference-between-uint-and-uintptr-in-golang
* https://go.dev/play/p/Wd7hWn9Zsu
* for go vet false positives, see: https://github.com/golang/go/issues/41205

Basically by convert to an uintptr, which is than just a number to an object existing at the moment of creation without
any other reference, the garbage collector will possible destroy the original object. Therefor uintptr should be avoided
as long as possible.
