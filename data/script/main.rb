
p "Hi from main.rb"
p global_variables
# p :$--TEST--
# p $"--TEST
log "hello to log"

class Client 
  attr_reader :id
  attr_reader :buffer
  attr_reader :in_lines
  
  def initialize(id)
    @id       = id
    @buffer   = ""
    @in_lines = []
  end
  
  def self.add(client_id)
    @clients2 ||= {}
    @clients[client_id] = Client.new(client_id)
  end
  
  def self.get(client_id)
    @clients ||= {}
    return @clients[client_id]
  end
  
  def self.remove(client_id)    
    @clients[client_id] = nil
  end
  
  def send(text) 
  log "Client #{@id}, send #{text}"
    Woe::Server.send_to_client(@id, text)
  end
  
  def on_input(str)
    @buffer ||= ""
    @buffer << str
    if @buffer["\r\n"]
      command, rest = @buffer.split("\r\n", 1)
      log "Client #{@id}, command #{command}"
      if (command.strip == "quit") 
        self.send("Bye bye!")
        Woe::Server.disconnect(@id)
        Client.remove(@id)
        return nil
      else 
        self.send("I don't know how to #{command}")
      end
      @buffer = rest
    end    
  end  
end



def woe_on_connect(client_id)
  p "Client #{client_id} connected"
  Client.add(client_id)
end

def woe_on_disconnect(client_id)
  p "Client #{client_id} disconnected"
  Client.remove(client_id)
end

def woe_on_input(client_id, buf)
  p "Client #{client_id} input #{buf}"
  client = Client.get(client_id)
  unless client
    log "Unknown client #{client_id} in woe_on_input."
    Woe::Server.disconnect(client_id)
  else
    p "Client #{client} #{client.id} ok."
    client.on_input(buf)
  end  
end

def woe_on_negotiate(client_id, how, option) 
  p "Client #{client} #{client.id} negotiating."
end

def woe_on_subnegotiate(client_id, option, buffer) 
  p "Client #{client} #{client.id} negotiating."
end


def woe_on_signal(signal)
  log "Received signal in script #{signal}"
  if signal !=28 
    Woe::Server.quit 
  end
end
