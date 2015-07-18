require 'socket'
require 'fiber'
require 'timeout'
require 'time'

require_relative '../monolog'



module Woe
class Server
  include Socket::Constants
  include Monolog
  
  def initialize(port = 7000, host='0.0.0.0')
    @port       = port
    @host       = host
    @timers     = {}
    @reading    = []
    @writing    = []
    @clients    = {}
    @client_id  = 0
    # Used for MSSP 
    @mssp       = {
      "NAME"        => "Workers Of Eruta",
      "UPTIME"      => Time.now.to_i.to_s,
      "PLAYERS"     => "0",
      "CRAWL DELAY" => "0",
      "CODEBASE"    => "WOE",
      "CONTACT"     => "beoran@gmail.com",
      "CREATED"     => "2015",
       "ICON"       => "None",
      "LANGUAGE"    => "English",
      "LOCATION"    => "USA",
      "MINIMUM AGE" => "18",
      "WEBSITE"     => "beoran.net",
      "FAMILY"      => "Custom",
      "GENRE"       => "Science Fiction",
      "GAMEPLAY"    => "Adventure",
      "STATUS"      => "Alpha",
      "GAMESYSTEM"  => "Custom",
      "INTERMUD"    => "",
      "SUBGENRE"    => "None",
      "AREAS"       => "0",
      "HELPFILES"   => "0",
      "MOBILES"     => "0",
      "OBJECTS"     => "0",
      "ROOMS"       => "1",
      "CLASSES"     => "0",
      "LEVELS"      => "0",
      "RACES"       => "3",
      "SKILLS"      => "900",
      "ANSI"        => "1",
      "MCCP"        => "1",
      "MCP"         => "0",
      "MSDP"        => "0",
      "MSP"         => "0",
      "MXP"         => "0",
      "PUEBLO"      => "0",
      "UTF-8"       => "1",
      "VT100"       => "1",
      "XTERM 255 COLORS" => "1",
      "PAY TO PLAY"      => "0",
      "PAY FOR PERKS"    => "0",
      "HIRING BUILDERS"  => "0",
      "HIRING CODERS"    => "0" 
    }    
  end
  
  
  # Returns the MSSP data
  def mssp
    @mssp["PLAYERS"] = @clients.size.to_s
    return @mssp
  end


  
  # Look for a free numeric client ID in the client hash.
  def get_free_client_id
    cli = 0
    @clients.each do |client|
      return cli if client.nil?
      cli += 1
    end
    return cli
  end
  
  def run
    log_info("Starting server on #{@host} #{@port}.")
    @server = TCPServer.new(@host, @port)
    @reading << @server
    serve
  end

  def add_client
    socket               = @server.accept_nonblock
    @reading          << socket
    client_id            = get_free_client_id 
    client               = Client.new(self, client_id, socket)
    @clients[client_id]  = client
    client.command(:start, nil)
    puts "Client #{client.id} connected on #{socket}."
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
  
  
  # Nodify timers that expired
  def handle_disconnected
   # Kick out clients that disconnected.
    disconnected = @clients.reject { |id, cl| cl.alive? }
    disconnected.each do |id, cl|
      rsock = cl.io
      @clients.delete(id)
      @reading.delete(rsock)
      cl.close
    end     
  end
  
  # Handle timers
  def handle_timers
    now     = Time.now
    expired = @timers.select { |k, t| t < now } 
    expired.each do |k, t|
      broadcast("Timer #{k} expired: #{t}, #{now}.\n\r")
      @timers.delete(k)
    end
  end
  
  def client_for_socket(sock)
    @clients.each do |k, c| 
      return c if c.io == sock 
    end
    return nil
  end
  
   # Notify clients that have a read timeout set
  def handle_timeouts
    now = Time.now
    @clients.each  do |id, cl| 
      if cl.timeout?
        cl.command(:timeout, nil)
      end
    end
  end
 
  def serve
    @busy = true
    while @busy
      handle_disconnected
      handle_timers
      handle_timeouts
            
      readable, writable = IO.select(@reading, @writing, nil, 0.1)
      if readable
        readable.each do | rsock |
          if rsock == @server
            add_client
          # Kick out clients with broken connections.
          elsif rsock.eof?
            cli = client_for_socket(rsock)
            @clients.delete(cli.id)
            @reading.delete(rsock)
            cli.close
          else
            # Tell the client to get their read on.              
            client  = client_for_socket(rsock)
            if client
              text    = client.command(:read, nil)
            end
          end
        end
      end
    end
  end  
  
  def self.run(port = 7000, host = '0.0.0.0', logname = 'woe.log')
    Monolog.setup_all(logname)
    server = self.new(port, host)
    server.run
    Monolog.close
  end
  
end # class Server
end # module Woe


