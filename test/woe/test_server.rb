require 'atto'
include Atto::Test

require_relative '../../lib/woe/client' 
require_relative '../../lib/woe/server' 

assert { Woe::Server }

# Fork off the server so the Net::Telnet tests can procede
pid = Process.fork do 
  Woe::Server.run
  # exit here to ge
  puts __FILE__ + ' Server Done'
  exit 0
end

p pid

assert { pid }

require 'net/telnet'


assert do
  sleep 1
  client = Net::Telnet.new('Host' => 'localhost', 'Port' => 7000, 
                            'Timeout' => 3)
  res    = client.waitfor(/Login:/)
  client.write("Axl\n")
  client.waitfor(/.*/)
  client.write("pass\n")
  client.waitfor(/.*/)
  client.write("hello\n")
  client.waitfor(/.*/)
  client.write("/quit\n")
  client.waitfor(/.*/)
  client.close
  ok = !!client
  ok  
end


# Finally stop the server and wait for it to finish
Process.kill(:TERM, pid)
Process.wait(pid)
 puts __FILE__ + ' Tests Done'
 
