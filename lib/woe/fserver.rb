require 'socket'
require 'fiber'
include Socket::Constants


class FServer
  def initialize
    @timers     = {}
    @reading    = []
    @writing    = []
    @clients    = {}
    @client_id  = 0
  end
  
  class Client
    attr_reader :io
    attr_reader :id
    
    def initialize(server, id, io)
      @server = server
      @id     = id
      @io     = io
      @fiber  = Fiber.new { serve } 
      @busy   = true
    end
    
    def alive?
      @fiber.alive? && @busy
    end
    
    def command(cmd, args)
      @fiber.resume(cmd, args)
    end
    
    def write(data)
      @io.write(data)
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
    
    def on_read    
      data = @io.readpartial(4096)
      p "After read: #{data}"
      @io.flush
      return data
    end

  
    def wait_for_input
      loop do
        cmd, arg  = Fiber.yield 
        data      = nil
        case cmd 
        when :start
          on_start
        when :read
          return on_read
        when :write 
          on_write(arg)
        else
          p "Unknown command #{cmd}" 
        end
      end
    end
            

    
    def ask_login
      @login = nil
      while  @login.nil? || @login.empty?
        write("Login:")
        @login = wait_for_input.chomp
      end
      true
    end

    def ask_password
      @password = nil
      while  @password.nil? || @password.empty?
        write("\r\nPassword:")
        @password = wait_for_input.chomp
      end
      true
    end
    
    def handle_command
      order = wait_for_input.chomp
      case order
      when "/quit"
        write("Byebye!\r\n")
        @busy = false
      else
        @server.broadcast("#@login said #{order}\r\n")
      end
    end
    
    
    def serve()
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
  end
  
  def start
    @server = TCPServer.new('localhost', 7000)
    @reading << @server
    serve
  end

  def add_client
    socket               = @server.accept_nonblock
    @reading          << socket
    @client_id          += 1
    client               = Client.new(self, @client_id, socket)
    @clients[socket]     = client
    client.command(:start, nil)
    puts "Client #{socket} connected"
    return client
  end

  def broadcast(message)
    @clients.each_pair do | id, client |
      client.write(message + "\r\n")
    end
  end
  
  def add_timer(id, delta = 1)
    now   = Time.now
    stop  = now + delta
    @timers[id.to_sym]  = stop
  end
  
  def serve
    add_timer(:test, 15)
    loop do
      # Kick out clients that disconnected.
      disconnected = @clients.reject { |c, cl| cl.alive? }
      disconnected.each do |c, cl|
        rsock = c
        @clients.delete(rsock)
        @reading.delete(rsock)
        rsock.close
      end
      
      # Nodify timers that expired
      now     = Time.now
      expired = @timers.select { |k, t| t < now } 
      expired.each do |k, t|
        broadcast("Timer #{k} expired: #{t}, #{now}.\n\r")
        @timers.delete(k)
      end
            
    
      readable, writable = IO.select(@reading, @writing, nil, 0.1)
      if readable
        readable.each do | rsock |
          if rsock == @server
            add_client
          # Kick out clients with broken connections.
          elsif rsock.eof?
            @clients.delete(rsock)
            @reading.delete(rsock)
            rsock.close
          else
            # Tell the client to get their read on.              
            client  = @clients[rsock]
            text    = client.command(:read, nil)
          end
        end
      end
    end
  end
  
end

server = FServer.new
server.start



