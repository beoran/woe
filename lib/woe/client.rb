require 'tempfile'
require 'fiber'
require 'timeout'

require_relative '../monolog'
require_relative '../telnet'

module Woe

class Client
  include Monolog
  include Telnet::Codes
  
  attr_reader :io
  attr_reader :id
  # to allow for read timeouts
  attr_reader :timeout_at
  
  def initialize(server, id, io)
    @server = server
    @id     = id
    @io     = io
    @fiber  = Fiber.new { serve }
    @telnet = Telnet.new(self)    
    @telnet.set_support(TELNET_TELOPT_NAWS)
    @telnet.set_support(TELNET_TELOPT_MSSP)
    @telnet.set_support(TELNET_TELOPT_TTYPE)
    @telnet.set_support(TELNET_TELOPT_ECHO)
    @telnet.set_support(TELNET_TELOPT_COMPRESS2)
    @busy   = true
    # telnet event queue
    @telnet_events = []
    @timeout_at   = nil    
  end
  
  # Closes up the client
  def close
    @telnet.close
    @io.close
  end
  
  def alive?
    @fiber.alive? && @busy
  end
  
  def command(cmd, args)
    @fiber.resume(cmd, args)
  end
  
  def write_raw(data)
    @io.write(data)
  end
  
  
  def write(data)
    @telnet.send_escaped(data)
  end
  
 
  def on_start
     p "Starting client fiber"
     return nil
  end
  
  def on_write(data)
      p "On write:"
      self.write("Client #{socket}:")
      self.write(args)
  end
  
  # Telnet event class
  class TelnetEvent
    attr_accessor :type
    attr_accessor :data
    def initialize(type, data)
      @type = type
      @data = data
    end
    
    def to_s
      "<TelnetEvent #{@type} #{@data}>"    
    end
  end

  def telnet_event(type, *data)
    # store in the event queue
    @telnet_events << TelnetEvent.new(type, data)
    p "Received tenet event #{@telnet_events}."
  end
  
  def telnet_send_data(buf)
    # @telnet_events << TelnetEvent.new(:command, buf)
    p "Sending telnet data."
    self.write_raw(buf)
  end
  
  def process_telnet_events
    
  end
  
  def on_read    
    data = @io.readpartial(4096)    
    @io.flush
    @telnet.telnet_receive(data) 
    # now, the data and any telnet events are in @telnet_events
    return data
  end


  # Waits for input from the client.
  # any
  # This is always wrapped as a TelnetEvent.
  # Pure commands have the field type == :command
  # consisting of a type and a data key in a hash
  # Pass in nloops to time out the loop a loop
  # has no definite timing. 
  def wait_for_input(timeout = nil)    
    loop do
      # Timout based on number of loops. 
      if timeout
        @timeout_at = Time.now + timeout
      else
        @timeout_at = nil
      end
      
      unless @telnet_events.empty?
        return @telnet_events.shift
      end

      cmd, arg  = Fiber.yield 
      data      = nil
      case cmd 
      when :start
        on_start
      when :timeout
        @timeout_at = nil
        return nil
      when :read
        data = on_read
        # all data ends up in he telnet_events queue
        unless @telnet_events.empty?
          return @telnet_events.shift
        end
      when :write 
        on_write(arg)
      else
        p "Unknown command #{cmd}" 
      end
    end
  end
  
  
  def autohandle_event(tev)
    case tev.type
    when :naws
      @window_h, @window_w = *tev.data
      log_info("Client #{@id} window size #{@window_w}x#{@window_h}") 
    else
      log_info("Telnet event #{tev} ignored")
    end
  end
  
  def wait_for_command(timeout = nil)
    loop do
      tevent = wait_for_input(timeout)
      return nil if tevent.nil?
      if tevent.type == :data
        return tevent.data.join('').strip
      else
        autohandle_event(tevent)
      end
    end
  end
          
  
  def ask_login
    @login = nil
    while  @login.nil? || @login.empty?
      write("Login:")
      @login = wait_for_command
    end
    @login.chomp!
    true
  end

  def ask_password
    @password = nil
    while  @password.nil? || @password.empty?
      write("\r\nPassword:")
      @password = wait_for_command
    end
    @password.chomp!
    true
  end
  
  def handle_command
    order = wait_for_command
    case order
    when "/quit"
      write("Byebye!\r\n")
      @busy = false
    else
      @server.broadcast("#@login said #{order}\r\n")
    end
  end
  
  # generic negotiation
  def setup_negotiate(command, option, yes_event, no_event)
    @telnet.telnet_send_negotiate(command, option)
    tev = wait_for_input(0.5)
    return false, nil unless tev
    return false, nil if tev.type == no_event
    return false, tev unless tev.type == yes_event && tev.data[0] == option
    return true, nil
  end
  
  # Negotiate COMPRESS2 support
  def setup_compress2
    ok, tev = setup_negotiate(TELNET_WILL, TELNET_TELOPT_COMPRESS2, :do, :dont)
    return tev unless ok    
    @telnet.telnet_begin_compress2
    log_info("Client #{@id} started COMPRESS2 compression")
  end
  
  # Negotiate NAWS (window size) support
  def setup_naws  
    ok, tev = setup_negotiate(TELNET_DO, TELNET_TELOPT_NAWS, :will, :wont)
    return tev unless ok
    tev2 = wait_for_input(0.5)
    return tev2 unless tev2 && tev2.type == :naws
    @window_h, @window_w = *tev2.data
    log_info("Client #{@id} window size #{@window_w}x#{@window_h}") 
    return nil
  end
  
  
  # Negotiate MSSP (mud server status protocol) support
  def setup_mssp
    ok, tev = setup_negotiate(TELNET_WILL, TELNET_TELOPT_MSSP, :do, :dont)    
    return tev unless ok
    mssp = @server.mssp
    @telnet.telnet_send_mssp(mssp)
    return nil
  end
  
  
  
  def setup_telnet
    loop do 
      tev = wait_for_input(0.5)
      if tev
        p "setup_telnet", tev
      else
        p "no telnet setup received..."
        break
      end
    end
    setup_mssp
    setup_compress2
    setup_naws
    
    #p "mssp ev #{tev}"
    # @telnet.telnet_send_negotiate(TELNET_WILL, TELNET_TELOPT_MSSP)        
    # tev = wait_for_input(0.5)
    # p "mssp ev #{tev}"
    
    # @telnet.telnet_ttype_send
    
    
  end
    
  def serve()
    setup_telnet
    data = nil
    lok  = ask_login
    return false unless lok    
    pok  = ask_password
    return false unless pok
    write("\r\nWelcome #{@login} #{@password}!\r\n")
    while @busy do
      handle_command
    end
  end

  
