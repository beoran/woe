
require 'singleton'
require 'optparse'

VERSION = '0.1.0'

module Woe
  class Settings
    include Singleton
    
    attr_reader :data_dir
    attr_reader :port
    
    def initialize
      @data_dir = 'data'
      @port     = 7000
    end
    
    def var_dir
      File.join(@data_dir, 'var')
    end
    
    def script_dir
      File.join(@data_dir, 'script')
    end
    
    def parse_args    
      oparser = OptionParser.new do |opt|
        opt.banner = "Usage woe [options]"
        opt.on('-p', '--port [PORT]',  "Port at which the WOE server should listen.") do |  port |
          @port = port.to_i
        end 
        opt.on('-d', '--data [DIR]', "Directory where the data and scripts of the server are installed.") do | dir |
          @data_dir = dir
        end   
        opt.on('-h', '--help', "Show this message.") do
          puts opt
          exit
        end
        opt.on('-v', '--version', "Show the version of WOE.") do
          puts VERSION      
          exit
        end  
      end      
      oparser.parse!    
    end
    
    def self.parse_args
      self.instance.parse_args
    end
    
    def self.port
      self.instance.port
    end
    
    def self.data_dir
      self.instance.data_dir
    end
    
    def self.var_dir
      self.instance.var_dir
    end
    
    def self.script_dir
      self.instance.script_dir
    end    
  end  # class
end # module Woe
  
  

