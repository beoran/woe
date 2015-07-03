
# Model and handle the clients of the server
class Client 
  include Telnet
  
  attr_reader   :id
  attr_reader   :buffer
  attr_reader   :in_lines
  attr_accessor :account
  attr_accessor :character
  attr_accessor :mode
  
  def initialize(id)
    @id       = id
    @buffer   = ""
    @in_lines = []
    @account  = nil
    @character= nil
    @mode     = Mode::Setup.new(self)
  end
  
  def self.add(client_id)
    @clients ||= {}
    @clients[client_id] = Client.new(client_id)
    return @clients[client_id]
  end
  
  def self.get(client_id)
    @clients ||= {}
    return @clients[client_id]
  end
  
  def self.remove(client_id)    
    @clients[client_id] = nil
  end
  
  def send_to_client(text)
    log "Client #{@id}, send_to_client #{text}"
    Woe::Server.send_to_client(@id, text)
  end
  
  def puts(text) 
    log "Client #{@id}, puts #{text}"
    Woe::Server.puts(@id, text)
  end
  
  def printf(fmt, *args)
    if args && !args.empty?
      text = fmt.format(*args) 
    else
      text = fmt
    end
    puts(text)
  end

  def raw_puts(text) 
    log "Client #{@id}, puts #{text}"
    Woe::Server.raw_puts(@id, text)
  end
  
  def raw_printf(fmt, *args)    
    if args
      text = fmt.format(*args) 
    else
      text = fmt
    end
    puts(text)
  end
  
  def negotiate(how, what)
    log "Client #{@id} negotiate #{how} #{what}"
    Woe::Server.negotiate(@id, how, what)
  end
  
  def password_mode
    self.negotiate(TELNET_WILL, TELNET_TELOPT_ECHO)
  end
  
  def normal_mode
    self.negotiate(TELNET_WONT, TELNET_TELOPT_ECHO)
  end
  
  
  def on_negotiate(how, opt)
    if (opt == TELNET_TELOPT_COMPRESS2) 
      if (how == TELNET_DO)
        log "Beginning compress2 mode"
        Woe::Server.begin_compress2(@id)
      end
    elsif (opt == TELNET_TELOPT_TELOPT_TTYPE) 
      if (how == TELNET_WILL) || (opt == TELNET_DO)
        
      end
    end
  end
  
  
  
  def ask_type   
    self.negotiate(TELNET_DO, TELNET_TELOPT_TTYPE
  end
  
  
  def on_start
    ask_ttype
    @mode.on_start
  end
  
  def on_input(str)
    @buffer ||= ""
    @buffer << str
    if @buffer["\r\n"]
      command, rest = @buffer.split("\r\n", 1)
      command.chomp!
      log "Client #{@id}, command #{command}"
      if (command.strip == "!quit") 
        self.send_to_client("Bye bye!")
        Woe::Server.disconnect(@id)
        Client.remove(@id)
        return nil
      elsif (command.strip == "!load") 
        log "Reloading main script."
        self.puts("Reloading main script.")
        script "main.rb"
      else 
        @mode.do_command(command);
      end
      command = nil
      @buffer = rest
    end    
  end  
end

log "Mruby client script loaded OK."