=begin

  attr_accessor :id
  attr_accessor :server
  
  def initialize(server, id, socket)        
    @id         = id
    @server     = server
    @connected  = true
    @socket     = socket
    @telnet     = Telnet.new(self) 
    @busy       = true
  end
  
  
  # Get some details about the telnet connection
  def setup_telnet
    @telnet.telnet_send_negotiate(TELNET_DO, TELNET_TELOPT_TTYPE)
    @telnet.telnet_ttype_send
    type, *args = wait_for_event
    p type, args
  end
  
  def post_init()    
    send_data("Welcome!\n")
    log_info("Client #{@id} connected.")
    self.send_data("Login:")
  end
  
  # Send data to the socket
  def send_data(data)
    @socket.write(data)
  end
  
  # Run the client's main loop
  def run
    post_init
    while @connected    
      data = @socket.readpartial(4096)
      unless data.nil? || data.empty?      
        receive_data(data) 
      end
      p data
    end
  end
  

    
  def save
    self.send_data("Saving...")
    
    do_save = proc do 
      begin
        f = Tempfile.new('random')
        sleep 3
        f.write("I'm saving data.")
      ensure 
        f.close
      end
    end
    
    on_save = proc do
      self.send_data("Saved.")
    end
    
    Celluloid.defer(do_save, on_save)
  end
  
  def wait_for_event
    return Fiber.yield
  end
  
  # Basically, this method yields the fiber, and will return
  # with the input that will cme later when the fiber is resumed, normally
  # when more input becomes available from the client.
  # Any telnet commands are dispatched to the related telnet handlers.
  def wait_for_input
    loop do
      type, *args = Fiber.yield
      if type == :data
        return args.first
      else
        telnet_dispatch(type, *args)
      end      
    end
  end
  
  def try
    self.send_data("\nOK, let's try. What do you say?:")
    try = wait_for_input
    self.send_data("\nOK, nice try #{try}.\n")
  end
    
  # Fake synchronous handing of input  
  def handle_input()        
    setup_telnet
    @login    = wait_for_input
    
    self.send_data([TELNET_IAC, TELNET_WILL, TELNET_TELOPT_ECHO].pack('c*'))    
    self.send_data("\nPassword for #{@login}:")
    @password = wait_for_input
    self.send_data([TELNET_IAC, TELNET_WONT, TELNET_TELOPT_ECHO].pack('c*'))    
    self.send_data("\nOK #{@password}, switching to command mode.\n")
      
    while @connected
      line = wait_for_input
      # If the user says 'quit', disconnect them
      if line =~ /^\/quit/
        @connected = false
        close_connection_after_writing
      # Shut down the server if we hear 'shutdown'
      elsif line =~ /^\/reload/
        @server.reload
      elsif line =~ /^\/shutdown/
        @connected = false
        @server.stop
      elsif line =~ /^\/save/      
        self.save
      elsif line =~ /^\/try/      
        self.try          
      else
        @server.broadcast("Client #{id} says #{line}")
      end
    end
  end  
  

  def receive_data(data)
    # Ignore any input if already requested disconnection
    return unless @connected
    @telnet.telnet_receive(data) 
  end
  
  
  
  def unbind
    log_info("Client #{@id} has left from #{@ip}:#{@port}")
    @server.disconnect(@id)
  end

  # Called when the telnet module wants to send data.
  def telnet_send_data(buffer)
    p "Sending telnet data #{buffer}"
    self.send_data(buffer)
  end
  
  # Dispatches a telnet event to a function named telnet_(event_name)
  def telnet_dispatch(type, *args)
    meth = "telnet_#{type}".to_sym
    self.send(meth, *args)
  end
  
  
  # Telnet event handler, called on incoming events.
  def telnet_event(type, *args)
    log_info("Telnet event received by client #{id}: #{type}, #{args}")
    if @fiber
      # restart the fiber if available
      @fiber.resume(type, *args)
    else      
      # set up a fiber to handle the events
      # Like that, the handle_input can be programmed in a fake-syncronous way
      @fiber = Fiber.new do      
        handle_input
      end
      # Must resume twice becaus of the way telnet_event_fiber works
      @fiber.resume()
      @fiber.resume(type, *args)
    end    
  end
  


  
  # Real handler, called inside a fiber
  def telnet_event_fiber()
    raise "not implemented"    
  end
    
  def telnet_environment(fb, vars)
    p fb,vars
  end
  
  def telnet_environment(fb, vars)  
    p fb,vars
  end
  
  
  def telnet_mssp(vars)
    @mssp_vars = vars
  end
  
  def telnet_ttype_is(term)
    @term = term
    self.send_data("\nYou have a #{@term} type terminal.\n")
    p "term #{@term}"
  end
    
  def telnet_ttype_send(term)
    p "term #{term} sent"
  end
  
  
  def telnet_compress(compress)  
    p "compress #{compress} set"
  end  
    
    
  def telnet_subnegotiate(sb_telopt, buffer)
    p "received subnegotiate #{sb_telopt} #{buffer}"
  end  
  
  
  
  def do_main
  end
  
  
  def telnet_data(data)
=begin  
    # send data over telnet protocol. Should arrive below in telnet_data
    if @fiber
      @fiber.resume(data)
    else      
      # set up a fiber to handle the input
      # Like that, the handle_input can be programmed in a fake-syncronous way
      @fiber = Fiber.new do      
        handle_input()
      end
      # Must resume twice becaus of the way handle_input works
      @fiber.resume()
      @fiber.resume(data)
    end    
  end  
  
  def telnet_iac(byte)  
    p "received  iac #{byte}"
  end  
=end

  end # class Client

end # module Woe


