
# Model and handle the clients of the server
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
    @clients ||= {}
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
  
  def puts(text) 
    log "Client #{@id}, puts #{text}"
    Woe::Server.puts(@id, text)
  end
  
  def printf(fmt, *args)
    text = fmt.format(*args) 
    puts(text)
  end

  def raw_puts(text) 
    log "Client #{@id}, puts #{text}"
    Woe::Server.raw_puts(@id, text)
  end
  
  def raw_printf(fmt, *args)
    text = fmt.format(*args) 
    puts(text)
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

log "Mruby client script loaded OK."

