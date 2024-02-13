package spherocommon

// LocatorConfig provides configuration for the Location api.
// For more information refer to the api specification of "Orbotix Communication API"
// see also: http://wiki.mark-toma.com/view/Sphero_API_Tutorial
// The current (X,Y) coordinates of Sphero on the ground plane in centimeters.
type LocatorConfig struct {
	// Determines whether calibrate commands automatically correct the yaw tare value
	Flags uint8
	// Controls how the X-plane is aligned with Sphero’s heading coordinate system.
	X int16
	// Controls how the Y-plane is aligned with Sphero’s heading coordinate system.
	Y int16
	// Controls how the X,Y-plane is aligned with Sphero’s heading coordinate system.
	YawTare int16
}

// CollisionConfig provides configuration for the collision detection alogorithm.
// For more information refer to the api specification of "Orbotix Communication API"
// see also: http://wiki.mark-toma.com/view/Sphero_API_Tutorial
type CollisionConfig struct {
	// Detection method type to use. Methods 01h and 02h are supported as
	// of FW ver 1.42. Use 00h to completely disable this service.
	Method uint8
	// An 8-bit settable threshold for the X (left/right) axes of Sphero.
	// A value of 00h disables the contribution of that axis.
	Xt uint8
	// An 8-bit settable threshold for the Y (front/back) axes of Sphero.
	// A value of 00h disables the contribution of that axis.
	Yt uint8
	// An 8-bit settable speed value for the X axes. This setting is ranged
	// by the speed, then added to Xt to generate the final threshold value.
	Xs uint8
	// An 8-bit settable speed value for the Y axes. This setting is ranged
	// by the speed, then added to Yt to generate the final threshold value.
	Ys uint8
	// An 8-bit post-collision dead time to prevent retriggering; specified
	// in 10ms increments.
	Dead uint8
}

// DataStreamingConfig provides configuration for Sensor Data Streaming.
// For more information refer to the api specification of "Orbotix Communication API"
// see also: http://wiki.mark-toma.com/view/Sphero_API_Tutorial
type DataStreamingConfig struct {
	// Divisor of the maximum sensor sampling rate
	N uint16
	// Number of sample frames emitted per packet
	M uint16
	// Bitwise selector of data sources to stream
	Mask uint32
	// Packet count 1-255 (or 0 for unlimited streaming)
	Pcnt uint8
	// Bitwise selector of more data sources to stream (optional)
	Mask2 uint32
}

// PowerStatePacket contains all data relevant to the power state of the sphero
type PowerStatePacket struct {
	// record Version Code
	RecVer uint8
	// High-Level State of the Battery; 1=charging, 2=battery ok, 3=battery low, 4=battery critical
	PowerState uint8
	// Battery Voltage, scaled in 100th of a Volt, 0x02EF would be 7.51 volts
	BattVoltage uint16
	// Number of charges in the total lifetime of the sphero
	NumCharges uint16
	// Seconds awake since last charge
	TimeSinceChg uint16
}

// CollisionPacket represents the response from a Collision event
type CollisionPacket struct {
	// Normalized impact components (direction of the collision event):
	X, Y, Z int16
	// Thresholds exceeded by X (1h) and/or Y (2h) axis (bitmask):
	Axis byte
	// Power that cross threshold Xt + Xs:
	XMagnitude, YMagnitude int16
	// Sphero's speed when impact detected:
	Speed uint8
	// Millisecond timer
	Timestamp uint32
}

// DataStreamingPacket represents the response from a Data Streaming event
type DataStreamingPacket struct {
	// 8000 0000h	accelerometer axis X, raw	-2048 to 2047	4mG
	RawAccX int16
	// 4000 0000h	accelerometer axis Y, raw	-2048 to 2047	4mG
	RawAccY int16
	// 2000 0000h	accelerometer axis Z, raw	-2048 to 2047	4mG
	RawAccZ int16
	// 1000 0000h	gyro axis X, raw	-32768 to 32767	0.068 degrees
	RawGyroX int16
	// 0800 0000h	gyro axis Y, raw	-32768 to 32767	0.068 degrees
	RawGyroY int16
	// 0400 0000h	gyro axis Z, raw	-32768 to 32767	0.068 degrees
	RawGyroZ int16
	// 0200 0000h	Reserved
	Rsrv1 int16
	// 0100 0000h	Reserved
	Rsrv2 int16
	// 0080 0000h	Reserved
	Rsrv3 int16
	// 0040 0000h	right motor back EMF, raw	-32768 to 32767	22.5 cm
	RawRMotorBack int16
	// 0020 0000h	left motor back EMF, raw	-32768 to 32767	22.5 cm
	RawLMotorBack int16
	// 0010 0000h	left motor, PWM, raw	-2048 to 2047	duty cycle
	RawLMotor int16
	// 0008 0000h	right motor, PWM raw	-2048 to 2047	duty cycle
	RawRMotor int16
	// 0004 0000h	IMU pitch angle, filtered	-179 to 180	degrees
	FiltPitch int16
	// 0002 0000h	IMU roll angle, filtered	-179 to 180	degrees
	FiltRoll int16
	// 0001 0000h	IMU yaw angle, filtered	-179 to 180	degrees
	FiltYaw int16
	// 0000 8000h	accelerometer axis X, filtered	-32768 to 32767	1/4096 G
	FiltAccX int16
	// 0000 4000h	accelerometer axis Y, filtered	-32768 to 32767	1/4096 G
	FiltAccY int16
	// 0000 2000h	accelerometer axis Z, filtered	-32768 to 32767	1/4096 G
	FiltAccZ int16
	// 0000 1000h	gyro axis X, filtered	-20000 to 20000	0.1 dps
	FiltGyroX int16
	// 0000 0800h	gyro axis Y, filtered	-20000 to 20000	0.1 dps
	FiltGyroY int16
	// 0000 0400h	gyro axis Z, filtered	-20000 to 20000	0.1 dps
	FiltGyroZ int16
	// 0000 0200h	Reserved
	Rsrv4 int16
	// 0000 0100h	Reserved
	Rsrv5 int16
	// 0000 0080h	Reserved
	Rsrv6 int16
	// 0000 0040h	right motor back EMF, filtered	-32768 to 32767	22.5 cm
	FiltRMotorBack int16
	// 0000 0020h	left motor back EMF, filtered	-32768 to 32767	22.5 cm
	FiltLMotorBack int16
	// 0000 0010h	Reserved 1
	Rsrv7 int16
	// 0000 0008h	Reserved 2
	Rsrv8 int16
	// 0000 0004h	Reserved 3
	Rsrv9 int16
	// 0000 0002h	Reserved 4
	Rsrv10 int16
	// 0000 0001h	Reserved 5
	Rsrv11 int16
	// 8000 0000h	Quaternion Q0	-10000 to 10000	1/10000 Q
	Quat0 int16
	// 4000 0000h	Quaternion Q1	-10000 to 10000	1/10000 Q
	Quat1 int16
	// 2000 0000h	Quaternion Q2	-10000 to 10000	1/10000 Q
	Quat2 int16
	// 1000 0000h	Quaternion Q3	-10000 to 10000	1/10000 Q
	Quat3 int16
	// 0800 0000h	Odometer X	-32768 to 32767	cm
	OdomX int16
	// 0400 0000h	Odometer Y	-32768 to 32767	cm
	OdomY int16
	// 0200 0000h	AccelOne	0 to 8000	1 mG
	AccelOne int16
	// 0100 0000h	Velocity X	-32768 to 32767	mm/s
	VeloX int16
	// 0080 0000h	Velocity Y	-32768 to 32767	mm/s
	VeloY int16
}
