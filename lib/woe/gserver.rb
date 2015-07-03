require 'gserver'


module Woe
class Client
  attr_reader :id
  attr_reader :io
  
  def initialize(server, id, io)
    @server = server
    @id     = id
    @io     = io
    @busy   = true
  end
  
  
  def on_input(line)
    # If the user says 'quit', disconnect them
    if line =~ /^\/quit/
      @busy = false
    # Shut down the server if we hear 'shutdown'
    elsif line =~ /^\/shutdown/
      @server.stop
    else
      @server.broadcast("Client #{id} says #{line}")
    end
  end
    
  def serve_once
      if IO.select([io], nil, nil, 0)
        # If so, retrieve the data and process it..
        line = io.gets
        on_input(line)
      else
        
      end
  end
  
  def serve
    while @busy
      serve_once
    end
  end
end


class Server < GServer
  def initialize(*args)
    super(*args)
    self.audit          = true
    # increase the connection limit
    @maxConnections     = 400 
    # Keep an overall record of the client IDs allocated
    # and the lines of chat
    @client_id = 0
    @clients   = []
  end
  
  
  
  def serve(io)
    # Increment the client ID so each client gets a unique ID
    @client_id += 1
    client      = Client.new(self, @client_id, io)
    @clients << client 
    client.io.puts("Welcome, client nr #{client.id}!")
    client.serve
  end
  
  def broadcast(msg)
    p msg
    @clients.each do |client|
      client.io.puts(msg)
    end
  end
  

  def self.run(port=7000)
    server = Woe::Server.new(port)
    server.start
    server.join
  end
    
end



end

Woe::Server.run



