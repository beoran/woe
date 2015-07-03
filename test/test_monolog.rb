require 'atto'
include Atto::Test

require_relative '../lib/monolog' 

LOG_NAME = '/tmp/monolog_test.log'

assert { Monolog } 
assert { Monolog.setup }
assert { Monolog.setup_all(LOG_NAME) }

assert { Monolog.get_log }

assert do
  lg = Monolog.get_log.loggers 
  !(lg.empty?)
end


assert do
  Monolog.log_info("bazz") 
  sleep 1
  res = File.read(LOG_NAME)
  res =~ /bazz/ && res =~ /INFO/
end

assert "Debug evel is not logged by default after setup_all" do
  Monolog.log_debug("frotz")
  sleep 1
  res = File.read(LOG_NAME)
  res !~ /frotz/
end


assert { Monolog.close }

