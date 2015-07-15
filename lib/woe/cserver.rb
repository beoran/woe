require 'tempfile'
require 'fiber'
require_relative '../monolog'

require 'celluloid/io'
require 'celluloid/autostart'


Celluloid.task_class = Celluloid::Task::Threaded

module Woe
  class Server     
    include Monolog
    
    include Celluloid::IO
    finalizer :shutdown
 
    def initialize(host = 'localhost', port =7000, logname="woe.log")
      Monolog.setup_all(logname)
      # Celluloid.logger = self
      @port      = port
      # Keep an overall record of the client IDs allocated
      # and the lines of chat
      @client_id = 0
      @clients   = {}
      @tick_id   = 0
      @host      = host
      p "Server listening on #@host port #@port"
      @server    = Celluloid::IO::TCPServer.new(@host, @port)
      async.run
    end
    
    def get_free_client_id
      cli = 0
      @clients.each do |client|
        return cli if client.nil?
        cli += 1
      end
      return cli
    end
    
    def run
      @busy = true
      p "Server main loop starts."
      while @busy
        begin
          p "Accepting"
          socket = @server.accept
          p socket
          async.handle_connection(socket)  
        rescue 
          p "exception #{$!}"
          @busy = false
        end
      end
    end    
    
    def handle_connection(socket)
      p "Connecting socket."
       _, port, host = socket.peeraddr
      p "*** Received connection from #{host}:#{port}"

      client_id            = get_free_client_id
      client               = Client.new(self, client_id, socket)   
      @clients[client_id]  = client
      begin
        client.run        
      rescue EOFError
        p "*** #{host}:#{port} disconnected"
      ensure
        disconnect(client.id)
        socket.close
      end
    end  
    
    
    def disconnect(id)
      log_info("Server disconnecting client #{@id}")
      @clients.delete(id)
    end
    
    def clients_stopped?
    end
    
    def reload
      log_info("Server reload")
      broadcast("Server reload NOW!\n")
      begin 
        load 'lib/telnet.rb'
        load 'lib/woe/client.rb'
        load 'lib/woe/server.rb'
        broadcast("Server reloaded OK.\n")
      rescue Exception => ex
        bt = ex.backtrace.join("\n")
        log_error("Server reload failed: #{ex}: #{bt}")
        broadcast("Server reload exception #{ex}: #{bt}!\n")
      end
    end
    
    def stop
      log_info("Server stop")
      shutdown
      log_info("Server stop OK.")      
    end
    
   
    def broadcast(msg)
      @clients.each do |id, client|
        client.send_data(msg)
      end
    end
    
   
    def shutdown
      log_info("Shuting down server.")
      @busy = false
      @server.close if @server
    end    
  end
end










