
script "log.rb"

log "Main mruby script loaded OK."
p Signal.constants

script "motd.rb"
script "sitef.rb"
script "account.rb"
script "security.rb"
script "mode.rb"
script "mode/setup.rb"
script "mode/login.rb"
script "mode/normal.rb"
script "mode/character.rb"
script "client.rb"
script "timer.rb"




# Return an array of symbols of constants of klass that match value
def const_syms(klass, value)
  res = []
  klass.constants.each do |c|
     cv = klass.const_get(c)
     res << c if cv == value
  end
  return res
end


def signal_syms(value)
  return const_syms(Signal, value) 
end

def woe_on_healing_tick
  # p "healing"
end

def woe_on_motion_tick
  # p "motion"
end

def woe_on_battle_tick
  # p "battle"
end

def woe_on_weather_tick
  # p "weather"
end

def woe_on_save_tick
  # p "weather"
end



def start_timers
  @started_timers ||= false
  if @started_timers
    log "Timers already started."
  else
    log "Staring timer(s)..."
    Timer.add("healing" , 30.0      , 30.0)       { woe_on_healing_tick }
    Timer.add("motion"  ,  5.0      ,  5.0)       { woe_on_motion_tick  }
    Timer.add("battle"  ,  1.0      ,  1.0)       { woe_on_battle_tick  }
    Timer.add("weather" , 90.0      , 90.0)       { woe_on_weather_tick }
    Timer.add("save"    , 15 * 60.0 , 15 * 60.0)  { woe_on_save_tick    }
    
    #@timer_id = Woe::Server.new_timer()
    #Woe::Server.set_timer(@timer_id, 1.0, 1.0);
  end
  @started_timers = true
end

def woe_on_connect(client_id)
  p "Client #{client_id} connected"
  client = Client.add(client_id)
  client.on_start
end

def woe_on_disconnect(client_id)
  p "Client #{client_id} disconnected"
  Client.remove(client_id)
end

def woe_forward_to_client(client_id, method, *args)
  client = Client.get(client_id)
  if client
    p "Client #{client} #{client.id} ok."
    if client.respond_to?(method)
      client.send(method, *args)
    else
      log "Client cannot handle #{method}."
    end
  else
    log "Unknown client #{client_id} for #{method}."
    Woe::Server.disconnect(client_id)
  end
end

def woe_on_input(client_id, buf)
  woe_forward_to_client(client_id, :on_input, buf)
end

def woe_on_negotiate(client_id, how, option) 
  woe_forward_to_client(client_id, :on_negotiate, how, option)
end

def woe_on_subnegotiate(client_id, option, buffer) 
  woe_forward_to_client(client_id, :on_subnegotiate, option, buffer)
end

def woe_on_iac(client_id, option, command) 
  woe_forward_to_client(client_id, :on_iac, option, command)
end

def woe_on_ttype(client_id, cmd, name) 
  woe_forward_to_client(client_id, :on_ttype, cmd, name)
end

def woe_on_error(client_id, code, message) 
  woe_forward_to_client(client_id, :on_error, code, message)
end

def woe_on_warning(client_id, code, message) 
  woe_forward_to_client(client_id, :on_warning, code, message)
end

def woe_begin_compress(client_id, state) 
  woe_forward_to_client(client_id, :on_compress, state)
end


def woe_begin_zmp(client_id, size) 
  woe_forward_to_client(client_id, :on_begin_zmp, size)
end

def woe_zmp_arg(client_id, index, value)  
  woe_forward_to_client(client_id, :on_zmp_arg, index, value)
end

def woe_finish_zmp(client_id, size) 
  woe_forward_to_client(client_id, :on_finish_zmp, size)
end

def woe_begin_environ(client_id, size) 
  woe_forward_to_client(client_id, :on_begin_environ, size)
end

def woe_environ_arg(client_id, index, type, key, value)  
  woe_forward_to_client(client_id, :on_environ_arg, index, type, key, value)
end

def woe_finish_environ(client_id, size) 
  woe_forward_to_client(client_id, :on_finish_environ, size)
end


def woe_begin_mssp(client_id, size) 
  woe_forward_to_client(client_id, :on_begin_mssp, size)
end

def woe_mssp_arg(client_id, index, type, key, value)  
  woe_forward_to_client(client_id, :on_mssp_arg, index, type, key, value)
end

def woe_finish_mssp(client_id, size) 
  woe_forward_to_client(client_id, :on_finish_mssp, size)
end



def woe_on_signal(signal)
  log "Received signal #{signal} #{signal_syms(signal)} in script"
  case signal 
    when 10 # SIGUSR1
      log "Reloading main script."
      script "main.rb"
    when 28 # SIGWINCH
      # ignore this signal
    else
      Woe::Server.quit 
  end
end


def woe_on_timer(timer, value, interval)
  # log "Timer #{timer} #{value} #{interval} passed."
  Timer.on_timer(timer)
end


start_timers

=begin
f = File.open("/account/B/Beoran/Beoran.account", "r");
if f
  while (!f.eof?)
    lin = f.gets(255)
    log "Read line #{lin}"
  end
  f.close
end

Dir.mkdir("/account/C")
Dir.mkdir("/account/C/Celia")



f = File.open("/account/C/Celia/Celia.account", "w");
if f
  f.puts("name=Celia\n")
  f.puts("algo=plain\n")
  f.puts("pass=hello1woe\n")
  f.close
end

f = File.open("/account/C/Celia/Celia.account", "r");
if f
  while (!f.eof?)
    lin = f.gets(255)
    log "Read line #{lin}"
  end
  f.close
end
=end

=begin
a = Account.new(:id => 'Dyon', :pass => 'DN33Fbe/OGrM6', 
  :algo => 'crypt'
)
p a.id
a.save

d = Account.serdes_fetch('Dyon')
p d
Account.serdes_forget('Dyon')
d = Account.serdes_fetch('Dyon')
p d
=end

p crypt("noyd8pass")


