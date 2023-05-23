package mavlink

//
// MAVLink comm protocol generated from common.xml
// http://qgroundcontrol.org/mavlink/
//
import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var messages = map[uint8]MAVLinkMessage{
	0:   &Heartbeat{},
	1:   &SysStatus{},
	2:   &SystemTime{},
	4:   &Ping{},
	5:   &ChangeOperatorControl{},
	6:   &ChangeOperatorControlAck{},
	7:   &AuthKey{},
	11:  &SetMode{},
	20:  &ParamRequestRead{},
	21:  &ParamRequestList{},
	22:  &ParamValue{},
	23:  &ParamSet{},
	24:  &GpsRawInt{},
	25:  &GpsStatus{},
	26:  &ScaledImu{},
	27:  &RawImu{},
	28:  &RawPressure{},
	29:  &ScaledPressure{},
	30:  &Attitude{},
	31:  &AttitudeQuaternion{},
	32:  &LocalPositionNed{},
	33:  &GlobalPositionInt{},
	34:  &RcChannelsScaled{},
	35:  &RcChannelsRaw{},
	36:  &ServoOutputRaw{},
	37:  &MissionRequestPartialList{},
	38:  &MissionWritePartialList{},
	39:  &MissionItem{},
	40:  &MissionRequest{},
	41:  &MissionSetCurrent{},
	42:  &MissionCurrent{},
	43:  &MissionRequestList{},
	44:  &MissionCount{},
	45:  &MissionClearAll{},
	46:  &MissionItemReached{},
	47:  &MissionAck{},
	48:  &SetGpsGlobalOrigin{},
	49:  &GpsGlobalOrigin{},
	50:  &SetLocalPositionSetpoint{},
	51:  &LocalPositionSetpoint{},
	52:  &GlobalPositionSetpointInt{},
	53:  &SetGlobalPositionSetpointInt{},
	54:  &SafetySetAllowedArea{},
	55:  &SafetyAllowedArea{},
	56:  &SetRollPitchYawThrust{},
	57:  &SetRollPitchYawSpeedThrust{},
	58:  &RollPitchYawThrustSetpoint{},
	59:  &RollPitchYawSpeedThrustSetpoint{},
	60:  &SetQuadMotorsSetpoint{},
	61:  &SetQuadSwarmRollPitchYawThrust{},
	62:  &NavControllerOutput{},
	63:  &SetQuadSwarmLedRollPitchYawThrust{},
	64:  &StateCorrection{},
	65:  &RcChannels{},
	66:  &RequestDataStream{},
	67:  &DataStream{},
	69:  &ManualControl{},
	70:  &RcChannelsOverride{},
	74:  &VfrHud{},
	76:  &CommandLong{},
	77:  &CommandAck{},
	80:  &RollPitchYawRatesThrustSetpoint{},
	81:  &ManualSetpoint{},
	82:  &AttitudeSetpointExternal{},
	83:  &LocalNedPositionSetpointExternal{},
	84:  &GlobalPositionSetpointExternalInt{},
	89:  &LocalPositionNedSystemGlobalOffset{},
	90:  &HilState{},
	91:  &HilControls{},
	92:  &HilRcInputsRaw{},
	100: &OpticalFlow{},
	101: &GlobalVisionPositionEstimate{},
	102: &VisionPositionEstimate{},
	103: &VisionSpeedEstimate{},
	104: &ViconPositionEstimate{},
	105: &HighresImu{},
	106: &OmnidirectionalFlow{},
	107: &HilSensor{},
	108: &SimState{},
	109: &RadioStatus{},
	110: &FileTransferStart{},
	111: &FileTransferDirList{},
	112: &FileTransferRes{},
	113: &HilGps{},
	114: &HilOpticalFlow{},
	115: &HilStateQuaternion{},
	116: &ScaledImu2{},
	117: &LogRequestList{},
	118: &LogEntry{},
	119: &LogRequestData{},
	120: &LogData{},
	121: &LogErase{},
	122: &LogRequestEnd{},
	123: &GpsInjectData{},
	124: &Gps2Raw{},
	125: &PowerStatus{},
	126: &SerialControl{},
	127: &GpsRtk{},
	128: &Gps2Rtk{},
	130: &DataTransmissionHandshake{},
	131: &EncapsulatedData{},
	132: &DistanceSensor{},
	133: &TerrainRequest{},
	134: &TerrainData{},
	135: &TerrainCheck{},
	136: &TerrainReport{},
	147: &BatteryStatus{},
	148: &Setpoint8Dof{},
	149: &Setpoint6Dof{},
	249: &MemoryVect{},
	250: &DebugVect{},
	251: &NamedValueFloat{},
	252: &NamedValueInt{},
	253: &Statustext{},
	254: &Debug{},
}

// NewMAVLinkMessage returns a new MAVLinkMessage or an error if it encounters an unknown Message ID
func NewMAVLinkMessage(msgid uint8, data []byte) (MAVLinkMessage, error) {
	message := messages[msgid]
	if message != nil {
		message.Decode(data)
		return message, nil
	}
	return nil, fmt.Errorf("Unknown Message ID: %v", msgid)
}

//
// MAV_AUTOPILOT
/*Micro air vehicle / autopilot classes. This identifies the individual model.*/
//
const (
	MAV_AUTOPILOT_GENERIC                                      = 0  // Generic autopilot, full support for everything |
	MAV_AUTOPILOT_PIXHAWK                                      = 1  // PIXHAWK autopilot, http://pixhawk.ethz.ch |
	MAV_AUTOPILOT_SLUGS                                        = 2  // SLUGS autopilot, http://slugsuav.soe.ucsc.edu |
	MAV_AUTOPILOT_ARDUPILOTMEGA                                = 3  // ArduPilotMega / ArduCopter, http://diydrones.com |
	MAV_AUTOPILOT_OPENPILOT                                    = 4  // OpenPilot, http://openpilot.org |
	MAV_AUTOPILOT_GENERIC_WAYPOINTS_ONLY                       = 5  // Generic autopilot only supporting simple waypoints |
	MAV_AUTOPILOT_GENERIC_WAYPOINTS_AND_SIMPLE_NAVIGATION_ONLY = 6  // Generic autopilot supporting waypoints and other simple navigation commands |
	MAV_AUTOPILOT_GENERIC_MISSION_FULL                         = 7  // Generic autopilot supporting the full mission command set |
	MAV_AUTOPILOT_INVALID                                      = 8  // No valid autopilot, e.g. a GCS or other MAVLink component |
	MAV_AUTOPILOT_PPZ                                          = 9  // PPZ UAV - http://nongnu.org/paparazzi |
	MAV_AUTOPILOT_UDB                                          = 10 // UAV Dev Board |
	MAV_AUTOPILOT_FP                                           = 11 // FlexiPilot |
	MAV_AUTOPILOT_PX4                                          = 12 // PX4 Autopilot - http://pixhawk.ethz.ch/px4/ |
	MAV_AUTOPILOT_SMACCMPILOT                                  = 13 // SMACCMPilot - http://smaccmpilot.org |
	MAV_AUTOPILOT_AUTOQUAD                                     = 14 // AutoQuad -- http://autoquad.org |
	MAV_AUTOPILOT_ARMAZILA                                     = 15 // Armazila -- http://armazila.com |
	MAV_AUTOPILOT_AEROB                                        = 16 // Aerob -- http://aerob.ru |
	MAV_AUTOPILOT_ENUM_END                                     = 17 //  |
)

//
// MAV_TYPE
/**/
//
const (
	MAV_TYPE_GENERIC            = 0  // Generic micro air vehicle. |
	MAV_TYPE_FIXED_WING         = 1  // Fixed wing aircraft. |
	MAV_TYPE_QUADROTOR          = 2  // Quadrotor |
	MAV_TYPE_COAXIAL            = 3  // Coaxial helicopter |
	MAV_TYPE_HELICOPTER         = 4  // Normal helicopter with tail rotor. |
	MAV_TYPE_ANTENNA_TRACKER    = 5  // Ground installation |
	MAV_TYPE_GCS                = 6  // Operator control unit / ground control station |
	MAV_TYPE_AIRSHIP            = 7  // Airship, controlled |
	MAV_TYPE_FREE_BALLOON       = 8  // Free balloon, uncontrolled |
	MAV_TYPE_ROCKET             = 9  // Rocket |
	MAV_TYPE_GROUND_ROVER       = 10 // Ground rover |
	MAV_TYPE_SURFACE_BOAT       = 11 // Surface vessel, boat, ship |
	MAV_TYPE_SUBMARINE          = 12 // Submarine |
	MAV_TYPE_HEXAROTOR          = 13 // Hexarotor |
	MAV_TYPE_OCTOROTOR          = 14 // Octorotor |
	MAV_TYPE_TRICOPTER          = 15 // Octorotor |
	MAV_TYPE_FLAPPING_WING      = 16 // Flapping wing |
	MAV_TYPE_KITE               = 17 // Flapping wing |
	MAV_TYPE_ONBOARD_CONTROLLER = 18 // Onboard companion controller |
	MAV_TYPE_ENUM_END           = 19 //  |
)

//
// MAV_MODE_FLAG
/*These flags encode the MAV mode.*/
//
const (
	MAV_MODE_FLAG_CUSTOM_MODE_ENABLED  = 1   // 0b00000001 Reserved for future use. |
	MAV_MODE_FLAG_TEST_ENABLED         = 2   // 0b00000010 system has a test mode enabled. This flag is intended for temporary system tests and should not be used for stable implementations. |
	MAV_MODE_FLAG_AUTO_ENABLED         = 4   // 0b00000100 autonomous mode enabled, system finds its own goal positions. Guided flag can be set or not, depends on the actual implementation. |
	MAV_MODE_FLAG_GUIDED_ENABLED       = 8   // 0b00001000 guided mode enabled, system flies MISSIONs / mission items. |
	MAV_MODE_FLAG_STABILIZE_ENABLED    = 16  // 0b00010000 system stabilizes electronically its attitude (and optionally position). It needs however further control inputs to move around. |
	MAV_MODE_FLAG_HIL_ENABLED          = 32  // 0b00100000 hardware in the loop simulation. All motors / actuators are blocked, but internal software is full operational. |
	MAV_MODE_FLAG_MANUAL_INPUT_ENABLED = 64  // 0b01000000 remote control input is enabled. |
	MAV_MODE_FLAG_SAFETY_ARMED         = 128 // 0b10000000 MAV safety set to armed. Motors are enabled / running / can start. Ready to fly. |
	MAV_MODE_FLAG_ENUM_END             = 129 //  |
)

//
// MAV_MODE_FLAG_DECODE_POSITION
/*These values encode the bit positions of the decode position. These values can be used to read the value of a flag bit by combining the base_mode variable with AND with the flag position value. The result will be either 0 or 1, depending on if the flag is set or not.*/
//
const (
	MAV_MODE_FLAG_DECODE_POSITION_CUSTOM_MODE = 1   // Eighth bit: 00000001 |
	MAV_MODE_FLAG_DECODE_POSITION_TEST        = 2   // Seventh bit: 00000010 |
	MAV_MODE_FLAG_DECODE_POSITION_AUTO        = 4   // Sixt bit:   00000100 |
	MAV_MODE_FLAG_DECODE_POSITION_GUIDED      = 8   // Fifth bit:  00001000 |
	MAV_MODE_FLAG_DECODE_POSITION_STABILIZE   = 16  // Fourth bit: 00010000 |
	MAV_MODE_FLAG_DECODE_POSITION_HIL         = 32  // Third bit:  00100000 |
	MAV_MODE_FLAG_DECODE_POSITION_MANUAL      = 64  // Second bit: 01000000 |
	MAV_MODE_FLAG_DECODE_POSITION_SAFETY      = 128 // First bit:  10000000 |
	MAV_MODE_FLAG_DECODE_POSITION_ENUM_END    = 129 //  |
)

//
// MAV_GOTO
/*Override command, pauses current mission execution and moves immediately to a position*/
//
const (
	MAV_GOTO_DO_HOLD                    = 0 // Hold at the current position. |
	MAV_GOTO_DO_CONTINUE                = 1 // Continue with the next item in mission execution. |
	MAV_GOTO_HOLD_AT_CURRENT_POSITION   = 2 // Hold at the current position of the system |
	MAV_GOTO_HOLD_AT_SPECIFIED_POSITION = 3 // Hold at the position specified in the parameters of the DO_HOLD action |
	MAV_GOTO_ENUM_END                   = 4 //  |
)

//
// MAV_MODE
/*These defines are predefined OR-combined mode flags. There is no need to use values from this enum, but it
  simplifies the use of the mode flags. Note that manual input is enabled in all modes as a safety override.*/
//
const (
	MAV_MODE_PREFLIGHT          = 0   // System is not ready to fly, booting, calibrating, etc. No flag is set. |
	MAV_MODE_MANUAL_DISARMED    = 64  // System is allowed to be active, under manual (RC) control, no stabilization |
	MAV_MODE_TEST_DISARMED      = 66  // UNDEFINED mode. This solely depends on the autopilot - use with caution, intended for developers only. |
	MAV_MODE_STABILIZE_DISARMED = 80  // System is allowed to be active, under assisted RC control. |
	MAV_MODE_GUIDED_DISARMED    = 88  // System is allowed to be active, under autonomous control, manual setpoint |
	MAV_MODE_AUTO_DISARMED      = 92  // System is allowed to be active, under autonomous control and navigation (the trajectory is decided onboard and not pre-programmed by MISSIONs) |
	MAV_MODE_MANUAL_ARMED       = 192 // System is allowed to be active, under manual (RC) control, no stabilization |
	MAV_MODE_TEST_ARMED         = 194 // UNDEFINED mode. This solely depends on the autopilot - use with caution, intended for developers only. |
	MAV_MODE_STABILIZE_ARMED    = 208 // System is allowed to be active, under assisted RC control. |
	MAV_MODE_GUIDED_ARMED       = 216 // System is allowed to be active, under autonomous control, manual setpoint |
	MAV_MODE_AUTO_ARMED         = 220 // System is allowed to be active, under autonomous control and navigation (the trajectory is decided onboard and not pre-programmed by MISSIONs) |
	MAV_MODE_ENUM_END           = 221 //  |
)

//
// MAV_STATE
/**/
//
const (
	MAV_STATE_UNINIT      = 0 // Uninitialized system, state is unknown. |
	MAV_STATE_BOOT        = 1 // System is booting up. |
	MAV_STATE_CALIBRATING = 2 // System is calibrating and not flight-ready. |
	MAV_STATE_STANDBY     = 3 // System is grounded and on standby. It can be launched any time. |
	MAV_STATE_ACTIVE      = 4 // System is active and might be already airborne. Motors are engaged. |
	MAV_STATE_CRITICAL    = 5 // System is in a non-normal flight mode. It can however still navigate. |
	MAV_STATE_EMERGENCY   = 6 // System is in a non-normal flight mode. It lost control over parts or over the whole airframe. It is in mayday and going down. |
	MAV_STATE_POWEROFF    = 7 // System just initialized its power-down sequence, will shut down now. |
	MAV_STATE_ENUM_END    = 8 //  |
)

//
// MAV_COMPONENT
/**/
//
const (
	MAV_COMP_ID_ALL            = 0   //  |
	MAV_COMP_ID_CAMERA         = 100 //  |
	MAV_COMP_ID_SERVO1         = 140 //  |
	MAV_COMP_ID_SERVO2         = 141 //  |
	MAV_COMP_ID_SERVO3         = 142 //  |
	MAV_COMP_ID_SERVO4         = 143 //  |
	MAV_COMP_ID_SERVO5         = 144 //  |
	MAV_COMP_ID_SERVO6         = 145 //  |
	MAV_COMP_ID_SERVO7         = 146 //  |
	MAV_COMP_ID_SERVO8         = 147 //  |
	MAV_COMP_ID_SERVO9         = 148 //  |
	MAV_COMP_ID_SERVO10        = 149 //  |
	MAV_COMP_ID_SERVO11        = 150 //  |
	MAV_COMP_ID_SERVO12        = 151 //  |
	MAV_COMP_ID_SERVO13        = 152 //  |
	MAV_COMP_ID_SERVO14        = 153 //  |
	MAV_COMP_ID_MAPPER         = 180 //  |
	MAV_COMP_ID_MISSIONPLANNER = 190 //  |
	MAV_COMP_ID_PATHPLANNER    = 195 //  |
	MAV_COMP_ID_IMU            = 200 //  |
	MAV_COMP_ID_IMU_2          = 201 //  |
	MAV_COMP_ID_IMU_3          = 202 //  |
	MAV_COMP_ID_GPS            = 220 //  |
	MAV_COMP_ID_UDP_BRIDGE     = 240 //  |
	MAV_COMP_ID_UART_BRIDGE    = 241 //  |
	MAV_COMP_ID_SYSTEM_CONTROL = 250 //  |
	MAV_COMPONENT_ENUM_END     = 251 //  |
)

//
// MAV_SYS_STATUS_SENSOR
/*These encode the sensors whose status is sent as part of the SYS_STATUS message.*/
//
const (
	MAV_SYS_STATUS_SENSOR_3D_GYRO                = 1       // 0x01 3D gyro |
	MAV_SYS_STATUS_SENSOR_3D_ACCEL               = 2       // 0x02 3D accelerometer |
	MAV_SYS_STATUS_SENSOR_3D_MAG                 = 4       // 0x04 3D magnetometer |
	MAV_SYS_STATUS_SENSOR_ABSOLUTE_PRESSURE      = 8       // 0x08 absolute pressure |
	MAV_SYS_STATUS_SENSOR_DIFFERENTIAL_PRESSURE  = 16      // 0x10 differential pressure |
	MAV_SYS_STATUS_SENSOR_GPS                    = 32      // 0x20 GPS |
	MAV_SYS_STATUS_SENSOR_OPTICAL_FLOW           = 64      // 0x40 optical flow |
	MAV_SYS_STATUS_SENSOR_VISION_POSITION        = 128     // 0x80 computer vision position |
	MAV_SYS_STATUS_SENSOR_LASER_POSITION         = 256     // 0x100 laser based position |
	MAV_SYS_STATUS_SENSOR_EXTERNAL_GROUND_TRUTH  = 512     // 0x200 external ground truth (Vicon or Leica) |
	MAV_SYS_STATUS_SENSOR_ANGULAR_RATE_CONTROL   = 1024    // 0x400 3D angular rate control |
	MAV_SYS_STATUS_SENSOR_ATTITUDE_STABILIZATION = 2048    // 0x800 attitude stabilization |
	MAV_SYS_STATUS_SENSOR_YAW_POSITION           = 4096    // 0x1000 yaw position |
	MAV_SYS_STATUS_SENSOR_Z_ALTITUDE_CONTROL     = 8192    // 0x2000 z/altitude control |
	MAV_SYS_STATUS_SENSOR_XY_POSITION_CONTROL    = 16384   // 0x4000 x/y position control |
	MAV_SYS_STATUS_SENSOR_MOTOR_OUTPUTS          = 32768   // 0x8000 motor outputs / control |
	MAV_SYS_STATUS_SENSOR_RC_RECEIVER            = 65536   // 0x10000 rc receiver |
	MAV_SYS_STATUS_SENSOR_3D_GYRO2               = 131072  // 0x20000 2nd 3D gyro |
	MAV_SYS_STATUS_SENSOR_3D_ACCEL2              = 262144  // 0x40000 2nd 3D accelerometer |
	MAV_SYS_STATUS_SENSOR_3D_MAG2                = 524288  // 0x80000 2nd 3D magnetometer |
	MAV_SYS_STATUS_GEOFENCE                      = 1048576 // 0x100000 geofence |
	MAV_SYS_STATUS_AHRS                          = 2097152 // 0x200000 AHRS subsystem health |
	MAV_SYS_STATUS_TERRAIN                       = 4194304 // 0x400000 Terrain subsystem health |
	MAV_SYS_STATUS_SENSOR_ENUM_END               = 4194305 //  |
)

//
// MAV_FRAME
/**/
//
const (
	MAV_FRAME_GLOBAL                  = 0  // Global coordinate frame, WGS84 coordinate system. First value / x: latitude, second value / y: longitude, third value / z: positive altitude over mean sea level (MSL) |
	MAV_FRAME_LOCAL_NED               = 1  // Local coordinate frame, Z-up (x: north, y: east, z: down). |
	MAV_FRAME_MISSION                 = 2  // NOT a coordinate frame, indicates a mission command. |
	MAV_FRAME_GLOBAL_RELATIVE_ALT     = 3  // Global coordinate frame, WGS84 coordinate system, relative altitude over ground with respect to the home position. First value / x: latitude, second value / y: longitude, third value / z: positive altitude with 0 being at the altitude of the home location. |
	MAV_FRAME_LOCAL_ENU               = 4  // Local coordinate frame, Z-down (x: east, y: north, z: up) |
	MAV_FRAME_GLOBAL_INT              = 5  // Global coordinate frame with some fields as scaled integers, WGS84 coordinate system. First value / x: latitude, second value / y: longitude, third value / z: positive altitude over mean sea level (MSL). Lat / Lon are scaled * 1E7 to avoid floating point accuracy limitations. |
	MAV_FRAME_GLOBAL_RELATIVE_ALT_INT = 6  // Global coordinate frame with some fields as scaled integers, WGS84 coordinate system, relative altitude over ground with respect to the home position. First value / x: latitude, second value / y: longitude, third value / z: positive altitude with 0 being at the altitude of the home location. Lat / Lon are scaled * 1E7 to avoid floating point accuracy limitations. |
	MAV_FRAME_LOCAL_OFFSET_NED        = 7  // Offset to the current local frame. Anything expressed in this frame should be added to the current local frame position. |
	MAV_FRAME_BODY_NED                = 8  // Setpoint in body NED frame. This makes sense if all position control is externalized - e.g. useful to command 2 m/s^2 acceleration to the right. |
	MAV_FRAME_BODY_OFFSET_NED         = 9  // Offset in body NED frame. This makes sense if adding setpoints to the current flight path, to avoid an obstacle - e.g. useful to command 2 m/s^2 acceleration to the east. |
	MAV_FRAME_GLOBAL_TERRAIN_ALT      = 10 // Global coordinate frame with above terrain level altitude. WGS84 coordinate system, relative altitude over terrain with respect to the waypoint coordinate. First value / x: latitude, second value / y: longitude, third value / z: positive altitude with 0 being at ground level in terrain model. |
	MAV_FRAME_ENUM_END                = 11 //  |
)

//
// MAVLINK_DATA_STREAM_TYPE
/**/
//
const (
	MAVLINK_DATA_STREAM_IMG_JPEG      = 1 //  |
	MAVLINK_DATA_STREAM_IMG_BMP       = 2 //  |
	MAVLINK_DATA_STREAM_IMG_RAW8U     = 3 //  |
	MAVLINK_DATA_STREAM_IMG_RAW32U    = 4 //  |
	MAVLINK_DATA_STREAM_IMG_PGM       = 5 //  |
	MAVLINK_DATA_STREAM_IMG_PNG       = 6 //  |
	MAVLINK_DATA_STREAM_TYPE_ENUM_END = 7 //  |
)

//
// FENCE_ACTION
/**/
//
const (
	FENCE_ACTION_NONE            = 0 // Disable fenced mode |
	FENCE_ACTION_GUIDED          = 1 // Switched to guided mode to return point (fence point 0) |
	FENCE_ACTION_REPORT          = 2 // Report fence breach, but don't take action |
	FENCE_ACTION_GUIDED_THR_PASS = 3 // Switched to guided mode to return point (fence point 0) with manual throttle control |
	FENCE_ACTION_ENUM_END        = 4 //  |
)

//
// FENCE_BREACH
/**/
//
const (
	FENCE_BREACH_NONE     = 0 // No last fence breach |
	FENCE_BREACH_MINALT   = 1 // Breached minimum altitude |
	FENCE_BREACH_MAXALT   = 2 // Breached maximum altitude |
	FENCE_BREACH_BOUNDARY = 3 // Breached fence boundary |
	FENCE_BREACH_ENUM_END = 4 //  |
)

//
// MAV_MOUNT_MODE
/*Enumeration of possible mount operation modes*/
//
const (
	MAV_MOUNT_MODE_RETRACT           = 0 // Load and keep safe position (Roll,Pitch,Yaw) from permant memory and stop stabilization |
	MAV_MOUNT_MODE_NEUTRAL           = 1 // Load and keep neutral position (Roll,Pitch,Yaw) from permanent memory. |
	MAV_MOUNT_MODE_MAVLINK_TARGETING = 2 // Load neutral position and start MAVLink Roll,Pitch,Yaw control with stabilization |
	MAV_MOUNT_MODE_RC_TARGETING      = 3 // Load neutral position and start RC Roll,Pitch,Yaw control with stabilization |
	MAV_MOUNT_MODE_GPS_POINT         = 4 // Load neutral position and start to point to Lat,Lon,Alt |
	MAV_MOUNT_MODE_ENUM_END          = 5 //  |
)

//
// MAV_CMD
/*Commands to be executed by the MAV. They can be executed on user request, or as part of a mission script. If the action is used in a mission, the parameter mapping to the waypoint/mission message is as follows: Param 1, Param 2, Param 3, Param 4, X: Param 5, Y:Param 6, Z:Param 7. This command list is similar what ARINC 424 is for commercial aircraft: A data format how to interpret waypoint/mission data.*/
//
const (
	MAV_CMD_NAV_WAYPOINT                 = 16  // Navigate to MISSION. | Hold time in decimal seconds. (ignored by fixed wing, time to stay at MISSION for rotary wing) | Acceptance radius in meters (if the sphere with this radius is hit, the MISSION counts as reached) | 0 to pass through the WP, if > 0 radius in meters to pass by WP. Positive value for clockwise orbit, negative value for counter-clockwise orbit. Allows trajectory control. | Desired yaw angle at MISSION (rotary wing) | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_LOITER_UNLIM             = 17  // Loiter around this MISSION an unlimited amount of time | Empty | Empty | Radius around MISSION, in meters. If positive loiter clockwise, else counter-clockwise | Desired yaw angle. | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_LOITER_TURNS             = 18  // Loiter around this MISSION for X turns | Turns | Empty | Radius around MISSION, in meters. If positive loiter clockwise, else counter-clockwise | Desired yaw angle. | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_LOITER_TIME              = 19  // Loiter around this MISSION for X seconds | Seconds (decimal) | Empty | Radius around MISSION, in meters. If positive loiter clockwise, else counter-clockwise | Desired yaw angle. | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_RETURN_TO_LAUNCH         = 20  // Return to launch location | Empty | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_NAV_LAND                     = 21  // Land at location | Empty | Empty | Empty | Desired yaw angle. | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_TAKEOFF                  = 22  // Takeoff from ground / hand | Minimum pitch (if airspeed sensor present), desired pitch without sensor | Empty | Empty | Yaw angle (if magnetometer present), ignored without magnetometer | Latitude | Longitude | Altitude |
	MAV_CMD_NAV_ROI                      = 80  // Sets the region of interest (ROI) for a sensor set or the vehicle itself. This can then be used by the vehicles control system to control the vehicle attitude and the attitude of various sensors such as cameras. | Region of intereset mode. (see MAV_ROI enum) | MISSION index/ target ID. (see MAV_ROI enum) | ROI index (allows a vehicle to manage multiple ROI's) | Empty | x the location of the fixed ROI (see MAV_FRAME) | y | z |
	MAV_CMD_NAV_PATHPLANNING             = 81  // Control autonomous path planning on the MAV. | 0: Disable local obstacle avoidance / local path planning (without resetting map), 1: Enable local path planning, 2: Enable and reset local path planning | 0: Disable full path planning (without resetting map), 1: Enable, 2: Enable and reset map/occupancy grid, 3: Enable and reset planned route, but not occupancy grid | Empty | Yaw angle at goal, in compass degrees, [0..360] | Latitude/X of goal | Longitude/Y of goal | Altitude/Z of goal |
	MAV_CMD_NAV_SPLINE_WAYPOINT          = 82  // Navigate to MISSION using a spline path. | Hold time in decimal seconds. (ignored by fixed wing, time to stay at MISSION for rotary wing) | Empty | Empty | Empty | Latitude/X of goal | Longitude/Y of goal | Altitude/Z of goal |
	MAV_CMD_NAV_GUIDED_ENABLE            = 92  // hand control over to an external controller | On / Off (> 0.5f on) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_NAV_LAST                     = 95  // NOP - This command is only used to mark the upper limit of the NAV/ACTION commands in the enumeration | Empty | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_CONDITION_DELAY              = 112 // Delay mission state machine. | Delay in seconds (decimal) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_CONDITION_CHANGE_ALT         = 113 // Ascend/descend at rate.  Delay mission state machine until desired altitude reached. | Descent / Ascend rate (m/s) | Empty | Empty | Empty | Empty | Empty | Finish Altitude |
	MAV_CMD_CONDITION_DISTANCE           = 114 // Delay mission state machine until within desired distance of next NAV point. | Distance (meters) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_CONDITION_YAW                = 115 // Reach a certain target angle. | target angle: [0-360], 0 is north | speed during yaw change:[deg per second] | direction: negative: counter clockwise, positive: clockwise [-1,1] | relative offset or absolute angle: [ 1,0] | Empty | Empty | Empty |
	MAV_CMD_CONDITION_LAST               = 159 // NOP - This command is only used to mark the upper limit of the CONDITION commands in the enumeration | Empty | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_SET_MODE                  = 176 // Set system mode. | Mode, as defined by ENUM MAV_MODE | Custom mode - this is system specific, please refer to the individual autopilot specifications for details. | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_JUMP                      = 177 // Jump to the desired command in the mission list.  Repeat this action only the specified number of times | Sequence number | Repeat count | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_CHANGE_SPEED              = 178 // Change speed and/or throttle set points. | Speed type (0=Airspeed, 1=Ground Speed) | Speed  (m/s, -1 indicates no change) | Throttle  ( Percent, -1 indicates no change) | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_SET_HOME                  = 179 // Changes the home location either to the current location or a specified location. | Use current (1=use current location, 0=use specified location) | Empty | Empty | Empty | Latitude | Longitude | Altitude |
	MAV_CMD_DO_SET_PARAMETER             = 180 // Set a system parameter.  Caution!  Use of this command requires knowledge of the numeric enumeration value of the parameter. | Parameter number | Parameter value | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_SET_RELAY                 = 181 // Set a relay to a condition. | Relay number | Setting (1=on, 0=off, others possible depending on system hardware) | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_REPEAT_RELAY              = 182 // Cycle a relay on and off for a desired number of cyles with a desired period. | Relay number | Cycle count | Cycle time (seconds, decimal) | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_SET_SERVO                 = 183 // Set a servo to a desired PWM value. | Servo number | PWM (microseconds, 1000 to 2000 typical) | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_REPEAT_SERVO              = 184 // Cycle a between its nominal setting and a desired PWM for a desired number of cycles with a desired period. | Servo number | PWM (microseconds, 1000 to 2000 typical) | Cycle count | Cycle time (seconds) | Empty | Empty | Empty |
	MAV_CMD_DO_FLIGHTTERMINATION         = 185 // Terminate flight immediately | Flight termination activated if > 0.5 | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_RALLY_LAND                = 190 // Mission command to perform a landing from a rally point. | Break altitude (meters) | Landing speed (m/s) | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_GO_AROUND                 = 191 // Mission command to safely abort an autonmous landing. | Altitude (meters) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_CONTROL_VIDEO             = 200 // Control onboard camera system. | Camera ID (-1 for all) | Transmission: 0: disabled, 1: enabled compressed, 2: enabled raw | Transmission mode: 0: video stream, >0: single images every n seconds (decimal) | Recording: 0: disabled, 1: enabled compressed, 2: enabled raw | Empty | Empty | Empty |
	MAV_CMD_DO_SET_ROI                   = 201 // Sets the region of interest (ROI) for a sensor set or the vehicle itself. This can then be used by the vehicles control system to control the vehicle attitude and the attitude of various sensors such as cameras. | Region of intereset mode. (see MAV_ROI enum) | MISSION index/ target ID. (see MAV_ROI enum) | ROI index (allows a vehicle to manage multiple ROI's) | Empty | x the location of the fixed ROI (see MAV_FRAME) | y | z |
	MAV_CMD_DO_DIGICAM_CONFIGURE         = 202 // Mission command to configure an on-board camera controller system. | Modes: P, TV, AV, M, Etc | Shutter speed: Divisor number for one second | Aperture: F stop number | ISO number e.g. 80, 100, 200, Etc | Exposure type enumerator | Command Identity | Main engine cut-off time before camera trigger in seconds/10 (0 means no cut-off) |
	MAV_CMD_DO_DIGICAM_CONTROL           = 203 // Mission command to control an on-board camera controller system. | Session control e.g. show/hide lens | Zoom's absolute position | Zooming step value to offset zoom from the current position | Focus Locking, Unlocking or Re-locking | Shooting Command | Command Identity | Empty |
	MAV_CMD_DO_MOUNT_CONFIGURE           = 204 // Mission command to configure a camera or antenna mount | Mount operation mode (see MAV_MOUNT_MODE enum) | stabilize roll? (1 = yes, 0 = no) | stabilize pitch? (1 = yes, 0 = no) | stabilize yaw? (1 = yes, 0 = no) | Empty | Empty | Empty |
	MAV_CMD_DO_MOUNT_CONTROL             = 205 // Mission command to control a camera or antenna mount | pitch or lat in degrees, depending on mount mode. | roll or lon in degrees depending on mount mode | yaw or alt (in meters) depending on mount mode | reserved | reserved | reserved | MAV_MOUNT_MODE enum value |
	MAV_CMD_DO_SET_CAM_TRIGG_DIST        = 206 // Mission command to set CAM_TRIGG_DIST for this flight | Camera trigger distance (meters) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_FENCE_ENABLE              = 207 // Mission command to enable the geofence | enable? (0=disable, 1=enable) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_PARACHUTE                 = 208 // Mission command to trigger a parachute | action (0=disable, 1=enable, 2=release, for some systems see PARACHUTE_ACTION enum, not in general message set.) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_INVERTED_FLIGHT           = 210 // Change to/from inverted flight | inverted (0=normal, 1=inverted) | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_MOUNT_CONTROL_QUAT        = 220 // Mission command to control a camera or antenna mount, using a quaternion as reference. | q1 - quaternion param #1, w (1 in null-rotation) | q2 - quaternion param #2, x (0 in null-rotation) | q3 - quaternion param #3, y (0 in null-rotation) | q4 - quaternion param #4, z (0 in null-rotation) | Empty | Empty | Empty |
	MAV_CMD_DO_GUIDED_MASTER             = 221 // set id of master controller | System ID | Component ID | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_DO_GUIDED_LIMITS             = 222 // set limits for external control | timeout - maximum time (in seconds) that external controller will be allowed to control vehicle. 0 means no timeout | absolute altitude min (in meters, WGS84) - if vehicle moves below this alt, the command will be aborted and the mission will continue.  0 means no lower altitude limit | absolute altitude max (in meters)- if vehicle moves above this alt, the command will be aborted and the mission will continue.  0 means no upper altitude limit | horizontal move limit (in meters, WGS84) - if vehicle moves more than this distance from it's location at the moment the command was executed, the command will be aborted and the mission will continue. 0 means no horizontal altitude limit | Empty | Empty | Empty |
	MAV_CMD_DO_LAST                      = 240 // NOP - This command is only used to mark the upper limit of the DO commands in the enumeration | Empty | Empty | Empty | Empty | Empty | Empty | Empty |
	MAV_CMD_PREFLIGHT_CALIBRATION        = 241 // Trigger calibration. This command will be only accepted if in pre-flight mode. | Gyro calibration: 0: no, 1: yes | Magnetometer calibration: 0: no, 1: yes | Ground pressure: 0: no, 1: yes | Radio calibration: 0: no, 1: yes | Accelerometer calibration: 0: no, 1: yes | Compass/Motor interference calibration: 0: no, 1: yes | Empty |
	MAV_CMD_PREFLIGHT_SET_SENSOR_OFFSETS = 242 // Set sensor offsets. This command will be only accepted if in pre-flight mode. | Sensor to adjust the offsets for: 0: gyros, 1: accelerometer, 2: magnetometer, 3: barometer, 4: optical flow, 5: second magnetometer | X axis offset (or generic dimension 1), in the sensor's raw units | Y axis offset (or generic dimension 2), in the sensor's raw units | Z axis offset (or generic dimension 3), in the sensor's raw units | Generic dimension 4, in the sensor's raw units | Generic dimension 5, in the sensor's raw units | Generic dimension 6, in the sensor's raw units |
	MAV_CMD_PREFLIGHT_STORAGE            = 245 // Request storage of different parameter values and logs. This command will be only accepted if in pre-flight mode. | Parameter storage: 0: READ FROM FLASH/EEPROM, 1: WRITE CURRENT TO FLASH/EEPROM | Mission storage: 0: READ FROM FLASH/EEPROM, 1: WRITE CURRENT TO FLASH/EEPROM | Reserved | Reserved | Empty | Empty | Empty |
	MAV_CMD_PREFLIGHT_REBOOT_SHUTDOWN    = 246 // Request the reboot or shutdown of system components. | 0: Do nothing for autopilot, 1: Reboot autopilot, 2: Shutdown autopilot. | 0: Do nothing for onboard computer, 1: Reboot onboard computer, 2: Shutdown onboard computer. | Reserved | Reserved | Empty | Empty | Empty |
	MAV_CMD_OVERRIDE_GOTO                = 252 // Hold / continue the current action | MAV_GOTO_DO_HOLD: hold MAV_GOTO_DO_CONTINUE: continue with next item in mission plan | MAV_GOTO_HOLD_AT_CURRENT_POSITION: Hold at current position MAV_GOTO_HOLD_AT_SPECIFIED_POSITION: hold at specified position | MAV_FRAME coordinate frame of hold point | Desired yaw angle in degrees | Latitude / X position | Longitude / Y position | Altitude / Z position |
	MAV_CMD_MISSION_START                = 300 // start running a mission | first_item: the first mission item to run | last_item:  the last mission item to run (after this item is run, the mission ends) |
	MAV_CMD_COMPONENT_ARM_DISARM         = 400 // Arms / Disarms a component | 1 to arm, 0 to disarm |
	MAV_CMD_START_RX_PAIR                = 500 // Starts receiver pairing | 0:Spektrum | 0:Spektrum DSM2, 1:Spektrum DSMX |
	MAV_CMD_ENUM_END                     = 501 //  |
)

//
// MAV_DATA_STREAM
/*Data stream IDs. A data stream is not a fixed set of messages, but rather a
  recommendation to the autopilot software. Individual autopilots may or may not obey
  the recommended messages.*/
//
const (
	MAV_DATA_STREAM_ALL             = 0  // Enable all data streams |
	MAV_DATA_STREAM_RAW_SENSORS     = 1  // Enable IMU_RAW, GPS_RAW, GPS_STATUS packets. |
	MAV_DATA_STREAM_EXTENDED_STATUS = 2  // Enable GPS_STATUS, CONTROL_STATUS, AUX_STATUS |
	MAV_DATA_STREAM_RC_CHANNELS     = 3  // Enable RC_CHANNELS_SCALED, RC_CHANNELS_RAW, SERVO_OUTPUT_RAW |
	MAV_DATA_STREAM_RAW_CONTROLLER  = 4  // Enable ATTITUDE_CONTROLLER_OUTPUT, POSITION_CONTROLLER_OUTPUT, NAV_CONTROLLER_OUTPUT. |
	MAV_DATA_STREAM_POSITION        = 6  // Enable LOCAL_POSITION, GLOBAL_POSITION/GLOBAL_POSITION_INT messages. |
	MAV_DATA_STREAM_EXTRA1          = 10 // Dependent on the autopilot |
	MAV_DATA_STREAM_EXTRA2          = 11 // Dependent on the autopilot |
	MAV_DATA_STREAM_EXTRA3          = 12 // Dependent on the autopilot |
	MAV_DATA_STREAM_ENUM_END        = 13 //  |
)

//
// MAV_ROI
/* The ROI (region of interest) for the vehicle. This can be
   be used by the vehicle for camera/vehicle attitude alignment (see
   MAV_CMD_NAV_ROI).*/
//
const (
	MAV_ROI_NONE     = 0 // No region of interest. |
	MAV_ROI_WPNEXT   = 1 // Point toward next MISSION. |
	MAV_ROI_WPINDEX  = 2 // Point toward given MISSION. |
	MAV_ROI_LOCATION = 3 // Point toward fixed location. |
	MAV_ROI_TARGET   = 4 // Point toward of given id. |
	MAV_ROI_ENUM_END = 5 //  |
)

//
// MAV_CMD_ACK
/*ACK / NACK / ERROR values as a result of MAV_CMDs and for mission item transmission.*/
//
const (
	MAV_CMD_ACK_OK                                 = 1  // Command / mission item is ok. |
	MAV_CMD_ACK_ERR_FAIL                           = 2  // Generic error message if none of the other reasons fails or if no detailed error reporting is implemented. |
	MAV_CMD_ACK_ERR_ACCESS_DENIED                  = 3  // The system is refusing to accept this command from this source / communication partner. |
	MAV_CMD_ACK_ERR_NOT_SUPPORTED                  = 4  // Command or mission item is not supported, other commands would be accepted. |
	MAV_CMD_ACK_ERR_COORDINATE_FRAME_NOT_SUPPORTED = 5  // The coordinate frame of this command / mission item is not supported. |
	MAV_CMD_ACK_ERR_COORDINATES_OUT_OF_RANGE       = 6  // The coordinate frame of this command is ok, but he coordinate values exceed the safety limits of this system. This is a generic error, please use the more specific error messages below if possible. |
	MAV_CMD_ACK_ERR_X_LAT_OUT_OF_RANGE             = 7  // The X or latitude value is out of range. |
	MAV_CMD_ACK_ERR_Y_LON_OUT_OF_RANGE             = 8  // The Y or longitude value is out of range. |
	MAV_CMD_ACK_ERR_Z_ALT_OUT_OF_RANGE             = 9  // The Z or altitude value is out of range. |
	MAV_CMD_ACK_ENUM_END                           = 10 //  |
)

//
// MAV_PARAM_TYPE
/*Specifies the datatype of a MAVLink parameter.*/
//
const (
	MAV_PARAM_TYPE_UINT8    = 1  // 8-bit unsigned integer |
	MAV_PARAM_TYPE_INT8     = 2  // 8-bit signed integer |
	MAV_PARAM_TYPE_UINT16   = 3  // 16-bit unsigned integer |
	MAV_PARAM_TYPE_INT16    = 4  // 16-bit signed integer |
	MAV_PARAM_TYPE_UINT32   = 5  // 32-bit unsigned integer |
	MAV_PARAM_TYPE_INT32    = 6  // 32-bit signed integer |
	MAV_PARAM_TYPE_UINT64   = 7  // 64-bit unsigned integer |
	MAV_PARAM_TYPE_INT64    = 8  // 64-bit signed integer |
	MAV_PARAM_TYPE_REAL32   = 9  // 32-bit floating-point |
	MAV_PARAM_TYPE_REAL64   = 10 // 64-bit floating-point |
	MAV_PARAM_TYPE_ENUM_END = 11 //  |
)

//
// MAV_RESULT
/*result from a mavlink command*/
//
const (
	MAV_RESULT_ACCEPTED             = 0 // Command ACCEPTED and EXECUTED |
	MAV_RESULT_TEMPORARILY_REJECTED = 1 // Command TEMPORARY REJECTED/DENIED |
	MAV_RESULT_DENIED               = 2 // Command PERMANENTLY DENIED |
	MAV_RESULT_UNSUPPORTED          = 3 // Command UNKNOWN/UNSUPPORTED |
	MAV_RESULT_FAILED               = 4 // Command executed, but failed |
	MAV_RESULT_ENUM_END             = 5 //  |
)

//
// MAV_MISSION_RESULT
/*result in a mavlink mission ack*/
//
const (
	MAV_MISSION_ACCEPTED          = 0  // mission accepted OK |
	MAV_MISSION_ERROR             = 1  // generic error / not accepting mission commands at all right now |
	MAV_MISSION_UNSUPPORTED_FRAME = 2  // coordinate frame is not supported |
	MAV_MISSION_UNSUPPORTED       = 3  // command is not supported |
	MAV_MISSION_NO_SPACE          = 4  // mission item exceeds storage space |
	MAV_MISSION_INVALID           = 5  // one of the parameters has an invalid value |
	MAV_MISSION_INVALID_PARAM1    = 6  // param1 has an invalid value |
	MAV_MISSION_INVALID_PARAM2    = 7  // param2 has an invalid value |
	MAV_MISSION_INVALID_PARAM3    = 8  // param3 has an invalid value |
	MAV_MISSION_INVALID_PARAM4    = 9  // param4 has an invalid value |
	MAV_MISSION_INVALID_PARAM5_X  = 10 // x/param5 has an invalid value |
	MAV_MISSION_INVALID_PARAM6_Y  = 11 // y/param6 has an invalid value |
	MAV_MISSION_INVALID_PARAM7    = 12 // param7 has an invalid value |
	MAV_MISSION_INVALID_SEQUENCE  = 13 // received waypoint out of sequence |
	MAV_MISSION_DENIED            = 14 // not accepting any mission commands from this communication partner |
	MAV_MISSION_RESULT_ENUM_END   = 15 //  |
)

//
// MAV_SEVERITY
/*Indicates the severity level, generally used for status messages to indicate their relative urgency. Based on RFC-5424 using expanded definitions at: http://www.kiwisyslog.com/kb/info:-syslog-message-levels/.*/
//
const (
	MAV_SEVERITY_EMERGENCY = 0 // System is unusable. This is a "panic" condition. |
	MAV_SEVERITY_ALERT     = 1 // Action should be taken immediately. Indicates error in non-critical systems. |
	MAV_SEVERITY_CRITICAL  = 2 // Action must be taken immediately. Indicates failure in a primary system. |
	MAV_SEVERITY_ERROR     = 3 // Indicates an error in secondary/redundant systems. |
	MAV_SEVERITY_WARNING   = 4 // Indicates about a possible future error if this is not resolved within a given timeframe. Example would be a low battery warning. |
	MAV_SEVERITY_NOTICE    = 5 // An unusual event has occurred, though not an error condition. This should be investigated for the root cause. |
	MAV_SEVERITY_INFO      = 6 // Normal operational messages. Useful for logging. No action is required for these messages. |
	MAV_SEVERITY_DEBUG     = 7 // Useful non-operational messages that can assist in debugging. These should not occur during normal operation. |
	MAV_SEVERITY_ENUM_END  = 8 //  |
)

//
// MAV_POWER_STATUS
/*Power supply status flags (bitmask)*/
//
const (
	MAV_POWER_STATUS_BRICK_VALID                = 1  // main brick power supply valid |
	MAV_POWER_STATUS_SERVO_VALID                = 2  // main servo power supply valid for FMU |
	MAV_POWER_STATUS_USB_CONNECTED              = 4  // USB power is connected |
	MAV_POWER_STATUS_PERIPH_OVERCURRENT         = 8  // peripheral supply is in over-current state |
	MAV_POWER_STATUS_PERIPH_HIPOWER_OVERCURRENT = 16 // hi-power peripheral supply is in over-current state |
	MAV_POWER_STATUS_CHANGED                    = 32 // Power status has changed since boot |
	MAV_POWER_STATUS_ENUM_END                   = 33 //  |
)

//
// SERIAL_CONTROL_DEV
/*SERIAL_CONTROL device types*/
//
const (
	SERIAL_CONTROL_DEV_TELEM1   = 0 // First telemetry port |
	SERIAL_CONTROL_DEV_TELEM2   = 1 // Second telemetry port |
	SERIAL_CONTROL_DEV_GPS1     = 2 // First GPS port |
	SERIAL_CONTROL_DEV_GPS2     = 3 // Second GPS port |
	SERIAL_CONTROL_DEV_ENUM_END = 4 //  |
)

//
// SERIAL_CONTROL_FLAG
/*SERIAL_CONTROL flags (bitmask)*/
//
const (
	SERIAL_CONTROL_FLAG_REPLY     = 1  // Set if this is a reply |
	SERIAL_CONTROL_FLAG_RESPOND   = 2  // Set if the sender wants the receiver to send a response as another SERIAL_CONTROL message |
	SERIAL_CONTROL_FLAG_EXCLUSIVE = 4  // Set if access to the serial port should be removed from whatever driver is currently using it, giving exclusive access to the SERIAL_CONTROL protocol. The port can be handed back by sending a request without this flag set |
	SERIAL_CONTROL_FLAG_BLOCKING  = 8  // Block on writes to the serial port |
	SERIAL_CONTROL_FLAG_MULTI     = 16 // Send multiple replies until port is drained |
	SERIAL_CONTROL_FLAG_ENUM_END  = 17 //  |
)

//
// MAV_DISTANCE_SENSOR
/*Enumeration of distance sensor types*/
//
const (
	MAV_DISTANCE_SENSOR_LASER      = 0 // Laser altimeter, e.g. LightWare SF02/F or PulsedLight units |
	MAV_DISTANCE_SENSOR_ULTRASOUND = 1 // Ultrasound altimeter, e.g. MaxBotix units |
	MAV_DISTANCE_SENSOR_ENUM_END   = 2 //  |
)

//
// MESSAGE HEARTBEAT
//
// MAVLINK_MSG_ID_HEARTBEAT 0
//
// MAVLINK_MSG_ID_HEARTBEAT_LEN 9
//
// MAVLINK_MSG_ID_HEARTBEAT_CRC 50
//
//
type Heartbeat struct {
	CUSTOM_MODE     uint32 // A bitfield for use for autopilot-specific flags.
	TYPE            uint8  // Type of the MAV (quadrotor, helicopter, etc., up to 15 types, defined in MAV_TYPE ENUM)
	AUTOPILOT       uint8  // Autopilot type / class. defined in MAV_AUTOPILOT ENUM
	BASE_MODE       uint8  // System mode bitfield, see MAV_MODE_FLAG ENUM in mavlink/include/mavlink_types.h
	SYSTEM_STATUS   uint8  // System status flag, see MAV_STATE ENUM
	MAVLINK_VERSION uint8  // MAVLink version, not writable by user, gets added by protocol because of magic data type: uint8_t_mavlink_version
}

// NewHeartbeat returns a new Heartbeat
func NewHeartbeat(CUSTOM_MODE uint32, TYPE uint8, AUTOPILOT uint8, BASE_MODE uint8, SYSTEM_STATUS uint8, MAVLINK_VERSION uint8) *Heartbeat {
	m := Heartbeat{}
	m.CUSTOM_MODE = CUSTOM_MODE
	m.TYPE = TYPE
	m.AUTOPILOT = AUTOPILOT
	m.BASE_MODE = BASE_MODE
	m.SYSTEM_STATUS = SYSTEM_STATUS
	m.MAVLINK_VERSION = MAVLINK_VERSION
	return &m
}

// Id returns the Heartbeat Message ID
func (*Heartbeat) Id() uint8 {
	return 0
}

// Len returns the Heartbeat Message Length
func (*Heartbeat) Len() uint8 {
	return 9
}

// Crc returns the Heartbeat Message CRC
func (*Heartbeat) Crc() uint8 {
	return 50
}

// Pack returns a packed byte array which represents a Heartbeat payload
func (m *Heartbeat) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.CUSTOM_MODE)
	binary.Write(data, binary.LittleEndian, m.TYPE)
	binary.Write(data, binary.LittleEndian, m.AUTOPILOT)
	binary.Write(data, binary.LittleEndian, m.BASE_MODE)
	binary.Write(data, binary.LittleEndian, m.SYSTEM_STATUS)
	binary.Write(data, binary.LittleEndian, m.MAVLINK_VERSION)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Heartbeat
func (m *Heartbeat) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.CUSTOM_MODE)
	binary.Read(data, binary.LittleEndian, &m.TYPE)
	binary.Read(data, binary.LittleEndian, &m.AUTOPILOT)
	binary.Read(data, binary.LittleEndian, &m.BASE_MODE)
	binary.Read(data, binary.LittleEndian, &m.SYSTEM_STATUS)
	binary.Read(data, binary.LittleEndian, &m.MAVLINK_VERSION)
}

//
// MESSAGE SYS_STATUS
//
// MAVLINK_MSG_ID_SYS_STATUS 1
//
// MAVLINK_MSG_ID_SYS_STATUS_LEN 31
//
// MAVLINK_MSG_ID_SYS_STATUS_CRC 124
//
//
type SysStatus struct {
	ONBOARD_CONTROL_SENSORS_PRESENT uint32 // Bitmask showing which onboard controllers and sensors are present. Value of 0: not present. Value of 1: present. Indices defined by ENUM MAV_SYS_STATUS_SENSOR
	ONBOARD_CONTROL_SENSORS_ENABLED uint32 // Bitmask showing which onboard controllers and sensors are enabled:  Value of 0: not enabled. Value of 1: enabled. Indices defined by ENUM MAV_SYS_STATUS_SENSOR
	ONBOARD_CONTROL_SENSORS_HEALTH  uint32 // Bitmask showing which onboard controllers and sensors are operational or have an error:  Value of 0: not enabled. Value of 1: enabled. Indices defined by ENUM MAV_SYS_STATUS_SENSOR
	LOAD                            uint16 // Maximum usage in percent of the mainloop time, (0%: 0, 100%: 1000) should be always below 1000
	VOLTAGE_BATTERY                 uint16 // Battery voltage, in millivolts (1 = 1 millivolt)
	CURRENT_BATTERY                 int16  // Battery current, in 10*milliamperes (1 = 10 milliampere), -1: autopilot does not measure the current
	DROP_RATE_COMM                  uint16 // Communication drops in percent, (0%: 0, 100%: 10'000), (UART, I2C, SPI, CAN), dropped packets on all links (packets that were corrupted on reception on the MAV)
	ERRORS_COMM                     uint16 // Communication errors (UART, I2C, SPI, CAN), dropped packets on all links (packets that were corrupted on reception on the MAV)
	ERRORS_COUNT1                   uint16 // Autopilot-specific errors
	ERRORS_COUNT2                   uint16 // Autopilot-specific errors
	ERRORS_COUNT3                   uint16 // Autopilot-specific errors
	ERRORS_COUNT4                   uint16 // Autopilot-specific errors
	BATTERY_REMAINING               int8   // Remaining battery energy: (0%: 0, 100%: 100), -1: autopilot estimate the remaining battery
}

// NewSysStatus returns a new SysStatus
func NewSysStatus(ONBOARD_CONTROL_SENSORS_PRESENT uint32, ONBOARD_CONTROL_SENSORS_ENABLED uint32, ONBOARD_CONTROL_SENSORS_HEALTH uint32, LOAD uint16, VOLTAGE_BATTERY uint16, CURRENT_BATTERY int16, DROP_RATE_COMM uint16, ERRORS_COMM uint16, ERRORS_COUNT1 uint16, ERRORS_COUNT2 uint16, ERRORS_COUNT3 uint16, ERRORS_COUNT4 uint16, BATTERY_REMAINING int8) *SysStatus {
	m := SysStatus{}
	m.ONBOARD_CONTROL_SENSORS_PRESENT = ONBOARD_CONTROL_SENSORS_PRESENT
	m.ONBOARD_CONTROL_SENSORS_ENABLED = ONBOARD_CONTROL_SENSORS_ENABLED
	m.ONBOARD_CONTROL_SENSORS_HEALTH = ONBOARD_CONTROL_SENSORS_HEALTH
	m.LOAD = LOAD
	m.VOLTAGE_BATTERY = VOLTAGE_BATTERY
	m.CURRENT_BATTERY = CURRENT_BATTERY
	m.DROP_RATE_COMM = DROP_RATE_COMM
	m.ERRORS_COMM = ERRORS_COMM
	m.ERRORS_COUNT1 = ERRORS_COUNT1
	m.ERRORS_COUNT2 = ERRORS_COUNT2
	m.ERRORS_COUNT3 = ERRORS_COUNT3
	m.ERRORS_COUNT4 = ERRORS_COUNT4
	m.BATTERY_REMAINING = BATTERY_REMAINING
	return &m
}

// Id returns the SysStatus Message ID
func (*SysStatus) Id() uint8 {
	return 1
}

// Len returns the SysStatus Message Length
func (*SysStatus) Len() uint8 {
	return 31
}

// Crc returns the SysStatus Message CRC
func (*SysStatus) Crc() uint8 {
	return 124
}

// Pack returns a packed byte array which represents a SysStatus payload
func (m *SysStatus) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ONBOARD_CONTROL_SENSORS_PRESENT)
	binary.Write(data, binary.LittleEndian, m.ONBOARD_CONTROL_SENSORS_ENABLED)
	binary.Write(data, binary.LittleEndian, m.ONBOARD_CONTROL_SENSORS_HEALTH)
	binary.Write(data, binary.LittleEndian, m.LOAD)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_BATTERY)
	binary.Write(data, binary.LittleEndian, m.CURRENT_BATTERY)
	binary.Write(data, binary.LittleEndian, m.DROP_RATE_COMM)
	binary.Write(data, binary.LittleEndian, m.ERRORS_COMM)
	binary.Write(data, binary.LittleEndian, m.ERRORS_COUNT1)
	binary.Write(data, binary.LittleEndian, m.ERRORS_COUNT2)
	binary.Write(data, binary.LittleEndian, m.ERRORS_COUNT3)
	binary.Write(data, binary.LittleEndian, m.ERRORS_COUNT4)
	binary.Write(data, binary.LittleEndian, m.BATTERY_REMAINING)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SysStatus
func (m *SysStatus) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ONBOARD_CONTROL_SENSORS_PRESENT)
	binary.Read(data, binary.LittleEndian, &m.ONBOARD_CONTROL_SENSORS_ENABLED)
	binary.Read(data, binary.LittleEndian, &m.ONBOARD_CONTROL_SENSORS_HEALTH)
	binary.Read(data, binary.LittleEndian, &m.LOAD)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_BATTERY)
	binary.Read(data, binary.LittleEndian, &m.CURRENT_BATTERY)
	binary.Read(data, binary.LittleEndian, &m.DROP_RATE_COMM)
	binary.Read(data, binary.LittleEndian, &m.ERRORS_COMM)
	binary.Read(data, binary.LittleEndian, &m.ERRORS_COUNT1)
	binary.Read(data, binary.LittleEndian, &m.ERRORS_COUNT2)
	binary.Read(data, binary.LittleEndian, &m.ERRORS_COUNT3)
	binary.Read(data, binary.LittleEndian, &m.ERRORS_COUNT4)
	binary.Read(data, binary.LittleEndian, &m.BATTERY_REMAINING)
}

//
// MESSAGE SYSTEM_TIME
//
// MAVLINK_MSG_ID_SYSTEM_TIME 2
//
// MAVLINK_MSG_ID_SYSTEM_TIME_LEN 12
//
// MAVLINK_MSG_ID_SYSTEM_TIME_CRC 137
//
//
type SystemTime struct {
	TIME_UNIX_USEC uint64 // Timestamp of the master clock in microseconds since UNIX epoch.
	TIME_BOOT_MS   uint32 // Timestamp of the component clock since boot time in milliseconds.
}

// NewSystemTime returns a new SystemTime
func NewSystemTime(TIME_UNIX_USEC uint64, TIME_BOOT_MS uint32) *SystemTime {
	m := SystemTime{}
	m.TIME_UNIX_USEC = TIME_UNIX_USEC
	m.TIME_BOOT_MS = TIME_BOOT_MS
	return &m
}

// Id returns the SystemTime Message ID
func (*SystemTime) Id() uint8 {
	return 2
}

// Len returns the SystemTime Message Length
func (*SystemTime) Len() uint8 {
	return 12
}

// Crc returns the SystemTime Message CRC
func (*SystemTime) Crc() uint8 {
	return 137
}

// Pack returns a packed byte array which represents a SystemTime payload
func (m *SystemTime) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_UNIX_USEC)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SystemTime
func (m *SystemTime) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_UNIX_USEC)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
}

//
// MESSAGE PING
//
// MAVLINK_MSG_ID_PING 4
//
// MAVLINK_MSG_ID_PING_LEN 14
//
// MAVLINK_MSG_ID_PING_CRC 237
//
//
type Ping struct {
	TIME_USEC        uint64 // Unix timestamp in microseconds
	SEQ              uint32 // PING sequence
	TARGET_SYSTEM    uint8  // 0: request ping from all receiving systems, if greater than 0: message is a ping response and number is the system id of the requesting system
	TARGET_COMPONENT uint8  // 0: request ping from all receiving components, if greater than 0: message is a ping response and number is the system id of the requesting system
}

// NewPing returns a new Ping
func NewPing(TIME_USEC uint64, SEQ uint32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *Ping {
	m := Ping{}
	m.TIME_USEC = TIME_USEC
	m.SEQ = SEQ
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the Ping Message ID
func (*Ping) Id() uint8 {
	return 4
}

// Len returns the Ping Message Length
func (*Ping) Len() uint8 {
	return 14
}

// Crc returns the Ping Message CRC
func (*Ping) Crc() uint8 {
	return 237
}

// Pack returns a packed byte array which represents a Ping payload
func (m *Ping) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Ping
func (m *Ping) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE CHANGE_OPERATOR_CONTROL
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL 5
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL_LEN 28
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL_CRC 217
//
//
type ChangeOperatorControl struct {
	TARGET_SYSTEM   uint8     // System the GCS requests control for
	CONTROL_REQUEST uint8     // 0: request control of this MAV, 1: Release control of this MAV
	VERSION         uint8     // 0: key as plaintext, 1-255: future, different hashing/encryption variants. The GCS should in general use the safest mode possible initially and then gradually move down the encryption level if it gets a NACK message indicating an encryption mismatch.
	PASSKEY         [25]uint8 // Password / Key, depending on version plaintext or encrypted. 25 or less characters, NULL terminated. The characters may involve A-Z, a-z, 0-9, and "!?,.-"
}

// NewChangeOperatorControl returns a new ChangeOperatorControl
func NewChangeOperatorControl(TARGET_SYSTEM uint8, CONTROL_REQUEST uint8, VERSION uint8, PASSKEY [25]uint8) *ChangeOperatorControl {
	m := ChangeOperatorControl{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.CONTROL_REQUEST = CONTROL_REQUEST
	m.VERSION = VERSION
	m.PASSKEY = PASSKEY
	return &m
}

// Id returns the ChangeOperatorControl Message ID
func (*ChangeOperatorControl) Id() uint8 {
	return 5
}

// Len returns the ChangeOperatorControl Message Length
func (*ChangeOperatorControl) Len() uint8 {
	return 28
}

// Crc returns the ChangeOperatorControl Message CRC
func (*ChangeOperatorControl) Crc() uint8 {
	return 217
}

// Pack returns a packed byte array which represents a ChangeOperatorControl payload
func (m *ChangeOperatorControl) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.CONTROL_REQUEST)
	binary.Write(data, binary.LittleEndian, m.VERSION)
	binary.Write(data, binary.LittleEndian, m.PASSKEY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ChangeOperatorControl
func (m *ChangeOperatorControl) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.CONTROL_REQUEST)
	binary.Read(data, binary.LittleEndian, &m.VERSION)
	binary.Read(data, binary.LittleEndian, &m.PASSKEY)
}

const (
	MAVLINK_MSG_CHANGE_OPERATOR_CONTROL_FIELD_passkey_LEN = 25
)

//
// MESSAGE CHANGE_OPERATOR_CONTROL_ACK
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL_ACK 6
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL_ACK_LEN 3
//
// MAVLINK_MSG_ID_CHANGE_OPERATOR_CONTROL_ACK_CRC 104
//
//
type ChangeOperatorControlAck struct {
	GCS_SYSTEM_ID   uint8 // ID of the GCS this message
	CONTROL_REQUEST uint8 // 0: request control of this MAV, 1: Release control of this MAV
	ACK             uint8 // 0: ACK, 1: NACK: Wrong passkey, 2: NACK: Unsupported passkey encryption method, 3: NACK: Already under control
}

// NewChangeOperatorControlAck returns a new ChangeOperatorControlAck
func NewChangeOperatorControlAck(GCS_SYSTEM_ID uint8, CONTROL_REQUEST uint8, ACK uint8) *ChangeOperatorControlAck {
	m := ChangeOperatorControlAck{}
	m.GCS_SYSTEM_ID = GCS_SYSTEM_ID
	m.CONTROL_REQUEST = CONTROL_REQUEST
	m.ACK = ACK
	return &m
}

// Id returns the ChangeOperatorControlAck Message ID
func (*ChangeOperatorControlAck) Id() uint8 {
	return 6
}

// Len returns the ChangeOperatorControlAck Message Length
func (*ChangeOperatorControlAck) Len() uint8 {
	return 3
}

// Crc returns the ChangeOperatorControlAck Message CRC
func (*ChangeOperatorControlAck) Crc() uint8 {
	return 104
}

// Pack returns a packed byte array which represents a ChangeOperatorControlAck payload
func (m *ChangeOperatorControlAck) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.GCS_SYSTEM_ID)
	binary.Write(data, binary.LittleEndian, m.CONTROL_REQUEST)
	binary.Write(data, binary.LittleEndian, m.ACK)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ChangeOperatorControlAck
func (m *ChangeOperatorControlAck) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.GCS_SYSTEM_ID)
	binary.Read(data, binary.LittleEndian, &m.CONTROL_REQUEST)
	binary.Read(data, binary.LittleEndian, &m.ACK)
}

//
// MESSAGE AUTH_KEY
//
// MAVLINK_MSG_ID_AUTH_KEY 7
//
// MAVLINK_MSG_ID_AUTH_KEY_LEN 32
//
// MAVLINK_MSG_ID_AUTH_KEY_CRC 119
//
//
type AuthKey struct {
	KEY [32]uint8 // key
}

// NewAuthKey returns a new AuthKey
func NewAuthKey(KEY [32]uint8) *AuthKey {
	m := AuthKey{}
	m.KEY = KEY
	return &m
}

// Id returns the AuthKey Message ID
func (*AuthKey) Id() uint8 {
	return 7
}

// Len returns the AuthKey Message Length
func (*AuthKey) Len() uint8 {
	return 32
}

// Crc returns the AuthKey Message CRC
func (*AuthKey) Crc() uint8 {
	return 119
}

// Pack returns a packed byte array which represents a AuthKey payload
func (m *AuthKey) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.KEY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the AuthKey
func (m *AuthKey) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.KEY)
}

const (
	MAVLINK_MSG_AUTH_KEY_FIELD_key_LEN = 32
)

//
// MESSAGE SET_MODE
//
// MAVLINK_MSG_ID_SET_MODE 11
//
// MAVLINK_MSG_ID_SET_MODE_LEN 6
//
// MAVLINK_MSG_ID_SET_MODE_CRC 89
//
//
type SetMode struct {
	CUSTOM_MODE   uint32 // The new autopilot-specific mode. This field can be ignored by an autopilot.
	TARGET_SYSTEM uint8  // The system setting the mode
	BASE_MODE     uint8  // The new base mode
}

// NewSetMode returns a new SetMode
func NewSetMode(CUSTOM_MODE uint32, TARGET_SYSTEM uint8, BASE_MODE uint8) *SetMode {
	m := SetMode{}
	m.CUSTOM_MODE = CUSTOM_MODE
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.BASE_MODE = BASE_MODE
	return &m
}

// Id returns the SetMode Message ID
func (*SetMode) Id() uint8 {
	return 11
}

// Len returns the SetMode Message Length
func (*SetMode) Len() uint8 {
	return 6
}

// Crc returns the SetMode Message CRC
func (*SetMode) Crc() uint8 {
	return 89
}

// Pack returns a packed byte array which represents a SetMode payload
func (m *SetMode) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.CUSTOM_MODE)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.BASE_MODE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetMode
func (m *SetMode) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.CUSTOM_MODE)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.BASE_MODE)
}

//
// MESSAGE PARAM_REQUEST_READ
//
// MAVLINK_MSG_ID_PARAM_REQUEST_READ 20
//
// MAVLINK_MSG_ID_PARAM_REQUEST_READ_LEN 20
//
// MAVLINK_MSG_ID_PARAM_REQUEST_READ_CRC 214
//
//
type ParamRequestRead struct {
	PARAM_INDEX      int16     // Parameter index. Send -1 to use the param ID field as identifier (else the param id will be ignored)
	TARGET_SYSTEM    uint8     // System ID
	TARGET_COMPONENT uint8     // Component ID
	PARAM_ID         [16]uint8 // Onboard parameter id, terminated by NULL if the length is less than 16 human-readable chars and WITHOUT null termination (NULL) byte if the length is exactly 16 chars - applications have to provide 16+1 bytes storage if the ID is stored as string
}

// NewParamRequestRead returns a new ParamRequestRead
func NewParamRequestRead(PARAM_INDEX int16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, PARAM_ID [16]uint8) *ParamRequestRead {
	m := ParamRequestRead{}
	m.PARAM_INDEX = PARAM_INDEX
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.PARAM_ID = PARAM_ID
	return &m
}

// Id returns the ParamRequestRead Message ID
func (*ParamRequestRead) Id() uint8 {
	return 20
}

// Len returns the ParamRequestRead Message Length
func (*ParamRequestRead) Len() uint8 {
	return 20
}

// Crc returns the ParamRequestRead Message CRC
func (*ParamRequestRead) Crc() uint8 {
	return 214
}

// Pack returns a packed byte array which represents a ParamRequestRead payload
func (m *ParamRequestRead) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.PARAM_INDEX)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.PARAM_ID)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ParamRequestRead
func (m *ParamRequestRead) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.PARAM_INDEX)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.PARAM_ID)
}

const (
	MAVLINK_MSG_PARAM_REQUEST_READ_FIELD_param_id_LEN = 16
)

//
// MESSAGE PARAM_REQUEST_LIST
//
// MAVLINK_MSG_ID_PARAM_REQUEST_LIST 21
//
// MAVLINK_MSG_ID_PARAM_REQUEST_LIST_LEN 2
//
// MAVLINK_MSG_ID_PARAM_REQUEST_LIST_CRC 159
//
//
type ParamRequestList struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewParamRequestList returns a new ParamRequestList
func NewParamRequestList(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *ParamRequestList {
	m := ParamRequestList{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the ParamRequestList Message ID
func (*ParamRequestList) Id() uint8 {
	return 21
}

// Len returns the ParamRequestList Message Length
func (*ParamRequestList) Len() uint8 {
	return 2
}

// Crc returns the ParamRequestList Message CRC
func (*ParamRequestList) Crc() uint8 {
	return 159
}

// Pack returns a packed byte array which represents a ParamRequestList payload
func (m *ParamRequestList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ParamRequestList
func (m *ParamRequestList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE PARAM_VALUE
//
// MAVLINK_MSG_ID_PARAM_VALUE 22
//
// MAVLINK_MSG_ID_PARAM_VALUE_LEN 25
//
// MAVLINK_MSG_ID_PARAM_VALUE_CRC 220
//
//
type ParamValue struct {
	PARAM_VALUE float32   // Onboard parameter value
	PARAM_COUNT uint16    // Total number of onboard parameters
	PARAM_INDEX uint16    // Index of this onboard parameter
	PARAM_ID    [16]uint8 // Onboard parameter id, terminated by NULL if the length is less than 16 human-readable chars and WITHOUT null termination (NULL) byte if the length is exactly 16 chars - applications have to provide 16+1 bytes storage if the ID is stored as string
	PARAM_TYPE  uint8     // Onboard parameter type: see the MAV_PARAM_TYPE enum for supported data types.
}

// NewParamValue returns a new ParamValue
func NewParamValue(PARAM_VALUE float32, PARAM_COUNT uint16, PARAM_INDEX uint16, PARAM_ID [16]uint8, PARAM_TYPE uint8) *ParamValue {
	m := ParamValue{}
	m.PARAM_VALUE = PARAM_VALUE
	m.PARAM_COUNT = PARAM_COUNT
	m.PARAM_INDEX = PARAM_INDEX
	m.PARAM_ID = PARAM_ID
	m.PARAM_TYPE = PARAM_TYPE
	return &m
}

// Id returns the ParamValue Message ID
func (*ParamValue) Id() uint8 {
	return 22
}

// Len returns the ParamValue Message Length
func (*ParamValue) Len() uint8 {
	return 25
}

// Crc returns the ParamValue Message CRC
func (*ParamValue) Crc() uint8 {
	return 220
}

// Pack returns a packed byte array which represents a ParamValue payload
func (m *ParamValue) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.PARAM_VALUE)
	binary.Write(data, binary.LittleEndian, m.PARAM_COUNT)
	binary.Write(data, binary.LittleEndian, m.PARAM_INDEX)
	binary.Write(data, binary.LittleEndian, m.PARAM_ID)
	binary.Write(data, binary.LittleEndian, m.PARAM_TYPE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ParamValue
func (m *ParamValue) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.PARAM_VALUE)
	binary.Read(data, binary.LittleEndian, &m.PARAM_COUNT)
	binary.Read(data, binary.LittleEndian, &m.PARAM_INDEX)
	binary.Read(data, binary.LittleEndian, &m.PARAM_ID)
	binary.Read(data, binary.LittleEndian, &m.PARAM_TYPE)
}

const (
	MAVLINK_MSG_PARAM_VALUE_FIELD_param_id_LEN = 16
)

//
// MESSAGE PARAM_SET
//
// MAVLINK_MSG_ID_PARAM_SET 23
//
// MAVLINK_MSG_ID_PARAM_SET_LEN 23
//
// MAVLINK_MSG_ID_PARAM_SET_CRC 168
//
//
type ParamSet struct {
	PARAM_VALUE      float32   // Onboard parameter value
	TARGET_SYSTEM    uint8     // System ID
	TARGET_COMPONENT uint8     // Component ID
	PARAM_ID         [16]uint8 // Onboard parameter id, terminated by NULL if the length is less than 16 human-readable chars and WITHOUT null termination (NULL) byte if the length is exactly 16 chars - applications have to provide 16+1 bytes storage if the ID is stored as string
	PARAM_TYPE       uint8     // Onboard parameter type: see the MAV_PARAM_TYPE enum for supported data types.
}

// NewParamSet returns a new ParamSet
func NewParamSet(PARAM_VALUE float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, PARAM_ID [16]uint8, PARAM_TYPE uint8) *ParamSet {
	m := ParamSet{}
	m.PARAM_VALUE = PARAM_VALUE
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.PARAM_ID = PARAM_ID
	m.PARAM_TYPE = PARAM_TYPE
	return &m
}

// Id returns the ParamSet Message ID
func (*ParamSet) Id() uint8 {
	return 23
}

// Len returns the ParamSet Message Length
func (*ParamSet) Len() uint8 {
	return 23
}

// Crc returns the ParamSet Message CRC
func (*ParamSet) Crc() uint8 {
	return 168
}

// Pack returns a packed byte array which represents a ParamSet payload
func (m *ParamSet) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.PARAM_VALUE)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.PARAM_ID)
	binary.Write(data, binary.LittleEndian, m.PARAM_TYPE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ParamSet
func (m *ParamSet) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.PARAM_VALUE)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.PARAM_ID)
	binary.Read(data, binary.LittleEndian, &m.PARAM_TYPE)
}

const (
	MAVLINK_MSG_PARAM_SET_FIELD_param_id_LEN = 16
)

//
// MESSAGE GPS_RAW_INT
//
// MAVLINK_MSG_ID_GPS_RAW_INT 24
//
// MAVLINK_MSG_ID_GPS_RAW_INT_LEN 30
//
// MAVLINK_MSG_ID_GPS_RAW_INT_CRC 24
//
//
type GpsRawInt struct {
	TIME_USEC          uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	LAT                int32  // Latitude (WGS84), in degrees * 1E7
	LON                int32  // Longitude (WGS84), in degrees * 1E7
	ALT                int32  // Altitude (WGS84), in meters * 1000 (positive for up)
	EPH                uint16 // GPS HDOP horizontal dilution of position in cm (m*100). If unknown, set to: UINT16_MAX
	EPV                uint16 // GPS VDOP vertical dilution of position in cm (m*100). If unknown, set to: UINT16_MAX
	VEL                uint16 // GPS ground speed (m/s * 100). If unknown, set to: UINT16_MAX
	COG                uint16 // Course over ground (NOT heading, but direction of movement) in degrees * 100, 0.0..359.99 degrees. If unknown, set to: UINT16_MAX
	FIX_TYPE           uint8  // 0-1: no fix, 2: 2D fix, 3: 3D fix, 4: DGPS, 5: RTK. Some applications will not use the value of this field unless it is at least two, so always correctly fill in the fix.
	SATELLITES_VISIBLE uint8  // Number of satellites visible. If unknown, set to 255
}

// NewGpsRawInt returns a new GpsRawInt
func NewGpsRawInt(TIME_USEC uint64, LAT int32, LON int32, ALT int32, EPH uint16, EPV uint16, VEL uint16, COG uint16, FIX_TYPE uint8, SATELLITES_VISIBLE uint8) *GpsRawInt {
	m := GpsRawInt{}
	m.TIME_USEC = TIME_USEC
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.EPH = EPH
	m.EPV = EPV
	m.VEL = VEL
	m.COG = COG
	m.FIX_TYPE = FIX_TYPE
	m.SATELLITES_VISIBLE = SATELLITES_VISIBLE
	return &m
}

// Id returns the GpsRawInt Message ID
func (*GpsRawInt) Id() uint8 {
	return 24
}

// Len returns the GpsRawInt Message Length
func (*GpsRawInt) Len() uint8 {
	return 30
}

// Crc returns the GpsRawInt Message CRC
func (*GpsRawInt) Crc() uint8 {
	return 24
}

// Pack returns a packed byte array which represents a GpsRawInt payload
func (m *GpsRawInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.EPH)
	binary.Write(data, binary.LittleEndian, m.EPV)
	binary.Write(data, binary.LittleEndian, m.VEL)
	binary.Write(data, binary.LittleEndian, m.COG)
	binary.Write(data, binary.LittleEndian, m.FIX_TYPE)
	binary.Write(data, binary.LittleEndian, m.SATELLITES_VISIBLE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GpsRawInt
func (m *GpsRawInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.EPH)
	binary.Read(data, binary.LittleEndian, &m.EPV)
	binary.Read(data, binary.LittleEndian, &m.VEL)
	binary.Read(data, binary.LittleEndian, &m.COG)
	binary.Read(data, binary.LittleEndian, &m.FIX_TYPE)
	binary.Read(data, binary.LittleEndian, &m.SATELLITES_VISIBLE)
}

//
// MESSAGE GPS_STATUS
//
// MAVLINK_MSG_ID_GPS_STATUS 25
//
// MAVLINK_MSG_ID_GPS_STATUS_LEN 101
//
// MAVLINK_MSG_ID_GPS_STATUS_CRC 23
//
//
type GpsStatus struct {
	SATELLITES_VISIBLE  uint8     // Number of satellites visible
	SATELLITE_PRN       [20]uint8 // Global satellite ID
	SATELLITE_USED      [20]uint8 // 0: Satellite not used, 1: used for localization
	SATELLITE_ELEVATION [20]uint8 // Elevation (0: right on top of receiver, 90: on the horizon) of satellite
	SATELLITE_AZIMUTH   [20]uint8 // Direction of satellite, 0: 0 deg, 255: 360 deg.
	SATELLITE_SNR       [20]uint8 // Signal to noise ratio of satellite
}

// NewGpsStatus returns a new GpsStatus
func NewGpsStatus(SATELLITES_VISIBLE uint8, SATELLITE_PRN [20]uint8, SATELLITE_USED [20]uint8, SATELLITE_ELEVATION [20]uint8, SATELLITE_AZIMUTH [20]uint8, SATELLITE_SNR [20]uint8) *GpsStatus {
	m := GpsStatus{}
	m.SATELLITES_VISIBLE = SATELLITES_VISIBLE
	m.SATELLITE_PRN = SATELLITE_PRN
	m.SATELLITE_USED = SATELLITE_USED
	m.SATELLITE_ELEVATION = SATELLITE_ELEVATION
	m.SATELLITE_AZIMUTH = SATELLITE_AZIMUTH
	m.SATELLITE_SNR = SATELLITE_SNR
	return &m
}

// Id returns the GpsStatus Message ID
func (*GpsStatus) Id() uint8 {
	return 25
}

// Len returns the GpsStatus Message Length
func (*GpsStatus) Len() uint8 {
	return 101
}

// Crc returns the GpsStatus Message CRC
func (*GpsStatus) Crc() uint8 {
	return 23
}

// Pack returns a packed byte array which represents a GpsStatus payload
func (m *GpsStatus) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SATELLITES_VISIBLE)
	binary.Write(data, binary.LittleEndian, m.SATELLITE_PRN)
	binary.Write(data, binary.LittleEndian, m.SATELLITE_USED)
	binary.Write(data, binary.LittleEndian, m.SATELLITE_ELEVATION)
	binary.Write(data, binary.LittleEndian, m.SATELLITE_AZIMUTH)
	binary.Write(data, binary.LittleEndian, m.SATELLITE_SNR)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GpsStatus
func (m *GpsStatus) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SATELLITES_VISIBLE)
	binary.Read(data, binary.LittleEndian, &m.SATELLITE_PRN)
	binary.Read(data, binary.LittleEndian, &m.SATELLITE_USED)
	binary.Read(data, binary.LittleEndian, &m.SATELLITE_ELEVATION)
	binary.Read(data, binary.LittleEndian, &m.SATELLITE_AZIMUTH)
	binary.Read(data, binary.LittleEndian, &m.SATELLITE_SNR)
}

const (
	MAVLINK_MSG_GPS_STATUS_FIELD_satellite_prn_LEN       = 20
	MAVLINK_MSG_GPS_STATUS_FIELD_satellite_used_LEN      = 20
	MAVLINK_MSG_GPS_STATUS_FIELD_satellite_elevation_LEN = 20
	MAVLINK_MSG_GPS_STATUS_FIELD_satellite_azimuth_LEN   = 20
	MAVLINK_MSG_GPS_STATUS_FIELD_satellite_snr_LEN       = 20
)

//
// MESSAGE SCALED_IMU
//
// MAVLINK_MSG_ID_SCALED_IMU 26
//
// MAVLINK_MSG_ID_SCALED_IMU_LEN 22
//
// MAVLINK_MSG_ID_SCALED_IMU_CRC 170
//
//
type ScaledImu struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	XACC         int16  // X acceleration (mg)
	YACC         int16  // Y acceleration (mg)
	ZACC         int16  // Z acceleration (mg)
	XGYRO        int16  // Angular speed around X axis (millirad /sec)
	YGYRO        int16  // Angular speed around Y axis (millirad /sec)
	ZGYRO        int16  // Angular speed around Z axis (millirad /sec)
	XMAG         int16  // X Magnetic field (milli tesla)
	YMAG         int16  // Y Magnetic field (milli tesla)
	ZMAG         int16  // Z Magnetic field (milli tesla)
}

// NewScaledImu returns a new ScaledImu
func NewScaledImu(TIME_BOOT_MS uint32, XACC int16, YACC int16, ZACC int16, XGYRO int16, YGYRO int16, ZGYRO int16, XMAG int16, YMAG int16, ZMAG int16) *ScaledImu {
	m := ScaledImu{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.XMAG = XMAG
	m.YMAG = YMAG
	m.ZMAG = ZMAG
	return &m
}

// Id returns the ScaledImu Message ID
func (*ScaledImu) Id() uint8 {
	return 26
}

// Len returns the ScaledImu Message Length
func (*ScaledImu) Len() uint8 {
	return 22
}

// Crc returns the ScaledImu Message CRC
func (*ScaledImu) Crc() uint8 {
	return 170
}

// Pack returns a packed byte array which represents a ScaledImu payload
func (m *ScaledImu) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.XMAG)
	binary.Write(data, binary.LittleEndian, m.YMAG)
	binary.Write(data, binary.LittleEndian, m.ZMAG)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ScaledImu
func (m *ScaledImu) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.XMAG)
	binary.Read(data, binary.LittleEndian, &m.YMAG)
	binary.Read(data, binary.LittleEndian, &m.ZMAG)
}

//
// MESSAGE RAW_IMU
//
// MAVLINK_MSG_ID_RAW_IMU 27
//
// MAVLINK_MSG_ID_RAW_IMU_LEN 26
//
// MAVLINK_MSG_ID_RAW_IMU_CRC 144
//
//
type RawImu struct {
	TIME_USEC uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	XACC      int16  // X acceleration (raw)
	YACC      int16  // Y acceleration (raw)
	ZACC      int16  // Z acceleration (raw)
	XGYRO     int16  // Angular speed around X axis (raw)
	YGYRO     int16  // Angular speed around Y axis (raw)
	ZGYRO     int16  // Angular speed around Z axis (raw)
	XMAG      int16  // X Magnetic field (raw)
	YMAG      int16  // Y Magnetic field (raw)
	ZMAG      int16  // Z Magnetic field (raw)
}

// NewRawImu returns a new RawImu
func NewRawImu(TIME_USEC uint64, XACC int16, YACC int16, ZACC int16, XGYRO int16, YGYRO int16, ZGYRO int16, XMAG int16, YMAG int16, ZMAG int16) *RawImu {
	m := RawImu{}
	m.TIME_USEC = TIME_USEC
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.XMAG = XMAG
	m.YMAG = YMAG
	m.ZMAG = ZMAG
	return &m
}

// Id returns the RawImu Message ID
func (*RawImu) Id() uint8 {
	return 27
}

// Len returns the RawImu Message Length
func (*RawImu) Len() uint8 {
	return 26
}

// Crc returns the RawImu Message CRC
func (*RawImu) Crc() uint8 {
	return 144
}

// Pack returns a packed byte array which represents a RawImu payload
func (m *RawImu) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.XMAG)
	binary.Write(data, binary.LittleEndian, m.YMAG)
	binary.Write(data, binary.LittleEndian, m.ZMAG)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RawImu
func (m *RawImu) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.XMAG)
	binary.Read(data, binary.LittleEndian, &m.YMAG)
	binary.Read(data, binary.LittleEndian, &m.ZMAG)
}

//
// MESSAGE RAW_PRESSURE
//
// MAVLINK_MSG_ID_RAW_PRESSURE 28
//
// MAVLINK_MSG_ID_RAW_PRESSURE_LEN 16
//
// MAVLINK_MSG_ID_RAW_PRESSURE_CRC 67
//
//
type RawPressure struct {
	TIME_USEC   uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	PRESS_ABS   int16  // Absolute pressure (raw)
	PRESS_DIFF1 int16  // Differential pressure 1 (raw)
	PRESS_DIFF2 int16  // Differential pressure 2 (raw)
	TEMPERATURE int16  // Raw Temperature measurement (raw)
}

// NewRawPressure returns a new RawPressure
func NewRawPressure(TIME_USEC uint64, PRESS_ABS int16, PRESS_DIFF1 int16, PRESS_DIFF2 int16, TEMPERATURE int16) *RawPressure {
	m := RawPressure{}
	m.TIME_USEC = TIME_USEC
	m.PRESS_ABS = PRESS_ABS
	m.PRESS_DIFF1 = PRESS_DIFF1
	m.PRESS_DIFF2 = PRESS_DIFF2
	m.TEMPERATURE = TEMPERATURE
	return &m
}

// Id returns the RawPressure Message ID
func (*RawPressure) Id() uint8 {
	return 28
}

// Len returns the RawPressure Message Length
func (*RawPressure) Len() uint8 {
	return 16
}

// Crc returns the RawPressure Message CRC
func (*RawPressure) Crc() uint8 {
	return 67
}

// Pack returns a packed byte array which represents a RawPressure payload
func (m *RawPressure) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.PRESS_ABS)
	binary.Write(data, binary.LittleEndian, m.PRESS_DIFF1)
	binary.Write(data, binary.LittleEndian, m.PRESS_DIFF2)
	binary.Write(data, binary.LittleEndian, m.TEMPERATURE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RawPressure
func (m *RawPressure) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.PRESS_ABS)
	binary.Read(data, binary.LittleEndian, &m.PRESS_DIFF1)
	binary.Read(data, binary.LittleEndian, &m.PRESS_DIFF2)
	binary.Read(data, binary.LittleEndian, &m.TEMPERATURE)
}

//
// MESSAGE SCALED_PRESSURE
//
// MAVLINK_MSG_ID_SCALED_PRESSURE 29
//
// MAVLINK_MSG_ID_SCALED_PRESSURE_LEN 14
//
// MAVLINK_MSG_ID_SCALED_PRESSURE_CRC 115
//
//
type ScaledPressure struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	PRESS_ABS    float32 // Absolute pressure (hectopascal)
	PRESS_DIFF   float32 // Differential pressure 1 (hectopascal)
	TEMPERATURE  int16   // Temperature measurement (0.01 degrees celsius)
}

// NewScaledPressure returns a new ScaledPressure
func NewScaledPressure(TIME_BOOT_MS uint32, PRESS_ABS float32, PRESS_DIFF float32, TEMPERATURE int16) *ScaledPressure {
	m := ScaledPressure{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.PRESS_ABS = PRESS_ABS
	m.PRESS_DIFF = PRESS_DIFF
	m.TEMPERATURE = TEMPERATURE
	return &m
}

// Id returns the ScaledPressure Message ID
func (*ScaledPressure) Id() uint8 {
	return 29
}

// Len returns the ScaledPressure Message Length
func (*ScaledPressure) Len() uint8 {
	return 14
}

// Crc returns the ScaledPressure Message CRC
func (*ScaledPressure) Crc() uint8 {
	return 115
}

// Pack returns a packed byte array which represents a ScaledPressure payload
func (m *ScaledPressure) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.PRESS_ABS)
	binary.Write(data, binary.LittleEndian, m.PRESS_DIFF)
	binary.Write(data, binary.LittleEndian, m.TEMPERATURE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ScaledPressure
func (m *ScaledPressure) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.PRESS_ABS)
	binary.Read(data, binary.LittleEndian, &m.PRESS_DIFF)
	binary.Read(data, binary.LittleEndian, &m.TEMPERATURE)
}

//
// MESSAGE ATTITUDE
//
// MAVLINK_MSG_ID_ATTITUDE 30
//
// MAVLINK_MSG_ID_ATTITUDE_LEN 28
//
// MAVLINK_MSG_ID_ATTITUDE_CRC 39
//
//
type Attitude struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	ROLL         float32 // Roll angle (rad, -pi..+pi)
	PITCH        float32 // Pitch angle (rad, -pi..+pi)
	YAW          float32 // Yaw angle (rad, -pi..+pi)
	ROLLSPEED    float32 // Roll angular speed (rad/s)
	PITCHSPEED   float32 // Pitch angular speed (rad/s)
	YAWSPEED     float32 // Yaw angular speed (rad/s)
}

// NewAttitude returns a new Attitude
func NewAttitude(TIME_BOOT_MS uint32, ROLL float32, PITCH float32, YAW float32, ROLLSPEED float32, PITCHSPEED float32, YAWSPEED float32) *Attitude {
	m := Attitude{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.ROLLSPEED = ROLLSPEED
	m.PITCHSPEED = PITCHSPEED
	m.YAWSPEED = YAWSPEED
	return &m
}

// Id returns the Attitude Message ID
func (*Attitude) Id() uint8 {
	return 30
}

// Len returns the Attitude Message Length
func (*Attitude) Len() uint8 {
	return 28
}

// Crc returns the Attitude Message CRC
func (*Attitude) Crc() uint8 {
	return 39
}

// Pack returns a packed byte array which represents a Attitude payload
func (m *Attitude) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.ROLLSPEED)
	binary.Write(data, binary.LittleEndian, m.PITCHSPEED)
	binary.Write(data, binary.LittleEndian, m.YAWSPEED)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Attitude
func (m *Attitude) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.ROLLSPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCHSPEED)
	binary.Read(data, binary.LittleEndian, &m.YAWSPEED)
}

//
// MESSAGE ATTITUDE_QUATERNION
//
// MAVLINK_MSG_ID_ATTITUDE_QUATERNION 31
//
// MAVLINK_MSG_ID_ATTITUDE_QUATERNION_LEN 32
//
// MAVLINK_MSG_ID_ATTITUDE_QUATERNION_CRC 246
//
//
type AttitudeQuaternion struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	Q1           float32 // Quaternion component 1, w (1 in null-rotation)
	Q2           float32 // Quaternion component 2, x (0 in null-rotation)
	Q3           float32 // Quaternion component 3, y (0 in null-rotation)
	Q4           float32 // Quaternion component 4, z (0 in null-rotation)
	ROLLSPEED    float32 // Roll angular speed (rad/s)
	PITCHSPEED   float32 // Pitch angular speed (rad/s)
	YAWSPEED     float32 // Yaw angular speed (rad/s)
}

// NewAttitudeQuaternion returns a new AttitudeQuaternion
func NewAttitudeQuaternion(TIME_BOOT_MS uint32, Q1 float32, Q2 float32, Q3 float32, Q4 float32, ROLLSPEED float32, PITCHSPEED float32, YAWSPEED float32) *AttitudeQuaternion {
	m := AttitudeQuaternion{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.Q1 = Q1
	m.Q2 = Q2
	m.Q3 = Q3
	m.Q4 = Q4
	m.ROLLSPEED = ROLLSPEED
	m.PITCHSPEED = PITCHSPEED
	m.YAWSPEED = YAWSPEED
	return &m
}

// Id returns the AttitudeQuaternion Message ID
func (*AttitudeQuaternion) Id() uint8 {
	return 31
}

// Len returns the AttitudeQuaternion Message Length
func (*AttitudeQuaternion) Len() uint8 {
	return 32
}

// Crc returns the AttitudeQuaternion Message CRC
func (*AttitudeQuaternion) Crc() uint8 {
	return 246
}

// Pack returns a packed byte array which represents a AttitudeQuaternion payload
func (m *AttitudeQuaternion) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.Q1)
	binary.Write(data, binary.LittleEndian, m.Q2)
	binary.Write(data, binary.LittleEndian, m.Q3)
	binary.Write(data, binary.LittleEndian, m.Q4)
	binary.Write(data, binary.LittleEndian, m.ROLLSPEED)
	binary.Write(data, binary.LittleEndian, m.PITCHSPEED)
	binary.Write(data, binary.LittleEndian, m.YAWSPEED)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the AttitudeQuaternion
func (m *AttitudeQuaternion) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.Q1)
	binary.Read(data, binary.LittleEndian, &m.Q2)
	binary.Read(data, binary.LittleEndian, &m.Q3)
	binary.Read(data, binary.LittleEndian, &m.Q4)
	binary.Read(data, binary.LittleEndian, &m.ROLLSPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCHSPEED)
	binary.Read(data, binary.LittleEndian, &m.YAWSPEED)
}

//
// MESSAGE LOCAL_POSITION_NED
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED 32
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED_LEN 28
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED_CRC 185
//
//
type LocalPositionNed struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	X            float32 // X Position
	Y            float32 // Y Position
	Z            float32 // Z Position
	VX           float32 // X Speed
	VY           float32 // Y Speed
	VZ           float32 // Z Speed
}

// NewLocalPositionNed returns a new LocalPositionNed
func NewLocalPositionNed(TIME_BOOT_MS uint32, X float32, Y float32, Z float32, VX float32, VY float32, VZ float32) *LocalPositionNed {
	m := LocalPositionNed{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.X = X
	m.Y = Y
	m.Z = Z
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	return &m
}

// Id returns the LocalPositionNed Message ID
func (*LocalPositionNed) Id() uint8 {
	return 32
}

// Len returns the LocalPositionNed Message Length
func (*LocalPositionNed) Len() uint8 {
	return 28
}

// Crc returns the LocalPositionNed Message CRC
func (*LocalPositionNed) Crc() uint8 {
	return 185
}

// Pack returns a packed byte array which represents a LocalPositionNed payload
func (m *LocalPositionNed) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LocalPositionNed
func (m *LocalPositionNed) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
}

//
// MESSAGE GLOBAL_POSITION_INT
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_INT 33
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_INT_LEN 28
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_INT_CRC 104
//
//
type GlobalPositionInt struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	LAT          int32  // Latitude, expressed as * 1E7
	LON          int32  // Longitude, expressed as * 1E7
	ALT          int32  // Altitude in meters, expressed as * 1000 (millimeters), above MSL
	RELATIVE_ALT int32  // Altitude above ground in meters, expressed as * 1000 (millimeters)
	VX           int16  // Ground X Speed (Latitude), expressed as m/s * 100
	VY           int16  // Ground Y Speed (Longitude), expressed as m/s * 100
	VZ           int16  // Ground Z Speed (Altitude), expressed as m/s * 100
	HDG          uint16 // Compass heading in degrees * 100, 0.0..359.99 degrees. If unknown, set to: UINT16_MAX
}

// NewGlobalPositionInt returns a new GlobalPositionInt
func NewGlobalPositionInt(TIME_BOOT_MS uint32, LAT int32, LON int32, ALT int32, RELATIVE_ALT int32, VX int16, VY int16, VZ int16, HDG uint16) *GlobalPositionInt {
	m := GlobalPositionInt{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.RELATIVE_ALT = RELATIVE_ALT
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	m.HDG = HDG
	return &m
}

// Id returns the GlobalPositionInt Message ID
func (*GlobalPositionInt) Id() uint8 {
	return 33
}

// Len returns the GlobalPositionInt Message Length
func (*GlobalPositionInt) Len() uint8 {
	return 28
}

// Crc returns the GlobalPositionInt Message CRC
func (*GlobalPositionInt) Crc() uint8 {
	return 104
}

// Pack returns a packed byte array which represents a GlobalPositionInt payload
func (m *GlobalPositionInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.RELATIVE_ALT)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	binary.Write(data, binary.LittleEndian, m.HDG)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GlobalPositionInt
func (m *GlobalPositionInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.RELATIVE_ALT)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
	binary.Read(data, binary.LittleEndian, &m.HDG)
}

//
// MESSAGE RC_CHANNELS_SCALED
//
// MAVLINK_MSG_ID_RC_CHANNELS_SCALED 34
//
// MAVLINK_MSG_ID_RC_CHANNELS_SCALED_LEN 22
//
// MAVLINK_MSG_ID_RC_CHANNELS_SCALED_CRC 237
//
//
type RcChannelsScaled struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	CHAN1_SCALED int16  // RC channel 1 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN2_SCALED int16  // RC channel 2 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN3_SCALED int16  // RC channel 3 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN4_SCALED int16  // RC channel 4 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN5_SCALED int16  // RC channel 5 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN6_SCALED int16  // RC channel 6 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN7_SCALED int16  // RC channel 7 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	CHAN8_SCALED int16  // RC channel 8 value scaled, (-100%) -10000, (0%) 0, (100%) 10000, (invalid) INT16_MAX.
	PORT         uint8  // Servo output port (set of 8 outputs = 1 port). Most MAVs will just use one, but this allows for more than 8 servos.
	RSSI         uint8  // Receive signal strength indicator, 0: 0%, 100: 100%, 255: invalid/unknown.
}

// NewRcChannelsScaled returns a new RcChannelsScaled
func NewRcChannelsScaled(TIME_BOOT_MS uint32, CHAN1_SCALED int16, CHAN2_SCALED int16, CHAN3_SCALED int16, CHAN4_SCALED int16, CHAN5_SCALED int16, CHAN6_SCALED int16, CHAN7_SCALED int16, CHAN8_SCALED int16, PORT uint8, RSSI uint8) *RcChannelsScaled {
	m := RcChannelsScaled{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.CHAN1_SCALED = CHAN1_SCALED
	m.CHAN2_SCALED = CHAN2_SCALED
	m.CHAN3_SCALED = CHAN3_SCALED
	m.CHAN4_SCALED = CHAN4_SCALED
	m.CHAN5_SCALED = CHAN5_SCALED
	m.CHAN6_SCALED = CHAN6_SCALED
	m.CHAN7_SCALED = CHAN7_SCALED
	m.CHAN8_SCALED = CHAN8_SCALED
	m.PORT = PORT
	m.RSSI = RSSI
	return &m
}

// Id returns the RcChannelsScaled Message ID
func (*RcChannelsScaled) Id() uint8 {
	return 34
}

// Len returns the RcChannelsScaled Message Length
func (*RcChannelsScaled) Len() uint8 {
	return 22
}

// Crc returns the RcChannelsScaled Message CRC
func (*RcChannelsScaled) Crc() uint8 {
	return 237
}

// Pack returns a packed byte array which represents a RcChannelsScaled payload
func (m *RcChannelsScaled) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.CHAN1_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN2_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN3_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN4_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN5_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN6_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN7_SCALED)
	binary.Write(data, binary.LittleEndian, m.CHAN8_SCALED)
	binary.Write(data, binary.LittleEndian, m.PORT)
	binary.Write(data, binary.LittleEndian, m.RSSI)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RcChannelsScaled
func (m *RcChannelsScaled) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.CHAN1_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN2_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN3_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN4_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN5_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN6_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN7_SCALED)
	binary.Read(data, binary.LittleEndian, &m.CHAN8_SCALED)
	binary.Read(data, binary.LittleEndian, &m.PORT)
	binary.Read(data, binary.LittleEndian, &m.RSSI)
}

//
// MESSAGE RC_CHANNELS_RAW
//
// MAVLINK_MSG_ID_RC_CHANNELS_RAW 35
//
// MAVLINK_MSG_ID_RC_CHANNELS_RAW_LEN 22
//
// MAVLINK_MSG_ID_RC_CHANNELS_RAW_CRC 244
//
//
type RcChannelsRaw struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	CHAN1_RAW    uint16 // RC channel 1 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN2_RAW    uint16 // RC channel 2 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN3_RAW    uint16 // RC channel 3 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN4_RAW    uint16 // RC channel 4 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN5_RAW    uint16 // RC channel 5 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN6_RAW    uint16 // RC channel 6 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN7_RAW    uint16 // RC channel 7 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN8_RAW    uint16 // RC channel 8 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	PORT         uint8  // Servo output port (set of 8 outputs = 1 port). Most MAVs will just use one, but this allows for more than 8 servos.
	RSSI         uint8  // Receive signal strength indicator, 0: 0%, 100: 100%, 255: invalid/unknown.
}

// NewRcChannelsRaw returns a new RcChannelsRaw
func NewRcChannelsRaw(TIME_BOOT_MS uint32, CHAN1_RAW uint16, CHAN2_RAW uint16, CHAN3_RAW uint16, CHAN4_RAW uint16, CHAN5_RAW uint16, CHAN6_RAW uint16, CHAN7_RAW uint16, CHAN8_RAW uint16, PORT uint8, RSSI uint8) *RcChannelsRaw {
	m := RcChannelsRaw{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.CHAN1_RAW = CHAN1_RAW
	m.CHAN2_RAW = CHAN2_RAW
	m.CHAN3_RAW = CHAN3_RAW
	m.CHAN4_RAW = CHAN4_RAW
	m.CHAN5_RAW = CHAN5_RAW
	m.CHAN6_RAW = CHAN6_RAW
	m.CHAN7_RAW = CHAN7_RAW
	m.CHAN8_RAW = CHAN8_RAW
	m.PORT = PORT
	m.RSSI = RSSI
	return &m
}

// Id returns the RcChannelsRaw Message ID
func (*RcChannelsRaw) Id() uint8 {
	return 35
}

// Len returns the RcChannelsRaw Message Length
func (*RcChannelsRaw) Len() uint8 {
	return 22
}

// Crc returns the RcChannelsRaw Message CRC
func (*RcChannelsRaw) Crc() uint8 {
	return 244
}

// Pack returns a packed byte array which represents a RcChannelsRaw payload
func (m *RcChannelsRaw) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.CHAN1_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN2_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN3_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN4_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN5_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN6_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN7_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN8_RAW)
	binary.Write(data, binary.LittleEndian, m.PORT)
	binary.Write(data, binary.LittleEndian, m.RSSI)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RcChannelsRaw
func (m *RcChannelsRaw) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.CHAN1_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN2_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN3_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN4_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN5_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN6_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN7_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN8_RAW)
	binary.Read(data, binary.LittleEndian, &m.PORT)
	binary.Read(data, binary.LittleEndian, &m.RSSI)
}

//
// MESSAGE SERVO_OUTPUT_RAW
//
// MAVLINK_MSG_ID_SERVO_OUTPUT_RAW 36
//
// MAVLINK_MSG_ID_SERVO_OUTPUT_RAW_LEN 21
//
// MAVLINK_MSG_ID_SERVO_OUTPUT_RAW_CRC 222
//
//
type ServoOutputRaw struct {
	TIME_USEC  uint32 // Timestamp (microseconds since system boot)
	SERVO1_RAW uint16 // Servo output 1 value, in microseconds
	SERVO2_RAW uint16 // Servo output 2 value, in microseconds
	SERVO3_RAW uint16 // Servo output 3 value, in microseconds
	SERVO4_RAW uint16 // Servo output 4 value, in microseconds
	SERVO5_RAW uint16 // Servo output 5 value, in microseconds
	SERVO6_RAW uint16 // Servo output 6 value, in microseconds
	SERVO7_RAW uint16 // Servo output 7 value, in microseconds
	SERVO8_RAW uint16 // Servo output 8 value, in microseconds
	PORT       uint8  // Servo output port (set of 8 outputs = 1 port). Most MAVs will just use one, but this allows to encode more than 8 servos.
}

// NewServoOutputRaw returns a new ServoOutputRaw
func NewServoOutputRaw(TIME_USEC uint32, SERVO1_RAW uint16, SERVO2_RAW uint16, SERVO3_RAW uint16, SERVO4_RAW uint16, SERVO5_RAW uint16, SERVO6_RAW uint16, SERVO7_RAW uint16, SERVO8_RAW uint16, PORT uint8) *ServoOutputRaw {
	m := ServoOutputRaw{}
	m.TIME_USEC = TIME_USEC
	m.SERVO1_RAW = SERVO1_RAW
	m.SERVO2_RAW = SERVO2_RAW
	m.SERVO3_RAW = SERVO3_RAW
	m.SERVO4_RAW = SERVO4_RAW
	m.SERVO5_RAW = SERVO5_RAW
	m.SERVO6_RAW = SERVO6_RAW
	m.SERVO7_RAW = SERVO7_RAW
	m.SERVO8_RAW = SERVO8_RAW
	m.PORT = PORT
	return &m
}

// Id returns the ServoOutputRaw Message ID
func (*ServoOutputRaw) Id() uint8 {
	return 36
}

// Len returns the ServoOutputRaw Message Length
func (*ServoOutputRaw) Len() uint8 {
	return 21
}

// Crc returns the ServoOutputRaw Message CRC
func (*ServoOutputRaw) Crc() uint8 {
	return 222
}

// Pack returns a packed byte array which represents a ServoOutputRaw payload
func (m *ServoOutputRaw) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.SERVO1_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO2_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO3_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO4_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO5_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO6_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO7_RAW)
	binary.Write(data, binary.LittleEndian, m.SERVO8_RAW)
	binary.Write(data, binary.LittleEndian, m.PORT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ServoOutputRaw
func (m *ServoOutputRaw) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.SERVO1_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO2_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO3_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO4_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO5_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO6_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO7_RAW)
	binary.Read(data, binary.LittleEndian, &m.SERVO8_RAW)
	binary.Read(data, binary.LittleEndian, &m.PORT)
}

//
// MESSAGE MISSION_REQUEST_PARTIAL_LIST
//
// MAVLINK_MSG_ID_MISSION_REQUEST_PARTIAL_LIST 37
//
// MAVLINK_MSG_ID_MISSION_REQUEST_PARTIAL_LIST_LEN 6
//
// MAVLINK_MSG_ID_MISSION_REQUEST_PARTIAL_LIST_CRC 212
//
//
type MissionRequestPartialList struct {
	START_INDEX      int16 // Start index, 0 by default
	END_INDEX        int16 // End index, -1 by default (-1: send list to end). Else a valid index of the list
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewMissionRequestPartialList returns a new MissionRequestPartialList
func NewMissionRequestPartialList(START_INDEX int16, END_INDEX int16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionRequestPartialList {
	m := MissionRequestPartialList{}
	m.START_INDEX = START_INDEX
	m.END_INDEX = END_INDEX
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionRequestPartialList Message ID
func (*MissionRequestPartialList) Id() uint8 {
	return 37
}

// Len returns the MissionRequestPartialList Message Length
func (*MissionRequestPartialList) Len() uint8 {
	return 6
}

// Crc returns the MissionRequestPartialList Message CRC
func (*MissionRequestPartialList) Crc() uint8 {
	return 212
}

// Pack returns a packed byte array which represents a MissionRequestPartialList payload
func (m *MissionRequestPartialList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.START_INDEX)
	binary.Write(data, binary.LittleEndian, m.END_INDEX)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionRequestPartialList
func (m *MissionRequestPartialList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.START_INDEX)
	binary.Read(data, binary.LittleEndian, &m.END_INDEX)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_WRITE_PARTIAL_LIST
//
// MAVLINK_MSG_ID_MISSION_WRITE_PARTIAL_LIST 38
//
// MAVLINK_MSG_ID_MISSION_WRITE_PARTIAL_LIST_LEN 6
//
// MAVLINK_MSG_ID_MISSION_WRITE_PARTIAL_LIST_CRC 9
//
//
type MissionWritePartialList struct {
	START_INDEX      int16 // Start index, 0 by default and smaller / equal to the largest index of the current onboard list.
	END_INDEX        int16 // End index, equal or greater than start index.
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewMissionWritePartialList returns a new MissionWritePartialList
func NewMissionWritePartialList(START_INDEX int16, END_INDEX int16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionWritePartialList {
	m := MissionWritePartialList{}
	m.START_INDEX = START_INDEX
	m.END_INDEX = END_INDEX
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionWritePartialList Message ID
func (*MissionWritePartialList) Id() uint8 {
	return 38
}

// Len returns the MissionWritePartialList Message Length
func (*MissionWritePartialList) Len() uint8 {
	return 6
}

// Crc returns the MissionWritePartialList Message CRC
func (*MissionWritePartialList) Crc() uint8 {
	return 9
}

// Pack returns a packed byte array which represents a MissionWritePartialList payload
func (m *MissionWritePartialList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.START_INDEX)
	binary.Write(data, binary.LittleEndian, m.END_INDEX)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionWritePartialList
func (m *MissionWritePartialList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.START_INDEX)
	binary.Read(data, binary.LittleEndian, &m.END_INDEX)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_ITEM
//
// MAVLINK_MSG_ID_MISSION_ITEM 39
//
// MAVLINK_MSG_ID_MISSION_ITEM_LEN 37
//
// MAVLINK_MSG_ID_MISSION_ITEM_CRC 254
//
//
type MissionItem struct {
	PARAM1           float32 // PARAM1, see MAV_CMD enum
	PARAM2           float32 // PARAM2, see MAV_CMD enum
	PARAM3           float32 // PARAM3, see MAV_CMD enum
	PARAM4           float32 // PARAM4, see MAV_CMD enum
	X                float32 // PARAM5 / local: x position, global: latitude
	Y                float32 // PARAM6 / y position: global: longitude
	Z                float32 // PARAM7 / z position: global: altitude (relative or absolute, depending on frame.
	SEQ              uint16  // Sequence
	COMMAND          uint16  // The scheduled action for the MISSION. see MAV_CMD in common.xml MAVLink specs
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
	FRAME            uint8   // The coordinate system of the MISSION. see MAV_FRAME in mavlink_types.h
	CURRENT          uint8   // false:0, true:1
	AUTOCONTINUE     uint8   // autocontinue to next wp
}

// NewMissionItem returns a new MissionItem
func NewMissionItem(PARAM1 float32, PARAM2 float32, PARAM3 float32, PARAM4 float32, X float32, Y float32, Z float32, SEQ uint16, COMMAND uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, FRAME uint8, CURRENT uint8, AUTOCONTINUE uint8) *MissionItem {
	m := MissionItem{}
	m.PARAM1 = PARAM1
	m.PARAM2 = PARAM2
	m.PARAM3 = PARAM3
	m.PARAM4 = PARAM4
	m.X = X
	m.Y = Y
	m.Z = Z
	m.SEQ = SEQ
	m.COMMAND = COMMAND
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.FRAME = FRAME
	m.CURRENT = CURRENT
	m.AUTOCONTINUE = AUTOCONTINUE
	return &m
}

// Id returns the MissionItem Message ID
func (*MissionItem) Id() uint8 {
	return 39
}

// Len returns the MissionItem Message Length
func (*MissionItem) Len() uint8 {
	return 37
}

// Crc returns the MissionItem Message CRC
func (*MissionItem) Crc() uint8 {
	return 254
}

// Pack returns a packed byte array which represents a MissionItem payload
func (m *MissionItem) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.PARAM1)
	binary.Write(data, binary.LittleEndian, m.PARAM2)
	binary.Write(data, binary.LittleEndian, m.PARAM3)
	binary.Write(data, binary.LittleEndian, m.PARAM4)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	binary.Write(data, binary.LittleEndian, m.COMMAND)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.FRAME)
	binary.Write(data, binary.LittleEndian, m.CURRENT)
	binary.Write(data, binary.LittleEndian, m.AUTOCONTINUE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionItem
func (m *MissionItem) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.PARAM1)
	binary.Read(data, binary.LittleEndian, &m.PARAM2)
	binary.Read(data, binary.LittleEndian, &m.PARAM3)
	binary.Read(data, binary.LittleEndian, &m.PARAM4)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
	binary.Read(data, binary.LittleEndian, &m.COMMAND)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.FRAME)
	binary.Read(data, binary.LittleEndian, &m.CURRENT)
	binary.Read(data, binary.LittleEndian, &m.AUTOCONTINUE)
}

//
// MESSAGE MISSION_REQUEST
//
// MAVLINK_MSG_ID_MISSION_REQUEST 40
//
// MAVLINK_MSG_ID_MISSION_REQUEST_LEN 4
//
// MAVLINK_MSG_ID_MISSION_REQUEST_CRC 230
//
//
type MissionRequest struct {
	SEQ              uint16 // Sequence
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewMissionRequest returns a new MissionRequest
func NewMissionRequest(SEQ uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionRequest {
	m := MissionRequest{}
	m.SEQ = SEQ
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionRequest Message ID
func (*MissionRequest) Id() uint8 {
	return 40
}

// Len returns the MissionRequest Message Length
func (*MissionRequest) Len() uint8 {
	return 4
}

// Crc returns the MissionRequest Message CRC
func (*MissionRequest) Crc() uint8 {
	return 230
}

// Pack returns a packed byte array which represents a MissionRequest payload
func (m *MissionRequest) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionRequest
func (m *MissionRequest) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_SET_CURRENT
//
// MAVLINK_MSG_ID_MISSION_SET_CURRENT 41
//
// MAVLINK_MSG_ID_MISSION_SET_CURRENT_LEN 4
//
// MAVLINK_MSG_ID_MISSION_SET_CURRENT_CRC 28
//
//
type MissionSetCurrent struct {
	SEQ              uint16 // Sequence
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewMissionSetCurrent returns a new MissionSetCurrent
func NewMissionSetCurrent(SEQ uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionSetCurrent {
	m := MissionSetCurrent{}
	m.SEQ = SEQ
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionSetCurrent Message ID
func (*MissionSetCurrent) Id() uint8 {
	return 41
}

// Len returns the MissionSetCurrent Message Length
func (*MissionSetCurrent) Len() uint8 {
	return 4
}

// Crc returns the MissionSetCurrent Message CRC
func (*MissionSetCurrent) Crc() uint8 {
	return 28
}

// Pack returns a packed byte array which represents a MissionSetCurrent payload
func (m *MissionSetCurrent) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionSetCurrent
func (m *MissionSetCurrent) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_CURRENT
//
// MAVLINK_MSG_ID_MISSION_CURRENT 42
//
// MAVLINK_MSG_ID_MISSION_CURRENT_LEN 2
//
// MAVLINK_MSG_ID_MISSION_CURRENT_CRC 28
//
//
type MissionCurrent struct {
	SEQ uint16 // Sequence
}

// NewMissionCurrent returns a new MissionCurrent
func NewMissionCurrent(SEQ uint16) *MissionCurrent {
	m := MissionCurrent{}
	m.SEQ = SEQ
	return &m
}

// Id returns the MissionCurrent Message ID
func (*MissionCurrent) Id() uint8 {
	return 42
}

// Len returns the MissionCurrent Message Length
func (*MissionCurrent) Len() uint8 {
	return 2
}

// Crc returns the MissionCurrent Message CRC
func (*MissionCurrent) Crc() uint8 {
	return 28
}

// Pack returns a packed byte array which represents a MissionCurrent payload
func (m *MissionCurrent) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionCurrent
func (m *MissionCurrent) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
}

//
// MESSAGE MISSION_REQUEST_LIST
//
// MAVLINK_MSG_ID_MISSION_REQUEST_LIST 43
//
// MAVLINK_MSG_ID_MISSION_REQUEST_LIST_LEN 2
//
// MAVLINK_MSG_ID_MISSION_REQUEST_LIST_CRC 132
//
//
type MissionRequestList struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewMissionRequestList returns a new MissionRequestList
func NewMissionRequestList(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionRequestList {
	m := MissionRequestList{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionRequestList Message ID
func (*MissionRequestList) Id() uint8 {
	return 43
}

// Len returns the MissionRequestList Message Length
func (*MissionRequestList) Len() uint8 {
	return 2
}

// Crc returns the MissionRequestList Message CRC
func (*MissionRequestList) Crc() uint8 {
	return 132
}

// Pack returns a packed byte array which represents a MissionRequestList payload
func (m *MissionRequestList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionRequestList
func (m *MissionRequestList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_COUNT
//
// MAVLINK_MSG_ID_MISSION_COUNT 44
//
// MAVLINK_MSG_ID_MISSION_COUNT_LEN 4
//
// MAVLINK_MSG_ID_MISSION_COUNT_CRC 221
//
//
type MissionCount struct {
	COUNT            uint16 // Number of mission items in the sequence
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewMissionCount returns a new MissionCount
func NewMissionCount(COUNT uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionCount {
	m := MissionCount{}
	m.COUNT = COUNT
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionCount Message ID
func (*MissionCount) Id() uint8 {
	return 44
}

// Len returns the MissionCount Message Length
func (*MissionCount) Len() uint8 {
	return 4
}

// Crc returns the MissionCount Message CRC
func (*MissionCount) Crc() uint8 {
	return 221
}

// Pack returns a packed byte array which represents a MissionCount payload
func (m *MissionCount) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.COUNT)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionCount
func (m *MissionCount) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.COUNT)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_CLEAR_ALL
//
// MAVLINK_MSG_ID_MISSION_CLEAR_ALL 45
//
// MAVLINK_MSG_ID_MISSION_CLEAR_ALL_LEN 2
//
// MAVLINK_MSG_ID_MISSION_CLEAR_ALL_CRC 232
//
//
type MissionClearAll struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewMissionClearAll returns a new MissionClearAll
func NewMissionClearAll(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *MissionClearAll {
	m := MissionClearAll{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the MissionClearAll Message ID
func (*MissionClearAll) Id() uint8 {
	return 45
}

// Len returns the MissionClearAll Message Length
func (*MissionClearAll) Len() uint8 {
	return 2
}

// Crc returns the MissionClearAll Message CRC
func (*MissionClearAll) Crc() uint8 {
	return 232
}

// Pack returns a packed byte array which represents a MissionClearAll payload
func (m *MissionClearAll) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionClearAll
func (m *MissionClearAll) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE MISSION_ITEM_REACHED
//
// MAVLINK_MSG_ID_MISSION_ITEM_REACHED 46
//
// MAVLINK_MSG_ID_MISSION_ITEM_REACHED_LEN 2
//
// MAVLINK_MSG_ID_MISSION_ITEM_REACHED_CRC 11
//
//
type MissionItemReached struct {
	SEQ uint16 // Sequence
}

// NewMissionItemReached returns a new MissionItemReached
func NewMissionItemReached(SEQ uint16) *MissionItemReached {
	m := MissionItemReached{}
	m.SEQ = SEQ
	return &m
}

// Id returns the MissionItemReached Message ID
func (*MissionItemReached) Id() uint8 {
	return 46
}

// Len returns the MissionItemReached Message Length
func (*MissionItemReached) Len() uint8 {
	return 2
}

// Crc returns the MissionItemReached Message CRC
func (*MissionItemReached) Crc() uint8 {
	return 11
}

// Pack returns a packed byte array which represents a MissionItemReached payload
func (m *MissionItemReached) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEQ)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionItemReached
func (m *MissionItemReached) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEQ)
}

//
// MESSAGE MISSION_ACK
//
// MAVLINK_MSG_ID_MISSION_ACK 47
//
// MAVLINK_MSG_ID_MISSION_ACK_LEN 3
//
// MAVLINK_MSG_ID_MISSION_ACK_CRC 153
//
//
type MissionAck struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
	TYPE             uint8 // See MAV_MISSION_RESULT enum
}

// NewMissionAck returns a new MissionAck
func NewMissionAck(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, TYPE uint8) *MissionAck {
	m := MissionAck{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.TYPE = TYPE
	return &m
}

// Id returns the MissionAck Message ID
func (*MissionAck) Id() uint8 {
	return 47
}

// Len returns the MissionAck Message Length
func (*MissionAck) Len() uint8 {
	return 3
}

// Crc returns the MissionAck Message CRC
func (*MissionAck) Crc() uint8 {
	return 153
}

// Pack returns a packed byte array which represents a MissionAck payload
func (m *MissionAck) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.TYPE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MissionAck
func (m *MissionAck) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.TYPE)
}

//
// MESSAGE SET_GPS_GLOBAL_ORIGIN
//
// MAVLINK_MSG_ID_SET_GPS_GLOBAL_ORIGIN 48
//
// MAVLINK_MSG_ID_SET_GPS_GLOBAL_ORIGIN_LEN 13
//
// MAVLINK_MSG_ID_SET_GPS_GLOBAL_ORIGIN_CRC 41
//
//
type SetGpsGlobalOrigin struct {
	LATITUDE      int32 // Latitude (WGS84), in degrees * 1E7
	LONGITUDE     int32 // Longitude (WGS84, in degrees * 1E7
	ALTITUDE      int32 // Altitude (WGS84), in meters * 1000 (positive for up)
	TARGET_SYSTEM uint8 // System ID
}

// NewSetGpsGlobalOrigin returns a new SetGpsGlobalOrigin
func NewSetGpsGlobalOrigin(LATITUDE int32, LONGITUDE int32, ALTITUDE int32, TARGET_SYSTEM uint8) *SetGpsGlobalOrigin {
	m := SetGpsGlobalOrigin{}
	m.LATITUDE = LATITUDE
	m.LONGITUDE = LONGITUDE
	m.ALTITUDE = ALTITUDE
	m.TARGET_SYSTEM = TARGET_SYSTEM
	return &m
}

// Id returns the SetGpsGlobalOrigin Message ID
func (*SetGpsGlobalOrigin) Id() uint8 {
	return 48
}

// Len returns the SetGpsGlobalOrigin Message Length
func (*SetGpsGlobalOrigin) Len() uint8 {
	return 13
}

// Crc returns the SetGpsGlobalOrigin Message CRC
func (*SetGpsGlobalOrigin) Crc() uint8 {
	return 41
}

// Pack returns a packed byte array which represents a SetGpsGlobalOrigin payload
func (m *SetGpsGlobalOrigin) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LATITUDE)
	binary.Write(data, binary.LittleEndian, m.LONGITUDE)
	binary.Write(data, binary.LittleEndian, m.ALTITUDE)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetGpsGlobalOrigin
func (m *SetGpsGlobalOrigin) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LATITUDE)
	binary.Read(data, binary.LittleEndian, &m.LONGITUDE)
	binary.Read(data, binary.LittleEndian, &m.ALTITUDE)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
}

//
// MESSAGE GPS_GLOBAL_ORIGIN
//
// MAVLINK_MSG_ID_GPS_GLOBAL_ORIGIN 49
//
// MAVLINK_MSG_ID_GPS_GLOBAL_ORIGIN_LEN 12
//
// MAVLINK_MSG_ID_GPS_GLOBAL_ORIGIN_CRC 39
//
//
type GpsGlobalOrigin struct {
	LATITUDE  int32 // Latitude (WGS84), in degrees * 1E7
	LONGITUDE int32 // Longitude (WGS84), in degrees * 1E7
	ALTITUDE  int32 // Altitude (WGS84), in meters * 1000 (positive for up)
}

// NewGpsGlobalOrigin returns a new GpsGlobalOrigin
func NewGpsGlobalOrigin(LATITUDE int32, LONGITUDE int32, ALTITUDE int32) *GpsGlobalOrigin {
	m := GpsGlobalOrigin{}
	m.LATITUDE = LATITUDE
	m.LONGITUDE = LONGITUDE
	m.ALTITUDE = ALTITUDE
	return &m
}

// Id returns the GpsGlobalOrigin Message ID
func (*GpsGlobalOrigin) Id() uint8 {
	return 49
}

// Len returns the GpsGlobalOrigin Message Length
func (*GpsGlobalOrigin) Len() uint8 {
	return 12
}

// Crc returns the GpsGlobalOrigin Message CRC
func (*GpsGlobalOrigin) Crc() uint8 {
	return 39
}

// Pack returns a packed byte array which represents a GpsGlobalOrigin payload
func (m *GpsGlobalOrigin) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LATITUDE)
	binary.Write(data, binary.LittleEndian, m.LONGITUDE)
	binary.Write(data, binary.LittleEndian, m.ALTITUDE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GpsGlobalOrigin
func (m *GpsGlobalOrigin) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LATITUDE)
	binary.Read(data, binary.LittleEndian, &m.LONGITUDE)
	binary.Read(data, binary.LittleEndian, &m.ALTITUDE)
}

//
// MESSAGE SET_LOCAL_POSITION_SETPOINT
//
// MAVLINK_MSG_ID_SET_LOCAL_POSITION_SETPOINT 50
//
// MAVLINK_MSG_ID_SET_LOCAL_POSITION_SETPOINT_LEN 19
//
// MAVLINK_MSG_ID_SET_LOCAL_POSITION_SETPOINT_CRC 214
//
//
type SetLocalPositionSetpoint struct {
	X                float32 // x position
	Y                float32 // y position
	Z                float32 // z position
	YAW              float32 // Desired yaw angle
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
	COORDINATE_FRAME uint8   // Coordinate frame - valid values are only MAV_FRAME_LOCAL_NED or MAV_FRAME_LOCAL_ENU
}

// NewSetLocalPositionSetpoint returns a new SetLocalPositionSetpoint
func NewSetLocalPositionSetpoint(X float32, Y float32, Z float32, YAW float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, COORDINATE_FRAME uint8) *SetLocalPositionSetpoint {
	m := SetLocalPositionSetpoint{}
	m.X = X
	m.Y = Y
	m.Z = Z
	m.YAW = YAW
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.COORDINATE_FRAME = COORDINATE_FRAME
	return &m
}

// Id returns the SetLocalPositionSetpoint Message ID
func (*SetLocalPositionSetpoint) Id() uint8 {
	return 50
}

// Len returns the SetLocalPositionSetpoint Message Length
func (*SetLocalPositionSetpoint) Len() uint8 {
	return 19
}

// Crc returns the SetLocalPositionSetpoint Message CRC
func (*SetLocalPositionSetpoint) Crc() uint8 {
	return 214
}

// Pack returns a packed byte array which represents a SetLocalPositionSetpoint payload
func (m *SetLocalPositionSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.COORDINATE_FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetLocalPositionSetpoint
func (m *SetLocalPositionSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.COORDINATE_FRAME)
}

//
// MESSAGE LOCAL_POSITION_SETPOINT
//
// MAVLINK_MSG_ID_LOCAL_POSITION_SETPOINT 51
//
// MAVLINK_MSG_ID_LOCAL_POSITION_SETPOINT_LEN 17
//
// MAVLINK_MSG_ID_LOCAL_POSITION_SETPOINT_CRC 223
//
//
type LocalPositionSetpoint struct {
	X                float32 // x position
	Y                float32 // y position
	Z                float32 // z position
	YAW              float32 // Desired yaw angle
	COORDINATE_FRAME uint8   // Coordinate frame - valid values are only MAV_FRAME_LOCAL_NED or MAV_FRAME_LOCAL_ENU
}

// NewLocalPositionSetpoint returns a new LocalPositionSetpoint
func NewLocalPositionSetpoint(X float32, Y float32, Z float32, YAW float32, COORDINATE_FRAME uint8) *LocalPositionSetpoint {
	m := LocalPositionSetpoint{}
	m.X = X
	m.Y = Y
	m.Z = Z
	m.YAW = YAW
	m.COORDINATE_FRAME = COORDINATE_FRAME
	return &m
}

// Id returns the LocalPositionSetpoint Message ID
func (*LocalPositionSetpoint) Id() uint8 {
	return 51
}

// Len returns the LocalPositionSetpoint Message Length
func (*LocalPositionSetpoint) Len() uint8 {
	return 17
}

// Crc returns the LocalPositionSetpoint Message CRC
func (*LocalPositionSetpoint) Crc() uint8 {
	return 223
}

// Pack returns a packed byte array which represents a LocalPositionSetpoint payload
func (m *LocalPositionSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.COORDINATE_FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LocalPositionSetpoint
func (m *LocalPositionSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.COORDINATE_FRAME)
}

//
// MESSAGE GLOBAL_POSITION_SETPOINT_INT
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_INT 52
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_INT_LEN 15
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_INT_CRC 141
//
//
type GlobalPositionSetpointInt struct {
	LATITUDE         int32 // Latitude (WGS84), in degrees * 1E7
	LONGITUDE        int32 // Longitude (WGS84), in degrees * 1E7
	ALTITUDE         int32 // Altitude (WGS84), in meters * 1000 (positive for up)
	YAW              int16 // Desired yaw angle in degrees * 100
	COORDINATE_FRAME uint8 // Coordinate frame - valid values are only MAV_FRAME_GLOBAL or MAV_FRAME_GLOBAL_RELATIVE_ALT
}

// NewGlobalPositionSetpointInt returns a new GlobalPositionSetpointInt
func NewGlobalPositionSetpointInt(LATITUDE int32, LONGITUDE int32, ALTITUDE int32, YAW int16, COORDINATE_FRAME uint8) *GlobalPositionSetpointInt {
	m := GlobalPositionSetpointInt{}
	m.LATITUDE = LATITUDE
	m.LONGITUDE = LONGITUDE
	m.ALTITUDE = ALTITUDE
	m.YAW = YAW
	m.COORDINATE_FRAME = COORDINATE_FRAME
	return &m
}

// Id returns the GlobalPositionSetpointInt Message ID
func (*GlobalPositionSetpointInt) Id() uint8 {
	return 52
}

// Len returns the GlobalPositionSetpointInt Message Length
func (*GlobalPositionSetpointInt) Len() uint8 {
	return 15
}

// Crc returns the GlobalPositionSetpointInt Message CRC
func (*GlobalPositionSetpointInt) Crc() uint8 {
	return 141
}

// Pack returns a packed byte array which represents a GlobalPositionSetpointInt payload
func (m *GlobalPositionSetpointInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LATITUDE)
	binary.Write(data, binary.LittleEndian, m.LONGITUDE)
	binary.Write(data, binary.LittleEndian, m.ALTITUDE)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.COORDINATE_FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GlobalPositionSetpointInt
func (m *GlobalPositionSetpointInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LATITUDE)
	binary.Read(data, binary.LittleEndian, &m.LONGITUDE)
	binary.Read(data, binary.LittleEndian, &m.ALTITUDE)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.COORDINATE_FRAME)
}

//
// MESSAGE SET_GLOBAL_POSITION_SETPOINT_INT
//
// MAVLINK_MSG_ID_SET_GLOBAL_POSITION_SETPOINT_INT 53
//
// MAVLINK_MSG_ID_SET_GLOBAL_POSITION_SETPOINT_INT_LEN 15
//
// MAVLINK_MSG_ID_SET_GLOBAL_POSITION_SETPOINT_INT_CRC 33
//
//
type SetGlobalPositionSetpointInt struct {
	LATITUDE         int32 // Latitude (WGS84), in degrees * 1E7
	LONGITUDE        int32 // Longitude (WGS84), in degrees * 1E7
	ALTITUDE         int32 // Altitude (WGS84), in meters * 1000 (positive for up)
	YAW              int16 // Desired yaw angle in degrees * 100
	COORDINATE_FRAME uint8 // Coordinate frame - valid values are only MAV_FRAME_GLOBAL or MAV_FRAME_GLOBAL_RELATIVE_ALT
}

// NewSetGlobalPositionSetpointInt returns a new SetGlobalPositionSetpointInt
func NewSetGlobalPositionSetpointInt(LATITUDE int32, LONGITUDE int32, ALTITUDE int32, YAW int16, COORDINATE_FRAME uint8) *SetGlobalPositionSetpointInt {
	m := SetGlobalPositionSetpointInt{}
	m.LATITUDE = LATITUDE
	m.LONGITUDE = LONGITUDE
	m.ALTITUDE = ALTITUDE
	m.YAW = YAW
	m.COORDINATE_FRAME = COORDINATE_FRAME
	return &m
}

// Id returns the SetGlobalPositionSetpointInt Message ID
func (*SetGlobalPositionSetpointInt) Id() uint8 {
	return 53
}

// Len returns the SetGlobalPositionSetpointInt Message Length
func (*SetGlobalPositionSetpointInt) Len() uint8 {
	return 15
}

// Crc returns the SetGlobalPositionSetpointInt Message CRC
func (*SetGlobalPositionSetpointInt) Crc() uint8 {
	return 33
}

// Pack returns a packed byte array which represents a SetGlobalPositionSetpointInt payload
func (m *SetGlobalPositionSetpointInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LATITUDE)
	binary.Write(data, binary.LittleEndian, m.LONGITUDE)
	binary.Write(data, binary.LittleEndian, m.ALTITUDE)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.COORDINATE_FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetGlobalPositionSetpointInt
func (m *SetGlobalPositionSetpointInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LATITUDE)
	binary.Read(data, binary.LittleEndian, &m.LONGITUDE)
	binary.Read(data, binary.LittleEndian, &m.ALTITUDE)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.COORDINATE_FRAME)
}

//
// MESSAGE SAFETY_SET_ALLOWED_AREA
//
// MAVLINK_MSG_ID_SAFETY_SET_ALLOWED_AREA 54
//
// MAVLINK_MSG_ID_SAFETY_SET_ALLOWED_AREA_LEN 27
//
// MAVLINK_MSG_ID_SAFETY_SET_ALLOWED_AREA_CRC 15
//
//
type SafetySetAllowedArea struct {
	P1X              float32 // x position 1 / Latitude 1
	P1Y              float32 // y position 1 / Longitude 1
	P1Z              float32 // z position 1 / Altitude 1
	P2X              float32 // x position 2 / Latitude 2
	P2Y              float32 // y position 2 / Longitude 2
	P2Z              float32 // z position 2 / Altitude 2
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
	FRAME            uint8   // Coordinate frame, as defined by MAV_FRAME enum in mavlink_types.h. Can be either global, GPS, right-handed with Z axis up or local, right handed, Z axis down.
}

// NewSafetySetAllowedArea returns a new SafetySetAllowedArea
func NewSafetySetAllowedArea(P1X float32, P1Y float32, P1Z float32, P2X float32, P2Y float32, P2Z float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, FRAME uint8) *SafetySetAllowedArea {
	m := SafetySetAllowedArea{}
	m.P1X = P1X
	m.P1Y = P1Y
	m.P1Z = P1Z
	m.P2X = P2X
	m.P2Y = P2Y
	m.P2Z = P2Z
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.FRAME = FRAME
	return &m
}

// Id returns the SafetySetAllowedArea Message ID
func (*SafetySetAllowedArea) Id() uint8 {
	return 54
}

// Len returns the SafetySetAllowedArea Message Length
func (*SafetySetAllowedArea) Len() uint8 {
	return 27
}

// Crc returns the SafetySetAllowedArea Message CRC
func (*SafetySetAllowedArea) Crc() uint8 {
	return 15
}

// Pack returns a packed byte array which represents a SafetySetAllowedArea payload
func (m *SafetySetAllowedArea) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.P1X)
	binary.Write(data, binary.LittleEndian, m.P1Y)
	binary.Write(data, binary.LittleEndian, m.P1Z)
	binary.Write(data, binary.LittleEndian, m.P2X)
	binary.Write(data, binary.LittleEndian, m.P2Y)
	binary.Write(data, binary.LittleEndian, m.P2Z)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SafetySetAllowedArea
func (m *SafetySetAllowedArea) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.P1X)
	binary.Read(data, binary.LittleEndian, &m.P1Y)
	binary.Read(data, binary.LittleEndian, &m.P1Z)
	binary.Read(data, binary.LittleEndian, &m.P2X)
	binary.Read(data, binary.LittleEndian, &m.P2Y)
	binary.Read(data, binary.LittleEndian, &m.P2Z)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.FRAME)
}

//
// MESSAGE SAFETY_ALLOWED_AREA
//
// MAVLINK_MSG_ID_SAFETY_ALLOWED_AREA 55
//
// MAVLINK_MSG_ID_SAFETY_ALLOWED_AREA_LEN 25
//
// MAVLINK_MSG_ID_SAFETY_ALLOWED_AREA_CRC 3
//
//
type SafetyAllowedArea struct {
	P1X   float32 // x position 1 / Latitude 1
	P1Y   float32 // y position 1 / Longitude 1
	P1Z   float32 // z position 1 / Altitude 1
	P2X   float32 // x position 2 / Latitude 2
	P2Y   float32 // y position 2 / Longitude 2
	P2Z   float32 // z position 2 / Altitude 2
	FRAME uint8   // Coordinate frame, as defined by MAV_FRAME enum in mavlink_types.h. Can be either global, GPS, right-handed with Z axis up or local, right handed, Z axis down.
}

// NewSafetyAllowedArea returns a new SafetyAllowedArea
func NewSafetyAllowedArea(P1X float32, P1Y float32, P1Z float32, P2X float32, P2Y float32, P2Z float32, FRAME uint8) *SafetyAllowedArea {
	m := SafetyAllowedArea{}
	m.P1X = P1X
	m.P1Y = P1Y
	m.P1Z = P1Z
	m.P2X = P2X
	m.P2Y = P2Y
	m.P2Z = P2Z
	m.FRAME = FRAME
	return &m
}

// Id returns the SafetyAllowedArea Message ID
func (*SafetyAllowedArea) Id() uint8 {
	return 55
}

// Len returns the SafetyAllowedArea Message Length
func (*SafetyAllowedArea) Len() uint8 {
	return 25
}

// Crc returns the SafetyAllowedArea Message CRC
func (*SafetyAllowedArea) Crc() uint8 {
	return 3
}

// Pack returns a packed byte array which represents a SafetyAllowedArea payload
func (m *SafetyAllowedArea) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.P1X)
	binary.Write(data, binary.LittleEndian, m.P1Y)
	binary.Write(data, binary.LittleEndian, m.P1Z)
	binary.Write(data, binary.LittleEndian, m.P2X)
	binary.Write(data, binary.LittleEndian, m.P2Y)
	binary.Write(data, binary.LittleEndian, m.P2Z)
	binary.Write(data, binary.LittleEndian, m.FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SafetyAllowedArea
func (m *SafetyAllowedArea) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.P1X)
	binary.Read(data, binary.LittleEndian, &m.P1Y)
	binary.Read(data, binary.LittleEndian, &m.P1Z)
	binary.Read(data, binary.LittleEndian, &m.P2X)
	binary.Read(data, binary.LittleEndian, &m.P2Y)
	binary.Read(data, binary.LittleEndian, &m.P2Z)
	binary.Read(data, binary.LittleEndian, &m.FRAME)
}

//
// MESSAGE SET_ROLL_PITCH_YAW_THRUST
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_THRUST 56
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_THRUST_LEN 18
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_THRUST_CRC 100
//
//
type SetRollPitchYawThrust struct {
	ROLL             float32 // Desired roll angle in radians
	PITCH            float32 // Desired pitch angle in radians
	YAW              float32 // Desired yaw angle in radians
	THRUST           float32 // Collective thrust, normalized to 0 .. 1
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
}

// NewSetRollPitchYawThrust returns a new SetRollPitchYawThrust
func NewSetRollPitchYawThrust(ROLL float32, PITCH float32, YAW float32, THRUST float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *SetRollPitchYawThrust {
	m := SetRollPitchYawThrust{}
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.THRUST = THRUST
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the SetRollPitchYawThrust Message ID
func (*SetRollPitchYawThrust) Id() uint8 {
	return 56
}

// Len returns the SetRollPitchYawThrust Message Length
func (*SetRollPitchYawThrust) Len() uint8 {
	return 18
}

// Crc returns the SetRollPitchYawThrust Message CRC
func (*SetRollPitchYawThrust) Crc() uint8 {
	return 100
}

// Pack returns a packed byte array which represents a SetRollPitchYawThrust payload
func (m *SetRollPitchYawThrust) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetRollPitchYawThrust
func (m *SetRollPitchYawThrust) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE SET_ROLL_PITCH_YAW_SPEED_THRUST
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_SPEED_THRUST 57
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_SPEED_THRUST_LEN 18
//
// MAVLINK_MSG_ID_SET_ROLL_PITCH_YAW_SPEED_THRUST_CRC 24
//
//
type SetRollPitchYawSpeedThrust struct {
	ROLL_SPEED       float32 // Desired roll angular speed in rad/s
	PITCH_SPEED      float32 // Desired pitch angular speed in rad/s
	YAW_SPEED        float32 // Desired yaw angular speed in rad/s
	THRUST           float32 // Collective thrust, normalized to 0 .. 1
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
}

// NewSetRollPitchYawSpeedThrust returns a new SetRollPitchYawSpeedThrust
func NewSetRollPitchYawSpeedThrust(ROLL_SPEED float32, PITCH_SPEED float32, YAW_SPEED float32, THRUST float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *SetRollPitchYawSpeedThrust {
	m := SetRollPitchYawSpeedThrust{}
	m.ROLL_SPEED = ROLL_SPEED
	m.PITCH_SPEED = PITCH_SPEED
	m.YAW_SPEED = YAW_SPEED
	m.THRUST = THRUST
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the SetRollPitchYawSpeedThrust Message ID
func (*SetRollPitchYawSpeedThrust) Id() uint8 {
	return 57
}

// Len returns the SetRollPitchYawSpeedThrust Message Length
func (*SetRollPitchYawSpeedThrust) Len() uint8 {
	return 18
}

// Crc returns the SetRollPitchYawSpeedThrust Message CRC
func (*SetRollPitchYawSpeedThrust) Crc() uint8 {
	return 24
}

// Pack returns a packed byte array which represents a SetRollPitchYawSpeedThrust payload
func (m *SetRollPitchYawSpeedThrust) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ROLL_SPEED)
	binary.Write(data, binary.LittleEndian, m.PITCH_SPEED)
	binary.Write(data, binary.LittleEndian, m.YAW_SPEED)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetRollPitchYawSpeedThrust
func (m *SetRollPitchYawSpeedThrust) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ROLL_SPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCH_SPEED)
	binary.Read(data, binary.LittleEndian, &m.YAW_SPEED)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE ROLL_PITCH_YAW_THRUST_SETPOINT
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_THRUST_SETPOINT 58
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_THRUST_SETPOINT_LEN 20
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_THRUST_SETPOINT_CRC 239
//
//
type RollPitchYawThrustSetpoint struct {
	TIME_BOOT_MS uint32  // Timestamp in milliseconds since system boot
	ROLL         float32 // Desired roll angle in radians
	PITCH        float32 // Desired pitch angle in radians
	YAW          float32 // Desired yaw angle in radians
	THRUST       float32 // Collective thrust, normalized to 0 .. 1
}

// NewRollPitchYawThrustSetpoint returns a new RollPitchYawThrustSetpoint
func NewRollPitchYawThrustSetpoint(TIME_BOOT_MS uint32, ROLL float32, PITCH float32, YAW float32, THRUST float32) *RollPitchYawThrustSetpoint {
	m := RollPitchYawThrustSetpoint{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.THRUST = THRUST
	return &m
}

// Id returns the RollPitchYawThrustSetpoint Message ID
func (*RollPitchYawThrustSetpoint) Id() uint8 {
	return 58
}

// Len returns the RollPitchYawThrustSetpoint Message Length
func (*RollPitchYawThrustSetpoint) Len() uint8 {
	return 20
}

// Crc returns the RollPitchYawThrustSetpoint Message CRC
func (*RollPitchYawThrustSetpoint) Crc() uint8 {
	return 239
}

// Pack returns a packed byte array which represents a RollPitchYawThrustSetpoint payload
func (m *RollPitchYawThrustSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RollPitchYawThrustSetpoint
func (m *RollPitchYawThrustSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
}

//
// MESSAGE ROLL_PITCH_YAW_SPEED_THRUST_SETPOINT
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_SPEED_THRUST_SETPOINT 59
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_SPEED_THRUST_SETPOINT_LEN 20
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_SPEED_THRUST_SETPOINT_CRC 238
//
//
type RollPitchYawSpeedThrustSetpoint struct {
	TIME_BOOT_MS uint32  // Timestamp in milliseconds since system boot
	ROLL_SPEED   float32 // Desired roll angular speed in rad/s
	PITCH_SPEED  float32 // Desired pitch angular speed in rad/s
	YAW_SPEED    float32 // Desired yaw angular speed in rad/s
	THRUST       float32 // Collective thrust, normalized to 0 .. 1
}

// NewRollPitchYawSpeedThrustSetpoint returns a new RollPitchYawSpeedThrustSetpoint
func NewRollPitchYawSpeedThrustSetpoint(TIME_BOOT_MS uint32, ROLL_SPEED float32, PITCH_SPEED float32, YAW_SPEED float32, THRUST float32) *RollPitchYawSpeedThrustSetpoint {
	m := RollPitchYawSpeedThrustSetpoint{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.ROLL_SPEED = ROLL_SPEED
	m.PITCH_SPEED = PITCH_SPEED
	m.YAW_SPEED = YAW_SPEED
	m.THRUST = THRUST
	return &m
}

// Id returns the RollPitchYawSpeedThrustSetpoint Message ID
func (*RollPitchYawSpeedThrustSetpoint) Id() uint8 {
	return 59
}

// Len returns the RollPitchYawSpeedThrustSetpoint Message Length
func (*RollPitchYawSpeedThrustSetpoint) Len() uint8 {
	return 20
}

// Crc returns the RollPitchYawSpeedThrustSetpoint Message CRC
func (*RollPitchYawSpeedThrustSetpoint) Crc() uint8 {
	return 238
}

// Pack returns a packed byte array which represents a RollPitchYawSpeedThrustSetpoint payload
func (m *RollPitchYawSpeedThrustSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.ROLL_SPEED)
	binary.Write(data, binary.LittleEndian, m.PITCH_SPEED)
	binary.Write(data, binary.LittleEndian, m.YAW_SPEED)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RollPitchYawSpeedThrustSetpoint
func (m *RollPitchYawSpeedThrustSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.ROLL_SPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCH_SPEED)
	binary.Read(data, binary.LittleEndian, &m.YAW_SPEED)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
}

//
// MESSAGE SET_QUAD_MOTORS_SETPOINT
//
// MAVLINK_MSG_ID_SET_QUAD_MOTORS_SETPOINT 60
//
// MAVLINK_MSG_ID_SET_QUAD_MOTORS_SETPOINT_LEN 9
//
// MAVLINK_MSG_ID_SET_QUAD_MOTORS_SETPOINT_CRC 30
//
//
type SetQuadMotorsSetpoint struct {
	MOTOR_FRONT_NW uint16 // Front motor in + configuration, front left motor in x configuration
	MOTOR_RIGHT_NE uint16 // Right motor in + configuration, front right motor in x configuration
	MOTOR_BACK_SE  uint16 // Back motor in + configuration, back right motor in x configuration
	MOTOR_LEFT_SW  uint16 // Left motor in + configuration, back left motor in x configuration
	TARGET_SYSTEM  uint8  // System ID of the system that should set these motor commands
}

// NewSetQuadMotorsSetpoint returns a new SetQuadMotorsSetpoint
func NewSetQuadMotorsSetpoint(MOTOR_FRONT_NW uint16, MOTOR_RIGHT_NE uint16, MOTOR_BACK_SE uint16, MOTOR_LEFT_SW uint16, TARGET_SYSTEM uint8) *SetQuadMotorsSetpoint {
	m := SetQuadMotorsSetpoint{}
	m.MOTOR_FRONT_NW = MOTOR_FRONT_NW
	m.MOTOR_RIGHT_NE = MOTOR_RIGHT_NE
	m.MOTOR_BACK_SE = MOTOR_BACK_SE
	m.MOTOR_LEFT_SW = MOTOR_LEFT_SW
	m.TARGET_SYSTEM = TARGET_SYSTEM
	return &m
}

// Id returns the SetQuadMotorsSetpoint Message ID
func (*SetQuadMotorsSetpoint) Id() uint8 {
	return 60
}

// Len returns the SetQuadMotorsSetpoint Message Length
func (*SetQuadMotorsSetpoint) Len() uint8 {
	return 9
}

// Crc returns the SetQuadMotorsSetpoint Message CRC
func (*SetQuadMotorsSetpoint) Crc() uint8 {
	return 30
}

// Pack returns a packed byte array which represents a SetQuadMotorsSetpoint payload
func (m *SetQuadMotorsSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.MOTOR_FRONT_NW)
	binary.Write(data, binary.LittleEndian, m.MOTOR_RIGHT_NE)
	binary.Write(data, binary.LittleEndian, m.MOTOR_BACK_SE)
	binary.Write(data, binary.LittleEndian, m.MOTOR_LEFT_SW)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetQuadMotorsSetpoint
func (m *SetQuadMotorsSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.MOTOR_FRONT_NW)
	binary.Read(data, binary.LittleEndian, &m.MOTOR_RIGHT_NE)
	binary.Read(data, binary.LittleEndian, &m.MOTOR_BACK_SE)
	binary.Read(data, binary.LittleEndian, &m.MOTOR_LEFT_SW)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
}

//
// MESSAGE SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST 61
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_LEN 34
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_CRC 240
//
//
type SetQuadSwarmRollPitchYawThrust struct {
	ROLL   [4]int16  // Desired roll angle in radians +-PI (+-INT16_MAX)
	PITCH  [4]int16  // Desired pitch angle in radians +-PI (+-INT16_MAX)
	YAW    [4]int16  // Desired yaw angle in radians, scaled to int16 +-PI (+-INT16_MAX)
	THRUST [4]uint16 // Collective thrust, scaled to uint16 (0..UINT16_MAX)
	GROUP  uint8     // ID of the quadrotor group (0 - 255, up to 256 groups supported)
	MODE   uint8     // ID of the flight mode (0 - 255, up to 256 modes supported)
}

// NewSetQuadSwarmRollPitchYawThrust returns a new SetQuadSwarmRollPitchYawThrust
func NewSetQuadSwarmRollPitchYawThrust(ROLL [4]int16, PITCH [4]int16, YAW [4]int16, THRUST [4]uint16, GROUP uint8, MODE uint8) *SetQuadSwarmRollPitchYawThrust {
	m := SetQuadSwarmRollPitchYawThrust{}
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.THRUST = THRUST
	m.GROUP = GROUP
	m.MODE = MODE
	return &m
}

// Id returns the SetQuadSwarmRollPitchYawThrust Message ID
func (*SetQuadSwarmRollPitchYawThrust) Id() uint8 {
	return 61
}

// Len returns the SetQuadSwarmRollPitchYawThrust Message Length
func (*SetQuadSwarmRollPitchYawThrust) Len() uint8 {
	return 34
}

// Crc returns the SetQuadSwarmRollPitchYawThrust Message CRC
func (*SetQuadSwarmRollPitchYawThrust) Crc() uint8 {
	return 240
}

// Pack returns a packed byte array which represents a SetQuadSwarmRollPitchYawThrust payload
func (m *SetQuadSwarmRollPitchYawThrust) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.GROUP)
	binary.Write(data, binary.LittleEndian, m.MODE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetQuadSwarmRollPitchYawThrust
func (m *SetQuadSwarmRollPitchYawThrust) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.GROUP)
	binary.Read(data, binary.LittleEndian, &m.MODE)
}

const (
	MAVLINK_MSG_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_FIELD_roll_LEN   = 4
	MAVLINK_MSG_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_FIELD_pitch_LEN  = 4
	MAVLINK_MSG_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_FIELD_yaw_LEN    = 4
	MAVLINK_MSG_SET_QUAD_SWARM_ROLL_PITCH_YAW_THRUST_FIELD_thrust_LEN = 4
)

//
// MESSAGE NAV_CONTROLLER_OUTPUT
//
// MAVLINK_MSG_ID_NAV_CONTROLLER_OUTPUT 62
//
// MAVLINK_MSG_ID_NAV_CONTROLLER_OUTPUT_LEN 26
//
// MAVLINK_MSG_ID_NAV_CONTROLLER_OUTPUT_CRC 183
//
//
type NavControllerOutput struct {
	NAV_ROLL       float32 // Current desired roll in degrees
	NAV_PITCH      float32 // Current desired pitch in degrees
	ALT_ERROR      float32 // Current altitude error in meters
	ASPD_ERROR     float32 // Current airspeed error in meters/second
	XTRACK_ERROR   float32 // Current crosstrack error on x-y plane in meters
	NAV_BEARING    int16   // Current desired heading in degrees
	TARGET_BEARING int16   // Bearing to current MISSION/target in degrees
	WP_DIST        uint16  // Distance to active MISSION in meters
}

// NewNavControllerOutput returns a new NavControllerOutput
func NewNavControllerOutput(NAV_ROLL float32, NAV_PITCH float32, ALT_ERROR float32, ASPD_ERROR float32, XTRACK_ERROR float32, NAV_BEARING int16, TARGET_BEARING int16, WP_DIST uint16) *NavControllerOutput {
	m := NavControllerOutput{}
	m.NAV_ROLL = NAV_ROLL
	m.NAV_PITCH = NAV_PITCH
	m.ALT_ERROR = ALT_ERROR
	m.ASPD_ERROR = ASPD_ERROR
	m.XTRACK_ERROR = XTRACK_ERROR
	m.NAV_BEARING = NAV_BEARING
	m.TARGET_BEARING = TARGET_BEARING
	m.WP_DIST = WP_DIST
	return &m
}

// Id returns the NavControllerOutput Message ID
func (*NavControllerOutput) Id() uint8 {
	return 62
}

// Len returns the NavControllerOutput Message Length
func (*NavControllerOutput) Len() uint8 {
	return 26
}

// Crc returns the NavControllerOutput Message CRC
func (*NavControllerOutput) Crc() uint8 {
	return 183
}

// Pack returns a packed byte array which represents a NavControllerOutput payload
func (m *NavControllerOutput) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.NAV_ROLL)
	binary.Write(data, binary.LittleEndian, m.NAV_PITCH)
	binary.Write(data, binary.LittleEndian, m.ALT_ERROR)
	binary.Write(data, binary.LittleEndian, m.ASPD_ERROR)
	binary.Write(data, binary.LittleEndian, m.XTRACK_ERROR)
	binary.Write(data, binary.LittleEndian, m.NAV_BEARING)
	binary.Write(data, binary.LittleEndian, m.TARGET_BEARING)
	binary.Write(data, binary.LittleEndian, m.WP_DIST)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the NavControllerOutput
func (m *NavControllerOutput) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.NAV_ROLL)
	binary.Read(data, binary.LittleEndian, &m.NAV_PITCH)
	binary.Read(data, binary.LittleEndian, &m.ALT_ERROR)
	binary.Read(data, binary.LittleEndian, &m.ASPD_ERROR)
	binary.Read(data, binary.LittleEndian, &m.XTRACK_ERROR)
	binary.Read(data, binary.LittleEndian, &m.NAV_BEARING)
	binary.Read(data, binary.LittleEndian, &m.TARGET_BEARING)
	binary.Read(data, binary.LittleEndian, &m.WP_DIST)
}

//
// MESSAGE SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST 63
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_LEN 46
//
// MAVLINK_MSG_ID_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_CRC 130
//
//
type SetQuadSwarmLedRollPitchYawThrust struct {
	ROLL      [4]int16  // Desired roll angle in radians +-PI (+-INT16_MAX)
	PITCH     [4]int16  // Desired pitch angle in radians +-PI (+-INT16_MAX)
	YAW       [4]int16  // Desired yaw angle in radians, scaled to int16 +-PI (+-INT16_MAX)
	THRUST    [4]uint16 // Collective thrust, scaled to uint16 (0..UINT16_MAX)
	GROUP     uint8     // ID of the quadrotor group (0 - 255, up to 256 groups supported)
	MODE      uint8     // ID of the flight mode (0 - 255, up to 256 modes supported)
	LED_RED   [4]uint8  // RGB red channel (0-255)
	LED_BLUE  [4]uint8  // RGB green channel (0-255)
	LED_GREEN [4]uint8  // RGB blue channel (0-255)
}

// NewSetQuadSwarmLedRollPitchYawThrust returns a new SetQuadSwarmLedRollPitchYawThrust
func NewSetQuadSwarmLedRollPitchYawThrust(ROLL [4]int16, PITCH [4]int16, YAW [4]int16, THRUST [4]uint16, GROUP uint8, MODE uint8, LED_RED [4]uint8, LED_BLUE [4]uint8, LED_GREEN [4]uint8) *SetQuadSwarmLedRollPitchYawThrust {
	m := SetQuadSwarmLedRollPitchYawThrust{}
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.THRUST = THRUST
	m.GROUP = GROUP
	m.MODE = MODE
	m.LED_RED = LED_RED
	m.LED_BLUE = LED_BLUE
	m.LED_GREEN = LED_GREEN
	return &m
}

// Id returns the SetQuadSwarmLedRollPitchYawThrust Message ID
func (*SetQuadSwarmLedRollPitchYawThrust) Id() uint8 {
	return 63
}

// Len returns the SetQuadSwarmLedRollPitchYawThrust Message Length
func (*SetQuadSwarmLedRollPitchYawThrust) Len() uint8 {
	return 46
}

// Crc returns the SetQuadSwarmLedRollPitchYawThrust Message CRC
func (*SetQuadSwarmLedRollPitchYawThrust) Crc() uint8 {
	return 130
}

// Pack returns a packed byte array which represents a SetQuadSwarmLedRollPitchYawThrust payload
func (m *SetQuadSwarmLedRollPitchYawThrust) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.GROUP)
	binary.Write(data, binary.LittleEndian, m.MODE)
	binary.Write(data, binary.LittleEndian, m.LED_RED)
	binary.Write(data, binary.LittleEndian, m.LED_BLUE)
	binary.Write(data, binary.LittleEndian, m.LED_GREEN)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SetQuadSwarmLedRollPitchYawThrust
func (m *SetQuadSwarmLedRollPitchYawThrust) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.GROUP)
	binary.Read(data, binary.LittleEndian, &m.MODE)
	binary.Read(data, binary.LittleEndian, &m.LED_RED)
	binary.Read(data, binary.LittleEndian, &m.LED_BLUE)
	binary.Read(data, binary.LittleEndian, &m.LED_GREEN)
}

const (
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_roll_LEN      = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_pitch_LEN     = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_yaw_LEN       = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_thrust_LEN    = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_led_red_LEN   = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_led_blue_LEN  = 4
	MAVLINK_MSG_SET_QUAD_SWARM_LED_ROLL_PITCH_YAW_THRUST_FIELD_led_green_LEN = 4
)

//
// MESSAGE STATE_CORRECTION
//
// MAVLINK_MSG_ID_STATE_CORRECTION 64
//
// MAVLINK_MSG_ID_STATE_CORRECTION_LEN 36
//
// MAVLINK_MSG_ID_STATE_CORRECTION_CRC 130
//
//
type StateCorrection struct {
	XERR     float32 // x position error
	YERR     float32 // y position error
	ZERR     float32 // z position error
	ROLLERR  float32 // roll error (radians)
	PITCHERR float32 // pitch error (radians)
	YAWERR   float32 // yaw error (radians)
	VXERR    float32 // x velocity
	VYERR    float32 // y velocity
	VZERR    float32 // z velocity
}

// NewStateCorrection returns a new StateCorrection
func NewStateCorrection(XERR float32, YERR float32, ZERR float32, ROLLERR float32, PITCHERR float32, YAWERR float32, VXERR float32, VYERR float32, VZERR float32) *StateCorrection {
	m := StateCorrection{}
	m.XERR = XERR
	m.YERR = YERR
	m.ZERR = ZERR
	m.ROLLERR = ROLLERR
	m.PITCHERR = PITCHERR
	m.YAWERR = YAWERR
	m.VXERR = VXERR
	m.VYERR = VYERR
	m.VZERR = VZERR
	return &m
}

// Id returns the StateCorrection Message ID
func (*StateCorrection) Id() uint8 {
	return 64
}

// Len returns the StateCorrection Message Length
func (*StateCorrection) Len() uint8 {
	return 36
}

// Crc returns the StateCorrection Message CRC
func (*StateCorrection) Crc() uint8 {
	return 130
}

// Pack returns a packed byte array which represents a StateCorrection payload
func (m *StateCorrection) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.XERR)
	binary.Write(data, binary.LittleEndian, m.YERR)
	binary.Write(data, binary.LittleEndian, m.ZERR)
	binary.Write(data, binary.LittleEndian, m.ROLLERR)
	binary.Write(data, binary.LittleEndian, m.PITCHERR)
	binary.Write(data, binary.LittleEndian, m.YAWERR)
	binary.Write(data, binary.LittleEndian, m.VXERR)
	binary.Write(data, binary.LittleEndian, m.VYERR)
	binary.Write(data, binary.LittleEndian, m.VZERR)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the StateCorrection
func (m *StateCorrection) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.XERR)
	binary.Read(data, binary.LittleEndian, &m.YERR)
	binary.Read(data, binary.LittleEndian, &m.ZERR)
	binary.Read(data, binary.LittleEndian, &m.ROLLERR)
	binary.Read(data, binary.LittleEndian, &m.PITCHERR)
	binary.Read(data, binary.LittleEndian, &m.YAWERR)
	binary.Read(data, binary.LittleEndian, &m.VXERR)
	binary.Read(data, binary.LittleEndian, &m.VYERR)
	binary.Read(data, binary.LittleEndian, &m.VZERR)
}

//
// MESSAGE RC_CHANNELS
//
// MAVLINK_MSG_ID_RC_CHANNELS 65
//
// MAVLINK_MSG_ID_RC_CHANNELS_LEN 42
//
// MAVLINK_MSG_ID_RC_CHANNELS_CRC 118
//
//
type RcChannels struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	CHAN1_RAW    uint16 // RC channel 1 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN2_RAW    uint16 // RC channel 2 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN3_RAW    uint16 // RC channel 3 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN4_RAW    uint16 // RC channel 4 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN5_RAW    uint16 // RC channel 5 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN6_RAW    uint16 // RC channel 6 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN7_RAW    uint16 // RC channel 7 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN8_RAW    uint16 // RC channel 8 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN9_RAW    uint16 // RC channel 9 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN10_RAW   uint16 // RC channel 10 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN11_RAW   uint16 // RC channel 11 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN12_RAW   uint16 // RC channel 12 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN13_RAW   uint16 // RC channel 13 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN14_RAW   uint16 // RC channel 14 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN15_RAW   uint16 // RC channel 15 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN16_RAW   uint16 // RC channel 16 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN17_RAW   uint16 // RC channel 17 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHAN18_RAW   uint16 // RC channel 18 value, in microseconds. A value of UINT16_MAX implies the channel is unused.
	CHANCOUNT    uint8  // Total number of RC channels being received. This can be larger than 18, indicating that more channels are available but not given in this message. This value should be 0 when no RC channels are available.
	RSSI         uint8  // Receive signal strength indicator, 0: 0%, 100: 100%, 255: invalid/unknown.
}

// NewRcChannels returns a new RcChannels
func NewRcChannels(TIME_BOOT_MS uint32, CHAN1_RAW uint16, CHAN2_RAW uint16, CHAN3_RAW uint16, CHAN4_RAW uint16, CHAN5_RAW uint16, CHAN6_RAW uint16, CHAN7_RAW uint16, CHAN8_RAW uint16, CHAN9_RAW uint16, CHAN10_RAW uint16, CHAN11_RAW uint16, CHAN12_RAW uint16, CHAN13_RAW uint16, CHAN14_RAW uint16, CHAN15_RAW uint16, CHAN16_RAW uint16, CHAN17_RAW uint16, CHAN18_RAW uint16, CHANCOUNT uint8, RSSI uint8) *RcChannels {
	m := RcChannels{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.CHAN1_RAW = CHAN1_RAW
	m.CHAN2_RAW = CHAN2_RAW
	m.CHAN3_RAW = CHAN3_RAW
	m.CHAN4_RAW = CHAN4_RAW
	m.CHAN5_RAW = CHAN5_RAW
	m.CHAN6_RAW = CHAN6_RAW
	m.CHAN7_RAW = CHAN7_RAW
	m.CHAN8_RAW = CHAN8_RAW
	m.CHAN9_RAW = CHAN9_RAW
	m.CHAN10_RAW = CHAN10_RAW
	m.CHAN11_RAW = CHAN11_RAW
	m.CHAN12_RAW = CHAN12_RAW
	m.CHAN13_RAW = CHAN13_RAW
	m.CHAN14_RAW = CHAN14_RAW
	m.CHAN15_RAW = CHAN15_RAW
	m.CHAN16_RAW = CHAN16_RAW
	m.CHAN17_RAW = CHAN17_RAW
	m.CHAN18_RAW = CHAN18_RAW
	m.CHANCOUNT = CHANCOUNT
	m.RSSI = RSSI
	return &m
}

// Id returns the RcChannels Message ID
func (*RcChannels) Id() uint8 {
	return 65
}

// Len returns the RcChannels Message Length
func (*RcChannels) Len() uint8 {
	return 42
}

// Crc returns the RcChannels Message CRC
func (*RcChannels) Crc() uint8 {
	return 118
}

// Pack returns a packed byte array which represents a RcChannels payload
func (m *RcChannels) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.CHAN1_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN2_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN3_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN4_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN5_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN6_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN7_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN8_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN9_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN10_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN11_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN12_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN13_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN14_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN15_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN16_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN17_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN18_RAW)
	binary.Write(data, binary.LittleEndian, m.CHANCOUNT)
	binary.Write(data, binary.LittleEndian, m.RSSI)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RcChannels
func (m *RcChannels) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.CHAN1_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN2_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN3_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN4_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN5_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN6_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN7_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN8_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN9_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN10_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN11_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN12_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN13_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN14_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN15_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN16_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN17_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN18_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHANCOUNT)
	binary.Read(data, binary.LittleEndian, &m.RSSI)
}

//
// MESSAGE REQUEST_DATA_STREAM
//
// MAVLINK_MSG_ID_REQUEST_DATA_STREAM 66
//
// MAVLINK_MSG_ID_REQUEST_DATA_STREAM_LEN 6
//
// MAVLINK_MSG_ID_REQUEST_DATA_STREAM_CRC 148
//
//
type RequestDataStream struct {
	REQ_MESSAGE_RATE uint16 // The requested interval between two messages of this type
	TARGET_SYSTEM    uint8  // The target requested to send the message stream.
	TARGET_COMPONENT uint8  // The target requested to send the message stream.
	REQ_STREAM_ID    uint8  // The ID of the requested data stream
	START_STOP       uint8  // 1 to start sending, 0 to stop sending.
}

// NewRequestDataStream returns a new RequestDataStream
func NewRequestDataStream(REQ_MESSAGE_RATE uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, REQ_STREAM_ID uint8, START_STOP uint8) *RequestDataStream {
	m := RequestDataStream{}
	m.REQ_MESSAGE_RATE = REQ_MESSAGE_RATE
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.REQ_STREAM_ID = REQ_STREAM_ID
	m.START_STOP = START_STOP
	return &m
}

// Id returns the RequestDataStream Message ID
func (*RequestDataStream) Id() uint8 {
	return 66
}

// Len returns the RequestDataStream Message Length
func (*RequestDataStream) Len() uint8 {
	return 6
}

// Crc returns the RequestDataStream Message CRC
func (*RequestDataStream) Crc() uint8 {
	return 148
}

// Pack returns a packed byte array which represents a RequestDataStream payload
func (m *RequestDataStream) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.REQ_MESSAGE_RATE)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.REQ_STREAM_ID)
	binary.Write(data, binary.LittleEndian, m.START_STOP)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RequestDataStream
func (m *RequestDataStream) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.REQ_MESSAGE_RATE)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.REQ_STREAM_ID)
	binary.Read(data, binary.LittleEndian, &m.START_STOP)
}

//
// MESSAGE DATA_STREAM
//
// MAVLINK_MSG_ID_DATA_STREAM 67
//
// MAVLINK_MSG_ID_DATA_STREAM_LEN 4
//
// MAVLINK_MSG_ID_DATA_STREAM_CRC 21
//
//
type DataStream struct {
	MESSAGE_RATE uint16 // The requested interval between two messages of this type
	STREAM_ID    uint8  // The ID of the requested data stream
	ON_OFF       uint8  // 1 stream is enabled, 0 stream is stopped.
}

// NewDataStream returns a new DataStream
func NewDataStream(MESSAGE_RATE uint16, STREAM_ID uint8, ON_OFF uint8) *DataStream {
	m := DataStream{}
	m.MESSAGE_RATE = MESSAGE_RATE
	m.STREAM_ID = STREAM_ID
	m.ON_OFF = ON_OFF
	return &m
}

// Id returns the DataStream Message ID
func (*DataStream) Id() uint8 {
	return 67
}

// Len returns the DataStream Message Length
func (*DataStream) Len() uint8 {
	return 4
}

// Crc returns the DataStream Message CRC
func (*DataStream) Crc() uint8 {
	return 21
}

// Pack returns a packed byte array which represents a DataStream payload
func (m *DataStream) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.MESSAGE_RATE)
	binary.Write(data, binary.LittleEndian, m.STREAM_ID)
	binary.Write(data, binary.LittleEndian, m.ON_OFF)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the DataStream
func (m *DataStream) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.MESSAGE_RATE)
	binary.Read(data, binary.LittleEndian, &m.STREAM_ID)
	binary.Read(data, binary.LittleEndian, &m.ON_OFF)
}

//
// MESSAGE MANUAL_CONTROL
//
// MAVLINK_MSG_ID_MANUAL_CONTROL 69
//
// MAVLINK_MSG_ID_MANUAL_CONTROL_LEN 11
//
// MAVLINK_MSG_ID_MANUAL_CONTROL_CRC 243
//
//
type ManualControl struct {
	X       int16  // X-axis, normalized to the range [-1000,1000]. A value of INT16_MAX indicates that this axis is invalid. Generally corresponds to forward(1000)-backward(-1000) movement on a joystick and the pitch of a vehicle.
	Y       int16  // Y-axis, normalized to the range [-1000,1000]. A value of INT16_MAX indicates that this axis is invalid. Generally corresponds to left(-1000)-right(1000) movement on a joystick and the roll of a vehicle.
	Z       int16  // Z-axis, normalized to the range [-1000,1000]. A value of INT16_MAX indicates that this axis is invalid. Generally corresponds to a separate slider movement with maximum being 1000 and minimum being -1000 on a joystick and the thrust of a vehicle.
	R       int16  // R-axis, normalized to the range [-1000,1000]. A value of INT16_MAX indicates that this axis is invalid. Generally corresponds to a twisting of the joystick, with counter-clockwise being 1000 and clockwise being -1000, and the yaw of a vehicle.
	BUTTONS uint16 // A bitfield corresponding to the joystick buttons' current state, 1 for pressed, 0 for released. The lowest bit corresponds to Button 1.
	TARGET  uint8  // The system to be controlled.
}

// NewManualControl returns a new ManualControl
func NewManualControl(X int16, Y int16, Z int16, R int16, BUTTONS uint16, TARGET uint8) *ManualControl {
	m := ManualControl{}
	m.X = X
	m.Y = Y
	m.Z = Z
	m.R = R
	m.BUTTONS = BUTTONS
	m.TARGET = TARGET
	return &m
}

// Id returns the ManualControl Message ID
func (*ManualControl) Id() uint8 {
	return 69
}

// Len returns the ManualControl Message Length
func (*ManualControl) Len() uint8 {
	return 11
}

// Crc returns the ManualControl Message CRC
func (*ManualControl) Crc() uint8 {
	return 243
}

// Pack returns a packed byte array which represents a ManualControl payload
func (m *ManualControl) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.R)
	binary.Write(data, binary.LittleEndian, m.BUTTONS)
	binary.Write(data, binary.LittleEndian, m.TARGET)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ManualControl
func (m *ManualControl) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.R)
	binary.Read(data, binary.LittleEndian, &m.BUTTONS)
	binary.Read(data, binary.LittleEndian, &m.TARGET)
}

//
// MESSAGE RC_CHANNELS_OVERRIDE
//
// MAVLINK_MSG_ID_RC_CHANNELS_OVERRIDE 70
//
// MAVLINK_MSG_ID_RC_CHANNELS_OVERRIDE_LEN 18
//
// MAVLINK_MSG_ID_RC_CHANNELS_OVERRIDE_CRC 124
//
//
type RcChannelsOverride struct {
	CHAN1_RAW        uint16 // RC channel 1 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN2_RAW        uint16 // RC channel 2 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN3_RAW        uint16 // RC channel 3 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN4_RAW        uint16 // RC channel 4 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN5_RAW        uint16 // RC channel 5 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN6_RAW        uint16 // RC channel 6 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN7_RAW        uint16 // RC channel 7 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	CHAN8_RAW        uint16 // RC channel 8 value, in microseconds. A value of UINT16_MAX means to ignore this field.
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewRcChannelsOverride returns a new RcChannelsOverride
func NewRcChannelsOverride(CHAN1_RAW uint16, CHAN2_RAW uint16, CHAN3_RAW uint16, CHAN4_RAW uint16, CHAN5_RAW uint16, CHAN6_RAW uint16, CHAN7_RAW uint16, CHAN8_RAW uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *RcChannelsOverride {
	m := RcChannelsOverride{}
	m.CHAN1_RAW = CHAN1_RAW
	m.CHAN2_RAW = CHAN2_RAW
	m.CHAN3_RAW = CHAN3_RAW
	m.CHAN4_RAW = CHAN4_RAW
	m.CHAN5_RAW = CHAN5_RAW
	m.CHAN6_RAW = CHAN6_RAW
	m.CHAN7_RAW = CHAN7_RAW
	m.CHAN8_RAW = CHAN8_RAW
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the RcChannelsOverride Message ID
func (*RcChannelsOverride) Id() uint8 {
	return 70
}

// Len returns the RcChannelsOverride Message Length
func (*RcChannelsOverride) Len() uint8 {
	return 18
}

// Crc returns the RcChannelsOverride Message CRC
func (*RcChannelsOverride) Crc() uint8 {
	return 124
}

// Pack returns a packed byte array which represents a RcChannelsOverride payload
func (m *RcChannelsOverride) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.CHAN1_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN2_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN3_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN4_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN5_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN6_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN7_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN8_RAW)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RcChannelsOverride
func (m *RcChannelsOverride) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.CHAN1_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN2_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN3_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN4_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN5_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN6_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN7_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN8_RAW)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE VFR_HUD
//
// MAVLINK_MSG_ID_VFR_HUD 74
//
// MAVLINK_MSG_ID_VFR_HUD_LEN 20
//
// MAVLINK_MSG_ID_VFR_HUD_CRC 20
//
//
type VfrHud struct {
	AIRSPEED    float32 // Current airspeed in m/s
	GROUNDSPEED float32 // Current ground speed in m/s
	ALT         float32 // Current altitude (MSL), in meters
	CLIMB       float32 // Current climb rate in meters/second
	HEADING     int16   // Current heading in degrees, in compass units (0..360, 0=north)
	THROTTLE    uint16  // Current throttle setting in integer percent, 0 to 100
}

// NewVfrHud returns a new VfrHud
func NewVfrHud(AIRSPEED float32, GROUNDSPEED float32, ALT float32, CLIMB float32, HEADING int16, THROTTLE uint16) *VfrHud {
	m := VfrHud{}
	m.AIRSPEED = AIRSPEED
	m.GROUNDSPEED = GROUNDSPEED
	m.ALT = ALT
	m.CLIMB = CLIMB
	m.HEADING = HEADING
	m.THROTTLE = THROTTLE
	return &m
}

// Id returns the VfrHud Message ID
func (*VfrHud) Id() uint8 {
	return 74
}

// Len returns the VfrHud Message Length
func (*VfrHud) Len() uint8 {
	return 20
}

// Crc returns the VfrHud Message CRC
func (*VfrHud) Crc() uint8 {
	return 20
}

// Pack returns a packed byte array which represents a VfrHud payload
func (m *VfrHud) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.AIRSPEED)
	binary.Write(data, binary.LittleEndian, m.GROUNDSPEED)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.CLIMB)
	binary.Write(data, binary.LittleEndian, m.HEADING)
	binary.Write(data, binary.LittleEndian, m.THROTTLE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the VfrHud
func (m *VfrHud) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.AIRSPEED)
	binary.Read(data, binary.LittleEndian, &m.GROUNDSPEED)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.CLIMB)
	binary.Read(data, binary.LittleEndian, &m.HEADING)
	binary.Read(data, binary.LittleEndian, &m.THROTTLE)
}

//
// MESSAGE COMMAND_LONG
//
// MAVLINK_MSG_ID_COMMAND_LONG 76
//
// MAVLINK_MSG_ID_COMMAND_LONG_LEN 33
//
// MAVLINK_MSG_ID_COMMAND_LONG_CRC 152
//
//
type CommandLong struct {
	PARAM1           float32 // Parameter 1, as defined by MAV_CMD enum.
	PARAM2           float32 // Parameter 2, as defined by MAV_CMD enum.
	PARAM3           float32 // Parameter 3, as defined by MAV_CMD enum.
	PARAM4           float32 // Parameter 4, as defined by MAV_CMD enum.
	PARAM5           float32 // Parameter 5, as defined by MAV_CMD enum.
	PARAM6           float32 // Parameter 6, as defined by MAV_CMD enum.
	PARAM7           float32 // Parameter 7, as defined by MAV_CMD enum.
	COMMAND          uint16  // Command ID, as defined by MAV_CMD enum.
	TARGET_SYSTEM    uint8   // System which should execute the command
	TARGET_COMPONENT uint8   // Component which should execute the command, 0 for all components
	CONFIRMATION     uint8   // 0: First transmission of this command. 1-255: Confirmation transmissions (e.g. for kill command)
}

// NewCommandLong returns a new CommandLong
func NewCommandLong(PARAM1 float32, PARAM2 float32, PARAM3 float32, PARAM4 float32, PARAM5 float32, PARAM6 float32, PARAM7 float32, COMMAND uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, CONFIRMATION uint8) *CommandLong {
	m := CommandLong{}
	m.PARAM1 = PARAM1
	m.PARAM2 = PARAM2
	m.PARAM3 = PARAM3
	m.PARAM4 = PARAM4
	m.PARAM5 = PARAM5
	m.PARAM6 = PARAM6
	m.PARAM7 = PARAM7
	m.COMMAND = COMMAND
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.CONFIRMATION = CONFIRMATION
	return &m
}

// Id returns the CommandLong Message ID
func (*CommandLong) Id() uint8 {
	return 76
}

// Len returns the CommandLong Message Length
func (*CommandLong) Len() uint8 {
	return 33
}

// Crc returns the CommandLong Message CRC
func (*CommandLong) Crc() uint8 {
	return 152
}

// Pack returns a packed byte array which represents a CommandLong payload
func (m *CommandLong) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.PARAM1)
	binary.Write(data, binary.LittleEndian, m.PARAM2)
	binary.Write(data, binary.LittleEndian, m.PARAM3)
	binary.Write(data, binary.LittleEndian, m.PARAM4)
	binary.Write(data, binary.LittleEndian, m.PARAM5)
	binary.Write(data, binary.LittleEndian, m.PARAM6)
	binary.Write(data, binary.LittleEndian, m.PARAM7)
	binary.Write(data, binary.LittleEndian, m.COMMAND)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.CONFIRMATION)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the CommandLong
func (m *CommandLong) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.PARAM1)
	binary.Read(data, binary.LittleEndian, &m.PARAM2)
	binary.Read(data, binary.LittleEndian, &m.PARAM3)
	binary.Read(data, binary.LittleEndian, &m.PARAM4)
	binary.Read(data, binary.LittleEndian, &m.PARAM5)
	binary.Read(data, binary.LittleEndian, &m.PARAM6)
	binary.Read(data, binary.LittleEndian, &m.PARAM7)
	binary.Read(data, binary.LittleEndian, &m.COMMAND)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.CONFIRMATION)
}

//
// MESSAGE COMMAND_ACK
//
// MAVLINK_MSG_ID_COMMAND_ACK 77
//
// MAVLINK_MSG_ID_COMMAND_ACK_LEN 3
//
// MAVLINK_MSG_ID_COMMAND_ACK_CRC 143
//
//
type CommandAck struct {
	COMMAND uint16 // Command ID, as defined by MAV_CMD enum.
	RESULT  uint8  // See MAV_RESULT enum
}

// NewCommandAck returns a new CommandAck
func NewCommandAck(COMMAND uint16, RESULT uint8) *CommandAck {
	m := CommandAck{}
	m.COMMAND = COMMAND
	m.RESULT = RESULT
	return &m
}

// Id returns the CommandAck Message ID
func (*CommandAck) Id() uint8 {
	return 77
}

// Len returns the CommandAck Message Length
func (*CommandAck) Len() uint8 {
	return 3
}

// Crc returns the CommandAck Message CRC
func (*CommandAck) Crc() uint8 {
	return 143
}

// Pack returns a packed byte array which represents a CommandAck payload
func (m *CommandAck) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.COMMAND)
	binary.Write(data, binary.LittleEndian, m.RESULT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the CommandAck
func (m *CommandAck) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.COMMAND)
	binary.Read(data, binary.LittleEndian, &m.RESULT)
}

//
// MESSAGE ROLL_PITCH_YAW_RATES_THRUST_SETPOINT
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_RATES_THRUST_SETPOINT 80
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_RATES_THRUST_SETPOINT_LEN 20
//
// MAVLINK_MSG_ID_ROLL_PITCH_YAW_RATES_THRUST_SETPOINT_CRC 127
//
//
type RollPitchYawRatesThrustSetpoint struct {
	TIME_BOOT_MS uint32  // Timestamp in milliseconds since system boot
	ROLL_RATE    float32 // Desired roll rate in radians per second
	PITCH_RATE   float32 // Desired pitch rate in radians per second
	YAW_RATE     float32 // Desired yaw rate in radians per second
	THRUST       float32 // Collective thrust, normalized to 0 .. 1
}

// NewRollPitchYawRatesThrustSetpoint returns a new RollPitchYawRatesThrustSetpoint
func NewRollPitchYawRatesThrustSetpoint(TIME_BOOT_MS uint32, ROLL_RATE float32, PITCH_RATE float32, YAW_RATE float32, THRUST float32) *RollPitchYawRatesThrustSetpoint {
	m := RollPitchYawRatesThrustSetpoint{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.ROLL_RATE = ROLL_RATE
	m.PITCH_RATE = PITCH_RATE
	m.YAW_RATE = YAW_RATE
	m.THRUST = THRUST
	return &m
}

// Id returns the RollPitchYawRatesThrustSetpoint Message ID
func (*RollPitchYawRatesThrustSetpoint) Id() uint8 {
	return 80
}

// Len returns the RollPitchYawRatesThrustSetpoint Message Length
func (*RollPitchYawRatesThrustSetpoint) Len() uint8 {
	return 20
}

// Crc returns the RollPitchYawRatesThrustSetpoint Message CRC
func (*RollPitchYawRatesThrustSetpoint) Crc() uint8 {
	return 127
}

// Pack returns a packed byte array which represents a RollPitchYawRatesThrustSetpoint payload
func (m *RollPitchYawRatesThrustSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.ROLL_RATE)
	binary.Write(data, binary.LittleEndian, m.PITCH_RATE)
	binary.Write(data, binary.LittleEndian, m.YAW_RATE)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RollPitchYawRatesThrustSetpoint
func (m *RollPitchYawRatesThrustSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.ROLL_RATE)
	binary.Read(data, binary.LittleEndian, &m.PITCH_RATE)
	binary.Read(data, binary.LittleEndian, &m.YAW_RATE)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
}

//
// MESSAGE MANUAL_SETPOINT
//
// MAVLINK_MSG_ID_MANUAL_SETPOINT 81
//
// MAVLINK_MSG_ID_MANUAL_SETPOINT_LEN 22
//
// MAVLINK_MSG_ID_MANUAL_SETPOINT_CRC 106
//
//
type ManualSetpoint struct {
	TIME_BOOT_MS           uint32  // Timestamp in milliseconds since system boot
	ROLL                   float32 // Desired roll rate in radians per second
	PITCH                  float32 // Desired pitch rate in radians per second
	YAW                    float32 // Desired yaw rate in radians per second
	THRUST                 float32 // Collective thrust, normalized to 0 .. 1
	MODE_SWITCH            uint8   // Flight mode switch position, 0.. 255
	MANUAL_OVERRIDE_SWITCH uint8   // Override mode switch position, 0.. 255
}

// NewManualSetpoint returns a new ManualSetpoint
func NewManualSetpoint(TIME_BOOT_MS uint32, ROLL float32, PITCH float32, YAW float32, THRUST float32, MODE_SWITCH uint8, MANUAL_OVERRIDE_SWITCH uint8) *ManualSetpoint {
	m := ManualSetpoint{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.THRUST = THRUST
	m.MODE_SWITCH = MODE_SWITCH
	m.MANUAL_OVERRIDE_SWITCH = MANUAL_OVERRIDE_SWITCH
	return &m
}

// Id returns the ManualSetpoint Message ID
func (*ManualSetpoint) Id() uint8 {
	return 81
}

// Len returns the ManualSetpoint Message Length
func (*ManualSetpoint) Len() uint8 {
	return 22
}

// Crc returns the ManualSetpoint Message CRC
func (*ManualSetpoint) Crc() uint8 {
	return 106
}

// Pack returns a packed byte array which represents a ManualSetpoint payload
func (m *ManualSetpoint) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.MODE_SWITCH)
	binary.Write(data, binary.LittleEndian, m.MANUAL_OVERRIDE_SWITCH)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ManualSetpoint
func (m *ManualSetpoint) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.MODE_SWITCH)
	binary.Read(data, binary.LittleEndian, &m.MANUAL_OVERRIDE_SWITCH)
}

//
// MESSAGE ATTITUDE_SETPOINT_EXTERNAL
//
// MAVLINK_MSG_ID_ATTITUDE_SETPOINT_EXTERNAL 82
//
// MAVLINK_MSG_ID_ATTITUDE_SETPOINT_EXTERNAL_LEN 39
//
// MAVLINK_MSG_ID_ATTITUDE_SETPOINT_EXTERNAL_CRC 147
//
//
type AttitudeSetpointExternal struct {
	TIME_BOOT_MS     uint32     // Timestamp in milliseconds since system boot
	Q                [4]float32 // Attitude quaternion (w, x, y, z order, zero-rotation is 1, 0, 0, 0)
	BODY_ROLL_RATE   float32    // Body roll rate in radians per second
	BODY_PITCH_RATE  float32    // Body roll rate in radians per second
	BODY_YAW_RATE    float32    // Body roll rate in radians per second
	THRUST           float32    // Collective thrust, normalized to 0 .. 1 (-1 .. 1 for vehicles capable of reverse trust)
	TARGET_SYSTEM    uint8      // System ID
	TARGET_COMPONENT uint8      // Component ID
	TYPE_MASK        uint8      // Mappings: If any of these bits are set, the corresponding input should be ignored: bit 1: body roll rate, bit 2: body pitch rate, bit 3: body yaw rate. bit 4-bit 7: reserved, bit 8: attitude
}

// NewAttitudeSetpointExternal returns a new AttitudeSetpointExternal
func NewAttitudeSetpointExternal(TIME_BOOT_MS uint32, Q [4]float32, BODY_ROLL_RATE float32, BODY_PITCH_RATE float32, BODY_YAW_RATE float32, THRUST float32, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, TYPE_MASK uint8) *AttitudeSetpointExternal {
	m := AttitudeSetpointExternal{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.Q = Q
	m.BODY_ROLL_RATE = BODY_ROLL_RATE
	m.BODY_PITCH_RATE = BODY_PITCH_RATE
	m.BODY_YAW_RATE = BODY_YAW_RATE
	m.THRUST = THRUST
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.TYPE_MASK = TYPE_MASK
	return &m
}

// Id returns the AttitudeSetpointExternal Message ID
func (*AttitudeSetpointExternal) Id() uint8 {
	return 82
}

// Len returns the AttitudeSetpointExternal Message Length
func (*AttitudeSetpointExternal) Len() uint8 {
	return 39
}

// Crc returns the AttitudeSetpointExternal Message CRC
func (*AttitudeSetpointExternal) Crc() uint8 {
	return 147
}

// Pack returns a packed byte array which represents a AttitudeSetpointExternal payload
func (m *AttitudeSetpointExternal) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.Q)
	binary.Write(data, binary.LittleEndian, m.BODY_ROLL_RATE)
	binary.Write(data, binary.LittleEndian, m.BODY_PITCH_RATE)
	binary.Write(data, binary.LittleEndian, m.BODY_YAW_RATE)
	binary.Write(data, binary.LittleEndian, m.THRUST)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.TYPE_MASK)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the AttitudeSetpointExternal
func (m *AttitudeSetpointExternal) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.Q)
	binary.Read(data, binary.LittleEndian, &m.BODY_ROLL_RATE)
	binary.Read(data, binary.LittleEndian, &m.BODY_PITCH_RATE)
	binary.Read(data, binary.LittleEndian, &m.BODY_YAW_RATE)
	binary.Read(data, binary.LittleEndian, &m.THRUST)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.TYPE_MASK)
}

const (
	MAVLINK_MSG_ATTITUDE_SETPOINT_EXTERNAL_FIELD_q_LEN = 4
)

//
// MESSAGE LOCAL_NED_POSITION_SETPOINT_EXTERNAL
//
// MAVLINK_MSG_ID_LOCAL_NED_POSITION_SETPOINT_EXTERNAL 83
//
// MAVLINK_MSG_ID_LOCAL_NED_POSITION_SETPOINT_EXTERNAL_LEN 45
//
// MAVLINK_MSG_ID_LOCAL_NED_POSITION_SETPOINT_EXTERNAL_CRC 211
//
//
type LocalNedPositionSetpointExternal struct {
	TIME_BOOT_MS     uint32  // Timestamp in milliseconds since system boot
	X                float32 // X Position in NED frame in meters
	Y                float32 // Y Position in NED frame in meters
	Z                float32 // Z Position in NED frame in meters (note, altitude is negative in NED)
	VX               float32 // X velocity in NED frame in meter / s
	VY               float32 // Y velocity in NED frame in meter / s
	VZ               float32 // Z velocity in NED frame in meter / s
	AFX              float32 // X acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	AFY              float32 // Y acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	AFZ              float32 // Z acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	TYPE_MASK        uint16  // Bitmask to indicate which dimensions should be ignored by the vehicle: a value of 0b0000000000000000 or 0b0000001000000000 indicates that none of the setpoint dimensions should be ignored. If bit 10 is set the floats afx afy afz should be interpreted as force instead of acceleration. Mapping: bit 1: x, bit 2: y, bit 3: z, bit 4: vx, bit 5: vy, bit 6: vz, bit 7: ax, bit 8: ay, bit 9: az, bit 10: is force setpoint
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
	COORDINATE_FRAME uint8   // Valid options are: MAV_FRAME_LOCAL_NED, MAV_FRAME_LOCAL_OFFSET_NED = 5, MAV_FRAME_BODY_NED = 6, MAV_FRAME_BODY_OFFSET_NED = 7
}

// NewLocalNedPositionSetpointExternal returns a new LocalNedPositionSetpointExternal
func NewLocalNedPositionSetpointExternal(TIME_BOOT_MS uint32, X float32, Y float32, Z float32, VX float32, VY float32, VZ float32, AFX float32, AFY float32, AFZ float32, TYPE_MASK uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, COORDINATE_FRAME uint8) *LocalNedPositionSetpointExternal {
	m := LocalNedPositionSetpointExternal{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.X = X
	m.Y = Y
	m.Z = Z
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	m.AFX = AFX
	m.AFY = AFY
	m.AFZ = AFZ
	m.TYPE_MASK = TYPE_MASK
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.COORDINATE_FRAME = COORDINATE_FRAME
	return &m
}

// Id returns the LocalNedPositionSetpointExternal Message ID
func (*LocalNedPositionSetpointExternal) Id() uint8 {
	return 83
}

// Len returns the LocalNedPositionSetpointExternal Message Length
func (*LocalNedPositionSetpointExternal) Len() uint8 {
	return 45
}

// Crc returns the LocalNedPositionSetpointExternal Message CRC
func (*LocalNedPositionSetpointExternal) Crc() uint8 {
	return 211
}

// Pack returns a packed byte array which represents a LocalNedPositionSetpointExternal payload
func (m *LocalNedPositionSetpointExternal) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	binary.Write(data, binary.LittleEndian, m.AFX)
	binary.Write(data, binary.LittleEndian, m.AFY)
	binary.Write(data, binary.LittleEndian, m.AFZ)
	binary.Write(data, binary.LittleEndian, m.TYPE_MASK)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.COORDINATE_FRAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LocalNedPositionSetpointExternal
func (m *LocalNedPositionSetpointExternal) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
	binary.Read(data, binary.LittleEndian, &m.AFX)
	binary.Read(data, binary.LittleEndian, &m.AFY)
	binary.Read(data, binary.LittleEndian, &m.AFZ)
	binary.Read(data, binary.LittleEndian, &m.TYPE_MASK)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.COORDINATE_FRAME)
}

//
// MESSAGE GLOBAL_POSITION_SETPOINT_EXTERNAL_INT
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_EXTERNAL_INT 84
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_EXTERNAL_INT_LEN 44
//
// MAVLINK_MSG_ID_GLOBAL_POSITION_SETPOINT_EXTERNAL_INT_CRC 198
//
//
type GlobalPositionSetpointExternalInt struct {
	TIME_BOOT_MS     uint32  // Timestamp in milliseconds since system boot. The rationale for the timestamp in the setpoint is to allow the system to compensate for the transport delay of the setpoint. This allows the system to compensate processing latency.
	LAT_INT          int32   // X Position in WGS84 frame in 1e7 * meters
	LON_INT          int32   // Y Position in WGS84 frame in 1e7 * meters
	ALT              float32 // Altitude in WGS84, not AMSL
	VX               float32 // X velocity in NED frame in meter / s
	VY               float32 // Y velocity in NED frame in meter / s
	VZ               float32 // Z velocity in NED frame in meter / s
	AFX              float32 // X acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	AFY              float32 // Y acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	AFZ              float32 // Z acceleration or force (if bit 10 of type_mask is set) in NED frame in meter / s^2 or N
	TYPE_MASK        uint16  // Bitmask to indicate which dimensions should be ignored by the vehicle: a value of 0b0000000000000000 or 0b0000001000000000 indicates that none of the setpoint dimensions should be ignored. If bit 10 is set the floats afx afy afz should be interpreted as force instead of acceleration. Mapping: bit 1: x, bit 2: y, bit 3: z, bit 4: vx, bit 5: vy, bit 6: vz, bit 7: ax, bit 8: ay, bit 9: az, bit 10: is force setpoint
	TARGET_SYSTEM    uint8   // System ID
	TARGET_COMPONENT uint8   // Component ID
}

// NewGlobalPositionSetpointExternalInt returns a new GlobalPositionSetpointExternalInt
func NewGlobalPositionSetpointExternalInt(TIME_BOOT_MS uint32, LAT_INT int32, LON_INT int32, ALT float32, VX float32, VY float32, VZ float32, AFX float32, AFY float32, AFZ float32, TYPE_MASK uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *GlobalPositionSetpointExternalInt {
	m := GlobalPositionSetpointExternalInt{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.LAT_INT = LAT_INT
	m.LON_INT = LON_INT
	m.ALT = ALT
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	m.AFX = AFX
	m.AFY = AFY
	m.AFZ = AFZ
	m.TYPE_MASK = TYPE_MASK
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the GlobalPositionSetpointExternalInt Message ID
func (*GlobalPositionSetpointExternalInt) Id() uint8 {
	return 84
}

// Len returns the GlobalPositionSetpointExternalInt Message Length
func (*GlobalPositionSetpointExternalInt) Len() uint8 {
	return 44
}

// Crc returns the GlobalPositionSetpointExternalInt Message CRC
func (*GlobalPositionSetpointExternalInt) Crc() uint8 {
	return 198
}

// Pack returns a packed byte array which represents a GlobalPositionSetpointExternalInt payload
func (m *GlobalPositionSetpointExternalInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.LAT_INT)
	binary.Write(data, binary.LittleEndian, m.LON_INT)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	binary.Write(data, binary.LittleEndian, m.AFX)
	binary.Write(data, binary.LittleEndian, m.AFY)
	binary.Write(data, binary.LittleEndian, m.AFZ)
	binary.Write(data, binary.LittleEndian, m.TYPE_MASK)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GlobalPositionSetpointExternalInt
func (m *GlobalPositionSetpointExternalInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.LAT_INT)
	binary.Read(data, binary.LittleEndian, &m.LON_INT)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
	binary.Read(data, binary.LittleEndian, &m.AFX)
	binary.Read(data, binary.LittleEndian, &m.AFY)
	binary.Read(data, binary.LittleEndian, &m.AFZ)
	binary.Read(data, binary.LittleEndian, &m.TYPE_MASK)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET 89
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET_LEN 28
//
// MAVLINK_MSG_ID_LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET_CRC 231
//
//
type LocalPositionNedSystemGlobalOffset struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	X            float32 // X Position
	Y            float32 // Y Position
	Z            float32 // Z Position
	ROLL         float32 // Roll
	PITCH        float32 // Pitch
	YAW          float32 // Yaw
}

// NewLocalPositionNedSystemGlobalOffset returns a new LocalPositionNedSystemGlobalOffset
func NewLocalPositionNedSystemGlobalOffset(TIME_BOOT_MS uint32, X float32, Y float32, Z float32, ROLL float32, PITCH float32, YAW float32) *LocalPositionNedSystemGlobalOffset {
	m := LocalPositionNedSystemGlobalOffset{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.X = X
	m.Y = Y
	m.Z = Z
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	return &m
}

// Id returns the LocalPositionNedSystemGlobalOffset Message ID
func (*LocalPositionNedSystemGlobalOffset) Id() uint8 {
	return 89
}

// Len returns the LocalPositionNedSystemGlobalOffset Message Length
func (*LocalPositionNedSystemGlobalOffset) Len() uint8 {
	return 28
}

// Crc returns the LocalPositionNedSystemGlobalOffset Message CRC
func (*LocalPositionNedSystemGlobalOffset) Crc() uint8 {
	return 231
}

// Pack returns a packed byte array which represents a LocalPositionNedSystemGlobalOffset payload
func (m *LocalPositionNedSystemGlobalOffset) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LocalPositionNedSystemGlobalOffset
func (m *LocalPositionNedSystemGlobalOffset) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
}

//
// MESSAGE HIL_STATE
//
// MAVLINK_MSG_ID_HIL_STATE 90
//
// MAVLINK_MSG_ID_HIL_STATE_LEN 56
//
// MAVLINK_MSG_ID_HIL_STATE_CRC 183
//
//
type HilState struct {
	TIME_USEC  uint64  // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	ROLL       float32 // Roll angle (rad)
	PITCH      float32 // Pitch angle (rad)
	YAW        float32 // Yaw angle (rad)
	ROLLSPEED  float32 // Body frame roll / phi angular speed (rad/s)
	PITCHSPEED float32 // Body frame pitch / theta angular speed (rad/s)
	YAWSPEED   float32 // Body frame yaw / psi angular speed (rad/s)
	LAT        int32   // Latitude, expressed as * 1E7
	LON        int32   // Longitude, expressed as * 1E7
	ALT        int32   // Altitude in meters, expressed as * 1000 (millimeters)
	VX         int16   // Ground X Speed (Latitude), expressed as m/s * 100
	VY         int16   // Ground Y Speed (Longitude), expressed as m/s * 100
	VZ         int16   // Ground Z Speed (Altitude), expressed as m/s * 100
	XACC       int16   // X acceleration (mg)
	YACC       int16   // Y acceleration (mg)
	ZACC       int16   // Z acceleration (mg)
}

// NewHilState returns a new HilState
func NewHilState(TIME_USEC uint64, ROLL float32, PITCH float32, YAW float32, ROLLSPEED float32, PITCHSPEED float32, YAWSPEED float32, LAT int32, LON int32, ALT int32, VX int16, VY int16, VZ int16, XACC int16, YACC int16, ZACC int16) *HilState {
	m := HilState{}
	m.TIME_USEC = TIME_USEC
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.ROLLSPEED = ROLLSPEED
	m.PITCHSPEED = PITCHSPEED
	m.YAWSPEED = YAWSPEED
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	return &m
}

// Id returns the HilState Message ID
func (*HilState) Id() uint8 {
	return 90
}

// Len returns the HilState Message Length
func (*HilState) Len() uint8 {
	return 56
}

// Crc returns the HilState Message CRC
func (*HilState) Crc() uint8 {
	return 183
}

// Pack returns a packed byte array which represents a HilState payload
func (m *HilState) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.ROLLSPEED)
	binary.Write(data, binary.LittleEndian, m.PITCHSPEED)
	binary.Write(data, binary.LittleEndian, m.YAWSPEED)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilState
func (m *HilState) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.ROLLSPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCHSPEED)
	binary.Read(data, binary.LittleEndian, &m.YAWSPEED)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
}

//
// MESSAGE HIL_CONTROLS
//
// MAVLINK_MSG_ID_HIL_CONTROLS 91
//
// MAVLINK_MSG_ID_HIL_CONTROLS_LEN 42
//
// MAVLINK_MSG_ID_HIL_CONTROLS_CRC 63
//
//
type HilControls struct {
	TIME_USEC      uint64  // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	ROLL_AILERONS  float32 // Control output -1 .. 1
	PITCH_ELEVATOR float32 // Control output -1 .. 1
	YAW_RUDDER     float32 // Control output -1 .. 1
	THROTTLE       float32 // Throttle 0 .. 1
	AUX1           float32 // Aux 1, -1 .. 1
	AUX2           float32 // Aux 2, -1 .. 1
	AUX3           float32 // Aux 3, -1 .. 1
	AUX4           float32 // Aux 4, -1 .. 1
	MODE           uint8   // System mode (MAV_MODE)
	NAV_MODE       uint8   // Navigation mode (MAV_NAV_MODE)
}

// NewHilControls returns a new HilControls
func NewHilControls(TIME_USEC uint64, ROLL_AILERONS float32, PITCH_ELEVATOR float32, YAW_RUDDER float32, THROTTLE float32, AUX1 float32, AUX2 float32, AUX3 float32, AUX4 float32, MODE uint8, NAV_MODE uint8) *HilControls {
	m := HilControls{}
	m.TIME_USEC = TIME_USEC
	m.ROLL_AILERONS = ROLL_AILERONS
	m.PITCH_ELEVATOR = PITCH_ELEVATOR
	m.YAW_RUDDER = YAW_RUDDER
	m.THROTTLE = THROTTLE
	m.AUX1 = AUX1
	m.AUX2 = AUX2
	m.AUX3 = AUX3
	m.AUX4 = AUX4
	m.MODE = MODE
	m.NAV_MODE = NAV_MODE
	return &m
}

// Id returns the HilControls Message ID
func (*HilControls) Id() uint8 {
	return 91
}

// Len returns the HilControls Message Length
func (*HilControls) Len() uint8 {
	return 42
}

// Crc returns the HilControls Message CRC
func (*HilControls) Crc() uint8 {
	return 63
}

// Pack returns a packed byte array which represents a HilControls payload
func (m *HilControls) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.ROLL_AILERONS)
	binary.Write(data, binary.LittleEndian, m.PITCH_ELEVATOR)
	binary.Write(data, binary.LittleEndian, m.YAW_RUDDER)
	binary.Write(data, binary.LittleEndian, m.THROTTLE)
	binary.Write(data, binary.LittleEndian, m.AUX1)
	binary.Write(data, binary.LittleEndian, m.AUX2)
	binary.Write(data, binary.LittleEndian, m.AUX3)
	binary.Write(data, binary.LittleEndian, m.AUX4)
	binary.Write(data, binary.LittleEndian, m.MODE)
	binary.Write(data, binary.LittleEndian, m.NAV_MODE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilControls
func (m *HilControls) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.ROLL_AILERONS)
	binary.Read(data, binary.LittleEndian, &m.PITCH_ELEVATOR)
	binary.Read(data, binary.LittleEndian, &m.YAW_RUDDER)
	binary.Read(data, binary.LittleEndian, &m.THROTTLE)
	binary.Read(data, binary.LittleEndian, &m.AUX1)
	binary.Read(data, binary.LittleEndian, &m.AUX2)
	binary.Read(data, binary.LittleEndian, &m.AUX3)
	binary.Read(data, binary.LittleEndian, &m.AUX4)
	binary.Read(data, binary.LittleEndian, &m.MODE)
	binary.Read(data, binary.LittleEndian, &m.NAV_MODE)
}

//
// MESSAGE HIL_RC_INPUTS_RAW
//
// MAVLINK_MSG_ID_HIL_RC_INPUTS_RAW 92
//
// MAVLINK_MSG_ID_HIL_RC_INPUTS_RAW_LEN 33
//
// MAVLINK_MSG_ID_HIL_RC_INPUTS_RAW_CRC 54
//
//
type HilRcInputsRaw struct {
	TIME_USEC  uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	CHAN1_RAW  uint16 // RC channel 1 value, in microseconds
	CHAN2_RAW  uint16 // RC channel 2 value, in microseconds
	CHAN3_RAW  uint16 // RC channel 3 value, in microseconds
	CHAN4_RAW  uint16 // RC channel 4 value, in microseconds
	CHAN5_RAW  uint16 // RC channel 5 value, in microseconds
	CHAN6_RAW  uint16 // RC channel 6 value, in microseconds
	CHAN7_RAW  uint16 // RC channel 7 value, in microseconds
	CHAN8_RAW  uint16 // RC channel 8 value, in microseconds
	CHAN9_RAW  uint16 // RC channel 9 value, in microseconds
	CHAN10_RAW uint16 // RC channel 10 value, in microseconds
	CHAN11_RAW uint16 // RC channel 11 value, in microseconds
	CHAN12_RAW uint16 // RC channel 12 value, in microseconds
	RSSI       uint8  // Receive signal strength indicator, 0: 0%, 255: 100%
}

// NewHilRcInputsRaw returns a new HilRcInputsRaw
func NewHilRcInputsRaw(TIME_USEC uint64, CHAN1_RAW uint16, CHAN2_RAW uint16, CHAN3_RAW uint16, CHAN4_RAW uint16, CHAN5_RAW uint16, CHAN6_RAW uint16, CHAN7_RAW uint16, CHAN8_RAW uint16, CHAN9_RAW uint16, CHAN10_RAW uint16, CHAN11_RAW uint16, CHAN12_RAW uint16, RSSI uint8) *HilRcInputsRaw {
	m := HilRcInputsRaw{}
	m.TIME_USEC = TIME_USEC
	m.CHAN1_RAW = CHAN1_RAW
	m.CHAN2_RAW = CHAN2_RAW
	m.CHAN3_RAW = CHAN3_RAW
	m.CHAN4_RAW = CHAN4_RAW
	m.CHAN5_RAW = CHAN5_RAW
	m.CHAN6_RAW = CHAN6_RAW
	m.CHAN7_RAW = CHAN7_RAW
	m.CHAN8_RAW = CHAN8_RAW
	m.CHAN9_RAW = CHAN9_RAW
	m.CHAN10_RAW = CHAN10_RAW
	m.CHAN11_RAW = CHAN11_RAW
	m.CHAN12_RAW = CHAN12_RAW
	m.RSSI = RSSI
	return &m
}

// Id returns the HilRcInputsRaw Message ID
func (*HilRcInputsRaw) Id() uint8 {
	return 92
}

// Len returns the HilRcInputsRaw Message Length
func (*HilRcInputsRaw) Len() uint8 {
	return 33
}

// Crc returns the HilRcInputsRaw Message CRC
func (*HilRcInputsRaw) Crc() uint8 {
	return 54
}

// Pack returns a packed byte array which represents a HilRcInputsRaw payload
func (m *HilRcInputsRaw) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.CHAN1_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN2_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN3_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN4_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN5_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN6_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN7_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN8_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN9_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN10_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN11_RAW)
	binary.Write(data, binary.LittleEndian, m.CHAN12_RAW)
	binary.Write(data, binary.LittleEndian, m.RSSI)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilRcInputsRaw
func (m *HilRcInputsRaw) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.CHAN1_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN2_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN3_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN4_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN5_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN6_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN7_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN8_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN9_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN10_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN11_RAW)
	binary.Read(data, binary.LittleEndian, &m.CHAN12_RAW)
	binary.Read(data, binary.LittleEndian, &m.RSSI)
}

//
// MESSAGE OPTICAL_FLOW
//
// MAVLINK_MSG_ID_OPTICAL_FLOW 100
//
// MAVLINK_MSG_ID_OPTICAL_FLOW_LEN 26
//
// MAVLINK_MSG_ID_OPTICAL_FLOW_CRC 175
//
//
type OpticalFlow struct {
	TIME_USEC       uint64  // Timestamp (UNIX)
	FLOW_COMP_M_X   float32 // Flow in meters in x-sensor direction, angular-speed compensated
	FLOW_COMP_M_Y   float32 // Flow in meters in y-sensor direction, angular-speed compensated
	GROUND_DISTANCE float32 // Ground distance in meters. Positive value: distance known. Negative value: Unknown distance
	FLOW_X          int16   // Flow in pixels * 10 in x-sensor direction (dezi-pixels)
	FLOW_Y          int16   // Flow in pixels * 10 in y-sensor direction (dezi-pixels)
	SENSOR_ID       uint8   // Sensor ID
	QUALITY         uint8   // Optical flow quality / confidence. 0: bad, 255: maximum quality
}

// NewOpticalFlow returns a new OpticalFlow
func NewOpticalFlow(TIME_USEC uint64, FLOW_COMP_M_X float32, FLOW_COMP_M_Y float32, GROUND_DISTANCE float32, FLOW_X int16, FLOW_Y int16, SENSOR_ID uint8, QUALITY uint8) *OpticalFlow {
	m := OpticalFlow{}
	m.TIME_USEC = TIME_USEC
	m.FLOW_COMP_M_X = FLOW_COMP_M_X
	m.FLOW_COMP_M_Y = FLOW_COMP_M_Y
	m.GROUND_DISTANCE = GROUND_DISTANCE
	m.FLOW_X = FLOW_X
	m.FLOW_Y = FLOW_Y
	m.SENSOR_ID = SENSOR_ID
	m.QUALITY = QUALITY
	return &m
}

// Id returns the OpticalFlow Message ID
func (*OpticalFlow) Id() uint8 {
	return 100
}

// Len returns the OpticalFlow Message Length
func (*OpticalFlow) Len() uint8 {
	return 26
}

// Crc returns the OpticalFlow Message CRC
func (*OpticalFlow) Crc() uint8 {
	return 175
}

// Pack returns a packed byte array which represents a OpticalFlow payload
func (m *OpticalFlow) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.FLOW_COMP_M_X)
	binary.Write(data, binary.LittleEndian, m.FLOW_COMP_M_Y)
	binary.Write(data, binary.LittleEndian, m.GROUND_DISTANCE)
	binary.Write(data, binary.LittleEndian, m.FLOW_X)
	binary.Write(data, binary.LittleEndian, m.FLOW_Y)
	binary.Write(data, binary.LittleEndian, m.SENSOR_ID)
	binary.Write(data, binary.LittleEndian, m.QUALITY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the OpticalFlow
func (m *OpticalFlow) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.FLOW_COMP_M_X)
	binary.Read(data, binary.LittleEndian, &m.FLOW_COMP_M_Y)
	binary.Read(data, binary.LittleEndian, &m.GROUND_DISTANCE)
	binary.Read(data, binary.LittleEndian, &m.FLOW_X)
	binary.Read(data, binary.LittleEndian, &m.FLOW_Y)
	binary.Read(data, binary.LittleEndian, &m.SENSOR_ID)
	binary.Read(data, binary.LittleEndian, &m.QUALITY)
}

//
// MESSAGE GLOBAL_VISION_POSITION_ESTIMATE
//
// MAVLINK_MSG_ID_GLOBAL_VISION_POSITION_ESTIMATE 101
//
// MAVLINK_MSG_ID_GLOBAL_VISION_POSITION_ESTIMATE_LEN 32
//
// MAVLINK_MSG_ID_GLOBAL_VISION_POSITION_ESTIMATE_CRC 102
//
//
type GlobalVisionPositionEstimate struct {
	USEC  uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	X     float32 // Global X position
	Y     float32 // Global Y position
	Z     float32 // Global Z position
	ROLL  float32 // Roll angle in rad
	PITCH float32 // Pitch angle in rad
	YAW   float32 // Yaw angle in rad
}

// NewGlobalVisionPositionEstimate returns a new GlobalVisionPositionEstimate
func NewGlobalVisionPositionEstimate(USEC uint64, X float32, Y float32, Z float32, ROLL float32, PITCH float32, YAW float32) *GlobalVisionPositionEstimate {
	m := GlobalVisionPositionEstimate{}
	m.USEC = USEC
	m.X = X
	m.Y = Y
	m.Z = Z
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	return &m
}

// Id returns the GlobalVisionPositionEstimate Message ID
func (*GlobalVisionPositionEstimate) Id() uint8 {
	return 101
}

// Len returns the GlobalVisionPositionEstimate Message Length
func (*GlobalVisionPositionEstimate) Len() uint8 {
	return 32
}

// Crc returns the GlobalVisionPositionEstimate Message CRC
func (*GlobalVisionPositionEstimate) Crc() uint8 {
	return 102
}

// Pack returns a packed byte array which represents a GlobalVisionPositionEstimate payload
func (m *GlobalVisionPositionEstimate) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.USEC)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GlobalVisionPositionEstimate
func (m *GlobalVisionPositionEstimate) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.USEC)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
}

//
// MESSAGE VISION_POSITION_ESTIMATE
//
// MAVLINK_MSG_ID_VISION_POSITION_ESTIMATE 102
//
// MAVLINK_MSG_ID_VISION_POSITION_ESTIMATE_LEN 32
//
// MAVLINK_MSG_ID_VISION_POSITION_ESTIMATE_CRC 158
//
//
type VisionPositionEstimate struct {
	USEC  uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	X     float32 // Global X position
	Y     float32 // Global Y position
	Z     float32 // Global Z position
	ROLL  float32 // Roll angle in rad
	PITCH float32 // Pitch angle in rad
	YAW   float32 // Yaw angle in rad
}

// NewVisionPositionEstimate returns a new VisionPositionEstimate
func NewVisionPositionEstimate(USEC uint64, X float32, Y float32, Z float32, ROLL float32, PITCH float32, YAW float32) *VisionPositionEstimate {
	m := VisionPositionEstimate{}
	m.USEC = USEC
	m.X = X
	m.Y = Y
	m.Z = Z
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	return &m
}

// Id returns the VisionPositionEstimate Message ID
func (*VisionPositionEstimate) Id() uint8 {
	return 102
}

// Len returns the VisionPositionEstimate Message Length
func (*VisionPositionEstimate) Len() uint8 {
	return 32
}

// Crc returns the VisionPositionEstimate Message CRC
func (*VisionPositionEstimate) Crc() uint8 {
	return 158
}

// Pack returns a packed byte array which represents a VisionPositionEstimate payload
func (m *VisionPositionEstimate) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.USEC)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the VisionPositionEstimate
func (m *VisionPositionEstimate) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.USEC)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
}

//
// MESSAGE VISION_SPEED_ESTIMATE
//
// MAVLINK_MSG_ID_VISION_SPEED_ESTIMATE 103
//
// MAVLINK_MSG_ID_VISION_SPEED_ESTIMATE_LEN 20
//
// MAVLINK_MSG_ID_VISION_SPEED_ESTIMATE_CRC 208
//
//
type VisionSpeedEstimate struct {
	USEC uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	X    float32 // Global X speed
	Y    float32 // Global Y speed
	Z    float32 // Global Z speed
}

// NewVisionSpeedEstimate returns a new VisionSpeedEstimate
func NewVisionSpeedEstimate(USEC uint64, X float32, Y float32, Z float32) *VisionSpeedEstimate {
	m := VisionSpeedEstimate{}
	m.USEC = USEC
	m.X = X
	m.Y = Y
	m.Z = Z
	return &m
}

// Id returns the VisionSpeedEstimate Message ID
func (*VisionSpeedEstimate) Id() uint8 {
	return 103
}

// Len returns the VisionSpeedEstimate Message Length
func (*VisionSpeedEstimate) Len() uint8 {
	return 20
}

// Crc returns the VisionSpeedEstimate Message CRC
func (*VisionSpeedEstimate) Crc() uint8 {
	return 208
}

// Pack returns a packed byte array which represents a VisionSpeedEstimate payload
func (m *VisionSpeedEstimate) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.USEC)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the VisionSpeedEstimate
func (m *VisionSpeedEstimate) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.USEC)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
}

//
// MESSAGE VICON_POSITION_ESTIMATE
//
// MAVLINK_MSG_ID_VICON_POSITION_ESTIMATE 104
//
// MAVLINK_MSG_ID_VICON_POSITION_ESTIMATE_LEN 32
//
// MAVLINK_MSG_ID_VICON_POSITION_ESTIMATE_CRC 56
//
//
type ViconPositionEstimate struct {
	USEC  uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	X     float32 // Global X position
	Y     float32 // Global Y position
	Z     float32 // Global Z position
	ROLL  float32 // Roll angle in rad
	PITCH float32 // Pitch angle in rad
	YAW   float32 // Yaw angle in rad
}

// NewViconPositionEstimate returns a new ViconPositionEstimate
func NewViconPositionEstimate(USEC uint64, X float32, Y float32, Z float32, ROLL float32, PITCH float32, YAW float32) *ViconPositionEstimate {
	m := ViconPositionEstimate{}
	m.USEC = USEC
	m.X = X
	m.Y = Y
	m.Z = Z
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	return &m
}

// Id returns the ViconPositionEstimate Message ID
func (*ViconPositionEstimate) Id() uint8 {
	return 104
}

// Len returns the ViconPositionEstimate Message Length
func (*ViconPositionEstimate) Len() uint8 {
	return 32
}

// Crc returns the ViconPositionEstimate Message CRC
func (*ViconPositionEstimate) Crc() uint8 {
	return 56
}

// Pack returns a packed byte array which represents a ViconPositionEstimate payload
func (m *ViconPositionEstimate) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.USEC)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ViconPositionEstimate
func (m *ViconPositionEstimate) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.USEC)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
}

//
// MESSAGE HIGHRES_IMU
//
// MAVLINK_MSG_ID_HIGHRES_IMU 105
//
// MAVLINK_MSG_ID_HIGHRES_IMU_LEN 62
//
// MAVLINK_MSG_ID_HIGHRES_IMU_CRC 93
//
//
type HighresImu struct {
	TIME_USEC      uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	XACC           float32 // X acceleration (m/s^2)
	YACC           float32 // Y acceleration (m/s^2)
	ZACC           float32 // Z acceleration (m/s^2)
	XGYRO          float32 // Angular speed around X axis (rad / sec)
	YGYRO          float32 // Angular speed around Y axis (rad / sec)
	ZGYRO          float32 // Angular speed around Z axis (rad / sec)
	XMAG           float32 // X Magnetic field (Gauss)
	YMAG           float32 // Y Magnetic field (Gauss)
	ZMAG           float32 // Z Magnetic field (Gauss)
	ABS_PRESSURE   float32 // Absolute pressure in millibar
	DIFF_PRESSURE  float32 // Differential pressure in millibar
	PRESSURE_ALT   float32 // Altitude calculated from pressure
	TEMPERATURE    float32 // Temperature in degrees celsius
	FIELDS_UPDATED uint16  // Bitmask for fields that have updated since last message, bit 0 = xacc, bit 12: temperature
}

// NewHighresImu returns a new HighresImu
func NewHighresImu(TIME_USEC uint64, XACC float32, YACC float32, ZACC float32, XGYRO float32, YGYRO float32, ZGYRO float32, XMAG float32, YMAG float32, ZMAG float32, ABS_PRESSURE float32, DIFF_PRESSURE float32, PRESSURE_ALT float32, TEMPERATURE float32, FIELDS_UPDATED uint16) *HighresImu {
	m := HighresImu{}
	m.TIME_USEC = TIME_USEC
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.XMAG = XMAG
	m.YMAG = YMAG
	m.ZMAG = ZMAG
	m.ABS_PRESSURE = ABS_PRESSURE
	m.DIFF_PRESSURE = DIFF_PRESSURE
	m.PRESSURE_ALT = PRESSURE_ALT
	m.TEMPERATURE = TEMPERATURE
	m.FIELDS_UPDATED = FIELDS_UPDATED
	return &m
}

// Id returns the HighresImu Message ID
func (*HighresImu) Id() uint8 {
	return 105
}

// Len returns the HighresImu Message Length
func (*HighresImu) Len() uint8 {
	return 62
}

// Crc returns the HighresImu Message CRC
func (*HighresImu) Crc() uint8 {
	return 93
}

// Pack returns a packed byte array which represents a HighresImu payload
func (m *HighresImu) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.XMAG)
	binary.Write(data, binary.LittleEndian, m.YMAG)
	binary.Write(data, binary.LittleEndian, m.ZMAG)
	binary.Write(data, binary.LittleEndian, m.ABS_PRESSURE)
	binary.Write(data, binary.LittleEndian, m.DIFF_PRESSURE)
	binary.Write(data, binary.LittleEndian, m.PRESSURE_ALT)
	binary.Write(data, binary.LittleEndian, m.TEMPERATURE)
	binary.Write(data, binary.LittleEndian, m.FIELDS_UPDATED)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HighresImu
func (m *HighresImu) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.XMAG)
	binary.Read(data, binary.LittleEndian, &m.YMAG)
	binary.Read(data, binary.LittleEndian, &m.ZMAG)
	binary.Read(data, binary.LittleEndian, &m.ABS_PRESSURE)
	binary.Read(data, binary.LittleEndian, &m.DIFF_PRESSURE)
	binary.Read(data, binary.LittleEndian, &m.PRESSURE_ALT)
	binary.Read(data, binary.LittleEndian, &m.TEMPERATURE)
	binary.Read(data, binary.LittleEndian, &m.FIELDS_UPDATED)
}

//
// MESSAGE OMNIDIRECTIONAL_FLOW
//
// MAVLINK_MSG_ID_OMNIDIRECTIONAL_FLOW 106
//
// MAVLINK_MSG_ID_OMNIDIRECTIONAL_FLOW_LEN 54
//
// MAVLINK_MSG_ID_OMNIDIRECTIONAL_FLOW_CRC 211
//
//
type OmnidirectionalFlow struct {
	TIME_USEC        uint64    // Timestamp (microseconds, synced to UNIX time or since system boot)
	FRONT_DISTANCE_M float32   // Front distance in meters. Positive value (including zero): distance known. Negative value: Unknown distance
	LEFT             [10]int16 // Flow in deci pixels (1 = 0.1 pixel) on left hemisphere
	RIGHT            [10]int16 // Flow in deci pixels (1 = 0.1 pixel) on right hemisphere
	SENSOR_ID        uint8     // Sensor ID
	QUALITY          uint8     // Optical flow quality / confidence. 0: bad, 255: maximum quality
}

// NewOmnidirectionalFlow returns a new OmnidirectionalFlow
func NewOmnidirectionalFlow(TIME_USEC uint64, FRONT_DISTANCE_M float32, LEFT [10]int16, RIGHT [10]int16, SENSOR_ID uint8, QUALITY uint8) *OmnidirectionalFlow {
	m := OmnidirectionalFlow{}
	m.TIME_USEC = TIME_USEC
	m.FRONT_DISTANCE_M = FRONT_DISTANCE_M
	m.LEFT = LEFT
	m.RIGHT = RIGHT
	m.SENSOR_ID = SENSOR_ID
	m.QUALITY = QUALITY
	return &m
}

// Id returns the OmnidirectionalFlow Message ID
func (*OmnidirectionalFlow) Id() uint8 {
	return 106
}

// Len returns the OmnidirectionalFlow Message Length
func (*OmnidirectionalFlow) Len() uint8 {
	return 54
}

// Crc returns the OmnidirectionalFlow Message CRC
func (*OmnidirectionalFlow) Crc() uint8 {
	return 211
}

// Pack returns a packed byte array which represents a OmnidirectionalFlow payload
func (m *OmnidirectionalFlow) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.FRONT_DISTANCE_M)
	binary.Write(data, binary.LittleEndian, m.LEFT)
	binary.Write(data, binary.LittleEndian, m.RIGHT)
	binary.Write(data, binary.LittleEndian, m.SENSOR_ID)
	binary.Write(data, binary.LittleEndian, m.QUALITY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the OmnidirectionalFlow
func (m *OmnidirectionalFlow) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.FRONT_DISTANCE_M)
	binary.Read(data, binary.LittleEndian, &m.LEFT)
	binary.Read(data, binary.LittleEndian, &m.RIGHT)
	binary.Read(data, binary.LittleEndian, &m.SENSOR_ID)
	binary.Read(data, binary.LittleEndian, &m.QUALITY)
}

const (
	MAVLINK_MSG_OMNIDIRECTIONAL_FLOW_FIELD_left_LEN  = 10
	MAVLINK_MSG_OMNIDIRECTIONAL_FLOW_FIELD_right_LEN = 10
)

//
// MESSAGE HIL_SENSOR
//
// MAVLINK_MSG_ID_HIL_SENSOR 107
//
// MAVLINK_MSG_ID_HIL_SENSOR_LEN 64
//
// MAVLINK_MSG_ID_HIL_SENSOR_CRC 108
//
//
type HilSensor struct {
	TIME_USEC      uint64  // Timestamp (microseconds, synced to UNIX time or since system boot)
	XACC           float32 // X acceleration (m/s^2)
	YACC           float32 // Y acceleration (m/s^2)
	ZACC           float32 // Z acceleration (m/s^2)
	XGYRO          float32 // Angular speed around X axis in body frame (rad / sec)
	YGYRO          float32 // Angular speed around Y axis in body frame (rad / sec)
	ZGYRO          float32 // Angular speed around Z axis in body frame (rad / sec)
	XMAG           float32 // X Magnetic field (Gauss)
	YMAG           float32 // Y Magnetic field (Gauss)
	ZMAG           float32 // Z Magnetic field (Gauss)
	ABS_PRESSURE   float32 // Absolute pressure in millibar
	DIFF_PRESSURE  float32 // Differential pressure (airspeed) in millibar
	PRESSURE_ALT   float32 // Altitude calculated from pressure
	TEMPERATURE    float32 // Temperature in degrees celsius
	FIELDS_UPDATED uint32  // Bitmask for fields that have updated since last message, bit 0 = xacc, bit 12: temperature
}

// NewHilSensor returns a new HilSensor
func NewHilSensor(TIME_USEC uint64, XACC float32, YACC float32, ZACC float32, XGYRO float32, YGYRO float32, ZGYRO float32, XMAG float32, YMAG float32, ZMAG float32, ABS_PRESSURE float32, DIFF_PRESSURE float32, PRESSURE_ALT float32, TEMPERATURE float32, FIELDS_UPDATED uint32) *HilSensor {
	m := HilSensor{}
	m.TIME_USEC = TIME_USEC
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.XMAG = XMAG
	m.YMAG = YMAG
	m.ZMAG = ZMAG
	m.ABS_PRESSURE = ABS_PRESSURE
	m.DIFF_PRESSURE = DIFF_PRESSURE
	m.PRESSURE_ALT = PRESSURE_ALT
	m.TEMPERATURE = TEMPERATURE
	m.FIELDS_UPDATED = FIELDS_UPDATED
	return &m
}

// Id returns the HilSensor Message ID
func (*HilSensor) Id() uint8 {
	return 107
}

// Len returns the HilSensor Message Length
func (*HilSensor) Len() uint8 {
	return 64
}

// Crc returns the HilSensor Message CRC
func (*HilSensor) Crc() uint8 {
	return 108
}

// Pack returns a packed byte array which represents a HilSensor payload
func (m *HilSensor) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.XMAG)
	binary.Write(data, binary.LittleEndian, m.YMAG)
	binary.Write(data, binary.LittleEndian, m.ZMAG)
	binary.Write(data, binary.LittleEndian, m.ABS_PRESSURE)
	binary.Write(data, binary.LittleEndian, m.DIFF_PRESSURE)
	binary.Write(data, binary.LittleEndian, m.PRESSURE_ALT)
	binary.Write(data, binary.LittleEndian, m.TEMPERATURE)
	binary.Write(data, binary.LittleEndian, m.FIELDS_UPDATED)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilSensor
func (m *HilSensor) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.XMAG)
	binary.Read(data, binary.LittleEndian, &m.YMAG)
	binary.Read(data, binary.LittleEndian, &m.ZMAG)
	binary.Read(data, binary.LittleEndian, &m.ABS_PRESSURE)
	binary.Read(data, binary.LittleEndian, &m.DIFF_PRESSURE)
	binary.Read(data, binary.LittleEndian, &m.PRESSURE_ALT)
	binary.Read(data, binary.LittleEndian, &m.TEMPERATURE)
	binary.Read(data, binary.LittleEndian, &m.FIELDS_UPDATED)
}

//
// MESSAGE SIM_STATE
//
// MAVLINK_MSG_ID_SIM_STATE 108
//
// MAVLINK_MSG_ID_SIM_STATE_LEN 84
//
// MAVLINK_MSG_ID_SIM_STATE_CRC 32
//
//
type SimState struct {
	Q1           float32 // True attitude quaternion component 1, w (1 in null-rotation)
	Q2           float32 // True attitude quaternion component 2, x (0 in null-rotation)
	Q3           float32 // True attitude quaternion component 3, y (0 in null-rotation)
	Q4           float32 // True attitude quaternion component 4, z (0 in null-rotation)
	ROLL         float32 // Attitude roll expressed as Euler angles, not recommended except for human-readable outputs
	PITCH        float32 // Attitude pitch expressed as Euler angles, not recommended except for human-readable outputs
	YAW          float32 // Attitude yaw expressed as Euler angles, not recommended except for human-readable outputs
	XACC         float32 // X acceleration m/s/s
	YACC         float32 // Y acceleration m/s/s
	ZACC         float32 // Z acceleration m/s/s
	XGYRO        float32 // Angular speed around X axis rad/s
	YGYRO        float32 // Angular speed around Y axis rad/s
	ZGYRO        float32 // Angular speed around Z axis rad/s
	LAT          float32 // Latitude in degrees
	LON          float32 // Longitude in degrees
	ALT          float32 // Altitude in meters
	STD_DEV_HORZ float32 // Horizontal position standard deviation
	STD_DEV_VERT float32 // Vertical position standard deviation
	VN           float32 // True velocity in m/s in NORTH direction in earth-fixed NED frame
	VE           float32 // True velocity in m/s in EAST direction in earth-fixed NED frame
	VD           float32 // True velocity in m/s in DOWN direction in earth-fixed NED frame
}

// NewSimState returns a new SimState
func NewSimState(Q1 float32, Q2 float32, Q3 float32, Q4 float32, ROLL float32, PITCH float32, YAW float32, XACC float32, YACC float32, ZACC float32, XGYRO float32, YGYRO float32, ZGYRO float32, LAT float32, LON float32, ALT float32, STD_DEV_HORZ float32, STD_DEV_VERT float32, VN float32, VE float32, VD float32) *SimState {
	m := SimState{}
	m.Q1 = Q1
	m.Q2 = Q2
	m.Q3 = Q3
	m.Q4 = Q4
	m.ROLL = ROLL
	m.PITCH = PITCH
	m.YAW = YAW
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.STD_DEV_HORZ = STD_DEV_HORZ
	m.STD_DEV_VERT = STD_DEV_VERT
	m.VN = VN
	m.VE = VE
	m.VD = VD
	return &m
}

// Id returns the SimState Message ID
func (*SimState) Id() uint8 {
	return 108
}

// Len returns the SimState Message Length
func (*SimState) Len() uint8 {
	return 84
}

// Crc returns the SimState Message CRC
func (*SimState) Crc() uint8 {
	return 32
}

// Pack returns a packed byte array which represents a SimState payload
func (m *SimState) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.Q1)
	binary.Write(data, binary.LittleEndian, m.Q2)
	binary.Write(data, binary.LittleEndian, m.Q3)
	binary.Write(data, binary.LittleEndian, m.Q4)
	binary.Write(data, binary.LittleEndian, m.ROLL)
	binary.Write(data, binary.LittleEndian, m.PITCH)
	binary.Write(data, binary.LittleEndian, m.YAW)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.STD_DEV_HORZ)
	binary.Write(data, binary.LittleEndian, m.STD_DEV_VERT)
	binary.Write(data, binary.LittleEndian, m.VN)
	binary.Write(data, binary.LittleEndian, m.VE)
	binary.Write(data, binary.LittleEndian, m.VD)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SimState
func (m *SimState) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.Q1)
	binary.Read(data, binary.LittleEndian, &m.Q2)
	binary.Read(data, binary.LittleEndian, &m.Q3)
	binary.Read(data, binary.LittleEndian, &m.Q4)
	binary.Read(data, binary.LittleEndian, &m.ROLL)
	binary.Read(data, binary.LittleEndian, &m.PITCH)
	binary.Read(data, binary.LittleEndian, &m.YAW)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.STD_DEV_HORZ)
	binary.Read(data, binary.LittleEndian, &m.STD_DEV_VERT)
	binary.Read(data, binary.LittleEndian, &m.VN)
	binary.Read(data, binary.LittleEndian, &m.VE)
	binary.Read(data, binary.LittleEndian, &m.VD)
}

//
// MESSAGE RADIO_STATUS
//
// MAVLINK_MSG_ID_RADIO_STATUS 109
//
// MAVLINK_MSG_ID_RADIO_STATUS_LEN 9
//
// MAVLINK_MSG_ID_RADIO_STATUS_CRC 185
//
//
type RadioStatus struct {
	RXERRORS uint16 // receive errors
	FIXED    uint16 // count of error corrected packets
	RSSI     uint8  // local signal strength
	REMRSSI  uint8  // remote signal strength
	TXBUF    uint8  // how full the tx buffer is as a percentage
	NOISE    uint8  // background noise level
	REMNOISE uint8  // remote background noise level
}

// NewRadioStatus returns a new RadioStatus
func NewRadioStatus(RXERRORS uint16, FIXED uint16, RSSI uint8, REMRSSI uint8, TXBUF uint8, NOISE uint8, REMNOISE uint8) *RadioStatus {
	m := RadioStatus{}
	m.RXERRORS = RXERRORS
	m.FIXED = FIXED
	m.RSSI = RSSI
	m.REMRSSI = REMRSSI
	m.TXBUF = TXBUF
	m.NOISE = NOISE
	m.REMNOISE = REMNOISE
	return &m
}

// Id returns the RadioStatus Message ID
func (*RadioStatus) Id() uint8 {
	return 109
}

// Len returns the RadioStatus Message Length
func (*RadioStatus) Len() uint8 {
	return 9
}

// Crc returns the RadioStatus Message CRC
func (*RadioStatus) Crc() uint8 {
	return 185
}

// Pack returns a packed byte array which represents a RadioStatus payload
func (m *RadioStatus) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.RXERRORS)
	binary.Write(data, binary.LittleEndian, m.FIXED)
	binary.Write(data, binary.LittleEndian, m.RSSI)
	binary.Write(data, binary.LittleEndian, m.REMRSSI)
	binary.Write(data, binary.LittleEndian, m.TXBUF)
	binary.Write(data, binary.LittleEndian, m.NOISE)
	binary.Write(data, binary.LittleEndian, m.REMNOISE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the RadioStatus
func (m *RadioStatus) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.RXERRORS)
	binary.Read(data, binary.LittleEndian, &m.FIXED)
	binary.Read(data, binary.LittleEndian, &m.RSSI)
	binary.Read(data, binary.LittleEndian, &m.REMRSSI)
	binary.Read(data, binary.LittleEndian, &m.TXBUF)
	binary.Read(data, binary.LittleEndian, &m.NOISE)
	binary.Read(data, binary.LittleEndian, &m.REMNOISE)
}

//
// MESSAGE FILE_TRANSFER_START
//
// MAVLINK_MSG_ID_FILE_TRANSFER_START 110
//
// MAVLINK_MSG_ID_FILE_TRANSFER_START_LEN 254
//
// MAVLINK_MSG_ID_FILE_TRANSFER_START_CRC 235
//
//
type FileTransferStart struct {
	TRANSFER_UID uint64     // Unique transfer ID
	FILE_SIZE    uint32     // File size in bytes
	DEST_PATH    [240]uint8 // Destination path
	DIRECTION    uint8      // Transfer direction: 0: from requester, 1: to requester
	FLAGS        uint8      // RESERVED
}

// NewFileTransferStart returns a new FileTransferStart
func NewFileTransferStart(TRANSFER_UID uint64, FILE_SIZE uint32, DEST_PATH [240]uint8, DIRECTION uint8, FLAGS uint8) *FileTransferStart {
	m := FileTransferStart{}
	m.TRANSFER_UID = TRANSFER_UID
	m.FILE_SIZE = FILE_SIZE
	m.DEST_PATH = DEST_PATH
	m.DIRECTION = DIRECTION
	m.FLAGS = FLAGS
	return &m
}

// Id returns the FileTransferStart Message ID
func (*FileTransferStart) Id() uint8 {
	return 110
}

// Len returns the FileTransferStart Message Length
func (*FileTransferStart) Len() uint8 {
	return 254
}

// Crc returns the FileTransferStart Message CRC
func (*FileTransferStart) Crc() uint8 {
	return 235
}

// Pack returns a packed byte array which represents a FileTransferStart payload
func (m *FileTransferStart) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TRANSFER_UID)
	binary.Write(data, binary.LittleEndian, m.FILE_SIZE)
	binary.Write(data, binary.LittleEndian, m.DEST_PATH)
	binary.Write(data, binary.LittleEndian, m.DIRECTION)
	binary.Write(data, binary.LittleEndian, m.FLAGS)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the FileTransferStart
func (m *FileTransferStart) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TRANSFER_UID)
	binary.Read(data, binary.LittleEndian, &m.FILE_SIZE)
	binary.Read(data, binary.LittleEndian, &m.DEST_PATH)
	binary.Read(data, binary.LittleEndian, &m.DIRECTION)
	binary.Read(data, binary.LittleEndian, &m.FLAGS)
}

const (
	MAVLINK_MSG_FILE_TRANSFER_START_FIELD_dest_path_LEN = 240
)

//
// MESSAGE FILE_TRANSFER_DIR_LIST
//
// MAVLINK_MSG_ID_FILE_TRANSFER_DIR_LIST 111
//
// MAVLINK_MSG_ID_FILE_TRANSFER_DIR_LIST_LEN 249
//
// MAVLINK_MSG_ID_FILE_TRANSFER_DIR_LIST_CRC 93
//
//
type FileTransferDirList struct {
	TRANSFER_UID uint64     // Unique transfer ID
	DIR_PATH     [240]uint8 // Directory path to list
	FLAGS        uint8      // RESERVED
}

// NewFileTransferDirList returns a new FileTransferDirList
func NewFileTransferDirList(TRANSFER_UID uint64, DIR_PATH [240]uint8, FLAGS uint8) *FileTransferDirList {
	m := FileTransferDirList{}
	m.TRANSFER_UID = TRANSFER_UID
	m.DIR_PATH = DIR_PATH
	m.FLAGS = FLAGS
	return &m
}

// Id returns the FileTransferDirList Message ID
func (*FileTransferDirList) Id() uint8 {
	return 111
}

// Len returns the FileTransferDirList Message Length
func (*FileTransferDirList) Len() uint8 {
	return 249
}

// Crc returns the FileTransferDirList Message CRC
func (*FileTransferDirList) Crc() uint8 {
	return 93
}

// Pack returns a packed byte array which represents a FileTransferDirList payload
func (m *FileTransferDirList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TRANSFER_UID)
	binary.Write(data, binary.LittleEndian, m.DIR_PATH)
	binary.Write(data, binary.LittleEndian, m.FLAGS)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the FileTransferDirList
func (m *FileTransferDirList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TRANSFER_UID)
	binary.Read(data, binary.LittleEndian, &m.DIR_PATH)
	binary.Read(data, binary.LittleEndian, &m.FLAGS)
}

const (
	MAVLINK_MSG_FILE_TRANSFER_DIR_LIST_FIELD_dir_path_LEN = 240
)

//
// MESSAGE FILE_TRANSFER_RES
//
// MAVLINK_MSG_ID_FILE_TRANSFER_RES 112
//
// MAVLINK_MSG_ID_FILE_TRANSFER_RES_LEN 9
//
// MAVLINK_MSG_ID_FILE_TRANSFER_RES_CRC 124
//
//
type FileTransferRes struct {
	TRANSFER_UID uint64 // Unique transfer ID
	RESULT       uint8  // 0: OK, 1: not permitted, 2: bad path / file name, 3: no space left on device
}

// NewFileTransferRes returns a new FileTransferRes
func NewFileTransferRes(TRANSFER_UID uint64, RESULT uint8) *FileTransferRes {
	m := FileTransferRes{}
	m.TRANSFER_UID = TRANSFER_UID
	m.RESULT = RESULT
	return &m
}

// Id returns the FileTransferRes Message ID
func (*FileTransferRes) Id() uint8 {
	return 112
}

// Len returns the FileTransferRes Message Length
func (*FileTransferRes) Len() uint8 {
	return 9
}

// Crc returns the FileTransferRes Message CRC
func (*FileTransferRes) Crc() uint8 {
	return 124
}

// Pack returns a packed byte array which represents a FileTransferRes payload
func (m *FileTransferRes) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TRANSFER_UID)
	binary.Write(data, binary.LittleEndian, m.RESULT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the FileTransferRes
func (m *FileTransferRes) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TRANSFER_UID)
	binary.Read(data, binary.LittleEndian, &m.RESULT)
}

//
// MESSAGE HIL_GPS
//
// MAVLINK_MSG_ID_HIL_GPS 113
//
// MAVLINK_MSG_ID_HIL_GPS_LEN 36
//
// MAVLINK_MSG_ID_HIL_GPS_CRC 124
//
//
type HilGps struct {
	TIME_USEC          uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	LAT                int32  // Latitude (WGS84), in degrees * 1E7
	LON                int32  // Longitude (WGS84), in degrees * 1E7
	ALT                int32  // Altitude (WGS84), in meters * 1000 (positive for up)
	EPH                uint16 // GPS HDOP horizontal dilution of position in cm (m*100). If unknown, set to: 65535
	EPV                uint16 // GPS VDOP vertical dilution of position in cm (m*100). If unknown, set to: 65535
	VEL                uint16 // GPS ground speed (m/s * 100). If unknown, set to: 65535
	VN                 int16  // GPS velocity in cm/s in NORTH direction in earth-fixed NED frame
	VE                 int16  // GPS velocity in cm/s in EAST direction in earth-fixed NED frame
	VD                 int16  // GPS velocity in cm/s in DOWN direction in earth-fixed NED frame
	COG                uint16 // Course over ground (NOT heading, but direction of movement) in degrees * 100, 0.0..359.99 degrees. If unknown, set to: 65535
	FIX_TYPE           uint8  // 0-1: no fix, 2: 2D fix, 3: 3D fix. Some applications will not use the value of this field unless it is at least two, so always correctly fill in the fix.
	SATELLITES_VISIBLE uint8  // Number of satellites visible. If unknown, set to 255
}

// NewHilGps returns a new HilGps
func NewHilGps(TIME_USEC uint64, LAT int32, LON int32, ALT int32, EPH uint16, EPV uint16, VEL uint16, VN int16, VE int16, VD int16, COG uint16, FIX_TYPE uint8, SATELLITES_VISIBLE uint8) *HilGps {
	m := HilGps{}
	m.TIME_USEC = TIME_USEC
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.EPH = EPH
	m.EPV = EPV
	m.VEL = VEL
	m.VN = VN
	m.VE = VE
	m.VD = VD
	m.COG = COG
	m.FIX_TYPE = FIX_TYPE
	m.SATELLITES_VISIBLE = SATELLITES_VISIBLE
	return &m
}

// Id returns the HilGps Message ID
func (*HilGps) Id() uint8 {
	return 113
}

// Len returns the HilGps Message Length
func (*HilGps) Len() uint8 {
	return 36
}

// Crc returns the HilGps Message CRC
func (*HilGps) Crc() uint8 {
	return 124
}

// Pack returns a packed byte array which represents a HilGps payload
func (m *HilGps) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.EPH)
	binary.Write(data, binary.LittleEndian, m.EPV)
	binary.Write(data, binary.LittleEndian, m.VEL)
	binary.Write(data, binary.LittleEndian, m.VN)
	binary.Write(data, binary.LittleEndian, m.VE)
	binary.Write(data, binary.LittleEndian, m.VD)
	binary.Write(data, binary.LittleEndian, m.COG)
	binary.Write(data, binary.LittleEndian, m.FIX_TYPE)
	binary.Write(data, binary.LittleEndian, m.SATELLITES_VISIBLE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilGps
func (m *HilGps) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.EPH)
	binary.Read(data, binary.LittleEndian, &m.EPV)
	binary.Read(data, binary.LittleEndian, &m.VEL)
	binary.Read(data, binary.LittleEndian, &m.VN)
	binary.Read(data, binary.LittleEndian, &m.VE)
	binary.Read(data, binary.LittleEndian, &m.VD)
	binary.Read(data, binary.LittleEndian, &m.COG)
	binary.Read(data, binary.LittleEndian, &m.FIX_TYPE)
	binary.Read(data, binary.LittleEndian, &m.SATELLITES_VISIBLE)
}

//
// MESSAGE HIL_OPTICAL_FLOW
//
// MAVLINK_MSG_ID_HIL_OPTICAL_FLOW 114
//
// MAVLINK_MSG_ID_HIL_OPTICAL_FLOW_LEN 26
//
// MAVLINK_MSG_ID_HIL_OPTICAL_FLOW_CRC 119
//
//
type HilOpticalFlow struct {
	TIME_USEC       uint64  // Timestamp (UNIX)
	FLOW_COMP_M_X   float32 // Flow in meters in x-sensor direction, angular-speed compensated
	FLOW_COMP_M_Y   float32 // Flow in meters in y-sensor direction, angular-speed compensated
	GROUND_DISTANCE float32 // Ground distance in meters. Positive value: distance known. Negative value: Unknown distance
	FLOW_X          int16   // Flow in pixels in x-sensor direction
	FLOW_Y          int16   // Flow in pixels in y-sensor direction
	SENSOR_ID       uint8   // Sensor ID
	QUALITY         uint8   // Optical flow quality / confidence. 0: bad, 255: maximum quality
}

// NewHilOpticalFlow returns a new HilOpticalFlow
func NewHilOpticalFlow(TIME_USEC uint64, FLOW_COMP_M_X float32, FLOW_COMP_M_Y float32, GROUND_DISTANCE float32, FLOW_X int16, FLOW_Y int16, SENSOR_ID uint8, QUALITY uint8) *HilOpticalFlow {
	m := HilOpticalFlow{}
	m.TIME_USEC = TIME_USEC
	m.FLOW_COMP_M_X = FLOW_COMP_M_X
	m.FLOW_COMP_M_Y = FLOW_COMP_M_Y
	m.GROUND_DISTANCE = GROUND_DISTANCE
	m.FLOW_X = FLOW_X
	m.FLOW_Y = FLOW_Y
	m.SENSOR_ID = SENSOR_ID
	m.QUALITY = QUALITY
	return &m
}

// Id returns the HilOpticalFlow Message ID
func (*HilOpticalFlow) Id() uint8 {
	return 114
}

// Len returns the HilOpticalFlow Message Length
func (*HilOpticalFlow) Len() uint8 {
	return 26
}

// Crc returns the HilOpticalFlow Message CRC
func (*HilOpticalFlow) Crc() uint8 {
	return 119
}

// Pack returns a packed byte array which represents a HilOpticalFlow payload
func (m *HilOpticalFlow) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.FLOW_COMP_M_X)
	binary.Write(data, binary.LittleEndian, m.FLOW_COMP_M_Y)
	binary.Write(data, binary.LittleEndian, m.GROUND_DISTANCE)
	binary.Write(data, binary.LittleEndian, m.FLOW_X)
	binary.Write(data, binary.LittleEndian, m.FLOW_Y)
	binary.Write(data, binary.LittleEndian, m.SENSOR_ID)
	binary.Write(data, binary.LittleEndian, m.QUALITY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilOpticalFlow
func (m *HilOpticalFlow) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.FLOW_COMP_M_X)
	binary.Read(data, binary.LittleEndian, &m.FLOW_COMP_M_Y)
	binary.Read(data, binary.LittleEndian, &m.GROUND_DISTANCE)
	binary.Read(data, binary.LittleEndian, &m.FLOW_X)
	binary.Read(data, binary.LittleEndian, &m.FLOW_Y)
	binary.Read(data, binary.LittleEndian, &m.SENSOR_ID)
	binary.Read(data, binary.LittleEndian, &m.QUALITY)
}

//
// MESSAGE HIL_STATE_QUATERNION
//
// MAVLINK_MSG_ID_HIL_STATE_QUATERNION 115
//
// MAVLINK_MSG_ID_HIL_STATE_QUATERNION_LEN 64
//
// MAVLINK_MSG_ID_HIL_STATE_QUATERNION_CRC 4
//
//
type HilStateQuaternion struct {
	TIME_USEC           uint64     // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	ATTITUDE_QUATERNION [4]float32 // Vehicle attitude expressed as normalized quaternion in w, x, y, z order (with 1 0 0 0 being the null-rotation)
	ROLLSPEED           float32    // Body frame roll / phi angular speed (rad/s)
	PITCHSPEED          float32    // Body frame pitch / theta angular speed (rad/s)
	YAWSPEED            float32    // Body frame yaw / psi angular speed (rad/s)
	LAT                 int32      // Latitude, expressed as * 1E7
	LON                 int32      // Longitude, expressed as * 1E7
	ALT                 int32      // Altitude in meters, expressed as * 1000 (millimeters)
	VX                  int16      // Ground X Speed (Latitude), expressed as m/s * 100
	VY                  int16      // Ground Y Speed (Longitude), expressed as m/s * 100
	VZ                  int16      // Ground Z Speed (Altitude), expressed as m/s * 100
	IND_AIRSPEED        uint16     // Indicated airspeed, expressed as m/s * 100
	TRUE_AIRSPEED       uint16     // True airspeed, expressed as m/s * 100
	XACC                int16      // X acceleration (mg)
	YACC                int16      // Y acceleration (mg)
	ZACC                int16      // Z acceleration (mg)
}

// NewHilStateQuaternion returns a new HilStateQuaternion
func NewHilStateQuaternion(TIME_USEC uint64, ATTITUDE_QUATERNION [4]float32, ROLLSPEED float32, PITCHSPEED float32, YAWSPEED float32, LAT int32, LON int32, ALT int32, VX int16, VY int16, VZ int16, IND_AIRSPEED uint16, TRUE_AIRSPEED uint16, XACC int16, YACC int16, ZACC int16) *HilStateQuaternion {
	m := HilStateQuaternion{}
	m.TIME_USEC = TIME_USEC
	m.ATTITUDE_QUATERNION = ATTITUDE_QUATERNION
	m.ROLLSPEED = ROLLSPEED
	m.PITCHSPEED = PITCHSPEED
	m.YAWSPEED = YAWSPEED
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.VX = VX
	m.VY = VY
	m.VZ = VZ
	m.IND_AIRSPEED = IND_AIRSPEED
	m.TRUE_AIRSPEED = TRUE_AIRSPEED
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	return &m
}

// Id returns the HilStateQuaternion Message ID
func (*HilStateQuaternion) Id() uint8 {
	return 115
}

// Len returns the HilStateQuaternion Message Length
func (*HilStateQuaternion) Len() uint8 {
	return 64
}

// Crc returns the HilStateQuaternion Message CRC
func (*HilStateQuaternion) Crc() uint8 {
	return 4
}

// Pack returns a packed byte array which represents a HilStateQuaternion payload
func (m *HilStateQuaternion) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.ATTITUDE_QUATERNION)
	binary.Write(data, binary.LittleEndian, m.ROLLSPEED)
	binary.Write(data, binary.LittleEndian, m.PITCHSPEED)
	binary.Write(data, binary.LittleEndian, m.YAWSPEED)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.VX)
	binary.Write(data, binary.LittleEndian, m.VY)
	binary.Write(data, binary.LittleEndian, m.VZ)
	binary.Write(data, binary.LittleEndian, m.IND_AIRSPEED)
	binary.Write(data, binary.LittleEndian, m.TRUE_AIRSPEED)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the HilStateQuaternion
func (m *HilStateQuaternion) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.ATTITUDE_QUATERNION)
	binary.Read(data, binary.LittleEndian, &m.ROLLSPEED)
	binary.Read(data, binary.LittleEndian, &m.PITCHSPEED)
	binary.Read(data, binary.LittleEndian, &m.YAWSPEED)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.VX)
	binary.Read(data, binary.LittleEndian, &m.VY)
	binary.Read(data, binary.LittleEndian, &m.VZ)
	binary.Read(data, binary.LittleEndian, &m.IND_AIRSPEED)
	binary.Read(data, binary.LittleEndian, &m.TRUE_AIRSPEED)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
}

const (
	MAVLINK_MSG_HIL_STATE_QUATERNION_FIELD_attitude_quaternion_LEN = 4
)

//
// MESSAGE SCALED_IMU2
//
// MAVLINK_MSG_ID_SCALED_IMU2 116
//
// MAVLINK_MSG_ID_SCALED_IMU2_LEN 22
//
// MAVLINK_MSG_ID_SCALED_IMU2_CRC 76
//
//
type ScaledImu2 struct {
	TIME_BOOT_MS uint32 // Timestamp (milliseconds since system boot)
	XACC         int16  // X acceleration (mg)
	YACC         int16  // Y acceleration (mg)
	ZACC         int16  // Z acceleration (mg)
	XGYRO        int16  // Angular speed around X axis (millirad /sec)
	YGYRO        int16  // Angular speed around Y axis (millirad /sec)
	ZGYRO        int16  // Angular speed around Z axis (millirad /sec)
	XMAG         int16  // X Magnetic field (milli tesla)
	YMAG         int16  // Y Magnetic field (milli tesla)
	ZMAG         int16  // Z Magnetic field (milli tesla)
}

// NewScaledImu2 returns a new ScaledImu2
func NewScaledImu2(TIME_BOOT_MS uint32, XACC int16, YACC int16, ZACC int16, XGYRO int16, YGYRO int16, ZGYRO int16, XMAG int16, YMAG int16, ZMAG int16) *ScaledImu2 {
	m := ScaledImu2{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.XACC = XACC
	m.YACC = YACC
	m.ZACC = ZACC
	m.XGYRO = XGYRO
	m.YGYRO = YGYRO
	m.ZGYRO = ZGYRO
	m.XMAG = XMAG
	m.YMAG = YMAG
	m.ZMAG = ZMAG
	return &m
}

// Id returns the ScaledImu2 Message ID
func (*ScaledImu2) Id() uint8 {
	return 116
}

// Len returns the ScaledImu2 Message Length
func (*ScaledImu2) Len() uint8 {
	return 22
}

// Crc returns the ScaledImu2 Message CRC
func (*ScaledImu2) Crc() uint8 {
	return 76
}

// Pack returns a packed byte array which represents a ScaledImu2 payload
func (m *ScaledImu2) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.XACC)
	binary.Write(data, binary.LittleEndian, m.YACC)
	binary.Write(data, binary.LittleEndian, m.ZACC)
	binary.Write(data, binary.LittleEndian, m.XGYRO)
	binary.Write(data, binary.LittleEndian, m.YGYRO)
	binary.Write(data, binary.LittleEndian, m.ZGYRO)
	binary.Write(data, binary.LittleEndian, m.XMAG)
	binary.Write(data, binary.LittleEndian, m.YMAG)
	binary.Write(data, binary.LittleEndian, m.ZMAG)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the ScaledImu2
func (m *ScaledImu2) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.XACC)
	binary.Read(data, binary.LittleEndian, &m.YACC)
	binary.Read(data, binary.LittleEndian, &m.ZACC)
	binary.Read(data, binary.LittleEndian, &m.XGYRO)
	binary.Read(data, binary.LittleEndian, &m.YGYRO)
	binary.Read(data, binary.LittleEndian, &m.ZGYRO)
	binary.Read(data, binary.LittleEndian, &m.XMAG)
	binary.Read(data, binary.LittleEndian, &m.YMAG)
	binary.Read(data, binary.LittleEndian, &m.ZMAG)
}

//
// MESSAGE LOG_REQUEST_LIST
//
// MAVLINK_MSG_ID_LOG_REQUEST_LIST 117
//
// MAVLINK_MSG_ID_LOG_REQUEST_LIST_LEN 6
//
// MAVLINK_MSG_ID_LOG_REQUEST_LIST_CRC 128
//
//
type LogRequestList struct {
	START            uint16 // First log id (0 for first available)
	END              uint16 // Last log id (0xffff for last available)
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewLogRequestList returns a new LogRequestList
func NewLogRequestList(START uint16, END uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *LogRequestList {
	m := LogRequestList{}
	m.START = START
	m.END = END
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the LogRequestList Message ID
func (*LogRequestList) Id() uint8 {
	return 117
}

// Len returns the LogRequestList Message Length
func (*LogRequestList) Len() uint8 {
	return 6
}

// Crc returns the LogRequestList Message CRC
func (*LogRequestList) Crc() uint8 {
	return 128
}

// Pack returns a packed byte array which represents a LogRequestList payload
func (m *LogRequestList) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.START)
	binary.Write(data, binary.LittleEndian, m.END)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogRequestList
func (m *LogRequestList) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.START)
	binary.Read(data, binary.LittleEndian, &m.END)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE LOG_ENTRY
//
// MAVLINK_MSG_ID_LOG_ENTRY 118
//
// MAVLINK_MSG_ID_LOG_ENTRY_LEN 14
//
// MAVLINK_MSG_ID_LOG_ENTRY_CRC 56
//
//
type LogEntry struct {
	TIME_UTC     uint32 // UTC timestamp of log in seconds since 1970, or 0 if not available
	SIZE         uint32 // Size of the log (may be approximate) in bytes
	ID           uint16 // Log id
	NUM_LOGS     uint16 // Total number of logs
	LAST_LOG_NUM uint16 // High log number
}

// NewLogEntry returns a new LogEntry
func NewLogEntry(TIME_UTC uint32, SIZE uint32, ID uint16, NUM_LOGS uint16, LAST_LOG_NUM uint16) *LogEntry {
	m := LogEntry{}
	m.TIME_UTC = TIME_UTC
	m.SIZE = SIZE
	m.ID = ID
	m.NUM_LOGS = NUM_LOGS
	m.LAST_LOG_NUM = LAST_LOG_NUM
	return &m
}

// Id returns the LogEntry Message ID
func (*LogEntry) Id() uint8 {
	return 118
}

// Len returns the LogEntry Message Length
func (*LogEntry) Len() uint8 {
	return 14
}

// Crc returns the LogEntry Message CRC
func (*LogEntry) Crc() uint8 {
	return 56
}

// Pack returns a packed byte array which represents a LogEntry payload
func (m *LogEntry) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_UTC)
	binary.Write(data, binary.LittleEndian, m.SIZE)
	binary.Write(data, binary.LittleEndian, m.ID)
	binary.Write(data, binary.LittleEndian, m.NUM_LOGS)
	binary.Write(data, binary.LittleEndian, m.LAST_LOG_NUM)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogEntry
func (m *LogEntry) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_UTC)
	binary.Read(data, binary.LittleEndian, &m.SIZE)
	binary.Read(data, binary.LittleEndian, &m.ID)
	binary.Read(data, binary.LittleEndian, &m.NUM_LOGS)
	binary.Read(data, binary.LittleEndian, &m.LAST_LOG_NUM)
}

//
// MESSAGE LOG_REQUEST_DATA
//
// MAVLINK_MSG_ID_LOG_REQUEST_DATA 119
//
// MAVLINK_MSG_ID_LOG_REQUEST_DATA_LEN 12
//
// MAVLINK_MSG_ID_LOG_REQUEST_DATA_CRC 116
//
//
type LogRequestData struct {
	OFS              uint32 // Offset into the log
	COUNT            uint32 // Number of bytes
	ID               uint16 // Log id (from LOG_ENTRY reply)
	TARGET_SYSTEM    uint8  // System ID
	TARGET_COMPONENT uint8  // Component ID
}

// NewLogRequestData returns a new LogRequestData
func NewLogRequestData(OFS uint32, COUNT uint32, ID uint16, TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *LogRequestData {
	m := LogRequestData{}
	m.OFS = OFS
	m.COUNT = COUNT
	m.ID = ID
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the LogRequestData Message ID
func (*LogRequestData) Id() uint8 {
	return 119
}

// Len returns the LogRequestData Message Length
func (*LogRequestData) Len() uint8 {
	return 12
}

// Crc returns the LogRequestData Message CRC
func (*LogRequestData) Crc() uint8 {
	return 116
}

// Pack returns a packed byte array which represents a LogRequestData payload
func (m *LogRequestData) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.OFS)
	binary.Write(data, binary.LittleEndian, m.COUNT)
	binary.Write(data, binary.LittleEndian, m.ID)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogRequestData
func (m *LogRequestData) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.OFS)
	binary.Read(data, binary.LittleEndian, &m.COUNT)
	binary.Read(data, binary.LittleEndian, &m.ID)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE LOG_DATA
//
// MAVLINK_MSG_ID_LOG_DATA 120
//
// MAVLINK_MSG_ID_LOG_DATA_LEN 97
//
// MAVLINK_MSG_ID_LOG_DATA_CRC 134
//
//
type LogData struct {
	OFS   uint32    // Offset into the log
	ID    uint16    // Log id (from LOG_ENTRY reply)
	COUNT uint8     // Number of bytes (zero for end of log)
	DATA  [90]uint8 // log data
}

// NewLogData returns a new LogData
func NewLogData(OFS uint32, ID uint16, COUNT uint8, DATA [90]uint8) *LogData {
	m := LogData{}
	m.OFS = OFS
	m.ID = ID
	m.COUNT = COUNT
	m.DATA = DATA
	return &m
}

// Id returns the LogData Message ID
func (*LogData) Id() uint8 {
	return 120
}

// Len returns the LogData Message Length
func (*LogData) Len() uint8 {
	return 97
}

// Crc returns the LogData Message CRC
func (*LogData) Crc() uint8 {
	return 134
}

// Pack returns a packed byte array which represents a LogData payload
func (m *LogData) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.OFS)
	binary.Write(data, binary.LittleEndian, m.ID)
	binary.Write(data, binary.LittleEndian, m.COUNT)
	binary.Write(data, binary.LittleEndian, m.DATA)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogData
func (m *LogData) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.OFS)
	binary.Read(data, binary.LittleEndian, &m.ID)
	binary.Read(data, binary.LittleEndian, &m.COUNT)
	binary.Read(data, binary.LittleEndian, &m.DATA)
}

const (
	MAVLINK_MSG_LOG_DATA_FIELD_data_LEN = 90
)

//
// MESSAGE LOG_ERASE
//
// MAVLINK_MSG_ID_LOG_ERASE 121
//
// MAVLINK_MSG_ID_LOG_ERASE_LEN 2
//
// MAVLINK_MSG_ID_LOG_ERASE_CRC 237
//
//
type LogErase struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewLogErase returns a new LogErase
func NewLogErase(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *LogErase {
	m := LogErase{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the LogErase Message ID
func (*LogErase) Id() uint8 {
	return 121
}

// Len returns the LogErase Message Length
func (*LogErase) Len() uint8 {
	return 2
}

// Crc returns the LogErase Message CRC
func (*LogErase) Crc() uint8 {
	return 237
}

// Pack returns a packed byte array which represents a LogErase payload
func (m *LogErase) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogErase
func (m *LogErase) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE LOG_REQUEST_END
//
// MAVLINK_MSG_ID_LOG_REQUEST_END 122
//
// MAVLINK_MSG_ID_LOG_REQUEST_END_LEN 2
//
// MAVLINK_MSG_ID_LOG_REQUEST_END_CRC 203
//
//
type LogRequestEnd struct {
	TARGET_SYSTEM    uint8 // System ID
	TARGET_COMPONENT uint8 // Component ID
}

// NewLogRequestEnd returns a new LogRequestEnd
func NewLogRequestEnd(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8) *LogRequestEnd {
	m := LogRequestEnd{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	return &m
}

// Id returns the LogRequestEnd Message ID
func (*LogRequestEnd) Id() uint8 {
	return 122
}

// Len returns the LogRequestEnd Message Length
func (*LogRequestEnd) Len() uint8 {
	return 2
}

// Crc returns the LogRequestEnd Message CRC
func (*LogRequestEnd) Crc() uint8 {
	return 203
}

// Pack returns a packed byte array which represents a LogRequestEnd payload
func (m *LogRequestEnd) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the LogRequestEnd
func (m *LogRequestEnd) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
}

//
// MESSAGE GPS_INJECT_DATA
//
// MAVLINK_MSG_ID_GPS_INJECT_DATA 123
//
// MAVLINK_MSG_ID_GPS_INJECT_DATA_LEN 113
//
// MAVLINK_MSG_ID_GPS_INJECT_DATA_CRC 250
//
//
type GpsInjectData struct {
	TARGET_SYSTEM    uint8      // System ID
	TARGET_COMPONENT uint8      // Component ID
	LEN              uint8      // data length
	DATA             [110]uint8 // raw data (110 is enough for 12 satellites of RTCMv2)
}

// NewGpsInjectData returns a new GpsInjectData
func NewGpsInjectData(TARGET_SYSTEM uint8, TARGET_COMPONENT uint8, LEN uint8, DATA [110]uint8) *GpsInjectData {
	m := GpsInjectData{}
	m.TARGET_SYSTEM = TARGET_SYSTEM
	m.TARGET_COMPONENT = TARGET_COMPONENT
	m.LEN = LEN
	m.DATA = DATA
	return &m
}

// Id returns the GpsInjectData Message ID
func (*GpsInjectData) Id() uint8 {
	return 123
}

// Len returns the GpsInjectData Message Length
func (*GpsInjectData) Len() uint8 {
	return 113
}

// Crc returns the GpsInjectData Message CRC
func (*GpsInjectData) Crc() uint8 {
	return 250
}

// Pack returns a packed byte array which represents a GpsInjectData payload
func (m *GpsInjectData) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	binary.Write(data, binary.LittleEndian, m.TARGET_COMPONENT)
	binary.Write(data, binary.LittleEndian, m.LEN)
	binary.Write(data, binary.LittleEndian, m.DATA)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GpsInjectData
func (m *GpsInjectData) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
	binary.Read(data, binary.LittleEndian, &m.TARGET_COMPONENT)
	binary.Read(data, binary.LittleEndian, &m.LEN)
	binary.Read(data, binary.LittleEndian, &m.DATA)
}

const (
	MAVLINK_MSG_GPS_INJECT_DATA_FIELD_data_LEN = 110
)

//
// MESSAGE GPS2_RAW
//
// MAVLINK_MSG_ID_GPS2_RAW 124
//
// MAVLINK_MSG_ID_GPS2_RAW_LEN 35
//
// MAVLINK_MSG_ID_GPS2_RAW_CRC 87
//
//
type Gps2Raw struct {
	TIME_USEC          uint64 // Timestamp (microseconds since UNIX epoch or microseconds since system boot)
	LAT                int32  // Latitude (WGS84), in degrees * 1E7
	LON                int32  // Longitude (WGS84), in degrees * 1E7
	ALT                int32  // Altitude (WGS84), in meters * 1000 (positive for up)
	DGPS_AGE           uint32 // Age of DGPS info
	EPH                uint16 // GPS HDOP horizontal dilution of position in cm (m*100). If unknown, set to: UINT16_MAX
	EPV                uint16 // GPS VDOP vertical dilution of position in cm (m*100). If unknown, set to: UINT16_MAX
	VEL                uint16 // GPS ground speed (m/s * 100). If unknown, set to: UINT16_MAX
	COG                uint16 // Course over ground (NOT heading, but direction of movement) in degrees * 100, 0.0..359.99 degrees. If unknown, set to: UINT16_MAX
	FIX_TYPE           uint8  // 0-1: no fix, 2: 2D fix, 3: 3D fix, 4: DGPS fix, 5: RTK Fix. Some applications will not use the value of this field unless it is at least two, so always correctly fill in the fix.
	SATELLITES_VISIBLE uint8  // Number of satellites visible. If unknown, set to 255
	DGPS_NUMCH         uint8  // Number of DGPS satellites
}

// NewGps2Raw returns a new Gps2Raw
func NewGps2Raw(TIME_USEC uint64, LAT int32, LON int32, ALT int32, DGPS_AGE uint32, EPH uint16, EPV uint16, VEL uint16, COG uint16, FIX_TYPE uint8, SATELLITES_VISIBLE uint8, DGPS_NUMCH uint8) *Gps2Raw {
	m := Gps2Raw{}
	m.TIME_USEC = TIME_USEC
	m.LAT = LAT
	m.LON = LON
	m.ALT = ALT
	m.DGPS_AGE = DGPS_AGE
	m.EPH = EPH
	m.EPV = EPV
	m.VEL = VEL
	m.COG = COG
	m.FIX_TYPE = FIX_TYPE
	m.SATELLITES_VISIBLE = SATELLITES_VISIBLE
	m.DGPS_NUMCH = DGPS_NUMCH
	return &m
}

// Id returns the Gps2Raw Message ID
func (*Gps2Raw) Id() uint8 {
	return 124
}

// Len returns the Gps2Raw Message Length
func (*Gps2Raw) Len() uint8 {
	return 35
}

// Crc returns the Gps2Raw Message CRC
func (*Gps2Raw) Crc() uint8 {
	return 87
}

// Pack returns a packed byte array which represents a Gps2Raw payload
func (m *Gps2Raw) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.ALT)
	binary.Write(data, binary.LittleEndian, m.DGPS_AGE)
	binary.Write(data, binary.LittleEndian, m.EPH)
	binary.Write(data, binary.LittleEndian, m.EPV)
	binary.Write(data, binary.LittleEndian, m.VEL)
	binary.Write(data, binary.LittleEndian, m.COG)
	binary.Write(data, binary.LittleEndian, m.FIX_TYPE)
	binary.Write(data, binary.LittleEndian, m.SATELLITES_VISIBLE)
	binary.Write(data, binary.LittleEndian, m.DGPS_NUMCH)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Gps2Raw
func (m *Gps2Raw) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.ALT)
	binary.Read(data, binary.LittleEndian, &m.DGPS_AGE)
	binary.Read(data, binary.LittleEndian, &m.EPH)
	binary.Read(data, binary.LittleEndian, &m.EPV)
	binary.Read(data, binary.LittleEndian, &m.VEL)
	binary.Read(data, binary.LittleEndian, &m.COG)
	binary.Read(data, binary.LittleEndian, &m.FIX_TYPE)
	binary.Read(data, binary.LittleEndian, &m.SATELLITES_VISIBLE)
	binary.Read(data, binary.LittleEndian, &m.DGPS_NUMCH)
}

//
// MESSAGE POWER_STATUS
//
// MAVLINK_MSG_ID_POWER_STATUS 125
//
// MAVLINK_MSG_ID_POWER_STATUS_LEN 6
//
// MAVLINK_MSG_ID_POWER_STATUS_CRC 203
//
//
type PowerStatus struct {
	VCC    uint16 // 5V rail voltage in millivolts
	VSERVO uint16 // servo rail voltage in millivolts
	FLAGS  uint16 // power supply status flags (see MAV_POWER_STATUS enum)
}

// NewPowerStatus returns a new PowerStatus
func NewPowerStatus(VCC uint16, VSERVO uint16, FLAGS uint16) *PowerStatus {
	m := PowerStatus{}
	m.VCC = VCC
	m.VSERVO = VSERVO
	m.FLAGS = FLAGS
	return &m
}

// Id returns the PowerStatus Message ID
func (*PowerStatus) Id() uint8 {
	return 125
}

// Len returns the PowerStatus Message Length
func (*PowerStatus) Len() uint8 {
	return 6
}

// Crc returns the PowerStatus Message CRC
func (*PowerStatus) Crc() uint8 {
	return 203
}

// Pack returns a packed byte array which represents a PowerStatus payload
func (m *PowerStatus) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.VCC)
	binary.Write(data, binary.LittleEndian, m.VSERVO)
	binary.Write(data, binary.LittleEndian, m.FLAGS)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the PowerStatus
func (m *PowerStatus) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.VCC)
	binary.Read(data, binary.LittleEndian, &m.VSERVO)
	binary.Read(data, binary.LittleEndian, &m.FLAGS)
}

//
// MESSAGE SERIAL_CONTROL
//
// MAVLINK_MSG_ID_SERIAL_CONTROL 126
//
// MAVLINK_MSG_ID_SERIAL_CONTROL_LEN 79
//
// MAVLINK_MSG_ID_SERIAL_CONTROL_CRC 220
//
//
type SerialControl struct {
	BAUDRATE uint32    // Baudrate of transfer. Zero means no change.
	TIMEOUT  uint16    // Timeout for reply data in milliseconds
	DEVICE   uint8     // See SERIAL_CONTROL_DEV enum
	FLAGS    uint8     // See SERIAL_CONTROL_FLAG enum
	COUNT    uint8     // how many bytes in this transfer
	DATA     [70]uint8 // serial data
}

// NewSerialControl returns a new SerialControl
func NewSerialControl(BAUDRATE uint32, TIMEOUT uint16, DEVICE uint8, FLAGS uint8, COUNT uint8, DATA [70]uint8) *SerialControl {
	m := SerialControl{}
	m.BAUDRATE = BAUDRATE
	m.TIMEOUT = TIMEOUT
	m.DEVICE = DEVICE
	m.FLAGS = FLAGS
	m.COUNT = COUNT
	m.DATA = DATA
	return &m
}

// Id returns the SerialControl Message ID
func (*SerialControl) Id() uint8 {
	return 126
}

// Len returns the SerialControl Message Length
func (*SerialControl) Len() uint8 {
	return 79
}

// Crc returns the SerialControl Message CRC
func (*SerialControl) Crc() uint8 {
	return 220
}

// Pack returns a packed byte array which represents a SerialControl payload
func (m *SerialControl) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.BAUDRATE)
	binary.Write(data, binary.LittleEndian, m.TIMEOUT)
	binary.Write(data, binary.LittleEndian, m.DEVICE)
	binary.Write(data, binary.LittleEndian, m.FLAGS)
	binary.Write(data, binary.LittleEndian, m.COUNT)
	binary.Write(data, binary.LittleEndian, m.DATA)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the SerialControl
func (m *SerialControl) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.BAUDRATE)
	binary.Read(data, binary.LittleEndian, &m.TIMEOUT)
	binary.Read(data, binary.LittleEndian, &m.DEVICE)
	binary.Read(data, binary.LittleEndian, &m.FLAGS)
	binary.Read(data, binary.LittleEndian, &m.COUNT)
	binary.Read(data, binary.LittleEndian, &m.DATA)
}

const (
	MAVLINK_MSG_SERIAL_CONTROL_FIELD_data_LEN = 70
)

//
// MESSAGE GPS_RTK
//
// MAVLINK_MSG_ID_GPS_RTK 127
//
// MAVLINK_MSG_ID_GPS_RTK_LEN 35
//
// MAVLINK_MSG_ID_GPS_RTK_CRC 25
//
//
type GpsRtk struct {
	TIME_LAST_BASELINE_MS uint32 // Time since boot of last baseline message received in ms.
	TOW                   uint32 // GPS Time of Week of last baseline
	BASELINE_A_MM         int32  // Current baseline in ECEF x or NED north component in mm.
	BASELINE_B_MM         int32  // Current baseline in ECEF y or NED east component in mm.
	BASELINE_C_MM         int32  // Current baseline in ECEF z or NED down component in mm.
	ACCURACY              uint32 // Current estimate of baseline accuracy.
	IAR_NUM_HYPOTHESES    int32  // Current number of integer ambiguity hypotheses.
	WN                    uint16 // GPS Week Number of last baseline
	RTK_RECEIVER_ID       uint8  // Identification of connected RTK receiver.
	RTK_HEALTH            uint8  // GPS-specific health report for RTK data.
	RTK_RATE              uint8  // Rate of baseline messages being received by GPS, in HZ
	NSATS                 uint8  // Current number of sats used for RTK calculation.
	BASELINE_COORDS_TYPE  uint8  // Coordinate system of baseline. 0 == ECEF, 1 == NED
}

// NewGpsRtk returns a new GpsRtk
func NewGpsRtk(TIME_LAST_BASELINE_MS uint32, TOW uint32, BASELINE_A_MM int32, BASELINE_B_MM int32, BASELINE_C_MM int32, ACCURACY uint32, IAR_NUM_HYPOTHESES int32, WN uint16, RTK_RECEIVER_ID uint8, RTK_HEALTH uint8, RTK_RATE uint8, NSATS uint8, BASELINE_COORDS_TYPE uint8) *GpsRtk {
	m := GpsRtk{}
	m.TIME_LAST_BASELINE_MS = TIME_LAST_BASELINE_MS
	m.TOW = TOW
	m.BASELINE_A_MM = BASELINE_A_MM
	m.BASELINE_B_MM = BASELINE_B_MM
	m.BASELINE_C_MM = BASELINE_C_MM
	m.ACCURACY = ACCURACY
	m.IAR_NUM_HYPOTHESES = IAR_NUM_HYPOTHESES
	m.WN = WN
	m.RTK_RECEIVER_ID = RTK_RECEIVER_ID
	m.RTK_HEALTH = RTK_HEALTH
	m.RTK_RATE = RTK_RATE
	m.NSATS = NSATS
	m.BASELINE_COORDS_TYPE = BASELINE_COORDS_TYPE
	return &m
}

// Id returns the GpsRtk Message ID
func (*GpsRtk) Id() uint8 {
	return 127
}

// Len returns the GpsRtk Message Length
func (*GpsRtk) Len() uint8 {
	return 35
}

// Crc returns the GpsRtk Message CRC
func (*GpsRtk) Crc() uint8 {
	return 25
}

// Pack returns a packed byte array which represents a GpsRtk payload
func (m *GpsRtk) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_LAST_BASELINE_MS)
	binary.Write(data, binary.LittleEndian, m.TOW)
	binary.Write(data, binary.LittleEndian, m.BASELINE_A_MM)
	binary.Write(data, binary.LittleEndian, m.BASELINE_B_MM)
	binary.Write(data, binary.LittleEndian, m.BASELINE_C_MM)
	binary.Write(data, binary.LittleEndian, m.ACCURACY)
	binary.Write(data, binary.LittleEndian, m.IAR_NUM_HYPOTHESES)
	binary.Write(data, binary.LittleEndian, m.WN)
	binary.Write(data, binary.LittleEndian, m.RTK_RECEIVER_ID)
	binary.Write(data, binary.LittleEndian, m.RTK_HEALTH)
	binary.Write(data, binary.LittleEndian, m.RTK_RATE)
	binary.Write(data, binary.LittleEndian, m.NSATS)
	binary.Write(data, binary.LittleEndian, m.BASELINE_COORDS_TYPE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the GpsRtk
func (m *GpsRtk) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_LAST_BASELINE_MS)
	binary.Read(data, binary.LittleEndian, &m.TOW)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_A_MM)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_B_MM)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_C_MM)
	binary.Read(data, binary.LittleEndian, &m.ACCURACY)
	binary.Read(data, binary.LittleEndian, &m.IAR_NUM_HYPOTHESES)
	binary.Read(data, binary.LittleEndian, &m.WN)
	binary.Read(data, binary.LittleEndian, &m.RTK_RECEIVER_ID)
	binary.Read(data, binary.LittleEndian, &m.RTK_HEALTH)
	binary.Read(data, binary.LittleEndian, &m.RTK_RATE)
	binary.Read(data, binary.LittleEndian, &m.NSATS)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_COORDS_TYPE)
}

//
// MESSAGE GPS2_RTK
//
// MAVLINK_MSG_ID_GPS2_RTK 128
//
// MAVLINK_MSG_ID_GPS2_RTK_LEN 35
//
// MAVLINK_MSG_ID_GPS2_RTK_CRC 226
//
//
type Gps2Rtk struct {
	TIME_LAST_BASELINE_MS uint32 // Time since boot of last baseline message received in ms.
	TOW                   uint32 // GPS Time of Week of last baseline
	BASELINE_A_MM         int32  // Current baseline in ECEF x or NED north component in mm.
	BASELINE_B_MM         int32  // Current baseline in ECEF y or NED east component in mm.
	BASELINE_C_MM         int32  // Current baseline in ECEF z or NED down component in mm.
	ACCURACY              uint32 // Current estimate of baseline accuracy.
	IAR_NUM_HYPOTHESES    int32  // Current number of integer ambiguity hypotheses.
	WN                    uint16 // GPS Week Number of last baseline
	RTK_RECEIVER_ID       uint8  // Identification of connected RTK receiver.
	RTK_HEALTH            uint8  // GPS-specific health report for RTK data.
	RTK_RATE              uint8  // Rate of baseline messages being received by GPS, in HZ
	NSATS                 uint8  // Current number of sats used for RTK calculation.
	BASELINE_COORDS_TYPE  uint8  // Coordinate system of baseline. 0 == ECEF, 1 == NED
}

// NewGps2Rtk returns a new Gps2Rtk
func NewGps2Rtk(TIME_LAST_BASELINE_MS uint32, TOW uint32, BASELINE_A_MM int32, BASELINE_B_MM int32, BASELINE_C_MM int32, ACCURACY uint32, IAR_NUM_HYPOTHESES int32, WN uint16, RTK_RECEIVER_ID uint8, RTK_HEALTH uint8, RTK_RATE uint8, NSATS uint8, BASELINE_COORDS_TYPE uint8) *Gps2Rtk {
	m := Gps2Rtk{}
	m.TIME_LAST_BASELINE_MS = TIME_LAST_BASELINE_MS
	m.TOW = TOW
	m.BASELINE_A_MM = BASELINE_A_MM
	m.BASELINE_B_MM = BASELINE_B_MM
	m.BASELINE_C_MM = BASELINE_C_MM
	m.ACCURACY = ACCURACY
	m.IAR_NUM_HYPOTHESES = IAR_NUM_HYPOTHESES
	m.WN = WN
	m.RTK_RECEIVER_ID = RTK_RECEIVER_ID
	m.RTK_HEALTH = RTK_HEALTH
	m.RTK_RATE = RTK_RATE
	m.NSATS = NSATS
	m.BASELINE_COORDS_TYPE = BASELINE_COORDS_TYPE
	return &m
}

// Id returns the Gps2Rtk Message ID
func (*Gps2Rtk) Id() uint8 {
	return 128
}

// Len returns the Gps2Rtk Message Length
func (*Gps2Rtk) Len() uint8 {
	return 35
}

// Crc returns the Gps2Rtk Message CRC
func (*Gps2Rtk) Crc() uint8 {
	return 226
}

// Pack returns a packed byte array which represents a Gps2Rtk payload
func (m *Gps2Rtk) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_LAST_BASELINE_MS)
	binary.Write(data, binary.LittleEndian, m.TOW)
	binary.Write(data, binary.LittleEndian, m.BASELINE_A_MM)
	binary.Write(data, binary.LittleEndian, m.BASELINE_B_MM)
	binary.Write(data, binary.LittleEndian, m.BASELINE_C_MM)
	binary.Write(data, binary.LittleEndian, m.ACCURACY)
	binary.Write(data, binary.LittleEndian, m.IAR_NUM_HYPOTHESES)
	binary.Write(data, binary.LittleEndian, m.WN)
	binary.Write(data, binary.LittleEndian, m.RTK_RECEIVER_ID)
	binary.Write(data, binary.LittleEndian, m.RTK_HEALTH)
	binary.Write(data, binary.LittleEndian, m.RTK_RATE)
	binary.Write(data, binary.LittleEndian, m.NSATS)
	binary.Write(data, binary.LittleEndian, m.BASELINE_COORDS_TYPE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Gps2Rtk
func (m *Gps2Rtk) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_LAST_BASELINE_MS)
	binary.Read(data, binary.LittleEndian, &m.TOW)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_A_MM)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_B_MM)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_C_MM)
	binary.Read(data, binary.LittleEndian, &m.ACCURACY)
	binary.Read(data, binary.LittleEndian, &m.IAR_NUM_HYPOTHESES)
	binary.Read(data, binary.LittleEndian, &m.WN)
	binary.Read(data, binary.LittleEndian, &m.RTK_RECEIVER_ID)
	binary.Read(data, binary.LittleEndian, &m.RTK_HEALTH)
	binary.Read(data, binary.LittleEndian, &m.RTK_RATE)
	binary.Read(data, binary.LittleEndian, &m.NSATS)
	binary.Read(data, binary.LittleEndian, &m.BASELINE_COORDS_TYPE)
}

//
// MESSAGE DATA_TRANSMISSION_HANDSHAKE
//
// MAVLINK_MSG_ID_DATA_TRANSMISSION_HANDSHAKE 130
//
// MAVLINK_MSG_ID_DATA_TRANSMISSION_HANDSHAKE_LEN 13
//
// MAVLINK_MSG_ID_DATA_TRANSMISSION_HANDSHAKE_CRC 29
//
//
type DataTransmissionHandshake struct {
	SIZE        uint32 // total data size in bytes (set on ACK only)
	WIDTH       uint16 // Width of a matrix or image
	HEIGHT      uint16 // Height of a matrix or image
	PACKETS     uint16 // number of packets being sent (set on ACK only)
	TYPE        uint8  // type of requested/acknowledged data (as defined in ENUM DATA_TYPES in mavlink/include/mavlink_types.h)
	PAYLOAD     uint8  // payload size per packet (normally 253 byte, see DATA field size in message ENCAPSULATED_DATA) (set on ACK only)
	JPG_QUALITY uint8  // JPEG quality out of [1,100]
}

// NewDataTransmissionHandshake returns a new DataTransmissionHandshake
func NewDataTransmissionHandshake(SIZE uint32, WIDTH uint16, HEIGHT uint16, PACKETS uint16, TYPE uint8, PAYLOAD uint8, JPG_QUALITY uint8) *DataTransmissionHandshake {
	m := DataTransmissionHandshake{}
	m.SIZE = SIZE
	m.WIDTH = WIDTH
	m.HEIGHT = HEIGHT
	m.PACKETS = PACKETS
	m.TYPE = TYPE
	m.PAYLOAD = PAYLOAD
	m.JPG_QUALITY = JPG_QUALITY
	return &m
}

// Id returns the DataTransmissionHandshake Message ID
func (*DataTransmissionHandshake) Id() uint8 {
	return 130
}

// Len returns the DataTransmissionHandshake Message Length
func (*DataTransmissionHandshake) Len() uint8 {
	return 13
}

// Crc returns the DataTransmissionHandshake Message CRC
func (*DataTransmissionHandshake) Crc() uint8 {
	return 29
}

// Pack returns a packed byte array which represents a DataTransmissionHandshake payload
func (m *DataTransmissionHandshake) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SIZE)
	binary.Write(data, binary.LittleEndian, m.WIDTH)
	binary.Write(data, binary.LittleEndian, m.HEIGHT)
	binary.Write(data, binary.LittleEndian, m.PACKETS)
	binary.Write(data, binary.LittleEndian, m.TYPE)
	binary.Write(data, binary.LittleEndian, m.PAYLOAD)
	binary.Write(data, binary.LittleEndian, m.JPG_QUALITY)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the DataTransmissionHandshake
func (m *DataTransmissionHandshake) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SIZE)
	binary.Read(data, binary.LittleEndian, &m.WIDTH)
	binary.Read(data, binary.LittleEndian, &m.HEIGHT)
	binary.Read(data, binary.LittleEndian, &m.PACKETS)
	binary.Read(data, binary.LittleEndian, &m.TYPE)
	binary.Read(data, binary.LittleEndian, &m.PAYLOAD)
	binary.Read(data, binary.LittleEndian, &m.JPG_QUALITY)
}

//
// MESSAGE ENCAPSULATED_DATA
//
// MAVLINK_MSG_ID_ENCAPSULATED_DATA 131
//
// MAVLINK_MSG_ID_ENCAPSULATED_DATA_LEN 255
//
// MAVLINK_MSG_ID_ENCAPSULATED_DATA_CRC 223
//
//
type EncapsulatedData struct {
	SEQNR uint16     // sequence number (starting with 0 on every transmission)
	DATA  [253]uint8 // image data bytes
}

// NewEncapsulatedData returns a new EncapsulatedData
func NewEncapsulatedData(SEQNR uint16, DATA [253]uint8) *EncapsulatedData {
	m := EncapsulatedData{}
	m.SEQNR = SEQNR
	m.DATA = DATA
	return &m
}

// Id returns the EncapsulatedData Message ID
func (*EncapsulatedData) Id() uint8 {
	return 131
}

// Len returns the EncapsulatedData Message Length
func (*EncapsulatedData) Len() uint8 {
	return 255
}

// Crc returns the EncapsulatedData Message CRC
func (*EncapsulatedData) Crc() uint8 {
	return 223
}

// Pack returns a packed byte array which represents a EncapsulatedData payload
func (m *EncapsulatedData) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEQNR)
	binary.Write(data, binary.LittleEndian, m.DATA)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the EncapsulatedData
func (m *EncapsulatedData) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEQNR)
	binary.Read(data, binary.LittleEndian, &m.DATA)
}

const (
	MAVLINK_MSG_ENCAPSULATED_DATA_FIELD_data_LEN = 253
)

//
// MESSAGE DISTANCE_SENSOR
//
// MAVLINK_MSG_ID_DISTANCE_SENSOR 132
//
// MAVLINK_MSG_ID_DISTANCE_SENSOR_LEN 14
//
// MAVLINK_MSG_ID_DISTANCE_SENSOR_CRC 85
//
//
type DistanceSensor struct {
	TIME_BOOT_MS     uint32 // Time since system boot
	MIN_DISTANCE     uint16 // Minimum distance the sensor can measure in centimeters
	MAX_DISTANCE     uint16 // Maximum distance the sensor can measure in centimeters
	CURRENT_DISTANCE uint16 // Current distance reading
	TYPE             uint8  // Type from MAV_DISTANCE_SENSOR enum.
	ID               uint8  // Onboard ID of the sensor
	ORIENTATION      uint8  // Direction the sensor faces from FIXME enum.
	COVARIANCE       uint8  // Measurement covariance in centimeters, 0 for unknown / invalid readings
}

// NewDistanceSensor returns a new DistanceSensor
func NewDistanceSensor(TIME_BOOT_MS uint32, MIN_DISTANCE uint16, MAX_DISTANCE uint16, CURRENT_DISTANCE uint16, TYPE uint8, ID uint8, ORIENTATION uint8, COVARIANCE uint8) *DistanceSensor {
	m := DistanceSensor{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.MIN_DISTANCE = MIN_DISTANCE
	m.MAX_DISTANCE = MAX_DISTANCE
	m.CURRENT_DISTANCE = CURRENT_DISTANCE
	m.TYPE = TYPE
	m.ID = ID
	m.ORIENTATION = ORIENTATION
	m.COVARIANCE = COVARIANCE
	return &m
}

// Id returns the DistanceSensor Message ID
func (*DistanceSensor) Id() uint8 {
	return 132
}

// Len returns the DistanceSensor Message Length
func (*DistanceSensor) Len() uint8 {
	return 14
}

// Crc returns the DistanceSensor Message CRC
func (*DistanceSensor) Crc() uint8 {
	return 85
}

// Pack returns a packed byte array which represents a DistanceSensor payload
func (m *DistanceSensor) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.MIN_DISTANCE)
	binary.Write(data, binary.LittleEndian, m.MAX_DISTANCE)
	binary.Write(data, binary.LittleEndian, m.CURRENT_DISTANCE)
	binary.Write(data, binary.LittleEndian, m.TYPE)
	binary.Write(data, binary.LittleEndian, m.ID)
	binary.Write(data, binary.LittleEndian, m.ORIENTATION)
	binary.Write(data, binary.LittleEndian, m.COVARIANCE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the DistanceSensor
func (m *DistanceSensor) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.MIN_DISTANCE)
	binary.Read(data, binary.LittleEndian, &m.MAX_DISTANCE)
	binary.Read(data, binary.LittleEndian, &m.CURRENT_DISTANCE)
	binary.Read(data, binary.LittleEndian, &m.TYPE)
	binary.Read(data, binary.LittleEndian, &m.ID)
	binary.Read(data, binary.LittleEndian, &m.ORIENTATION)
	binary.Read(data, binary.LittleEndian, &m.COVARIANCE)
}

//
// MESSAGE TERRAIN_REQUEST
//
// MAVLINK_MSG_ID_TERRAIN_REQUEST 133
//
// MAVLINK_MSG_ID_TERRAIN_REQUEST_LEN 18
//
// MAVLINK_MSG_ID_TERRAIN_REQUEST_CRC 6
//
//
type TerrainRequest struct {
	MASK         uint64 // Bitmask of requested 4x4 grids (row major 8x7 array of grids, 56 bits)
	LAT          int32  // Latitude of SW corner of first grid (degrees *10^7)
	LON          int32  // Longitude of SW corner of first grid (in degrees *10^7)
	GRID_SPACING uint16 // Grid spacing in meters
}

// NewTerrainRequest returns a new TerrainRequest
func NewTerrainRequest(MASK uint64, LAT int32, LON int32, GRID_SPACING uint16) *TerrainRequest {
	m := TerrainRequest{}
	m.MASK = MASK
	m.LAT = LAT
	m.LON = LON
	m.GRID_SPACING = GRID_SPACING
	return &m
}

// Id returns the TerrainRequest Message ID
func (*TerrainRequest) Id() uint8 {
	return 133
}

// Len returns the TerrainRequest Message Length
func (*TerrainRequest) Len() uint8 {
	return 18
}

// Crc returns the TerrainRequest Message CRC
func (*TerrainRequest) Crc() uint8 {
	return 6
}

// Pack returns a packed byte array which represents a TerrainRequest payload
func (m *TerrainRequest) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.MASK)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.GRID_SPACING)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the TerrainRequest
func (m *TerrainRequest) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.MASK)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.GRID_SPACING)
}

//
// MESSAGE TERRAIN_DATA
//
// MAVLINK_MSG_ID_TERRAIN_DATA 134
//
// MAVLINK_MSG_ID_TERRAIN_DATA_LEN 43
//
// MAVLINK_MSG_ID_TERRAIN_DATA_CRC 229
//
//
type TerrainData struct {
	LAT          int32     // Latitude of SW corner of first grid (degrees *10^7)
	LON          int32     // Longitude of SW corner of first grid (in degrees *10^7)
	GRID_SPACING uint16    // Grid spacing in meters
	DATA         [16]int16 // Terrain data in meters AMSL
	GRIDBIT      uint8     // bit within the terrain request mask
}

// NewTerrainData returns a new TerrainData
func NewTerrainData(LAT int32, LON int32, GRID_SPACING uint16, DATA [16]int16, GRIDBIT uint8) *TerrainData {
	m := TerrainData{}
	m.LAT = LAT
	m.LON = LON
	m.GRID_SPACING = GRID_SPACING
	m.DATA = DATA
	m.GRIDBIT = GRIDBIT
	return &m
}

// Id returns the TerrainData Message ID
func (*TerrainData) Id() uint8 {
	return 134
}

// Len returns the TerrainData Message Length
func (*TerrainData) Len() uint8 {
	return 43
}

// Crc returns the TerrainData Message CRC
func (*TerrainData) Crc() uint8 {
	return 229
}

// Pack returns a packed byte array which represents a TerrainData payload
func (m *TerrainData) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.GRID_SPACING)
	binary.Write(data, binary.LittleEndian, m.DATA)
	binary.Write(data, binary.LittleEndian, m.GRIDBIT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the TerrainData
func (m *TerrainData) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.GRID_SPACING)
	binary.Read(data, binary.LittleEndian, &m.DATA)
	binary.Read(data, binary.LittleEndian, &m.GRIDBIT)
}

const (
	MAVLINK_MSG_TERRAIN_DATA_FIELD_data_LEN = 16
)

//
// MESSAGE TERRAIN_CHECK
//
// MAVLINK_MSG_ID_TERRAIN_CHECK 135
//
// MAVLINK_MSG_ID_TERRAIN_CHECK_LEN 8
//
// MAVLINK_MSG_ID_TERRAIN_CHECK_CRC 203
//
//
type TerrainCheck struct {
	LAT int32 // Latitude (degrees *10^7)
	LON int32 // Longitude (degrees *10^7)
}

// NewTerrainCheck returns a new TerrainCheck
func NewTerrainCheck(LAT int32, LON int32) *TerrainCheck {
	m := TerrainCheck{}
	m.LAT = LAT
	m.LON = LON
	return &m
}

// Id returns the TerrainCheck Message ID
func (*TerrainCheck) Id() uint8 {
	return 135
}

// Len returns the TerrainCheck Message Length
func (*TerrainCheck) Len() uint8 {
	return 8
}

// Crc returns the TerrainCheck Message CRC
func (*TerrainCheck) Crc() uint8 {
	return 203
}

// Pack returns a packed byte array which represents a TerrainCheck payload
func (m *TerrainCheck) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the TerrainCheck
func (m *TerrainCheck) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
}

//
// MESSAGE TERRAIN_REPORT
//
// MAVLINK_MSG_ID_TERRAIN_REPORT 136
//
// MAVLINK_MSG_ID_TERRAIN_REPORT_LEN 22
//
// MAVLINK_MSG_ID_TERRAIN_REPORT_CRC 1
//
//
type TerrainReport struct {
	LAT            int32   // Latitude (degrees *10^7)
	LON            int32   // Longitude (degrees *10^7)
	TERRAIN_HEIGHT float32 // Terrain height in meters AMSL
	CURRENT_HEIGHT float32 // Current vehicle height above lat/lon terrain height (meters)
	SPACING        uint16  // grid spacing (zero if terrain at this location unavailable)
	PENDING        uint16  // Number of 4x4 terrain blocks waiting to be received or read from disk
	LOADED         uint16  // Number of 4x4 terrain blocks in memory
}

// NewTerrainReport returns a new TerrainReport
func NewTerrainReport(LAT int32, LON int32, TERRAIN_HEIGHT float32, CURRENT_HEIGHT float32, SPACING uint16, PENDING uint16, LOADED uint16) *TerrainReport {
	m := TerrainReport{}
	m.LAT = LAT
	m.LON = LON
	m.TERRAIN_HEIGHT = TERRAIN_HEIGHT
	m.CURRENT_HEIGHT = CURRENT_HEIGHT
	m.SPACING = SPACING
	m.PENDING = PENDING
	m.LOADED = LOADED
	return &m
}

// Id returns the TerrainReport Message ID
func (*TerrainReport) Id() uint8 {
	return 136
}

// Len returns the TerrainReport Message Length
func (*TerrainReport) Len() uint8 {
	return 22
}

// Crc returns the TerrainReport Message CRC
func (*TerrainReport) Crc() uint8 {
	return 1
}

// Pack returns a packed byte array which represents a TerrainReport payload
func (m *TerrainReport) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.LAT)
	binary.Write(data, binary.LittleEndian, m.LON)
	binary.Write(data, binary.LittleEndian, m.TERRAIN_HEIGHT)
	binary.Write(data, binary.LittleEndian, m.CURRENT_HEIGHT)
	binary.Write(data, binary.LittleEndian, m.SPACING)
	binary.Write(data, binary.LittleEndian, m.PENDING)
	binary.Write(data, binary.LittleEndian, m.LOADED)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the TerrainReport
func (m *TerrainReport) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.LAT)
	binary.Read(data, binary.LittleEndian, &m.LON)
	binary.Read(data, binary.LittleEndian, &m.TERRAIN_HEIGHT)
	binary.Read(data, binary.LittleEndian, &m.CURRENT_HEIGHT)
	binary.Read(data, binary.LittleEndian, &m.SPACING)
	binary.Read(data, binary.LittleEndian, &m.PENDING)
	binary.Read(data, binary.LittleEndian, &m.LOADED)
}

//
// MESSAGE BATTERY_STATUS
//
// MAVLINK_MSG_ID_BATTERY_STATUS 147
//
// MAVLINK_MSG_ID_BATTERY_STATUS_LEN 24
//
// MAVLINK_MSG_ID_BATTERY_STATUS_CRC 177
//
//
type BatteryStatus struct {
	CURRENT_CONSUMED  int32  // Consumed charge, in milliampere hours (1 = 1 mAh), -1: autopilot does not provide mAh consumption estimate
	ENERGY_CONSUMED   int32  // Consumed energy, in 100*Joules (integrated U*I*dt)  (1 = 100 Joule), -1: autopilot does not provide energy consumption estimate
	VOLTAGE_CELL_1    uint16 // Battery voltage of cell 1, in millivolts (1 = 1 millivolt)
	VOLTAGE_CELL_2    uint16 // Battery voltage of cell 2, in millivolts (1 = 1 millivolt), -1: no cell
	VOLTAGE_CELL_3    uint16 // Battery voltage of cell 3, in millivolts (1 = 1 millivolt), -1: no cell
	VOLTAGE_CELL_4    uint16 // Battery voltage of cell 4, in millivolts (1 = 1 millivolt), -1: no cell
	VOLTAGE_CELL_5    uint16 // Battery voltage of cell 5, in millivolts (1 = 1 millivolt), -1: no cell
	VOLTAGE_CELL_6    uint16 // Battery voltage of cell 6, in millivolts (1 = 1 millivolt), -1: no cell
	CURRENT_BATTERY   int16  // Battery current, in 10*milliamperes (1 = 10 milliampere), -1: autopilot does not measure the current
	ACCU_ID           uint8  // Accupack ID
	BATTERY_REMAINING int8   // Remaining battery energy: (0%: 0, 100%: 100), -1: autopilot does not estimate the remaining battery
}

// NewBatteryStatus returns a new BatteryStatus
func NewBatteryStatus(CURRENT_CONSUMED int32, ENERGY_CONSUMED int32, VOLTAGE_CELL_1 uint16, VOLTAGE_CELL_2 uint16, VOLTAGE_CELL_3 uint16, VOLTAGE_CELL_4 uint16, VOLTAGE_CELL_5 uint16, VOLTAGE_CELL_6 uint16, CURRENT_BATTERY int16, ACCU_ID uint8, BATTERY_REMAINING int8) *BatteryStatus {
	m := BatteryStatus{}
	m.CURRENT_CONSUMED = CURRENT_CONSUMED
	m.ENERGY_CONSUMED = ENERGY_CONSUMED
	m.VOLTAGE_CELL_1 = VOLTAGE_CELL_1
	m.VOLTAGE_CELL_2 = VOLTAGE_CELL_2
	m.VOLTAGE_CELL_3 = VOLTAGE_CELL_3
	m.VOLTAGE_CELL_4 = VOLTAGE_CELL_4
	m.VOLTAGE_CELL_5 = VOLTAGE_CELL_5
	m.VOLTAGE_CELL_6 = VOLTAGE_CELL_6
	m.CURRENT_BATTERY = CURRENT_BATTERY
	m.ACCU_ID = ACCU_ID
	m.BATTERY_REMAINING = BATTERY_REMAINING
	return &m
}

// Id returns the BatteryStatus Message ID
func (*BatteryStatus) Id() uint8 {
	return 147
}

// Len returns the BatteryStatus Message Length
func (*BatteryStatus) Len() uint8 {
	return 24
}

// Crc returns the BatteryStatus Message CRC
func (*BatteryStatus) Crc() uint8 {
	return 177
}

// Pack returns a packed byte array which represents a BatteryStatus payload
func (m *BatteryStatus) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.CURRENT_CONSUMED)
	binary.Write(data, binary.LittleEndian, m.ENERGY_CONSUMED)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_1)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_2)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_3)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_4)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_5)
	binary.Write(data, binary.LittleEndian, m.VOLTAGE_CELL_6)
	binary.Write(data, binary.LittleEndian, m.CURRENT_BATTERY)
	binary.Write(data, binary.LittleEndian, m.ACCU_ID)
	binary.Write(data, binary.LittleEndian, m.BATTERY_REMAINING)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the BatteryStatus
func (m *BatteryStatus) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.CURRENT_CONSUMED)
	binary.Read(data, binary.LittleEndian, &m.ENERGY_CONSUMED)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_1)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_2)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_3)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_4)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_5)
	binary.Read(data, binary.LittleEndian, &m.VOLTAGE_CELL_6)
	binary.Read(data, binary.LittleEndian, &m.CURRENT_BATTERY)
	binary.Read(data, binary.LittleEndian, &m.ACCU_ID)
	binary.Read(data, binary.LittleEndian, &m.BATTERY_REMAINING)
}

//
// MESSAGE SETPOINT_8DOF
//
// MAVLINK_MSG_ID_SETPOINT_8DOF 148
//
// MAVLINK_MSG_ID_SETPOINT_8DOF_LEN 33
//
// MAVLINK_MSG_ID_SETPOINT_8DOF_CRC 241
//
//
type Setpoint8Dof struct {
	VAL1          float32 // Value 1
	VAL2          float32 // Value 2
	VAL3          float32 // Value 3
	VAL4          float32 // Value 4
	VAL5          float32 // Value 5
	VAL6          float32 // Value 6
	VAL7          float32 // Value 7
	VAL8          float32 // Value 8
	TARGET_SYSTEM uint8   // System ID
}

// NewSetpoint8Dof returns a new Setpoint8Dof
func NewSetpoint8Dof(VAL1 float32, VAL2 float32, VAL3 float32, VAL4 float32, VAL5 float32, VAL6 float32, VAL7 float32, VAL8 float32, TARGET_SYSTEM uint8) *Setpoint8Dof {
	m := Setpoint8Dof{}
	m.VAL1 = VAL1
	m.VAL2 = VAL2
	m.VAL3 = VAL3
	m.VAL4 = VAL4
	m.VAL5 = VAL5
	m.VAL6 = VAL6
	m.VAL7 = VAL7
	m.VAL8 = VAL8
	m.TARGET_SYSTEM = TARGET_SYSTEM
	return &m
}

// Id returns the Setpoint8Dof Message ID
func (*Setpoint8Dof) Id() uint8 {
	return 148
}

// Len returns the Setpoint8Dof Message Length
func (*Setpoint8Dof) Len() uint8 {
	return 33
}

// Crc returns the Setpoint8Dof Message CRC
func (*Setpoint8Dof) Crc() uint8 {
	return 241
}

// Pack returns a packed byte array which represents a Setpoint8Dof payload
func (m *Setpoint8Dof) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.VAL1)
	binary.Write(data, binary.LittleEndian, m.VAL2)
	binary.Write(data, binary.LittleEndian, m.VAL3)
	binary.Write(data, binary.LittleEndian, m.VAL4)
	binary.Write(data, binary.LittleEndian, m.VAL5)
	binary.Write(data, binary.LittleEndian, m.VAL6)
	binary.Write(data, binary.LittleEndian, m.VAL7)
	binary.Write(data, binary.LittleEndian, m.VAL8)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Setpoint8Dof
func (m *Setpoint8Dof) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.VAL1)
	binary.Read(data, binary.LittleEndian, &m.VAL2)
	binary.Read(data, binary.LittleEndian, &m.VAL3)
	binary.Read(data, binary.LittleEndian, &m.VAL4)
	binary.Read(data, binary.LittleEndian, &m.VAL5)
	binary.Read(data, binary.LittleEndian, &m.VAL6)
	binary.Read(data, binary.LittleEndian, &m.VAL7)
	binary.Read(data, binary.LittleEndian, &m.VAL8)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
}

//
// MESSAGE SETPOINT_6DOF
//
// MAVLINK_MSG_ID_SETPOINT_6DOF 149
//
// MAVLINK_MSG_ID_SETPOINT_6DOF_LEN 25
//
// MAVLINK_MSG_ID_SETPOINT_6DOF_CRC 15
//
//
type Setpoint6Dof struct {
	TRANS_X       float32 // Translational Component in x
	TRANS_Y       float32 // Translational Component in y
	TRANS_Z       float32 // Translational Component in z
	ROT_X         float32 // Rotational Component in x
	ROT_Y         float32 // Rotational Component in y
	ROT_Z         float32 // Rotational Component in z
	TARGET_SYSTEM uint8   // System ID
}

// NewSetpoint6Dof returns a new Setpoint6Dof
func NewSetpoint6Dof(TRANS_X float32, TRANS_Y float32, TRANS_Z float32, ROT_X float32, ROT_Y float32, ROT_Z float32, TARGET_SYSTEM uint8) *Setpoint6Dof {
	m := Setpoint6Dof{}
	m.TRANS_X = TRANS_X
	m.TRANS_Y = TRANS_Y
	m.TRANS_Z = TRANS_Z
	m.ROT_X = ROT_X
	m.ROT_Y = ROT_Y
	m.ROT_Z = ROT_Z
	m.TARGET_SYSTEM = TARGET_SYSTEM
	return &m
}

// Id returns the Setpoint6Dof Message ID
func (*Setpoint6Dof) Id() uint8 {
	return 149
}

// Len returns the Setpoint6Dof Message Length
func (*Setpoint6Dof) Len() uint8 {
	return 25
}

// Crc returns the Setpoint6Dof Message CRC
func (*Setpoint6Dof) Crc() uint8 {
	return 15
}

// Pack returns a packed byte array which represents a Setpoint6Dof payload
func (m *Setpoint6Dof) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TRANS_X)
	binary.Write(data, binary.LittleEndian, m.TRANS_Y)
	binary.Write(data, binary.LittleEndian, m.TRANS_Z)
	binary.Write(data, binary.LittleEndian, m.ROT_X)
	binary.Write(data, binary.LittleEndian, m.ROT_Y)
	binary.Write(data, binary.LittleEndian, m.ROT_Z)
	binary.Write(data, binary.LittleEndian, m.TARGET_SYSTEM)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Setpoint6Dof
func (m *Setpoint6Dof) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TRANS_X)
	binary.Read(data, binary.LittleEndian, &m.TRANS_Y)
	binary.Read(data, binary.LittleEndian, &m.TRANS_Z)
	binary.Read(data, binary.LittleEndian, &m.ROT_X)
	binary.Read(data, binary.LittleEndian, &m.ROT_Y)
	binary.Read(data, binary.LittleEndian, &m.ROT_Z)
	binary.Read(data, binary.LittleEndian, &m.TARGET_SYSTEM)
}

//
// MESSAGE MEMORY_VECT
//
// MAVLINK_MSG_ID_MEMORY_VECT 249
//
// MAVLINK_MSG_ID_MEMORY_VECT_LEN 36
//
// MAVLINK_MSG_ID_MEMORY_VECT_CRC 204
//
//
type MemoryVect struct {
	ADDRESS uint16   // Starting address of the debug variables
	VER     uint8    // Version code of the type variable. 0=unknown, type ignored and assumed int16_t. 1=as below
	TYPE    uint8    // Type code of the memory variables. for ver = 1: 0=16 x int16_t, 1=16 x uint16_t, 2=16 x Q15, 3=16 x 1Q14
	VALUE   [32]int8 // Memory contents at specified address
}

// NewMemoryVect returns a new MemoryVect
func NewMemoryVect(ADDRESS uint16, VER uint8, TYPE uint8, VALUE [32]int8) *MemoryVect {
	m := MemoryVect{}
	m.ADDRESS = ADDRESS
	m.VER = VER
	m.TYPE = TYPE
	m.VALUE = VALUE
	return &m
}

// Id returns the MemoryVect Message ID
func (*MemoryVect) Id() uint8 {
	return 249
}

// Len returns the MemoryVect Message Length
func (*MemoryVect) Len() uint8 {
	return 36
}

// Crc returns the MemoryVect Message CRC
func (*MemoryVect) Crc() uint8 {
	return 204
}

// Pack returns a packed byte array which represents a MemoryVect payload
func (m *MemoryVect) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.ADDRESS)
	binary.Write(data, binary.LittleEndian, m.VER)
	binary.Write(data, binary.LittleEndian, m.TYPE)
	binary.Write(data, binary.LittleEndian, m.VALUE)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the MemoryVect
func (m *MemoryVect) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.ADDRESS)
	binary.Read(data, binary.LittleEndian, &m.VER)
	binary.Read(data, binary.LittleEndian, &m.TYPE)
	binary.Read(data, binary.LittleEndian, &m.VALUE)
}

const (
	MAVLINK_MSG_MEMORY_VECT_FIELD_value_LEN = 32
)

//
// MESSAGE DEBUG_VECT
//
// MAVLINK_MSG_ID_DEBUG_VECT 250
//
// MAVLINK_MSG_ID_DEBUG_VECT_LEN 30
//
// MAVLINK_MSG_ID_DEBUG_VECT_CRC 49
//
//
type DebugVect struct {
	TIME_USEC uint64    // Timestamp
	X         float32   // x
	Y         float32   // y
	Z         float32   // z
	NAME      [10]uint8 // Name
}

// NewDebugVect returns a new DebugVect
func NewDebugVect(TIME_USEC uint64, X float32, Y float32, Z float32, NAME [10]uint8) *DebugVect {
	m := DebugVect{}
	m.TIME_USEC = TIME_USEC
	m.X = X
	m.Y = Y
	m.Z = Z
	m.NAME = NAME
	return &m
}

// Id returns the DebugVect Message ID
func (*DebugVect) Id() uint8 {
	return 250
}

// Len returns the DebugVect Message Length
func (*DebugVect) Len() uint8 {
	return 30
}

// Crc returns the DebugVect Message CRC
func (*DebugVect) Crc() uint8 {
	return 49
}

// Pack returns a packed byte array which represents a DebugVect payload
func (m *DebugVect) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_USEC)
	binary.Write(data, binary.LittleEndian, m.X)
	binary.Write(data, binary.LittleEndian, m.Y)
	binary.Write(data, binary.LittleEndian, m.Z)
	binary.Write(data, binary.LittleEndian, m.NAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the DebugVect
func (m *DebugVect) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_USEC)
	binary.Read(data, binary.LittleEndian, &m.X)
	binary.Read(data, binary.LittleEndian, &m.Y)
	binary.Read(data, binary.LittleEndian, &m.Z)
	binary.Read(data, binary.LittleEndian, &m.NAME)
}

const (
	MAVLINK_MSG_DEBUG_VECT_FIELD_name_LEN = 10
)

//
// MESSAGE NAMED_VALUE_FLOAT
//
// MAVLINK_MSG_ID_NAMED_VALUE_FLOAT 251
//
// MAVLINK_MSG_ID_NAMED_VALUE_FLOAT_LEN 18
//
// MAVLINK_MSG_ID_NAMED_VALUE_FLOAT_CRC 170
//
//
type NamedValueFloat struct {
	TIME_BOOT_MS uint32    // Timestamp (milliseconds since system boot)
	VALUE        float32   // Floating point value
	NAME         [10]uint8 // Name of the debug variable
}

// NewNamedValueFloat returns a new NamedValueFloat
func NewNamedValueFloat(TIME_BOOT_MS uint32, VALUE float32, NAME [10]uint8) *NamedValueFloat {
	m := NamedValueFloat{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.VALUE = VALUE
	m.NAME = NAME
	return &m
}

// Id returns the NamedValueFloat Message ID
func (*NamedValueFloat) Id() uint8 {
	return 251
}

// Len returns the NamedValueFloat Message Length
func (*NamedValueFloat) Len() uint8 {
	return 18
}

// Crc returns the NamedValueFloat Message CRC
func (*NamedValueFloat) Crc() uint8 {
	return 170
}

// Pack returns a packed byte array which represents a NamedValueFloat payload
func (m *NamedValueFloat) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.VALUE)
	binary.Write(data, binary.LittleEndian, m.NAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the NamedValueFloat
func (m *NamedValueFloat) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.VALUE)
	binary.Read(data, binary.LittleEndian, &m.NAME)
}

const (
	MAVLINK_MSG_NAMED_VALUE_FLOAT_FIELD_name_LEN = 10
)

//
// MESSAGE NAMED_VALUE_INT
//
// MAVLINK_MSG_ID_NAMED_VALUE_INT 252
//
// MAVLINK_MSG_ID_NAMED_VALUE_INT_LEN 18
//
// MAVLINK_MSG_ID_NAMED_VALUE_INT_CRC 44
//
//
type NamedValueInt struct {
	TIME_BOOT_MS uint32    // Timestamp (milliseconds since system boot)
	VALUE        int32     // Signed integer value
	NAME         [10]uint8 // Name of the debug variable
}

// NewNamedValueInt returns a new NamedValueInt
func NewNamedValueInt(TIME_BOOT_MS uint32, VALUE int32, NAME [10]uint8) *NamedValueInt {
	m := NamedValueInt{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.VALUE = VALUE
	m.NAME = NAME
	return &m
}

// Id returns the NamedValueInt Message ID
func (*NamedValueInt) Id() uint8 {
	return 252
}

// Len returns the NamedValueInt Message Length
func (*NamedValueInt) Len() uint8 {
	return 18
}

// Crc returns the NamedValueInt Message CRC
func (*NamedValueInt) Crc() uint8 {
	return 44
}

// Pack returns a packed byte array which represents a NamedValueInt payload
func (m *NamedValueInt) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.VALUE)
	binary.Write(data, binary.LittleEndian, m.NAME)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the NamedValueInt
func (m *NamedValueInt) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.VALUE)
	binary.Read(data, binary.LittleEndian, &m.NAME)
}

const (
	MAVLINK_MSG_NAMED_VALUE_INT_FIELD_name_LEN = 10
)

//
// MESSAGE STATUSTEXT
//
// MAVLINK_MSG_ID_STATUSTEXT 253
//
// MAVLINK_MSG_ID_STATUSTEXT_LEN 51
//
// MAVLINK_MSG_ID_STATUSTEXT_CRC 83
//
//
type Statustext struct {
	SEVERITY uint8     // Severity of status. Relies on the definitions within RFC-5424. See enum MAV_SEVERITY.
	TEXT     [50]uint8 // Status text message, without null termination character
}

// NewStatustext returns a new Statustext
func NewStatustext(SEVERITY uint8, TEXT [50]uint8) *Statustext {
	m := Statustext{}
	m.SEVERITY = SEVERITY
	m.TEXT = TEXT
	return &m
}

// Id returns the Statustext Message ID
func (*Statustext) Id() uint8 {
	return 253
}

// Len returns the Statustext Message Length
func (*Statustext) Len() uint8 {
	return 51
}

// Crc returns the Statustext Message CRC
func (*Statustext) Crc() uint8 {
	return 83
}

// Pack returns a packed byte array which represents a Statustext payload
func (m *Statustext) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.SEVERITY)
	binary.Write(data, binary.LittleEndian, m.TEXT)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Statustext
func (m *Statustext) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.SEVERITY)
	binary.Read(data, binary.LittleEndian, &m.TEXT)
}

const (
	MAVLINK_MSG_STATUSTEXT_FIELD_text_LEN = 50
)

//
// MESSAGE DEBUG
//
// MAVLINK_MSG_ID_DEBUG 254
//
// MAVLINK_MSG_ID_DEBUG_LEN 9
//
// MAVLINK_MSG_ID_DEBUG_CRC 46
//
//
type Debug struct {
	TIME_BOOT_MS uint32  // Timestamp (milliseconds since system boot)
	VALUE        float32 // DEBUG value
	IND          uint8   // index of debug variable
}

// NewDebug returns a new Debug
func NewDebug(TIME_BOOT_MS uint32, VALUE float32, IND uint8) *Debug {
	m := Debug{}
	m.TIME_BOOT_MS = TIME_BOOT_MS
	m.VALUE = VALUE
	m.IND = IND
	return &m
}

// Id returns the Debug Message ID
func (*Debug) Id() uint8 {
	return 254
}

// Len returns the Debug Message Length
func (*Debug) Len() uint8 {
	return 9
}

// Crc returns the Debug Message CRC
func (*Debug) Crc() uint8 {
	return 46
}

// Pack returns a packed byte array which represents a Debug payload
func (m *Debug) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.TIME_BOOT_MS)
	binary.Write(data, binary.LittleEndian, m.VALUE)
	binary.Write(data, binary.LittleEndian, m.IND)
	return data.Bytes()
}

// Decode accepts a packed byte array and populates the fields of the Debug
func (m *Debug) Decode(buf []byte) {
	data := bytes.NewBuffer(buf)
	binary.Read(data, binary.LittleEndian, &m.TIME_BOOT_MS)
	binary.Read(data, binary.LittleEndian, &m.VALUE)
	binary.Read(data, binary.LittleEndian, &m.IND)
}
