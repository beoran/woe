require 'eventmachine'
require 'tempfile'
require 'fiber'


module Woe

class Minifilter

  def initialize(give)
    @give   = give
    @fiber  = nil
  end 
  
  def wait_for_input(val=nil)
    return Fiber.yield(val)
  end
  
  def filter_input(line)
    if line =~ /2/
      return (line + line)
    end
    
    if line =~ /4/
      res = wait_for_input
      return (res * 4)      
    end
    
    if line =~ /0/
      return (nil)      
    end
    
    return line
  end
    
end

class Client < EventMachine::Connection
  attr_accessor :id
  attr_accessor :server
  
  def initialize(*args)    
    super(*args)
    @id         = nil
    @server     = nil
    @connected  = false
    @port       = nil
    @ip         = nil
    @fiber      = nil
    @account    = nil
    @filter     = ::Woe::Minifilter.new(self)
  end
  
  def post_init()
    send_data("Welcome!\n")
    pn          = self.get_peername
    @port, @ip  = Socket.unpack_sockaddr_in(pn)
    send_data("You are connecting from #{@ip}:#{@port}\n")
    @connected  = true
    self.send_data("Login:")
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
    
    EM.defer(do_save, on_save)    
  end
  
  
  # Basically, this method yields the fiber, and will return
  # with the input that will cme later when the fiber is resumed, normally
  # when more input becomes available from the client.
  # The 
  def wait_for_input
    data = Fiber.yield
    # the filters MUST be aplied here, since then it can also be 
    # fake-syncronous and use Fiber.yield to wait for additional input if 
    # needed 
    line = @filter.filter_input(data)
    return line
  end
  
  def try
    self.send_data("\nOK, let's try. What do you say?:")
    try = wait_for_input
    self.send_data("\nOK, nice try #{try}.\n")
  end
    
  # Fake synchronous handing of input  
  def handle_input()
    @login    = wait_for_input
    self.send_data("\nPassword for #{@login}:")
    @password = wait_for_input
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
    # 
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
  
  
  
  
  def unbind
    $stderr.puts("Client #{id} has left")
    @server.disconnect(@id)
    
  end
end

end


