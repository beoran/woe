
# Monolog, an easy to use logger for ruby 
module Monolog 
  
  module Logger
    attr_reader :data
    
    def initialize(data = nil)
      @data       = data
    end 
    
    def log(file, line, name, level, format, *args)
    end
    
    def close
      @data.close if @data && @data.respond_to?(:close)
    end
    
  end
  
  class FileLogger 
    include Logger
    def initialize(filename)
      @data = File.open(filename, "at")
    end
    
    def log(file, line, name, format, *args)
      @data.printf("%s: %s: %s: %d: ", Time.now.to_s, name, file, line)
      @data.printf(format, *args)
      @data.printf("\n")
    end
  end
  
  class StdinLogger < FileLogger
    def initialize
       @data = $stdin
    end
  end
    
  class StderrLogger < FileLogger
    def initialize
       @data = $stderr
    end
  end
         
  class Log 
    attr_reader :loggers
    attr_reader :levels
    
    def initialize
      @loggers = []
      @levels  = {} 
    end
    
    def add_logger(logger)  
      @loggers << logger
    end
    
    def enable_level(name)
      @levels[name.to_sym] = true 
    end
    
    def disable_level(name)
      @levels[name.to_sym] = false 
    end
    
    def log_va(file, line, name, format, *args)
      level = @levels[name.to_sym]
      return nil unless level   
      @loggers.each do | logger |
        logger.log(file, line, name, format, *args)
      end
    end
    
    def close
      @loggers.each do | logger |
        logger.close()
      end
    end
  end
  
  def self.setup
    @log = Log.new
  end
  
  def self.get_log
    return @log
  end
  
  def self.setup_all(name = nil, err = true, out = false) 
    setup
    add_stderr_logger if err
    add_stdout_logger if out    
    add_file_logger(name) if name
    enable_level(:INFO)
    enable_level(:WARNING)
    enable_level(:ERROR)
    enable_level(:FATAL)
  end
  
  def self.enable_level(l)
    @log ||= nil
    return unless @log
    @log.enable_level(l)
  end

  def self.disable_level(l)      
    @log ||= nil
    return unless @log
    @log.disable_level(l)
  end
  
  def self.add_logger(l)
    @log ||= nil
    return unless @log
    @log.add_logger(l)
  end
  
  def self.add_stdin_logger
    self.add_logger(StdinLogger.new)
  end
  
  def self.add_stderr_logger
    self.add_logger(StderrLogger.new)
  end
  
  def self.add_file_logger(filename = "log.log")
    self.add_logger(FileLogger.new(filename))
  end
  
  def self.close
    @log ||= nil
    return unless @log
    @log.close
  end

  
  def self.log_va(file, line, name, format, *args)
    @log ||= nil
    return unless @log
    @log.log_va(file, line, name, format, *args)
  end
  
  def log(name, format, * args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, name, format, *args)
  end
  
  def log_error(format, *args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, :ERROR, format, *args)
  end

  def log_warning(format, *args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, :WARNING, format, *args)
  end

  def log_info(format, *args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, :INFO, format, *args)
  end

  def log_debug(format, *args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, :DEBUG, format, *args)
  end
  
  def log_fatal(format, *args)
    file, line, fun = caller.first.to_s.split(':')
    Monolog.log_va(file, line, :FATAL, format, *args)
  end
  
  alias error log_error
  alias warn log_warning
  alias info log_info
  alias error log_error
  
  extend(self)
 end
 
 
 
 
