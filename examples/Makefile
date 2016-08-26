
ip 				:= 192.168.0.2
build_dcmotor 	:= GOARM=7 GOARCH=arm GOOS=linux go build -v raspi_adafruit_dcmotor.go
build_stepper 	:= GOARM=7 GOARCH=arm GOOS=linux go build -v raspi_adafruit_stepper.go
build_servo 	:= GOARM=7 GOARCH=arm GOOS=linux go build -v raspi_adafruit_servo.go

all:
	$(build_stepper)	
	$(build_dcmotor)	
	$(build_servo)	
	scp raspi_adafruit_stepper raspi_adafruit_dcmotor raspi_adafruit_servo pi@$(ip):~

stepper:
	$(build_stepper)	
	scp raspi_adafruit_stepper pi@$(ip):~

dc:
	$(build_dcmotor)	
	scp raspi_adafruit_dcmotor pi@$(ip):~

servo:
	$(build_servo)	
	scp raspi_adafruit_servo pi@$(ip):~
