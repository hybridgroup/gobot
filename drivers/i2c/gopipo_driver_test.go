package i2c

import (
	"testing"
	"time"

	"log"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/drivers/i2c"
	"github.com/hybridgroup/gobot/gobottest"
	"github.com/hybridgroup/gobot/platforms/gopigo"
)

// # This is the address for the GoPiGo
// address = 0x08
// debug=0
// #GoPiGo Commands
// fwd_cmd				=[119]		#Move forward with PID
// motor_fwd_cmd		=[105]		#Move forward without PID
// bwd_cmd				=[115]		#Move back with PID
// motor_bwd_cmd		=[107]		#Move back without PID
// left_cmd			=[97]		#Turn Left by turning off one motor
// left_rot_cmd		=[98]		#Rotate left by running both motors is opposite direction
// right_cmd			=[100]		#Turn Right by turning off one motor
// right_rot_cmd		=[110]		#Rotate Right by running both motors is opposite direction
// stop_cmd			=[120]		#Stop the GoPiGo
// ispd_cmd			=[116]		#Increase the speed by 10
// dspd_cmd			=[103]		#Decrease the speed by 10
// m1_cmd      		=[111]     	#Control motor1
// m2_cmd    			=[112]     	#Control motor2
// read_motor_speed_cmd=[114]		#Get motor speed back
//
// volt_cmd			=[118]		#Read the voltage of the batteries
// us_cmd				=[117]		#Read the distance from the ultrasonic sensor
// led_cmd				=[108]		#Turn On/Off the LED's
// servo_cmd			=[101]		#Rotate the servo
// enc_tgt_cmd			=[50]		#Set the encoder targeting
// fw_ver_cmd			=[20]		#Read the firmware version
// en_enc_cmd			=[51]		#Enable the encoders
// dis_enc_cmd			=[52]		#Disable the encoders
// read_enc_status_cmd	=[53]		#Read encoder status
// en_servo_cmd		=[61]		#Enable the servo's
// dis_servo_cmd		=[60]		#Disable the servo's
// set_left_speed_cmd	=[70]		#Set the speed of the right motor
// set_right_speed_cmd	=[71]		#Set the speed of the left motor
// en_com_timeout_cmd	=[80]		#Enable communication timeout
// dis_com_timeout_cmd	=[81]		#Disable communication timeout
// timeout_status_cmd	=[82]		#Read the timeout status
// enc_read_cmd		=[53]		#Read encoder values
// trim_test_cmd		=[30]		#Test the trim values
// trim_write_cmd		=[31]		#Write the trim values
// trim_read_cmd		=[32]
//
// digital_write_cmd   =[12]      	#Digital write on a port
// digital_read_cmd    =[13]      	#Digital read on a port
// analog_read_cmd     =[14]      	#Analog read on a port
// analog_write_cmd    =[15]      	#Analog read on a port
// pin_mode_cmd        =[16]      	#Set up the pin mode on a port
//
// ir_read_cmd			=[21]
// ir_recv_pin_cmd		=[22]
// cpu_speed_cmd		=[25]

type i2cGoPiGoTestAdaptor struct {
	t       *testing.T
	asserts []byte
}

func (g *i2cGoPiGoTestAdaptor) I2cWrite(address int, data []byte) error {
	gobottest.Assert(g.t, 0x08, address)

	for i, a := range g.asserts {
		gobottest.Assert(g.t, data[i], a)
	}

	return nil
}

func (g *i2cGoPiGoTestAdaptor) I2cRead(address int, len int) ([]byte, error) {
	gobottest.Assert(g.t, 0x08, address)
	return asserts, nil
}

func initTestGoPiGo(t *testing.T, asserts []byte) *GoPiGoDriver {
	return &i2cGoPiGoTestAdaptor{
		t:       t,
		asserts: asserts,
	}
}

func TestGoPiGoMotor1(t *testing.T) {
	asserts := []byte{byte(111), 1, 10, 0}
	g := initTestGoPiGo(t, asserts)
	g.Motor1(1, 10)
}

func ExampleNewGoPiGoDriver() {
	gpg, err := gopigo.NewAdaptor()
	if err != nil {
		log.Fatalf("failed to crate gopigo adaptor: %s", err)
	}
	err = gpg.PinMode(gopigo.PIN_LED_LEFT, gopigo.PIN_MODE_OUTPUT)
	if err != nil {
		log.Fatalf("failed to set output mode: %s", err)
	}

	led := gpio.NewLedDriver(gpg, gopigo.PIN_LED_LEFT)

	gpgDriver := i2c.NewGoPiGoDriver(gpg)
	work := func() {
		v, err := gpgDriver.FirmwareVersion()
		if err != nil {
			log.Errorf("failed to read firmware version: %s", err)
		}
		log.Infof("firmware version: %d", v)

		gobot.Every(1*time.Second, func() {
			err = led.Toggle()
			if err != nil {
				log.Error("failed to toggle led: %s", err)
			}
		})
	}
	dex := gobot.NewRobot("moe",
		[]gobot.Connection{gpg},
		[]gobot.Device{gpgDriver},
		work)
	err = dex.Start()
	if err != nil {
		log.Fatal(err)
	}
}
