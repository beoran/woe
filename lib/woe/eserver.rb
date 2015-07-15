require 'eventmachine'
require 'tempfile'
require 'fiber'
require_relative '../monolog'



module Woe
  class Server 
    include Monolog
  
    def initialize(port =7000, logname="woe.log")
      Monolog.setup_all(logname)
      @port      = port
      # Keep an overall record of the client IDs allocated
      # and the lines of chat
      @client_id = 0
      @clients   = {}
      @tick_id   = 0
      @fiber     = nil
    end
    
    def get_free_client_id
      cli = 0
      @clients.each do |client|
        return cli if client.nil?
        cli += 1
      end
      return cli
    end
    
    def start() 
      log_info("Server listening on port #@port")
      @signature = EventMachine.start_server("0.0.0.0", @port, Client) do |client|
        client_id            = get_free_client_id
        client.id            = client_id
        client.server        = self   
        @clients[client_id]  = client
      end
      EventMachine.add_periodic_timer(1) do 
        @tick_id            += 1
        # self.broadcast("Tick tock #{@tick_id}\n")
      end  
    end
    
    
    def run
      log_info("Server main loop starts.")
      EventMachine.run do       
        self.start
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
      EventMachine.stop_server(@signature)
      EventMachine.add_timer(1) do 
        EventMachine.stop
        log_info("Server stop OK.")
      end
    end
    
   
    def broadcast(msg)
      @clients.each do |id, client|
        client.send_data(msg)
      end
    end
    

    def self.run(port=7000, logname="woe.log")    
      server = Woe::Server.new(port, logname)
      server.run
    end
      
  end
end








