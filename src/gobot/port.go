package gobot

type Port struct {
  Name string
}

func (Port) NewPort(p string) *Port{
  return new(Port)
}

func (p *Port) ToString() string {
  return p.Name
}
/*
module Artoo
  # The Artoo::Port class represents port and/or host to be used to connect
  # tp a specific individual hardware device.
  class Port
    attr_reader :port, :host

    # Create new port
    # @param [Object] data
    def initialize(data=nil)
      @is_tcp, @is_serial, @is_portless = false
      parse(data)
    end

    # @return [Boolean] True if serial port
    def is_serial?
      @is_serial == true
    end

    # @return [Boolean] True if tcp port
    def is_tcp?
      @is_tcp == true
    end

    # @return [Boolean] True if does not have real port
    def is_portless?
      @is_portless == true
    end

    # @return [String] port
    def to_s
      if is_portless?
        "none"
      elsif is_serial?
        port
      else
        "#{host}:#{port}"
      end
    end

    private

    def parse(data)
      case
      # portless
      when data.nil?
        @port = nil
        @is_portless = true

      # is TCP host/port?
      when m = /(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})/.match(data)
        @port = m[2]
        @host = m[1]
        @is_tcp = true

      # is it a numeric port for localhost tcp?
      when /^[0-9]{1,5}$/.match(data)
        @port = data
        @host = "localhost"
        @is_tcp = true

      # must be a serial port
      else
        @port = data
        @host = nil
        @is_serial = true
      end
    end
  end
end
*/