require 'eventmachine'
require 'tempfile'
require 'fiber'


module Woe
  class Server 
    def initialize(port =7000)
      @port      = port
      # Keep an overall record of the client IDs allocated
      # and the lines of chat
      @client_id = 0
      @clients   = {}
      @tick_id   = 0
      @fiber     = nil
    end
    
    def start()    
      @signature = EventMachine.start_server("0.0.0.0", @port, Client) do |client|
        @client_id          += 1
        client.id            = @client_id
        client.server        = self   
        @clients[@client_id] = client
      end
      EventMachine.add_periodic_timer(1) do 
        @tick_id            += 1
        # self.broadcast("Tick tock #{@tick_id}\n")
      end  
    end
    
    def run
      EventMachine.run do       
        self.start
      end  
    end
    
    
    def disconnect(id)
      @clients.delete(id)
    end
    
    def clients_stopped?    
    end
    
    def reload
      broadcast("Reloading\n")
      begin 
        load 'lib/woe/server.rb'
        broadcast("Reloaded\n")
      rescue Exception => ex
        broadcast("Exception #{ex}: #{ex.backtrace.join("\n")}!\n")
      end
    end
    
    def stop
      EventMachine.stop_server(@signature)
      EventMachine.add_timer(1) { EventMachine.stop }
    end
    
   
    def broadcast(msg)
      @clients.each do |id, client|
        client.send_data(msg)
      end
    end
    

    def self.run(port=7000)    
      server = Woe::Server.new(port)
      server.run
    end
      
  end
end








