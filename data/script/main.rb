
script "log.rb"

log "Main mruby script loaded OK."
p Signal.constants

script "sitef.rb"
script "client.rb"
script "timer.rb"
script "account.rb"

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




def start_timers
  @started_timers ||= false
  if @started_timers
    log "Timers already started."
  else
    log "Staring timer(s)..."
    Timer.add("healing" , 30.0, 30.0) { woe_on_healing_tick }
    Timer.add("motion"  ,  5.0,  5.0) { woe_on_motion_tick }
    Timer.add("battle"  ,  1.0,  1.0) { woe_on_battle_tick }
    Timer.add("weather" , 90.0, 90.0) { woe_on_weather_tick }
    
    #@timer_id = Woe::Server.new_timer()
    #Woe::Server.set_timer(@timer_id, 1.0, 1.0);
  end
  @started_timers = true
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
  p "Client #{client_id} negotiating."
end

def woe_on_subnegotiate(client_id, option, buffer) 
  p "Client #{client_id} subnegotiating."
end

def woe_on_iac(client_id, option, command) 
  p "Client #{client_id} iac #{command}."
end


def woe_on_ttype(client_id, cmd, name) 
  p "Client #{client_id} ttype #{cmd} #{name}."
end

def woe_on_error(client_id, code, message) 
end

def woe_on_warning(client_id, code, message) 
end


def woe_begin_compress(client_id, state) 
end


def woe_begin_zmp(client_id, size) 
end

def woe_zmp_arg(client_id, index, value)  
end

def woe_finish_zmp(client_id, size) 
end


def woe_begin_environ(client_id, size) 
end

def woe_environ_arg(client_id, index, type, key, value)  
end

def woe_finish_environ(client_id, size) 
end


def woe_begin_mssp(client_id, size) 
end

def woe_mssp_arg(client_id, index, type, key, value)  
end

def woe_finish_mssp(client_id, size) 
end



def woe_on_signal(signal)
  log "Received signal #{signal} #{signal_syms(signal)} in script"
  case signal 
    when 10 # SIGUSR1
      log "Reloading main script."
      script "main.rb"
    when 28 
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

rec = { :name => 'Dyon', :algo => 'plain', :pass => 'noyd8pass',
:desc => "A buff, tanned man\n with blonde hair."
} 
recs = [ rec ]
Sitef.save_filename('Dyon.test', recs);

rrecs, rerrs = Sitef.load_filename('Dyon.test');
p rrecs, rerrs

a = Account.new('Dyon', 'noyd8pass')
p a.id
a.save
